package keys

// Key represents a key in the keystore.
type Key struct {
	Name    string
	ChainID string
	PrivKey string
	PubKey  string
}

// NewKey creates a new key instance.
func NewKey(name, chainID, privKey, pubKey string) *Key {
	return &Key{
		Name:    name,
		ChainID: chainID,
		PrivKey: privKey,
		PubKey:  pubKey,
	}
}
