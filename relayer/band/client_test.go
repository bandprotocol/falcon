package band_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	tmbytes "github.com/cometbft/cometbft/libs/bytes"

	"github.com/cosmos/gogoproto/proto"

	"cosmossdk.io/math"

	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/falcon/internal/relayertest/mocks"
	"github.com/bandprotocol/falcon/relayer/band"
	"github.com/bandprotocol/falcon/relayer/band/params"
	bandclienttypes "github.com/bandprotocol/falcon/relayer/band/types"

	bandtsstypes "github.com/bandprotocol/chain/v3/x/bandtss/types"
	tsstypes "github.com/bandprotocol/chain/v3/x/tss/types"
	tunneltypes "github.com/bandprotocol/chain/v3/x/tunnel/types"
)

type AppTestSuite struct {
	suite.Suite

	ctx             context.Context
	bandQueryClient *mocks.MockBandQueryClient
	client          band.Client
	log             *zap.Logger
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
	s.bandQueryClient = mocks.NewMockBandQueryClient(ctrl)
	s.ctx = context.Background()
}

func (s *AppTestSuite) TestGetTunnelTssRoute() {
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
		Interval:         60, // 60 seconds
		TotalDeposit:     types.NewCoins(types.NewCoin("uband", math.NewInt(1000))),
		IsActive:         false,
		CreatedAt:        time.Now().Unix(),
		Creator:          "cosmos1abc...",
	}
	queryResponse := &tunneltypes.QueryTunnelResponse{
		Tunnel: tunnel,
	}

	s.bandQueryClient.EXPECT().Tunnel(s.ctx, uint64(1)).Return(queryResponse, nil)
	encodingConfig := params.MakeEncodingConfig()
	s.client = band.NewClient(
		cosmosclient.Context{}.
			WithCodec(encodingConfig.Marshaler).
			WithInterfaceRegistry(encodingConfig.InterfaceRegistry),
		s.bandQueryClient,
		s.log,
		[]string{})

	expected := bandclienttypes.NewTunnel(1, 100, "0xe00F1f85abDB2aF6760759547d450da68CE66Bb1", "eth")

	actual, err := s.client.GetTunnel(s.ctx, uint64(1))
	s.Require().NoError(err)
	s.Require().Equal(expected, actual)
}

func (s *AppTestSuite) TestGetTunnelAxelarRoute() {
	// mock query response
	destinationChainID := "eth"
	destinationContractAddress := "0xe00F1f85abDB2aF6760759547d450da68CE66Bb1"
	r := &tunneltypes.AxelarRoute{
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
		Interval:         60, // 60 seconds
		TotalDeposit:     types.NewCoins(types.NewCoin("uband", math.NewInt(1000))),
		IsActive:         false,
		CreatedAt:        time.Now().Unix(),
		Creator:          "cosmos1abc...",
	}
	queryResponse := &tunneltypes.QueryTunnelResponse{
		Tunnel: tunnel,
	}
	s.bandQueryClient.EXPECT().Tunnel(s.ctx, uint64(1)).Return(queryResponse, nil)

	// expected result
	expected := bandclienttypes.NewTunnel(1, 100, "0xe00F1f85abDB2aF6760759547d450da68CE66Bb1", "eth")

	// actual result
	encodingConfig := params.MakeEncodingConfig()
	s.client = band.NewClient(
		cosmosclient.Context{}.
			WithCodec(encodingConfig.Marshaler).
			WithInterfaceRegistry(encodingConfig.InterfaceRegistry),
		s.bandQueryClient,
		s.log,
		[]string{})

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
		TunnelID:      1,
		Sequence:      100,
		SignalPrices:  []tunneltypes.SignalPrice{},
		PacketContent: any,
		CreatedAt:     time.Now().Unix(),
	}

	queryResponse := &tunneltypes.QueryPacketResponse{
		Packet: &packet,
	}
	s.bandQueryClient.EXPECT().Packet(s.ctx, uint64(1), uint64(100)).Return(queryResponse, nil)

	// expected result
	expected := bandclienttypes.NewPacket(1, 100, 2)

	// actual result
	encodingConfig := params.MakeEncodingConfig()
	s.client = band.NewClient(
		cosmosclient.Context{}.
			WithCodec(encodingConfig.Marshaler).
			WithInterfaceRegistry(encodingConfig.InterfaceRegistry),
		s.bandQueryClient,
		s.log,
		[]string{})
	actual, err := s.client.GetTunnelPacket(s.ctx, uint64(1), uint64(100))
	s.Require().NoError(err)
	s.Require().Equal(expected, actual)
}

func (s *AppTestSuite) TestGetSigningCurrentGroup() {
	// mock query response
	message := tmbytes.HexBytes{0x7a, 0x8b, 0x9c}
	rAddress := tmbytes.HexBytes{0x1a, 0x2b, 0x3c}
	signature := tmbytes.HexBytes{0x4d, 0x5e, 0x6f}
	signing := tsstypes.Signing{
		ID:      1,
		Message: message,
	}
	evmSignature := &tsstypes.EVMSignature{
		RAddress:  rAddress,
		Signature: signature,
	}
	queryResponse := &bandtsstypes.QuerySigningResponse{
		CurrentGroupSigningResult: &tsstypes.SigningResult{
			Signing:      signing,
			EVMSignature: evmSignature,
		},
		IncomingGroupSigningResult: nil,
	}
	s.bandQueryClient.EXPECT().Signing(s.ctx, uint64(1)).Return(queryResponse, nil)

	// expected result
	expected := bandclienttypes.NewSigning(1, message, bandclienttypes.NewEVMSignature(rAddress, signature))

	// actual result
	encodingConfig := params.MakeEncodingConfig()
	s.client = band.NewClient(
		cosmosclient.Context{}.
			WithCodec(encodingConfig.Marshaler).
			WithInterfaceRegistry(encodingConfig.InterfaceRegistry),
		s.bandQueryClient,
		s.log,
		[]string{})
	actual, err := s.client.GetSigning(s.ctx, uint64(1))
	s.Require().NoError(err)
	s.Require().Equal(expected, actual)
}

func (s *AppTestSuite) TestGetSigningIncomingGroup() {
	// mock query response
	message := tmbytes.HexBytes{0x7a, 0x8b, 0x9c}
	rAddress := tmbytes.HexBytes{0x1a, 0x2b, 0x3c}
	signature := tmbytes.HexBytes{0x4d, 0x5e, 0x6f}
	signing := tsstypes.Signing{
		ID:      1,
		Message: message,
	}
	evmSignature := &tsstypes.EVMSignature{
		RAddress:  rAddress,
		Signature: signature,
	}
	queryResponse := &bandtsstypes.QuerySigningResponse{
		CurrentGroupSigningResult: nil,
		IncomingGroupSigningResult: &tsstypes.SigningResult{
			Signing:      signing,
			EVMSignature: evmSignature,
		},
	}
	s.bandQueryClient.EXPECT().Signing(s.ctx, uint64(1)).Return(queryResponse, nil)

	// expected result
	expected := bandclienttypes.NewSigning(1, message, bandclienttypes.NewEVMSignature(rAddress, signature))

	// actual result
	encodingConfig := params.MakeEncodingConfig()
	s.client = band.NewClient(
		cosmosclient.Context{}.
			WithCodec(encodingConfig.Marshaler).
			WithInterfaceRegistry(encodingConfig.InterfaceRegistry),
		s.bandQueryClient,
		s.log,
		[]string{})
	actual, err := s.client.GetSigning(s.ctx, uint64(1))
	s.Require().NoError(err)
	s.Require().Equal(expected, actual)
}
