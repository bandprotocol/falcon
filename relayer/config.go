package relayer

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pelletier/go-toml/v2"

	"github.com/bandprotocol/falcon/internal/datasource"
	"github.com/bandprotocol/falcon/relayer/band"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/chains/evm"
	evmgas "github.com/bandprotocol/falcon/relayer/chains/evm/gas"
)

// GlobalConfig is the global configuration for the falcon tunnel relayer
type GlobalConfig struct {
	LogLevel                         string        `mapstructure:"log_level"                            toml:"log_level"`
	CheckingPacketInterval           time.Duration `mapstructure:"checking_packet_interval"             toml:"checking_packet_interval"`
	MaxCheckingPacketPenaltyDuration time.Duration `mapstructure:"max_checking_packet_penalty_duration" toml:"max_checking_packet_penalty_duration"`
	PenaltyExponentialFactor         float64       `mapstructure:"penalty_exponential_factor"           toml:"penalty_exponential_factor"`
}

// Config defines the configuration for the falcon tunnel relayer.
type Config struct {
	Global       GlobalConfig                `mapstructure:"global"        toml:"global"`
	BandChain    band.Config                 `mapstructure:"bandchain"     toml:"bandchain"`
	TargetChains chains.ChainProviderConfigs `mapstructure:"target_chains" toml:"target_chains"`
}

// ChainProviderConfigWrapper is an intermediary type for parsing any object from config.toml file
type ChainProviderConfigWrapper map[string]interface{}

// ConfigInputWrapper is an intermediary type for parsing the config.toml file
type ConfigInputWrapper struct {
	Global       GlobalConfig                          `mapstructure:"global"`
	BandChain    band.Config                           `mapstructure:"bandchain"`
	TargetChains map[string]ChainProviderConfigWrapper `mapstructure:"target_chains"`
}

// ParseChainProviderConfig converts a TOMLWrapper object to a ChainProviderConfig object.
func ParseChainProviderConfig(w ChainProviderConfigWrapper) (chains.ChainProviderConfig, error) {
	typeName, ok := w["chain_type"].(string)
	if !ok {
		return nil, fmt.Errorf("chain_type is required")
	}
	chainType := chains.ToChainType(typeName)

	decoderConfig := mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			decodeTimeHook,
			chains.DecodeChainTypeHook,
			evmgas.DecodeGasTypeHook,
			datasource.DecodeDataSourceConfigHook,
		),
	}

	var cfg chains.ChainProviderConfig
	switch chainType {
	case chains.ChainTypeEVM:
		var newCfg evm.EVMChainProviderConfig

		decoderConfig.Result = &newCfg
		decoder, err := mapstructure.NewDecoder(&decoderConfig)
		if err != nil {
			return nil, err
		}

		if err := decoder.Decode(w); err != nil {
			return nil, err
		}

		cfg = &newCfg
	default:
		return cfg, fmt.Errorf("unsupported chain type: %s", typeName)
	}

	return cfg, nil
}

// DecodeConfigInputWrapperTOML decodes a TOML bytes into a ConfigInputWrapper object.
func DecodeConfigInputWrapperTOML(data []byte, cw *ConfigInputWrapper) error {
	var input map[string]interface{}
	err := toml.Unmarshal(data, &input)
	if err != nil {
		return err
	}

	decoderConfig := mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			decodeTimeHook,
		),
		Result: cw,
	}

	decoder, err := mapstructure.NewDecoder(&decoderConfig)
	if err != nil {
		return err
	}

	if err := decoder.Decode(input); err != nil {
		return err
	}

	return nil
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
			Timeout:      3 * time.Second,
		},
		TargetChains: make(map[string]chains.ChainProviderConfig),
		Global: GlobalConfig{
			LogLevel:                         "info",
			CheckingPacketInterval:           time.Minute,
			MaxCheckingPacketPenaltyDuration: time.Hour,
			PenaltyExponentialFactor:         1.0,
		},
	}
}

// LoadConfig reads config file from given path and return config object
func LoadConfig(cfgPath string) (*Config, error) {
	b, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	var cfgWrapper ConfigInputWrapper
	if err := DecodeConfigInputWrapperTOML(b, &cfgWrapper); err != nil {
		return nil, err
	}

	// convert ConfigWrapperInput to Config
	cfg, err := ParseConfig(&cfgWrapper)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// decodeTimeHook is a custom function to decode time.Duration using mapstructure
func decodeTimeHook(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
	if to == reflect.TypeOf(time.Duration(0)) && from.Kind() == reflect.String {
		return time.ParseDuration(data.(string))
	}

	return data, nil
}

// LoadChainConfig reads chain config file from given path and return chain config object
func LoadChainConfig(cfgPath string) (chains.ChainProviderConfig, error) {
	byt, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	var chainProviderConfigWrapper ChainProviderConfigWrapper
	// unmarshal them with Config into struct
	err = toml.Unmarshal(byt, &chainProviderConfigWrapper)
	if err != nil {
		return nil, err
	}

	chainProviderConfig, err := ParseChainProviderConfig(chainProviderConfigWrapper)
	if err != nil {
		return nil, err
	}

	return chainProviderConfig, nil
}
