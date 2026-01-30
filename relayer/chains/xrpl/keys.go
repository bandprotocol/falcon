package xrpl

import (
	"fmt"

	xrplwallet "github.com/Peersyst/xrpl-go/xrpl/wallet"
	"github.com/bsv-blockchain/go-sdk/compat/bip39"

	chainstypes "github.com/bandprotocol/falcon/relayer/chains/types"
)

const (
	xrplMnemonicEntropyBits = 256
	xrplDefaultCoinType     = 144
)

// AddKeyByMnemonic adds a key using a mnemonic phrase.
func (cp *XRPLChainProvider) AddKeyByMnemonic(
	keyName string,
	mnemonic string,
	coinType uint32,
	account uint,
	index uint,
) (*chainstypes.Key, error) {
	if coinType != xrplDefaultCoinType || account != 0 || index != 0 {
		return nil, fmt.Errorf("xrpl mnemonic derivation only supports m/44'/144'/0'/0/0")
	}

	generatedMnemonic := ""
	if mnemonic == "" {
		entropy, err := bip39.NewEntropy(xrplMnemonicEntropyBits)
		if err != nil {
			return nil, err
		}
		mnemonic, err = bip39.NewMnemonic(entropy)
		if err != nil {
			return nil, err
		}
		generatedMnemonic = mnemonic
	}

	w, err := xrplwallet.FromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	addr, err := cp.Wallet.SavePrivateKey(keyName, w.PrivateKey)
	if err != nil {
		return nil, err
	}

	return chainstypes.NewKey(generatedMnemonic, addr, ""), nil
}
