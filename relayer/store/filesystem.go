package store

import (
	"fmt"
	"os"
	"path"

	"github.com/pelletier/go-toml/v2"

	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/config"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

var _ Store = &FileSystem{}

const (
	cfgDir             = "config"
	cfgFileName        = "config.toml"
	passphraseFileName = "passphrase.hash"
)

type FileSystem struct {
	HomePath   string
	Passphrase string
}

// NewFileSystem creates a new filesystem store.
func NewFileSystem(homePath, passphrase string) *FileSystem {
	return &FileSystem{
		HomePath:   homePath,
		Passphrase: passphrase,
	}
}

// HasConfig checks if the config file exists in the filesystem.
func (fs *FileSystem) HasConfig() (bool, error) {
	cfgPath := fs.getConfigPath()

	// check if file doesn't exist, exit the function as the config may not be initialized.
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

// GetConfig reads the config file from the filesystem and returns the config object.
func (fs *FileSystem) GetConfig() (*config.Config, error) {
	cfgPath := fs.getConfigPath()

	if ok, err := fs.HasConfig(); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}

	b, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	return config.ParseConfig(b)
}

// SaveConfig saves the given config object to the filesystem.
func (fs *FileSystem) SaveConfig(cfg *config.Config) error {
	// Marshal config object into bytes
	b, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}

	// Create the home and config folder if doesn't exist
	if err := checkAndCreateFolder(fs.HomePath); err != nil {
		return err
	}
	if err := checkAndCreateFolder(path.Join(fs.HomePath, cfgDir)); err != nil {
		return err
	}

	// Create the file and write the config to the given location.
	f, err := os.Create(fs.getConfigPath())
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.Write(b); err != nil {
		return err
	}

	return nil
}

// GetPassphrase reads the passphrase from the filesystem and returns it.
func (fs *FileSystem) GetPassphrase() ([]byte, error) {
	return os.ReadFile(fs.getPassphrasePath())
}

// SavePassphrase saves the given passphrase to the filesystem.
func (fs *FileSystem) SavePassphrase(passphrase []byte) error {
	// Create the home and config folder if doesn't exist
	if err := checkAndCreateFolder(fs.HomePath); err != nil {
		return err
	}
	if err := checkAndCreateFolder(path.Join(fs.HomePath, cfgDir)); err != nil {
		return err
	}

	// Create the file and write the passphrase to the given location.
	f, err := os.Create(fs.getPassphrasePath())
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.Write(passphrase); err != nil {
		return err
	}

	return nil
}

// getConfigPath returns the path to the config file.
func (fs *FileSystem) getConfigPath() string {
	return path.Join(fs.HomePath, cfgDir, cfgFileName)
}

// getPassphrasePath returns the path to the passphrase file.
func (fs *FileSystem) getPassphrasePath() string {
	return path.Join(fs.HomePath, cfgDir, passphraseFileName)
}

func (fs *FileSystem) NewWallet(chainType chains.ChainType, chainName string) (wallet.Wallet, error) {
	switch chainType {
	case chains.ChainTypeEVM:
		return wallet.NewGethKeyStoreWallet(fs.Passphrase, fs.HomePath, chainName)
	default:
		return nil, fmt.Errorf("unsupported chain type: %s", chainType)
	}
}

// checkAndCreateFolder checks if the folder exists and creates it if it doesn't.
func checkAndCreateFolder(path string) error {
	// If the folder exists and no error, return nil
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}

	// If the folder does not exist, create it.
	if os.IsNotExist(err) {
		return os.Mkdir(path, os.ModePerm)
	}

	return err
}
