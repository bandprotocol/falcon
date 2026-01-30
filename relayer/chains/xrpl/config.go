package xrpl

import (
	"fmt"

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

	OracleID   uint16 `mapstructure:"oracle_id"   toml:"oracle_id"`
	Fee        string `mapstructure:"fee"         toml:"fee"`
	PriceScale uint32 `mapstructure:"price_scale" toml:"price_scale"`
}

// NewChainProvider creates a new XRPL chain provider.
func (cpc *XRPLChainProviderConfig) NewChainProvider(
	chainName string,
	log logger.Logger,
	wallet wallet.Wallet,
	alert alert.Alert,
) (chains.ChainProvider, error) {
	client := NewClient(chainName, cpc, log, alert)

	return NewXRPLChainProvider(chainName, client, cpc, log, wallet, alert)
}

// Validate validates the XRPL chain provider configuration.
func (cpc *XRPLChainProviderConfig) Validate() error {
	if len(cpc.Endpoints) == 0 {
		return fmt.Errorf("endpoints is required")
	}
	if cpc.OracleID == 0 {
		return fmt.Errorf("oracle_id is required")
	}
	if cpc.Fee == "" {
		return fmt.Errorf("fee is required")
	}
	return nil
}

func (cpc *XRPLChainProviderConfig) GetChainType() types.ChainType {
	return types.ChainTypeXRPL
}
