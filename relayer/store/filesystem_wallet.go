package store

import (
	"fmt"

	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

type WalletFactory interface {
	NewWallet(chainType chains.ChainType, chainName string) (wallet.Wallet, error)
}

var _ WalletFactory = &FileSystemWalletFactory{}

type FileSystemWalletFactory struct {
	HomePath   string
	Passphrase string
}

func NewFileSystemWalletFactory(homePath, passphrase string) *FileSystemWalletFactory {
	return &FileSystemWalletFactory{
		HomePath:   homePath,
		Passphrase: passphrase,
	}
}

func (fs *FileSystemWalletFactory) NewWallet(chainType chains.ChainType, chainName string) (wallet.Wallet, error) {
	switch chainType {
	case chains.ChainTypeEVM:
		return wallet.NewGethKeyStoreWallet(fs.Passphrase, fs.HomePath, chainName)
	default:
		return nil, fmt.Errorf("unsupported chain type: %s", chainType)
	}
}
