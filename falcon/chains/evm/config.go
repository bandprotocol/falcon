package evm

import (
	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/falcon/chains"
)

var _ chains.ChainProviderConfig = &EVMChainProviderConfig{}

// EVMChainProviderConfig is the configuration for the EVM chain provider.
type EVMChainProviderConfig struct {
	chains.BaseChainProviderConfig
}

// NewProvider creates a new EVM chain provider.
func (cpc *EVMChainProviderConfig) NewChainProvider(
	chainName string,
	log *zap.Logger,
	homePath string,
	debug bool,
) (chains.ChainProvider, error) {
	client := NewClient(chainName, cpc.RpcEndpoints, log)

	return NewEVMChainProvider(chainName, client, cpc, log)
}

// Validate validates the EVM chain provider configuration.
func (cpc *EVMChainProviderConfig) Validate() error {
	return nil
}
