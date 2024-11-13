package evm

import "go.uber.org/zap"

// Key represents a key in the keystore.
type Key struct {
	Name    string
	ChainID string
	PrivKey string
	PubKey  string
}

// NewKey creates a new key instance.
func NewKey(name, chainID, privKey, pubKey string) Key {
	return Key{
		Name:    name,
		ChainID: chainID,
		PrivKey: privKey,
		PubKey:  pubKey,
	}
}

// KeyStore stores information for interacting with the key store.
type KeyStore struct {
	Log  *zap.Logger
	Path string
}

// NewKeyStore creates a new key store instance.
func NewKeyStore(log *zap.Logger, cfgPath string) KeyStore {
	return KeyStore{
		Log:  log,
		Path: cfgPath,
	}
}

// SaveKeys saves keys to the store.
func (ks *KeyStore) SaveKeys(keys []Key, cfgPath string) error {
	return nil
}

// GetKeys gets keys from the store.
func (ks *KeyStore) GetKeys() ([]Key, error) {
	return nil, nil
}

// GetKey gets a key by name from the store.
func (ks *KeyStore) GetKey(name string) (Key, error) {
	return Key{}, nil
}
