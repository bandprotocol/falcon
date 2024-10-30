package evm

import (
	"crypto/ecdsa"
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Sender is the struct that represents the sender of the transaction.
type Sender struct {
	privateKey *ecdsa.PrivateKey
	address    gethcommon.Address
}

// NewSender creates a new sender object.
func NewSender(privateKeyHex string) (Sender, error) {
	// if private key is empty, return empty object
	if privateKeyHex == "" {
		return Sender{}, nil
	}

	// Convert the private key hex string to a private key object
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return Sender{}, err
	}

	// Get the public key from the private key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return Sender{}, fmt.Errorf("cannot assert type to *ecdsa.PublicKey")
	}

	return Sender{
		privateKey: privateKey,
		address:    crypto.PubkeyToAddress(*publicKeyECDSA),
	}, nil
}
