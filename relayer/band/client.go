package band

import (
	"context"

	"fmt"

	"go.uber.org/zap"

	"github.com/bandprotocol/chain/v3/app/params"
	"github.com/bandprotocol/falcon/relayer/band/types"

	httpclient "github.com/cometbft/cometbft/rpc/client/http"
	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	tunneltypes "github.com/bandprotocol/chain/v3/x/tunnel/types"
	tsstypes "github.com/bandprotocol/chain/v3/x/tss/types"
	bandtsstypes "github.com/bandprotocol/chain/v3/x/bandtss/types"
	"github.com/cosmos/cosmos-sdk/std"
)

var _ Client = &client{}

// Client is the interface to interact with the BandChain.
type Client interface {
	// GetTunnelPacket returns the packet with the given tunnelID and sequence.
	GetTunnelPacket(ctx context.Context, tunnelID uint64, sequence uint64) (*types.Packet, error)

	// GetTunnel returns the tunnel with the given tunnelID.
	GetTunnel(ctx context.Context, tunnelID uint64) (*types.Tunnel, error)

	// GetSigning returns the signing with the given signingID.
	GetSigning(ctx context.Context, signingID uint64) (*types.Signing, error)

	// Connect will establish connection to rpc endpoints
	Connect() error
}

// client is the BandChain client struct.
type client struct {
	Context      cosmosclient.Context
	QueryClient  queryClient
	Log          *zap.Logger
	RpcEndpoints []string

}

type queryClient struct {
	TunnelClient tunneltypes.QueryClient
	BandtssClient bandtsstypes.QueryClient
}

// EncodingConfig specifies the concrete encoding types to use for a given app.
// This is provided for compatibility between protobuf and amino implementations.
type EncodingConfig struct {
	InterfaceRegistry codectypes.InterfaceRegistry
	Marshaler         codec.Codec
	TxConfig          cosmosclient.TxConfig
	Amino             *codec.LegacyAmino
}

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() params.EncodingConfig {
	encodingConfig := params.MakeEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}

// NewClient creates a new BandChain client instance.
func NewClient(log *zap.Logger, rpcEndpoints []string) Client {
	return &client{
		Log:          log,
		RpcEndpoints: rpcEndpoints,
	}
}

// Connect connects to the Band chain using the provided RPC endpoints.
func (c *client) Connect() error {
	for _, rpcEndpoint := range c.RpcEndpoints {
		// Create a new HTTP client for the specified node URI
		client, err := httpclient.New(rpcEndpoint, "/websocket")
		if err != nil {
			c.Log.Error("Failed to create HTTP client", zap.String("rpcEndpoint", rpcEndpoint), zap.Error(err))
			continue // Try the next endpoint if there's an error
		}

		// Start the client to establish a connection
		if err = client.Start(); err != nil {
			c.Log.Error("Failed to start client", zap.String("rpcEndpoint", rpcEndpoint), zap.Error(err))
			continue // Try the next endpoint if starting the client fails
		}

		// Create a new client context and configure it with necessary parameters
		encodingConfig := MakeEncodingConfig()
		ctx := cosmosclient.Context{}.
			WithClient(client).
			WithChainID("bandchain").
			WithCodec(encodingConfig.Marshaler).
			WithInterfaceRegistry(encodingConfig.InterfaceRegistry)
		
		c.Context = ctx
		// ...
		c.QueryClient.TunnelClient = tunneltypes.NewQueryClient(ctx)
		c.QueryClient.BandtssClient = bandtsstypes.NewQueryClient(ctx)
		
		return nil
	}

	return nil
}

func (c *client) GetTunnelPacket(ctx context.Context, tunnelID uint64, sequence uint64) (*types.Packet, error) {
	queryClient := c.QueryClient.TunnelClient

	res, err := queryClient.Packet(ctx, &tunneltypes.QueryPacketRequest{
		TunnelId: tunnelID,
		Sequence: sequence,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query tunnel packet: %w", err)
	}

	// Extract tunnel packet information
	var pc tunneltypes.PacketContentI
	err = c.UnpackAny(res.Packet.PacketContent, &pc)
	if err != nil {
		return nil, fmt.Errorf("failed to extract packet content info: %w", err)
	}

	var sigingID uint64
	switch tmp := pc.(type) {
	case *tunneltypes.TSSPacketContent:
		sigingID = uint64(tmp.SigningID)
	case *tunneltypes.AxelarPacketContent:
		sigingID = tmp.IBCQueueID
	default:
		return nil, fmt.Errorf("unsupported packet content type: %T", tmp)
	}	

	return &types.Packet{
		TunnelID:  res.Packet.TunnelID,
		Sequence:  res.Packet.Sequence,
		SigningID: sigingID,
	}, nil
}

func (c *client) GetTunnel(ctx context.Context, tunnelID uint64) (*types.Tunnel, error) {
	queryClient := c.QueryClient.TunnelClient

	// mockTunnelClient := mocks.New...
	// c.querClient.TunnelClient = mockTunnel

	res, err := queryClient.Tunnel(ctx, &tunneltypes.QueryTunnelRequest{
		TunnelId: tunnelID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tunnel: %w", err)
	}

	// Extract route information
	var route tunneltypes.RouteI
	err = c.UnpackAny(res.Tunnel.Route, &route)
	if err != nil {
		return nil, fmt.Errorf("failed to extract route info: %w", err)
	}

	var targetChainID, targetAddress string
	switch tmp := route.(type) {
	case *tunneltypes.TSSRoute:
		targetChainID = tmp.DestinationChainID
		targetAddress = tmp.DestinationContractAddress
	case *tunneltypes.AxelarRoute:
		targetChainID = tmp.DestinationChainID
		targetAddress = tmp.DestinationContractAddress
	default:
		return nil, fmt.Errorf("unsupported route type: %T", route)
	}	

	return types.NewTunnel(
		res.Tunnel.ID, 
		res.Tunnel.Sequence, 
		targetAddress, 
		targetChainID,
	), nil
}

func (c *client) GetSigning(ctx context.Context, signingID uint64) (*types.Signing, error) {
	queryClient := c.QueryClient.BandtssClient

	res, err := queryClient.Signing(context.Background(), &bandtsstypes.QuerySigningRequest{
		SigningId: signingID,
	})
	if err != nil {
		return nil, err
	}

	var signingResult *tsstypes.SigningResult

	// Handle the possible values of CurrentGroupSigningResult
	switch {
	case res.CurrentGroupSigningResult != nil:
		signingResult = res.CurrentGroupSigningResult
	case res.IncomingGroupSigningResult != nil:
		signingResult = res.IncomingGroupSigningResult
	default:
		return nil, fmt.Errorf("no signing result available")
	}

	signing := signingResult.Signing
	evmSignature := signingResult.EVMSignature

	return &types.Signing{
		ID:        uint64(signing.ID),
		Message:   signing.Message,
		Signature: signing.Signature,
		EVMSignature: &types.EVMSignature{
			RAddress:  evmSignature.RAddress,
			Signature: evmSignature.Signature,
		},
		CreatedAt: signing.CreatedTimestamp,
	}, nil
}

// unpackAny unpacks the provided *codectypes.Any into the specified interface.
func (c *client) UnpackAny(any *codectypes.Any, target interface{}) error {
	err := c.Context.InterfaceRegistry.UnpackAny(any, target)
	if err != nil {
		return fmt.Errorf("error unpacking into %T: %w", target, err)
	}
	return nil
}
