package geth_test

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"

	internalOs "github.com/bandprotocol/falcon/internal/os"
	"github.com/bandprotocol/falcon/relayer/wallet/geth"
)

type WalletTestSuite struct {
	suite.Suite
	passphrase string
	chainName  string
}

func TestWalletTestSuite(t *testing.T) {
	suite.Run(t, new(WalletTestSuite))
}

func (s *WalletTestSuite) SetupTest() {
	s.passphrase = ""
	s.chainName = "testnet"
}

// newWallet creates a fresh GethWallet with its own temp directory.
func (s *WalletTestSuite) newWallet() (*geth.GethWallet, string) {
	home := s.T().TempDir()
	w, err := geth.NewGethWallet(s.passphrase, home, s.chainName)
	s.Require().NoError(err)
	return w, home
}

func (s *WalletTestSuite) TestSavePrivateKey() {
	priv, err := crypto.GenerateKey()
	s.Require().NoError(err)
	addrHex := crypto.PubkeyToAddress(priv.PublicKey).Hex()

	tests := []struct {
		name      string
		keyName   string
		setup     func(w *geth.GethWallet)
		wantErr   bool
		errSubstr string
	}{
		{"first import succeeds", "alice", nil, false, ""},
		{
			"duplicate name fails", "alice",
			func(w *geth.GethWallet) {
				_, err := w.SavePrivateKey("alice", priv)
				s.Require().NoError(err)
			},
			true, "key name exists",
		},
		{
			"duplicate address fails", "bob",
			func(w *geth.GethWallet) {
				_, err := w.SavePrivateKey("a", priv)
				s.Require().NoError(err)
			},
			true, "address exists",
		},
	}

	for _, tc := range tests {
		s.T().Run(tc.name, func(t *testing.T) {
			w, home := s.newWallet()
			if tc.setup != nil {
				tc.setup(w)
				// reload to pick up on-disk records
				w, _ = geth.NewGethWallet(s.passphrase, home, s.chainName)
			}

			gotAddr, err := w.SavePrivateKey(tc.keyName, priv)
			if tc.wantErr {
				s.Error(err)
				s.Contains(err.Error(), tc.errSubstr)
				return
			}
			s.NoError(err)
			s.Equal(addrHex, gotAddr)

			ksFiles, err := internalOs.ListFilePaths(path.Join(home, "keys", s.chainName, "priv"))
			s.Require().NoError(err)
			s.NotEmpty(ksFiles)

			recFiles, err := internalOs.ListFilePaths(path.Join(home, "keys", s.chainName, "signer"))
			s.Require().NoError(err)
			s.Len(recFiles, 1)
			s.Equal(tc.keyName+".toml", filepath.Base(recFiles[0]))
		})
	}
}

func (s *WalletTestSuite) TestSaveRemoteSignerKey() {
	priv, err := crypto.GenerateKey()
	s.Require().NoError(err)
	validAddr := crypto.PubkeyToAddress(priv.PublicKey).Hex()

	testKey := "testKey"

	tests := []struct {
		name      string
		keyName   string
		addr      string
		url       string
		key       *string
		setup     func(w *geth.GethWallet)
		wantErr   bool
		errSubstr string
	}{
		{"first remote succeeds", "remote1", validAddr, "http://example.com", &testKey, nil, false, ""},
		{
			"duplicate name fails", "dup", validAddr, "http://x", &testKey,
			func(w *geth.GethWallet) {
				s.Require().NoError(w.SaveRemoteSignerKey("dup", validAddr, "http://x", &testKey))
			},
			true, "key name exists",
		},
		{"invalid address fails", "bad", "not-an-addr", "url", &testKey, nil, true, "invalid address"},
		{
			"duplicate address fails", "another", validAddr, "http://y", &testKey,
			func(w *geth.GethWallet) {
				s.Require().NoError(w.SaveRemoteSignerKey("orig", validAddr, "http://orig", &testKey))
			},
			true, "address exists",
		},
	}

	for _, tc := range tests {
		s.T().Run(tc.name, func(t *testing.T) {
			w, home := s.newWallet()
			if tc.setup != nil {
				tc.setup(w)
				w, _ = geth.NewGethWallet(s.passphrase, home, s.chainName)
			}

			err := w.SaveRemoteSignerKey(tc.keyName, tc.addr, tc.url, tc.key)
			if tc.wantErr {
				s.Error(err)
				s.Contains(err.Error(), tc.errSubstr)
				return
			}
			s.NoError(err)

			signerFiles, err := internalOs.ListFilePaths(path.Join(home, "keys", s.chainName, "signer"))
			s.Require().NoError(err)
			s.Len(signerFiles, 1)
			s.Equal(tc.keyName+".toml", filepath.Base(signerFiles[0]))

			remoteFiles, err := internalOs.ListFilePaths(path.Join(home, "keys", s.chainName, "remote"))
			s.Require().NoError(err)
			s.Len(remoteFiles, 1)
			s.Equal(tc.keyName+".toml", filepath.Base(remoteFiles[0]))
		})
	}
}

func (s *WalletTestSuite) TestDeleteKey() {
	priv, err := crypto.GenerateKey()
	s.Require().NoError(err)
	addrHex := crypto.PubkeyToAddress(priv.PublicKey).Hex()

	testKey := "testKey"

	tests := []struct {
		name      string
		setup     func(w *geth.GethWallet)
		keyToDel  string
		wantErr   bool
		errSubstr string
	}{
		{
			"delete local succeeds",
			func(w *geth.GethWallet) {
				_, err := w.SavePrivateKey("alice", priv)
				s.Require().NoError(err)
			},
			"alice", false, "",
		},
		{
			"delete remote succeeds",
			func(w *geth.GethWallet) {
				s.Require().NoError(w.SaveRemoteSignerKey("bob", addrHex, "http://u", &testKey))
			},
			"bob", false, "",
		},
		{"delete non-existent fails", nil, "none", true, "does not exist"},
	}

	for _, tc := range tests {
		s.T().Run(tc.name, func(t *testing.T) {
			w, home := s.newWallet()
			if tc.setup != nil {
				tc.setup(w)
				w, _ = geth.NewGethWallet(s.passphrase, home, s.chainName)
			}

			err := w.DeleteKey(tc.keyToDel)
			if tc.wantErr {
				s.Error(err)
				s.Contains(err.Error(), tc.errSubstr)
			} else {
				s.NoError(err)

				sigs, err := internalOs.ListFilePaths(path.Join(home, "keys", s.chainName, "signer"))
				s.Require().NoError(err)
				s.Empty(sigs)

				if tc.keyToDel == "alice" {
					entries, err := os.ReadDir(path.Join(home, "keys", s.chainName, "priv"))
					s.Require().NoError(err)
					s.Empty(entries)
				} else {
					rem, err := internalOs.ListFilePaths(path.Join(home, "keys", s.chainName, "remote"))
					s.Require().NoError(err)
					s.Empty(rem)
				}
			}
		})
	}
}
