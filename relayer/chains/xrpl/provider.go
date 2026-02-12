package xrpl

import (
	"context"
	"encoding/hex"
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

const (
	provider   = "Band Protocol"
	dataClass  = "currency"
	priceScale = 9
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
		Config:    cfg,
		ChainName: chainName,
		Client:    client,
		Log:       log.With("chain_name", chainName),
		Wallet:    wallet,
		Alert:     alert,
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

		tx, err := cp.buildOracleSetTx(packet, freeSigner.GetAddress(), sequence)
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

		// use binarycodec.Encode instead of []byte to prevent type change when decode again
		encodedTx, err := binarycodec.Encode(tx)
		if err != nil {
			log.Error("Encode transaction error", "retry_count", retryCount, err)
			lastErr = err
			continue
		}

		signing, err := chains.SelectSigning(packet)
		if err != nil {
			log.Error("Select signing error", "retry_count", retryCount, err)
			lastErr = err
			continue
		}

		preSignPayload := wallet.NewPreSignPayload(
			signing.Message,
			signing.EVMSignature.RAddress,
			signing.EVMSignature.Signature,
		)

		txBlobBytes, err := freeSigner.Sign([]byte(encodedTx), preSignPayload)
		if err != nil {
			log.Error("Sign transaction error", "retry_count", retryCount, err)
			lastErr = err
			continue
		}

		txBlobStr := hex.EncodeToString(txBlobBytes)

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

		txResult, err := cp.Client.BroadcastTx(ctx, txBlobStr)
		if err != nil {
			log.Error("Broadcast transaction error", "retry_count", retryCount, err)
			lastErr = err

			// save failed tx in db
			if cp.DB != nil {
				tx := cp.prepareTransaction(ctx, txResult, types.TX_STATUS_FAILED, freeSigner.GetAddress(), packet, balance, log, retryCount)
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
			tx := cp.prepareTransaction(ctx, txResult, types.TX_STATUS_SUCCESS, freeSigner.GetAddress(), packet, balance, log, retryCount)
			chains.HandleSaveTransaction(cp.DB, cp.Alert, tx, log)
		}

		relayermetrics.IncTxsCount(packet.TunnelID, cp.ChainName, types.ChainTypeXRPL.String(), types.TX_STATUS_SUCCESS.String())
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

// LoadSigners loads signers to prepare to relay the packet.
func (cp *XRPLChainProvider) LoadSigners() error {
	cp.FreeSigners = chains.LoadSigners(cp.Wallet)
	return nil
}

func (cp *XRPLChainProvider) buildOracleSetTx(
	packet *bandtypes.Packet,
	signerAddress string,
	sequence uint32,
) (transaction.FlatTransaction, error) {
	providerHex, err := StringToHex(provider, 0)
	if err != nil {
		return transaction.FlatTransaction{}, err
	}
	dataClassHex, err := StringToHex(dataClass, 0)
	if err != nil {
		return transaction.FlatTransaction{}, err
	}

	priceDataSeries := make([]ledger.PriceDataWrapper, 0, len(packet.SignalPrices))

	for _, p := range packet.SignalPrices {
		baseAsset, quoteAsset, err := ParseAssetsFromSignal(p.SignalID)
		if err != nil {
			return transaction.FlatTransaction{}, err
		}

		priceDataSeries = append(priceDataSeries, ledger.PriceDataWrapper{
			PriceData: ledger.PriceData{
				BaseAsset:  baseAsset,
				QuoteAsset: quoteAsset,
				AssetPrice: p.Price,
				Scale:      priceScale,
			},
		})
	}

	tx := &transaction.OracleSet{
		BaseTx: transaction.BaseTx{
			Account:         xrpltypes.Address(signerAddress),
			TransactionType: transaction.OracleSetTx,
			Sequence:        sequence,
			Fee:             xrpltypes.XRPCurrencyAmount(cp.Config.Fee),
		},
		OracleDocumentID: uint32(packet.TunnelID),
		LastUpdatedTime:  uint32(time.Now().Unix()),
		Provider:         providerHex,
		AssetClass:       dataClassHex,
		PriceDataSeries:  priceDataSeries,
	}

	flattenedTx := tx.Flatten()

	formattedPriceTx, err := FormatAssetPrice(flattenedTx)
	if err != nil {
		return transaction.FlatTransaction{}, err
	}

	flattenedTx["PriceDataSeries"] = formattedPriceTx

	return flattenedTx, nil
}

func FormatAssetPrice(tx map[string]any) ([]map[string]any, error) {
	// Look for the PriceDataSeries in the flattened map
	priceDataSeries, ok := tx["PriceDataSeries"].([]map[string]any)
	if !ok {
		// If it's not there, it's either not an OracleSet or already empty
		return nil, nil
	}

	for i, priceDataWrapper := range priceDataSeries {
		// Access the inner PriceData object
		priceData, ok := priceDataWrapper["PriceData"].(map[string]any)
		if !ok {
			return nil, fmt.Errorf("failed to get PriceData at index %d", i)
		}

		// Ensure AssetPrice is a uint64 hex string
		// In Go, after transaction flattening, AssetPrice might be a uint64 or a string
		if rawPrice, ok := priceData["AssetPrice"]; ok {
			var hexStr string
			var err error
			switch v := rawPrice.(type) {
			case string:
				hexStr, err = Uint64StrToHexStr(v)
				if err != nil {
					return nil, fmt.Errorf("failed to convert AssetPrice string at index %d: %w", i, err)
				}
			case uint64:
				hexStr = fmt.Sprintf("%016X", v)
			default:
				return nil, fmt.Errorf("unexpected type for AssetPrice at index %d: %T", i, v)
			}

			// Re-assign as a clean uint64 hex string.
			// The XRPL Binary Codec will encode this into the correct 8-byte big-endian format.
			priceData["AssetPrice"] = hexStr
		} else {
			return nil, fmt.Errorf("failed to get AssetPrice at index %d", i)
		}
	}

	return priceDataSeries, nil
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
