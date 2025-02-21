package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"path"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pelletier/go-toml/v2"

	"github.com/bandprotocol/falcon/internal/os"
)

var _ Wallet = &GethWallet{}

type GethWallet struct {
	Passphrase  string
	Store       *keystore.KeyStore
	KeyNameInfo map[string]string
	HomePath    string
	ChainName   string
}

// NewGethWallet creates a new GethWallet instance
func NewGethWallet(passphrase, homePath, chainName string) (*GethWallet, error) {
	// create keystore
	keyStoreDir := path.Join(getEVMKeyStoreDir(homePath, chainName)...)
	store := keystore.NewKeyStore(keyStoreDir, keystore.StandardScryptN, keystore.StandardScryptP)

	// load keyNameToHexAddress map
	keyNameInfoPath := path.Join(getEVMKeyNameInfoPath(homePath, chainName)...)
	b, err := os.ReadFileIfExist(keyNameInfoPath)
	if err != nil {
		return nil, err
	}

	// unmarshal them with Config into struct
	keyNameInfo := make(map[string]string)
	err = toml.Unmarshal(b, &keyNameInfo)
	if err != nil {
		return nil, err
	}

	return &GethWallet{
		Passphrase:  passphrase,
		Store:       store,
		HomePath:    homePath,
		KeyNameInfo: keyNameInfo,
		ChainName:   chainName,
	}, nil
}

// SavePrivateKey saves the private key to the keystore and returns the account and update the keyNameToHexAddress map
func (w *GethWallet) SavePrivateKey(name string, privKey *ecdsa.PrivateKey) (string, error) {
	acc, err := w.Store.ImportECDSA(privKey, w.Passphrase)
	if err != nil {
		return "", err
	}

	addr := acc.Address.Hex()
	w.KeyNameInfo[name] = addr
	if err := w.SaveKeyNameInfo(); err != nil {
		return "", err
	}

	return addr, nil
}

// DeletePrivateKey deletes the private key from the keystore and returns the address
func (w *GethWallet) DeletePrivateKey(name string) error {
	hexAddr, ok := w.KeyNameInfo[name]
	if !ok {
		return fmt.Errorf("key name does not exist: %s", name)
	}

	addr, err := HexToETHAddress(hexAddr)
	if err != nil {
		return err
	}

	if err := w.Store.Delete(accounts.Account{Address: addr}, w.Passphrase); err != nil {
		return err
	}

	delete(w.KeyNameInfo, name)

	return w.SaveKeyNameInfo()
}

// GetAddress returns the address of the given key name
func (w *GethWallet) GetAddress(name string) (string, bool) {
	addr, ok := w.KeyNameInfo[name]
	return addr, ok
}

// GetNames returns the list of key names
func (w *GethWallet) GetNames() []string {
	names := make([]string, 0, len(w.KeyNameInfo))
	for name := range w.KeyNameInfo {
		names = append(names, name)
	}

	return names
}

// SaveKeyNameInfo writes the keyNameInfo map to the file
func (w *GethWallet) SaveKeyNameInfo() error {
	b, err := toml.Marshal(w.KeyNameInfo)
	if err != nil {
		return err
	}

	return os.Write(b, getEVMKeyNameInfoPath(w.HomePath, w.ChainName))
}

// ExportPrivateKey exports the private key of the given key name
func (w *GethWallet) ExportPrivateKey(name string) (string, error) {
	privKey, err := w.getPrivateKey(name)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(crypto.FromECDSA(privKey)), nil
}

// Sign signs the data with the private key of the given key name
func (w *GethWallet) Sign(name string, data []byte) ([]byte, error) {
	privKey, err := w.getPrivateKey(name)
	if err != nil {
		return nil, err
	}

	return crypto.Sign(data, privKey)
}

// GetPrivateKey returns the private key of the given key name
func (w *GethWallet) getPrivateKey(name string) (*ecdsa.PrivateKey, error) {
	hexAddr, ok := w.KeyNameInfo[name]
	if !ok {
		return nil, fmt.Errorf("key name does not exist: %s", name)
	}

	gethAddr, err := HexToETHAddress(hexAddr)
	if err != nil {
		return nil, err
	}

	accs, err := w.Store.Find(accounts.Account{Address: gethAddr})
	if err != nil {
		return nil, err
	}

	// need to export the key due to no direct access to the private key
	b, err := w.Store.Export(accs, w.Passphrase, w.Passphrase)
	if err != nil {
		return nil, err
	}

	gethKey, err := keystore.DecryptKey(b, w.Passphrase)
	if err != nil {
		return nil, err
	}

	return gethKey.PrivateKey, nil
}
