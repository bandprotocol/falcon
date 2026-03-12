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
	// GetKeyDir returns the directory path components for key record storage.
	GetKeyDir() []string

	// ValidateAddress checks if the address is valid for this chain.
	ValidateAddress(address string) error

	// CompareAddresses returns true if two addresses refer to the same account.
	CompareAddresses(a, b string) bool

	// DeriveFromPrivateKey parses the private key and creates a signer
	// without persisting the secret. Returns the derived address and signer.
	DeriveFromPrivateKey(name, privateKey string) (addr string, signer Signer, err error)

	// DeriveFromMnemonic derives from a mnemonic and creates a signer
	// without persisting the secret. Returns the derived address, signer, and save method.
	DeriveFromMnemonic(
		name, mnemonic string,
		coinType uint32,
		account uint,
		index uint,
	) (addr string, signer Signer, saveMethod string, err error)

	// PersistKey stores the secret in chain-specific secure storage.
	// The secret parameter carries the original input (private key or mnemonic).
	PersistKey(name string, signer Signer, secret string) error

	// LoadLocalSigner reconstructs a local signer from a persisted KeyRecord.
	LoadLocalSigner(name string, record KeyRecord) (Signer, error)

	// NewRemoteSigner creates a remote signer that delegates signing to an external KMS.
	NewRemoteSigner(name, address, url string, key *string) (Signer, error)

	// DeleteLocalSecret removes the locally stored secret for the named key.
	// It should be a no-op for remote signers.
	DeleteLocalSecret(name string, signer Signer) error
}

// BaseWallet provides shared wallet logic for all chain types.
// It delegates chain-specific operations to the embedded WalletAdapter.
type BaseWallet struct {
	Adapter WalletAdapter
	Signers map[string]Signer
}

var _ Wallet = (*BaseWallet)(nil)

// NewBaseWallet creates a BaseWallet by loading existing key records via the adapter.
func NewBaseWallet(adapter WalletAdapter) (*BaseWallet, error) {
	keyDir := path.Join(adapter.GetKeyDir()...)
	records, err := LoadKeyRecords(keyDir)
	if err != nil {
		return nil, err
	}

	signers := make(map[string]Signer)
	for name, record := range records {
		var signer Signer
		switch record.Type {
		case LocalSignerType:
			signer, err = adapter.LoadLocalSigner(name, record)
			if err != nil {
				return nil, err
			}
		case RemoteSignerType:
			signer, err = adapter.NewRemoteSigner(name, record.Address, record.Url, record.Key)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unsupported signer type: %s for key %s", record.Type, name)
		}
		signers[name] = signer
	}

	return &BaseWallet{Adapter: adapter, Signers: signers}, nil
}

// SaveByPrivateKey imports from a raw private key.
func (w *BaseWallet) SaveByPrivateKey(name string, privateKey string) (string, error) {
	if _, ok := w.Signers[name]; ok {
		return "", fmt.Errorf("key name exists: %s", name)
	}

	addr, signer, err := w.Adapter.DeriveFromPrivateKey(name, privateKey)
	if err != nil {
		return "", err
	}

	if w.IsAddressExist(addr) {
		return "", fmt.Errorf("address exists: %s", addr)
	}

	if err := w.Adapter.PersistKey(name, signer, privateKey); err != nil {
		return "", err
	}

	record := NewKeyRecord(LocalSignerType, addr, "", nil, "")
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

	addr, signer, saveMethod, err := w.Adapter.DeriveFromMnemonic(name, mnemonic, coinType, account, index)
	if err != nil {
		return "", err
	}

	if w.IsAddressExist(addr) {
		return "", fmt.Errorf("address exists: %s", addr)
	}

	if err := w.Adapter.PersistKey(name, signer, mnemonic); err != nil {
		return "", err
	}

	record := NewKeyRecord(LocalSignerType, addr, "", nil, saveMethod)
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

	if err := w.Adapter.ValidateAddress(address); err != nil {
		return err
	}

	if w.IsAddressExist(address) {
		return fmt.Errorf("address exists: %s", address)
	}

	record := NewKeyRecord(RemoteSignerType, address, url, key, "")
	if err := w.saveKeyRecord(name, record); err != nil {
		return err
	}

	signer, err := w.Adapter.NewRemoteSigner(name, address, url, key)
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
	for _, signer := range w.Signers {
		if w.Adapter.CompareAddresses(signer.GetAddress(), address) {
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

	return os.Write(b, append(w.Adapter.GetKeyDir(), fmt.Sprintf("%s.toml", name)))
}

// deleteKeyRecord deletes the KeyRecord file from disk.
func (w *BaseWallet) deleteKeyRecord(name string) error {
	dir := path.Join(w.Adapter.GetKeyDir()...)
	filePath := path.Join(dir, fmt.Sprintf("%s.toml", name))
	return os.DeletePath(filePath)
}
