package types

import (
	cmbytes "github.com/cometbft/cometbft/libs/bytes"

	tsstypes "github.com/bandprotocol/chain/v3/x/tss/types"
)

// EVMSignature defines a signature in the EVM format.
type EVMSignature struct {
	RAddress  cmbytes.HexBytes `json:"r_address"`
	Signature cmbytes.HexBytes `json:"signature"`
}

// NewEVMSignature creates a new EVMSignature instance.
func NewEVMSignature(
	rAddress cmbytes.HexBytes,
	signature cmbytes.HexBytes,
) *EVMSignature {
	return &EVMSignature{
		RAddress:  rAddress,
		Signature: signature,
	}
}

// Signing contains information of a requested message and group signature.
type Signing struct {
	ID           uint64           `json:"id"`
	Message      cmbytes.HexBytes `json:"message"`
	EVMSignature *EVMSignature    `json:"evm_signature"`
	Status       string           `json:"signing_status"`
}

// ConvertSigning converts tsstypes.SigningResult and return Signing type.
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
		tsstypes.SigningStatus_name[int32(res.Signing.Status)],
	)
}

// NewSigningResult creates a new Signing instance.
func NewSigning(
	id uint64,
	message cmbytes.HexBytes,
	evmSignature *EVMSignature,
	status string,
) *Signing {
	return &Signing{
		ID:           id,
		Message:      message,
		EVMSignature: evmSignature,
		Status:       status,
	}
}
