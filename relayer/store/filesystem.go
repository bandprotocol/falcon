package store

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"path"

	"github.com/pelletier/go-toml/v2"

	"github.com/bandprotocol/falcon/internal/os"
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
	HomePath string

	hashedPassphrase []byte
}

// NewFileSystem creates a new filesystem store.
func NewFileSystem(homePath string) (*FileSystem, error) {
	passphrasePath := path.Join(getPassphrasePath(homePath)...)
	hashedPassphrase, err := os.ReadFileIfExist(passphrasePath)
	if err != nil {
		return nil, err
	}

	return &FileSystem{
		HomePath:         homePath,
		hashedPassphrase: hashedPassphrase,
	}, nil
}

// HasConfig checks if the config file exists in the filesystem.
func (fs *FileSystem) HasConfig() (bool, error) {
	cfgPath := path.Join(getConfigPath(fs.HomePath)...)
	return os.IsPathExist(cfgPath)
}

// GetConfig reads the config file from the filesystem and returns the config object.
func (fs *FileSystem) GetConfig() (*config.Config, error) {
	cfgPath := path.Join(getConfigPath(fs.HomePath)...)
	b, err := os.ReadFileIfExist(cfgPath)
	if err != nil {
		return nil, err
	} else if b == nil {
		return nil, nil
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

	return os.Write(b, getConfigPath(fs.HomePath))
}

// GetHashedPassphrase reads the passphrase from the filesystem and returns it.
func (fs *FileSystem) GetHashedPassphrase() ([]byte, error) {
	return fs.hashedPassphrase, nil
}

// SaveHashedPassphrase saves the hashedPassphrase to the filesystem.
func (fs *FileSystem) SaveHashedPassphrase(hashedPassphrase []byte) error {
	fs.hashedPassphrase = hashedPassphrase

	return os.Write(hashedPassphrase, getPassphrasePath(fs.HomePath))
}

// ValidatePassphrase validates the given passphrase with the stored hashed passphrase.
func (fs *FileSystem) ValidatePassphrase(passphrase string) error {
	// prepare bytes slices of hashed env passphrase
	h := sha256.New()
	h.Write([]byte(passphrase))
	hashedPassphrase := h.Sum(nil)

	// load passphrase from local disk
	storedHashedPassphrase, err := fs.GetHashedPassphrase()
	if err != nil {
		return err
	}

	if !bytes.Equal(hashedPassphrase, storedHashedPassphrase) {
		return fmt.Errorf("invalid passphrase: the provided passphrase does not match the stored hashed passphrase")
	}

	return nil
}

// NewWallet creates a new wallet object based on the chain type and chain name.
func (fs *FileSystem) NewWallet(chainType chains.ChainType, chainName, passphrase string) (wallet.Wallet, error) {
	switch chainType {
	case chains.ChainTypeEVM:
		return wallet.NewGethWallet(passphrase, fs.HomePath, chainName)
	default:
		return nil, fmt.Errorf("unsupported chain type: %s", chainType)
	}
}

// getConfigPath returns the directories of the config file and config file name.
func getConfigPath(homePath string) []string {
	return []string{homePath, cfgDir, cfgFileName}
}

// getPassphrasePath returns the directories of the passphrase file and passphrase file name.
func getPassphrasePath(homePath string) []string {
	return []string{homePath, cfgDir, passphraseFileName}
}
