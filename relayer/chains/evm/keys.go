package evm

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
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
	if mnemonic == "" {
		mnemonic, err = hdwallet.NewMnemonic(mnemonicSize)
		if err != nil {
			return nil, err
		}
	}

	// Generate private key using mnemonic
	priv, err := generatePrivateKey(mnemonic, coinType, account, index)
	if err != nil {
		return nil, err
	}

	return cp.finalizeKeyAddition(keyName, priv, mnemonic)
}

// AddKeyByPrivateKey adds a key using a raw private key.
func (cp *EVMChainProvider) AddKeyByPrivateKey(keyName, privateKey string) (*chainstypes.Key, error) {
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

// AddRemoteSignerKey adds a remote signer with the given name, address, and URL.
func (cp *EVMChainProvider) AddRemoteSignerKey(keyName, addr, url string) (*chainstypes.Key, error) {
	if err := cp.Wallet.SaveRemoteSignerKey(keyName, addr, url); err != nil {
		return nil, err
	}
	return chainstypes.NewKey("", addr, ""), nil
}

// DeleteKey deletes the given key name from the key store and removes its information.
func (cp *EVMChainProvider) DeleteKey(keyName string) error {
	return cp.Wallet.DeleteKey(keyName)
}

// ExportPrivateKey exports private key of given key name.
func (cp *EVMChainProvider) ExportPrivateKey(keyName string) (string, error) {
	signer, ok := cp.Wallet.GetSigner(keyName)
	if !ok {
		return "", fmt.Errorf("key name does not exist: %s", keyName)
	}

	return signer.ExportPrivateKey()
}

// ListKeys lists all keys.
func (cp *EVMChainProvider) ListKeys() []*chainstypes.Key {
	signers := cp.Wallet.GetSigners()

	res := make([]*chainstypes.Key, 0, len(signers))
	for _, signer := range signers {
		key := chainstypes.NewKey("", signer.GetAddress(), signer.GetName())
		res = append(res, key)
	}

	return res
}

// ShowKey shows key by the given name.
func (cp *EVMChainProvider) ShowKey(keyName string) (string, error) {
	signer, ok := cp.Wallet.GetSigner(keyName)
	if !ok {
		return "", fmt.Errorf("key name does not exist: %s", keyName)
	}

	return signer.GetAddress(), nil
}

// generatePrivateKey generates private key from given mnemonic.
func generatePrivateKey(
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
