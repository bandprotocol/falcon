package icon

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

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

	ContractAddress string
	NetworkID       string
	FreeSigners     chan wallet.Signer
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
		Config:          cfg,
		ChainName:       chainName,
		Client:          client,
		Log:             log.With("chain_name", chainName),
		Wallet:          wallet,
		Alert:           alert,
		ContractAddress: cfg.ContractAddress,
		NetworkID:       cfg.NetworkID,
		FreeSigners:     chains.LoadSigners(wallet),
	}
}

// Init connects to the Icon chain.
func (cp *IconChainProvider) Init(ctx context.Context) error {
	if err := cp.Client.Connect(); err != nil {
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
			cp.ContractAddress,
			uint64(cp.Config.StepLimit),
			cp.NetworkID,
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

			// save failed tx in db
			if cp.DB != nil {
				tx := cp.prepareTransaction(ctx, txHash, freeSigner.GetAddress(), packet, balance, log, retryCount)
				chains.HandleSaveTransaction(cp.DB, cp.Alert, tx, log)
			}

			continue
		}

		log.Info(
			"Packet is successfully relayed",
			"tx_hash", txHash,
			"retry_count", retryCount,
		)

		// save success tx in db
		if cp.DB != nil {
			tx := cp.prepareTransaction(ctx, txHash, freeSigner.GetAddress(), packet, balance, log, retryCount)
			chains.HandleSaveTransaction(cp.DB, cp.Alert, tx, log)
		}

		relayermetrics.IncTxsCount(packet.TunnelID, cp.ChainName, types.ChainTypeIcon.String(), types.TX_STATUS_SUCCESS.String())
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

// LoadSigners loads signers to prepare to relay the packet.
func (cp *IconChainProvider) LoadSigners() error {
	cp.FreeSigners = chains.LoadSigners(cp.Wallet)
	return nil
}

// prepareTransaction prepares the transaction to be stored in the database.
func (cp *IconChainProvider) prepareTransaction(
	ctx context.Context,
	txHash string,
	signerAddress string,
	packet *bandtypes.Packet,
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

	// For ICON, we don't have detailed fee/gas info like EVM
	// Using placeholder values
	fee := decimal.NullDecimal{}
	balanceDelta := decimal.NullDecimal{}

	// Compute new balance if old balance was provided
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

	tx := db.NewTransaction(
		txHash,
		packet.TunnelID,
		packet.Sequence,
		cp.ChainName,
		types.ChainTypeIcon,
		signerAddress,
		types.TX_STATUS_SUCCESS, // Assume success for now, could be updated later
		decimal.NewNullDecimal(decimal.NewFromInt(1)), // gasUsed placeholder
		fee,
		balanceDelta,
		signalPrices,
		nil, // blockTimestamp
	)

	return tx
}
