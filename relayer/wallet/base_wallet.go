package wallet

import (
	"fmt"
	"path"

	toml "github.com/pelletier/go-toml/v2"

	"github.com/bandprotocol/falcon/internal/os"
)

// WalletAdapter defines chain-specific wallet operations.
// Implement this interface to add support for a new chain type.
type WalletAdapter interface {
	// NormalizeAddress returns the canonical form of the address, or an error if it is invalid.
	// Use the returned form for storage and comparison.
	NormalizeAddress(addr string) (string, error)

	// DeriveFromPrivateKey parses the private key and creates a signer
	// without persisting the secret.
	DeriveFromPrivateKey(name, privateKey string) (Signer, error)

	// DeriveFromMnemonic derives from a mnemonic and creates a signer
	// without persisting the secret.
	DeriveFromMnemonic(
		name, mnemonic string,
		coinType uint32,
		account uint,
		index uint,
	) (Signer, error)

	// PersistKey stores the secret in chain-specific secure storage.
	// The secret parameter carries the original input (private key or mnemonic).
	PersistKey(name string, signer Signer, secret string) error

	// LoadSigner reconstructs a Signer from a persisted KeyRecord, handling both local and remote types.
	LoadSigner(name string, record KeyRecord) (Signer, error)

	// DeleteLocalSecret removes the locally stored secret for the named key.
	// It should be a no-op for remote signers.
	DeleteLocalSecret(name string, signer Signer) error
}

// BaseWallet provides shared wallet logic for all chain types.
// It delegates chain-specific operations to the embedded WalletAdapter.
type BaseWallet struct {
	Adapter WalletAdapter
	keyDir  []string
	Signers map[string]Signer
}

var _ Wallet = (*BaseWallet)(nil)

// NewBaseWallet creates a BaseWallet for the given chain, using the standard
// key record directory: {homePath}/keys/{chainName}/metadata.
func NewBaseWallet(homePath, chainName string, adapter WalletAdapter) (*BaseWallet, error) {
	keyDir := []string{homePath, "keys", chainName, "metadata"}
	records, err := LoadKeyRecords(path.Join(keyDir...))
	if err != nil {
		return nil, err
	}

	signers := make(map[string]Signer)
	for name, record := range records {
		signer, err := adapter.LoadSigner(name, record)
		if err != nil {
			return nil, err
		}
		signers[name] = signer
	}

	return &BaseWallet{Adapter: adapter, keyDir: keyDir, Signers: signers}, nil
}

// SaveByPrivateKey imports from a raw private key.
func (w *BaseWallet) SaveByPrivateKey(name string, privateKey string) (string, error) {
	if _, ok := w.Signers[name]; ok {
		return "", fmt.Errorf("key name exists: %s", name)
	}

	signer, err := w.Adapter.DeriveFromPrivateKey(name, privateKey)
	if err != nil {
		return "", err
	}
	addr := signer.GetAddress()

	if w.IsAddressExist(addr) {
		return "", fmt.Errorf("address exists: %s", addr)
	}

	if err := w.Adapter.PersistKey(name, signer, privateKey); err != nil {
		return "", err
	}

	record := NewKeyRecord(LocalSignerType, addr, "", nil)
	if err := w.saveKeyRecord(name, record); err != nil {
		return "", err
	}

	w.Signers[name] = signer
	return addr, nil
}

// SaveByMnemonic imports from a mnemonic phrase.
func (w *BaseWallet) SaveByMnemonic(
	name string,
	mnemonic string,
	coinType uint32,
	account uint,
	index uint,
) (string, error) {
	if _, ok := w.Signers[name]; ok {
		return "", fmt.Errorf("key name exists: %s", name)
	}

	if mnemonic == "" {
		return "", fmt.Errorf("mnemonic is empty")
	}

	signer, err := w.Adapter.DeriveFromMnemonic(name, mnemonic, coinType, account, index)
	if err != nil {
		return "", err
	}
	addr := signer.GetAddress()

	if w.IsAddressExist(addr) {
		return "", fmt.Errorf("address exists: %s", addr)
	}

	if err := w.Adapter.PersistKey(name, signer, mnemonic); err != nil {
		return "", err
	}

	record := NewKeyRecord(LocalSignerType, addr, "", nil)
	if err := w.saveKeyRecord(name, record); err != nil {
		return "", err
	}

	w.Signers[name] = signer
	return addr, nil
}

// SaveRemoteSignerKey registers a remote signer.
func (w *BaseWallet) SaveRemoteSignerKey(name, address, url string, key *string) error {
	if _, ok := w.Signers[name]; ok {
		return fmt.Errorf("key name exists: %s", name)
	}

	normalizedAddr, err := w.Adapter.NormalizeAddress(address)
	if err != nil {
		return err
	}

	if w.IsAddressExist(normalizedAddr) {
		return fmt.Errorf("address exists: %s", normalizedAddr)
	}

	record := NewKeyRecord(RemoteSignerType, normalizedAddr, url, key)
	if err := w.saveKeyRecord(name, record); err != nil {
		return err
	}

	signer, err := w.Adapter.LoadSigner(name, record)
	if err != nil {
		return err
	}

	w.Signers[name] = signer
	return nil
}

// DeleteKey removes the signer and its records.
func (w *BaseWallet) DeleteKey(name string) error {
	signer, ok := w.Signers[name]
	if !ok {
		return fmt.Errorf("key name does not exist: %s", name)
	}

	if err := w.Adapter.DeleteLocalSecret(name, signer); err != nil {
		return err
	}

	if err := w.deleteKeyRecord(name); err != nil {
		return err
	}

	delete(w.Signers, name)
	return nil
}

// GetSigners lists all signers.
func (w *BaseWallet) GetSigners() []Signer {
	signers := make([]Signer, 0, len(w.Signers))
	for _, signer := range w.Signers {
		signers = append(signers, signer)
	}

	return signers
}

// GetSigner returns the signer with the given name and a flag indicating if it was found.
func (w *BaseWallet) GetSigner(name string) (Signer, bool) {
	signer, ok := w.Signers[name]
	return signer, ok
}

// IsAddressExist returns true if the given address is already registered.
func (w *BaseWallet) IsAddressExist(address string) bool {
	normNew, err := w.Adapter.NormalizeAddress(address)
	if err != nil {
		return false
	}
	for _, signer := range w.Signers {
		normExist, _ := w.Adapter.NormalizeAddress(signer.GetAddress())
		if normExist == normNew {
			return true
		}
	}
	return false
}

// saveKeyRecord writes the KeyRecord to disk.
func (w *BaseWallet) saveKeyRecord(name string, record KeyRecord) error {
	b, err := toml.Marshal(record)
	if err != nil {
		return err
	}

	return os.Write(b, append(w.keyDir, fmt.Sprintf("%s.toml", name)))
}

// deleteKeyRecord deletes the KeyRecord file from disk.
func (w *BaseWallet) deleteKeyRecord(name string) error {
	filePath := path.Join(append(w.keyDir, fmt.Sprintf("%s.toml", name))...)
	return os.DeletePath(filePath)
}
