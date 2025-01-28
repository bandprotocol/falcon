package relayer

import "fmt"

var (
	ErrConfigNotExist = func(homePath string) error {
		return fmt.Errorf("config does not exist: %s", homePath)
	}

	ErrConfigExist = func(cfgPath string) error {
		return fmt.Errorf("config already exists: %s", cfgPath)
	}

	ErrChainNameExist = func(chainName string) error {
		return fmt.Errorf("chain name already exists: %s", chainName)
	}

	ErrChainNameNotExist = func(chainName string) error {
		return fmt.Errorf("chain name does not exist: %s", chainName)
	}

	ErrUnsupportedChainType = func(typeName string) error {
		return fmt.Errorf("unsupported chain type: %s", typeName)
	}
)
