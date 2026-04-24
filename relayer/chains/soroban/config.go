package soroban

import (
	"time"

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

	HorizonEndpoints   []string      `mapstructure:"horizon_endpoints"    toml:"horizon_endpoints"`
	Fee                string        `mapstructure:"fee"                  toml:"fee"`
	NetworkPassphrase  string        `mapstructure:"network_passphrase"   toml:"network_passphrase"`
	WaitingTxDuration  time.Duration `mapstructure:"waiting_tx_duration"  toml:"waiting_tx_duration"`
	CheckingTxInterval time.Duration `mapstructure:"checking_tx_interval" toml:"checking_tx_interval"`
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
