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

	"github.com/bandprotocol/falcon/internal/bandchain/tss"
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
	defaultTargetChainID          = ""
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
	), nil)
}

// Helper function to mock QueryTunnelInfo.
func (s *TunnelRelayerTestSuite) mockQueryTunnelInfo(sequence uint64, isActive bool, contractAddress string) {
	s.chainProvider.EXPECT().
		QueryTunnelInfo(s.ctx, s.tunnelRelayer.TunnelID, contractAddress).
		Return(&chaintypes.Tunnel{
			ID:             s.tunnelRelayer.TunnelID,
			TargetAddress:  contractAddress,
			IsActive:       isActive,
			LatestSequence: sequence,
			Balance:        big.NewInt(1),
		}, nil)
}

// Helper function to create a mock Packet.
func createMockPacket(
	tunnelID, sequence uint64,
	currentStatus int32,
	incomingStatus int32,
) *bandtypes.Packet {
	signalPrices := []bandtypes.SignalPrice{
		{SignalID: "signal1", Price: 100},
		{SignalID: "signal2", Price: 200},
	}
	evmSignature := bandtypes.NewEVMSignature(
		cmbytes.HexBytes("0x1234"),
		cmbytes.HexBytes("0xabcd"),
	)
	var currentGroupSigning *bandtypes.Signing
	var incomingGroupSigning *bandtypes.Signing

	if currentStatus != -1 {
		currentGroupSigning = bandtypes.NewSigning(
			1,
			cmbytes.HexBytes("0xdeadbeef"),
			evmSignature,
			tss.SigningStatus(currentStatus),
		)
	}

	if incomingStatus != -1 {
		incomingGroupSigning = bandtypes.NewSigning(
			1,
			cmbytes.HexBytes("0xdeadbeef"),
			evmSignature,
			tss.SigningStatus(incomingStatus),
		)
	}

	return bandtypes.NewPacket(
		tunnelID,
		sequence,
		signalPrices,
		currentGroupSigning,
		incomingGroupSigning,
	)
}

