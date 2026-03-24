package flow_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	cmbytes "github.com/cometbft/cometbft/libs/bytes"
	flowsdk "github.com/onflow/flow-go-sdk"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	tsstypes "github.com/bandprotocol/falcon/internal/bandchain/tss"
	"github.com/bandprotocol/falcon/internal/relayertest/mocks"
	"github.com/bandprotocol/falcon/relayer/alert"
	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/chains/flow"
	"github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/logger"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

const testSignerAddress = "0x1234567890abcdef"

type FlowProviderTestSuite struct {
	suite.Suite
	ctrl *gomock.Controller

	chainProvider *flow.FlowChainProvider
	client        *mocks.MockFlowClient
	wallet        *mocks.MockWallet
	alert         alert.Alert
	log           logger.Logger
}

func TestFlowProviderTestSuite(t *testing.T) {
	suite.Run(t, new(FlowProviderTestSuite))
}

func (s *FlowProviderTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.client = mocks.NewMockFlowClient(s.ctrl)
	s.wallet = mocks.NewMockWallet(s.ctrl)
	s.alert = nil
	s.log = logger.NewZapLogWrapper(zap.Must(zap.NewDevelopment()).Sugar())

	cfg := &flow.FlowChainProviderConfig{
		BaseChainProviderConfig: chains.BaseChainProviderConfig{
			MaxRetry: 3,
		},
		ComputeLimit:       1000,
		WaitingTxDuration:  time.Second,
		CheckingTxInterval: time.Millisecond,
	}

	s.wallet.EXPECT().GetSigners().Return(nil) // consumed by LoadSigners in NewFlowChainProvider

	cp, err := flow.NewFlowChainProvider("flow-test", s.client, cfg, s.log, s.wallet, s.alert)
	s.Require().NoError(err)
	s.chainProvider = cp
}

// newMockSigner creates a new mock signer with GetAddress returning testSignerAddress.
func (s *FlowProviderTestSuite) newMockSigner() *mocks.MockSigner {
	signer := mocks.NewMockSigner(s.ctrl)
	signer.EXPECT().GetAddress().Return(testSignerAddress).AnyTimes()
	return signer
}

// newTestPacket creates a minimal Packet with a non-nil CurrentGroupSigning.
func newTestPacket() *bandtypes.Packet {
	return &bandtypes.Packet{
		TunnelID:      1,
		Sequence:      1,
		TargetAddress: "0xContractAddress",
		SignalPrices: []bandtypes.SignalPrice{
			{SignalID: "CS:BAND-USD", Price: 5000},
		},
		CurrentGroupSigning: bandtypes.NewSigning(
			1,
			cmbytes.HexBytes("deadbeef"),
			bandtypes.NewEVMSignature(
				cmbytes.HexBytes("raddress"),
				cmbytes.HexBytes("signature"),
			),
			tsstypes.SIGNING_STATUS_SUCCESS,
		),
	}
}

// newTestAccount creates a flow.Account with one key.
func newTestAccount() *flowsdk.Account {
	return &flowsdk.Account{
		Keys: []*flowsdk.AccountKey{
			{Index: 0, SequenceNumber: 10},
		},
	}
}

// --- Init ---

func (s *FlowProviderTestSuite) TestInit() {
	s.client.EXPECT().Connect(gomock.Any()).Return(nil)
	s.client.EXPECT().StartLivelinessCheck(gomock.Any(), gomock.Any())

	err := s.chainProvider.Init(context.Background())
	s.Require().NoError(err)
	time.Sleep(10 * time.Millisecond) // wait for goroutine
}

func (s *FlowProviderTestSuite) TestInit_ConnectError() {
	s.client.EXPECT().Connect(gomock.Any()).Return(fmt.Errorf("connection failed"))

	err := s.chainProvider.Init(context.Background())
	s.Require().Error(err)
	s.Contains(err.Error(), "connection failed")
}

// --- QueryTunnelInfo ---

func (s *FlowProviderTestSuite) TestQueryTunnelInfo() {
	tunnel, err := s.chainProvider.QueryTunnelInfo(context.Background(), 1, "0xContractAddress")
	s.Require().NoError(err)
	s.Require().NotNil(tunnel)
	s.Equal(uint64(1), tunnel.ID)
	s.Equal("0xContractAddress", tunnel.TargetAddress)
	s.True(tunnel.IsActive)
}

// --- RelayPacket ---

