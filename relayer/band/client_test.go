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

type AppTestSuite struct {
	suite.Suite

	ctx                context.Context
	tunnelQueryClient  *mocks.MockTunnelQueryClient
	bandtssQueryClient *mocks.MockBandtssQueryClient
	client             band.Client
	log                *zap.Logger
}

func TestAppTestSuite(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}

// SetupTest sets up the test suite by creating a temporary directory and declare mock objects.
func (s *AppTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())

	log, err := zap.NewDevelopment()
	s.Require().NoError(err)

	// mock objects.
	s.log = log
	s.tunnelQueryClient = mocks.NewMockTunnelQueryClient(ctrl)
	s.bandtssQueryClient = mocks.NewMockBandtssQueryClient(ctrl)
	s.ctx = context.Background()
}

func (s *AppTestSuite) TestGetTunnel() {
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
		Encoder:          0,
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
	encodingConfig := band.MakeEncodingConfig()
	s.client = band.NewClient(
		cosmosclient.Context{}.
			WithCodec(encodingConfig.Marshaler).
			WithInterfaceRegistry(encodingConfig.InterfaceRegistry),
		band.NewQueryClient(s.tunnelQueryClient, s.bandtssQueryClient),
		s.log,
		[]string{})

	expected := bandclienttypes.NewTunnel(1, 100, "0xe00F1f85abDB2aF6760759547d450da68CE66Bb1", "eth", false)

	actual, err := s.client.GetTunnel(s.ctx, uint64(1))
	s.Require().NoError(err)
	s.Require().Equal(expected, actual)
}

func (s *AppTestSuite) TestGetTunnelPacket() {
	// mock query response
	destinationChainID := "eth"
	destinationContractAddress := "0xe00F1f85abDB2aF6760759547d450da68CE66Bb1"
	pc := &tunneltypes.TSSPacketContent{
		SigningID:                  2,
		DestinationChainID:         destinationChainID,
		DestinationContractAddress: destinationContractAddress,
	}

	var packetContentI tunneltypes.PacketContentI = pc
	msg, ok := packetContentI.(proto.Message)
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
		PacketContent: any,
		CreatedAt:     time.Now().Unix(),
	}
	signingResult := &tsstypes.SigningResult{
		Signing: tsstypes.Signing{
			ID:      2,
			Message: cmbytes.HexBytes("0xdeadbeef"),
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
	encodingConfig := band.MakeEncodingConfig()
	s.client = band.NewClient(
		cosmosclient.Context{}.
			WithCodec(encodingConfig.Marshaler).
			WithInterfaceRegistry(encodingConfig.InterfaceRegistry),
		band.NewQueryClient(s.tunnelQueryClient, s.bandtssQueryClient),
		s.log,
		[]string{})
	actual, err := s.client.GetTunnelPacket(s.ctx, uint64(1), uint64(100))
	s.Require().NoError(err)
	s.Require().Equal(expected, actual)
}

func (s *AppTestSuite) TestGetTunnels() {
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

	// mock first request
	tunnels1 := make([]*tunneltypes.Tunnel, 0, 100)
	for i := 1; i <= 100; i++ {
		tunnel := &tunneltypes.Tunnel{
			ID:               uint64(i),
			Sequence:         100,
			Route:            any,
			Encoder:          0,
			FeePayer:         "cosmos1xyz...",
			SignalDeviations: []tunneltypes.SignalDeviation{},
			Interval:         60,
			TotalDeposit:     types.NewCoins(types.NewCoin("uband", math.NewInt(1000))),
			IsActive:         false,
			CreatedAt:        time.Now().Unix(),
			Creator:          "cosmos1abc...",
		}
		tunnels1 = append(tunnels1, tunnel)
	}
	// mock first response
	queryResponse1 := &tunneltypes.QueryTunnelsResponse{
		Tunnels: tunnels1,
		Pagination: &querytypes.PageResponse{
			NextKey: []byte("next-key"),
		},
	}

	// expect first response from bandQueryClient
	s.tunnelQueryClient.EXPECT().Tunnels(s.ctx, &tunneltypes.QueryTunnelsRequest{
		Pagination: &querytypes.PageRequest{
			Key: nil,
		},
	}).Return(queryResponse1, nil)

	// mock second request
	tunnels2 := make([]*tunneltypes.Tunnel, 0, 100)
	for i := 101; i <= 120; i++ {
		tunnel := &tunneltypes.Tunnel{
			ID:               uint64(i),
			Sequence:         100,
			Route:            any,
			Encoder:          0,
			FeePayer:         "cosmos1xyz...",
			SignalDeviations: []tunneltypes.SignalDeviation{},
			Interval:         60,
			TotalDeposit:     types.NewCoins(types.NewCoin("uband", math.NewInt(1000))),
			IsActive:         false,
			CreatedAt:        time.Now().Unix(),
			Creator:          "cosmos1abc...",
		}
		tunnels2 = append(tunnels2, tunnel)
	}
	// mock second response
	queryResponse2 := &tunneltypes.QueryTunnelsResponse{
		Tunnels: tunnels2,
		Pagination: &querytypes.PageResponse{
			NextKey: []byte(""),
		},
	}

	// expect second response from bandQueryClient
	s.tunnelQueryClient.EXPECT().Tunnels(s.ctx, &tunneltypes.QueryTunnelsRequest{
		Pagination: &querytypes.PageRequest{
			Key: []byte("next-key"),
		},
	}).Return(queryResponse2, nil)

	encodingConfig := band.MakeEncodingConfig()
	s.client = band.NewClient(
		cosmosclient.Context{}.
			WithCodec(encodingConfig.Marshaler).
			WithInterfaceRegistry(encodingConfig.InterfaceRegistry),
		band.NewQueryClient(s.tunnelQueryClient, s.bandtssQueryClient),
		s.log,
		[]string{})

	expected := make([]bandclienttypes.Tunnel, 0, 120)
	for i := 1; i <= 120; i++ {
		expected = append(
			expected,
			*bandclienttypes.NewTunnel(uint64(i), 100, "0xe00F1f85abDB2aF6760759547d450da68CE66Bb1", "eth", false),
		)
	}

	actual, err := s.client.GetTunnels(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal(expected, actual)
}
