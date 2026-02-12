package xrpl

import (
	"fmt"
	"path"

	"github.com/99designs/keyring"
)

const (
	xrplKeyringService = "falcon-xrpl"
	xrplKeyringDirName = "priv"
)

func openXRPLKeyring(passphrase, homePath, chainName string) (keyring.Keyring, error) {
	return keyring.Open(keyring.Config{
		ServiceName:     xrplKeyringService,
		AllowedBackends: []keyring.BackendType{keyring.FileBackend},
		FileDir:         path.Join(homePath, "keys", chainName, xrplKeyringDirName),
		FilePasswordFunc: func(_ string) (string, error) {
			return passphrase, nil
		},
	})
}

func xrplKeyringKey(chainName, name string) string {
	return fmt.Sprintf("xrpl/%s/%s", chainName, name)
}

func getXRPLSecret(kr keyring.Keyring, chainName, name string) (string, error) {
	item, err := kr.Get(xrplKeyringKey(chainName, name))
	if err != nil {
		if err == keyring.ErrKeyNotFound {
			return "", fmt.Errorf("missing secret for key %s", name)
		}
		return "", err
	}
	return string(item.Data), nil
}

func setXRPLSecret(kr keyring.Keyring, chainName, name, secret string) error {
	return kr.Set(keyring.Item{
		Key:         xrplKeyringKey(chainName, name),
		Data:        []byte(secret),
		Label:       fmt.Sprintf("XRPL secret %s/%s", chainName, name),
		Description: "XRPL signer seed",
	})
}

func deleteXRPLSecret(kr keyring.Keyring, chainName, name string) error {
	err := kr.Remove(xrplKeyringKey(chainName, name))
	if err == keyring.ErrKeyNotFound {
		return nil
	}
	return err
}
