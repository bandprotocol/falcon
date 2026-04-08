package secret

import "github.com/bandprotocol/falcon/relayer/wallet"

// NewWallet creates a wallet for the Secret chain.
// Secret signing is remote-only via fkms.SignSecret, so passphrase is unused.
func NewWallet(passphrase, homePath, chainName string) (*wallet.BaseWallet, error) {
	return wallet.NewBaseWallet(
		passphrase,
		homePath,
		chainName,
		wallet.NewRemoteOnlyAdapter(NewRemoteSigner),
	)
}
