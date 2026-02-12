package xrpl

import (
	"time"

	"github.com/bandprotocol/falcon/relayer/alert"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/logger"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

var _ chains.ChainProviderConfig = &XRPLChainProviderConfig{}

// XRPLChainProviderConfig is the configuration for the XRPL chain provider.
type XRPLChainProviderConfig struct {
	chains.BaseChainProviderConfig `mapstructure:",squash"`

	Fee           uint64        `mapstructure:"fee"         toml:"fee"`
	NonceInterval time.Duration `mapstructure:"nonce_interval" toml:"nonce_interval"`
}

// NewChainProvider creates a new XRPL chain provider.
func (cpc *XRPLChainProviderConfig) NewChainProvider(
	chainName string,
	log logger.Logger,
	wallet wallet.Wallet,
	alert alert.Alert,
) (chains.ChainProvider, error) {
	client := NewClient(chainName, cpc, log, alert)

	return NewXRPLChainProvider(chainName, client, cpc, log, wallet, alert), nil
}

// Validate validates the XRPL chain provider configuration.
func (cpc *XRPLChainProviderConfig) Validate() error {
	return nil
}

func (cpc *XRPLChainProviderConfig) GetChainType() types.ChainType {
	return types.ChainTypeXRPL
}
