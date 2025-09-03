package chains

import (
	"time"

	"github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/logger"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

// BaseChainProviderConfig contains common field for particular chain provider.
type BaseChainProviderConfig struct {
	Endpoints      []string        `mapstructure:"endpoints"       toml:"endpoints"`
	ChainType      types.ChainType `mapstructure:"chain_type"      toml:"chain_type"`
	MaxRetry       int             `mapstructure:"max_retry"       toml:"max_retry"`
	QueryTimeout   time.Duration   `mapstructure:"query_timeout"   toml:"query_timeout"`
	ExecuteTimeout time.Duration   `mapstructure:"execute_timeout" toml:"execute_timeout"`
	ChainID        uint64          `mapstructure:"chain_id"        toml:"chain_id"`

	TunnelRouterAddress string `mapstructure:"tunnel_router_address" toml:"tunnel_router_address"`
}

// ChainProviderConfigs is a collection of ChainProviderConfig interfaces (mapped by chainName)
type ChainProviderConfigs map[string]ChainProviderConfig

// ChainProviderConfig defines the interface for creating a chain provider object.
type ChainProviderConfig interface {
	NewChainProvider(
		chainName string,
		log logger.Logger,
		wallet wallet.Wallet,
	) (ChainProvider, error)

	GetChainType() types.ChainType
	Validate() error
}
