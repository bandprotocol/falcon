package xrpl

import "github.com/bandprotocol/falcon/relayer/wallet"

// NewWallet creates a new wallet for the given XRPL chain.
func NewWallet(passphrase, homePath, chainName string) (*wallet.BaseWallet, error) {
	return wallet.NewBaseWallet(passphrase, homePath, chainName, &Adapter{})
}
