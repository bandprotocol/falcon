package falcon

import (
	"time"

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
func DefaultConfig() Config {
	return Config{
		BandChainConfig: band.Config{
			RpcEndpoints: []string{"http://localhost:26657"},
			Timeout:      5,
		},
		TargetChains:           []TargetChainConfig{},
		CheckingPacketInterval: time.Minute,
	}
}
