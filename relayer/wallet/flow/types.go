package flow

// SignerPayload is the payload passed (JSON-serialized) to the Flow signer.
type SignerPayload struct {
	Address         string
	ComputeLimit    uint64
	BlockID         string
	KeyIndex        uint32
	Sequence        uint64
	ContractAddress string
}

// NewSignerPayload creates a new SignerPayload.
func NewSignerPayload(address string, computeLimit uint64, blockID string, keyIndex uint32, sequence uint64, contractAddress string) SignerPayload {
	return SignerPayload{
		Address:         address,
		ComputeLimit:    computeLimit,
		BlockID:         blockID,
		KeyIndex:        keyIndex,
		Sequence:        sequence,
		ContractAddress: contractAddress,
	}
}
