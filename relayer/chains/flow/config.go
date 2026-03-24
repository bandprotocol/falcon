package flow

import (
	"time"

	"github.com/bandprotocol/falcon/relayer/alert"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/logger"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

var _ chains.ChainProviderConfig = &FlowChainProviderConfig{}

// FlowChainProviderConfig is the configuration for the Flow chain provider.
type FlowChainProviderConfig struct {
	chains.BaseChainProviderConfig `mapstructure:",squash"`

	ComputeLimit       uint64        `mapstructure:"compute_limit"        toml:"compute_limit"`
	WaitingTxDuration  time.Duration `mapstructure:"waiting_tx_duration"  toml:"waiting_tx_duration"`
	CheckingTxInterval time.Duration `mapstructure:"checking_tx_interval" toml:"checking_tx_interval"`
}

// NewChainProvider creates a new Flow chain provider.
func (cpc *FlowChainProviderConfig) NewChainProvider(
	chainName string,
	log logger.Logger,
	w wallet.Wallet,
	a alert.Alert,
) (chains.ChainProvider, error) {
	c := NewClient(chainName, cpc, log, a)
	cp, err := NewFlowChainProvider(chainName, c, cpc, log, w, a)
	if err != nil {
		return nil, err
	}

	return cp, nil
}

// GetChainType returns the chain type for Flow.
func (cpc *FlowChainProviderConfig) GetChainType() types.ChainType {
	return types.ChainTypeFlow
}
