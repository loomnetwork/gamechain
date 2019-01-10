package main

import (
	"github.com/loomnetwork/gamechain/sergen/generator"
	"os"
)

func main() {
	if err := generator.Execute(); err != nil {
		os.Exit(1)
	}
}