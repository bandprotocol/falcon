package icon

import "github.com/bandprotocol/falcon/relayer/wallet"

// NewWallet creates a new wallet.BaseWallet for the ICON chain.
// ICON is remote-only — passphrase is unused but kept for a consistent signature.
func NewWallet(passphrase, homePath, chainName string) (*wallet.BaseWallet, error) {
	return wallet.NewBaseWallet(homePath, chainName, wallet.NewRemoteOnlyAdapter(NewRemoteSigner))
}
