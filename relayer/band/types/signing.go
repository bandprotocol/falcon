package types

import (
	"time"

	tmbytes "github.com/cometbft/cometbft/libs/bytes"
)

// EVMSignature defines a signature in the EVM format.
type EVMSignature struct {
	RAddress  tmbytes.HexBytes
	Signature tmbytes.HexBytes
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
	ID           uint64
	Message      tmbytes.HexBytes
	Signature    []byte
	EVMSignature *EVMSignature
	CreatedAt    time.Time
}

// NewSigning creates a new Signing instance.
func NewSigning(
	id uint64,
	message tmbytes.HexBytes,
	evmSignature *EVMSignature,
) *Signing {
	return &Signing{
		ID:           id,
		Message:      message,
		EVMSignature: evmSignature,
	}
}
