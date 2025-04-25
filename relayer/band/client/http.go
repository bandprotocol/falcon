package client

import (
	"context"
	"fmt"
	"time"

	"github.com/cometbft/cometbft/rpc/client/http"
	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/internal/bandchain/bandtss"
	"github.com/bandprotocol/falcon/internal/bandchain/tunnel"
	"github.com/bandprotocol/falcon/relayer/band"
	"github.com/bandprotocol/falcon/relayer/band/types"
)

var _ Client = &HttpClient{}

// HttpClient is a BandChain http client struct.
type HttpClient struct {
	Context     cosmosclient.Context
	QueryClient band.QueryClient
	Log         *zap.Logger
	Config      *band.Config
}

// NewClient creates a new BandChain client instance.
func NewClient(queryClient band.QueryClient, log *zap.Logger, bandChainCfg *band.Config) *HttpClient {
	encodingConfig := band.MakeEncodingConfig()
	ctx := cosmosclient.Context{}.
		WithCodec(encodingConfig.Marshaller).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry)

	return &HttpClient{
		Context:     ctx,
		QueryClient: queryClient,
		Log:         log,
		Config:      bandChainCfg,
	}
}

// Init initializes the BandChain client by connecting to the chain and starting
// periodic liveliness checks.
func (c *HttpClient) Init(ctx context.Context) error {
	if err := c.connect(true); err != nil {
		c.Log.Error("Failed to connect to BandChain", zap.Error(err))
		return err
	}

	go c.startLivelinessCheck(ctx)
	return nil
}

// connect connects to the BandChain using the provided RPC endpoints.
func (c *HttpClient) connect(onStartup bool) error {
	timeout := uint(c.Config.Timeout)
	for _, rpcEndpoint := range c.Config.RpcEndpoints {
		// Create a new HTTP client for the specified node URI
		client, err := http.NewWithTimeout(rpcEndpoint, "/websocket", timeout)
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
		c.QueryClient = band.NewBandQueryClient(c.Context)

		c.Log.Info("Connected to BandChain", zap.String("endpoint", rpcEndpoint))

		return nil
	}

	return fmt.Errorf("failed to connect to BandChain")
}

// GetTunnel gets tunnel info from band client
func (c *HttpClient) GetTunnel(ctx context.Context, tunnelID uint64) (types.Tunnel, error) {
	// check connection to bandchain
	if c.QueryClient == nil {
		return types.Tunnel{}, fmt.Errorf("cannot connect to bandchain")
	}

	res, err := c.QueryClient.Tunnel(
		ctx, &tunnel.QueryTunnelRequest{
			TunnelId: tunnelID,
		},
	)
	if err != nil {
		return types.Tunnel{}, err
	}

	if res.Tunnel.Route.TypeUrl != "/band.tunnel.v1beta1.TSSRoute" {
		return types.Tunnel{}, fmt.Errorf("unsupported route type: %s", res.Tunnel.Route.TypeUrl)
	}

	// Extract route information
	var route tunnel.RouteI
	err = c.UnpackAny(res.Tunnel.Route, &route)
	if err != nil {
		return types.Tunnel{}, err
	}

	tssRoute, ok := route.(*tunnel.TSSRoute)
	if !ok {
		return types.Tunnel{}, fmt.Errorf("unsupported route type: %T", route)
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
func (c *HttpClient) GetTunnelPacket(ctx context.Context, tunnelID uint64, sequence uint64) (types.Packet, error) {
	// check connection to bandchain
	if c.QueryClient == nil {
		return types.Packet{}, fmt.Errorf("cannot connect to bandchain")
	}

	// Get packet information by given tunnel ID and sequence
	resPacket, err := c.QueryClient.Packet(
		ctx, &tunnel.QueryPacketRequest{
			TunnelId: tunnelID,
			Sequence: sequence,
		},
	)
	if err != nil {
		return types.Packet{}, err
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
	var packetReceipt tunnel.PacketReceiptI
	err = c.UnpackAny(resPacket.Packet.Receipt, &packetReceipt)
	if err != nil {
		return types.Packet{}, err
	}

	tssPacketReceipt, ok := packetReceipt.(*tunnel.TSSPacketReceipt)
	if !ok {
		return types.Packet{}, fmt.Errorf("unsupported packet content type: %T", packetReceipt)
	}
	signingID := uint64(tssPacketReceipt.SigningID)

	// Get tss signing information by given signing ID
	resSigning, err := c.QueryClient.Signing(
		ctx, &bandtss.QuerySigningRequest{
			SigningId: signingID,
		},
	)
	if err != nil {
		return types.Packet{}, err
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
func (c *HttpClient) GetTunnels(ctx context.Context) ([]types.Tunnel, error) {
	// check connection to bandchain
	if c.QueryClient == nil {
		return nil, fmt.Errorf("cannot connect to BandChain")
	}

	tunnels := make([]types.Tunnel, 0)
	var nextKey []byte

	for {
		res, err := c.QueryClient.Tunnels(
			ctx, &tunnel.QueryTunnelsRequest{
				Pagination: &query.PageRequest{
					Key: nextKey,
				},
			},
		)
		if err != nil {
			return nil, err
		}

		for _, t := range res.Tunnels {
			// Extract route information and filter out non-TSS tunnels
			if t.Route.TypeUrl != "/band.tunnel.v1beta1.TSSRoute" {
				continue
			}

			var route tunnel.RouteI
			if err := c.UnpackAny(t.Route, &route); err != nil {
				return nil, err
			}

			// if not tssRoute, skip this tunnel
			tssRoute, ok := route.(*tunnel.TSSRoute)
			if !ok {
				continue
			}

			tunnels = append(
				tunnels, types.NewTunnel(
					t.ID,
					t.Sequence,
					tssRoute.DestinationContractAddress,
					tssRoute.DestinationChainID,
					t.IsActive,
				),
			)
		}

		nextKey = res.GetPagination().GetNextKey()
		if len(nextKey) == 0 {
			break
		}
	}

	return tunnels, nil
}

// UnpackAny unpacks the provided *codectypes.Any into the specified interface.
func (c *HttpClient) UnpackAny(any *codectypes.Any, target interface{}) error {
	err := c.Context.InterfaceRegistry.UnpackAny(any, target)
	if err != nil {
		return fmt.Errorf("error unpacking into %T: %w", target, err)
	}
	return nil
}

// startLivelinessCheck starts the liveliness check for the BandChain.
func (c *HttpClient) startLivelinessCheck(ctx context.Context) {
	ticker := time.NewTicker(c.Config.LivelinessCheckingInterval)
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
				if err := c.connect(false); err != nil {
					c.Log.Error("Liveliness check: unable to reconnect to any endpoints", zap.Error(err))
				}
			}
		}
	}
}
