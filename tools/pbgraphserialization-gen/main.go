package main

import (
	"github.com/loomnetwork/gamechain/tools/pbgraphserialization-gen/generator"
	"os"
)

func main() {
	if err := generator.Execute(); err != nil {
		os.Exit(1)
	}
}