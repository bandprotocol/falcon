package band_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"
	cmbytes "github.com/cometbft/cometbft/libs/bytes"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types"
	querytypes "github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	bandtsstypes "github.com/bandprotocol/falcon/internal/bandchain/bandtss"
	feedstypes "github.com/bandprotocol/falcon/internal/bandchain/feeds"
	tsstypes "github.com/bandprotocol/falcon/internal/bandchain/tss"
	tunneltypes "github.com/bandprotocol/falcon/internal/bandchain/tunnel"
	"github.com/bandprotocol/falcon/internal/relayertest/mocks"
	"github.com/bandprotocol/falcon/relayer/band"
	bandclienttypes "github.com/bandprotocol/falcon/relayer/band/types"
)

type ClientTestSuite struct {
	suite.Suite

	ctx             context.Context
	bandQueryClient *mocks.MockQueryClient
	client          band.Client
	log             *zap.Logger
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

// SetupTest sets up the test suite by creating a temporary directory and declare mock objects.
func (s *ClientTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())

	// mock objects.
	s.log = zap.NewNop()
	s.bandQueryClient = mocks.NewMockQueryClient(ctrl)
	s.client = band.NewClient(
		s.bandQueryClient,
		s.log,
		&band.Config{LivelinessCheckingInterval: 15 * time.Minute},
	)
	s.ctx = context.Background()
}

// GetMockIBCTunnel returns a mock IBC tunnel.
func (s *ClientTestSuite) GetMockIBCTunnel(tunnelID uint64) (tunneltypes.Tunnel, error) {
	ibcRoute := tunneltypes.IBCRoute{ChannelID: "test"}
	var routeI tunneltypes.RouteI = &ibcRoute

	msg, ok := routeI.(proto.Message)
	if !ok {
		return tunneltypes.Tunnel{}, fmt.Errorf("cannot convert route to proto.Message")
	}

	routeAny, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return tunneltypes.Tunnel{}, err
	}

	return tunneltypes.Tunnel{
		ID:               tunnelID,
		Sequence:         100,
		Route:            routeAny,
		FeePayer:         "cosmos1xyz...",
		SignalDeviations: []tunneltypes.SignalDeviation{},
		Interval:         60,
		TotalDeposit:     types.NewCoins(types.NewCoin("uband", math.NewInt(1000))),
		IsActive:         false,
		CreatedAt:        1736145613,
		Creator:          "cosmos1abc...",
	}, nil
}

// GetMockTSSTunnel returns a mock TSS tunnel.
func (s *ClientTestSuite) GetMockTSSTunnel(tunnelID uint64) (tunneltypes.Tunnel, error) {
	r := &tunneltypes.TSSRoute{
		DestinationChainID:         "eth",
		DestinationContractAddress: "0xe00F1f85abDB2aF6760759547d450da68CE66Bb1",
	}
	var routeI tunneltypes.RouteI = r

	msg, ok := routeI.(proto.Message)
	if !ok {
		return tunneltypes.Tunnel{}, fmt.Errorf("cannot convert route to proto.Message")
	}

	routeAny, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return tunneltypes.Tunnel{}, err
	}

	return tunneltypes.Tunnel{
		ID:               tunnelID,
		Sequence:         100,
		Route:            routeAny,
		FeePayer:         "cosmos1xyz...",
		SignalDeviations: []tunneltypes.SignalDeviation{},
		Interval:         60,
		TotalDeposit:     types.NewCoins(types.NewCoin("uband", math.NewInt(1000))),
		IsActive:         false,
		CreatedAt:        1736145613,
		Creator:          "cosmos1abc...",
	}, nil
}

