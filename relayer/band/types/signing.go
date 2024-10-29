package types

import (
	tmbytes "github.com/cometbft/cometbft/libs/bytes"

	tsstypes "github.com/bandprotocol/chain/v3/x/tss/types"
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
	ID           uint64           `json:"id"`
	Message      tmbytes.HexBytes `json:"messsage"`
	EVMSignature *EVMSignature    `json:"evm_signature"`
}

// ConvertSigning converts tsstypes.SigningResult and return .Signing
func ConvertSigning(res *tsstypes.SigningResult) *Signing {
	if res == nil {
		return nil
	}

	var evmSignature *EVMSignature
	if res.EVMSignature != nil {
		evmSignature = NewEVMSignature(
			res.EVMSignature.RAddress,
			res.EVMSignature.Signature,
		)
	}

	return NewSigning(
		uint64(res.Signing.ID),
		res.Signing.Message,
		evmSignature,
	)
}

// NewSigningResult creates a new Signing instance.
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
