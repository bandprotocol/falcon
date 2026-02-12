package chains

import (
	"context"
	"math/big"

	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
	chainstypes "github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/db"
)

// ChainProviders is a collection of ChainProvider interfaces (mapped by chainName)
type ChainProviders map[string]ChainProvider

// ChainProvider defines the interface for the chain interaction with the destination chain.
type ChainProvider interface {
	// Init initialize to the chain.
	Init(ctx context.Context) error

	// SetDatabase assigns the given database instance.
	SetDatabase(database db.Database)

	// QueryTunnelInfo queries the tunnel information from the destination chain.
	QueryTunnelInfo(
		ctx context.Context,
		tunnelID uint64,
		tunnelDestinationAddr string,
	) (*chainstypes.Tunnel, error)

	// RelayPacket relays the packet from the source chain to the destination chain.
	RelayPacket(ctx context.Context, packet *bandtypes.Packet) error

	// QueryBalance queries balance by given key name from the destination chain.
	QueryBalance(ctx context.Context, keyName string) (*big.Int, error)

	// GetChainName retrieves the chain name from the chain provider.
	GetChainName() string

	// GetChainType retrieves the chain type from the chain provider.
	ChainType() chainstypes.ChainType

	// LoadSigners loads signers to prepare to relay the packet
	LoadSigners() error
}
