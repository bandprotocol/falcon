package evm_test

import (
	"encoding/hex"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"

	internalOs "github.com/bandprotocol/falcon/internal/os"
	"github.com/bandprotocol/falcon/relayer/wallet"
	"github.com/bandprotocol/falcon/relayer/wallet/evm"
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

// newWallet creates a fresh wallet with its own temp directory.
func (s *WalletTestSuite) newWallet() (*wallet.BaseWallet, string) {
	home := s.T().TempDir()
	w, err := evm.NewWallet(s.passphrase, home, s.chainName)
	s.Require().NoError(err)
	return w, home
}

func (s *WalletTestSuite) TestSaveBySecret() {
	priv, err := crypto.GenerateKey()
	s.Require().NoError(err)
	addrHex := crypto.PubkeyToAddress(priv.PublicKey).Hex()
	privHex := "0x" + hex.EncodeToString(crypto.FromECDSA(priv))

	tests := []struct {
		name      string
		keyName   string
		setup     func(w *wallet.BaseWallet)
		wantErr   bool
		errSubstr string
	}{
		{"first import succeeds", "alice", nil, false, ""},
		{
			"duplicate name fails", "alice",
			func(w *wallet.BaseWallet) {
				_, err := w.SaveByPrivateKey("alice", privHex)
				s.Require().NoError(err)
			},
			true, "key name exists",
		},
		{
			"duplicate address fails", "bob",
			func(w *wallet.BaseWallet) {
				_, err := w.SaveByPrivateKey("a", privHex)
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
				w, _ = evm.NewWallet(s.passphrase, home, s.chainName)
			}

			gotAddr, err := w.SaveByPrivateKey(tc.keyName, privHex)
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

			recFiles, err := internalOs.ListFilePaths(path.Join(home, "keys", s.chainName, "metadata"))
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
		setup     func(w *wallet.BaseWallet)
		wantErr   bool
		errSubstr string
	}{
		{"first remote succeeds", "remote1", validAddr, "http://example.com", &testKey, nil, false, ""},
		{
			"duplicate name fails", "dup", validAddr, "http://x", &testKey,
			func(w *wallet.BaseWallet) {
				s.Require().NoError(w.SaveRemoteSignerKey("dup", validAddr, "http://x", &testKey))
			},
			true, "key name exists",
		},
		{"invalid address fails", "bad", "not-an-addr", "url", &testKey, nil, true, "invalid EVM address"},
		{
			"duplicate address fails", "another", validAddr, "http://y", &testKey,
			func(w *wallet.BaseWallet) {
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
				w, _ = evm.NewWallet(s.passphrase, home, s.chainName)
			}

			err := w.SaveRemoteSignerKey(tc.keyName, tc.addr, tc.url, tc.key)
			if tc.wantErr {
				s.Error(err)
				s.Contains(err.Error(), tc.errSubstr)
				return
			}
			s.NoError(err)

			metaFiles, err := internalOs.ListFilePaths(path.Join(home, "keys", s.chainName, "metadata"))
			s.Require().NoError(err)
			s.Len(metaFiles, 1)
			s.Equal(tc.keyName+".toml", filepath.Base(metaFiles[0]))
		})
	}
}

func (s *WalletTestSuite) TestDeleteKey() {
	priv, err := crypto.GenerateKey()
	s.Require().NoError(err)
	addrHex := crypto.PubkeyToAddress(priv.PublicKey).Hex()
	privHex := hex.EncodeToString(crypto.FromECDSA(priv))

	testKey := "testKey"

	tests := []struct {
		name      string
		setup     func(w *wallet.BaseWallet)
		keyToDel  string
		wantErr   bool
		errSubstr string
	}{
		{
			"delete local succeeds",
			func(w *wallet.BaseWallet) {
				_, err := w.SaveByPrivateKey("alice", privHex)
				s.Require().NoError(err)
			},
			"alice", false, "",
		},
		{
			"delete remote succeeds",
			func(w *wallet.BaseWallet) {
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
				w, _ = evm.NewWallet(s.passphrase, home, s.chainName)
			}

			err := w.DeleteKey(tc.keyToDel)
			if tc.wantErr {
				s.Error(err)
				s.Contains(err.Error(), tc.errSubstr)
			} else {
				s.NoError(err)

				metaFiles, err := internalOs.ListFilePaths(path.Join(home, "keys", s.chainName, "metadata"))
				s.Require().NoError(err)
				s.Empty(metaFiles)

				if tc.keyToDel == "alice" {
					entries, err := os.ReadDir(path.Join(home, "keys", s.chainName, "priv"))
					s.Require().NoError(err)
					s.Empty(entries)
				}
			}
		})
	}
}
