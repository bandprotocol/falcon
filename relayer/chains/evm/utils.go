package evm

import (
	"fmt"
	"math/big"
	"strings"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

const privateKeyPrefix = "0x"

// HexToAddress checks a given string and converts it to an geth address. The string must
// be align with the ^(0x)?[0-9a-fA-F]{40}$ regex format, e.g. 0xe688b84b23f322a994A53dbF8E15FA82CDB71127.
func HexToAddress(s string) (gethcommon.Address, error) {
	if !gethcommon.IsHexAddress(s) {
		return gethcommon.Address{}, fmt.Errorf("invalid address: %s", s)
	}

	return gethcommon.HexToAddress(s), nil
}

// StripPrivateKeyPrefix removes the "0x" prefix from the given private key string, if present.
func StripPrivateKeyPrefix(privateKey string) string {
	return strings.TrimPrefix(privateKey, privateKeyPrefix)
}

// MultiplyWithFloat64 multiplies a big.Int value with a float64 multiplier and convert back to big.Int.
func MultiplyBigIntWithFloat64(value *big.Int, multiplier float64) *big.Int {
	// Define precision scale
	scale := 1_000_000

	multiplierScaled := int64(multiplier * float64(scale))
	valueScaled := new(big.Int).Mul(value, big.NewInt(multiplierScaled))
	result := new(big.Int).Div(valueScaled, big.NewInt(int64(scale)))

	return result
}
