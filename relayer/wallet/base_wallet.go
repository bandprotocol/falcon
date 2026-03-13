package wallet

import (
	"fmt"
	"path"

	toml "github.com/pelletier/go-toml/v2"

	"github.com/bandprotocol/falcon/internal/os"
)

var _ Wallet = (*BaseWallet)(nil)

// BaseWallet provides shared wallet logic for all chain types.
// It implements the full Wallet interface; chain-specific code lives only in
// the WalletAdapter passed to NewBaseWallet.
type BaseWallet struct {
	adapter WalletAdapter
	kr      *signerKeyring
	signers map[string]Signer
	keyDir  []string
}

// NewBaseWallet creates a BaseWallet for the given chain, using the per-chain keyring
// at {homePath}/keys/{chainName}/keyring and key records at {homePath}/keys/{chainName}/metadata.
func NewBaseWallet(passphrase, homePath, chainName string, adapter WalletAdapter) (*BaseWallet, error) {
	kr, err := newSignerKeyring(passphrase, homePath, chainName)
	if err != nil {
		return nil, err
	}

	keyDir := []string{homePath, "keys", chainName, "metadata"}
	records, err := LoadKeyRecords(path.Join(keyDir...))
	if err != nil {
		return nil, err
	}

	signers := make(map[string]Signer)
	for name, record := range records {
		secret, err := kr.load(name)
		if err != nil {
			return nil, err
		}
		signer, err := adapter.LoadSigner(name, record, secret)
		if err != nil {
			return nil, err
		}
		signers[name] = signer
	}

	return &BaseWallet{
		adapter: adapter,
		kr:      kr,
		signers: signers,
		keyDir:  keyDir,
	}, nil
}

// SaveByPrivateKey imports a key from a raw private key.
func (w *BaseWallet) SaveByPrivateKey(name, privateKey string) (string, error) {
	if _, ok := w.signers[name]; ok {
		return "", fmt.Errorf("key name exists: %s", name)
	}

	if privateKey == "" {
		return "", fmt.Errorf("private key is empty")
	}

	record := NewKeyRecord(PrivKeySignerType, "", "")
	signer, err := w.adapter.LoadSigner(name, record, privateKey)
	if err != nil {
		return "", err
	}

	return w.persistLocalSigner(name, record, signer, privateKey)
}

// SaveByMnemonic imports a key derived from a mnemonic phrase.
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

	secret, err := EncodeMnemonicSecret(mnemonic, coinType, account, index)
	if err != nil {
		return "", err
	}

	record := NewKeyRecord(MnemonicSignerType, "", "")
	signer, err := w.adapter.LoadSigner(name, record, secret)
	if err != nil {
		return "", err
	}

	return w.persistLocalSigner(name, record, signer, secret)
}

// persistLocalSigner stores the secret, fills in the address on the record,
// writes it to disk, and registers the signer.
func (w *BaseWallet) persistLocalSigner(name string, record KeyRecord, signer Signer, secret string) (string, error) {
	addr := signer.GetAddress()

	if w.IsAddressExist(addr) {
		return "", fmt.Errorf("address exists: %s", addr)
	}

	if err := w.kr.store(name, secret); err != nil {
		return "", err
	}

	record.Address = addr
	if err := w.saveKeyRecord(name, record); err != nil {
		// Roll back the keyring write so no orphaned entry is left.
		_ = w.kr.delete(name)
		return "", err
	}

	w.signers[name] = signer
	return addr, nil
}

// SaveRemoteSignerKey registers a remote signer.
// key is the optional API key used to authenticate with the KMS; it is stored
// in the shared keyring (not the metadata file). Pass an empty string if no key is required.
func (w *BaseWallet) SaveRemoteSignerKey(name, address, url string, key string) error {
	if _, ok := w.signers[name]; ok {
		return fmt.Errorf("key name exists: %s", name)
	}

	record := NewKeyRecord(RemoteSignerType, address, url)
	signer, err := w.adapter.LoadSigner(name, record, key)
	if err != nil {
		return err
	}

	_, err = w.persistLocalSigner(name, record, signer, key)
	return err
}

// DeleteKey removes the signer and its records.
func (w *BaseWallet) DeleteKey(name string) error {
	if _, ok := w.signers[name]; !ok {
		return fmt.Errorf("key name does not exist: %s", name)
	}

	// Delete the metadata file first. If this fails nothing has changed.
	// If we deleted the keyring entry first and the file delete failed, the
	// TOML record would remain but the secret would be gone, causing a startup
	// error on the next run.
	if err := w.deleteKeyRecord(name); err != nil {
		return err
	}

	if err := w.kr.delete(name); err != nil {
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
