package evm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
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

//go:generate go run github.com/fjl/gencodec -type TxReceipt -field-override txResultMarshaling -out gen_tx_receipt_json.go

// TxReceipt is struct that represents the transaction receipt that will be marshaled to/from JSON.
type TxReceipt struct {
	Status            uint64   `json:"status"`
	GasUsed           uint64   `json:"gasUsed"`
	EffectiveGasPrice *big.Int `json:"effectiveGasPrice"`
	BlockNumber       *big.Int `json:"blockNumber,omitempty"`
}

// txResultMarshaling is an internal struct for JSON marshaling of TxReceipt.
type txResultMarshaling struct {
	Status            hexutil.Uint64
	GasUsed           hexutil.Uint64
	EffectiveGasPrice *hexutil.Big
	BlockNumber       *hexutil.Big
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
