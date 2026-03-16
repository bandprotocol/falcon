package xrpl_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
	"time"

	xrplwallet "github.com/Peersyst/xrpl-go/xrpl/wallet"
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
	walletxrpl "github.com/bandprotocol/falcon/relayer/wallet/xrpl"
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
		Fee:           "100",
		NonceInterval: time.Millisecond, // reduce wait time in tests
	}

	s.wallet.EXPECT().GetSigners().Return(nil).AnyTimes()

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
	seed := "sEdVeuhfwHB6dMxgSBccJ7ZYGyLfySa"
	w, _ := xrplwallet.FromSecret(seed)
	mockSigner := walletxrpl.NewLocalSigner("test-local-signer", &w)

	s.chainProvider.FreeSigners = make(chan wallet.Signer, 1)
	s.chainProvider.FreeSigners <- mockSigner

	// Valid FixedPointABI-encoded TSS message (from chain's encoding_tss_test.go):
	// sequence=3, prices=[{CS:BAND-USD, price=2}], createdAt=123
	rawHex := ("cba0ad5a" +
		"0000000000000000000000000000000000000000000000000000000000000020" +
		"0000000000000000000000000000000000000000000000000000000000000003" +
		"0000000000000000000000000000000000000000000000000000000000000060" +
		"000000000000000000000000000000000000000000000000000000000000007b" +
		"0000000000000000000000000000000000000000000000000000000000000001" +
		"00000000000000000000000000000000000000000043533a42414e442d555344" +
		"0000000000000000000000000000000000000000000000000000000000000002")
	tssMsg, err := hex.DecodeString(rawHex)
	s.Require().NoError(err)

	packet := &bandtypes.Packet{
		TunnelID: 1,
		Sequence: 1,
		SignalPrices: []bandtypes.SignalPrice{
			{SignalID: "CS:BAND-USD", Price: 2},
		},
		CurrentGroupSigning: bandtypes.NewSigning(
			1,
			cmbytes.HexBytes(tssMsg),
			bandtypes.NewEVMSignature(
				cmbytes.HexBytes("0xraddress"),
				cmbytes.HexBytes("0xsignature"),
			),
			tsstypes.SIGNING_STATUS_SUCCESS,
		),
	}

	// s.client expectations
	s.client.EXPECT().CheckAndConnect().Return(nil)
	s.client.EXPECT().GetAccountSequenceNumber(gomock.Any(), mockSigner.GetAddress()).Return(uint32(10), nil)
	s.client.EXPECT().BroadcastTx(gomock.Any(), gomock.Any()).Return(
		xrpl.TxResult{TxHash: "HASH", Fee: "100"}, nil,
	)

	// Execute
	err = s.chainProvider.RelayPacket(context.Background(), packet)
	s.Require().NoError(err)
}

func (s *XRPLProviderTestSuite) TestRelayPacket_ConnectionError() {
	s.client.EXPECT().CheckAndConnect().Return(fmt.Errorf("connection error"))

	err := s.chainProvider.RelayPacket(context.Background(), &bandtypes.Packet{})
	s.Require().Error(err)
	s.Contains(err.Error(), "connection error")
}

func (s *XRPLProviderTestSuite) TestQueryBalance() {
	address := "rHb9CJAW8f5rjR5juUs6K3mJtr47MS9f2"
	expectedBalance := big.NewInt(1000)
	s.client.EXPECT().GetBalance(gomock.Any(), address).Return(expectedBalance, nil)

	bal, err := s.chainProvider.QueryBalance(context.Background(), address)
	s.Require().NoError(err)
	s.Equal(expectedBalance, bal)
}
