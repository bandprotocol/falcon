package band

import (
	"context"
	"fmt"
	"time"

	httpclient "github.com/cometbft/cometbft/rpc/client/http"
	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	querytypes "github.com/cosmos/cosmos-sdk/types/query"
	"go.uber.org/zap"

	bandtsstypes "github.com/bandprotocol/falcon/internal/bandchain/bandtss"
	tunneltypes "github.com/bandprotocol/falcon/internal/bandchain/tunnel"
	"github.com/bandprotocol/falcon/relayer/band/types"
)

var _ Client = &client{}

// Client is the interface to interact with the BandChain.
type Client interface {
	// Init initializes the BandChain client by connecting to the chain and starting
	// periodic liveliness checks.
	Init(ctx context.Context) error

	// GetTunnelPacket returns the packet with the given tunnelID and sequence.
	GetTunnelPacket(ctx context.Context, tunnelID uint64, sequence uint64) (*types.Packet, error)

	// GetTunnel returns the tunnel with the given tunnelID.
	GetTunnel(ctx context.Context, tunnelID uint64) (*types.Tunnel, error)

	// GetTunnels returns all tunnel in BandChain.
	GetTunnels(ctx context.Context) ([]types.Tunnel, error)
}

// client is the BandChain client struct.
type client struct {
	Context     cosmosclient.Context
	QueryClient QueryClient
	Log         *zap.Logger
	Config      *Config
}

// NewClient creates a new BandChain client instance.
func NewClient(queryClient QueryClient, log *zap.Logger, bandChainCfg *Config) Client {
	encodingConfig := MakeEncodingConfig()
	ctx := cosmosclient.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry)

	return &client{
		Context:     ctx,
		QueryClient: queryClient,
		Log:         log,
		Config:      bandChainCfg,
	}
}

// Init initializes the BandChain client by connecting to the chain and starting
// periodic liveliness checks.
func (c *client) Init(ctx context.Context) error {
	timeout := uint(c.Config.Timeout)
	if err := c.connect(timeout, true); err != nil {
		c.Log.Error("Failed to connect to BandChain", zap.Error(err))
		return err
	}

	go c.startLivelinessCheck(ctx, timeout, c.Config.LivelinessCheckingInterval)
	return nil
}

// connect connects to the BandChain using the provided RPC endpoints.
func (c *client) connect(timeout uint, onStartup bool) error {
	for _, rpcEndpoint := range c.Config.RpcEndpoints {
		// Create a new HTTP client for the specified node URI
		client, err := httpclient.NewWithTimeout(rpcEndpoint, "/websocket", timeout)
		if err != nil {
			c.Log.Error("Failed to create HTTP client", zap.String("rpcEndpoint", rpcEndpoint), zap.Error(err))
			continue // Try the next endpoint if there's an error
		}

		// skip status check on startup to avoid blocking relayer initialization
		// perform status checks later to ensure endpoint health and rotation
		if !onStartup {
			if _, err := c.Context.Client.Status(context.Background()); err != nil {
				continue
			}
		}

		c.Context.Client = client
		c.Context.NodeURI = rpcEndpoint
		c.QueryClient = NewBandQueryClient(c.Context)

		c.Log.Info("Connected to BandChain", zap.String("endpoint", rpcEndpoint))

		return nil
	}

	return fmt.Errorf("failed to connect to BandChain")
}

// startLivelinessCheck starts the liveliness check for the BandChain.
func (c *client) startLivelinessCheck(ctx context.Context, timeout uint, interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ctx.Done():
			c.Log.Info("Stopping liveliness check")

			ticker.Stop()

			return
		case <-ticker.C:
			if _, err := c.Context.Client.Status(ctx); err != nil {
				c.Log.Error(
					"BandChain client disconnected",
					zap.String("rpcEndpoint", c.Context.NodeURI),
					zap.Error(err),
				)
				if err := c.connect(timeout, false); err != nil {
					c.Log.Error("Liveliness check: unable to reconnect to any endpoints", zap.Error(err))
				}
			}
		}
	}
}

