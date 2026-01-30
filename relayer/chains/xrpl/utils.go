package xrpl

import (
	"encoding/hex"
	"fmt"
	"strings"
)

func stringToHex(str string, length int) (string, error) {
	encoded := strings.ToUpper(hex.EncodeToString([]byte(str)))
	if length != 0 && len(encoded) > length {
		return "", fmt.Errorf("hex string length %d exceeds expected length %d", len(encoded), length)
	}
	for length != 0 && len(encoded) < length {
		encoded += "0"
	}
	return encoded, nil
}

func parseAssetsFromSignal(signalID string) (string, string, error) {
	parts := strings.Split(signalID, ":")
	core := parts[len(parts)-1]
	assets := strings.Split(core, "-")
	if len(assets) != 2 {
		return "", "", fmt.Errorf("invalid signal_id format: %s", signalID)
	}
	base := strings.TrimSpace(assets[0])
	quote := strings.TrimSpace(assets[1])
	if base == "" || quote == "" {
		return "", "", fmt.Errorf("invalid signal_id format: %s", signalID)
	}

	baseAsset := base
	if len(base) != 3 {
		var err error
		baseAsset, err = stringToHex(base, 40)
		if err != nil {
			return "", "", err
		}
	}

	quoteAsset := quote
	if len(quote) != 3 {
		var err error
		quoteAsset, err = stringToHex(quote, 40)
		if err != nil {
			return "", "", err
		}
	}

	return baseAsset, quoteAsset, nil
}
