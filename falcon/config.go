package falcon

import (
	"os"
	"time"

	"github.com/pelletier/go-toml/v2"

	"github.com/bandprotocol/falcon/falcon/band"
	"github.com/bandprotocol/falcon/falcon/chains"
)

// TargetChainConfig is the metadata of the configured target chain.
type TargetChainConfig struct {
	Name       string           `toml:"id"`
	Type       chains.ChainType `toml:"type"`
	ConfigPath string           `toml:"path"`
}

// Config defines the configuration for the falcon tunnel relayer.
type Config struct {
	BandChainConfig        band.Config         `toml:"bandchain"`
	TargetChains           []TargetChainConfig `toml:"target_chains"`
	CheckingPacketInterval time.Duration       `toml:"checking_packet_interval"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		BandChainConfig: band.Config{
			RpcEndpoints: []string{"http://localhost:26657"},
			Timeout:      5,
		},
		TargetChains:           []TargetChainConfig{},
		CheckingPacketInterval: time.Minute,
	}
}

// LoadConfig reads config file from given path and return config object
func LoadConfig(cfgPath string) (*Config, error) {
	byt, err := os.ReadFile(cfgPath)
	if err != nil {
		return &Config{}, err
	}

	// unmarshall them with Config into struct
	cfg := &Config{}
	err = toml.Unmarshal(byt, cfg)
	if err != nil {
		return &Config{}, err
	}
	return cfg, nil
}
