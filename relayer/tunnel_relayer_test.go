package relayer_test

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

	"github.com/bandprotocol/falcon/internal/relayertest/mocks"
	"github.com/bandprotocol/falcon/relayer"
	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
	chaintypes "github.com/bandprotocol/falcon/relayer/chains/types"
)

type TunnelRelayerTestSuite struct {
	suite.Suite

	app           *relayer.App
	ctx           context.Context
	chainProvider *mocks.MockChainProvider
	client        *mocks.MockClient
	tunnelRelayer *relayer.TunnelRelayer
}

const (
	defaultTunnelID               = uint64(1)
	defaultContractAddress        = ""
	defaultCheckingPacketInterval = time.Minute
	defaultBandLatestSequence     = uint64(1)
	defaultTargetChainSequence    = uint64(0)
)

// SetupTest sets up the test suite by creating mock objects and initializing the TunnelRelayer.
func (s *TunnelRelayerTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())

	s.chainProvider = mocks.NewMockChainProvider(ctrl)
	s.client = mocks.NewMockClient(ctrl)
	s.ctx = context.Background()

	tunnelRelayer := relayer.NewTunnelRelayer(
		zap.NewNop(),
		defaultTunnelID,
		defaultContractAddress,
		defaultCheckingPacketInterval,
		s.client,
		s.chainProvider,
	)
	s.tunnelRelayer = &tunnelRelayer
}

func TestTunnelRelayerTestSuite(t *testing.T) {
	suite.Run(t, new(TunnelRelayerTestSuite))
}

// Helper function to mock GetTunnel.
func (s *TunnelRelayerTestSuite) mockGetTunnel(bandLatestSequence uint64) {
	s.client.EXPECT().GetTunnel(s.ctx, s.tunnelRelayer.TunnelID).Return(bandtypes.NewTunnel(
		s.tunnelRelayer.TunnelID,
		bandLatestSequence,
		"",
		"",
		true,
	), nil).AnyTimes()
}

// Helper function to mock QueryTunnelInfo.
func (s *TunnelRelayerTestSuite) mockQueryTunnelInfo(sequence uint64, isActive bool) {
	s.chainProvider.EXPECT().
		QueryTunnelInfo(s.ctx, s.tunnelRelayer.TunnelID, s.tunnelRelayer.ContractAddress).
		Return(&chaintypes.Tunnel{
			ID:             s.tunnelRelayer.TunnelID,
			TargetAddress:  s.tunnelRelayer.ContractAddress,
			IsActive:       isActive,
			LatestSequence: sequence,
			Balance:        big.NewInt(1),
		}, nil)
}

// Helper function to create a mock Packet.
func createMockPacket(tunnelID, sequence uint64, status string) *bandtypes.Packet {
	signalPrices := []bandtypes.SignalPrice{
		{SignalID: "signal1", Price: 100},
		{SignalID: "signal2", Price: 200},
	}
	evmSignature := bandtypes.NewEVMSignature(
		cmbytes.HexBytes("0x1234"),
		cmbytes.HexBytes("0xabcd"),
	)

	signing := bandtypes.NewSigning(
		1,
		cmbytes.HexBytes("0xdeadbeef"),
		evmSignature,
		status,
	)

	return bandtypes.NewPacket(
		tunnelID,
		sequence,
		signalPrices,
		signing,
		nil,
	)
}

func (s *TunnelRelayerTestSuite) TestCheckAndRelaySuccess() {
	s.mockGetTunnel(defaultBandLatestSequence)
	s.mockQueryTunnelInfo(defaultTargetChainSequence, true)
	s.mockQueryTunnelInfo(defaultTargetChainSequence+1, true)
	packet := createMockPacket(s.tunnelRelayer.TunnelID, defaultTargetChainSequence+1, "SIGNING_STATUS_SUCCESS")

	s.client.EXPECT().
		GetTunnelPacket(s.ctx, s.tunnelRelayer.TunnelID, defaultTargetChainSequence+1).
		Return(packet, nil)
	s.chainProvider.EXPECT().RelayPacket(s.ctx, packet).Return(nil)

	err := s.tunnelRelayer.CheckAndRelay(s.ctx)
	s.Require().NoError(err)
}

