package flow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"text/template"
	"time"

	"github.com/onflow/flow-go-sdk"
	"github.com/shopspring/decimal"

	"github.com/bandprotocol/falcon/internal/relayermetrics"
	"github.com/bandprotocol/falcon/relayer/alert"
	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/db"
	"github.com/bandprotocol/falcon/relayer/logger"
	"github.com/bandprotocol/falcon/relayer/wallet"
	flowwallet "github.com/bandprotocol/falcon/relayer/wallet/flow"
)

const (
	flowFeeEvent = "FlowFees.FeesDeducted"
	flowToWeiExp = 10
)

var _ chains.ChainProvider = (*FlowChainProvider)(nil)

// FlowChainProvider handles interactions with the Flow blockchain.
type FlowChainProvider struct {
	Config    *FlowChainProviderConfig
	ChainName string

	Client Client

	Log logger.Logger

	DB db.Database

	Alert alert.Alert

	FreeSigners    chan wallet.Signer
	Wallet         wallet.Wallet
	scriptTemplate *template.Template
}

// NewFlowChainProvider creates a new Flow chain provider.
func NewFlowChainProvider(
	chainName string,
	client Client,
	cfg *FlowChainProviderConfig,
	log logger.Logger,
	w wallet.Wallet,
	a alert.Alert,
) (*FlowChainProvider, error) {
	tmpl, err := template.New("relayer_rates").Parse(RelayScript)
	if err != nil {
		return nil, err
	}

	return &FlowChainProvider{
		Config:         cfg,
		ChainName:      chainName,
		Client:         client,
		Log:            log.With("chain_name", chainName),
		Alert:          a,
		FreeSigners:    chains.LoadSigners(w),
		Wallet:         w,
		scriptTemplate: tmpl,
	}, nil
}

// Init connects to the Flow chain.
func (cp *FlowChainProvider) Init(ctx context.Context) error {
	if err := cp.Client.Connect(ctx); err != nil {
		return err
	}

	go cp.Client.StartLivelinessCheck(ctx, cp.Config.LivelinessCheckingInterval)

	return nil
}

// SetDatabase assigns the given database instance.
func (cp *FlowChainProvider) SetDatabase(database db.Database) {
	cp.DB = database
}

// QueryTunnelInfo returns an active Flow tunnel
// Flow does not require on-chain sequence tracking. If a database is configured,
// the latest successful sequence for the given tunnel is read from the DB and
// attached to the returned Tunnel; otherwise no sequence is provided.
func (cp *FlowChainProvider) QueryTunnelInfo(
	_ context.Context,
	tunnelID uint64,
	tunnelDestinationAddr string,
) (*types.Tunnel, error) {
	var seq *uint64
	// Keep the latest successful sequence in the database for reference
	// because it cannot be tracked on-chain.
	if cp.DB != nil {
		latestSequence, err := cp.DB.GetLatestSuccessSequence(tunnelID)
		if err != nil {
			return nil, fmt.Errorf("[FlowProvider] failed to get latest success sequence: %w", err)
		}
		seq = &latestSequence
	}
	// If db is not set, seq will be nil, which means relayer doesn't know the latest successful sequence.
	tunnel := types.NewTunnel(tunnelID, tunnelDestinationAddr, true, seq, nil)
	return tunnel, nil
}

