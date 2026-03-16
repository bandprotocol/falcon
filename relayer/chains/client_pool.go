package chains

import (
	"fmt"
	"sync"
)

// ClientPool is a thread-safe pool of RPC clients keyed by endpoint URL.
// T is the underlying client type (e.g. *ethclient.Client, *rpc.Client).
type ClientPool[T any] struct {
	mu               sync.RWMutex
	selectedEndpoint string
	clients          map[string]*T
}

// NewClientPool creates and returns a new ClientPool with no entries.
func NewClientPool[T any]() ClientPool[T] {
	return ClientPool[T]{
		clients: make(map[string]*T),
	}
}

// GetClient returns the client for a given endpoint and whether it exists.
func (p *ClientPool[T]) GetClient(endpoint string) (*T, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	client, exists := p.clients[endpoint]
	return client, exists
}

// SetClient stores the client for a given endpoint.
func (p *ClientPool[T]) SetClient(endpoint string, client *T) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.clients[endpoint] = client
}

// SetSelectedEndpoint sets the currently active endpoint.
func (p *ClientPool[T]) SetSelectedEndpoint(endpoint string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.selectedEndpoint = endpoint
}

// GetSelectedEndpoint returns the currently active endpoint.
func (p *ClientPool[T]) GetSelectedEndpoint() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.selectedEndpoint
}

// GetSelectedClient returns the client for the selected endpoint.
// Returns an error if no endpoint is selected or if the selected client does not exist.
func (p *ClientPool[T]) GetSelectedClient() (*T, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.selectedEndpoint == "" {
		return nil, fmt.Errorf("no selected endpoint")
	}

	selectedClient, exists := p.clients[p.selectedEndpoint]
	if !exists {
		return nil, fmt.Errorf("selected endpoint client not found: %s", p.selectedEndpoint)
	}

	return selectedClient, nil
}
