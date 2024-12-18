package band_test

import (
	"context"
	"testing"
	"time"

	"cosmossdk.io/math"
	cmbytes "github.com/cometbft/cometbft/libs/bytes"
	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types"
	querytypes "github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	bandtsstypes "github.com/bandprotocol/chain/v3/x/bandtss/types"
	feedstypes "github.com/bandprotocol/chain/v3/x/feeds/types"
	tsstypes "github.com/bandprotocol/chain/v3/x/tss/types"
	tunneltypes "github.com/bandprotocol/chain/v3/x/tunnel/types"

	"github.com/bandprotocol/falcon/internal/relayertest/mocks"
	"github.com/bandprotocol/falcon/relayer/band"
	bandclienttypes "github.com/bandprotocol/falcon/relayer/band/types"
)

type ClientTestSuite struct {
	suite.Suite

	ctx                context.Context
	tunnelQueryClient  *mocks.MockTunnelQueryClient
	bandtssQueryClient *mocks.MockBandtssQueryClient
	client             band.Client
	log                *zap.Logger
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

// SetupTest sets up the test suite by creating a temporary directory and declare mock objects.
func (s *ClientTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())

	log, err := zap.NewDevelopment()
	s.Require().NoError(err)

	// mock objects.
	s.log = log
	s.tunnelQueryClient = mocks.NewMockTunnelQueryClient(ctrl)
	s.bandtssQueryClient = mocks.NewMockBandtssQueryClient(ctrl)
	s.ctx = context.Background()

	encodingConfig := band.MakeEncodingConfig()
	s.client = band.NewClient(
		cosmosclient.Context{}.
			WithCodec(encodingConfig.Marshaler).
			WithInterfaceRegistry(encodingConfig.InterfaceRegistry),
		band.NewQueryClient(s.tunnelQueryClient, s.bandtssQueryClient),
		s.log,
		[]string{})
}

func (s *ClientTestSuite) TestGetTSSTunnel() {
	// mock route value
	destinationChainID := "eth"
	destinationContractAddress := "0xe00F1f85abDB2aF6760759547d450da68CE66Bb1"
	r := &tunneltypes.TSSRoute{
		DestinationChainID:         destinationChainID,
		DestinationContractAddress: destinationContractAddress,
	}
	var routeI tunneltypes.RouteI = r
	msg, ok := routeI.(proto.Message)
	s.Require().Equal(true, ok)

	any, err := codectypes.NewAnyWithValue(msg)
	s.Require().NoError(err)

	tunnel := tunneltypes.Tunnel{
		ID:               uint64(1),
		Sequence:         100,
		Route:            any,
		FeePayer:         "cosmos1xyz...",
		SignalDeviations: []tunneltypes.SignalDeviation{},
		Interval:         60,
		TotalDeposit:     types.NewCoins(types.NewCoin("uband", math.NewInt(1000))),
		IsActive:         false,
		CreatedAt:        time.Now().Unix(),
		Creator:          "cosmos1abc...",
	}
	queryResponse := &tunneltypes.QueryTunnelResponse{
		Tunnel: tunnel,
	}

	// expect response from bandQueryClient
	s.tunnelQueryClient.EXPECT().Tunnel(s.ctx, &tunneltypes.QueryTunnelRequest{
		TunnelId: uint64(1),
	}).Return(queryResponse, nil)

	expected := bandclienttypes.NewTunnel(1, 100, "0xe00F1f85abDB2aF6760759547d450da68CE66Bb1", "eth", false)

	actual, err := s.client.GetTunnel(s.ctx, uint64(1))
	s.Require().NoError(err)
	s.Require().Equal(expected, actual)
}

