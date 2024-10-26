package gas

import (
	"context"
	"math"
	"time"

	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/internal/datasource"
)

var (
	_ Gas      = EIP1559Gas{}
	_ GasModel = EIP1559GasModel{}
)

// EIP1559Gas defines the gas parameters for EIP-1559 fee transaction submission.
type EIP1559Gas struct {
	MaxBaseFee     uint64
	MaxPriorityFee uint64
}

// NewEIP1559Gas returns a new EIP-1559 gas info.
func NewEIP1559Gas(maxBaseFee, maxPriorityFee uint64) EIP1559Gas {
	return EIP1559Gas{
		MaxBaseFee:     maxBaseFee,
		MaxPriorityFee: maxPriorityFee,
	}
}

// Param returns the gas parameters.
func (g EIP1559Gas) Param() Param {
	return Param{
		GasPrice:       0,
		MaxBaseFee:     g.MaxBaseFee,
		MaxPriorityFee: g.MaxPriorityFee,
	}
}

// Bump returns a new gas model with the gas price bumped by the given factor.
func (g EIP1559Gas) Bump(factor float64) Gas {
	return EIP1559Gas{
		MaxBaseFee:     g.MaxBaseFee,
		MaxPriorityFee: uint64(math.Round(float64(g.MaxPriorityFee) * factor)),
	}
}

// EIP1559GasModel defines the gas model for EIP-1559 fee model.
type EIP1559GasModel struct {
	MaxBaseFee      uint64
	MaxPriorityFee  uint64
	DataSources     datasource.DataSources
	QueryGasTimeout time.Duration

	log *zap.Logger
}

// NewEIP1559GasModel returns a new EIP-1559 gas model.
func NewEIP1559GasModel(
	maxBaseFee uint64,
	maxPriorityFee uint64,
	sources []datasource.Source,
	queryGasTimeout time.Duration,
	log *zap.Logger,
) EIP1559GasModel {
	dataSources := datasource.NewDataSources(sources, log)

	return EIP1559GasModel{
		MaxBaseFee:      maxBaseFee,
		MaxPriorityFee:  maxPriorityFee,
		DataSources:     dataSources,
		QueryGasTimeout: queryGasTimeout,
		log:             log,
	}
}

// GasType returns the gas type.
func (m EIP1559GasModel) GasType() GasType {
	return GasTypeEIP1559
}

// GetGas returns the gas parameters from the predefined sources.
func (m EIP1559GasModel) GetGas(ctx context.Context) Gas {
	newCtx, cancel := context.WithTimeout(ctx, m.QueryGasTimeout)
	defer cancel()

	priorityFee, err := m.DataSources.GetData(newCtx)
	if err != nil {
		m.log.Debug("Failed to get priority fee", zap.Error(err))
		priorityFee = m.MaxPriorityFee
	}

	return NewEIP1559Gas(m.MaxBaseFee, priorityFee)
}
