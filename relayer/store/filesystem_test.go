package store_test

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bandprotocol/falcon/relayer/config"
	"github.com/bandprotocol/falcon/relayer/store"
)

type FileSystemTestSuite struct {
	suite.Suite
	store *store.FileSystem
}

func (s *FileSystemTestSuite) SetupTest() {
	tmpDir := s.T().TempDir()
	s.store = &store.FileSystem{HomePath: tmpDir}
}

func TestAppTestSuite(t *testing.T) {
	suite.Run(t, new(FileSystemTestSuite))
}

func (s *FileSystemTestSuite) TestGetConfig() {
	// empty config
	exist, err := s.store.HasConfig()
	s.NoError(err)
	s.False(exist)

	cfg1, err := s.store.GetConfig()
	s.NoError(err)
	s.Equal((*config.Config)(nil), cfg1)

	// create a config file
	defaultCfg := config.DefaultConfig()
	err = s.store.SaveConfig(defaultCfg)
	s.NoError(err)

	exist, err = s.store.HasConfig()
	s.NoError(err)
	s.True(exist)

	// read the config file
	cfg2, err := s.store.GetConfig()
	s.NoError(err)
	s.Equal(defaultCfg, cfg2)
}

func (s *FileSystemTestSuite) TestSaveNilConfig() {
	err := s.store.SaveConfig(nil)
	s.NoError(err)
}

func (s *FileSystemTestSuite) TestGetEmptyHashedPassphrase() {
	hashedPassphrase, err := s.store.GetHashedPassphrase()
	s.NoError(err)
	s.Equal([]byte(nil), hashedPassphrase)
}

func (s *FileSystemTestSuite) TestGetHashedPassphrase() {
	err := s.store.SaveHashedPassphrase([]byte("test"))
	s.NoError(err)

	// overwrite the passphrase shouldn't cause any error
	err = s.store.SaveHashedPassphrase([]byte("new passphrase"))
	s.NoError(err)

	// create a new store to read the new passphrase
	newStore, err := store.NewFileSystem(s.store.HomePath)
	s.NoError(err)

	hashedPassphrase, err := newStore.GetHashedPassphrase()
	s.NoError(err)
	s.Require().Equal([]byte("new passphrase"), hashedPassphrase)
}

func (s *FileSystemTestSuite) TestValidatePassphraseInvalidPassphrase() {
	// prepare bytes slices of hashed env passphrase
	h := sha256.New()
	h.Write([]byte("secret"))
	hashedPassphrase := h.Sum(nil)

	err := s.store.SaveHashedPassphrase(hashedPassphrase)
	s.NoError(err)

	testcases := []struct {
		name          string
		envPassphrase string
		err           error
	}{
		{name: "valid", envPassphrase: "secret", err: nil},
		{name: "invalid", envPassphrase: "invalid", err: fmt.Errorf("invalid passphrase")},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			err := s.store.ValidatePassphrase(tc.envPassphrase)
			if tc.err != nil {
				s.Require().ErrorContains(err, tc.err.Error())
			} else {
				s.Require().NoError(err)
			}
		})
	}
}
