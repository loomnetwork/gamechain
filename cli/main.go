package main

import (
	"os"

	"github.com/loomnetwork/gamechain/cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
