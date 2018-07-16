package card

import (
	"github.com/loomnetwork/go-loom/plugin"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/loomnetwork/zombie_battleground/types/zbcard"
)

type Card struct {
}

func (c *Card) Meta() (plugin.Meta, error) {
	return plugin.Meta{
		Name:    "ZombieBattlegroundCard",
		Version: "1.0.0",
	}, nil
}

func (c *Card) Init(ctx contract.Context, req *zbcard.InitRequest) error {
	return saveCardList(ctx, req.Cardlist)
}

var Contract plugin.Contract = contract.MakePluginContract(&Card{})
