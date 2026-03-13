package wallet

import "fmt"

// WalletAdapter defines chain-specific wallet operations.
// Implement this interface to add support for a new chain type.
type WalletAdapter interface {
	// LoadSigner reconstructs a Signer from a persisted KeyRecord using the secret
	// retrieved from the shared keyring.
	// For local signers, secret is the raw private key or mnemonic.
	// For remote signers, secret is the API key (empty string means no key).
	LoadSigner(name string, record KeyRecord, secret string) (Signer, error)
}

// RemoteSignerFactory is a function that constructs a chain-specific remote Signer
// from the fields stored in a KeyRecord.
type RemoteSignerFactory func(name, address, url string, key string) (Signer, error)

// RemoteOnlyAdapter is a WalletAdapter for chain types that only support remote signing.
// Add a new remote-only chain by defining a RemoteSigner that implements Sign, then
// passing its constructor to NewRemoteOnlyAdapter — no custom Adapter type is needed.
type RemoteOnlyAdapter struct {
	newRemote RemoteSignerFactory
}

var _ WalletAdapter = RemoteOnlyAdapter{}

// NewRemoteOnlyAdapter creates a RemoteOnlyAdapter backed by the given factory.
func NewRemoteOnlyAdapter(f RemoteSignerFactory) RemoteOnlyAdapter {
	return RemoteOnlyAdapter{newRemote: f}
}

func (a RemoteOnlyAdapter) LoadSigner(name string, record KeyRecord, secret string) (Signer, error) {
	switch record.Type {
	case RemoteSignerType:
		return a.newRemote(name, record.Address, record.URL, secret)
	case PrivKeySignerType, MnemonicSignerType:
		return nil, fmt.Errorf("chain does not support local key import")
	default:
		return nil, fmt.Errorf("unsupported signer type: %s for key %s", record.Type, name)
	}
}