func (s *ClientTestSuite) TestGetTunnel() {
	tssTunnel, err := s.GetMockTSSTunnel(1)
	s.Require().NoError(err)

	ibcTunnel, err := s.GetMockIBCTunnel(2)
	s.Require().NoError(err)

	testcases := []struct {
		name       string
		in         uint64
		preprocess func(c context.Context)
		out        *bandclienttypes.Tunnel
		err        error
	}{
		{
			name: "success",
			in:   1,
			out: bandclienttypes.NewTunnel(
				1,
				100,
				"0xe00F1f85abDB2aF6760759547d450da68CE66Bb1",
				"eth",
				false,
				"cosmos1abc...",
			),
			preprocess: func(c context.Context) {
				s.bandQueryClient.EXPECT().Tunnel(s.ctx, &tunneltypes.QueryTunnelRequest{
					TunnelId: uint64(1),
				}).Return(&tunneltypes.QueryTunnelResponse{Tunnel: tssTunnel}, nil)
			},
		},
		{
			name: "unsupported route type",
			in:   2,
			err:  fmt.Errorf("unsupported route type"),
			preprocess: func(c context.Context) {
				s.bandQueryClient.EXPECT().Tunnel(s.ctx, &tunneltypes.QueryTunnelRequest{
					TunnelId: uint64(2),
				}).Return(&tunneltypes.QueryTunnelResponse{Tunnel: ibcTunnel}, nil)
			},
		},
	}

	for _, tc := range testcases {
		s.T().Run(tc.name, func(t *testing.T) {
			if tc.preprocess != nil {
				tc.preprocess(s.ctx)
			}

			actual, err := s.client.GetTunnel(s.ctx, tc.in)
			if tc.err != nil {
				s.Require().ErrorContains(err, tc.err.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.out, actual)
			}
		})
	}
}

func (s *ClientTestSuite) TestGetTSSTunnelPacket() {
	// mock query response
	pc := &tunneltypes.TSSPacketReceipt{
		SigningID: 2,
	}

	var packetReceiptI tunneltypes.PacketReceiptI = pc
	msg, ok := packetReceiptI.(proto.Message)
	s.Require().Equal(true, ok)

	any, err := codectypes.NewAnyWithValue(msg)
	s.Require().NoError(err)

	packet := tunneltypes.Packet{
		TunnelID: 1,
		Sequence: 100,
		Prices: []feedstypes.Price{
			{SignalID: "signal1", Price: 100},
			{SignalID: "signal2", Price: 200},
		},
		Receipt:   any,
		CreatedAt: time.Now().Unix(),
	}
	signingResult := &tsstypes.SigningResult{
		Signing: tsstypes.Signing{
			ID:      2,
			Message: cmbytes.HexBytes("0xdeadbeef"),
			Status:  tsstypes.SIGNING_STATUS_SUCCESS,
		},
		EVMSignature: &tsstypes.EVMSignature{
			RAddress:  cmbytes.HexBytes("0x1234"),
			Signature: cmbytes.HexBytes("0xabcd"),
		},
	}

	queryPacketResponse := &tunneltypes.QueryPacketResponse{
		Packet: &packet,
	}
	querySigningResponse := &bandtsstypes.QuerySigningResponse{
		CurrentGroupSigningResult:  signingResult,
		IncomingGroupSigningResult: nil,
	}

	// expect response from bandQueryClient
	s.bandQueryClient.EXPECT().Packet(s.ctx, &tunneltypes.QueryPacketRequest{
		TunnelId: uint64(1),
		Sequence: uint64(100),
	}).Return(queryPacketResponse, nil)
	s.bandQueryClient.EXPECT().Signing(s.ctx, &bandtsstypes.QuerySigningRequest{
		SigningId: uint64(2),
	}).Return(querySigningResponse, nil)

	// expected result
	expectedSignalPrices := []bandclienttypes.SignalPrice{
		{SignalID: "signal1", Price: 100},
		{SignalID: "signal2", Price: 200},
	}

	expectedCurrentGroupSigning := bandclienttypes.NewSigning(
		2,
		cmbytes.HexBytes("0xdeadbeef"),
		bandclienttypes.NewEVMSignature(
			cmbytes.HexBytes("0x1234"),
			cmbytes.HexBytes("0xabcd"),
		),
		tsstypes.SIGNING_STATUS_SUCCESS,
	)
	expectedCurrentGroupSigning.SigningStatusString = tsstypes.SigningStatus_name[int32(tsstypes.SIGNING_STATUS_SUCCESS)]

	// Define the expected Packet result
	expected := bandclienttypes.NewPacket(
		1,
		100,
		expectedSignalPrices,
		expectedCurrentGroupSigning,
		nil,
	)

	// actual result
	actual, err := s.client.GetTunnelPacket(s.ctx, uint64(1), uint64(100))
	s.Require().NoError(err)
	s.Require().Equal(expected, actual)
}

