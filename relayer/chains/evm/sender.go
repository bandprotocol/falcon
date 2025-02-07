package evm

import (
	"crypto/ecdsa"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// Sender is the struct that represents the sender of the transaction.
type Sender struct {
	PrivateKey *ecdsa.PrivateKey
	Address    gethcommon.Address
}

// NewSender creates a new sender object.
func NewSender(privateKey *ecdsa.PrivateKey, address gethcommon.Address) *Sender {
	return &Sender{
		PrivateKey: privateKey,
		Address:    address,
	}
}

// LoadFreeSenders loads the FreeSenders channel with sender instances.
// derived from the keys stored in the keystore located at the specified homePath.
func (cp *EVMChainProvider) LoadFreeSenders() error {
	if cp.FreeSenders != nil {
		return nil
	}

	keyNames := cp.Wallet.GetNames()
	freeSenders := make(chan *Sender, len(keyNames))

	for _, keyName := range keyNames {
		key, err := cp.Wallet.GetKey(keyName)
		if err != nil {
			return err
		}

		addr := gethcommon.HexToAddress(key.Address)
		freeSenders <- NewSender(key.PrivateKey, addr)
	}

	cp.FreeSenders = freeSenders
	return nil
}
