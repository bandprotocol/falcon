package soroban

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/shopspring/decimal"

	"github.com/bandprotocol/falcon/internal/relayermetrics"
	"github.com/bandprotocol/falcon/relayer/alert"
	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/db"
	"github.com/bandprotocol/falcon/relayer/logger"
	"github.com/bandprotocol/falcon/relayer/wallet"
	"github.com/bandprotocol/falcon/relayer/wallet/soroban"
)

const (
	sorobanToWeiExp       = 11 // 100 stroops (1 stroop = 0.0000001 xlm)
	DefaultFee      int64 = 100
)

var _ chains.ChainProvider = (*SorobanChainProvider)(nil)

// TxResult holds the outcome of a relayed transaction.
type TxResult struct {
	Status        types.TxStatus
	TxHash        string
	LedgerIndex   uint64
	Fee           decimal.NullDecimal
	FailureReason string
}

// NewTxResult creates a TxResult with the given fields.
func NewTxResult(
	status types.TxStatus,
	txHash string,
	ledgerIndex uint64,
	fee decimal.NullDecimal,
	failureReason string,
) TxResult {
	return TxResult{
		Status:        status,
		TxHash:        txHash,
		LedgerIndex:   ledgerIndex,
		Fee:           fee,
		FailureReason: failureReason,
	}
}

// SorobanChainProvider handles interactions with Soroban.
type SorobanChainProvider struct {
	Config    *SorobanChainProviderConfig
	ChainName string

	Client Client
	Log    logger.Logger
	DB     db.Database
	Alert  alert.Alert

	FreeSigners chan wallet.Signer
	Wallet      wallet.Wallet
}

// NewSorobanChainProvider creates a new Soroban chain provider.
func NewSorobanChainProvider(
	chainName string,
	client Client,
	cfg *SorobanChainProviderConfig,
	log logger.Logger,
	w wallet.Wallet,
	a alert.Alert,
) *SorobanChainProvider {
	return &SorobanChainProvider{
		Config:      cfg,
		ChainName:   chainName,
		Client:      client,
		Log:         log.With("chain_name", chainName),
		Alert:       a,
		FreeSigners: chains.LoadSigners(w),
		Wallet:      w,
	}
}

func (cp *SorobanChainProvider) Init(ctx context.Context) error {
	if err := cp.Client.Connect(ctx); err != nil {
		return err
	}
	go cp.Client.StartLivelinessCheck(ctx, cp.Config.LivelinessCheckingInterval)
	return nil
}

func (cp *SorobanChainProvider) SetDatabase(database db.Database) {
	cp.DB = database
}

func (cp *SorobanChainProvider) QueryTunnelInfo(
	_ context.Context,
	tunnelID uint64,
	tunnelDestinationAddr string,
) (*types.Tunnel, error) {
	// Soroban uses Skippable tunnels without sequence tracking similar to XRPL
	tunnel := types.NewTunnel(tunnelID, tunnelDestinationAddr, true, nil, nil)
	return tunnel, nil
}

