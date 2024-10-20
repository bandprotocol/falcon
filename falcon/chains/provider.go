package chains

import (
	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/falcon/chains/types"
)

// ChainProvider defines the interface for the chain interaction with the destination chain.
type ChainProvider interface {
	// Connect connects to the chain.
	Connect() error

	// QueryTunnelInfo queries the tunnel information from the destination chain.
	QueryTunnelInfo(tunnelID uint64, tunnelDestinationAddr string) (*types.Tunnel, error)
}

// BaseChainProvider is a base object for connecting with the chain network.
type BaseChainProvider struct {
	log *zap.Logger

	Config    ChainProviderConfig
	ChainName string
	ChainID   string

	debug bool
}
