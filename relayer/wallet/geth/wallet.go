package geth

import (
	"path"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/bandprotocol/falcon/relayer/wallet"
)

// NewWallet creates a new wallet.BaseWallet for the given EVM chain.
func NewWallet(passphrase, homePath, chainName string) (*wallet.BaseWallet, error) {
	keyStoreDir := path.Join(homePath, "keys", chainName, "priv")
	store := keystore.NewKeyStore(keyStoreDir, keystore.StandardScryptN, keystore.StandardScryptP)

	adapter := &Adapter{
		passphrase: passphrase,
		store:      store,
	}

	return wallet.NewBaseWallet(homePath, chainName, adapter)
}
