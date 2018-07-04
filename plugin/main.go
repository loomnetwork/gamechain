package main

import (
	"github.com/loomnetwork/go-loom/plugin"
	"github.com/loomnetwork/zombie_battleground/battleground"
)

var Contract = battleground.Contract

func main() {
	plugin.Serve(Contract)
}
