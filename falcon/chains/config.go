package chains

import (
	"go.uber.org/zap"
)

// ChainProviderConfigs is a collection of ChainProviderConfig interfaces (mapped by chainName)
type ChainProviderConfigs map[string]ChainProviderConfig

// ChainProviderConfig defines the interface for creating a chain provider object.
type ChainProviderConfig interface {
	NewProvider(log *zap.Logger, homePath string, debug bool) (ChainProvider, error)
	Validate() error
}
