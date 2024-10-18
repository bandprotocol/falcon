package falcon

import (
	"fmt"
	"os"
	"time"

	"github.com/pelletier/go-toml/v2"

	"github.com/bandprotocol/falcon/falcon/band"
	"github.com/bandprotocol/falcon/falcon/chains"
	"github.com/bandprotocol/falcon/falcon/chains/evm"
)

// GlobalConfig is the global configuration for the falcon tunnel relayer
type GlobalConfig struct {
	CheckingPacketInterval time.Duration `toml:"checking_packet_interval"`
}

// Config defines the configuration for the falcon tunnel relayer.
type Config struct {
	Global       GlobalConfig                `toml:"global"`
	BandChain    band.Config                 `toml:"bandchain"`
	TargetChains chains.ChainProviderConfigs `toml:"target_chains"`
}

// TOMLWrapper is an intermediary type for parsing any object from config.toml file
type TOMLWrapper map[string]interface{}

// ConfigInputWrapper is an intermediary type for parsing the config.toml file
type ConfigInputWrapper struct {
	Global       GlobalConfig           `toml:"global"`
	BandChain    band.Config            `toml:"bandchain"`
	TargetChains map[string]TOMLWrapper `toml:"target_chains"`
}

// ParseChainProviderConfig converts a TOMLWrapper object to a ChainProviderConfig object.
func ParseChainProviderConfig(w TOMLWrapper) (chains.ChainProviderConfig, error) {
	typeName, ok := w["chain_type"].(string)
	if !ok {
		return nil, fmt.Errorf("type field is required")
	}
	chainType := chains.ToChainType(typeName)

	b, err := toml.Marshal(w)
	if err != nil {
		return nil, err
	}

	var cfg chains.ChainProviderConfig
	switch chainType {
	case chains.ChainTypeEVM:
		var newCfg evm.EVMProviderConfig
		if err := toml.Unmarshal(b, &newCfg); err != nil {
			return nil, err
		}
		cfg = &newCfg
	default:
		return cfg, fmt.Errorf("unsupported chain type: %s", chainType)
	}

	return cfg, nil
}

// ParseConfig converts a ConfigInputWrapper object to a Config object.
func ParseConfig(wrappedCfg *ConfigInputWrapper) (*Config, error) {
	targetChains := make(chains.ChainProviderConfigs)
	for name, provCfg := range wrappedCfg.TargetChains {
		newProvCfg, err := ParseChainProviderConfig(provCfg)
		if err != nil {
			return nil, err
		}
		targetChains[name] = newProvCfg
	}

	return &Config{
		Global:       wrappedCfg.Global,
		BandChain:    wrappedCfg.BandChain,
		TargetChains: targetChains,
	}, nil
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		BandChain: band.Config{
			RpcEndpoints: []string{"http://localhost:26657"},
			Timeout:      5,
		},
		TargetChains: make(map[string]chains.ChainProviderConfig),
		Global:       GlobalConfig{CheckingPacketInterval: time.Minute},
	}
}

// LoadConfig reads config file from given path and return config object
func LoadConfig(cfgPath string) (*Config, error) {
	byt, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	// unmarshal them with Config into struct
	cfgWrapper := &ConfigInputWrapper{}
	err = toml.Unmarshal(byt, cfgWrapper)
	if err != nil {
		return nil, err
	}

	cfg, err := ParseConfig(cfgWrapper)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
