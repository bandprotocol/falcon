package wallet

import (
	"fmt"

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

// getEVMKeyStoreDir returns the key store directory
func getEVMKeyStoreDir(homePath, chainName string) []string {
	return []string{homePath, "keys", chainName, "priv"}
}

// getEVMKeyNameInfoPath returns the keyNameInfo file path
func getEVMKeyNameInfoPath(homePath, chainName string) []string {
	return []string{homePath, "keys", chainName, "info", "info.toml"}
}