func (cp *SorobanChainProvider) RelayPacket(ctx context.Context, packet *bandtypes.Packet) error {
	if err := cp.Client.CheckAndConnect(ctx); err != nil {
		return fmt.Errorf("[SorobanProvider] failed to connect client: %w", err)
	}

	// get a free signer
	cp.Log.Debug("Waiting for a free signer...")
	freeSigner := <-cp.FreeSigners
	defer func() { cp.FreeSigners <- freeSigner }()

	log := cp.Log.With(
		"tunnel_id", packet.TunnelID,
		"sequence", packet.Sequence,
		"signer_address", freeSigner.GetAddress(),
	)

	var lastErr error
	for retryCount := 1; retryCount <= cp.Config.MaxRetry; retryCount++ {
		log.Info("Relaying a message", "retry_count", retryCount)

		sequence, err := cp.Client.GetAccountSequenceNumber(freeSigner.GetAddress())
		if err != nil {
			log.Error("Get account sequence number error", "retry_count", retryCount, err)
			lastErr = err
			time.Sleep(2 * time.Second)
			continue
		}
		sequence++

		signing, err := chains.SelectSigning(packet)
		if err != nil {
			log.Error("Select signing error", "retry_count", retryCount, err)
			lastErr = err
			continue
		}

		// Currently we use default fee of 100 stroops (1 stroop = 0.0000001 xlm) for all transactions.
		// The precise network fee rates will be defined in the upcoming mainnet release. The Futurenet
		// preview release uses placeholder values.
		// https://soroban.stellar.org/docs/fundamentals-and-concepts/fees-and-metering
		//
		// TODO: Update fee model upon mainnet release
		fee := DefaultFee

		feeStats, err := cp.Client.GetFeeStats()
		if err == nil {
			fee = max(fee, feeStats.LastLedgerBaseFee)
		}

		// Add a 10% buffer to the fee to help ensure the transaction is processed.
		// The user will pay the lesser of this value or the network's required fee.
		fee = fee * 110 / 100

		signerPayload := soroban.NewSignerPayload(
			freeSigner.GetAddress(),
			packet.TargetAddress,
			fmt.Sprintf("%d", fee),
			sequence,
			cp.Config.NetworkPassphrase,
			cp.Config.Endpoints,
		)

		payloadBytes, err := json.Marshal(signerPayload)
		if err != nil {
			log.Error("Marshal signer payload error", "retry_count", retryCount, err)
			lastErr = err
			continue
		}

		rAddress, signature := chains.ExtractEVMSignature(signing.EVMSignature)
		tssPayload := wallet.NewTssPayload(
			signing.Message,
			rAddress,
			signature,
		)

		result, err := freeSigner.Sign(payloadBytes, tssPayload)
		if err != nil {
			log.Error("Sign transaction error", "retry_count", retryCount, err)
			lastErr = err
			continue
		}
		txBlob := string(result)

		var balance *big.Int
		if cp.DB != nil {
			balance, err = cp.Client.GetBalance(freeSigner.GetAddress())
			if err != nil {
				log.Error("Failed to get balance", "retry_count", retryCount, err)
				alert.HandleAlert(cp.Alert, alert.NewTopic(alert.GetBalanceErrorMsg).
					WithTunnelID(packet.TunnelID).
					WithChainName(cp.ChainName), err.Error())
			} else {
				alert.HandleReset(cp.Alert, alert.NewTopic(alert.GetBalanceErrorMsg).
					WithTunnelID(packet.TunnelID).
					WithChainName(cp.ChainName))
			}
		}

		broadcastResult, err := cp.Client.BroadcastTx(txBlob)
		if err != nil {
			log.Error("Broadcast transaction error", "retry_count", retryCount, err)
			lastErr = err
			time.Sleep(2 * time.Second)
			continue
		}

		createdAt := time.Now()

		log.Info(
			"Submitted a message; checking transaction status",
			"tx_hash", broadcastResult.TxHash,
			"retry_count", retryCount,
		)

		// save pending tx in db
		if cp.DB != nil {
			pending := NewTxResult(types.TX_STATUS_PENDING, broadcastResult.TxHash, 0, decimal.NullDecimal{}, "")
			cp.handleSaveTransaction(pending, freeSigner.GetAddress(), packet, balance, log, retryCount)
		}

		txResult := cp.WaitForConfirmedTx(broadcastResult.TxHash, log)

		cp.handleMetrics(packet.TunnelID, createdAt, txResult)
		cp.handleSaveTransaction(txResult, freeSigner.GetAddress(), packet, balance, log, retryCount)

		if txResult.Status == types.TX_STATUS_SUCCESS {
			log.Info(
				"Packet is successfully relayed",
				"tx_hash", txResult.TxHash,
				"retry_count", retryCount,
			)
			alert.HandleReset(
				cp.Alert,
				alert.NewTopic(alert.RelayTxErrorMsg).WithTunnelID(packet.TunnelID).WithChainName(cp.ChainName),
			)
			return nil
		}

		lastErr = fmt.Errorf("%s", txResult.FailureReason)
		log.Error(
			"Failed to relay packet",
			"status", txResult.Status.String(),
			"tx_hash", txResult.TxHash,
			"retry_count", retryCount,
			lastErr,
		)

		time.Sleep(2 * time.Second)
	}

	alert.HandleAlert(
		cp.Alert,
		alert.NewTopic(alert.RelayTxErrorMsg).WithTunnelID(packet.TunnelID).WithChainName(cp.ChainName),
		lastErr.Error(),
	)
	return fmt.Errorf("[SorobanProvider] failed to relay packet after %d attempts", cp.Config.MaxRetry)
}

// CheckConfirmedTx checks whether the submitted tx is confirmed on-chain.
// Returns PENDING if not yet found, SUCCESS if confirmed, FAILED if tx failed.
func (cp *SorobanChainProvider) CheckConfirmedTx(txHash string) (TxResult, error) {
	tx, err := cp.Client.GetTransactionStatus(txHash)
	if err != nil {
		return NewTxResult(types.TX_STATUS_PENDING, txHash, 0, decimal.NullDecimal{}, err.Error()), err
	}

	fee := decimal.NewNullDecimal(decimal.NewFromInt(tx.FeeCharged))

	if tx.Successful {
		return NewTxResult(types.TX_STATUS_SUCCESS, tx.Hash, uint64(tx.Ledger), fee, ""), nil
	}
	return NewTxResult(types.TX_STATUS_FAILED, tx.Hash, uint64(tx.Ledger), fee, "transaction failed on-chain"), nil
}