func (s *FlowProviderTestSuite) TestRelayPacket() {
	signer := s.newMockSigner()
	s.chainProvider.FreeSigners = make(chan wallet.Signer, 1)
	s.chainProvider.FreeSigners <- signer

	packet := newTestPacket()
	txBlob := []byte("signed-tx-blob")

	s.client.EXPECT().CheckAndConnect(gomock.Any()).Return(nil)
	s.client.EXPECT().GetLatestBlockID(gomock.Any()).Return("abc123", nil)
	s.client.EXPECT().GetAccount(gomock.Any(), testSignerAddress).Return(newTestAccount(), nil)
	signer.EXPECT().Sign(gomock.Any(), gomock.Any()).Return(txBlob, nil)
	s.client.EXPECT().BroadcastTx(gomock.Any(), txBlob).Return("txhash-001", nil)
	s.client.EXPECT().GetTxResult(gomock.Any(), "txhash-001").Return(&flowsdk.TransactionResult{
		Status: flowsdk.TransactionStatusSealed,
		Error:  nil,
	}, nil)

	err := s.chainProvider.RelayPacket(context.Background(), packet)
	s.Require().NoError(err)
}

func (s *FlowProviderTestSuite) TestRelayPacket_ConnectionError() {
	signer := s.newMockSigner()
	s.chainProvider.FreeSigners = make(chan wallet.Signer, 1)
	s.chainProvider.FreeSigners <- signer

	s.client.EXPECT().CheckAndConnect(gomock.Any()).Return(fmt.Errorf("connection error"))

	err := s.chainProvider.RelayPacket(context.Background(), newTestPacket())
	s.Require().Error(err)
	s.Contains(err.Error(), "connection error")
}

func (s *FlowProviderTestSuite) TestRelayPacket_GetLatestBlockIDError() {
	signer := s.newMockSigner()
	s.chainProvider.FreeSigners = make(chan wallet.Signer, 1)
	s.chainProvider.FreeSigners <- signer

	s.client.EXPECT().CheckAndConnect(gomock.Any()).Return(nil)
	s.client.EXPECT().GetLatestBlockID(gomock.Any()).Return("", fmt.Errorf("rpc error")).Times(3)

	err := s.chainProvider.RelayPacket(context.Background(), newTestPacket())
	s.Require().Error(err)
	s.Contains(err.Error(), "failed to relay packet after 3 attempts")
}

func (s *FlowProviderTestSuite) TestRelayPacket_GetAccountError() {
	signer := s.newMockSigner()
	s.chainProvider.FreeSigners = make(chan wallet.Signer, 1)
	s.chainProvider.FreeSigners <- signer

	s.client.EXPECT().CheckAndConnect(gomock.Any()).Return(nil)
	s.client.EXPECT().GetLatestBlockID(gomock.Any()).Return("abc123", nil).Times(3)
	s.client.EXPECT().GetAccount(gomock.Any(), testSignerAddress).Return(nil, fmt.Errorf("account not found")).Times(3)

	err := s.chainProvider.RelayPacket(context.Background(), newTestPacket())
	s.Require().Error(err)
	s.Contains(err.Error(), "failed to relay packet after 3 attempts")
}

func (s *FlowProviderTestSuite) TestRelayPacket_AccountHasNoKeys() {
	signer := s.newMockSigner()
	s.chainProvider.FreeSigners = make(chan wallet.Signer, 1)
	s.chainProvider.FreeSigners <- signer

	emptyAccount := &flowsdk.Account{Keys: []*flowsdk.AccountKey{}}

	s.client.EXPECT().CheckAndConnect(gomock.Any()).Return(nil)
	s.client.EXPECT().GetLatestBlockID(gomock.Any()).Return("abc123", nil).Times(3)
	s.client.EXPECT().GetAccount(gomock.Any(), testSignerAddress).Return(emptyAccount, nil).Times(3)

	err := s.chainProvider.RelayPacket(context.Background(), newTestPacket())
	s.Require().Error(err)
	s.Contains(err.Error(), "failed to relay packet after 3 attempts")
}

func (s *FlowProviderTestSuite) TestRelayPacket_SignError() {
	signer := s.newMockSigner()
	s.chainProvider.FreeSigners = make(chan wallet.Signer, 1)
	s.chainProvider.FreeSigners <- signer

	s.client.EXPECT().CheckAndConnect(gomock.Any()).Return(nil)
	s.client.EXPECT().GetLatestBlockID(gomock.Any()).Return("abc123", nil).Times(3)
	s.client.EXPECT().GetAccount(gomock.Any(), testSignerAddress).Return(newTestAccount(), nil).Times(3)
	signer.EXPECT().Sign(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("sign failed")).Times(3)

	err := s.chainProvider.RelayPacket(context.Background(), newTestPacket())
	s.Require().Error(err)
	s.Contains(err.Error(), "failed to relay packet after 3 attempts")
}

