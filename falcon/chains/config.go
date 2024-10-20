package chains

import (
	"go.uber.org/zap"
)

// BaseChainProviderConfig contains common field for particular chain provider.
type BaseChainProviderConfig struct {
	RpcEndpoints []string  `toml:"rpc_endpoints"`
	ChainType    ChainType `toml:"chain_type"`
	MaxRetry     int       `toml:"max_retry"`
	ChainID      uint64    `toml:"chain_id"`

	TunnelRouterAddress string `toml:"tunnel_router_address"`
}

// ChainProviderConfigs is a collection of ChainProviderConfig interfaces (mapped by chainName)
type ChainProviderConfigs map[string]ChainProviderConfig

// ChainProviderConfig defines the interface for creating a chain provider object.
type ChainProviderConfig interface {
	NewChainProvider(chainName string, log *zap.Logger, homePath string, debug bool) (ChainProvider, error)
	Validate() error
}
