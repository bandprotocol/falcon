package geth

import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/bandprotocol/falcon/relayer/wallet"
)

var _ wallet.Signer = (*LocalSigner)(nil)

// LocalSigner is signer that locally stored ECDSA private key.
type LocalSigner struct {
	Name       string
	privateKey *ecdsa.PrivateKey
}

// NewLocalSigner creates a new LocalSigner instance
func NewLocalSigner(
	name string,
	privateKey *ecdsa.PrivateKey,
) *LocalSigner {
	return &LocalSigner{
		Name:       name,
		privateKey: privateKey,
	}
}

// ExportPrivateKey returns the signer's private key.
func (l *LocalSigner) ExportPrivateKey() (string, error) {
	b := crypto.FromECDSA(l.privateKey)

	return hex.EncodeToString(b), nil
}

// GetName returns the signer's key name.
func (l *LocalSigner) GetName() (addr string) {
	return l.Name
}

// GetAddress returns the signer's address.
func (l *LocalSigner) GetAddress() (addr string) {
	return crypto.PubkeyToAddress(l.privateKey.PublicKey).String()
}

// Sign hashes the input data which is RLP encoded with Keccak256, signs it, and returns the signature.
func (l *LocalSigner) Sign(data []byte) ([]byte, error) {
	hash := crypto.Keccak256(data)

	return crypto.Sign(hash, l.privateKey)
}