func (s *FlowProviderTestSuite) TestRelayPacket_BroadcastError() {
	signer := s.newMockSigner()
	s.chainProvider.FreeSigners = make(chan wallet.Signer, 1)
	s.chainProvider.FreeSigners <- signer

	txBlob := []byte("signed-tx-blob")

	s.client.EXPECT().CheckAndConnect(gomock.Any()).Return(nil)
	s.client.EXPECT().GetLatestBlockID(gomock.Any()).Return("abc123", nil).Times(3)
	s.client.EXPECT().GetAccount(gomock.Any(), testSignerAddress).Return(newTestAccount(), nil).Times(3)
	signer.EXPECT().Sign(gomock.Any(), gomock.Any()).Return(txBlob, nil).Times(3)
	s.client.EXPECT().BroadcastTx(gomock.Any(), txBlob).Return("", fmt.Errorf("broadcast failed")).Times(3)

	err := s.chainProvider.RelayPacket(context.Background(), newTestPacket())
	s.Require().Error(err)
	s.Contains(err.Error(), "failed to relay packet after 3 attempts")
}

func (s *FlowProviderTestSuite) TestRelayPacket_TxFailed() {
	signer := s.newMockSigner()
	s.chainProvider.FreeSigners = make(chan wallet.Signer, 1)
	s.chainProvider.FreeSigners <- signer

	txBlob := []byte("signed-tx-blob")

	s.client.EXPECT().CheckAndConnect(gomock.Any()).Return(nil)
	s.client.EXPECT().GetLatestBlockID(gomock.Any()).Return("abc123", nil).Times(3)
	s.client.EXPECT().GetAccount(gomock.Any(), testSignerAddress).Return(newTestAccount(), nil).Times(3)
	signer.EXPECT().Sign(gomock.Any(), gomock.Any()).Return(txBlob, nil).Times(3)
	s.client.EXPECT().BroadcastTx(gomock.Any(), txBlob).Return("txhash-fail", nil).Times(3)
	s.client.EXPECT().GetTxResult(gomock.Any(), "txhash-fail").Return(&flowsdk.TransactionResult{
		Status: flowsdk.TransactionStatusSealed,
		Error:  fmt.Errorf("execution reverted"),
	}, nil).Times(3)

	err := s.chainProvider.RelayPacket(context.Background(), newTestPacket())
	s.Require().Error(err)
	s.Contains(err.Error(), "failed to relay packet after 3 attempts")
}

func (s *FlowProviderTestSuite) TestRelayPacket_MissingSigning() {
	signer := s.newMockSigner()
	s.chainProvider.FreeSigners = make(chan wallet.Signer, 1)
	s.chainProvider.FreeSigners <- signer

	packet := &bandtypes.Packet{
		TunnelID:      1,
		Sequence:      1,
		TargetAddress: "0xContractAddress",
		// No signing set
	}

	txBlob := []byte("signed-tx-blob")
	_ = txBlob

	s.client.EXPECT().CheckAndConnect(gomock.Any()).Return(nil)
	s.client.EXPECT().GetLatestBlockID(gomock.Any()).Return("abc123", nil).Times(3)
	s.client.EXPECT().GetAccount(gomock.Any(), testSignerAddress).Return(newTestAccount(), nil).Times(3)

	err := s.chainProvider.RelayPacket(context.Background(), packet)
	s.Require().Error(err)
	s.Contains(err.Error(), "failed to relay packet after 3 attempts")
}

// --- QueryBalance ---

func (s *FlowProviderTestSuite) TestQueryBalance() {
	address := testSignerAddress
	expected := big.NewInt(1_000_000_000)
	s.client.EXPECT().GetBalance(gomock.Any(), address).Return(expected, nil)

	bal, err := s.chainProvider.QueryBalance(context.Background(), address)
	s.Require().NoError(err)
	s.Equal(expected, bal)
}

func (s *FlowProviderTestSuite) TestQueryBalance_Error() {
	s.client.EXPECT().GetBalance(gomock.Any(), testSignerAddress).Return(nil, fmt.Errorf("node unavailable"))

	bal, err := s.chainProvider.QueryBalance(context.Background(), testSignerAddress)
	s.Require().Error(err)
	s.Nil(bal)
}

// --- GetChainName / ChainType ---

func (s *FlowProviderTestSuite) TestGetChainName() {
	s.Equal("flow-test", s.chainProvider.GetChainName())
}

func (s *FlowProviderTestSuite) TestChainType() {
	s.Equal(types.ChainTypeFlow, s.chainProvider.ChainType())
}