// WaitForConfirmedTx polls CheckConfirmedTx until the transaction is confirmed or
// WaitingTxDuration elapses. It never returns an error — on timeout the result has
// Status == TX_STATUS_TIMEOUT and FailureReason set accordingly.
func (cp *SorobanChainProvider) WaitForConfirmedTx(txHash string, log logger.Logger) TxResult {
	createdAt := time.Now()
	var lastErr error
	for time.Since(createdAt) <= cp.Config.WaitingTxDuration {
		result, err := cp.CheckConfirmedTx(txHash)
		if err != nil {
			lastErr = err
			log.Debug("Failed to check tx status", "tx_hash", txHash, err)
		}

		switch result.Status {
		case types.TX_STATUS_SUCCESS, types.TX_STATUS_FAILED:
			return result
		case types.TX_STATUS_PENDING:
			log.Debug("Waiting for tx to be confirmed", "tx_hash", txHash)
			time.Sleep(cp.Config.CheckingTxInterval)
		}
	}

	failureReason := fmt.Sprintf("timed out waiting %s for tx %s", cp.Config.WaitingTxDuration, txHash)
	if lastErr != nil {
		failureReason = fmt.Sprintf("%s: %v", failureReason, lastErr)
	}
	return NewTxResult(types.TX_STATUS_TIMEOUT, txHash, 0, decimal.NullDecimal{}, failureReason)
}

// handleMetrics increments tx count and records processing time for confirmed txs.
func (cp *SorobanChainProvider) handleMetrics(tunnelID uint64, createdAt time.Time, txResult TxResult) {
	relayermetrics.IncTxsCount(tunnelID, cp.ChainName, types.ChainTypeSoroban.String(), txResult.Status.String())

	switch txResult.Status {
	case types.TX_STATUS_SUCCESS, types.TX_STATUS_FAILED:
		relayermetrics.ObserveTxProcessTime(
			tunnelID,
			cp.ChainName,
			types.ChainTypeSoroban.String(),
			txResult.Status.String(),
			time.Since(createdAt).Milliseconds(),
		)
	}
}

func (cp *SorobanChainProvider) QueryBalance(ctx context.Context, address string) (*big.Int, error) {
	return cp.Client.GetBalance(address)
}

func (cp *SorobanChainProvider) GetChainName() string       { return cp.ChainName }
func (cp *SorobanChainProvider) ChainType() types.ChainType { return types.ChainTypeSoroban }
func (cp *SorobanChainProvider) GetWallet() wallet.Wallet   { return cp.Wallet }

func (cp *SorobanChainProvider) handleSaveTransaction(
	txResult TxResult,
	signerAddress string,
	packet *bandtypes.Packet,
	oldBalance *big.Int,
	log logger.Logger,
	retryCount int,
) {
	if cp.DB == nil || txResult.TxHash == "" {
		return
	}

	var signalPrices []db.SignalPrice
	for _, p := range packet.SignalPrices {
		signalPrices = append(signalPrices, *db.NewSignalPrice(p.SignalID, p.Price))
	}

	balanceDelta := decimal.NullDecimal{}
	gasUsed := decimal.NullDecimal{}
	fee := decimal.NullDecimal{}
	if oldBalance != nil && (txResult.Status == types.TX_STATUS_SUCCESS || txResult.Status == types.TX_STATUS_FAILED) {
		gasUsed = decimal.NewNullDecimal(decimal.NewFromInt(1))
		newBalance, err := cp.Client.GetBalance(signerAddress)
		if txResult.Fee.Valid {
			fee = decimal.NewNullDecimal(txResult.Fee.Decimal.Shift(sorobanToWeiExp))
		}
		if err != nil {
			log.Error("Failed to get balance", "retry_count", retryCount, err)
			alert.HandleAlert(cp.Alert, alert.NewTopic(alert.GetBalanceErrorMsg).
				WithTunnelID(packet.TunnelID).
				WithChainName(cp.ChainName), err.Error())
		} else {
			diff := new(big.Int).Sub(newBalance, oldBalance)
			balanceDelta = decimal.NewNullDecimal(decimal.NewFromBigInt(diff, sorobanToWeiExp))
		}
	}

	var closeTime *time.Time
	if txResult.LedgerIndex != 0 {
		ledgerCloseTime, err := cp.Client.GetLedgerCloseTime(txResult.LedgerIndex)
		if err != nil {
			log.Error(
				"Failed to get ledger close time",
				"retry_count",
				retryCount,
				"ledger_index",
				txResult.LedgerIndex,
				err,
			)
			alert.HandleAlert(cp.Alert, alert.NewTopic(alert.GetLedgerCloseTimeErrorMsg).
				WithTunnelID(packet.TunnelID).
				WithChainName(cp.ChainName), err.Error())
		} else {
			closeTime = ledgerCloseTime
			alert.HandleReset(cp.Alert, alert.NewTopic(alert.GetLedgerCloseTimeErrorMsg).
				WithTunnelID(packet.TunnelID).
				WithChainName(cp.ChainName))
		}
	}

	packetTimestamp := time.Unix(packet.CreatedAt, 0).UTC()

	tx := db.NewTransaction(
		txResult.TxHash,
		packet.TunnelID,
		packet.Sequence,
		cp.ChainName,
		types.ChainTypeSoroban,
		signerAddress,
		txResult.Status,
		gasUsed,
		fee,
		balanceDelta,
		signalPrices,
		closeTime,
		&packetTimestamp,
	)

	chains.HandleSaveTransaction(cp.DB, cp.Alert, tx, log)
}
