package evm_test

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/bandprotocol/falcon/internal/relayertest/mocks"
	fkmsv1 "github.com/bandprotocol/falcon/proto/fkms/v1"
	"github.com/bandprotocol/falcon/relayer/wallet"
	"github.com/bandprotocol/falcon/relayer/wallet/evm"
)

const (
	name    = "remote"
	address = "0x1234567890abcdef1234567890abcdef12345678"
	url     = "0.0.0.0:50051"
)

// RemoteSignerTestSuite runs tests for evm.RemoteSigner.
type RemoteSignerTestSuite struct {
	suite.Suite

	ctrl       *gomock.Controller
	mockClient *mocks.MockFkmsServiceClient
	rs         *evm.RemoteSigner
}

func TestRemoteSignerTestSuite(t *testing.T) {
	suite.Run(t, new(RemoteSignerTestSuite))
}

func (s *RemoteSignerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockClient = mocks.NewMockFkmsServiceClient(s.ctrl)

	testKey := "testKey"
	rs, err := evm.NewRemoteSigner(
		name,
		address,
		url,
		testKey,
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

func (s *RemoteSignerTestSuite) TestSign() {
	signerPayload := evm.NewSignerPayload(
		"0x1234567890abcdef1234567890abcdef12345678",
		1,
		42,
		"0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
		21000,
		[]byte{0x01}, nil, nil,
	)

	payloadJSON, err := json.Marshal(signerPayload)
	s.Require().NoError(err)

	tssPayload := wallet.TssPayload{
		TssMessage: []byte("tss-msg"),
		RandomAddr: []byte("random-addr"),
		Signature:  []byte("signature"),
	}

	expectedTxBlob := []byte("signed-evm-tx")

	s.mockClient.EXPECT().SignEvm(
		gomock.Any(),
		&fkmsv1.SignEvmRequest{
			SignerPayload: &fkmsv1.EvmSignerPayload{
				Address:  signerPayload.Address,
				ChainId:  signerPayload.ChainID,
				Nonce:    signerPayload.Nonce,
				To:       signerPayload.To,
				GasLimit: signerPayload.GasLimit,
				GasPrice: signerPayload.GasPrice,
			},
			Tss: &fkmsv1.Tss{
				Message:    tssPayload.TssMessage,
				RandomAddr: tssPayload.RandomAddr,
				SignatureS: tssPayload.Signature,
			},
		},
	).Return(&fkmsv1.SignEvmResponse{TxBlob: expectedTxBlob}, nil)

	result, err := s.rs.Sign(payloadJSON, tssPayload)
	s.Require().NoError(err)
	s.Equal(expectedTxBlob, result)
}
