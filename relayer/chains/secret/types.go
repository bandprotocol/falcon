package secret

import (
	"math/big"

	"github.com/shopspring/decimal"

	"github.com/bandprotocol/falcon/relayer/chains/types"
)

// TxResult is the result of transaction confirmation.
type TxResult struct {
	Status            types.TxStatus
	GasUsed           decimal.NullDecimal
	EffectiveGasPrice decimal.NullDecimal
	BlockHeight       *big.Int
	FailureReason     string
}

func NewTxResult(
	status types.TxStatus,
	gasUsed decimal.NullDecimal,
	effectiveGasPrice decimal.NullDecimal,
	blockHeight *big.Int,
	failureReason string,
) TxResult {
	return TxResult{
		Status:            status,
		GasUsed:           gasUsed,
		EffectiveGasPrice: effectiveGasPrice,
		BlockHeight:       blockHeight,
		FailureReason:     failureReason,
	}
}
