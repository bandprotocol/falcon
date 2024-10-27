package types

import (
	"time"

	tmbytes "github.com/cometbft/cometbft/libs/bytes"
)

// EVMSignature defines a signature in the EVM format.
type EVMSignature struct {
	RAddress  tmbytes.HexBytes `json:"r_address"`
	Signature tmbytes.HexBytes `json:"signature"`
}

// NewEVMSignature creates a new EVMSignature instance.
func NewEVMSignature(
	rAddress tmbytes.HexBytes,
	signature tmbytes.HexBytes,
) *EVMSignature {
	return &EVMSignature{
		RAddress:  rAddress,
		Signature: signature,
	}
}

// Signing contains information of a requested message and group signature.
type Signing struct {
	ID           uint64           `json:"tunnel_id"`
	Message      tmbytes.HexBytes `json:"message"`
	Signature    []byte           `json:"signature"`
	EVMSignature *EVMSignature    `json:"evm_signature"`
	CreatedAt    time.Time        `json:"created_at"`
}

// NewSigning creates a new Signing instance.
func NewSigning(
	id uint64,
	message tmbytes.HexBytes,
	evmSignature *EVMSignature,
	createdAt time.Time,
) *Signing {
	return &Signing{
		ID:           id,
		Message:      message,
		EVMSignature: evmSignature,
		CreatedAt:    createdAt,
	}
}
