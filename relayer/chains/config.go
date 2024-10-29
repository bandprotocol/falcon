package chains

import (
	"time"

	"go.uber.org/zap"
)

// BaseChainProviderConfig contains common field for particular chain provider.
type BaseChainProviderConfig struct {
	Endpoints    []string      `toml:"endpoints"`
	ChainType    ChainType     `toml:"chain_type"`
	MaxRetry     int           `toml:"max_retry"`
	QueryTimeout time.Duration `toml:"query_timeout"`
	ChainID      uint64        `toml:"chain_id"`

	TunnelRouterAddress string `toml:"tunnel_router_address"`
}

// ChainProviderConfigs is a collection of ChainProviderConfig interfaces (mapped by chainName)
type ChainProviderConfigs map[string]ChainProviderConfig

// ChainProviders is a collection of ChainProvider interfaces (mapped by chainName)
type ChainProviders map[string]ChainProvider

// ChainProviderConfig defines the interface for creating a chain provider object.
type ChainProviderConfig interface {
	NewChainProvider(
		chainName string,
		log *zap.Logger,
		homePath string,
		debug bool,
	) (ChainProvider, error)
	Validate() error
}
