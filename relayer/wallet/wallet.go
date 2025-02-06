package wallet

import (
	"crypto/ecdsa"
)

type Key struct {
	Address    string
	PrivateKey *ecdsa.PrivateKey
}

type Wallet interface {
	SavePrivateKey(name string, privKey *ecdsa.PrivateKey) (addr string, err error)
	DeletePrivateKey(name string) error
	GetNames() []string
	GetAddress(name string) (addr string, ok bool)
	GetKey(name string) (key *Key, err error)
}
