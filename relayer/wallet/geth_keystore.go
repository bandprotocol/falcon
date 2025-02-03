package wallet

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"path"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/pelletier/go-toml/v2"

	"github.com/bandprotocol/falcon/internal"
)

var _ Wallet = &GethKeyStoreWallet{}

type GethKeyStoreWallet struct {
	passphrase          string
	store               *keystore.KeyStore
	keyNameToHexAddress map[string]string
	keyNamePath         string
}

// NewGethKeyStoreWallet creates a new GethKeyStoreWallet instance
func NewGethKeyStoreWallet(passphrase, homeDir, chainName string) (*GethKeyStoreWallet, error) {
	// create folders if not exists
	if err := internal.CheckAndCreateFolder(homeDir); err != nil {
		return nil, err
	}

	keyDir := path.Join(homeDir, "keys")
	if err := internal.CheckAndCreateFolder(keyDir); err != nil {
		return nil, err
	}

	keyChainDir := path.Join(keyDir, chainName)
	if err := internal.CheckAndCreateFolder(keyChainDir); err != nil {
		return nil, err
	}

	keyStoreDir := path.Join(keyChainDir, "priv")
	if err := internal.CheckAndCreateFolder(keyStoreDir); err != nil {
		return nil, err
	}

	keyNameDir := path.Join(keyChainDir, "info")
	if err := internal.CheckAndCreateFolder(keyNameDir); err != nil {
		return nil, err
	}

	keyNamePath := path.Join(keyNameDir, "info.toml")

	// create keystore
	store := keystore.NewKeyStore(keyStoreDir, keystore.StandardScryptN, keystore.StandardScryptP)

	// load keyNameToHexAddress map
	keyNameToHexAddress := make(map[string]string)
	if _, err := os.Stat(keyNamePath); err != nil && !os.IsNotExist(err) {
		return nil, err
	} else if err == nil {
		b, err := os.ReadFile(keyNamePath)
		if err != nil {
			return nil, err
		}

		// unmarshal them with Config into struct
		err = toml.Unmarshal(b, &keyNameToHexAddress)
		if err != nil {
			return nil, err
		}
	}

	return &GethKeyStoreWallet{
		passphrase:          passphrase,
		store:               store,
		keyNamePath:         keyNamePath,
		keyNameToHexAddress: keyNameToHexAddress,
	}, nil
}

// SavePrivateKey saves the private key to the keystore and returns the account and update the keyNameToHexAddress map
func (w *GethKeyStoreWallet) SavePrivateKey(name string, privKey *ecdsa.PrivateKey) (string, error) {
	acc, err := w.store.ImportECDSA(privKey, w.passphrase)
	if err != nil {
		return "", err
	}

	addr := acc.Address.Hex()
	w.keyNameToHexAddress[name] = addr
	if err := w.saveKeyNameToHexAddresses(); err != nil {
		return "", err
	}

	return addr, nil
}

// DeletePrivateKey deletes the private key from the keystore and returns the address
func (w *GethKeyStoreWallet) DeletePrivateKey(name string) error {
	hexAddr, ok := w.keyNameToHexAddress[name]
	if !ok {
		return fmt.Errorf("key name does not exist: %s", name)
	}

	addr, err := HexToAddress(hexAddr)
	if err != nil {
		return err
	}

	if err := w.store.Delete(accounts.Account{Address: addr}, w.passphrase); err != nil {
		return err
	}

	delete(w.keyNameToHexAddress, name)

	return w.saveKeyNameToHexAddresses()
}

// GetAddress returns the address of the given key name
func (w *GethKeyStoreWallet) GetAddress(name string) (string, bool) {
	addr, ok := w.keyNameToHexAddress[name]
	return addr, ok
}

// GetNames returns the list of key names
func (w *GethKeyStoreWallet) GetNames() []string {
	names := make([]string, 0, len(w.keyNameToHexAddress))
	for name := range w.keyNameToHexAddress {
		names = append(names, name)
	}

	return names
}

// GetKey returns the private key and address of the given key name
func (w *GethKeyStoreWallet) GetKey(name string) (*Key, error) {
	hexAddr, ok := w.keyNameToHexAddress[name]
	if !ok {
		return nil, fmt.Errorf("key name does not exist: %s", name)
	}

	gethAddr, err := HexToAddress(hexAddr)
	if err != nil {
		return nil, err
	}

	accs, err := w.store.Find(accounts.Account{Address: gethAddr})
	if err != nil {
		return nil, err
	}

	// need to export the key due to no direct access to the private key
	b, err := w.store.Export(accs, w.passphrase, w.passphrase)
	if err != nil {
		return nil, err
	}

	gethKey, err := keystore.DecryptKey(b, w.passphrase)
	if err != nil {
		return nil, err
	}

	return &Key{
		Address:    gethAddr.Hex(),
		PrivateKey: gethKey.PrivateKey,
	}, nil
}

// saveKeyNameToHexAddresses writes the keyNameToHexAddress map to the file
func (w *GethKeyStoreWallet) saveKeyNameToHexAddresses() error {
	f, err := os.Create(w.keyNamePath)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	return encoder.Encode(w.keyNameToHexAddress)
}
