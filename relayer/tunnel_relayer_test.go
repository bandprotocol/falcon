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

// SetupTest sets up the test suite by creating a temporary directory and declare mock objects.
func (s *TunnelRelayerTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())

	// mock objects.
	s.chainProvider = mocks.NewMockChainProvider(ctrl)
	s.client = mocks.NewMockClient(ctrl)

	s.ctx = context.Background()

	log, err := zap.NewDevelopment()
	s.Require().NoError(err)

	tunnelID := uint64(1)
	contractAddress := ""
	checkingPacketInterval := time.Minute

	tunnelRelayer := relayer.NewTunnelRelayer(
		log,
		tunnelID,
		contractAddress,
		checkingPacketInterval,
		s.client,
		s.chainProvider,
	)
	s.tunnelRelayer = &tunnelRelayer
}

func TestTunnelRelayerTestSuite(t *testing.T) {
	suite.Run(t, new(TunnelRelayerTestSuite))
}

func (s *TunnelRelayerTestSuite) TestCheckAndRelay() {
	bandLatestSequence := uint64(1)
	targetChainLatestSequence := uint64(0)
	s.client.EXPECT().GetTunnel(s.ctx, s.tunnelRelayer.TunnelID).Return(bandtypes.NewTunnel(
		s.tunnelRelayer.TunnelID,
		bandLatestSequence,
		"",
		"",
		true,
	), nil)

	s.chainProvider.EXPECT().
		QueryTunnelInfo(s.ctx, s.tunnelRelayer.TunnelID, s.tunnelRelayer.ContractAddress).
		Return(&chaintypes.Tunnel{
			ID:             s.tunnelRelayer.TunnelID,
			TargetAddress:  s.tunnelRelayer.ContractAddress,
			IsActive:       true,
			LatestSequence: targetChainLatestSequence,
			Balance:        big.NewInt(1),
		}, nil)

	signalPrices := []bandtypes.SignalPrice{
		{SignalID: "signal1", Price: 100},
		{SignalID: "signal2", Price: 200},
	}

	// Create a mock EVMSignature
	evmSignature := bandtypes.NewEVMSignature(
		cmbytes.HexBytes("0x1234"),
		cmbytes.HexBytes("0xabcd"),
	)

	// Create mock signing information
	signing := bandtypes.NewSigning(
		1,
		cmbytes.HexBytes("0xdeadbeef"),
		evmSignature,
		"SIGNING_STATUS_SUCCESS",
	)

	// Create the expected Packet object
	packet := bandtypes.NewPacket(
		s.tunnelRelayer.TunnelID,
		targetChainLatestSequence+1,
		signalPrices,
		signing,
		nil,
	)
	s.client.EXPECT().GetTunnelPacket(s.ctx, s.tunnelRelayer.TunnelID, targetChainLatestSequence+1).Return(
		packet, nil,
	)

	s.chainProvider.EXPECT().RelayPacket(s.ctx, packet).Return(
		nil,
	)

	err := s.tunnelRelayer.CheckAndRelay(s.ctx)
	s.Require().NoError(err)
}

func (s *TunnelRelayerTestSuite) TestCheckAndRelayFailedGetTunnel() {
	s.client.EXPECT().GetTunnel(s.ctx, s.tunnelRelayer.TunnelID).Return(nil, fmt.Errorf("failed to get tunnel"))

	err := s.tunnelRelayer.CheckAndRelay(s.ctx)

	s.Require().ErrorContains(err, "failed to get tunnel")
}

func (s *TunnelRelayerTestSuite) TestCheckAndRelayFailedGetTargetChainTunnel() {
	bandLatestSequence := uint64(1)
	s.client.EXPECT().GetTunnel(s.ctx, s.tunnelRelayer.TunnelID).Return(bandtypes.NewTunnel(
		s.tunnelRelayer.TunnelID,
		bandLatestSequence,
		"",
		"",
		true,
	), nil)

	s.chainProvider.EXPECT().
		QueryTunnelInfo(s.ctx, s.tunnelRelayer.TunnelID, s.tunnelRelayer.ContractAddress).
		Return(nil, fmt.Errorf("failed to get target chain tunnel"))

	err := s.tunnelRelayer.CheckAndRelay(s.ctx)
	s.Require().ErrorContains(err, "failed to get target chain tunnel")
}

