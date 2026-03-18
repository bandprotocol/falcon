package wallet

import (
	"fmt"
	"path"

	"github.com/99designs/keyring"
)

const (
	keyringService = "falcon"
	keyringDir     = "keyring"
)

// signerKeyring wraps the shared file-backed keyring and scopes all operations
// to a single chain, so call sites never need to pass chainName explicitly.
type signerKeyring struct {
	kr        keyring.Keyring
	chainName string
}

// newSignerKeyring opens the chain-specific keyring at {homePath}/keys/{chainName}/keyring
// and returns a signerKeyring scoped to that chain.
func newSignerKeyring(passphrase, homePath, chainName string) (*signerKeyring, error) {
	kr, err := keyring.Open(keyring.Config{
		ServiceName:     keyringService,
		AllowedBackends: []keyring.BackendType{keyring.FileBackend},
		FileDir:         path.Join(homePath, "keys", chainName, keyringDir),
		FilePasswordFunc: func(_ string) (string, error) {
			return passphrase, nil
		},
	})
	if err != nil {
		return nil, err
	}
	return &signerKeyring{kr: kr, chainName: chainName}, nil
}

// key returns the keyring key for the given signer name.
// The keyring directory is already chain-scoped, so no prefix is needed.
func (k *signerKeyring) key(name string) string {
	return name
}

// store writes the raw secret for the given signer into the keyring.
func (k *signerKeyring) store(name, secret string) error {
	return k.kr.Set(keyring.Item{
		Key:         k.key(name),
		Data:        []byte(secret),
		Label:       fmt.Sprintf("falcon %s/%s", k.chainName, name),
		Description: "falcon signer secret",
	})
}

// load reads the raw secret for the given signer from the keyring.
func (k *signerKeyring) load(name string) (string, error) {
	item, err := k.kr.Get(k.key(name))
	if err != nil {
		if err == keyring.ErrKeyNotFound {
			return "", fmt.Errorf("missing secret for key %s/%s", k.chainName, name)
		}
		return "", err
	}
	return string(item.Data), nil
}

// delete removes the secret for the given signer from the keyring.
// It is a no-op when the key is not present.
func (k *signerKeyring) delete(name string) error {
	err := k.kr.Remove(k.key(name))
	if err == keyring.ErrKeyNotFound {
		return nil
	}
	return err
}
