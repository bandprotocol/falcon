package wallet

import (
	"crypto/ecdsa"
	"fmt"
	"path"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/pelletier/go-toml/v2"

	"github.com/bandprotocol/falcon/internal"
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
	keyStoreDir := path.Join(getKeyStoreDir(homePath, chainName)...)
	store := keystore.NewKeyStore(keyStoreDir, keystore.StandardScryptN, keystore.StandardScryptP)

	// load keyNameToHexAddress map
	keyNameInfoPath := path.Join(getKeyNameInfoPath(homePath, chainName)...)
	b, err := internal.ReadFileIfExist(keyNameInfoPath)
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

	addr, err := HexToAddress(hexAddr)
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

// GetKey returns the private key and address of the given key name
func (w *GethWallet) GetKey(name string) (*Key, error) {
	hexAddr, ok := w.KeyNameInfo[name]
	if !ok {
		return nil, fmt.Errorf("key name does not exist: %s", name)
	}

	gethAddr, err := HexToAddress(hexAddr)
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

	return &Key{
		Address:    gethAddr.Hex(),
		PrivateKey: gethKey.PrivateKey,
	}, nil
}

// SaveKeyNameInfo writes the keyNameInfo map to the file
func (w *GethWallet) SaveKeyNameInfo() error {
	b, err := toml.Marshal(w.KeyNameInfo)
	if err != nil {
		return err
	}

	return internal.Write(b, getKeyNameInfoPath(w.HomePath, w.ChainName))
}

// getKeyStoreDir returns the key store directory
func getKeyStoreDir(homePath, chainName string) []string {
	return []string{homePath, "keys", chainName, "priv"}
}

func getKeyNameInfoPath(homePath, chainName string) []string {
	return []string{homePath, "keys", chainName, "info", "info.toml"}
}