// GetTunnel gets tunnel info from band client
func (c *client) GetTunnel(ctx context.Context, tunnelID uint64) (*types.Tunnel, error) {
	// check connection to bandchain
	if c.QueryClient == nil {
		return nil, fmt.Errorf("cannot connect to bandchain")
	}

	res, err := c.QueryClient.Tunnel(ctx, &tunneltypes.QueryTunnelRequest{
		TunnelId: tunnelID,
	})
	if err != nil {
		return nil, err
	}

	// Extract route information
	var route tunneltypes.RouteI
	err = c.UnpackAny(res.Tunnel.Route, &route)
	if err != nil {
		return nil, err
	}

	tssRoute, ok := route.(*tunneltypes.TSSRoute)
	if !ok {
		return nil, fmt.Errorf("unsupported route type: %T", route)
	}

	return types.NewTunnel(
		res.Tunnel.ID,
		res.Tunnel.Sequence,
		tssRoute.DestinationContractAddress,
		tssRoute.DestinationChainID,
		res.Tunnel.IsActive,
	), nil
}

// GetTunnelPacket gets tunnel packet info from band client
func (c *client) GetTunnelPacket(ctx context.Context, tunnelID uint64, sequence uint64) (*types.Packet, error) {
	// check connection to bandchain
	if c.QueryClient == nil {
		return nil, fmt.Errorf("cannot connect to bandchain")
	}

	// Get packet information by given tunnel ID and sequence
	resPacket, err := c.QueryClient.Packet(ctx, &tunneltypes.QueryPacketRequest{
		TunnelId: tunnelID,
		Sequence: sequence,
	})
	if err != nil {
		return nil, err
	}

	// Convert resPacket.Packet.SignalPrices to []types.SignalPrice
	signalPrices := make([]types.SignalPrice, len(resPacket.Packet.Prices))
	for i, sp := range resPacket.Packet.Prices {
		signalPrices[i] = types.SignalPrice{
			SignalID: sp.SignalID,
			Price:    sp.Price,
		}
	}

	// Extract tunnel packet information
	var packetReceipt tunneltypes.PacketReceiptI
	err = c.UnpackAny(resPacket.Packet.Receipt, &packetReceipt)
	if err != nil {
		return nil, err
	}

	tssPacketReceipt, ok := packetReceipt.(*tunneltypes.TSSPacketReceipt)
	if !ok {
		return nil, fmt.Errorf("unsupported packet content type: %T", packetReceipt)
	}
	signingID := uint64(tssPacketReceipt.SigningID)

	// Get tss signing information by given signing ID
	resSigning, err := c.QueryClient.Signing(ctx, &bandtsstypes.QuerySigningRequest{
		SigningId: signingID,
	})
	if err != nil {
		return nil, err
	}

	currentGroupSigning := types.ConvertSigning(resSigning.CurrentGroupSigningResult)
	incomingGroupSigning := types.ConvertSigning(resSigning.IncomingGroupSigningResult)

	return types.NewPacket(
		resPacket.Packet.TunnelID,
		resPacket.Packet.Sequence,
		signalPrices,
		currentGroupSigning,
		incomingGroupSigning,
	), nil
}

// GetTunnels returns every tss-route tunnels in BandChain.
func (c *client) GetTunnels(ctx context.Context) ([]types.Tunnel, error) {
	// check connection to bandchain
	if c.QueryClient == nil {
		return nil, fmt.Errorf("cannot connect to BandChain")
	}

	tunnels := make([]types.Tunnel, 0)
	var nextKey []byte

	for {
		res, err := c.QueryClient.Tunnels(ctx, &tunneltypes.QueryTunnelsRequest{
			Pagination: &querytypes.PageRequest{
				Key: nextKey,
			},
		})
		if err != nil {
			return nil, err
		}

		for _, tunnel := range res.Tunnels {
			// Extract route information
			var route tunneltypes.RouteI
			if err := c.UnpackAny(tunnel.Route, &route); err != nil {
				return nil, err
			}

			// if not tssRoute, skip this tunnel
			tssRoute, ok := route.(*tunneltypes.TSSRoute)
			if !ok {
				continue
			}

			tunnels = append(tunnels, *types.NewTunnel(
				tunnel.ID,
				tunnel.Sequence,
				tssRoute.DestinationContractAddress,
				tssRoute.DestinationChainID,
				tunnel.IsActive,
			))
		}

		nextKey = res.GetPagination().GetNextKey()
		if len(nextKey) == 0 {
			break
		}
	}

	return tunnels, nil
}

// UnpackAny unpacks the provided *codectypes.Any into the specified interface.
func (c *client) UnpackAny(any *codectypes.Any, target interface{}) error {
	err := c.Context.InterfaceRegistry.UnpackAny(any, target)
	if err != nil {
		return fmt.Errorf("error unpacking into %T: %w", target, err)
	}
	return nil
}
