package xrpl

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

// xrplAssetHexLen is the fixed byte length of an XRPL currency hex code (20 bytes = 40 hex chars).
const xrplAssetHexLen = 40

// ParseAssetsFromSignal parses a signal ID into base, quote assets and then convert them to hex string if length != 3
func ParseAssetsFromSignal(signalID string) (string, string, error) {
	parts := strings.SplitN(signalID, ":", 2)
	if len(parts) < 2 || parts[1] == "" {
		return "", "", fmt.Errorf("invalid signal format (expected '<prefix>:<pair>'): %s", signalID)
	}
	pair := parts[1]
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
		return StringToHex(asset, xrplAssetHexLen)
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

// StringToHex converts a string to a hex string of a specified length.
func StringToHex(str string, length int) (string, error) {
	encoded := strings.ToUpper(hex.EncodeToString([]byte(str)))
	if length > 0 {
		if len(encoded) > length {
			return "", fmt.Errorf("hex string length %d exceeds expected length %d", len(encoded), length)
		}
		encoded += strings.Repeat("0", length-len(encoded))
	}
	return encoded, nil
}

// Uint64StrToHexStr converts a string to a hex string of exactly 16 characters.
func Uint64StrToHexStr(uint64Str string) (string, error) {
	n, err := strconv.ParseUint(uint64Str, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid or out-of-range uint64 string: %s", uint64Str)
	}
	return fmt.Sprintf("%016X", n), nil
}
