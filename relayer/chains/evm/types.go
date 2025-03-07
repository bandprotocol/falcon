package evm

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
)

// ClientConnectionResult is the struct that contains the result of connecting to the specific endpoint.
type ClientConnectionResult struct {
	Endpoint    string
	Client      *ethclient.Client
	BlockHeight uint64
}

// TxStatus is the status of the transaction
type TxStatus int

const (
	TX_STATUS_UNDEFINED TxStatus = iota
	TX_STATUS_UNMINED
	TX_STATUS_SUCCESS
	TX_STATUS_FAILED
)

var txStatusNameMap = map[TxStatus]string{
	TX_STATUS_UNDEFINED: "undefined",
	TX_STATUS_UNMINED:   "unmined",
	TX_STATUS_SUCCESS:   "success",
	TX_STATUS_FAILED:    "failed",
}

func (t TxStatus) String() string {
	return txStatusNameMap[t]
}

// ConfirmTxResult is the result of confirming a transaction
type ConfirmTxResult struct {
	TxHash  string
	Status  TxStatus
	GasUsed decimal.NullDecimal
}

func NewConfirmTxResult(
	txHash string,
	status TxStatus,
	gasUsed decimal.NullDecimal,
) *ConfirmTxResult {
	return &ConfirmTxResult{
		TxHash:  txHash,
		Status:  status,
		GasUsed: gasUsed,
	}
}
