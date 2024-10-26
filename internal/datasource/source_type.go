package datasource

import (
	"fmt"
	"reflect"
)

type SourceType int

const (
	SourceTypeUndefined SourceType = iota
	SourceTypeFix
	SourceTypeWeb3Legacy
	SourceTypeWeb3EIP1559
)

var SourceTypeNameMap = map[SourceType]string{
	SourceTypeFix:         "fix",
	SourceTypeWeb3Legacy:  "web3_legacy",
	SourceTypeWeb3EIP1559: "web3_eip1559",
}

var nameToSourceTypeMap map[string]SourceType

func init() {
	nameToSourceTypeMap = make(map[string]SourceType)
	for k, v := range SourceTypeNameMap {
		nameToSourceTypeMap[v] = k
	}
}

// String returns the string representation of the SourceType.
func (c SourceType) String() string {
	return SourceTypeNameMap[c]
}

// DecodeSourceTypeHook decode hook to convert strings to SourceType
func DecodeSourceTypeHook(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
	if from.Kind() != reflect.String || to != reflect.TypeOf(SourceType(0)) {
		return data, nil
	}

	SourceTypeStr, ok := data.(string)
	if !ok {
		return data, fmt.Errorf("expected string, got %T", data)
	}

	return ToSourceType(SourceTypeStr), nil
}

// MarshalText is used for toml encoding.
func (c SourceType) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

// ToSourceType converts a string to a SourceType.
func ToSourceType(s string) SourceType {
	if t, ok := nameToSourceTypeMap[s]; ok {
		return t
	}

	return SourceTypeUndefined
}
