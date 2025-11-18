package geth

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/bandprotocol/falcon/relayer/wallet"
)

var _ wallet.Signer = (*LocalSigner)(nil)

// LocalSigner is signer that locally stored ECDSA private key.
type LocalSigner struct {
	Name       string
	address    common.Address
	store      *keystore.KeyStore
	passphrase string

	privateKey *ecdsa.PrivateKey
	mu         sync.Mutex
	loaded     bool
}

// NewLocalSigner creates a new LocalSigner instance
func NewLocalSigner(
	name string,
	address common.Address,
	store *keystore.KeyStore,
	passphrase string,
) *LocalSigner {
	return &LocalSigner{
		Name:       name,
		address:    address,
		store:      store,
		passphrase: passphrase,
		loaded:     false,
	}
}

// loadPrivateKey decrypts and loads the private key (called once, on first use)
func (l *LocalSigner) loadPrivateKey() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.loaded {
		return nil
	}
	start := time.Now()
	fmt.Printf("Loading private key for LocalSigner:\n")
	fmt.Printf("  Name: %s\n", l.Name)
	fmt.Printf("  Address: %s\n", l.address.String())

	acc := accounts.Account{Address: l.address}
	exportStart := time.Now()
	fmt.Printf("Starting keystore export...\n")
	b, err := l.store.Export(acc, l.passphrase, l.passphrase)
	exportDuration := time.Since(exportStart)
	if err != nil {
		fmt.Printf("Export failed after %v: %v\n", exportDuration, err)
		return err
	}
	fmt.Printf("Export completed in %v\n", exportDuration)

	decryptStart := time.Now()
	fmt.Printf("Starting key decryption...\n")
	gethKey, err := keystore.DecryptKey(b, l.passphrase)
	decryptDuration := time.Since(decryptStart)
	if err != nil {
		fmt.Printf("Decryption failed after %v: %v\n", decryptDuration, err)
		return err
	}
	fmt.Printf("Decryption completed in %v\n", decryptDuration)

	l.privateKey = gethKey.PrivateKey
	l.loaded = true

	totalDuration := time.Since(start)
	fmt.Printf("Private key loading completed for %s (%s) in %v\n", l.Name, l.address.String(), totalDuration)

	return nil
}

// ExportPrivateKey returns the signer's private key.
func (l *LocalSigner) ExportPrivateKey() (string, error) {
	if err := l.loadPrivateKey(); err != nil {
		return "", err
	}

	b := crypto.FromECDSA(l.privateKey)
	return hex.EncodeToString(b), nil
}

// GetName returns the signer's key name.
func (l *LocalSigner) GetName() (addr string) {
	return l.Name
}

// GetAddress returns the signer's address.
func (l *LocalSigner) GetAddress() (addr string) {
	return l.address.String()
}

// Sign hashes the input data which is RLP encoded with Keccak256, signs it, and returns the signature.
func (l *LocalSigner) Sign(data []byte) ([]byte, error) {
	if err := l.loadPrivateKey(); err != nil {
		return nil, err
	}

	hash := crypto.Keccak256(data)
	return crypto.Sign(hash, l.privateKey)
}
