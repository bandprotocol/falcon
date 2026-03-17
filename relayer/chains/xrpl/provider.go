package xrpl

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
	"github.com/bandprotocol/falcon/relayer/wallet/xrpl"
)

var _ chains.ChainProvider = (*XRPLChainProvider)(nil)

// XRPLChainProvider handles interactions with XRPL.
type XRPLChainProvider struct {
	Config    *XRPLChainProviderConfig
	ChainName string

	Client Client

	Log logger.Logger

	DB db.Database

	Alert alert.Alert

	FreeSigners chan wallet.Signer
	Wallet      wallet.Wallet
}

// NewXRPLChainProvider creates a new XRPL chain provider.
func NewXRPLChainProvider(
	chainName string,
	client Client,
	cfg *XRPLChainProviderConfig,
	log logger.Logger,
	w wallet.Wallet,
	a alert.Alert,
) *XRPLChainProvider {
	return &XRPLChainProvider{
		Config:      cfg,
		ChainName:   chainName,
		Client:      client,
		Log:         log.With("chain_name", chainName),
		Alert:       a,
		FreeSigners: chains.LoadSigners(w),
		Wallet:      w,
	}
}

// Init connects to the XRPL chain.
func (cp *XRPLChainProvider) Init(ctx context.Context) error {
	if err := cp.Client.Connect(); err != nil {
		return err
	}

	go cp.Client.StartLivelinessCheck(ctx, cp.Config.LivelinessCheckingInterval)

	return nil
}

// SetDatabase assigns the given database instance.
func (cp *XRPLChainProvider) SetDatabase(database db.Database) {
	cp.DB = database
}

// QueryTunnelInfo returns an active tunnel for XRPL. Since XRPL does not track
// sequence on-chain, the latest sequence is sourced from the database when
// available, otherwise 0 is returned and the caller is responsible for the
// fallback logic.
func (cp *XRPLChainProvider) QueryTunnelInfo(
	_ context.Context,
	tunnelID uint64,
	tunnelDestinationAddr string,
) (*types.Tunnel, error) {
	var latestSeqPtr *uint64
	seq := uint64(0)
	if cp.DB != nil {
		if latestTx := cp.DB.GetLatestTransaction(tunnelID); latestTx != nil {
			seq = latestTx.Sequence
		}
		latestSeqPtr = &seq

	}
	tunnel := types.NewTunnel(tunnelID, "", true, latestSeqPtr, nil)
	return tunnel, nil
}

