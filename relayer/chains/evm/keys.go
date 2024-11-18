package evm

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"os"
	"path"

	"github.com/ethereum/go-ethereum/accounts"
	keyStore "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/pelletier/go-toml/v2"

	chainstypes "github.com/bandprotocol/falcon/relayer/chains/types"
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
	passphrase    = ""
)

func (cp *EVMChainProvider) AddKey(
	keyName string,
	mnemonic string,
	privateKey string,
	homePath string,
	coinType uint32,
	account uint,
	index uint,
) (*chainstypes.Key, error) {
	if privateKey != "" {
		return cp.AddKeyWithPrivateKey(keyName, privateKey, homePath)
	}
	return cp.AddKeyWithMnemonic(keyName, mnemonic, homePath, coinType, account, index)
}

// AddKeyWithMnemonic adds a key using a mnemonic phrase.
func (cp *EVMChainProvider) AddKeyWithMnemonic(
	keyName string,
	mnemonic string,
	homePath string,
	coinType uint32,
	account uint,
	index uint,
) (*chainstypes.Key, error) {
	var err error

	// Generate mnemonic if not provided
	if mnemonic == "" {
		mnemonic, err = hdwallet.NewMnemonic(mnemonicSize)
		if err != nil {
			return nil, err
		}
	}

	// Generate private key using mnemonic
	priv, err := cp.generatePrivateKey(mnemonic, coinType, account, index)
	if err != nil {
		return nil, err
	}

	return cp.finalizeKeyAddition(keyName, priv, mnemonic, homePath)
}

// AddKeyWithPrivateKey adds a key using a raw private key.
func (cp *EVMChainProvider) AddKeyWithPrivateKey(
	keyName string,
	privateKey string,
	homePath string,
) (*chainstypes.Key, error) {
	// Convert private key from hex
	priv, err := crypto.HexToECDSA(ConvertPrivateKeyStrToHex(privateKey))
	if err != nil {
		return nil, err
	}

	// No mnemonic is used, so pass an empty string
	return cp.finalizeKeyAddition(keyName, priv, "", homePath)
}

// finalizeKeyAddition stores the private key and initializes the sender.
func (cp *EVMChainProvider) finalizeKeyAddition(
	keyName string,
	priv *ecdsa.PrivateKey,
	mnemonic string,
	homePath string,
) (*chainstypes.Key, error) {
	// Get public key from private key
	publicKeyECDSA, ok := priv.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("cannot assert type to *ecdsa.PublicKey")
	}

	// Store private key and get account info
	_, err := cp.storePrivateKey(priv)
	if err != nil {
		return nil, err
	}

	addressHex := crypto.PubkeyToAddress(*publicKeyECDSA).String()

	// Store key info and finalize
	cp.KeyInfo[keyName] = addressHex
	if err := cp.storeKeyInfo(homePath); err != nil {
		return nil, err
	}

	return chainstypes.NewKey(mnemonic, addressHex, ""), nil
}

// IsKeyNameExist checks whether the given key name is already in use.
func (cp *EVMChainProvider) IsKeyNameExist(keyName string) bool {
	_, ok := cp.KeyInfo[keyName]
	return ok
}

// DeleteKey deletes the given key name from the key store and removes its information.
func (cp *EVMChainProvider) DeleteKey(homePath, keyName string) error {
	address, err := HexToAddress(cp.KeyInfo[keyName])
	if err != nil {
		return err
	}
	if err := cp.KeyStore.Delete(accounts.Account{Address: address}, passphrase); err != nil {
		return err
	}

	delete(cp.KeyInfo, keyName)

	return cp.storeKeyInfo(homePath)
}

// ExportPrivateKey exports private key of given key name.
func (cp *EVMChainProvider) ExportPrivateKey(keyName string) (string, error) {
	key, err := cp.getKeyFromKeyName(keyName)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(crypto.FromECDSA(key.PrivateKey)), nil
}

// Listkeys lists all keys.
func (cp *EVMChainProvider) Listkeys() []*chainstypes.Key {
	res := make([]*chainstypes.Key, 0, len(cp.KeyInfo))
	for keyName, address := range cp.KeyInfo {
		key := chainstypes.NewKey("", address, keyName)
		res = append(res, key)
	}
	return res
}

// ShowKey shows key by the given name.
func (cp *EVMChainProvider) ShowKey(keyName string) string {
	return cp.KeyInfo[keyName]
}

// storePrivateKey stores private key to keyStore.
func (cp *EVMChainProvider) storePrivateKey(
	priv *ecdsa.PrivateKey,
) (*accounts.Account, error) {
	accs, err := cp.KeyStore.ImportECDSA(priv, passphrase)
	if err != nil {
		return nil, err
	}
	return &accs, nil
}

// storeKeyInfo stores key information.
func (cp *EVMChainProvider) storeKeyInfo(homePath string) error {
	b, err := toml.Marshal(cp.KeyInfo)
	if err != nil {
		return err
	}

	keyInfoDir := path.Join(homePath, keyDir, cp.ChainName, infoDir)
	keyInfoPath := path.Join(keyInfoDir, infoFileName)
	// Create the info folder if doesn't exist
	if _, err := os.Stat(keyInfoDir); os.IsNotExist(err) {
		if err = os.Mkdir(keyInfoDir, os.ModePerm); err != nil {
			return err
		}
	}
	// Create the file and write the default config to the given location.
	f, err := os.Create(keyInfoPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.Write(b); err != nil {
		return err
	}

	return nil
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

func (cp *EVMChainProvider) getKeyFromKeyName(
	keyName string,
) (*keyStore.Key, error) {
	address, err := HexToAddress(cp.KeyInfo[keyName])
	if err != nil {
		return nil, err
	}

	accs, err := cp.KeyStore.Find(accounts.Account{Address: address})
	if err != nil {
		return nil, err
	}
	b, err := cp.KeyStore.Export(accs, passphrase, passphrase)
	if err != nil {
		return nil, err
	}
	key, err := keyStore.DecryptKey(b, passphrase)
	if err != nil {
		return nil, err
	}
	return key, nil
}
