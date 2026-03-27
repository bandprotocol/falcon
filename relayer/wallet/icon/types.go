package icon

type SignerPayload struct {
	Relayer         string
	ContractAddress string
	StepLimit       uint64
	NetworkID       string
}

func NewSignerPayload(relayer string, contractAddress string, stepLimit uint64, networkID string) SignerPayload {
	return SignerPayload{
		Relayer:         relayer,
		ContractAddress: contractAddress,
		StepLimit:       stepLimit,
		NetworkID:       networkID,
	}
}
