package evm

import (
	"fmt"
	"math/big"
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

// GasInfo contains the gas type and gas information being used for submitting a transaction.
type GasInfo struct {
	Type           GasType
	GasPrice       *big.Int
	GasPriorityFee *big.Int
	GasBaseFee     *big.Int
}

// NewGasLegacyInfo creates a new GasInfo instance with gas type legacy.
func NewGasLegacyInfo(gasPrice *big.Int) GasInfo {
	return GasInfo{
		Type:     GasTypeLegacy,
		GasPrice: gasPrice,
	}
}

// NewGasEIP1559Info creates a new GasInfo instance with gas type EIP1559.
func NewGasEIP1559Info(gasPriorityFee, gasBaseFee *big.Int) GasInfo {
	return GasInfo{
		Type:           GasTypeEIP1559,
		GasPriorityFee: gasPriorityFee,
		GasBaseFee:     gasBaseFee,
	}
}
