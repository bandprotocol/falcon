package xrpl

import (
	"encoding/json"
	"strconv"

	"github.com/Peersyst/xrpl-go/xrpl/ledger-entry-types"
	"github.com/Peersyst/xrpl-go/xrpl/transaction"
	xrpltypes "github.com/Peersyst/xrpl-go/xrpl/transaction/types"
	xrplwallet "github.com/Peersyst/xrpl-go/xrpl/wallet"

	"github.com/bandprotocol/falcon/relayer/wallet"
)

const (
	provider   = "Band Protocol"
	dataClass  = "Currency"
	priceScale = 9
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
func (l *LocalSigner) Sign(payload []byte, _ wallet.TssPayload) ([]byte, error) {
	var signerPayload SignerPayload
	if err := json.Unmarshal(payload, &signerPayload); err != nil {
		return nil, err
	}
	providerHex, err := StringToHex(provider, 0)
	if err != nil {
		return nil, err
	}
	dataClassHex, err := StringToHex(dataClass, 0)
	if err != nil {
		return nil, err
	}

	priceDataSeries := make([]ledger.PriceDataWrapper, 0, len(signerPayload.SignalPrices))

	for _, p := range signerPayload.SignalPrices {
		baseAsset, quoteAsset, err := ParseAssetsFromSignal(p.SignalID)
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
			Sequence:        uint32(signerPayload.Sequence),
			Fee:             xrpltypes.XRPCurrencyAmount(feeUint64),
		},
		OracleDocumentID: uint32(signerPayload.OracleId),
		LastUpdatedTime:  uint32(signerPayload.LastUpdatedTime),
		Provider:         providerHex,
		AssetClass:       dataClassHex,
		PriceDataSeries:  priceDataSeries,
	}

	flattenedTx := tx.Flatten()

	txBlob, _, err := l.Wallet.Sign(flattenedTx)
	if err != nil {
		return nil, err
	}

	return []byte(txBlob), nil
}
