package flow

import "github.com/bandprotocol/falcon/relayer/wallet"

// NewWallet creates a new wallet for the Flow chain.
// Flow is remote-only — local key import is not supported.
func NewWallet(passphrase, homePath, chainName string) (*wallet.BaseWallet, error) {
	return wallet.NewBaseWallet(passphrase, homePath, chainName, wallet.NewRemoteOnlyAdapter(NewRemoteSigner))
}
