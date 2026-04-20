package evm

// SignerPayload carries the EVM transaction parameters needed to build, sign,
// and return the complete EIP-2718 encoded transaction.
// Data holds the calldata; it is used by LocalSigner and ignored by the remote KMS
// (which reconstructs calldata from the TSS fields instead).
type SignerPayload struct {
	Address  string `json:"address"`
	ChainID  uint64 `json:"chain_id"`
	Nonce    uint64 `json:"nonce"`
	To       string `json:"to"`
	GasLimit uint64 `json:"gas_limit"`
	Data     []byte `json:"data,omitempty"`
	// Legacy tx field
	GasPrice []byte `json:"gas_price,omitempty"`
	// EIP-1559 fields
	GasFeeCap []byte `json:"gas_fee_cap,omitempty"`
	GasTipCap []byte `json:"gas_tip_cap,omitempty"`
}

// NewSignerPayload constructs a SignerPayload with the given transaction parameters.
func NewSignerPayload(
	address string,
	chainID uint64,
	nonce uint64,
	to string,
	gasLimit uint64,
	gasPrice []byte,
	gasFeeCap []byte,
	gasTipCap []byte,
) SignerPayload {
	return SignerPayload{
		Address:   address,
		ChainID:   chainID,
		Nonce:     nonce,
		To:        to,
		GasLimit:  gasLimit,
		GasPrice:  gasPrice,
		GasFeeCap: gasFeeCap,
		GasTipCap: gasTipCap,
	}
}
