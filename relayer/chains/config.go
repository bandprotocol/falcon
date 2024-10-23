package chains

import (
	"time"

	"go.uber.org/zap"
)

// BaseChainProviderConfig contains common field for particular chain provider.
type BaseChainProviderConfig struct {
	Endpoints    []string      `mapstructure:"endpoints"     toml:"endpoints"`
	ChainType    ChainType     `mapstructure:"chain_type"    toml:"chain_type"`
	MaxRetry     int           `mapstructure:"max_retry"     toml:"max_retry"`
	QueryTimeout time.Duration `mapstructure:"query_timeout" toml:"query_timeout"`
	ChainID      uint64        `mapstructure:"chain_id"      toml:"chain_id"`

	TunnelRouterAddress string `mapstructure:"tunnel_router_address" toml:"tunnel_router_address"`
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