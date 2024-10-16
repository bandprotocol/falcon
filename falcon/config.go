package falcon

import (
	"fmt"
	"os"
	"reflect"
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

func LoadConfig(file string) (Config, error) {
	byt, err := os.ReadFile(file)
	if err != nil {
		return Config{}, err
	}

	// unmarshall them with Config into struct
	cfgWrapper := &Config{}
	err = toml.Unmarshal(byt, cfgWrapper)
	if err != nil {
		return Config{}, err
	}

	// unmarshall them with raw map[string] into struct
	var rawData map[string]interface{}
	err = toml.Unmarshal(byt, &rawData)
	if err != nil {
		return Config{}, err
	}

	// validate if there is invalid field in toml file
	err = validateConfigFields(rawData, *cfgWrapper)
	if err != nil {
		return Config{}, err
	}

	return *cfgWrapper, nil
}

// Function to validate invalid fields in the TOML
func validateConfigFields(rawData map[string]interface{}, cfg Config) error {
	// Use reflection to get the struct field names
	expectedFields := make(map[string]bool)
	val := reflect.ValueOf(cfg)
	typ := val.Type()

	// Build a set of expected field names from the struct tags
	for i := 0; i < val.NumField(); i++ {
		tag := typ.Field(i).Tag.Get("toml")
		if tag != "" {
			expectedFields[tag] = true
		}
	}

	// Compare the map keys (raw TOML fields) with the expected field names
	for field := range rawData {
		if !expectedFields[field] {
			return fmt.Errorf("invalid field in TOML: %s", field)
		}
	}
	return nil
}
