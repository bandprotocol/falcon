package config

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pelletier/go-toml/v2"

	"github.com/bandprotocol/falcon/relayer/band"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/chains/evm"
	chainstypes "github.com/bandprotocol/falcon/relayer/chains/types"
)

// ChainProviderConfigs is a collection of ChainProviderConfig interfaces (mapped by chainName)
type ChainProviderConfigs map[string]chains.ChainProviderConfig

// GlobalConfig is the global configuration for the falcon tunnel relayer
type GlobalConfig struct {
	LogLevel               string        `mapstructure:"log_level"                toml:"log_level"`
	CheckingPacketInterval time.Duration `mapstructure:"checking_packet_interval" toml:"checking_packet_interval"`
	SyncTunnelsInterval    time.Duration `mapstructure:"sync_tunnels_interval"    toml:"sync_tunnels_interval"`
	PenaltySkipRounds      uint          `mapstructure:"penalty_skip_rounds"      toml:"penalty_skip_rounds"`
	MetricsListenAddr      string        `mapstructure:"metrics_listen_addr"      toml:"metrics_listen_addr"`
	DBPath                 string        `mapstructure:"db_path"                  toml:"db_path"`
}

// Config defines the configuration for the falcon tunnel relayer.
type Config struct {
	Global       GlobalConfig         `mapstructure:"global"        toml:"global"`
	BandChain    band.Config          `mapstructure:"bandchain"     toml:"bandchain"`
	TargetChains ChainProviderConfigs `mapstructure:"target_chains" toml:"target_chains"`
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
	chainType := chainstypes.ToChainType(typeName)

	decoderConfig := mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			decodeTimeHook,
			chainstypes.DecodeChainTypeHook,
			evm.DecodeGasTypeHook,
		),
	}

	var cfg chains.ChainProviderConfig
	switch chainType {
	case chainstypes.ChainTypeEVM:
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

// ParseConfigInputWrapper converts a ConfigInputWrapper object to a Config object.
func ParseConfigInputWrapper(wrappedCfg *ConfigInputWrapper) (*Config, error) {
	targetChains := make(ChainProviderConfigs)
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
			RpcEndpoints:               []string{"http://localhost:26657"},
			Timeout:                    3 * time.Second,
			LivelinessCheckingInterval: 5 * time.Minute,
		},
		TargetChains: make(map[string]chains.ChainProviderConfig),
		Global: GlobalConfig{
			LogLevel:               "info",
			CheckingPacketInterval: time.Minute,
			PenaltySkipRounds:      3,
			SyncTunnelsInterval:    5 * time.Minute,
		},
	}
}

func ParseConfig(data []byte) (*Config, error) {
	var cfgWrapper ConfigInputWrapper
	if err := DecodeConfigInputWrapperTOML(data, &cfgWrapper); err != nil {
		return nil, err
	}

	// convert ConfigWrapperInput to Config
	cfg, err := ParseConfigInputWrapper(&cfgWrapper)
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
	b, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	var chainProviderConfigWrapper ChainProviderConfigWrapper
	// unmarshal them with Config into struct
	err = toml.Unmarshal(b, &chainProviderConfigWrapper)
	if err != nil {
		return nil, err
	}

	chainProviderConfig, err := ParseChainProviderConfig(chainProviderConfigWrapper)
	if err != nil {
		return nil, err
	}

	return chainProviderConfig, nil
}
