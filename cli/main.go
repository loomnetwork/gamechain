package main

import (
	"os"

	"github.com/loomnetwork/zombie_battleground/cli/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
