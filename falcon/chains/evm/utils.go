package evm

import (
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// HexToAddress checks a given string and converts it to an geth address.
func HexToAddress(s string) (gethcommon.Address, error) {
	if !gethcommon.IsHexAddress(s) {
		return gethcommon.Address{}, fmt.Errorf("invalid address: %s", s)
	}

	return gethcommon.HexToAddress(s), nil
}
