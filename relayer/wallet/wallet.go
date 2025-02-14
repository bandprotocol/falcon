package wallet

import (
	"crypto/ecdsa"
)

type Wallet interface {
	SavePrivateKey(name string, privKey *ecdsa.PrivateKey) (addr string, err error)
	DeletePrivateKey(name string) error
	ExportPrivateKey(name string) (string, error)
	Sign(name string, data []byte) ([]byte, error)
	GetNames() []string
	GetAddress(name string) (addr string, ok bool)
}
