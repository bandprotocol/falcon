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

	Wallet wallet.Wallet
	DB     db.Database

	Alert alert.Alert

	FreeSigners chan wallet.Signer
}

// NewXRPLChainProvider creates a new XRPL chain provider.
func NewXRPLChainProvider(
	chainName string,
	client Client,
	cfg *XRPLChainProviderConfig,
	log logger.Logger,
	wallet wallet.Wallet,
	alert alert.Alert,
) *XRPLChainProvider {
	return &XRPLChainProvider{
		Config:      cfg,
		ChainName:   chainName,
		Client:      client,
		Log:         log.With("chain_name", chainName),
		Wallet:      wallet,
		Alert:       alert,
		FreeSigners: chains.LoadSigners(wallet),
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

// QueryTunnelInfo always returns active tunnel and 0 sequence
// as XRPL Oracleset doesn't have the concept of tunnel and sequence.
func (cp *XRPLChainProvider) QueryTunnelInfo(
	ctx context.Context,
	tunnelID uint64,
	tunnelDestinationAddr string,
) (*types.Tunnel, error) {
	tunnel := types.NewTunnel(tunnelID, tunnelDestinationAddr, true, 0, nil)
	return tunnel, nil
}

// RelayPacket relays the packet to XRPL OracleSet transaction.
func (cp *XRPLChainProvider) RelayPacket(ctx context.Context, packet *bandtypes.Packet) error {
	if err := cp.Client.CheckAndConnect(); err != nil {
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
	var err error
	sequence := uint32(0)
	for retryCount := 1; retryCount <= cp.Config.MaxRetry; retryCount++ {
		log.Info("Relaying a message", "retry_count", retryCount)

		// If it is the first attempt or previous attempt failed due to sequence error, fetch the latest account sequence number.
		if sequence == 0 {
			sequence, err = cp.Client.GetAccountSequenceNumber(ctx, freeSigner.GetAddress())
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

		rAddress := []byte{}
		signature := []byte{}
		if signing.EVMSignature != nil {
			rAddress = signing.EVMSignature.RAddress
			signature = signing.EVMSignature.Signature
		}
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

		txResult, err := cp.Client.BroadcastTx(ctx, txBlob)
		if err != nil {
			log.Error("Broadcast transaction error", "retry_count", retryCount, err)
			lastErr = err

			// save failed tx in db
			if cp.DB != nil {
				tx := cp.prepareTransaction(
					ctx,
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
				ctx,
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
	return fmt.Errorf("failed to relay packet after %d attempts", cp.Config.MaxRetry)
}

// QueryBalance queries balance by given key name from the destination chain.
func (cp *XRPLChainProvider) QueryBalance(ctx context.Context, keyName string) (*big.Int, error) {
	signer, ok := cp.Wallet.GetSigner(keyName)
	if !ok {
		cp.Log.Error("Key name does not exist", "key_name", keyName)
		return nil, fmt.Errorf("key name does not exist: %s", keyName)
	}

	return cp.Client.GetBalance(ctx, signer.GetAddress())
}

// GetChainName retrieves the chain name from the chain provider.
func (cp *XRPLChainProvider) GetChainName() string { return cp.ChainName }

// ChainType retrieves the chain type from the chain provider.
func (cp *XRPLChainProvider) ChainType() types.ChainType {
	return types.ChainTypeXRPL
}

// prepareTransaction prepares the transaction to be stored in the database.
func (cp *XRPLChainProvider) prepareTransaction(
	ctx context.Context,
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
		nil, // blockTimestamp - XRPL doesn't provide this easily
	)

	return tx
}
