package icon

import (
	"time"

	"github.com/bandprotocol/falcon/relayer/alert"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/logger"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

var _ chains.ChainProviderConfig = &IconChainProviderConfig{}

// IconChainProviderConfig is the configuration for the Icon chain provider.
type IconChainProviderConfig struct {
	chains.BaseChainProviderConfig `mapstructure:",squash"`

	NetworkID          string        `mapstructure:"network_id"           toml:"network_id"`
	StepLimit          uint64        `mapstructure:"step_limit"           toml:"step_limit"`
	WaitingTxDuration  time.Duration `mapstructure:"waiting_tx_duration"  toml:"waiting_tx_duration"`
	CheckingTxInterval time.Duration `mapstructure:"checking_tx_interval" toml:"checking_tx_interval"`
}

// NewChainProvider creates a new Icon chain provider.
func (cpc *IconChainProviderConfig) NewChainProvider(
	chainName string,
	log logger.Logger,
	wallet wallet.Wallet,
	alert alert.Alert,
) (chains.ChainProvider, error) {
	client := NewClient(chainName, cpc, log, alert)

	return NewIconChainProvider(chainName, client, cpc, log, wallet, alert), nil
}

// Validate validates the Icon chain provider configuration.
func (cpc *IconChainProviderConfig) Validate() error {
	return nil
}

func (cpc *IconChainProviderConfig) GetChainType() types.ChainType {
	return types.ChainTypeIcon
}
