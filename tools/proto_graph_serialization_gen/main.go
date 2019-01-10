package main

import (
	"github.com/loomnetwork/gamechain/tools/proto_graph_serialization_gen/generator"
	"os"
)

func main() {
	if err := generator.Execute(); err != nil {
		os.Exit(1)
	}
}