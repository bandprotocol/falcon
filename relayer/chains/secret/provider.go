package secret

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
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
	secretwallet "github.com/bandprotocol/falcon/relayer/wallet/secret"
)

var _ chains.ChainProvider = (*SecretChainProvider)(nil)

// SecretChainProvider handles interactions with the Secret Network chain.
type SecretChainProvider struct {
	Config *SecretChainProviderConfig

	ChainName string
	Client    Client
	Log       logger.Logger

	Wallet      wallet.Wallet
	DB          db.Database
	Alert       alert.Alert
	FreeSigners chan wallet.Signer
}

func NewSecretChainProvider(
	chainName string,
	client Client,
	cfg *SecretChainProviderConfig,
	log logger.Logger,
	wallet wallet.Wallet,
	alert alert.Alert,
) *SecretChainProvider {
	return &SecretChainProvider{
		Config:      cfg,
		ChainName:   chainName,
		Client:      client,
		Log:         log.With("chain_name", chainName),
		Wallet:      wallet,
		Alert:       alert,
		FreeSigners: chains.LoadSigners(wallet),
	}
}

func (cp *SecretChainProvider) Init(ctx context.Context) error {
	if err := cp.Client.Connect(ctx); err != nil {
		return err
	}

	go cp.Client.StartLivelinessCheck(ctx, cp.Config.LivelinessCheckingInterval)
	return nil
}

func (cp *SecretChainProvider) SetDatabase(database db.Database) {
	cp.DB = database
}

func (cp *SecretChainProvider) QueryTunnelInfo(
	_ context.Context,
	tunnelID uint64,
	tunnelDestinationAddr string,
) (*types.Tunnel, error) {
	tunnel := types.NewTunnel(tunnelID, tunnelDestinationAddr, true, nil, nil)
	return tunnel, nil
}

func countValidSecretSymbolRates(signalPrices []bandtypes.SignalPrice) uint64 {
	// Mirrors fkms/cosmwasm_secret expectations: parse "CS:BTC-USD" and only count "USD" quote with price > 0.
	var count uint64
	for _, sp := range signalPrices {
		if sp.Price == 0 {
			continue
		}
		parts := strings.Split(sp.SignalID, ":")
		if len(parts) != 2 {
			continue
		}
		right := strings.Split(parts[1], "-")
		if len(right) != 2 {
			continue
		}
		if right[1] != "USD" {
			continue
		}
		count++
	}
	return count
}