func (s *TunnelRelayerTestSuite) TestCheckAndRelayTargetChainNotActive() {
	bandLatestSequence := uint64(1)
	targetChainLatestSequence := uint64(0)
	s.client.EXPECT().GetTunnel(s.ctx, s.tunnelRelayer.TunnelID).Return(bandtypes.NewTunnel(
		s.tunnelRelayer.TunnelID,
		bandLatestSequence,
		"",
		"",
		true,
	), nil)

	s.chainProvider.EXPECT().
		QueryTunnelInfo(s.ctx, s.tunnelRelayer.TunnelID, s.tunnelRelayer.ContractAddress).
		Return(&chaintypes.Tunnel{
			ID:             s.tunnelRelayer.TunnelID,
			TargetAddress:  s.tunnelRelayer.ContractAddress,
			IsActive:       false,
			LatestSequence: targetChainLatestSequence,
			Balance:        big.NewInt(1),
		}, nil)

	err := s.tunnelRelayer.CheckAndRelay(s.ctx)

	s.Require().ErrorContains(err, "tunnel is not active on target chain")
}

func (s *TunnelRelayerTestSuite) TestCheckAndRelayNoNewPackets() {
	bandLatestSequence := uint64(1)
	targetChainLatestSequence := uint64(1)

	// Mock BandClient to return the same sequence as TargetChain
	s.client.EXPECT().GetTunnel(s.ctx, s.tunnelRelayer.TunnelID).Return(bandtypes.NewTunnel(
		s.tunnelRelayer.TunnelID,
		bandLatestSequence,
		"",
		"",
		true,
	), nil)

	s.chainProvider.EXPECT().
		QueryTunnelInfo(s.ctx, s.tunnelRelayer.TunnelID, s.tunnelRelayer.ContractAddress).
		Return(&chaintypes.Tunnel{
			ID:             s.tunnelRelayer.TunnelID,
			TargetAddress:  s.tunnelRelayer.ContractAddress,
			IsActive:       true,
			LatestSequence: targetChainLatestSequence,
			Balance:        big.NewInt(1),
		}, nil)

	// Run CheckAndRelay
	err := s.tunnelRelayer.CheckAndRelay(s.ctx)

	// Assertions
	s.Require().NoError(err)
}

func (s *TunnelRelayerTestSuite) TestCheckAndRelayFailedGetPacket() {
	bandLatestSequence := uint64(1)
	targetChainLatestSequence := uint64(0)

	// Mock BandClient to return tunnel info
	s.client.EXPECT().GetTunnel(s.ctx, s.tunnelRelayer.TunnelID).Return(bandtypes.NewTunnel(
		s.tunnelRelayer.TunnelID,
		bandLatestSequence,
		"",
		"",
		true,
	), nil)

	// Mock TargetChainProvider to return chain tunnel info
	s.chainProvider.EXPECT().
		QueryTunnelInfo(s.ctx, s.tunnelRelayer.TunnelID, s.tunnelRelayer.ContractAddress).
		Return(&chaintypes.Tunnel{
			ID:             s.tunnelRelayer.TunnelID,
			TargetAddress:  s.tunnelRelayer.ContractAddress,
			IsActive:       true,
			LatestSequence: targetChainLatestSequence,
			Balance:        big.NewInt(1),
		}, nil)

	s.client.EXPECT().GetTunnelPacket(s.ctx, s.tunnelRelayer.TunnelID, targetChainLatestSequence+1).
		Return(nil, fmt.Errorf("failed to get packet"))

	// Run CheckAndRelay
	err := s.tunnelRelayer.CheckAndRelay(s.ctx)

	// Assertions
	s.Require().Error(err)
}

func (s *TunnelRelayerTestSuite) TestCheckAndRelaySigningStatusFallen() {
	bandLatestSequence := uint64(1)
	targetChainLatestSequence := uint64(0)
	s.client.EXPECT().GetTunnel(s.ctx, s.tunnelRelayer.TunnelID).Return(bandtypes.NewTunnel(
		s.tunnelRelayer.TunnelID,
		bandLatestSequence,
		"",
		"",
		true,
	), nil)

	s.chainProvider.EXPECT().
		QueryTunnelInfo(s.ctx, s.tunnelRelayer.TunnelID, s.tunnelRelayer.ContractAddress).
		Return(&chaintypes.Tunnel{
			ID:             s.tunnelRelayer.TunnelID,
			TargetAddress:  s.tunnelRelayer.ContractAddress,
			IsActive:       true,
			LatestSequence: targetChainLatestSequence,
			Balance:        big.NewInt(1),
		}, nil)

	signalPrices := []bandtypes.SignalPrice{
		{SignalID: "signal1", Price: 100},
		{SignalID: "signal2", Price: 200},
	}

	// Create a mock EVMSignature
	evmSignature := bandtypes.NewEVMSignature(
		cmbytes.HexBytes("0x1234"),
		cmbytes.HexBytes("0xabcd"),
	)

	// Create mock signing information
	signing := bandtypes.NewSigning(
		1,
		cmbytes.HexBytes("0xdeadbeef"),
		evmSignature,
		"SIGNING_STATUS_FALLEN",
	)

	// Create the expected Packet object
	packet := bandtypes.NewPacket(
		s.tunnelRelayer.TunnelID,
		targetChainLatestSequence+1,
		signalPrices,
		signing,
		nil,
	)
	s.client.EXPECT().GetTunnelPacket(s.ctx, s.tunnelRelayer.TunnelID, targetChainLatestSequence+1).Return(
		packet, nil,
	)

	err := s.tunnelRelayer.CheckAndRelay(s.ctx)

	s.Require().Error(err)
}

