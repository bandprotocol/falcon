package datasource

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	_ Source = &Web3EIP1559Source{}
	_ Config = &Web3EIP1559SourceConfig{}
)

// Web3EIP1559Source defines a source that retrieves data from an Ethereum node.
type Web3EIP1559Source struct {
	Endpoint string

	Client *ethclient.Client
}

// NewWeb3EIP1559Source returns a new wWeb3EIP1559 source.
func NewWeb3EIP1559Source(endpoint string) (*Web3EIP1559Source, error) {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return nil, err
	}

	return &Web3EIP1559Source{
		Endpoint: endpoint,
		Client:   client,
	}, nil
}

// GetName returns the source name.
func (s Web3EIP1559Source) GetName() string {
	return fmt.Sprintf("Web3EIP1559Source:%s", s.Endpoint)
}

// GetData returns the data from the Ethereum node.
func (s Web3EIP1559Source) GetData(ctx context.Context) (uint64, error) {
	gasPrice, err := s.Client.SuggestGasPrice(ctx)
	if err != nil {
		return 0, err
	}
	return gasPrice.Uint64(), nil
}

// Web3EIP1559SourceConfig defines the configuration for the web3 eip 1559 source.
type Web3EIP1559SourceConfig struct {
	SourceType SourceType `mapstructure:"source_type" toml:"source_type"`
	Endpoint   string     `mapstructure:"endpoint"    toml:"endpoint"`
}

// Validate validates the web3 eip 1559 source configuration.
func (c Web3EIP1559SourceConfig) Validate() error {
	if c.Endpoint == "" {
		return fmt.Errorf("invalid endpoint")
	}

	return nil
}

// NewSource creates a new Web3EIP1559Source object.
func (c Web3EIP1559SourceConfig) NewSource() (Source, error) {
	return NewWeb3EIP1559Source(c.Endpoint)
}