// RelayPacket relays the packet to the Flow chain.
func (cp *FlowChainProvider) RelayPacket(ctx context.Context, packet *bandtypes.Packet) error {
	if err := cp.Client.CheckAndConnect(ctx); err != nil {
		cp.Log.Error("Connect client error", err)
		return fmt.Errorf("[FlowProvider] failed to connect client: %w", err)
	}

	// Get a free signer matching the target address.
	freeSigner := <-cp.FreeSigners
	defer func() { cp.FreeSigners <- freeSigner }()

	log := cp.Log.With(
		"tunnel_id", packet.TunnelID,
		"sequence", packet.Sequence,
		"signer_address", freeSigner.GetAddress(),
	)

	var script bytes.Buffer
	err := cp.scriptTemplate.Execute(&script, struct{ Contract string }{Contract: packet.TargetAddress})
	if err != nil {
		cp.Log.Error("Execute script template error", err)
		alert.HandleAlert(
			cp.Alert,
			alert.NewTopic(alert.ExecuteScriptErrorMsg).WithTunnelID(packet.TunnelID).WithChainName(cp.ChainName),
			err.Error(),
		)
		return fmt.Errorf("[FlowProvider] failed to execute script template: %w", err)
	}
	scriptBytes := script.Bytes()

	var lastErr error
	for retryCount := 1; retryCount <= cp.Config.MaxRetry; retryCount++ {
		log.Info("Relaying a message", "retry_count", retryCount)

		blockID, err := cp.Client.GetLatestBlockID(ctx)
		if err != nil {
			log.Error("Get latest block ID error", "retry_count", retryCount, err)
			lastErr = err
			continue
		}

		acc, err := cp.Client.GetAccount(ctx, freeSigner.GetAddress())
		if err != nil {
			log.Error("Get account error", "retry_count", retryCount, err)
			lastErr = err
			continue
		}

		if len(acc.Keys) == 0 {
			lastErr = fmt.Errorf("account %s has no keys", freeSigner.GetAddress())
			log.Error("Account has no keys", "retry_count", retryCount, lastErr)
			continue
		}

		keyIndex := uint32(acc.Keys[0].Index)
		sequence := acc.Keys[0].SequenceNumber

		signerPayload := flowwallet.NewSignerPayload(
			freeSigner.GetAddress(),
			cp.Config.ComputeLimit,
			blockID,
			keyIndex,
			sequence,
			scriptBytes,
		)

		payloadBytes, err := json.Marshal(signerPayload)
		if err != nil {
			log.Error("Marshal signer payload error", "retry_count", retryCount, err)
			lastErr = err
			continue
		}

		signing, err := chains.SelectSigning(packet)
		if err != nil {
			log.Error("Select signing error", "retry_count", retryCount, err)
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

		txHash, err := cp.Client.BroadcastTx(ctx, txBlob)
		if err != nil {
			log.Error("Broadcast transaction error", "retry_count", retryCount, err)
			lastErr = err
			continue
		}

		createdAt := time.Now()

		log.Info(
			"Submitted a message; checking transaction status",
			"tx_hash", txHash,
			"retry_count", retryCount,
		)

		// save pending tx in db
		cp.handleSaveTransaction(
			ctx,
			txHash,
			types.TX_STATUS_PENDING,
			freeSigner.GetAddress(),
			packet,
			balance,
			nil,
			log,
			retryCount,
		)

		// Poll for transaction confirmation.
		txStatus, fee := cp.waitForTx(ctx, txHash, log)

		cp.handleMetrics(packet.TunnelID, createdAt, txStatus)
		cp.handleSaveTransaction(
			ctx,
			txHash,
			txStatus,
			freeSigner.GetAddress(),
			packet,
			balance,
			fee,
			log,
			retryCount,
		)

		if txStatus == types.TX_STATUS_SUCCESS {
			alert.HandleReset(
				cp.Alert,
				alert.NewTopic(alert.RelayTxErrorMsg).WithTunnelID(packet.TunnelID).WithChainName(cp.ChainName),
			)

			return nil
		}

		lastErr = fmt.Errorf("transaction %s ended with status %s", txHash, txStatus)
	}

	alert.HandleAlert(
		cp.Alert,
		alert.NewTopic(alert.RelayTxErrorMsg).WithTunnelID(packet.TunnelID).WithChainName(cp.ChainName),
		lastErr.Error(),
	)

	return fmt.Errorf("[FlowProvider] failed to relay packet after %d attempts", cp.Config.MaxRetry)
}

// QueryBalance queries the FLOW balance for the given address.
func (cp *FlowChainProvider) QueryBalance(ctx context.Context, address string) (*big.Int, error) {
	return cp.Client.GetBalance(ctx, address)
}

// GetChainName retrieves the chain name from the chain provider.
func (cp *FlowChainProvider) GetChainName() string { return cp.ChainName }

// ChainType retrieves the chain type from the chain provider.
func (cp *FlowChainProvider) ChainType() types.ChainType {
	return types.ChainTypeFlow
}

// GetWallet retrieves the wallet from the chain provider.
func (cp *FlowChainProvider) GetWallet() wallet.Wallet {
	return cp.Wallet
}

// waitForTx polls until the transaction is sealed, failed, or times out.
// Returns the final TxStatus and the fee (if extractable).
func (cp *FlowChainProvider) waitForTx(
	ctx context.Context,
	txHash string,
	log logger.Logger,
) (types.TxStatus, *uint64) {
	deadline := time.Now().Add(cp.Config.WaitingTxDuration)
	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return types.TX_STATUS_TIMEOUT, nil
		default:
		}

		result, err := cp.Client.GetTxResult(ctx, txHash)
		if err != nil {
			log.Debug("GetTxResult error, will retry", "tx_hash", txHash, err)
			time.Sleep(cp.Config.CheckingTxInterval)
			continue
		}

		if result.Error != nil {
			log.Error("Transaction failed on-chain", "tx_hash", txHash, result.Error)
			return types.TX_STATUS_FAILED, extractFee(result)
		}

		if result.Status == flow.TransactionStatusSealed {
			fee := extractFee(result)
			return types.TX_STATUS_SUCCESS, fee
		}

		time.Sleep(cp.Config.CheckingTxInterval)
	}

	return types.TX_STATUS_TIMEOUT, nil
}

