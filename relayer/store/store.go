package store

import (
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/config"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

type Store interface {
	HasConfig() (bool, error)
	GetConfig() (*config.Config, error)
	SaveConfig(cfg *config.Config) error
	GetHashedPassphrase() ([]byte, error)
	SaveHashedPassphrase(hashedPassphrase []byte) error
	NewWallet(chainType chains.ChainType, chainName, passphrase string) (wallet.Wallet, error)
}
