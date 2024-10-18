package evm

import (
	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/falcon/chains"
)

var _ chains.ChainProviderConfig = &EVMProviderConfig{}

// EVMProviderConfig is the configuration for the EVM chain provider.
type EVMProviderConfig struct {
	RpcEndpoints []string         `toml:"rpc_endpoints"`
	ChainType    chains.ChainType `toml:"chain_type"`
}

// NewProvider creates a new EVM chain provider.
func (pc *EVMProviderConfig) NewProvider(
	log *zap.Logger,
	homePath string,
	debug bool,
) (chains.ChainProvider, error) {
	return EVMChainProvider{}, nil
}

// Validate validates the EVM chain provider configuration.
func (pc *EVMProviderConfig) Validate() error {
	return nil
}
