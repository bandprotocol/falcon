package types

import (
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
type SigningInfo struct {
	ID           uint64           `json:"id"`
	Message      tmbytes.HexBytes `json:"messsage"`
	EVMSignature *EVMSignature    `json:"evm_signature"`
}

// NewSigning creates a new Signing instance.
func NewSigningInfo(
	id uint64,
	message tmbytes.HexBytes,
	evmSignature *EVMSignature,
) *SigningInfo {
	return &SigningInfo{
		ID:           id,
		Message:      message,
		EVMSignature: evmSignature,
	}
}
