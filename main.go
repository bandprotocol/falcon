package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/bandprotocol/falcon/cmd"
)

func main() {
	// loading .env file
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		panic(fmt.Sprintf("Error due to loading .env file; %v", err))
	}

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
		os.Exit(1)
	}
}
