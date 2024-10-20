package chains

import (
	"encoding"
	"fmt"
)

var _ encoding.TextUnmarshaler = (*ChainType)(nil)

// ChainType represents the type of chain.
type ChainType int

const (
	ChainTypeUndefined ChainType = iota
	ChainTypeEVM
	ChainTypeCosmwasm
)

var chainTypeNameMap = map[ChainType]string{
	ChainTypeEVM:      "evm",
	ChainTypeCosmwasm: "cosmwasm",
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

// MarshalText is used for toml encoding.
func (c ChainType) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

// ToChainType converts a string to a ChainType.
func ToChainType(s string) ChainType {
	if t, ok := nameToChainTypeMap[s]; ok {
		return t
	}

	return ChainTypeUndefined
}
