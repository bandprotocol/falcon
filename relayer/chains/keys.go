package chains

import (
	"fmt"

	"github.com/bsv-blockchain/go-sdk/compat/bip39"

	chainstypes "github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

// GenerateMnemonic creates a BIP-39 mnemonic with the requested entropy size.
func GenerateMnemonic(bitSize int) (string, error) {
	entropy, err := bip39.NewEntropy(bitSize)
	if err != nil {
		return "", err
	}

	return bip39.NewMnemonic(entropy)
}

// AddKeyByPrivateKey adds a key using a raw private key.
func AddKeyByPrivateKey(w wallet.Wallet, keyName, privateKey string) (*chainstypes.Key, error) {
	addr, err := w.SaveBySecret(keyName, privateKey)
	if err != nil {
		return nil, err
	}

	return chainstypes.NewKey("", addr, ""), nil
}

// AddRemoteSignerKey adds a remote signer with the given name, address, and URL.
func AddRemoteSignerKey(w wallet.Wallet, keyName, addr, url string, key *string) (*chainstypes.Key, error) {
	if err := w.SaveRemoteSignerKey(keyName, addr, url, key); err != nil {
		return nil, err
	}

	return chainstypes.NewKey("", addr, ""), nil
}

// DeleteKey deletes the given key name from the key store and removes its information.
func DeleteKey(w wallet.Wallet, keyName string) error {
	return w.DeleteKey(keyName)
}

// ExportPrivateKey exports private key of given key name.
func ExportPrivateKey(w wallet.Wallet, keyName string) (string, error) {
	signer, ok := w.GetSigner(keyName)
	if !ok {
		return "", fmt.Errorf("key name does not exist: %s", keyName)
	}

	return signer.ExportPrivateKey()
}

// ListKeys lists all keys.
func ListKeys(w wallet.Wallet) []*chainstypes.Key {
	signers := w.GetSigners()

	res := make([]*chainstypes.Key, 0, len(signers))
	for _, signer := range signers {
		key := chainstypes.NewKey("", signer.GetAddress(), signer.GetName())
		res = append(res, key)
	}

	return res
}

// ShowKey shows key by the given name.
func ShowKey(w wallet.Wallet, keyName string) (string, error) {
	signer, ok := w.GetSigner(keyName)
	if !ok {
		return "", fmt.Errorf("key name does not exist: %s", keyName)
	}

	return signer.GetAddress(), nil
}