// RelayPacket relays the packet to XRPL OracleSet transaction.
func (cp *XRPLChainProvider) RelayPacket(_ context.Context, packet *bandtypes.Packet) error {
	if err := cp.Client.CheckAndConnect(); err != nil {
		cp.Log.Error("Connect client error", err)
		return fmt.Errorf("[XRPLProvider] failed to connect client: %w", err)
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
	var err error
	sequence := uint32(0)
	for retryCount := 1; retryCount <= cp.Config.MaxRetry; retryCount++ {
		log.Info("Relaying a message", "retry_count", retryCount)

		// If it is the first attempt or previous attempt failed due to sequence error, fetch the latest account sequence number.
		if sequence == 0 {
			sequence, err = cp.Client.GetAccountSequenceNumber(freeSigner.GetAddress())
			if err != nil {
				log.Error("Get account sequence number error", "retry_count", retryCount, err)
				lastErr = err
				time.Sleep(cp.Config.NonceInterval)
				continue
			}
		}

		signing, err := chains.SelectSigning(packet)
		if err != nil {
			log.Error("Select signing error", "retry_count", retryCount, err)
			lastErr = err
			continue
		}

		signerPayload := xrpl.NewSignerPayload(
			freeSigner.GetAddress(),
			packet.TunnelID,
			cp.Config.Fee,
			sequence,
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

		txResult, err := cp.Client.BroadcastTx(txBlob)
		if err != nil {
			log.Error("Broadcast transaction error", "retry_count", retryCount, err)
			lastErr = err

			// save failed tx in db
			if cp.DB != nil {
				tx := cp.prepareTransaction(
					txResult,
					types.TX_STATUS_FAILED,
					freeSigner.GetAddress(),
					packet,
					balance,
					log,
					retryCount,
				)
				chains.HandleSaveTransaction(cp.DB, cp.Alert, tx, log)
			}

			// Set sequence to 0 to fetch the latest account sequence number in the next attempt
			sequence = 0
			continue
		}

		log.Info(
			"Packet is successfully relayed",
			"tx_hash", txResult.TxHash,
			"retry_count", retryCount,
		)

		// save success tx in db
		if cp.DB != nil {
			tx := cp.prepareTransaction(
				txResult,
				types.TX_STATUS_SUCCESS,
				freeSigner.GetAddress(),
				packet,
				balance,
				log,
				retryCount,
			)
			chains.HandleSaveTransaction(cp.DB, cp.Alert, tx, log)
		}

		relayermetrics.IncTxsCount(
			packet.TunnelID,
			cp.ChainName,
			types.ChainTypeXRPL.String(),
			types.TX_STATUS_SUCCESS.String(),
		)
		alert.HandleReset(
			cp.Alert,
			alert.NewTopic(alert.RelayTxErrorMsg).WithTunnelID(packet.TunnelID).WithChainName(cp.ChainName),
		)

		return nil
	}

	alert.HandleAlert(
		cp.Alert,
		alert.NewTopic(alert.RelayTxErrorMsg).WithTunnelID(packet.TunnelID).WithChainName(cp.ChainName),
		lastErr.Error(),
	)
	return fmt.Errorf("[XRPLProvider] failed to relay packet after %d attempts", cp.Config.MaxRetry)
}

// QueryBalance queries balance by given address from the destination chain.
func (cp *XRPLChainProvider) QueryBalance(ctx context.Context, address string) (*big.Int, error) {
	return cp.Client.GetBalance(address)
}

// GetChainName retrieves the chain name from the chain provider.
func (cp *XRPLChainProvider) GetChainName() string { return cp.ChainName }

// ChainType retrieves the chain type from the chain provider.
func (cp *XRPLChainProvider) ChainType() types.ChainType {
	return types.ChainTypeXRPL
}

func (cp *XRPLChainProvider) GetWallet() wallet.Wallet {
	return cp.Wallet
}

// PacketStaleDuration returns 5 minutes: XRPL rejects packets older than this.
func (cp *XRPLChainProvider) PacketStaleDuration() time.Duration {
	return 5 * time.Minute
}

// prepareTransaction prepares the transaction to be stored in the database.
func (cp *XRPLChainProvider) prepareTransaction(
	txResult TxResult,
	txStatus types.TxStatus,
	signerAddress string,
	packet *bandtypes.Packet,
	oldBalance *big.Int,
	log logger.Logger,
	retryCount int,
) *db.Transaction {
	if txResult.TxHash == "" {
		return nil
	}

	var signalPrices []db.SignalPrice
	for _, p := range packet.SignalPrices {
		signalPrices = append(signalPrices, *db.NewSignalPrice(p.SignalID, p.Price))
	}

	fee := decimal.NullDecimal{}
	balanceDelta := decimal.NullDecimal{}

	// Convert fee from string to decimal
	if txResult.Fee != "" {
		feeDecimal, err := decimal.NewFromString(txResult.Fee)
		if err != nil {
			log.Error("Failed to parse fee", "fee", txResult.Fee, "retry_count", retryCount, err)
		} else {
			fee = decimal.NewNullDecimal(feeDecimal)
		}
	}

	// Compute new balance
	// Note: this may be incorrect if other transactions affected the user's balance during this period.
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

	closeTime, err := cp.Client.GetLedgerCloseTime(txResult.LedgerIndex)
	if err != nil {
		log.Error("Failed to get ledger close time", "tx_hash", txResult.TxHash, "retry_count", retryCount, err)
		alert.HandleAlert(cp.Alert, alert.NewTopic(alert.GetLedgerCloseTimeErrorMsg).
			WithTunnelID(packet.TunnelID).
			WithChainName(cp.ChainName), err.Error())
	} else {
		alert.HandleReset(cp.Alert, alert.NewTopic(alert.GetLedgerCloseTimeErrorMsg).
			WithTunnelID(packet.TunnelID).
			WithChainName(cp.ChainName),
		)
	}

	tx := db.NewTransaction(
		txResult.TxHash,
		packet.TunnelID,
		packet.Sequence,
		cp.ChainName,
		types.ChainTypeXRPL,
		signerAddress,
		txStatus,
		decimal.NewNullDecimal(decimal.NewFromInt(1)), // gasUsed - XRPL doesn't have gas, using 1 as placeholder
		fee,
		balanceDelta,
		signalPrices,
		closeTime,
	)

	return tx
}
