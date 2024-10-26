package datasource

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	_ Source = &Web3LegacySource{}
	_ Config = &Web3LegacySourceConfig{}
)

// Web3LegacySource defines a source that retrieves data from an Ethereum node.
type Web3LegacySource struct {
	Endpoint string

	Client *ethclient.Client
}

// NewWeb3LegacySource returns a new Web3Legacy source.
func NewWeb3LegacySource(endpoint string) (*Web3LegacySource, error) {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return nil, err
	}

	return &Web3LegacySource{
		Endpoint: endpoint,
		Client:   client,
	}, nil
}

// GetName returns the source name.
func (s Web3LegacySource) GetName() string {
	return fmt.Sprintf("Web3LegacySource:%s", s.Endpoint)
}

// GetData returns the data from the Ethereum node.
func (s Web3LegacySource) GetData(ctx context.Context) (uint64, error) {
	gasPrice, err := s.Client.SuggestGasPrice(ctx)
	if err != nil {
		return 0, err
	}
	return gasPrice.Uint64(), nil
}

// Web3LegacySourceConfig
type Web3LegacySourceConfig struct {
	SourceType SourceType `mapstructure:"source_type" toml:"source_type"`
	Endpoint   string     `mapstructure:"endpoint"    toml:"endpoint"`
}

// Validate validates the web3 legacy source configuration.
func (c Web3LegacySourceConfig) Validate() error {
	if c.Endpoint == "" {
		return fmt.Errorf("invalid endpoint")
	}

	return nil
}

// NewSource creates a new Web3LegacySource object.
func (c Web3LegacySourceConfig) NewSource() (Source, error) {
	return NewWeb3LegacySource(c.Endpoint)
}
