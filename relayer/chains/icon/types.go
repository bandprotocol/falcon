package icon

import (
	"math/big"

	iconclient "github.com/icon-project/goloop/client"
	"github.com/shopspring/decimal"

	"github.com/bandprotocol/falcon/relayer/chains/types"
)

// ClientConnectionResult is the struct that contains the result of connecting to the specific endpoint.
type ClientConnectionResult struct {
	Endpoint    string
	Client      *iconclient.ClientV3
	BlockHeight uint64
}

// TxResult is the result of transaction.
type TxResult struct {
	Status            types.TxStatus
	GasUsed           decimal.NullDecimal
	EffectiveGasPrice decimal.NullDecimal
	BlockHeight       *big.Int

	// empty when Status == SUCCESS
	FailureReason string
}

// NewTxResult creates a new TxResult instance.
func NewTxResult(
	status types.TxStatus,
	stepUsed decimal.NullDecimal,
	stepPrice decimal.NullDecimal,
	blockHeight *big.Int,
	failureReason string,
) TxResult {
	return TxResult{
		Status:            status,
		GasUsed:           stepUsed,
		EffectiveGasPrice: stepPrice,
		BlockHeight:       blockHeight,
		FailureReason:     failureReason,
	}
}
