package evm

import "github.com/bandprotocol/falcon/relayer/wallet"

// NewWallet creates a new wallet for the given EVM chain.
func NewWallet(passphrase, homePath, chainName string) (*wallet.BaseWallet, error) {
	return wallet.NewBaseWallet(passphrase, homePath, chainName, &Adapter{})
}
