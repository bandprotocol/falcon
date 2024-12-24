package cmd

import "fmt"

func ErrConfigNotExist(homePath string) error {
	return fmt.Errorf("config does not exist: %s", homePath)
}
