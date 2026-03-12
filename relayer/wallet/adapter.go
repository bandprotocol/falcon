package wallet

import "fmt"

// WalletAdapter defines chain-specific wallet operations.
// Implement this interface to add support for a new chain type.
type WalletAdapter interface {
	// DeriveFromPrivateKey parses the private key and creates a signer
	// without persisting the secret.
	DeriveFromPrivateKey(name, privateKey string) (Signer, error)

	// DeriveFromMnemonic derives from a mnemonic and creates a signer
	// without persisting the secret.
	DeriveFromMnemonic(
		name, mnemonic string,
		coinType uint32,
		account uint,
		index uint,
	) (Signer, error)

	// PersistKey stores the secret in chain-specific secure storage.
	// The secret parameter carries the original input (private key or mnemonic).
	PersistKey(name string, signer Signer, secret string) error

	// LoadSigner reconstructs a Signer from a persisted KeyRecord, handling both local and remote types.
	LoadSigner(name string, record KeyRecord) (Signer, error)

	// DeleteLocalSecret removes the locally stored secret for the named key.
	// It should be a no-op for remote signers.
	DeleteLocalSecret(name string, signer Signer) error
}

// RemoteSignerFactory is a function that constructs a chain-specific remote Signer
// from the fields stored in a KeyRecord.
type RemoteSignerFactory func(name, address, url string, key *string) (Signer, error)

// RemoteOnlyAdapter is a complete WalletAdapter for chain types that only support
// remote signing. Add a new remote-only chain by defining a RemoteSigner that
// implements Sign, then passing its constructor to NewRemoteOnlyAdapter — no
// custom Adapter type is needed.
type RemoteOnlyAdapter struct {
	newRemote RemoteSignerFactory
}

var _ WalletAdapter = RemoteOnlyAdapter{}

// NewRemoteOnlyAdapter creates a RemoteOnlyAdapter backed by the given factory.
func NewRemoteOnlyAdapter(f RemoteSignerFactory) RemoteOnlyAdapter {
	return RemoteOnlyAdapter{newRemote: f}
}

func (RemoteOnlyAdapter) DeriveFromPrivateKey(name, privateKey string) (Signer, error) {
	return nil, fmt.Errorf("chain does not support private key import")
}

func (RemoteOnlyAdapter) DeriveFromMnemonic(name, mnemonic string, coinType uint32, account, index uint) (Signer, error) {
	return nil, fmt.Errorf("chain does not support mnemonic import")
}

func (RemoteOnlyAdapter) PersistKey(name string, signer Signer, secret string) error {
	return nil
}

func (a RemoteOnlyAdapter) LoadSigner(name string, record KeyRecord) (Signer, error) {
	if record.Type != RemoteSignerType {
		return nil, fmt.Errorf("unsupported signer type: %s for key %s", record.Type, name)
	}
	return a.newRemote(name, record.Address, record.Url, record.Key)
}

func (RemoteOnlyAdapter) DeleteLocalSecret(name string, signer Signer) error {
	return nil
}
