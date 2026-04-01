package soroban

import (
	"github.com/bandprotocol/falcon/relayer/alert"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/logger"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

var _ chains.ChainProviderConfig = &SorobanChainProviderConfig{}

// SorobanChainProviderConfig is the configuration for the Soroban chain provider.
type SorobanChainProviderConfig struct {
	chains.BaseChainProviderConfig `mapstructure:",squash"`

	HorizonEndpoint   string `mapstructure:"horizon_endpoint"   toml:"horizon_endpoint"`
	Fee               string `mapstructure:"fee"                toml:"fee"`
	NetworkPassphrase string `mapstructure:"network_passphrase" toml:"network_passphrase"`
}

// NewChainProvider creates a new Soroban chain provider.
func (cpc *SorobanChainProviderConfig) NewChainProvider(
	chainName string,
	log logger.Logger,
	wallet wallet.Wallet,
	alert alert.Alert,
) (chains.ChainProvider, error) {
	client := NewClient(chainName, cpc, log, alert)

	return NewSorobanChainProvider(chainName, client, cpc, log, wallet, alert), nil
}

func (cpc *SorobanChainProviderConfig) GetChainType() types.ChainType {
	return types.ChainTypeSoroban
}
