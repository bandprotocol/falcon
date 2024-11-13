package evm

import (
	"crypto/ecdsa"
	"os"
	"path"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pelletier/go-toml/v2"
)

// KeyInfo struct is the struct that represents mapping of address -> key name
type KeyInfo map[string]string

// Sender is the struct that represents the sender of the transaction.
type Sender struct {
	PrivateKey *ecdsa.PrivateKey
	Address    gethcommon.Address
}

// NewSender creates a new sender object.
func NewSender(privateKey *ecdsa.PrivateKey, address gethcommon.Address) *Sender {
	return &Sender{
		PrivateKey: privateKey,
		Address:    address,
	}
}

// SenderChannels struct is the struct that represents mapping of key name and sender
type FreeSenders map[string]chan *Sender

// LoadFreeSenders loads all sender account from the keystore and key information from the local disk.
func LoadFreeSenders(homePath string, chainName string, keyStore *keystore.KeyStore) (FreeSenders, error) {
	// load mapping address -> key name
	keyInfo, err := loadKeyInfo(homePath, chainName)
	if err != nil {
		return nil, err
	}
	freeSenders := make(FreeSenders)

	for _, account := range keyStore.Accounts() {
		accs, err := keyStore.Export(account, passphrase, passphrase)
		if err != nil {
			return nil, err
		}
		key, err := keystore.DecryptKey(accs, passphrase)
		if err != nil {
			return nil, err
		}

		keyName := keyInfo[key.Address.Hex()]
		sender := make(chan *Sender, 1)
		sender <- NewSender(key.PrivateKey, key.Address)

		freeSenders[keyName] = sender
	}

	return freeSenders, nil
}

// loadKeyInfo loads key information from local disk.
func loadKeyInfo(homePath, chainName string) (KeyInfo, error) {
	var keyInfo KeyInfo

	keyInfoDir := path.Join(homePath, keyDir, chainName, infoDir)
	keyInfoPath := path.Join(keyInfoDir, infoFile)
	if _, err := os.Stat(keyInfoPath); err != nil {
		// don't return error if file doesn't exist
		keyInfo = make(KeyInfo)
		return keyInfo, nil
	}

	b, err := os.ReadFile(keyInfoPath)
	if err != nil {
		return nil, err
	}

	// unmarshal them with Config into struct
	err = toml.Unmarshal(b, &keyInfo)
	if err != nil {
		return nil, err
	}

	return keyInfo, nil
}
