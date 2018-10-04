package main

import (
	"github.com/loomnetwork/go-loom/plugin"
	"github.com/loomnetwork/gamechain/battleground"
)

var Contract = battleground.Contract

func main() {
	plugin.Serve(Contract)
}
