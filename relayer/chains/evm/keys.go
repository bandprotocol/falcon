package evm

import (
	"fmt"

	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"

	chainstypes "github.com/bandprotocol/falcon/relayer/chains/types"
)

const (
	// hd path template for eth coin
	hdPathTemplate = "m/44'/%d'/%d'/0/%d"

	// mnemonic size is 256 bits
	mnemonicSize = 256
)

// AddKeyByMnemonic adds a key using a mnemonic phrase.
func (cp *EVMChainProvider) AddKeyByMnemonic(
	keyName string,
	mnemonic string,
	coinType uint32,
	account uint,
	index uint,
) (*chainstypes.Key, error) {
	var err error
	generatedMnemonic := ""
	if mnemonic == "" {
		mnemonic, err = hdwallet.NewMnemonic(mnemonicSize)
		if err != nil {
			return nil, err
		}
		generatedMnemonic = mnemonic
	}

	// Generate private key using mnemonic
	privHex, err := generatePrivateKeyHex(mnemonic, coinType, account, index)
	if err != nil {
		return nil, err
	}

	addr, err := cp.Wallet.SavePrivateKey(keyName, privHex)
	if err != nil {
		return nil, err
	}

	return chainstypes.NewKey(generatedMnemonic, addr, ""), nil
}

// generatePrivateKeyHex generates private key hex from given mnemonic.
func generatePrivateKeyHex(
	mnemonic string,
	coinType uint32,
	account uint,
	index uint,
) (string, error) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return "", err
	}

	hdPath := fmt.Sprintf(hdPathTemplate, coinType, account, index)
	path := hdwallet.MustParseDerivationPath(hdPath)
	accs, err := wallet.Derive(path, true)
	if err != nil {
		return "", err
	}
	privatekeyHex, err := wallet.PrivateKeyHex(accs)
	if err != nil {
		return "", err
	}

	return privatekeyHex, nil
}
