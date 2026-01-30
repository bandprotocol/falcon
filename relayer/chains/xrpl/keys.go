package xrpl

import (
	"fmt"

	"github.com/bandprotocol/falcon/relayer/chains"
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
	var err error
	if mnemonic == "" {
		mnemonic, err = chains.GenerateMnemonic(xrplMnemonicEntropyBits)
		if err != nil {
			return nil, err
		}
		generatedMnemonic = mnemonic
	}

	addr, err := cp.Wallet.SaveByMnemonic(keyName, mnemonic)
	if err != nil {
		return nil, err
	}

	return chainstypes.NewKey(generatedMnemonic, addr, ""), nil
}
