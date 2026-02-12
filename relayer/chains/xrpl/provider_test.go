package xrpl_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	cmbytes "github.com/cometbft/cometbft/libs/bytes"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	tsstypes "github.com/bandprotocol/falcon/internal/bandchain/tss"
	"github.com/bandprotocol/falcon/internal/relayertest/mocks"
	"github.com/bandprotocol/falcon/relayer/alert"
	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/chains/xrpl"
	"github.com/bandprotocol/falcon/relayer/logger"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

type XRPLProviderTestSuite struct {
	suite.Suite
	ctrl *gomock.Controller

	chainProvider *xrpl.XRPLChainProvider
	client        *mocks.MockXRPLClient
	wallet        *mocks.MockWallet
	alert         alert.Alert
	log           logger.Logger
}

func TestXRPLProviderTestSuite(t *testing.T) {
	suite.Run(t, new(XRPLProviderTestSuite))
}

func (s *XRPLProviderTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.client = mocks.NewMockXRPLClient(s.ctrl)
	s.wallet = mocks.NewMockWallet(s.ctrl)
	s.alert = nil
	s.log = logger.NewZapLogWrapper(zap.Must(zap.NewDevelopment()).Sugar())

	cfg := &xrpl.XRPLChainProviderConfig{
		BaseChainProviderConfig: chains.BaseChainProviderConfig{
			ChainID:  1,
			MaxRetry: 3,
		},
		Fee:           100,
		NonceInterval: time.Millisecond, // reduce wait time in tests
	}

	cp := xrpl.NewXRPLChainProvider("xrpl-test", s.client, cfg, s.log, s.wallet, s.alert)
	s.chainProvider = cp
}

func (s *XRPLProviderTestSuite) TestInit() {
	s.client.EXPECT().Connect().Return(nil)
	s.client.EXPECT().StartLivelinessCheck(gomock.Any(), gomock.Any()).Return()

	err := s.chainProvider.Init(context.Background())
	s.Require().NoError(err)
	time.Sleep(10 * time.Millisecond) // wait for goroutine
}

func (s *XRPLProviderTestSuite) TestInit_ConnectError() {
	s.client.EXPECT().Connect().Return(fmt.Errorf("connection failed"))

	err := s.chainProvider.Init(context.Background())
	s.Require().Error(err)
	s.Equal("connection failed", err.Error())
}

func (s *XRPLProviderTestSuite) TestRelayPacket() {
	// Setup test data
	mockSigner := mocks.NewMockSigner(s.ctrl)
	mockSigner.EXPECT().GetAddress().Return("rHb9CJAW8f5rjR5juUs6K3mJtr47MS9f2").AnyTimes()

	s.chainProvider.FreeSigners = make(chan wallet.Signer, 1)
	s.chainProvider.FreeSigners <- mockSigner

	packet := &bandtypes.Packet{
		TunnelID: 1,
		Sequence: 1,
		SignalPrices: []bandtypes.SignalPrice{
			{SignalID: "CS:XRP-USD", Price: 1000000},
		},
		CurrentGroupSigning: bandtypes.NewSigning(
			1,
			cmbytes.HexBytes("0xmessage"),
			bandtypes.NewEVMSignature(
				cmbytes.HexBytes("0xraddress"),
				cmbytes.HexBytes("0xsignature"),
			),
			tsstypes.SIGNING_STATUS_SUCCESS,
		),
	}

	// s.client expectations
	s.client.EXPECT().CheckAndConnect().Return(nil)
	s.client.EXPECT().GetAccountSequenceNumber(gomock.Any(), "rHb9CJAW8f5rjR5juUs6K3mJtr47MS9f2").Return(uint32(10), nil)
	s.client.EXPECT().Autofill(gomock.Any()).Return(nil)
	s.client.EXPECT().BroadcastTx(gomock.Any(), gomock.Any()).Return(
		xrpl.TxResult{TxHash: "HASH", Fee: "100"}, nil,
	)

	// mockSigner expectations
	mockSigner.EXPECT().Sign(gomock.Any(), gomock.Any()).Return([]byte("signed_tx"), nil)

	// Execute
	err := s.chainProvider.RelayPacket(context.Background(), packet)
	s.Require().NoError(err)
}

func (s *XRPLProviderTestSuite) TestRelayPacket_ConnectionError() {
	s.client.EXPECT().CheckAndConnect().Return(fmt.Errorf("connection error"))

	err := s.chainProvider.RelayPacket(context.Background(), &bandtypes.Packet{})
	s.Require().Error(err)
	s.Contains(err.Error(), "connection error")
}

func (s *XRPLProviderTestSuite) TestQueryBalance() {
	mockSigner := mocks.NewMockSigner(s.ctrl)
	mockSigner.EXPECT().GetAddress().Return("rHb9CJAW8f5rjR5juUs6K3mJtr47MS9f2").AnyTimes()

	s.wallet.EXPECT().GetSigner("keyName").Return(mockSigner, true)

	expectedBalance := big.NewInt(1000)
	s.client.EXPECT().GetBalance(gomock.Any(), "rHb9CJAW8f5rjR5juUs6K3mJtr47MS9f2").Return(expectedBalance, nil)

	bal, err := s.chainProvider.QueryBalance(context.Background(), "keyName")
	s.Require().NoError(err)
	s.Equal(expectedBalance, bal)
}

func (s *XRPLProviderTestSuite) TestQueryBalance_KeyNotFound() {
	s.wallet.EXPECT().GetSigner("unknown").Return(nil, false)

	_, err := s.chainProvider.QueryBalance(context.Background(), "unknown")
	s.Require().Error(err)
	s.Contains(err.Error(), "key name does not exist")
}
