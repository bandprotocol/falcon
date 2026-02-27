package geth_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/bandprotocol/falcon/internal/relayertest/mocks"
	"github.com/bandprotocol/falcon/relayer/wallet/geth"
)

const (
	name    = "remote"
	address = "0x1234567890abcdef1234567890abcdef12345678"
	url     = "0.0.0.0:50051"
)

// RemoteSignerTestSuite runs tests for geth.RemoteSigner.
type RemoteSignerTestSuite struct {
	suite.Suite

	ctrl       *gomock.Controller
	mockClient *mocks.MockFkmsServiceClient
	rs         *geth.RemoteSigner
}

func TestRemoteSignerTestSuite(t *testing.T) {
	suite.Run(t, new(RemoteSignerTestSuite))
}

func (s *RemoteSignerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockClient = mocks.NewMockFkmsServiceClient(s.ctrl)

	testKey := "testKey"
	rs, err := geth.NewRemoteSigner(
		name,
		common.HexToAddress(address),
		url,
		&testKey,
	)
	s.Require().NoError(err)

	s.rs = rs

	s.rs.FkmsClient = s.mockClient
}

func (s *RemoteSignerTestSuite) TestExportPrivateKey() {
	_, err := s.rs.ExportPrivateKey()
	s.Error(err)
	s.Contains(err.Error(), "cannot extract private key")
}

func (s *RemoteSignerTestSuite) TestGetName() {
	s.Equal(name, s.rs.GetName())
}

func (s *RemoteSignerTestSuite) TestGetAddress() {
	got := s.rs.GetAddress()
	s.Equal(common.HexToAddress(address).Hex(), got)
}
