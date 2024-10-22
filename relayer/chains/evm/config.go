package evm

import (
	"context"

	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/relayer/chains"
)

var _ chains.ChainProviderConfig = &EVMChainProviderConfig{}

// EVMChainProviderConfig is the configuration for the EVM chain provider.
type EVMChainProviderConfig struct {
	chains.BaseChainProviderConfig
}

// NewProvider creates a new EVM chain provider.
func (cpc *EVMChainProviderConfig) NewChainProvider(
	ctx context.Context,
	chainName string,
	log *zap.Logger,
	homePath string,
	debug bool,
) (chains.ChainProvider, error) {
	client := NewClient(chainName, cpc, log)

	return NewEVMChainProvider(ctx, chainName, client, cpc, log)
}

// Validate validates the EVM chain provider configuration.
func (cpc *EVMChainProviderConfig) Validate() error {
	return nil
}
