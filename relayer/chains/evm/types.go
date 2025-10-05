package evm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"

	"github.com/bandprotocol/falcon/relayer/chains/types"
)

// ClientConnectionResult is the struct that contains the result of connecting to the specific endpoint.
type ClientConnectionResult struct {
	Endpoint    string
	Client      *ethclient.Client
	BlockHeight uint64
}

// TxResult is the result of transaction.
type TxResult struct {
	TxHash            string
	Status            types.TxStatus
	GasUsed           decimal.NullDecimal
	EffectiveGasPrice decimal.NullDecimal
	BlockNumber       *big.Int

	// empty when Status == SUCCESS
	FailureReason string
}

// NewTxResult creates a new TxResult instance.
func NewTxResult(
	txHash string,
	status types.TxStatus,
	gasUsed decimal.NullDecimal,
	effectiveGasPrice decimal.NullDecimal,
	blockNumber *big.Int,
	failureReason string,
) TxResult {
	return TxResult{
		TxHash:            txHash,
		Status:            status,
		GasUsed:           gasUsed,
		EffectiveGasPrice: effectiveGasPrice,
		BlockNumber:       blockNumber,
		FailureReason:     failureReason,
	}
}
