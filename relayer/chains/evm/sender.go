package evm

import (
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// Sender is the struct that represents the sender of the transaction.
type Sender struct {
	Name    string
	Address gethcommon.Address
}

// NewSender creates a new sender object.
func NewSender(name string, address gethcommon.Address) *Sender {
	return &Sender{
		Name:    name,
		Address: address,
	}
}

// LoadFreeSenders loads the FreeSenders channel with sender instances.
// derived from the keys stored in the keystore located.
func (cp *EVMChainProvider) LoadFreeSenders() error {
	if cp.FreeSenders != nil {
		return nil
	}

	keyNames := cp.Wallet.GetNames()
	freeSenders := make(chan *Sender, len(keyNames))

	for _, keyName := range keyNames {
		addrHex, ok := cp.Wallet.GetAddress(keyName)
		if !ok {
			return fmt.Errorf("key name does not exist: %s", keyName)
		}

		addr := gethcommon.HexToAddress(addrHex)
		freeSenders <- NewSender(keyName, addr)
	}

	cp.FreeSenders = freeSenders
	return nil
}
