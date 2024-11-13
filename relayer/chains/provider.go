package chains

import (
	"context"

	"go.uber.org/zap"

	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
	chainstypes "github.com/bandprotocol/falcon/relayer/chains/types"
)

// ChainProvider defines the interface for the chain interaction with the destination chain.
type ChainProvider interface {
	// Init initialize to the chain.
	Init(ctx context.Context) error

	// QueryTunnelInfo queries the tunnel information from the destination chain.
	QueryTunnelInfo(
		ctx context.Context,
		tunnelID uint64,
		tunnelDestinationAddr string,
	) (*chainstypes.Tunnel, error)

	// RelayPacket relays the packet from the source chain to the destination chain.
	RelayPacket(ctx context.Context, packet *bandtypes.Packet) error
}

// BaseChainProvider is a base object for connecting with the chain network.
type BaseChainProvider struct {
	log *zap.Logger

	Config    ChainProviderConfig
	ChainName string
	ChainID   string

	debug bool
}
