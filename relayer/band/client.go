package band

import (
	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/relayer/band/types"
)

var _ Client = &client{}

// Client is the interface to interact with the BandChain.
type Client interface {
	// GetTunnelPacket returns the packet with the given tunnelID and sequence.
	GetTunnelPacket(tunnelID uint64, sequence uint64) (*types.Packet, error)

	// GetTunnel returns the tunnel with the given tunnelID.
	GetTunnel(tunnelID uint64) (*types.Tunnel, error)

	// GetSigning returns the signing with the given signingID.
	GetSigning(signingID uint64) (*types.Signing, error)
}

// client is the BandChain client struct.
type client struct {
	Log          *zap.Logger
	RpcEndpoints []string
}

// NewClient creates a new BandChain client instance.
func NewClient(log *zap.Logger, rpcEndpoints []string) Client {
	return &client{
		Log:          log,
		RpcEndpoints: rpcEndpoints,
	}
}

func (c *client) GetTunnelPacket(tunnelID uint64, sequence uint64) (*types.Packet, error) {
	return nil, nil
}

func (c *client) GetTunnel(tunnelID uint64) (*types.Tunnel, error) {
	return nil, nil
}

func (c *client) GetSigning(signingID uint64) (*types.Signing, error) {
	return nil, nil
}
