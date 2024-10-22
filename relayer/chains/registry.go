package chains

import "fmt"

// Registry is a collection of chain clients.
type Registry struct {
	Chains map[string]Client
}

// NewRegistry creates a new chain registry.
func NewRegistry() *Registry {
	return &Registry{
		Chains: make(map[string]Client),
	}
}

// Register registers a chain client to the registry.
func (r *Registry) Register(chainID string, client Client) error {
	if _, ok := r.Chains[chainID]; !ok {
		return fmt.Errorf("chain %s already registered", chainID)
	}

	r.Chains[chainID] = client
	return nil
}
