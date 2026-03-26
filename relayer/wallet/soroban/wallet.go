package soroban

import "github.com/bandprotocol/falcon/relayer/wallet"

// NewWallet creates a new wallet for the given Soroban chain.
func NewWallet(passphrase, homePath, chainName string) (*wallet.BaseWallet, error) {
	return wallet.NewBaseWallet(passphrase, homePath, chainName, wallet.NewRemoteOnlyAdapter(NewRemoteSigner))
}