func (s *TunnelRelayerTestSuite) TestCheckAndRelay() {
	testcases := []struct {
		name       string
		preprocess func()
		err        error
	}{
		{
			name: "success",
			preprocess: func() {
				s.mockGetTunnel(defaultBandLatestSequence)
				s.mockQueryTunnelInfo(defaultTargetChainSequence, true, "")

				packet := createMockPacket(
					s.tunnelRelayer.TunnelID,
					defaultTargetChainSequence+1,
					int32(tss.SIGNING_STATUS_SUCCESS),
					-1,
				)
				s.client.EXPECT().
					GetTunnelPacket(gomock.Any(), s.tunnelRelayer.TunnelID, defaultTargetChainSequence+1).
					Return(packet, nil)
				s.chainProvider.EXPECT().RelayPacket(gomock.Any(), packet).Return(nil)

				// Check and relay the packet for the second time
				s.mockGetTunnel(defaultBandLatestSequence)
				s.mockQueryTunnelInfo(defaultTargetChainSequence+1, true, defaultContractAddress)
			},
		},
		{
			name: "failed to get tunnel on band client",
			preprocess: func() {
				s.client.EXPECT().
					GetTunnel(s.ctx, s.tunnelRelayer.TunnelID).
					Return(nil, fmt.Errorf("failed to get tunnel"))
			},
			err: fmt.Errorf("failed to get tunnel"),
		},
		{
			name: "failed to query chain tunnel info",
			preprocess: func() {
				s.mockGetTunnel(defaultBandLatestSequence)
				s.chainProvider.EXPECT().
					QueryTunnelInfo(gomock.Any(), s.tunnelRelayer.TunnelID, defaultContractAddress).
					Return(nil, fmt.Errorf("failed to query tunnel info"))
			},
			err: fmt.Errorf("failed to query tunnel info"),
		},
		{
			name: "target chain not active",
			preprocess: func() {
				s.mockGetTunnel(defaultBandLatestSequence)
				s.mockQueryTunnelInfo(defaultTargetChainSequence, false, defaultContractAddress)
			},
			err: nil,
		},
		{
			name: "no new packet to relay",
			preprocess: func() {
				s.mockGetTunnel(defaultBandLatestSequence)
				s.mockQueryTunnelInfo(defaultTargetChainSequence+1, true, defaultContractAddress)
			},
			err: nil,
		},
		{
			name: "fail to get a new packet",
			preprocess: func() {
				s.mockGetTunnel(defaultBandLatestSequence)
				s.mockQueryTunnelInfo(defaultTargetChainSequence, true, defaultContractAddress)

				s.client.EXPECT().
					GetTunnelPacket(gomock.Any(), s.tunnelRelayer.TunnelID, defaultTargetChainSequence+1).
					Return(nil, fmt.Errorf("failed to get packet"))
			},
			err: fmt.Errorf("failed to get packet"),
		},
		{
			name: "fallen signing status of current group but incoming group success",
			preprocess: func() {
				s.mockGetTunnel(defaultBandLatestSequence)
				s.mockQueryTunnelInfo(defaultTargetChainSequence, true, defaultContractAddress)

				packet := createMockPacket(
					s.tunnelRelayer.TunnelID,
					defaultTargetChainSequence+1,
					int32(tss.SIGNING_STATUS_FALLEN),
					int32(tss.SIGNING_STATUS_SUCCESS),
				)
				s.client.EXPECT().
					GetTunnelPacket(gomock.Any(), s.tunnelRelayer.TunnelID, defaultTargetChainSequence+1).
					Return(packet, nil)
				s.chainProvider.EXPECT().RelayPacket(gomock.Any(), packet).Return(nil)

				// Check and relay the packet for the second time
				s.mockGetTunnel(defaultBandLatestSequence)
				s.mockQueryTunnelInfo(defaultTargetChainSequence+1, true, defaultContractAddress)
			},
		},
		{
			name: "incoming group signing status fallen",
			preprocess: func() {
				s.mockGetTunnel(defaultBandLatestSequence)
				s.mockQueryTunnelInfo(defaultTargetChainSequence, true, defaultContractAddress)

				packet := createMockPacket(
					s.tunnelRelayer.TunnelID,
					defaultTargetChainSequence+1,
					int32(tss.SIGNING_STATUS_FALLEN),
					int32(tss.SIGNING_STATUS_FALLEN),
				)

				s.client.EXPECT().
					GetTunnelPacket(s.ctx, s.tunnelRelayer.TunnelID, defaultTargetChainSequence+1).
					Return(packet, nil)
			},
			err: fmt.Errorf(("signing status is not success")),
		},
		{
			name: "signing status is waiting",
			preprocess: func() {
				s.mockGetTunnel(defaultBandLatestSequence)
				s.mockQueryTunnelInfo(defaultTargetChainSequence, true, defaultContractAddress)

				packet := createMockPacket(
					s.tunnelRelayer.TunnelID,
					defaultTargetChainSequence+1,
					int32(tss.SIGNING_STATUS_WAITING),
					-1,
				)

				s.client.EXPECT().
					GetTunnelPacket(s.ctx, s.tunnelRelayer.TunnelID, defaultTargetChainSequence+1).
					Return(packet, nil)
			},
			err: nil,
		},
		{
			name: "failed to relay packet",
			preprocess: func() {
				s.mockGetTunnel(defaultBandLatestSequence)
				s.mockQueryTunnelInfo(defaultTargetChainSequence, true, defaultContractAddress)

				packet := createMockPacket(
					s.tunnelRelayer.TunnelID,
					defaultTargetChainSequence+1,
					int32(tss.SIGNING_STATUS_SUCCESS),
					-1,
				)

				s.client.EXPECT().
					GetTunnelPacket(s.ctx, s.tunnelRelayer.TunnelID, defaultTargetChainSequence+1).
					Return(packet, nil)
				s.chainProvider.EXPECT().RelayPacket(s.ctx, packet).Return(fmt.Errorf("failed to relay packet"))
			},
			err: fmt.Errorf("failed to relay packet"),
		},
	}

	for _, tc := range testcases {
		s.T().Run(tc.name, func(t *testing.T) {
			if tc.preprocess != nil {
				tc.preprocess()
			}

			isExecuting, err := s.tunnelRelayer.CheckAndRelay(s.ctx)
			if tc.err != nil {
				s.Require().ErrorContains(err, tc.err.Error())
				s.Require().False(isExecuting)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}