func (s *ClientTestSuite) TestGetOtherTunnelPacket() {
	// mock query response
	pc := &tunneltypes.IBCPacketReceipt{
		Sequence: 2,
	}

	var packetReceiptI tunneltypes.PacketReceiptI = pc
	msg, ok := packetReceiptI.(proto.Message)
	s.Require().Equal(true, ok)

	any, err := codectypes.NewAnyWithValue(msg)
	s.Require().NoError(err)

	packet := tunneltypes.Packet{
		TunnelID: 1,
		Sequence: 100,
		Prices: []feedstypes.Price{
			{SignalID: "signal1", Price: 100},
			{SignalID: "signal2", Price: 200},
		},
		Receipt:   any,
		CreatedAt: time.Now().Unix(),
	}

	queryPacketResponse := &tunneltypes.QueryPacketResponse{
		Packet: &packet,
	}

	// expect response from bandQueryClient
	s.bandQueryClient.EXPECT().Packet(s.ctx, &tunneltypes.QueryPacketRequest{
		TunnelId: uint64(1),
		Sequence: uint64(100),
	}).Return(queryPacketResponse, nil)

	// actual result
	_, err = s.client.GetTunnelPacket(s.ctx, uint64(1), uint64(100))
	s.Require().ErrorContains(err, "unsupported packet content type")
}

func (s *ClientTestSuite) TestGetTunnels() {
	// mock tunnels result
	tssTunnels := make([]*tunneltypes.Tunnel, 0, 120)
	for i := 1; i <= cap(tssTunnels); i++ {
		tunnel, err := s.GetMockTSSTunnel(uint64(i))
		s.Require().NoError(err)

		tssTunnels = append(tssTunnels, &tunnel)
	}

	// expected result from tssTunnels
	expectedRes := make([]bandclienttypes.Tunnel, 0, len(tssTunnels))
	for _, tunnel := range tssTunnels {
		routeI, ok := tunnel.Route.GetCachedValue().(tunneltypes.RouteI)
		s.Require().True(ok)

		tssRoute, ok := routeI.(*tunneltypes.TSSRoute)
		s.Require().True(ok)

		expectedRes = append(expectedRes, *bandclienttypes.NewTunnel(
			tunnel.ID,
			tunnel.Sequence,
			tssRoute.DestinationContractAddress,
			tssRoute.DestinationChainID,
			tunnel.IsActive,
			tunnel.Creator,
		))
	}

	// create mock ibc tunnel
	ibcTunnel, err := s.GetMockIBCTunnel(uint64(121))
	s.Require().NoError(err)

	testcases := []struct {
		name       string
		preprocess func(c context.Context)
		out        []bandclienttypes.Tunnel
		err        error
	}{
		{
			name: "success",
			preprocess: func(c context.Context) {
				// expect response from bandQueryClient
				s.bandQueryClient.EXPECT().Tunnels(s.ctx, &tunneltypes.QueryTunnelsRequest{
					Pagination: &querytypes.PageRequest{
						Key: nil,
					},
				}).Return(&tunneltypes.QueryTunnelsResponse{
					Tunnels: tssTunnels[:100],
					Pagination: &querytypes.PageResponse{
						NextKey: []byte("next-key"),
					},
				}, nil)

				s.bandQueryClient.EXPECT().Tunnels(s.ctx, &tunneltypes.QueryTunnelsRequest{
					Pagination: &querytypes.PageRequest{
						Key: []byte("next-key"),
					},
				}).Return(&tunneltypes.QueryTunnelsResponse{
					Tunnels: tssTunnels[100:],
					Pagination: &querytypes.PageResponse{
						NextKey: []byte(""),
					},
				}, nil)
			},
			out: expectedRes,
		},
		{
			name: "filter out unrelated tunnel",
			preprocess: func(c context.Context) {
				s.bandQueryClient.EXPECT().Tunnels(s.ctx, &tunneltypes.QueryTunnelsRequest{
					Pagination: &querytypes.PageRequest{
						Key: nil,
					},
				}).Return(&tunneltypes.QueryTunnelsResponse{
					Tunnels: []*tunneltypes.Tunnel{tssTunnels[0], &ibcTunnel},
					Pagination: &querytypes.PageResponse{
						NextKey: []byte(""),
					},
				}, nil)
			},
			out: []bandclienttypes.Tunnel{expectedRes[0]},
		},
	}

	for _, tc := range testcases {
		s.T().Run(tc.name, func(t *testing.T) {
			if tc.preprocess != nil {
				tc.preprocess(s.ctx)
			}

			actual, err := s.client.GetTunnels(s.ctx)
			if tc.err != nil {
				s.Require().ErrorContains(err, tc.err.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.out, actual)
			}
		})
	}
}