// RelayPacket relays the given Band packet to Secret Network by:
// - building a SignSecretRequest payload for fkms (including per-tunnel contract address)
// - delegating tx signing & encryption to fkms via remote signer
// - broadcasting the returned signed tx blob
func (cp *SecretChainProvider) RelayPacket(ctx context.Context, packet *bandtypes.Packet) error {
	if err := cp.Client.CheckAndConnect(ctx); err != nil {
		return err
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

		signing, err := chains.SelectSigning(packet)
		if err != nil {
			log.Error("Select signing error", "retry_count", retryCount, err)
			lastErr = err
			continue
		}

		accountNumber, sequence, err := cp.Client.GetAccount(ctx, freeSigner.GetAddress())
		if err != nil {
			log.Error("Get account error", "retry_count", retryCount, err)
			lastErr = err
			continue
		}

		// fkms signs using gas_limit + gas_prices, and the fee impact depends on the number of symbol rates.
		validSymbolRates := countValidSecretSymbolRates(packet.SignalPrices)
		gasLimit := cp.Config.GasLimitBase + cp.Config.GasLimitEach*validSymbolRates

		signerPayload := secretwallet.NewSignerPayload(
			freeSigner.GetAddress(),
			packet.TargetAddress,
			cp.Config.CosmosChainID,
			accountNumber,
			sequence,
			gasLimit,
			cp.Config.GasPrice,
			"Relayer",
			cp.Config.CodeHash,
			cp.Config.ChainPubkey,
		)

		payloadBytes, err := json.Marshal(signerPayload)
		if err != nil {
			log.Error("Marshal signer payload error", "retry_count", retryCount, err)
			lastErr = err
			continue
		}

		rAddress, signature := chains.ExtractEVMSignature(signing.EVMSignature)
		tssPayload := wallet.NewTssPayload(signing.Message, rAddress, signature)

		txBlob, err := freeSigner.Sign(payloadBytes, tssPayload)
		if err != nil {
			log.Error("Sign transaction error", "retry_count", retryCount, err)
			lastErr = err
			continue
		}

		var balance *big.Int
		if cp.DB != nil {
			balance, err = cp.Client.GetBalance(ctx, freeSigner.GetAddress())
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

		txHash, err := cp.Client.BroadcastTx(txBlob)
		if err != nil {
			log.Error("Broadcast transaction error", "retry_count", retryCount, err)
			lastErr = err
			continue
		}

		createdAt := time.Now()
		log.Info("Submitted a message; checking transaction status", "tx_hash", txHash, "retry_count", retryCount)

		// save pending tx in db
		if cp.DB != nil {
			tx := cp.prepareTransaction(ctx, txHash, freeSigner.GetAddress(), packet, nil, balance, log, retryCount)
			chains.HandleSaveTransaction(cp.DB, cp.Alert, tx, log)
		}

		txResult := cp.WaitForConfirmedTx(ctx, txHash, log)
		cp.handleMetrics(packet.TunnelID, createdAt, txResult)

		if cp.DB != nil {
			tx := cp.prepareTransaction(ctx, txHash, freeSigner.GetAddress(), packet, &txResult, balance, log, retryCount)
			chains.HandleSaveTransaction(cp.DB, cp.Alert, tx, log)
		}

		if txResult.Status == types.TX_STATUS_SUCCESS {
			alert.HandleReset(cp.Alert, alert.NewTopic(alert.RelayTxErrorMsg).
				WithTunnelID(packet.TunnelID).
				WithChainName(cp.ChainName))
			return nil
		}

		lastErr = fmt.Errorf("transaction failed: %s", txResult.FailureReason)
	}

	alert.HandleAlert(
		cp.Alert,
		alert.NewTopic(alert.RelayTxErrorMsg).WithTunnelID(packet.TunnelID).WithChainName(cp.ChainName),
		lastErr.Error(),
	)
	return fmt.Errorf("failed to relay packet after %d attempts", cp.Config.MaxRetry)
}

func (cp *SecretChainProvider) QueryBalance(ctx context.Context, keyName string) (*big.Int, error) {
	signer, ok := cp.Wallet.GetSigner(keyName)
	if !ok {
		cp.Log.Error("Key name does not exist", "key_name", keyName)
		return nil, fmt.Errorf("key name does not exist: %s", keyName)
	}
	return cp.Client.GetBalance(ctx, signer.GetAddress())
}

func (cp *SecretChainProvider) GetChainName() string       { return cp.ChainName }
func (cp *SecretChainProvider) ChainType() types.ChainType { return types.ChainTypeSecret }
func (cp *SecretChainProvider) GetWallet() wallet.Wallet   { return cp.Wallet }

func (cp *SecretChainProvider) prepareTransaction(
	ctx context.Context,
	txHash string,
	signerAddress string,
	packet *bandtypes.Packet,
	txResult *TxResult,
	oldBalance *big.Int,
	log logger.Logger,
	retryCount int,
) *db.Transaction {
	if txHash == "" {
		return nil
	}

	var signalPrices []db.SignalPrice
	for _, p := range packet.SignalPrices {
		signalPrices = append(signalPrices, *db.NewSignalPrice(p.SignalID, p.Price))
	}

	txStatus := types.TX_STATUS_PENDING
	gasUsed := decimal.NullDecimal{}
	effectiveGasPrice := decimal.NullDecimal{}
	balanceDelta := decimal.NullDecimal{}

	var blockTimestamp *time.Time
	if txResult != nil {
		txStatus = txResult.Status
		gasUsed = txResult.GasUsed
		effectiveGasPrice = txResult.EffectiveGasPrice

		if (txResult.Status == types.TX_STATUS_SUCCESS || txResult.Status == types.TX_STATUS_FAILED) && txResult.BlockHeight != nil {
			block, err := cp.Client.GetBlockByHeight(ctx, txResult.BlockHeight)
			if err != nil {
				log.Error("Failed to get block by height", "retry_count", retryCount, "block_height", txResult.BlockHeight, err)
				alert.HandleAlert(cp.Alert, alert.NewTopic(alert.GetHeaderBlockErrorMsg).
					WithTunnelID(packet.TunnelID).
					WithChainName(cp.ChainName), err.Error())
			} else {
				timestamp := block.Time.UTC()
				blockTimestamp = &timestamp
				alert.HandleReset(cp.Alert, alert.NewTopic(alert.GetHeaderBlockErrorMsg).
					WithTunnelID(packet.TunnelID).
					WithChainName(cp.ChainName))
			}
		}

		if oldBalance != nil {
			newBalance, err := cp.Client.GetBalance(ctx, signerAddress)
			if err != nil {
				log.Error("Failed to get balance", "retry_count", retryCount, err)
				alert.HandleAlert(cp.Alert, alert.NewTopic(alert.GetBalanceErrorMsg).
					WithTunnelID(packet.TunnelID).
					WithChainName(cp.ChainName), err.Error())
			} else {
				diff := new(big.Int).Sub(newBalance, oldBalance)
				balanceDelta = decimal.NewNullDecimal(decimal.NewFromBigInt(diff, 0))
				alert.HandleReset(cp.Alert, alert.NewTopic(alert.GetBalanceErrorMsg).
					WithTunnelID(packet.TunnelID).
					WithChainName(cp.ChainName))
			}
		}
	}

	tx := db.NewTransaction(
		txHash,
		packet.TunnelID,
		packet.Sequence,
		cp.ChainName,
		types.ChainTypeSecret,
		signerAddress,
		txStatus,
		gasUsed,
		effectiveGasPrice,
		balanceDelta,
		signalPrices,
		blockTimestamp,
	)

	return tx
}

func (cp *SecretChainProvider) CheckConfirmedTx(
	ctx context.Context,
	txHash string,
) (TxResult, error) {
	tx, err := cp.Client.GetTx(ctx, txHash)
	if err != nil {
		err = fmt.Errorf("failed to get transaction: %w", err)
		return NewTxResult(
			types.TX_STATUS_PENDING,
			decimal.NullDecimal{},
			decimal.NullDecimal{},
			nil,
			err.Error(),
		), err
	}

	gasUsed := decimal.NewNullDecimal(decimal.New(tx.GasUsed, 0))
	fee := decimal.NewNullDecimal(decimal.NewFromInt(tx.Fee))
	var effectiveGasPrice decimal.NullDecimal
	if gasUsed.Decimal.IsZero() {
		effectiveGasPrice = decimal.NullDecimal{}
	} else {
		effectiveGasPrice = decimal.NewNullDecimal(fee.Decimal.Div(gasUsed.Decimal))
	}
	blockHeight := big.NewInt(tx.Height)

	if tx.StatusCode == 0 {
		return NewTxResult(types.TX_STATUS_SUCCESS, gasUsed, effectiveGasPrice, blockHeight, ""), nil
	}
	return NewTxResult(
		types.TX_STATUS_FAILED,
		gasUsed,
		effectiveGasPrice,
		blockHeight,
		fmt.Sprintf("transaction failed with failure message %s", tx.Log),
	), nil
}

func (cp *SecretChainProvider) WaitForConfirmedTx(
	ctx context.Context,
	txHash string,
	log logger.Logger,
) TxResult {
	createdAt := time.Now()
	var lastErr error
	for time.Since(createdAt) <= cp.Config.WaitingTxDuration {
		result, err := cp.CheckConfirmedTx(ctx, txHash)
		if err != nil {
			lastErr = err
			log.Debug("Failed to check tx status", "tx_hash", txHash, "err", err)
		}

		switch result.Status {
		case types.TX_STATUS_SUCCESS, types.TX_STATUS_FAILED:
			return result
		case types.TX_STATUS_PENDING:
			time.Sleep(cp.Config.CheckingTxInterval)
		}
	}

	failureReason := fmt.Sprintf("timed out waiting %s for tx %s to be confirmed", cp.Config.WaitingTxDuration, txHash)
	if lastErr != nil {
		failureReason = fmt.Sprintf("%s: %v", failureReason, lastErr)
	}

	return NewTxResult(types.TX_STATUS_TIMEOUT, decimal.NullDecimal{}, decimal.NullDecimal{}, nil, failureReason)
}

func (cp *SecretChainProvider) handleMetrics(tunnelID uint64, createdAt time.Time, txResult TxResult) {
	relayermetrics.IncTxsCount(tunnelID, cp.ChainName, types.ChainTypeSecret.String(), txResult.Status.String())

	switch txResult.Status {
	case types.TX_STATUS_SUCCESS, types.TX_STATUS_FAILED:
		relayermetrics.ObserveTxProcessTime(
			tunnelID,
			cp.ChainName,
			types.ChainTypeSecret.String(),
			txResult.Status.String(),
			time.Since(createdAt).Milliseconds(),
		)
		relayermetrics.ObserveGasUsed(
			tunnelID,
			cp.ChainName,
			types.ChainTypeSecret.String(),
			txResult.Status.String(),
			txResult.GasUsed.Decimal.InexactFloat64(),
		)
	}
}
