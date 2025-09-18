package db

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/bandprotocol/falcon/relayer/chains/types"
)

// Transaction represents transaction information sent to the target chain contract that will be stored in the database.
type Transaction struct {
	ID                uint                `gorm:"primarykey"`
	TxHash            string              `gorm:"unique"`
	TunnelID          uint64              `gorm:"not null"`
	Sequence          uint64              `gorm:"not null"`
	ChainName         string              `gorm:"not null"`
	ChainType         types.ChainType     `gorm:"type:chain_type;not null"`
	Status            types.TxStatus      `gorm:"type:tx_status;not null"`
	GasUsed           decimal.NullDecimal `gorm:"type:decimal"`
	EffectiveGasPrice decimal.NullDecimal `gorm:"type:decimal"`
	BalanceDelta      decimal.NullDecimal `gorm:"type:decimal"`

	SignalPrices   []SignalPrice
	BlockTimestamp time.Time `gorm:"default:NULL"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NewUnconfirmedTransaction creates a new pending Transaction instance.
func NewUnconfirmedTransaction(
	txHash string,
	tunnelID uint64,
	sequence uint64,
	chainName string,
	chainType types.ChainType,
	status types.TxStatus,
	signalPrices []SignalPrice,
) *Transaction {
	return &Transaction{
		TxHash:       txHash,
		TunnelID:     tunnelID,
		Sequence:     sequence,
		ChainName:    chainName,
		ChainType:    chainType,
		Status:       status,
		SignalPrices: signalPrices,
	}
}

// NewConfirmedTransaction creates a new Transaction instance.
func NewConfirmedTransaction(
	txHash string,
	tunnelID uint64,
	sequence uint64,
	chainName string,
	chainType types.ChainType,
	status types.TxStatus,
	gasUsed decimal.NullDecimal,
	effectiveGasPrice decimal.NullDecimal,
	balanceDelta decimal.NullDecimal,
	signalPrices []SignalPrice,
	blockTimestamp time.Time,
) *Transaction {
	return &Transaction{
		TxHash:            txHash,
		TunnelID:          tunnelID,
		Sequence:          sequence,
		ChainName:         chainName,
		ChainType:         chainType,
		Status:            status,
		GasUsed:           gasUsed,
		EffectiveGasPrice: effectiveGasPrice,
		BalanceDelta:      balanceDelta,
		SignalPrices:      signalPrices,
		BlockTimestamp:    blockTimestamp,
	}
}

// SignalPrice represents the price of a signal for a given transaction.
type SignalPrice struct {
	TransactionID uint   `gorm:"primarykey"`
	SignalID      string `gorm:"primarykey"`
	Price         uint64 `gorm:"not null"`
}

// NewSignalPrice creates a new SignalPrice instance.
func NewSignalPrice(
	signalID string,
	price uint64,
) *SignalPrice {
	return &SignalPrice{
		SignalID: signalID,
		Price:    price,
	}
}
