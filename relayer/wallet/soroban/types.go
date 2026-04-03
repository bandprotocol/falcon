package soroban

type SignerPayload struct {
	SourceAccount     string
	ContractAddress   string
	Fee               string
	Sequence          int64
	NetworkPassphrase string
	RpcUrls           []string
}

func NewSignerPayload(
	sourceAccount string,
	contractAddress string,
	fee string,
	sequence int64,
	networkPassphrase string,
	rpcUrls []string,
) SignerPayload {
	return SignerPayload{
		SourceAccount:     sourceAccount,
		ContractAddress:   contractAddress,
		Fee:               fee,
		Sequence:          sequence,
		NetworkPassphrase: networkPassphrase,
		RpcUrls:           rpcUrls,
	}
}