func (s *ClientTestSuite) TestGetOtherTunnel() {
	// mock ibc tunnel
	ibcRoute := tunneltypes.IBCRoute{ChannelID: "test"}
	var routeI tunneltypes.RouteI = &ibcRoute
	msg, ok := routeI.(proto.Message)
	s.Require().Equal(true, ok)

	routeIBCAny, err := codectypes.NewAnyWithValue(msg)
	s.Require().NoError(err)

	tunnel := tunneltypes.Tunnel{
		ID:               2,
		Sequence:         100,
		Route:            routeIBCAny,
		FeePayer:         "cosmos1xyz...",
		SignalDeviations: []tunneltypes.SignalDeviation{},
		Interval:         60,
		TotalDeposit:     types.NewCoins(types.NewCoin("uband", math.NewInt(1000))),
		IsActive:         false,
		CreatedAt:        time.Now().Unix(),
		Creator:          "cosmos1abc...",
	}

	// expect response from bandQueryClient
	s.tunnelQueryClient.EXPECT().Tunnel(s.ctx, &tunneltypes.QueryTunnelRequest{
		TunnelId: uint64(1),
	}).Return(&tunneltypes.QueryTunnelResponse{Tunnel: tunnel}, nil)

	_, err = s.client.GetTunnel(s.ctx, uint64(1))
	s.Require().ErrorContains(err, "unsupported route type")
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
	s.tunnelQueryClient.EXPECT().Packet(s.ctx, &tunneltypes.QueryPacketRequest{
		TunnelId: uint64(1),
		Sequence: uint64(100),
	}).Return(queryPacketResponse, nil)
	s.bandtssQueryClient.EXPECT().Signing(s.ctx, &bandtsstypes.QuerySigningRequest{
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
		"SIGNING_STATUS_SUCCESS",
	)

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
	s.tunnelQueryClient.EXPECT().Packet(s.ctx, &tunneltypes.QueryPacketRequest{
		TunnelId: uint64(1),
		Sequence: uint64(100),
	}).Return(queryPacketResponse, nil)

	// actual result
	_, err = s.client.GetTunnelPacket(s.ctx, uint64(1), uint64(100))
	s.Require().ErrorContains(err, "unsupported packet content type")
}

func (s *ClientTestSuite) TestGetTunnels() {
	// mock route value
	destinationChainID := "eth"
	destinationContractAddress := "0xe00F1f85abDB2aF6760759547d450da68CE66Bb1"
	r := &tunneltypes.TSSRoute{
		DestinationChainID:         destinationChainID,
		DestinationContractAddress: destinationContractAddress,
	}
	var routeI tunneltypes.RouteI = r
	msg, ok := routeI.(proto.Message)
	s.Require().Equal(true, ok)

	routeAny, err := codectypes.NewAnyWithValue(msg)
	s.Require().NoError(err)

	// mock tunnels result
	tunnels := make([]*tunneltypes.Tunnel, 0, 120)
	for i := 1; i <= cap(tunnels); i++ {
		tunnel := &tunneltypes.Tunnel{
			ID:               uint64(i),
			Sequence:         100,
			Route:            routeAny,
			FeePayer:         "cosmos1xyz...",
			SignalDeviations: []tunneltypes.SignalDeviation{},
			Interval:         60,
			TotalDeposit:     types.NewCoins(types.NewCoin("uband", math.NewInt(1000))),
			IsActive:         false,
			CreatedAt:        time.Now().Unix(),
			Creator:          "cosmos1abc...",
		}
		tunnels = append(tunnels, tunnel)
	}

	// mock responses
	queryResponse1 := &tunneltypes.QueryTunnelsResponse{
		Tunnels: tunnels[:100],
		Pagination: &querytypes.PageResponse{
			NextKey: []byte("next-key"),
		},
	}
	queryResponse2 := &tunneltypes.QueryTunnelsResponse{
		Tunnels: tunnels[100:],
		Pagination: &querytypes.PageResponse{
			NextKey: []byte(""),
		},
	}

	// expect response from bandQueryClient
	s.tunnelQueryClient.EXPECT().Tunnels(s.ctx, &tunneltypes.QueryTunnelsRequest{
		Pagination: &querytypes.PageRequest{
			Key: nil,
		},
	}).Return(queryResponse1, nil)
	s.tunnelQueryClient.EXPECT().Tunnels(s.ctx, &tunneltypes.QueryTunnelsRequest{
		Pagination: &querytypes.PageRequest{
			Key: []byte("next-key"),
		},
	}).Return(queryResponse2, nil)

	expected := make([]bandclienttypes.Tunnel, 0, len(tunnels))
	for i := 1; i <= len(tunnels); i++ {
		expected = append(
			expected,
			*bandclienttypes.NewTunnel(uint64(i), 100, "0xe00F1f85abDB2aF6760759547d450da68CE66Bb1", "eth", false),
		)
	}

	actual, err := s.client.GetTunnels(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal(expected, actual)
}

func (s *ClientTestSuite) TestGetTunnelsGetMultipleRouteTypes() {
	// mock tss tunnel
	destinationChainID := "eth"
	destinationContractAddress := "0xe00F1f85abDB2aF6760759547d450da68CE66Bb1"
	r := &tunneltypes.TSSRoute{
		DestinationChainID:         destinationChainID,
		DestinationContractAddress: destinationContractAddress,
	}
	var routeI tunneltypes.RouteI = r
	msg, ok := routeI.(proto.Message)
	s.Require().Equal(true, ok)

	routeTSSAny, err := codectypes.NewAnyWithValue(msg)
	s.Require().NoError(err)

	tssTunnel := tunneltypes.Tunnel{
		ID:               1,
		Sequence:         100,
		Route:            routeTSSAny,
		FeePayer:         "cosmos1xyz...",
		SignalDeviations: []tunneltypes.SignalDeviation{},
		Interval:         60,
		TotalDeposit:     types.NewCoins(types.NewCoin("uband", math.NewInt(1000))),
		IsActive:         false,
		CreatedAt:        time.Now().Unix(),
		Creator:          "cosmos1abc...",
	}

	// mock ibc tunnel
	ibcRoute := tunneltypes.IBCRoute{ChannelID: "test"}
	routeI = &ibcRoute
	msg, ok = routeI.(proto.Message)
	s.Require().Equal(true, ok)

	routeIBCAny, err := codectypes.NewAnyWithValue(msg)
	s.Require().NoError(err)

	ibcTunnel := tunneltypes.Tunnel{
		ID:               2,
		Sequence:         100,
		Route:            routeIBCAny,
		FeePayer:         "cosmos1xyz...",
		SignalDeviations: []tunneltypes.SignalDeviation{},
		Interval:         60,
		TotalDeposit:     types.NewCoins(types.NewCoin("uband", math.NewInt(1000))),
		IsActive:         false,
		CreatedAt:        time.Now().Unix(),
		Creator:          "cosmos1abc...",
	}

	// mock responses
	queryResponse := &tunneltypes.QueryTunnelsResponse{
		Tunnels:    []*tunneltypes.Tunnel{&tssTunnel, &ibcTunnel},
		Pagination: &querytypes.PageResponse{NextKey: []byte("")},
	}

	// expect response from bandQueryClient
	s.tunnelQueryClient.EXPECT().Tunnels(s.ctx, &tunneltypes.QueryTunnelsRequest{
		Pagination: &querytypes.PageRequest{
			Key: nil,
		},
	}).Return(queryResponse, nil)

	encodingConfig := band.MakeEncodingConfig()
	s.client = band.NewClient(
		cosmosclient.Context{}.
			WithCodec(encodingConfig.Marshaler).
			WithInterfaceRegistry(encodingConfig.InterfaceRegistry),
		band.NewQueryClient(s.tunnelQueryClient, s.bandtssQueryClient),
		s.log,
		[]string{})

	expected := make([]bandclienttypes.Tunnel, 1)
	expected[0] = *bandclienttypes.NewTunnel(1, 100, "0xe00F1f85abDB2aF6760759547d450da68CE66Bb1", "eth", false)

	actual, err := s.client.GetTunnels(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal(expected, actual)
}
