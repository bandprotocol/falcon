package xrpl

import (
	"github.com/bandprotocol/falcon/relayer/wallet"
)

// NewWallet creates a new wallet.BaseWallet for the given XRPL chain.
func NewWallet(passphrase, homePath, chainName string) (*wallet.BaseWallet, error) {
	adapter := &Adapter{
		passphrase: passphrase,
		homePath:   homePath,
		chainName:  chainName,
	}

	return wallet.NewBaseWallet(homePath, chainName, adapter)
}
