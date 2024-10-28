package band

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	httpclient "github.com/cometbft/cometbft/rpc/client/http"

	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/bandprotocol/falcon/relayer/band/params"
	"github.com/bandprotocol/falcon/relayer/band/types"

	tsstypes "github.com/bandprotocol/chain/v3/x/tss/types"
	tunneltypes "github.com/bandprotocol/chain/v3/x/tunnel/types"
)

var _ Client = &client{}

// Client is the interface to interact with the BandChain.
type Client interface {
	// GetTunnelPacket returns the packet with the given tunnelID and sequence.
	GetTunnelPacket(ctx context.Context, tunnelID uint64, sequence uint64) (*types.Packet, error)

	// GetTunnel returns the tunnel with the given tunnelID.
	GetTunnel(ctx context.Context, tunnelID uint64) (*types.Tunnel, error)

	// Connect will establish connection to rpc endpoints
	Connect(timeout uint) error
}

// client is the BandChain client struct.
type client struct {
	Context      cosmosclient.Context
	QueryClient  BandQueryClient
	Log          *zap.Logger
	RpcEndpoints []string
}

// NewClient creates a new BandChain client instance.
func NewClient(ctx cosmosclient.Context, queryClient BandQueryClient, log *zap.Logger, rpcEndpoints []string) Client {
	return &client{
		Context:      ctx,
		QueryClient:  queryClient,
		Log:          log,
		RpcEndpoints: rpcEndpoints,
	}
}

// Connect connects to the Band chain using the provided RPC endpoints.
func (c *client) Connect(timeout uint) error {
	for _, rpcEndpoint := range c.RpcEndpoints {
		// Create a new HTTP client for the specified node URI
		client, err := httpclient.NewWithTimeout(rpcEndpoint, "/websocket", timeout)
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
		encodingConfig := params.MakeEncodingConfig()
		ctx := cosmosclient.Context{}.
			WithClient(client).
			WithCodec(encodingConfig.Marshaler).
			WithInterfaceRegistry(encodingConfig.InterfaceRegistry)

		c.Context = ctx
		c.QueryClient = *NewQueryClientProvider(ctx)

		return nil
	}

	return nil
}

// GetTunnel gets tunnel info from band client
func (c *client) GetTunnel(ctx context.Context, tunnelID uint64) (*types.Tunnel, error) {
	res, err := c.QueryClient.Tunnel(ctx, tunnelID)
	if err != nil {
		return nil, err
	}

	// Extract route information
	var route tunneltypes.RouteI
	err = c.UnpackAny(res.Tunnel.Route, &route)
	if err != nil {
		return nil, err
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
		res.Tunnel.IsActive,
	), nil
}

// GetTunnelPacket gets tunnel packet info from band client
func (c *client) GetTunnelPacket(ctx context.Context, tunnelID uint64, sequence uint64) (*types.Packet, error) {
	// Get packet information by given tunnel ID and sequence
	resPacket, err := c.QueryClient.Packet(ctx, tunnelID, sequence)
	if err != nil {
		return nil, err
	}

	// Extract tunnel packet information
	var pc tunneltypes.PacketContentI
	err = c.UnpackAny(resPacket.Packet.PacketContent, &pc)
	if err != nil {
		return nil, err
	}

	var signingID uint64
	switch tmp := pc.(type) {
	case *tunneltypes.TSSPacketContent:
		signingID = uint64(tmp.SigningID)
	default:
		return nil, fmt.Errorf("unsupported packet content type: %T", tmp)
	}

	// Get tss signing information by given signing ID
	resSigning, err := c.QueryClient.Signing(ctx, signingID)
	if err != nil {
		return nil, err
	}

	var signingResult *tsstypes.SigningResult

	// Handle the possible values of CurrentGroupSigningResult
	switch {
	case resSigning.CurrentGroupSigningResult != nil:
		signingResult = resSigning.CurrentGroupSigningResult
	case resSigning.IncomingGroupSigningResult != nil:
		signingResult = resSigning.IncomingGroupSigningResult
	default:
		return nil, fmt.Errorf("no signing result available")
	}

	signing := signingResult.Signing
	evmSignature := signingResult.EVMSignature

	// Convert resPacket.Packet.SignalPrices to []types.SignalPrice
	signalPrices := make([]types.SignalPrice, len(resPacket.Packet.SignalPrices))
	for i, sp := range resPacket.Packet.SignalPrices {
		signalPrices[i] = types.SignalPrice{
			SignalID: sp.SignalID,
			Price:    sp.Price,
		}
	}

	return types.NewPacket(
		resPacket.Packet.TunnelID,
		resPacket.Packet.Sequence,
		signalPrices,
		types.NewSigningInfo(
			signingID,
			signing.Message,
			types.NewEVMSignature(
				evmSignature.RAddress,
				evmSignature.Signature,
			),
		),
	), nil
}

// unpackAny unpacks the provided *codectypes.Any into the specified interface.
func (c *client) UnpackAny(any *codectypes.Any, target interface{}) error {
	err := c.Context.InterfaceRegistry.UnpackAny(any, target)
	if err != nil {
		return fmt.Errorf("error unpacking into %T: %w", target, err)
	}
	return nil
}
