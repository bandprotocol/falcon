package store

import (
	chainstypes "github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/config"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

type Store interface {
	HasConfig() (bool, error)
	GetConfig() (*config.Config, error)
	SaveConfig(cfg *config.Config) error
	GetHashedPassphrase() ([]byte, error)
	SavePassphrase(passphrase string) error
	ValidatePassphrase(passphrase string) error
	NewWallet(chainType chainstypes.ChainType, chainName, passphrase string) (wallet.Wallet, error)
}
