package main

import (
	"github.com/loomnetwork/go-loom/plugin"
	"github.com/loomnetwork/zombie_battleground/card"
)

var Contract = card.Contract

func main() {
	plugin.Serve(Contract)
}
