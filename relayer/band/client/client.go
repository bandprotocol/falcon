package client

import (
	"context"

	"github.com/bandprotocol/falcon/relayer/band/types"
)

// Client is the interface to interact with the BandChain.
type Client interface {
	// Init initializes the BandChain client by connecting to the chain and starting
	// periodic liveliness checks.
	Init(ctx context.Context) error

	// GetTunnelPacket returns the packet with the given tunnelID and sequence.
	GetTunnelPacket(ctx context.Context, tunnelID uint64, sequence uint64) (types.Packet, error)

	// GetTunnel returns the tunnel with the given tunnelID.
	GetTunnel(ctx context.Context, tunnelID uint64) (types.Tunnel, error)

	// GetTunnels returns all tunnel in BandChain.
	GetTunnels(ctx context.Context) ([]types.Tunnel, error)
}
