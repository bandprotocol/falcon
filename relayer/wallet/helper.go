package wallet

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// HexToAddress checks a given string and converts it to an geth address. The string must
// be align with the ^(0x)?[0-9a-fA-F]{40}$ regex format, e.g. 0xe688b84b23f322a994A53dbF8E15FA82CDB71127.
func HexToAddress(s string) (common.Address, error) {
	if !common.IsHexAddress(s) {
		return common.Address{}, fmt.Errorf("invalid address: %s", s)
	}

	return common.HexToAddress(s), nil
}
