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
	GetPassphrase() ([]byte, error)
	SavePassphrase(passphrase []byte) error
	NewWallet(chainType chains.ChainType, chainName string) (wallet.Wallet, error)
}
