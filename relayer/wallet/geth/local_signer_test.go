package geth_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"

	"github.com/bandprotocol/falcon/relayer/wallet"
	"github.com/bandprotocol/falcon/relayer/wallet/geth"
)

const (
	testPrivateKeyHex = "72d4772a70645a5a5ec3fdc27afda98d2860a6f7903bff5fd45c0a23d7982121"
	testName          = "testSigner"
)

type LocalSignerTestSuite struct {
	suite.Suite
	ls *geth.LocalSigner
}

func TestLocalSignerTestSuite(t *testing.T) {
	suite.Run(t, new(LocalSignerTestSuite))
}

func (s *LocalSignerTestSuite) SetupTest() {
	priv, err := crypto.HexToECDSA(testPrivateKeyHex)
	s.Require().NoError(err)
	s.ls = geth.NewLocalSigner(testName, priv)
}

func (s *LocalSignerTestSuite) TestExportPrivateKey() {
	privHex, err := s.ls.ExportPrivateKey()
	s.Require().NoError(err)
	s.Equal(testPrivateKeyHex, privHex)
}

func (s *LocalSignerTestSuite) TestGetName() {
	s.Equal(testName, s.ls.GetName())
}

func (s *LocalSignerTestSuite) TestGetAddress() {
	priv, err := crypto.HexToECDSA(testPrivateKeyHex)
	s.Require().NoError(err)
	expected := crypto.PubkeyToAddress(priv.PublicKey).Hex()

	s.Equal(expected, s.ls.GetAddress())
}

func (s *LocalSignerTestSuite) TestSign() {
	data := []byte("hello world")

	sig, err := s.ls.Sign(data, wallet.PreSignPayload{})
	s.Require().NoError(err)

	hash := crypto.Keccak256(data)
	pubkey, err := crypto.SigToPub(hash, sig)
	s.Require().NoError(err)
	recovered := crypto.PubkeyToAddress(*pubkey).Hex()

	s.Equal(s.ls.GetAddress(), recovered)
}
