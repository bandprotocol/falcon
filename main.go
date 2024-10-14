package main

import (
	"os"

	"github.com/bandprotocol/falcon/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
