package xrpl

import (
	"encoding/hex"
	"encoding/json"
	"strconv"
	"strings"

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
	providerHex  = strings.ToUpper(hex.EncodeToString([]byte("Band Protocol")))
	dataClassHex = strings.ToUpper(hex.EncodeToString([]byte("Currency")))
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

	txBlob, _, err := l.Wallet.Sign(flattenedTx)
	if err != nil {
		return nil, err
	}

	return []byte(txBlob), nil
}
