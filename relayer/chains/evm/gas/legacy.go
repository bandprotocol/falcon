package gas

import (
	"context"
	"math"
	"time"

	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/internal/datasource"
)

var (
	_ Gas      = LegacyGas{}
	_ GasModel = LegacyGasModel{}
)

// LegacyGas defines the gas parameters for legacy fee transaction submission.
type LegacyGas struct {
	GasPrice uint64
}

// NewLegacyGas returns a new legacy gas info.
func NewLegacyGas(gasPrice uint64) LegacyGas {
	return LegacyGas{
		GasPrice: gasPrice,
	}
}

// Param returns the gas parameters.
func (g LegacyGas) Param() Param {
	return Param{
		GasPrice:       g.GasPrice,
		MaxBaseFee:     0,
		MaxPriorityFee: 0,
	}
}

// Bump returns a new gas model with the gas price bumped by the given factor.
func (g LegacyGas) Bump(factor float64) Gas {
	return LegacyGas{
		GasPrice: uint64(math.Round(float64(g.GasPrice) * factor)),
	}
}

// LegacyGasModel defines the gas model for legacy fee model.
type LegacyGasModel struct {
	GasPrice        uint64
	DataSources     datasource.DataSources
	QueryGasTimeout time.Duration

	log *zap.Logger
}

// NewLegacyGasModel returns a new legacy gas model.
func NewLegacyGasModel(
	gasPrice uint64,
	sources []datasource.Source,
	queryGasTimeout time.Duration,
	log *zap.Logger,
) LegacyGasModel {
	dataSources := datasource.NewDataSources(sources, log)

	return LegacyGasModel{
		GasPrice:        gasPrice,
		DataSources:     dataSources,
		QueryGasTimeout: queryGasTimeout,
		log:             log,
	}
}

// GasType returns the gas type.
func (m LegacyGasModel) GasType() GasType {
	return GasTypeLegacy
}

// GetGas returns the gas parameters from the predefined sources.
func (m LegacyGasModel) GetGas(ctx context.Context) Gas {
	newCtx, cancel := context.WithTimeout(ctx, m.QueryGasTimeout)
	defer cancel()

	gasPrice, err := m.DataSources.GetData(newCtx)
	if err != nil {
		m.log.Debug("Failed to get gas price", zap.Error(err))
		gasPrice = m.GasPrice
	}

	return NewLegacyGas(gasPrice)
}