func (s *TunnelRelayerTestSuite) TestCheckAndRelaySigningStatusWaiting() {
	bandLatestSequence := uint64(1)
	targetChainLatestSequence := uint64(0)
	s.client.EXPECT().GetTunnel(s.ctx, s.tunnelRelayer.TunnelID).Return(bandtypes.NewTunnel(
		s.tunnelRelayer.TunnelID,
		bandLatestSequence,
		"",
		"",
		true,
	), nil)

	s.chainProvider.EXPECT().
		QueryTunnelInfo(s.ctx, s.tunnelRelayer.TunnelID, s.tunnelRelayer.ContractAddress).
		Return(&chaintypes.Tunnel{
			ID:             s.tunnelRelayer.TunnelID,
			TargetAddress:  s.tunnelRelayer.ContractAddress,
			IsActive:       true,
			LatestSequence: targetChainLatestSequence,
			Balance:        big.NewInt(1),
		}, nil)

	signalPrices := []bandtypes.SignalPrice{
		{SignalID: "signal1", Price: 100},
		{SignalID: "signal2", Price: 200},
	}

	// Create a mock EVMSignature
	evmSignature := bandtypes.NewEVMSignature(
		cmbytes.HexBytes("0x1234"),
		cmbytes.HexBytes("0xabcd"),
	)

	// Create mock signing information
	signing := bandtypes.NewSigning(
		1,
		cmbytes.HexBytes("0xdeadbeef"),
		evmSignature,
		"SIGNING_STATUS_WAITING",
	)

	// Create the expected Packet object
	packet := bandtypes.NewPacket(
		s.tunnelRelayer.TunnelID,
		targetChainLatestSequence+1,
		signalPrices,
		signing,
		nil,
	)
	s.client.EXPECT().GetTunnelPacket(s.ctx, s.tunnelRelayer.TunnelID, targetChainLatestSequence+1).Return(
		packet, nil,
	)

	err := s.tunnelRelayer.CheckAndRelay(s.ctx)
	s.Require().NoError(err)
}

func (s *TunnelRelayerTestSuite) TestCheckAndRelayFailedToRelayPacket() {
	bandLatestSequence := uint64(1)
	targetChainLatestSequence := uint64(0)
	s.client.EXPECT().GetTunnel(s.ctx, s.tunnelRelayer.TunnelID).Return(bandtypes.NewTunnel(
		s.tunnelRelayer.TunnelID,
		bandLatestSequence,
		"",
		"",
		true,
	), nil)

	s.chainProvider.EXPECT().
		QueryTunnelInfo(s.ctx, s.tunnelRelayer.TunnelID, s.tunnelRelayer.ContractAddress).
		Return(&chaintypes.Tunnel{
			ID:             s.tunnelRelayer.TunnelID,
			TargetAddress:  s.tunnelRelayer.ContractAddress,
			IsActive:       true,
			LatestSequence: targetChainLatestSequence,
			Balance:        big.NewInt(1),
		}, nil)

	signalPrices := []bandtypes.SignalPrice{
		{SignalID: "signal1", Price: 100},
		{SignalID: "signal2", Price: 200},
	}

	// Create a mock EVMSignature
	evmSignature := bandtypes.NewEVMSignature(
		cmbytes.HexBytes("0x1234"),
		cmbytes.HexBytes("0xabcd"),
	)

	// Create mock signing information
	signing := bandtypes.NewSigning(
		1,
		cmbytes.HexBytes("0xdeadbeef"),
		evmSignature,
		"SIGNING_STATUS_SUCCESS",
	)

	// Create the expected Packet object
	packet := bandtypes.NewPacket(
		s.tunnelRelayer.TunnelID,
		targetChainLatestSequence+1,
		signalPrices,
		signing,
		nil,
	)
	s.client.EXPECT().GetTunnelPacket(s.ctx, s.tunnelRelayer.TunnelID, targetChainLatestSequence+1).Return(
		packet, nil,
	)

	s.chainProvider.EXPECT().RelayPacket(s.ctx, packet).Return(
		fmt.Errorf("failed to relay packet"),
	)

	err := s.tunnelRelayer.CheckAndRelay(s.ctx)
	s.Require().ErrorContains(err, "failed to relay packet")
}