func (s *TunnelRelayerTestSuite) TestCheckAndRelayFailedGetTunnel() {
	s.client.EXPECT().GetTunnel(s.ctx, s.tunnelRelayer.TunnelID).Return(nil, fmt.Errorf("failed to get tunnel"))

	err := s.tunnelRelayer.CheckAndRelay(s.ctx)
	s.Require().ErrorContains(err, "failed to get tunnel")
}

func (s *TunnelRelayerTestSuite) TestCheckAndRelayFailedQueryTunnelInfo() {
	s.mockGetTunnel(defaultBandLatestSequence)
	s.chainProvider.EXPECT().
		QueryTunnelInfo(s.ctx, s.tunnelRelayer.TunnelID, s.tunnelRelayer.ContractAddress).
		Return(nil, fmt.Errorf("failed to query tunnel info"))

	err := s.tunnelRelayer.CheckAndRelay(s.ctx)
	s.Require().ErrorContains(err, "failed to query tunnel info")
}

func (s *TunnelRelayerTestSuite) TestCheckAndRelayTargetChainNotActive() {
	s.mockGetTunnel(defaultBandLatestSequence)
	s.mockQueryTunnelInfo(defaultTargetChainSequence, false)

	err := s.tunnelRelayer.CheckAndRelay(s.ctx)
	s.Require().NoError(err)
}

func (s *TunnelRelayerTestSuite) TestCheckAndRelayNoNewPackets() {
	s.mockGetTunnel(defaultBandLatestSequence)
	s.mockQueryTunnelInfo(defaultBandLatestSequence, true)

	err := s.tunnelRelayer.CheckAndRelay(s.ctx)
	s.Require().NoError(err)
}

func (s *TunnelRelayerTestSuite) TestCheckAndRelayFailedGetPacket() {
	s.mockGetTunnel(defaultBandLatestSequence)
	s.mockQueryTunnelInfo(defaultTargetChainSequence, true)

	s.client.EXPECT().
		GetTunnelPacket(s.ctx, s.tunnelRelayer.TunnelID, defaultTargetChainSequence+1).
		Return(nil, fmt.Errorf("failed to get packet"))

	err := s.tunnelRelayer.CheckAndRelay(s.ctx)
	s.Require().Error(err)
}

func (s *TunnelRelayerTestSuite) TestCheckAndRelaySigningStatusFallen() {
	s.mockGetTunnel(defaultBandLatestSequence)
	s.mockQueryTunnelInfo(defaultTargetChainSequence, true)

	packet := createMockPacket(s.tunnelRelayer.TunnelID, defaultTargetChainSequence+1, "SIGNING_STATUS_FALLEN")

	s.client.EXPECT().
		GetTunnelPacket(s.ctx, s.tunnelRelayer.TunnelID, defaultTargetChainSequence+1).
		Return(packet, nil)

	err := s.tunnelRelayer.CheckAndRelay(s.ctx)
	s.Require().Error(err)
}

func (s *TunnelRelayerTestSuite) TestCheckAndRelaySigningStatusWaiting() {
	s.mockGetTunnel(defaultBandLatestSequence)
	s.mockQueryTunnelInfo(defaultTargetChainSequence, true)

	packet := createMockPacket(s.tunnelRelayer.TunnelID, defaultTargetChainSequence+1, "SIGNING_STATUS_WAITING")

	s.client.EXPECT().
		GetTunnelPacket(s.ctx, s.tunnelRelayer.TunnelID, defaultTargetChainSequence+1).
		Return(packet, nil)

	err := s.tunnelRelayer.CheckAndRelay(s.ctx)
	s.Require().NoError(err)
}

func (s *TunnelRelayerTestSuite) TestCheckAndRelayFailedToRelayPacket() {
	s.mockGetTunnel(defaultBandLatestSequence)
	s.mockQueryTunnelInfo(defaultTargetChainSequence, true)

	packet := createMockPacket(s.tunnelRelayer.TunnelID, defaultTargetChainSequence+1, "SIGNING_STATUS_SUCCESS")

	s.client.EXPECT().
		GetTunnelPacket(s.ctx, s.tunnelRelayer.TunnelID, defaultTargetChainSequence+1).
		Return(packet, nil)
	s.chainProvider.EXPECT().RelayPacket(s.ctx, packet).Return(fmt.Errorf("failed to relay packet"))

	err := s.tunnelRelayer.CheckAndRelay(s.ctx)
	s.Require().ErrorContains(err, "failed to relay packet")
}
