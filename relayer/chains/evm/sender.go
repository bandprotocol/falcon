package evm

import (
	"crypto/ecdsa"
	"os"
	"path"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pelletier/go-toml/v2"
)

// KeyInfo struct is the struct that represents mapping of key name -> address
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

// LoadFreeSenders loads the FreeSenders channel with sender instances.
// derived from the keys stored in the keystore located at the specified homePath.
func (cp *EVMChainProvider) LoadFreeSenders(
	homePath string,
	passphrase string,
) error {
	freeSenders := make(chan *Sender, len(cp.KeyInfo))

	for keyName := range cp.KeyInfo {
		key, err := cp.getKeyFromKeyName(keyName, passphrase)
		if err != nil {
			return err
		}

		freeSenders <- NewSender(key.PrivateKey, key.Address)
	}

	cp.FreeSenders = freeSenders
	return nil
}

// loadKeyInfo loads key information from local disk.
func LoadKeyInfo(homePath, chainName string) (KeyInfo, error) {
	keyInfo := make(KeyInfo)

	keyInfoDir := path.Join(homePath, keyDir, chainName, infoDir)
	keyInfoPath := path.Join(keyInfoDir, infoFileName)

	if _, err := os.Stat(keyInfoPath); err != nil {
		// don't return error if file doesn't exist
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
