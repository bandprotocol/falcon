package evm

import (
	"time"

	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/logger"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

var _ chains.ChainProviderConfig = &EVMChainProviderConfig{}

// EVMChainProviderConfig is the configuration for the EVM chain provider.
type EVMChainProviderConfig struct {
	chains.BaseChainProviderConfig `mapstructure:",squash"`

	BlockConfirmation          uint64        `mapstructure:"block_confirmation"           toml:"block_confirmation"`
	WaitingTxDuration          time.Duration `mapstructure:"waiting_tx_duration"          toml:"waiting_tx_duration"`
	LivelinessCheckingInterval time.Duration `mapstructure:"liveliness_checking_interval" toml:"liveliness_checking_interval"`
	CheckingTxInterval         time.Duration `mapstructure:"checking_tx_interval"         toml:"checking_tx_interval"`
	GasLimit                   uint64        `mapstructure:"gas_limit"                    toml:"gas_limit,omitempty"`

	GasType         GasType       `mapstructure:"gas_type"          toml:"gas_type"`
	MaxGasPrice     uint64        `mapstructure:"max_gas_price"     toml:"max_gas_price,omitempty"`
	MaxBaseFee      uint64        `mapstructure:"max_base_fee"      toml:"max_base_fee,omitempty"`
	MaxPriorityFee  uint64        `mapstructure:"max_priority_fee"  toml:"max_priority_fee,omitempty"`
	GasMultiplier   float64       `mapstructure:"gas_multiplier"    toml:"gas_multiplier"`
	QueryGasTimeout time.Duration `mapstructure:"query_gas_timeout" toml:"query_gas_timeout,omitempty"`
}

// NewProvider creates a new EVM chain provider.
func (cpc *EVMChainProviderConfig) NewChainProvider(
	chainName string,
	log logger.Logger,
	wallet wallet.Wallet,
) (chains.ChainProvider, error) {
	client := NewClient(chainName, cpc, log)

	return NewEVMChainProvider(chainName, client, cpc, log, wallet)
}

// Validate validates the EVM chain provider configuration.
func (cpc *EVMChainProviderConfig) Validate() error {
	return nil
}

func (cpc *EVMChainProviderConfig) GetChainType() types.ChainType {
	return types.ChainTypeEVM
}
