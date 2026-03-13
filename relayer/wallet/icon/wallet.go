package icon

import "github.com/bandprotocol/falcon/relayer/wallet"

// NewWallet creates a new wallet for the ICON chain.
// ICON is remote-only — passphrase is unused but kept for a consistent signature.
func NewWallet(passphrase, homePath, chainName string) (*wallet.BaseWallet, error) {
	return wallet.NewBaseWallet(passphrase, homePath, chainName, wallet.NewRemoteOnlyAdapter(NewRemoteSigner))
}
