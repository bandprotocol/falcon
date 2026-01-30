package xrpl

import (
	binarycodec "github.com/Peersyst/xrpl-go/binary-codec"
	xrplwallet "github.com/Peersyst/xrpl-go/xrpl/wallet"

	"github.com/bandprotocol/falcon/relayer/wallet"
)

var _ wallet.Signer = (*LocalSigner)(nil)

// LocalSigner uses a local XRPL secret for signing.
type LocalSigner struct {
	Name   string
	Wallet xrplwallet.Wallet
}

// NewLocalSigner creates a new LocalSigner.
func NewLocalSigner(name string, w xrplwallet.Wallet) *LocalSigner {
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
func (l *LocalSigner) GetAddress() (addr string) {
	return l.Wallet.ClassicAddress.String()
}

// Sign signs the provided transaction payload and returns the signed tx blob.
func (l *LocalSigner) Sign(data []byte) ([]byte, error) {
	tx, err := binarycodec.Decode(string(data))
	if err != nil {
		return nil, err
	}

	txBlob, _, err := l.Wallet.Sign(tx)
	if err != nil {
		return nil, err
	}

	return []byte(txBlob), nil
}
