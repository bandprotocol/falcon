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
	// hd path template for eth coin
	hdPathTemplate = "m/44'/%d'/%d'/0/%d"

	// mnemonic size is 256 bits
	mnemonicSize = 256
)

const (
	keyDir        = "keys"
	infoDir       = "info"
	privateKeyDir = "priv"
	infoFile      = "info.toml"
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
	var err error
	var priv *ecdsa.PrivateKey
	if privateKey == "" {
		if mnemonic == "" {
			mnemonic, err = hdwallet.NewMnemonic(mnemonicSize)
			if err != nil {
				return nil, err
			}
		}
		priv, err = cp.generatePrivateKey(mnemonic, coinType, account, index)
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

	accs, err := cp.storePrivateKey(priv)
	if err != nil {
		return nil, err
	}

	sender := make(chan *Sender, 1)
	sender <- NewSender(priv, accs.Address)

	cp.FreeSenders[keyName] = sender

	if err := cp.storeKeyInfo(homePath); err != nil {
		return nil, err
	}

	return chainstypes.NewKey(mnemonic, crypto.PubkeyToAddress(*publicKeyECDSA).String(), ""), nil
}

// IsKeyNameExist checks whether the given key name is already in use.
func (cp *EVMChainProvider) IsKeyNameExist(keyName string) bool {
	_, ok := cp.FreeSenders[keyName]
	return ok
}

// DeleteKey deletes the given key name from the key store and removes its information.
func (cp *EVMChainProvider) DeleteKey(homePath, keyName string) error {
	senderChannel := cp.FreeSenders[keyName]

	select {
	case sender := <-senderChannel:
		err := cp.KeyStore.Delete(accounts.Account{Address: sender.Address}, passphrase)
		if err != nil {
			return err
		}

		delete(cp.FreeSenders, keyName)

		return cp.storeKeyInfo(homePath)
	default:
		return fmt.Errorf("unavailable key name %s", keyName)
	}
}

// ExportPrivateKey exports private key of given key name.
func (cp *EVMChainProvider) ExportPrivateKey(keyName string) (string, error) {
	senderChannel := cp.FreeSenders[keyName]
	select {
	case sender := <-senderChannel:
		privateKeyByte := crypto.FromECDSA(sender.PrivateKey)
		privateKeyHex := hex.EncodeToString(privateKeyByte)
		return privateKeyHex, nil
	default:
		return "", fmt.Errorf("unavailable key name %s", keyName)
	}
}

// Listkeys lists all keys.
func (cp *EVMChainProvider) Listkeys() []*chainstypes.Key {
	res := make([]*chainstypes.Key, 0, len(cp.FreeSenders))
	for keyName, senderChannel := range cp.FreeSenders {
		select {
		case sender := <-senderChannel:
			address := sender.Address.Hex()
			key := chainstypes.NewKey("", address, keyName)
			res = append(res, key)
		default:
		}
	}
	return res
}

// Showkey shows key by the given name.
func (cp *EVMChainProvider) Showkey(keyName string) string {
	senderChannel := cp.FreeSenders[keyName]
	select {
	case sender := <-senderChannel:
		return sender.Address.Hex()
	default:
		return ""
	}
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
	keyInfo := make(KeyInfo)

	for keyName, senderChannel := range cp.FreeSenders {
		select {
		case sender := <-senderChannel:
			keyInfo[sender.Address.Hex()] = keyName
		default:
			return fmt.Errorf("unavailable key name %s", keyName)
		}
	}

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
