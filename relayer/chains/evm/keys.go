package evm

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"

	chainstypes "github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

const (
	// hd path template for eth coin
	hdPathTemplate = "m/44'/%d'/%d'/0/%d"

	// mnemonic size is 256 bits
	mnemonicSize = 256
)

const (
	keyDir        = "keys"
	infoDir       = "info"
	privateKeyDir = "priv"
	infoFileName  = "info.toml"
)

func (cp *EVMChainProvider) AddKey(
	keyName string,
	mnemonic string,
	privateKey string,
	homePath string,
	coinType uint32,
	account uint,
	index uint,
	passphrase string,
) (*chainstypes.Key, error) {
	if cp.IsKeyNameExist(keyName) {
		return nil, fmt.Errorf("duplicate key name")
	}

	if privateKey != "" {
		return cp.AddKeyWithPrivateKey(keyName, privateKey, homePath, passphrase)
	}

	var err error
	// Generate mnemonic if not provided
	if mnemonic == "" {
		mnemonic, err = hdwallet.NewMnemonic(mnemonicSize)
		if err != nil {
			return nil, err
		}
	}
	return cp.AddKeyWithMnemonic(keyName, mnemonic, homePath, coinType, account, index, passphrase)
}

// AddKeyWithMnemonic adds a key using a mnemonic phrase.
func (cp *EVMChainProvider) AddKeyWithMnemonic(
	keyName string,
	mnemonic string,
	homePath string,
	coinType uint32,
	account uint,
	index uint,
	passphrase string,
) (*chainstypes.Key, error) {
	// Generate private key using mnemonic
	priv, err := cp.generatePrivateKey(mnemonic, coinType, account, index)
	if err != nil {
		return nil, err
	}

	return cp.finalizeKeyAddition(keyName, priv, mnemonic)
}

// AddKeyWithPrivateKey adds a key using a raw private key.
func (cp *EVMChainProvider) AddKeyWithPrivateKey(
	keyName,
	privateKey,
	homePath,
	passphrase string,
) (*chainstypes.Key, error) {
	// Convert private key from hex
	priv, err := crypto.HexToECDSA(StripPrivateKeyPrefix(privateKey))
	if err != nil {
		return nil, err
	}

	// No mnemonic is used, so pass an empty string
	return cp.finalizeKeyAddition(keyName, priv, "")
}

// finalizeKeyAddition stores the private key and initializes the sender.
func (cp *EVMChainProvider) finalizeKeyAddition(
	keyName string,
	priv *ecdsa.PrivateKey,
	mnemonic string,
) (*chainstypes.Key, error) {
	addr, err := cp.Wallet.SavePrivateKey(keyName, priv)
	if err != nil {
		return nil, err
	}

	return chainstypes.NewKey(mnemonic, addr, ""), nil
}

// DeleteKey deletes the given key name from the key store and removes its information.
func (cp *EVMChainProvider) DeleteKey(homePath, keyName, passphrase string) error {
	return cp.Wallet.DeletePrivateKey(keyName)
}

// ExportPrivateKey exports private key of given key name.
func (cp *EVMChainProvider) ExportPrivateKey(keyName, passphrase string) (string, error) {
	if !cp.IsKeyNameExist(keyName) {
		return "", fmt.Errorf("key name does not exist: %s", keyName)
	}

	key, err := cp.GetKeyFromKeyName(keyName)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(crypto.FromECDSA(key.PrivateKey)), nil
}

// ListKeys lists all keys.
func (cp *EVMChainProvider) ListKeys() []*chainstypes.Key {
	keyNames := cp.Wallet.GetNames()

	res := make([]*chainstypes.Key, 0, len(keyNames))
	for _, keyName := range keyNames {
		address, _ := cp.Wallet.GetAddress(keyName)
		key := chainstypes.NewKey("", address, keyName)
		res = append(res, key)
	}

	return res
}

// ShowKey shows key by the given name.
func (cp *EVMChainProvider) ShowKey(keyName string) (string, error) {
	address, ok := cp.Wallet.GetAddress(keyName)
	if !ok {
		return "", fmt.Errorf("key name does not exist: %s", keyName)
	}

	return address, nil
}

// IsKeyNameExist checks whether the given key name is already in use.
func (cp *EVMChainProvider) IsKeyNameExist(keyName string) bool {
	_, ok := cp.Wallet.GetAddress(keyName)
	return ok
}

// generatePrivateKey generates private key from given mnemonic.
func (cp *EVMChainProvider) generatePrivateKey(
	mnemonic string,
	coinType uint32,
	account uint,
	index uint,
) (*ecdsa.PrivateKey, error) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}
	hdPath := fmt.Sprintf(hdPathTemplate, coinType, account, index)
	path := hdwallet.MustParseDerivationPath(hdPath)

	accs, err := wallet.Derive(path, true)
	if err != nil {
		return nil, err
	}
	privatekey, err := wallet.PrivateKey(accs)
	if err != nil {
		return nil, err
	}
	return privatekey, nil
}

func (cp *EVMChainProvider) GetKeyFromKeyName(keyName string) (*wallet.Key, error) {
	return cp.Wallet.GetKey(keyName)
}
