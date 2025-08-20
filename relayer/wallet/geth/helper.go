package geth

import (
	"fmt"
	"path"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// HexToETHAddress checks a given string and converts it to an geth address. The string must
// be align with the ^(0x)?[0-9a-fA-F]{40}$ regex format, e.g. 0xe688b84b23f322a994A53dbF8E15FA82CDB71127.
func HexToETHAddress(s string) (common.Address, error) {
	if !common.IsHexAddress(s) {
		return common.Address{}, fmt.Errorf("invalid address: %s", s)
	}

	return common.HexToAddress(s), nil
}

// ExtractKeyName returns the filename (the key name) without its extension, or an error if empty.
func ExtractKeyName(filePath string) (string, error) {
	fileName := path.Base(filePath)

	keyName := strings.TrimSuffix(fileName, path.Ext(fileName))
	if keyName == "" {
		return "", fmt.Errorf("wrong keyname format")
	}

	return keyName, nil
}

// getEVMKeyStoreDir returns the key store directory
func getEVMKeyStoreDir(homePath, chainName string) []string {
	return []string{homePath, "keys", chainName, "priv"}
}

// getEVMSignerRecordDir returns the signer record directory
func getEVMSignerRecordDir(homePath, chainName string) []string {
	return []string{homePath, "keys", chainName, "signer"}
}

// getEVMRemoteSignerRecordDir returns the remote signer directory
func getEVMRemoteSignerRecordDir(homePath, chainName string) []string {
	return []string{homePath, "keys", chainName, "remote"}
}
