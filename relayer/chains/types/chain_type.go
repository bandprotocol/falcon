package types

import (
	"database/sql/driver"
	"encoding"
	"fmt"
	"reflect"
)

var _ encoding.TextUnmarshaler = (*ChainType)(nil)

// ChainType represents the type of chain.
type ChainType int

const (
	ChainTypeUndefined ChainType = iota
	ChainTypeEVM
	ChainTypeXRPL
)

var chainTypeNameMap = map[ChainType]string{
	ChainTypeEVM:  "evm",
	ChainTypeXRPL: "xrpl",
}

var nameToChainTypeMap map[string]ChainType

func init() {
	nameToChainTypeMap = make(map[string]ChainType)
	for k, v := range chainTypeNameMap {
		nameToChainTypeMap[v] = k
	}
}

// String returns the string representation of the ChainType.
func (c ChainType) String() string {
	return chainTypeNameMap[c]
}

// UnmarshalText is used for toml decoding.
func (c *ChainType) UnmarshalText(text []byte) error {
	v, ok := nameToChainTypeMap[string(text)]
	if !ok {
		return fmt.Errorf("invalid chain type %s %v", text, v)
	}

	*c = v
	return nil
}

// DecodeChainTypeHook decode hook to convert strings to ChainType
func DecodeChainTypeHook(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
	if from.Kind() != reflect.String || to != reflect.TypeOf(ChainType(0)) {
		return data, nil
	}

	chainTypeStr, ok := data.(string)
	if !ok {
		return data, fmt.Errorf("expected string, got %T", data)
	}

	return ToChainType(chainTypeStr), nil
}

// MarshalText is used for toml encoding.
func (c ChainType) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

// Scan scans string value into ChainType, implements sql.Scanner interface.
// (needs to manually creates `chain_type` type in a database first
// by "CREATE TYPE chain_type AS ENUM ('evm', 'xrpl')")
func (c *ChainType) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("ChainType.Scan: expected string, got %T", value)
	}

	*c = ToChainType(str)
	return nil
}

// Value converts ChainType to a driver.Value (string form).
func (c ChainType) Value() (driver.Value, error) { return c.String(), nil }

// ToChainType converts a string to a ChainType.
func ToChainType(s string) ChainType {
	if t, ok := nameToChainTypeMap[s]; ok {
		return t
	}

	return ChainTypeUndefined
}
