package wallet

import (
	"fmt"
	"path"

	toml "github.com/pelletier/go-toml/v2"

	"github.com/bandprotocol/falcon/internal/os"
)

// BaseWallet provides shared wallet logic for all chain types.
// It delegates chain-specific operations to the embedded WalletAdapter.
type BaseWallet struct {
	adapter WalletAdapter
	signers map[string]Signer
	keyDir  []string
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

	return &BaseWallet{adapter: adapter, keyDir: keyDir, signers: signers}, nil
}

// SaveByPrivateKey imports from a raw private key.
func (w *BaseWallet) SaveByPrivateKey(name string, privateKey string) (string, error) {
	if _, ok := w.signers[name]; ok {
		return "", fmt.Errorf("key name exists: %s", name)
	}

	signer, err := w.adapter.DeriveFromPrivateKey(name, privateKey)
	if err != nil {
		return "", err
	}

	return w.saveLocalSigner(name, signer, privateKey)
}

// SaveByMnemonic imports from a mnemonic phrase.
func (w *BaseWallet) SaveByMnemonic(
	name string,
	mnemonic string,
	coinType uint32,
	account uint,
	index uint,
) (string, error) {
	if _, ok := w.signers[name]; ok {
		return "", fmt.Errorf("key name exists: %s", name)
	}

	if mnemonic == "" {
		return "", fmt.Errorf("mnemonic is empty")
	}

	signer, err := w.adapter.DeriveFromMnemonic(name, mnemonic, coinType, account, index)
	if err != nil {
		return "", err
	}

	return w.saveLocalSigner(name, signer, mnemonic)
}

// saveLocalSigner persists a derived local signer and registers it.
func (w *BaseWallet) saveLocalSigner(name string, signer Signer, secret string) (string, error) {
	addr := signer.GetAddress()

	if w.IsAddressExist(addr) {
		return "", fmt.Errorf("address exists: %s", addr)
	}

	if err := w.adapter.PersistKey(name, signer, secret); err != nil {
		return "", err
	}

	record := NewKeyRecord(LocalSignerType, addr, "", nil)
	if err := w.saveKeyRecord(name, record); err != nil {
		return "", err
	}

	w.signers[name] = signer
	return addr, nil
}

// SaveRemoteSignerKey registers a remote signer.
func (w *BaseWallet) SaveRemoteSignerKey(name, address, url string, key *string) error {
	if _, ok := w.signers[name]; ok {
		return fmt.Errorf("key name exists: %s", name)
	}

	if w.IsAddressExist(address) {
		return fmt.Errorf("address exists: %s", address)
	}

	record := NewKeyRecord(RemoteSignerType, address, url, key)
	if err := w.saveKeyRecord(name, record); err != nil {
		return err
	}

	signer, err := w.adapter.LoadSigner(name, record)
	if err != nil {
		return err
	}

	w.signers[name] = signer
	return nil
}

// DeleteKey removes the signer and its records.
func (w *BaseWallet) DeleteKey(name string) error {
	signer, ok := w.signers[name]
	if !ok {
		return fmt.Errorf("key name does not exist: %s", name)
	}

	if err := w.adapter.DeleteLocalSecret(name, signer); err != nil {
		return err
	}

	if err := w.deleteKeyRecord(name); err != nil {
		return err
	}

	delete(w.signers, name)
	return nil
}

// GetSigners lists all signers.
func (w *BaseWallet) GetSigners() []Signer {
	signers := make([]Signer, 0, len(w.signers))
	for _, signer := range w.signers {
		signers = append(signers, signer)
	}

	return signers
}

// GetSigner returns the signer with the given name and a flag indicating if it was found.
func (w *BaseWallet) GetSigner(name string) (Signer, bool) {
	signer, ok := w.signers[name]
	return signer, ok
}

// IsAddressExist returns true if the given address is already registered.
func (w *BaseWallet) IsAddressExist(address string) bool {
	for _, signer := range w.signers {
		if signer.GetAddress() == address {
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
