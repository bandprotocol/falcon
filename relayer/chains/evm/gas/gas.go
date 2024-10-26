package gas

import (
	"context"
	"fmt"
	"reflect"
)

type GasType int

const (
	GasTypeUndefined GasType = iota
	GasTypeEIP1559
	GasTypeLegacy
)

var gasTypeNameMap = map[GasType]string{
	GasTypeEIP1559: "eip1559",
	GasTypeLegacy:  "legacy",
}

var nameToGasTypeMap map[string]GasType

func init() {
	nameToGasTypeMap = make(map[string]GasType)
	for k, v := range gasTypeNameMap {
		nameToGasTypeMap[v] = k
	}
}

// String returns the string representation of the gas type.
func (g GasType) String() string {
	return gasTypeNameMap[g]
}

// DecodeGasTypeHook decode hook to convert strings to GasType
func DecodeGasTypeHook(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
	if from.Kind() != reflect.String || to != reflect.TypeOf(GasType(0)) {
		return data, nil
	}

	gasTypeStr, ok := data.(string)
	if !ok {
		return data, fmt.Errorf("expected string, got %T", data)
	}

	return ToGasType(gasTypeStr), nil
}

// MarshalText is used for toml encoding.
func (g GasType) MarshalText() ([]byte, error) {
	return []byte(g.String()), nil
}

// ToGasType converts a string to a GasType.
func ToGasType(s string) GasType {
	if t, ok := nameToGasTypeMap[s]; ok {
		return t
	}

	return GasTypeUndefined
}

type Param struct {
	GasPrice       uint64
	MaxBaseFee     uint64
	MaxPriorityFee uint64
}

type Gas interface {
	Param() Param
	Bump(float64) Gas
}

type GasModel interface {
	GetGas(ctx context.Context) Gas
	GasType() GasType
}
