package relayer

import "fmt"

func ErrConfigNotExist(homePath string) error {
	return fmt.Errorf("config does not exist: %s", homePath)
}

func ErrConfigExist(cfgPath string) error {
	return fmt.Errorf("config already exists: %s", cfgPath)
}

func ErrChainNameExist(chainName string) error {
	return fmt.Errorf("chain name already exists: %s", chainName)
}

func ErrChainNameNotExist(chainName string) error {
	return fmt.Errorf("chain name does not exist: %s", chainName)
}

func ErrUnsupportedChainType(typeName string) error {
	return fmt.Errorf("unsupported chain type: %s", typeName)
}
