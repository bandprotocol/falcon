package secret

// SignerPayload is the JSON payload sent to the remote signer via WalletSigner.Sign().
// It is then mapped 1:1 onto fkms.v1.SecretSignerPayload.
type SignerPayload struct {
	Sender          string `json:"sender"`
	ContractAddress string `json:"contract_address"`
	ChainID         string `json:"chain_id"`
	AccountNumber   uint64 `json:"account_number"`
	Sequence        uint64 `json:"sequence"`
	GasLimit        uint64 `json:"gas_limit"`
	GasPrices       string `json:"gas_prices"`
	Memo            string `json:"memo"`
	CodeHash        string `json:"code_hash"`
	PubKey          string `json:"pubkey"`
}

func NewSignerPayload(
	sender string,
	contractAddress string,
	chainID string,
	accountNumber uint64,
	sequence uint64,
	gasLimit uint64,
	gasPrices string,
	memo string,
	codeHash string,
	pubKey string,
) SignerPayload {
	return SignerPayload{
		Sender:          sender,
		ContractAddress: contractAddress,
		ChainID:         chainID,
		AccountNumber:   accountNumber,
		Sequence:        sequence,
		GasLimit:        gasLimit,
		GasPrices:       gasPrices,
		Memo:            memo,
		CodeHash:        codeHash,
		PubKey:          pubKey,
	}
}
