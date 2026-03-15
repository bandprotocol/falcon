package xrpl

type SignerPayload struct {
	Account  string
	OracleID uint64
	Fee      string
	Sequence uint64
}

func NewSignerPayload(account string, oracleID uint64, fee string, sequence uint64) SignerPayload {
	return SignerPayload{
		Account:  account,
		OracleID: oracleID,
		Fee:      fee,
		Sequence: sequence,
	}
}
