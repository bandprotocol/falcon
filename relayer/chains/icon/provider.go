package icon

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	v3 "github.com/icon-project/goloop/server/v3"
	"github.com/shopspring/decimal"

	"github.com/bandprotocol/falcon/internal/relayermetrics"
	"github.com/bandprotocol/falcon/relayer/alert"
	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/db"
	"github.com/bandprotocol/falcon/relayer/logger"
	"github.com/bandprotocol/falcon/relayer/wallet"
	iconwallet "github.com/bandprotocol/falcon/relayer/wallet/icon"
)

var _ chains.ChainProvider = (*IconChainProvider)(nil)

// IconChainProvider handles interactions with Icon.
type IconChainProvider struct {
	Config    *IconChainProviderConfig
	ChainName string

	Client Client

	Log logger.Logger

	Wallet wallet.Wallet
	DB     db.Database

	Alert alert.Alert

	FreeSigners chan wallet.Signer
}

// NewIconChainProvider creates a new Icon chain provider.
func NewIconChainProvider(
	chainName string,
	client Client,
	cfg *IconChainProviderConfig,
	log logger.Logger,
	wallet wallet.Wallet,
	alert alert.Alert,
) *IconChainProvider {
	return &IconChainProvider{
		Config:      cfg,
		ChainName:   chainName,
		Client:      client,
		Log:         log.With("chain_name", chainName),
		Wallet:      wallet,
		Alert:       alert,
		FreeSigners: chains.LoadSigners(wallet),
	}
}

// Init connects to the Icon chain.
func (cp *IconChainProvider) Init(ctx context.Context) error {
	if err := cp.Client.Connect(ctx); err != nil {
		return err
	}

	go cp.Client.StartLivelinessCheck(ctx, cp.Config.LivelinessCheckingInterval)

	return nil
}

// SetDatabase assigns the given database instance.
func (cp *IconChainProvider) SetDatabase(database db.Database) {
	cp.DB = database
}

// QueryTunnelInfo returns a best-effort tunnel info for Icon.
func (cp *IconChainProvider) QueryTunnelInfo(
	ctx context.Context,
	tunnelID uint64,
	tunnelDestinationAddr string,
) (*types.Tunnel, error) {
	tunnel := types.NewTunnel(tunnelID, tunnelDestinationAddr, true, nil, nil)
	return tunnel, nil
}

// RelayPacket relays the packet to Icon.
func (cp *IconChainProvider) RelayPacket(ctx context.Context, packet *bandtypes.Packet) error {
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

		signerPayload := iconwallet.NewSignerPayload(
			freeSigner.GetAddress(),
			packet.TargetAddress,
			cp.Config.StepLimit,
			cp.Config.NetworkID,
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

		res, err := freeSigner.Sign(payloadBytes, tssPayload)
		if err != nil {
			log.Error("Sign transaction error", "retry_count", retryCount, err)
			lastErr = err
			continue
		}

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

		txParams := &v3.TransactionParam{}
		if err := json.Unmarshal(res, txParams); err != nil {
			log.Error("Unmarshal transaction result error", "retry_count", retryCount, err)
			lastErr = err
			continue
		}

		txHash, err := cp.Client.BroadcastTx(*txParams)
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
		if cp.DB != nil {
			tx := cp.prepareTransaction(ctx, txHash, freeSigner.GetAddress(), packet, nil, balance, log, retryCount)
			chains.HandleSaveTransaction(cp.DB, cp.Alert, tx, log)
		}

		txResult := cp.WaitForConfirmedTx(ctx, txHash, log)

		cp.handleMetrics(packet.TunnelID, createdAt, txResult)

		if cp.DB != nil {
			tx := cp.prepareTransaction(
				ctx,
				txHash,
				freeSigner.GetAddress(),
				packet,
				&txResult,
				balance,
				log,
				retryCount,
			)
			chains.HandleSaveTransaction(cp.DB, cp.Alert, tx, log)
		}

		relayermetrics.IncTxsCount(
			packet.TunnelID,
			cp.ChainName,
			types.ChainTypeIcon.String(),
			txResult.Status.String(),
		)

		if txResult.Status == types.TX_STATUS_SUCCESS {
			log.Info(
				"Packet is successfully relayed",
				"tx_hash", txHash,
				"retry_count", retryCount,
			)

			alert.HandleReset(
				cp.Alert,
				alert.NewTopic(alert.RelayTxErrorMsg).WithTunnelID(packet.TunnelID).WithChainName(cp.ChainName),
			)

			return nil
		} else {
			log.Error("Transaction failed", "tx_hash", txHash, "failure_reason", txResult.FailureReason)
			lastErr = fmt.Errorf("transaction failed: %s", txResult.FailureReason)
			continue
		}
	}

	alert.HandleAlert(
		cp.Alert,
		alert.NewTopic(alert.RelayTxErrorMsg).WithTunnelID(packet.TunnelID).WithChainName(cp.ChainName),
		lastErr.Error(),
	)
	return fmt.Errorf("failed to relay packet after %d attempts", cp.Config.MaxRetry)
}

