package secret

import (
	"time"

	"github.com/bandprotocol/falcon/relayer/alert"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/logger"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

var _ chains.ChainProviderConfig = &SecretChainProviderConfig{}

const bech32PrefixAccAddr = "secret"

// SecretChainProviderConfig is the configuration for the Secret (CosmWasm) chain provider.
type SecretChainProviderConfig struct {
	chains.BaseChainProviderConfig `mapstructure:",squash"`

	// Cosmos SDK chain id used for signing (passed to fkms).
	CosmosChainID string `mapstructure:"cosmos_chain_id" toml:"cosmos_chain_id"`

	// GasPrice and denom are used for gas/fee and balance queries.
	GasPrice string `mapstructure:"gas_price" toml:"gas_price"`
	Denom    string `mapstructure:"denom" toml:"denom"`

	GasLimitBase uint64 `mapstructure:"gas_limit_base" toml:"gas_limit_base"`
	GasLimitEach uint64 `mapstructure:"gas_limit_each" toml:"gas_limit_each"`

	// Secret contract encryption parameters (passed to fkms).
	CodeHash string `mapstructure:"code_hash" toml:"code_hash"`
	// Secret chain's public key (used for signing, passed to fkms).
	ChainPubkey string `mapstructure:"chain_pubkey" toml:"chain_pubkey"`

	WaitingTxDuration  time.Duration `mapstructure:"waiting_tx_duration" toml:"waiting_tx_duration"`
	CheckingTxInterval time.Duration `mapstructure:"checking_tx_interval" toml:"checking_tx_interval"`
}

func (cpc *SecretChainProviderConfig) NewChainProvider(
	chainName string,
	log logger.Logger,
	wallet wallet.Wallet,
	alert alert.Alert,
) (chains.ChainProvider, error) {
	client := NewClient(chainName, cpc, log, alert)

	return NewSecretChainProvider(chainName, client, cpc, log, wallet, alert), nil
}

func (cpc *SecretChainProviderConfig) GetChainType() types.ChainType {
	return types.ChainTypeSecret
}
