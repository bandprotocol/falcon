package datasource

import (
	"fmt"
	"math"
	"math/big"
	"sort"
	"strconv"

	"github.com/shopspring/decimal"
)

var maxUint64 = new(big.Int).SetUint64(math.MaxUint64)

// median returns the median of a slice of uint64s.
func median(values []uint64) (uint64, error) {
	// Sort the slice.
	sort.Slice(values, func(i, j int) bool { return values[i] < values[j] })

	// Find the median.
	n := len(values)
	if n == 0 {
		return 0, fmt.Errorf("input must not be empty")
	} else if n%2 == 0 {
		return (values[n/2-1] + values[n/2]) / 2, nil
	} else {
		return values[n/2], nil
	}
}

// parseFloat converts the input data to a float64.
func parseFloat(data interface{}) (float64, error) {
	switch d := data.(type) {
	case float64:
		return d, nil
	case string:
		return strconv.ParseFloat(d, 64)
	default:
		return 0, fmt.Errorf("conversion to float64 from %T not supported", d)
	}
}

// multiplyDecimals multiplies two float64 numbers and returns the result as a uint64.
func multiplyDecimals(a float64, b float64) (uint64, error) {
	dGasBase := decimal.NewFromFloat(a)
	dMultiplier := decimal.NewFromFloat(b)

	result := dGasBase.Mul(dMultiplier)

	bInt := new(big.Int)

	_, ok := bInt.SetString(result.String(), 10)
	if !ok {
		return 0, fmt.Errorf("cannot convert to big.Int")
	}

	if bInt.Cmp(maxUint64) > 0 {
		return 0, fmt.Errorf("result is too large")
	}

	return bInt.Uint64(), nil
}
