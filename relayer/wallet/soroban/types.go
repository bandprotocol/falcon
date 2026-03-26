package soroban

type SignerPayload struct {
	SourceAccount     string
	ContractAddress   string
	Fee               string
	Sequence          uint64
	NetworkPassphrase string
}

func NewSignerPayload(
	sourceAccount string,
	contractAddress string,
	fee string,
	sequence uint64,
	networkPassphrase string,
) SignerPayload {
	return SignerPayload{
		SourceAccount:     sourceAccount,
		ContractAddress:   contractAddress,
		Fee:               fee,
		Sequence:          sequence,
		NetworkPassphrase: networkPassphrase,
	}
}
