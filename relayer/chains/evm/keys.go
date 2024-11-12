package evm

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"os"
	"path"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/pelletier/go-toml/v2"

	chainstypes "github.com/bandprotocol/falcon/relayer/chains/types"
)

const (
	// amount of addresses to generate
	addressAmount = 1

	// hd path template for eth coin
	hdPathTemplate = "m/44'/69'/0'/0/0"

	// mnemonic size is 256 bits
	mnemonicSize = 256
)

const (
	keyDir        = "keys"
	infoDir       = "info"
	privateKeyDir = "priv"
	infoFile      = "info.toml"
	passphrase    = ""
)

func (cp *EVMChainProvider) AddKey(
	keyName string,
	mnemonic string,
	privateKey string,
	homePath string,
) (*chainstypes.Key, error) {
	var err error
	var priv *ecdsa.PrivateKey
	var m string
	if privateKey == "" {
		m = mnemonic
		if m == "" {
			m, err = hdwallet.NewMnemonic(mnemonicSize)
			if err != nil {
				return nil, err
			}
		}
		priv, err = cp.generatePrivateKey(m)
		if err != nil {
			return nil, err
		}
	} else {
		priv, err = crypto.HexToECDSA(ConvertPrivateKeyStrToHex(privateKey))
		if err != nil {
			return nil, err
		}
	}

	// Get the public key from the private key
	publicKey := priv.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("cannot assert type to *ecdsa.PublicKey")
	}

	if err := cp.storePrivateKey(homePath, priv, keyName); err != nil {
		return nil, err
	}

	return chainstypes.NewKey(m, crypto.PubkeyToAddress(*publicKeyECDSA).String(), ""), nil
}

// IsKeyNameExist checks whether the given key name is already in use.
func (cp *EVMChainProvider) IsKeyNameExist(keyName string) bool {
	_, ok := cp.FreeSenders.Senders[keyName]
	return ok
}

// DeleteKey deletes the given key name from the key store and removes its information.
func (cp *EVMChainProvider) DeleteKey(homePath, keyName string) error {
	sender := cp.FreeSenders.Senders[keyName]

	select {
	case availableSender := <-sender:

		err := cp.KeyStore.Delete(accounts.Account{Address: availableSender.Address}, "")
		if err != nil {
			return err
		}
		keyInfo := cp.FreeSenders.KeyInfo
		delete(keyInfo, availableSender.Address.Hex())

		return cp.storeKeyInfo(homePath, &keyInfo)
	default:
		return fmt.Errorf("sender is not avaiblble at key name %s", keyName)
	}
}

// ExportPrivateKey exports private key of given key name.
func (cp *EVMChainProvider) ExportPrivateKey(keyName string) (string, error) {
	sender := cp.FreeSenders.Senders[keyName]

	select {
	case availableSender := <-sender:
		privateKeyByte := crypto.FromECDSA(availableSender.PrivateKey)
		privateKeyHex := hex.EncodeToString(privateKeyByte)
		return privateKeyHex, nil
	default:
		return "", fmt.Errorf("sender is not avaiblble at key name %s", keyName)
	}
}

// Listkeys lists all keys.
func (cp *EVMChainProvider) Listkeys() []*chainstypes.Key {
	res := make([]*chainstypes.Key, 0, len(cp.FreeSenders.Senders))
	for keyName, sender := range cp.FreeSenders.Senders {
		select {
		case availableSender := <-sender:
			address := availableSender.Address.Hex()
			key := chainstypes.NewKey("", address, keyName)
			res = append(res, key)
		default:
		}
	}

	return res
}

// Showkey shows
func (cp *EVMChainProvider) Showkey(keyName string) string {
	sender := cp.FreeSenders.Senders[keyName]
	select {
	case availableSender := <-sender:
		return availableSender.Address.Hex()
	default:
		return ""
	}
}

// storePrivateKey stores private key to keyStore.
func (cp *EVMChainProvider) storePrivateKey(
	homePath string,
	priv *ecdsa.PrivateKey,
	keyName string,
) error {
	accs, err := cp.KeyStore.ImportECDSA(priv, passphrase)
	if err != nil {
		return err
	}

	if cp.FreeSenders.KeyInfo == nil {
		cp.FreeSenders.KeyInfo = make(KeyInfo)
	}

	cp.FreeSenders.KeyInfo[accs.Address.Hex()] = keyName

	return cp.storeKeyInfo(homePath, &cp.FreeSenders.KeyInfo)
}

// storeKeyInfo stores key information.
func (cp *EVMChainProvider) storeKeyInfo(homePath string, keyInfo *KeyInfo) error {
	b, err := toml.Marshal(keyInfo)
	if err != nil {
		return err
	}

	keyInfoDir := path.Join(homePath, keyDir, cp.ChainName, infoDir)
	keyInfoPath := path.Join(keyInfoDir, infoFile)
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
func (cp *EVMChainProvider) generatePrivateKey(mnemonic string) (*ecdsa.PrivateKey, error) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}
	path := hdwallet.MustParseDerivationPath(hdPathTemplate)

	account, err := wallet.Derive(path, true)
	if err != nil {
		return nil, err
	}
	privatekey, err := wallet.PrivateKey(account)
	if err != nil {
		return nil, err
	}
	return privatekey, nil
}