// extractFee extracts the fee amount from FlowFees.FeesDeducted event.
func extractFee(result *flow.TransactionResult) *uint64 {
	for _, e := range result.Events {
		if e.Value.EventType.QualifiedIdentifier == flowFeeEvent {
			for i, f := range e.Value.GetFields() {
				if f.Identifier == "amount" {
					fee, ok := e.Value.GetFieldValues()[i].ToGoValue().(uint64)
					if ok {
						return &fee
					}
				}
			}
		}
	}

	return nil
}

// handleMetrics increments tx count and, for success/failed, records processing time (ms).
func (cp *FlowChainProvider) handleMetrics(tunnelID uint64, createdAt time.Time, txStatus types.TxStatus) {
	// increment the transactions count metric
	relayermetrics.IncTxsCount(tunnelID, cp.ChainName, cp.ChainType().String(), txStatus.String())

	switch txStatus {
	case types.TX_STATUS_SUCCESS, types.TX_STATUS_FAILED:
		// track transaction processing time (ms)
		relayermetrics.ObserveTxProcessTime(
			tunnelID,
			cp.ChainName,
			cp.ChainType().String(),
			txStatus.String(),
			time.Since(createdAt).Milliseconds(),
		)
	}
}

// handleSaveTransaction saves the transaction result to the database.
func (cp *FlowChainProvider) handleSaveTransaction(
	ctx context.Context,
	txHash string,
	txStatus types.TxStatus,
	signerAddress string,
	packet *bandtypes.Packet,
	oldBalance *big.Int,
	fee *uint64,
	log logger.Logger,
	retryCount int,
) {
	if cp.DB == nil || txHash == "" {
		return
	}

	var signalPrices []db.SignalPrice
	for _, p := range packet.SignalPrices {
		signalPrices = append(signalPrices, *db.NewSignalPrice(p.SignalID, p.Price))
	}

	var blockTimestamp *time.Time
	var err error

	feeDecimal := decimal.NullDecimal{}
	balanceDelta := decimal.NullDecimal{}

	if txStatus == types.TX_STATUS_SUCCESS || txStatus == types.TX_STATUS_FAILED {
		blockTimestamp, err = cp.Client.GetBlockTimestamp(ctx, txHash)
		if err != nil {
			log.Warn("Failed to get block timestamp", "tx_hash", txHash, err)
			alert.HandleAlert(cp.Alert, alert.NewTopic(alert.GetHeaderBlockErrorMsg).
				WithTunnelID(packet.TunnelID).
				WithChainName(cp.ChainName), err.Error())
		} else {
			alert.HandleReset(cp.Alert, alert.NewTopic(alert.GetHeaderBlockErrorMsg).
				WithTunnelID(packet.TunnelID).
				WithChainName(cp.ChainName))
		}
		weiScale := decimal.New(1, flowToWeiExp)
		if fee != nil {
			feeDecimal = decimal.NewNullDecimal(decimal.NewFromInt(int64(*fee)).Mul(weiScale))
		}
		if oldBalance != nil {
			newBalance, err := cp.Client.GetBalance(ctx, signerAddress)
			if err != nil {
				log.Error("Failed to get balance after tx", "retry_count", retryCount, err)
				alert.HandleAlert(cp.Alert, alert.NewTopic(alert.GetBalanceErrorMsg).
					WithTunnelID(packet.TunnelID).
					WithChainName(cp.ChainName), err.Error())
			} else {
				diff := new(big.Int).Sub(newBalance, oldBalance)
				balanceDelta = decimal.NewNullDecimal(decimal.NewFromBigInt(diff, 0).Mul(weiScale))
				alert.HandleReset(cp.Alert, alert.NewTopic(alert.GetBalanceErrorMsg).
					WithTunnelID(packet.TunnelID).
					WithChainName(cp.ChainName))
			}
		}
	}

	packetTimestamp := time.Unix(packet.CreatedAt, 0).UTC()
	tx := db.NewTransaction(
		txHash,
		packet.TunnelID,
		packet.Sequence,
		cp.ChainName,
		types.ChainTypeFlow,
		signerAddress,
		txStatus,
		decimal.NewNullDecimal(decimal.NewFromInt(1)), // Flow doesn't expose gas units
		feeDecimal,
		balanceDelta,
		signalPrices,
		blockTimestamp,
		&packetTimestamp,
	)

	chains.HandleSaveTransaction(cp.DB, cp.Alert, tx, log)
}
