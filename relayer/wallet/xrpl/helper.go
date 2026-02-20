package xrpl

import (
	"fmt"
	"path"
	"strings"
)

// ExtractKeyName returns the filename (the key name) without its extension, or an error if empty.
func ExtractKeyName(filePath string) (string, error) {
	fileName := path.Base(filePath)

	keyName := strings.TrimSuffix(fileName, path.Ext(fileName))
	if keyName == "" {
		return "", fmt.Errorf("wrong keyname format")
	}

	return keyName, nil
}

// getXRPLKeyDir returns the key record directory.
func getXRPLKeyDir(homePath, chainName string) []string {
	return []string{homePath, "keys", chainName, "metadata"}
}
