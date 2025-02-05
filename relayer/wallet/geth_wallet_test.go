package wallet_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"

	"github.com/bandprotocol/falcon/relayer/chains/evm"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

const (
	privKeyHex   = "0x72d4772a70645a5a5ec3fdc27afda98d2860a6f7903bff5fd45c0a23d7982121"
	expectedAddr = "0x990Ec0f6dFc9e8eE20dec3Ab855D03007A9dD946"
)

type GethWalletTestSuite struct {
	suite.Suite

	homePath   string
	passphrase string
	chainName  string
	wallet     *wallet.GethWallet
}

func (s *GethWalletTestSuite) SetupTest() {
	s.homePath = s.T().TempDir()

	s.passphrase = "secret"
	s.chainName = "testnet"

	var err error
	s.wallet, err = wallet.NewGethWallet(s.passphrase, s.homePath, s.chainName)
	s.Require().NoError(err)
}

func TestAppTestSuite(t *testing.T) {
	suite.Run(t, new(GethWalletTestSuite))
}

func (s *GethWalletTestSuite) TestGetKey() {
	// query key before adding
	_, err := s.wallet.GetKey("key1")
	s.Require().ErrorContains(err, "key name does not exist")

	_, ok := s.wallet.GetAddress("key1")
	s.Require().False(ok)

	names := s.wallet.GetNames()
	s.Require().Empty(names)

	// add key
	privKey, err := crypto.HexToECDSA(evm.StripPrivateKeyPrefix(privKeyHex))
	s.Require().NoError(err)

	addr, err := s.wallet.SavePrivateKey("key1", privKey)
	s.Require().NoError(err)
	s.Require().Equal(expectedAddr, addr)

	// query key after adding
	key, err := s.wallet.GetKey("key1")
	s.Require().NoError(err)
	s.Require().Equal(expectedAddr, key.Address)
	s.Require().Equal(privKey, key.PrivateKey)

	addr, ok = s.wallet.GetAddress("key1")
	s.Require().True(ok)
	s.Require().Equal(expectedAddr, addr)

	names = s.wallet.GetNames()
	s.Require().Equal([]string{"key1"}, names)

	// query key from another wallet
	anotherWallet, err := wallet.NewGethWallet(s.passphrase, s.homePath, s.chainName)
	s.Require().NoError(err)

	key, err = anotherWallet.GetKey("key1")
	s.Require().NoError(err)
	s.Require().Equal(expectedAddr, key.Address)
	s.Require().Equal(privKey, key.PrivateKey)

	addr, ok = anotherWallet.GetAddress("key1")
	s.Require().True(ok)
	s.Require().Equal(expectedAddr, addr)

	names = anotherWallet.GetNames()
	s.Require().Equal([]string{"key1"}, names)
}

func (s *GethWalletTestSuite) TestDeleteKey() {
	// delete key before adding
	err := s.wallet.DeletePrivateKey("key1")
	s.Require().ErrorContains(err, "key name does not exist")

	// add key
	privKey, err := crypto.HexToECDSA(evm.StripPrivateKeyPrefix(privKeyHex))
	s.Require().NoError(err)

	addr, err := s.wallet.SavePrivateKey("key1", privKey)
	s.Require().NoError(err)
	s.Require().Equal(expectedAddr, addr)

	// delete key after adding
	err = s.wallet.DeletePrivateKey("key1")
	s.Require().NoError(err)

	// delete key again
	err = s.wallet.DeletePrivateKey("key1")
	s.Require().ErrorContains(err, "key name does not exist")
}

func (s *GethWalletTestSuite) TestDeleteExistingKey() {
	// add key
	privKey, err := crypto.HexToECDSA(evm.StripPrivateKeyPrefix(privKeyHex))
	s.Require().NoError(err)

	addr, err := s.wallet.SavePrivateKey("key1", privKey)
	s.Require().NoError(err)
	s.Require().Equal(expectedAddr, addr)

	// delete key from another wallet
	anotherWallet, err := wallet.NewGethWallet(s.passphrase, s.homePath, s.chainName)
	s.Require().NoError(err)

	err = anotherWallet.DeletePrivateKey("key1")
	s.Require().NoError(err)

	_, err = anotherWallet.GetKey("key1")
	s.Require().ErrorContains(err, "key name does not exist")
}

func (s *GethWalletTestSuite) TestCreateDuplicatedKey() {
	// add key
	privKey, err := crypto.HexToECDSA(evm.StripPrivateKeyPrefix(privKeyHex))
	s.Require().NoError(err)

	addr, err := s.wallet.SavePrivateKey("key1", privKey)
	s.Require().NoError(err)
	s.Require().Equal(expectedAddr, addr)

	// add key with the same name
	_, err = s.wallet.SavePrivateKey("key1", privKey)
	s.Require().ErrorContains(err, "account already exists")

	// add key with the same address
	_, err = s.wallet.SavePrivateKey("key2", privKey)
	s.Require().ErrorContains(err, "account already exists")
}
