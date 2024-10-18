package chains

import (
	"go.uber.org/zap"
)

// ChainProvider defines the interface for the chain interaction with the destination chain.
type ChainProvider interface{}

// BaseChainProvider is a base object for connecting with the chain network.
type BaseChainProvider struct {
	log *zap.Logger

	Config    ChainProviderConfig
	ChainName string
	ChainID   string

	debug bool
}
