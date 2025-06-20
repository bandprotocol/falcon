package recorder

import (
	"github.com/bandprotocol/falcon/relayer/band/types"
)

type TransactionRecorder interface {
	RecordTransaction(tx Transaction) error
}

type Transaction struct {
	TxHash    string `json:"tx_hash"`
	TunnelID  uint64 `json:"tunnel_id"`
	ChainName string `json:"chain_name"`

	Sequence          uint64              `json:"sequence"`
	GasUsed           float64             `json:"gas_used"`
	EffectiveGasPrice float64             `json:"effective_gas_price"`
	SignalPrices      []types.SignalPrice `json:"signal_prices"`
	Timestamp         uint64              `json:"timestamp"`
	Status            string              `json:"status"`
}

func NewTransaction(
	txHash string,
	tunnelID uint64,
	chainName string,
	sequence uint64,
	gasUsed float64,
	effectiveGasPrice float64,
	signalPrices []types.SignalPrice,
	timestamp uint64,
	status string,
) Transaction {
	return Transaction{
		TxHash:            txHash,
		TunnelID:          tunnelID,
		ChainName:         chainName,
		Sequence:          sequence,
		GasUsed:           gasUsed,
		EffectiveGasPrice: effectiveGasPrice,
		SignalPrices:      signalPrices,
		Timestamp:         timestamp,
		Status:            status,
	}
}