// QueryBalance queries balance by given key name from the destination chain.
func (cp *IconChainProvider) QueryBalance(ctx context.Context, keyName string) (*big.Int, error) {
	signer, ok := cp.Wallet.GetSigner(keyName)
	if !ok {
		cp.Log.Error("Key name does not exist", "key_name", keyName)
		return nil, fmt.Errorf("key name does not exist: %s", keyName)
	}

	return cp.Client.GetBalance(signer.GetAddress())
}

// GetChainName retrieves the chain name from the chain provider.
func (cp *IconChainProvider) GetChainName() string { return cp.ChainName }

// ChainType retrieves the chain type from the chain provider.
func (cp *IconChainProvider) ChainType() types.ChainType {
	return types.ChainTypeIcon
}

// GetWallet retrieves the wallet from the chain provider.
func (cp *IconChainProvider) GetWallet() wallet.Wallet {
	return cp.Wallet
}

// prepareTransaction prepares the transaction to be stored in the database.
func (cp *IconChainProvider) prepareTransaction(
	_ context.Context,
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

		if txResult.Status == types.TX_STATUS_SUCCESS || txResult.Status == types.TX_STATUS_FAILED {
			if txResult.BlockHeight != nil {
				block, err := cp.Client.GetBlockByHeight(txResult.BlockHeight)
				if err != nil {
					log.Error(
						"Failed to get block by height",
						"retry_count",
						retryCount,
						"block_height",
						txResult.BlockHeight,
						err,
					)
					alert.HandleAlert(cp.Alert, alert.NewTopic(alert.GetHeaderBlockErrorMsg).
						WithTunnelID(packet.TunnelID).
						WithChainName(cp.ChainName), err.Error())
				} else {
					// Block.Timestamp is in microseconds.
					timestamp := time.Unix(0, block.Timestamp*int64(time.Microsecond)).UTC()
					blockTimestamp = &timestamp
					alert.HandleReset(cp.Alert, alert.NewTopic(alert.GetHeaderBlockErrorMsg).
						WithTunnelID(packet.TunnelID).
						WithChainName(cp.ChainName))
				}
			}
		}

		// Compute new balance
		if oldBalance != nil {
			newBalance, err := cp.Client.GetBalance(signerAddress)
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

	packetTimestamp := time.Unix(packet.CreatedAt, 0).UTC()

	tx := db.NewTransaction(
		txHash,
		packet.TunnelID,
		packet.Sequence,
		cp.ChainName,
		types.ChainTypeIcon,
		signerAddress,
		txStatus,
		gasUsed,
		effectiveGasPrice,
		balanceDelta,
		signalPrices,
		blockTimestamp,
		&packetTimestamp,
	)

	return tx
}

// CheckConfirmedTx checks the confirmed transaction status.
func (cp *IconChainProvider) CheckConfirmedTx(
	ctx context.Context,
	txHash string,
) (TxResult, error) {
	tx, err := cp.Client.GetTx(txHash)
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

	// calculate gas used and effective gas price
	gasUsed := decimal.NewNullDecimal(decimal.NewFromInt(tx.StepUsed.Value()))
	effectiveGasPrice := decimal.NewNullDecimal(decimal.NewFromInt(tx.StepPrice.Value()))
	blockHeight, err := tx.BlockHeight.BigInt()
	if err != nil {
		err = fmt.Errorf("failed to parse block height: %w", err)
		return NewTxResult(
			types.TX_STATUS_PENDING,
			decimal.NullDecimal{},
			decimal.NullDecimal{},
			nil,
			err.Error(),
		), err
	}

	if tx.Status.Value() == 0 {
		return NewTxResult(
			types.TX_STATUS_FAILED,
			gasUsed,
			effectiveGasPrice,
			blockHeight,
			fmt.Sprintf("transaction failed with failure message %s", tx.Failure.MessageValue),
		), nil
	}

	return NewTxResult(
		types.TX_STATUS_SUCCESS,
		gasUsed,
		effectiveGasPrice,
		blockHeight,
		"",
	), nil
}

// WaitForConfirmedTx polls the transaction until it reaches a terminal state.
// It NEVER returns an error. Instead, it always returns a TxResult where:
//   - Status == TX_STATUS_SUCCESS or TX_STATUS_FAILED when confirmed.
//   - Status == TX_STATUS_TIMEOUT if it did not reach the required confirmations
//     within WaitingTxDuration (or the context was canceled); in this case,
//     the result's FailureReason field is populated with details.
//
// The function sleeps for CheckingTxInterval between polls.
func (cp *IconChainProvider) WaitForConfirmedTx(
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
			log.Debug(
				"Failed to check tx status",
				"tx_hash", txHash,
				err,
			)
		}

		switch result.Status {
		case types.TX_STATUS_SUCCESS, types.TX_STATUS_FAILED:
			return result
		case types.TX_STATUS_PENDING:
			log.Debug(
				"Waiting for tx to be mined",
				"tx_hash", txHash,
			)
			time.Sleep(cp.Config.CheckingTxInterval)
		}
	}

	failureReason := fmt.Sprintf("timed out waiting %s for tx %s to be confirmed",
		cp.Config.WaitingTxDuration, txHash)

	if lastErr != nil {
		failureReason = fmt.Sprintf("%s: %v", failureReason, lastErr)
	}

	return NewTxResult(
		types.TX_STATUS_TIMEOUT,
		decimal.NullDecimal{},
		decimal.NullDecimal{},
		nil,
		failureReason,
	)
}

// handleMetrics increments tx count and, for success/failed, records processing time (ms) and gas used.
func (cp *IconChainProvider) handleMetrics(tunnelID uint64, createdAt time.Time, txResult TxResult) {
	// increment the transactions count metric
	relayermetrics.IncTxsCount(tunnelID, cp.ChainName, types.ChainTypeIcon.String(), txResult.Status.String())

	switch txResult.Status {
	case types.TX_STATUS_SUCCESS, types.TX_STATUS_FAILED:
		// track transaction processing time (ms)
		relayermetrics.ObserveTxProcessTime(
			tunnelID,
			cp.ChainName,
			types.ChainTypeIcon.String(),
			txResult.Status.String(),
			time.Since(createdAt).Milliseconds(),
		)

		// track gas used for the relayed transaction
		relayermetrics.ObserveGasUsed(
			tunnelID,
			cp.ChainName,
			types.ChainTypeIcon.String(),
			txResult.Status.String(),
			txResult.GasUsed.Decimal.InexactFloat64(),
		)
	}
}
