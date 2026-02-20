package xrpl

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
)

// StringToHex converts a string to a hex string of a specified length.
func StringToHex(str string, length int) (string, error) {
	encoded := strings.ToUpper(hex.EncodeToString([]byte(str)))
	if length != 0 && len(encoded) > length {
		return "", fmt.Errorf("hex string length %d exceeds expected length %d", len(encoded), length)
	}
	for length != 0 && len(encoded) < length {
		encoded += "0"
	}
	return encoded, nil
}

// ParseAssetsFromSignal parses a signal ID into base, quote assets and then convert them to hex string if length != 3
func ParseAssetsFromSignal(signalID string) (string, string, error) {
	var pair string
	parts := strings.SplitN(signalID, ":", 2)
	if len(parts) < 2 || parts[1] == "" {
		return "", "", fmt.Errorf("invalid signal format (expected '<prefix>:<pair>'): %s", signalID)
	} else {
		pair = parts[1]
	}
	assets := strings.SplitN(pair, "-", 2)
	if len(assets) != 2 {
		return "", "", fmt.Errorf("invalid base/quote format (expected '<base>-<quote>'): %s", signalID)
	}

	processAsset := func(asset string) (string, error) {
		asset = strings.TrimSpace(asset)
		if asset == "" {
			return "", fmt.Errorf("asset cannot be empty")
		}
		if len(asset) == 3 {
			return asset, nil
		}
		return StringToHex(asset, 40)
	}

	base, err := processAsset(assets[0])
	if err != nil {
		return "", "", err
	}

	quote, err := processAsset(assets[1])
	if err != nil {
		return "", "", err
	}

	return base, quote, nil
}

// Uint64StrToHexStr converts a string to a hex string of exactly 16 characters.
func Uint64StrToHexStr(uint64Str string) (string, error) {
	n := new(big.Int)
	n, ok := n.SetString(uint64Str, 10)
	if !ok {
		return "", fmt.Errorf("invalid numeric string: %s", uint64Str)
	}

	// Check for negative values
	if n.Sign() < 0 {
		return "", fmt.Errorf("value must be non-negative: %s", uint64Str)
	}

	// Check if the value fits within 64 bits
	if n.BitLen() > 64 {
		return "", fmt.Errorf("value exceeds uint64 limit: %s", uint64Str)
	}

	// Convert to native uint64 and format to exactly 16 hex chars
	return fmt.Sprintf("%016X", n.Uint64()), nil
}
