package xrpl

import (
	"context"
	"fmt"
	"math/big"
	"time"

	binarycodec "github.com/Peersyst/xrpl-go/binary-codec"
	ledger "github.com/Peersyst/xrpl-go/xrpl/ledger-entry-types"
	"github.com/Peersyst/xrpl-go/xrpl/transaction"
	xrpltypes "github.com/Peersyst/xrpl-go/xrpl/transaction/types"
	"github.com/shopspring/decimal"

	"github.com/bandprotocol/falcon/internal/relayermetrics"
	"github.com/bandprotocol/falcon/relayer/alert"
	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/db"
	"github.com/bandprotocol/falcon/relayer/logger"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

var _ chains.ChainProvider = (*XRPLChainProvider)(nil)

// XRPLChainProvider handles interactions with XRPL.
type XRPLChainProvider struct {
	Config    *XRPLChainProviderConfig
	ChainName string
	// OracleAccount is derived from the XRPL wallet signers at runtime.
	OracleAccount string

	Client *Client

	Log logger.Logger

	Wallet wallet.Wallet
	DB     db.Database

	Alert alert.Alert

	FreeSigners chan wallet.Signer

	nonceInterval time.Duration
}

// NewXRPLChainProvider creates a new XRPL chain provider.
func NewXRPLChainProvider(
	chainName string,
	client *Client,
	cfg *XRPLChainProviderConfig,
	log logger.Logger,
	wallet wallet.Wallet,
	alert alert.Alert,
) (*XRPLChainProvider, error) {
	if cfg.PriceScale == 0 {
		cfg.PriceScale = 9
	}
	if cfg.PriceScale > uint32(ledger.PriceDataScaleMax) {
		return nil, fmt.Errorf(
			"price_scale %d exceeds max %d",
			cfg.PriceScale,
			ledger.PriceDataScaleMax,
		)
	}

	return &XRPLChainProvider{
		Config:        cfg,
		ChainName:     chainName,
		Client:        client,
		Log:           log.With("chain_name", chainName),
		Wallet:        wallet,
		Alert:         alert,
		nonceInterval: time.Second,
	}, nil
}

// Init connects to the XRPL chain.
func (cp *XRPLChainProvider) Init(ctx context.Context) error {
	if err := cp.Client.Connect(ctx); err != nil {
		return err
	}

	return nil
}

// SetDatabase assigns the given database instance.
func (cp *XRPLChainProvider) SetDatabase(database db.Database) {
	cp.DB = database
}

// QueryTunnelInfo returns a best-effort tunnel info for XRPL.
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
	if cp.FreeSigners == nil {
		return fmt.Errorf("signers not loaded")
	}
	signer := <-cp.FreeSigners
	defer func() {
		cp.FreeSigners <- signer
	}()

	log := cp.Log.With(
		"tunnel_id", packet.TunnelID,
		"sequence", packet.Sequence,
		"signer_address", signer.GetAddress(),
	)

	var lastErr error
	var err error
	sequence := uint64(0)
	for retryCount := 1; retryCount <= cp.Config.MaxRetry; retryCount++ {
		log.Info("Relaying a message", "retry_count", retryCount)

		// If it is the first attempt or previous attempt failed due to sequence error, fetch the latest account sequence number.
		if sequence == 0 {
			sequence, err = cp.Client.GetAccountSequenceNumber(ctx, signer.GetAddress())
			if err != nil {
				log.Error("Get account sequence number error", "retry_count", retryCount, err)
				lastErr = err
				time.Sleep(cp.nonceInterval)
				continue
			}
		}

		tx, err := cp.buildOracleSetTx(packet, signer.GetAddress(), sequence)
		if err != nil {
			log.Error("Build OracleSet transaction error", "retry_count", retryCount, err)
			lastErr = err
			continue
		}

		if err := cp.Client.Autofill(&tx); err != nil {
			log.Error("Autofill transaction error", "retry_count", retryCount, err)
			lastErr = err
			continue
		}

		encodedTx, err := binarycodec.Encode(tx)
		if err != nil {
			log.Error("Encode transaction error", "retry_count", retryCount, err)
			lastErr = err
			continue
		}

		txBlobBytes, err := signer.Sign([]byte(encodedTx))
		if err != nil {
			log.Error("Sign transaction error", "retry_count", retryCount, err)
			lastErr = err
			continue
		}

		txHash, err := cp.Client.BroadcastTx(ctx, string(txBlobBytes))
		if err != nil {
			log.Error("Broadcast transaction error", "retry_count", retryCount, err)
			lastErr = err
			continue
		}

		log.Info(
			"Packet is successfully relayed",
			"tx_hash", txHash,
			"retry_count", retryCount,
		)

		cp.saveRelayTx(packet, txHash)
		relayermetrics.IncTxsCount(packet.TunnelID, cp.ChainName, types.TX_STATUS_SUCCESS.String())

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

// LoadSigners loads signers to prepare to relay the packet.
func (cp *XRPLChainProvider) LoadSigners() error {
	cp.FreeSigners = chains.LoadSigners(cp.Wallet)
	return nil
}

func (cp *XRPLChainProvider) buildOracleSetTx(
	packet *bandtypes.Packet,
	signerAddress string,
	sequence uint64,
) (transaction.FlatTransaction, error) {
	providerHex, err := stringToHex("Band Protocol", 0)
	if err != nil {
		return transaction.FlatTransaction{}, err
	}
	dataClassHex, err := stringToHex("currency", 0)
	if err != nil {
		return transaction.FlatTransaction{}, err
	}

	priceDataSeries := make([]ledger.PriceDataWrapper, 0, len(packet.SignalPrices))

	for _, p := range packet.SignalPrices {
		baseAsset, quoteAsset, err := parseAssetsFromSignal(p.SignalID)
		if err != nil {
			return transaction.FlatTransaction{}, err
		}

		priceDataSeries = append(priceDataSeries, ledger.PriceDataWrapper{
			PriceData: ledger.PriceData{
				BaseAsset:  baseAsset,
				QuoteAsset: quoteAsset,
				AssetPrice: p.Price,
				Scale:      uint8(cp.Config.PriceScale),
			},
		})
	}

	tx := &transaction.OracleSet{
		BaseTx: transaction.BaseTx{
			Account:         xrpltypes.Address(signerAddress),
			TransactionType: transaction.OracleSetTx,
			Sequence:        uint32(sequence),
			Fee:             xrpltypes.XRPCurrencyAmount(12),
		},
		OracleDocumentID: uint32(cp.Config.OracleID),
		LastUpdatedTime:  uint32(time.Now().Unix()),
		Provider:         providerHex,
		AssetClass:       dataClassHex,
		PriceDataSeries:  priceDataSeries,
	}

	return tx.Flatten(), nil
}

func (cp *XRPLChainProvider) saveRelayTx(packet *bandtypes.Packet, txHash string) {
	signalPrices := make([]db.SignalPrice, 0, len(packet.SignalPrices))
	for _, p := range packet.SignalPrices {
		signalPrices = append(signalPrices, *db.NewSignalPrice(p.SignalID, p.Price))
	}

	tx := db.NewTransaction(
		txHash,
		packet.TunnelID,
		packet.Sequence,
		cp.ChainName,
		types.ChainTypeXRPL,
		cp.OracleAccount,
		types.TX_STATUS_SUCCESS,
		decimal.NullDecimal{},
		decimal.NullDecimal{},
		decimal.NullDecimal{},
		signalPrices,
		nil,
	)

	if cp.DB == nil {
		return
	}

	if err := cp.DB.AddOrUpdateTransaction(tx); err != nil {
		cp.Log.Error("Save transaction error", err)
		alert.HandleAlert(cp.Alert, alert.NewTopic(alert.SaveDatabaseErrorMsg).
			WithTunnelID(tx.TunnelID).
			WithChainName(cp.ChainName), err.Error())
	} else {
		alert.HandleReset(cp.Alert, alert.NewTopic(alert.SaveDatabaseErrorMsg).
			WithTunnelID(tx.TunnelID).
			WithChainName(cp.ChainName))
	}
}
