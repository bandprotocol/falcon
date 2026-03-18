package xrpl

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Peersyst/xrpl-go/xrpl/ledger-entry-types"
	"github.com/Peersyst/xrpl-go/xrpl/transaction"
	xrpltypes "github.com/Peersyst/xrpl-go/xrpl/transaction/types"
	xrplwallet "github.com/Peersyst/xrpl-go/xrpl/wallet"

	"github.com/bandprotocol/falcon/relayer/wallet"
)

const (
	priceScale = 9
)

// Precomputed hex encodings of the constant provider and dataClass strings.
var (
	providerHex  = hex.EncodeToString([]byte("Band Protocol"))
	dataClassHex = hex.EncodeToString([]byte("currency"))
)

var _ wallet.Signer = (*LocalSigner)(nil)

// LocalSigner uses a local XRPL secret for signing.
type LocalSigner struct {
	Name   string
	Wallet *xrplwallet.Wallet
}

// NewLocalSigner creates a new LocalSigner.
func NewLocalSigner(name string, w *xrplwallet.Wallet) *LocalSigner {
	return &LocalSigner{
		Name:   name,
		Wallet: w,
	}
}

// ExportPrivateKey returns the decrypted XRPL secret.
func (l *LocalSigner) ExportPrivateKey() (string, error) {
	return l.Wallet.PrivateKey, nil
}

// GetName returns the signer's key name.
func (l *LocalSigner) GetName() string {
	return l.Name
}

// GetAddress returns the signer's classic address.
func (l *LocalSigner) GetAddress() string {
	return l.Wallet.ClassicAddress.String()
}

// Sign signs the provided transaction payload and returns the signed tx blob.
func (l *LocalSigner) Sign(payload []byte, tssPayload wallet.TssPayload) ([]byte, error) {
	var signerPayload SignerPayload
	if err := json.Unmarshal(payload, &signerPayload); err != nil {
		return nil, err
	}

	tssPacket, err := tssPayload.ToTSSPacket()
	if err != nil {
		return nil, err
	}

	priceDataSeries := make([]ledger.PriceDataWrapper, 0, len(tssPacket.RelayPrices))

	for _, p := range tssPacket.RelayPrices {
		signalID := wallet.Bytes32ToString(p.SignalID)
		baseAsset, quoteAsset, err := ParseAssetsFromSignal(signalID)
		if err != nil {
			return nil, err
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

	feeUint64, err := strconv.ParseUint(signerPayload.Fee, 10, 64)
	if err != nil {
		return nil, err
	}

	tx := &transaction.OracleSet{
		BaseTx: transaction.BaseTx{
			Account:         xrpltypes.Address(signerPayload.Account),
			TransactionType: transaction.OracleSetTx,
			Sequence:        signerPayload.Sequence,
			Fee:             xrpltypes.XRPCurrencyAmount(feeUint64),
		},
		OracleDocumentID: uint32(
			signerPayload.OracleID,
		), // expect oracle ID (tunnel ID) to fit in 32 bits, which is sufficient for our use case
		LastUpdatedTime: uint32(
			tssPacket.CreatedAt,
		), // expect createdAt to fit in 32 bits, which is sufficient until the year 2106
		Provider:        providerHex,
		AssetClass:      dataClassHex,
		PriceDataSeries: priceDataSeries,
	}

	flattenedTx := tx.Flatten()

	formattedPriceTx, err := formatAssetPrice(flattenedTx)
	if err != nil {
		return nil, err
	}

	flattenedTx["PriceDataSeries"] = formattedPriceTx

	txBlob, _, err := l.Wallet.Sign(flattenedTx)
	if err != nil {
		return nil, err
	}

	return []byte(txBlob), nil
}

func formatAssetPrice(tx map[string]any) ([]map[string]any, error) {
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
