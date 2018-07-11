package main

import (
	"github.com/loomnetwork/go-loom/plugin"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/loomnetwork/zombie_battleground/types/zb"

	"github.com/loomnetwork/zombie_battleground/battleground"
)

type ZombieBattleground struct {
}

func (z *ZombieBattleground) Meta() (plugin.Meta, error) {
	return plugin.Meta{
		Name:    "ZombieBattleground",
		Version: "1.0.0",
	}, nil
}

func (z *ZombieBattleground) Init(ctx contract.Context, req *plugin.Request) error {
	return nil
}

func (z *ZombieBattleground) GetAccount(ctx contract.StaticContext, req *zb.GetAccountRequest) (*zb.Account, error) {
	return battleground.GetAccount(ctx, req)
}

func UpdateAccount(ctx contract.Context, req *zb.UpsertAccountRequest) (*zb.Account, error) {
	return battleground.UpdateAccount(ctx, req)
}

func CreateAccount(ctx contract.Context, req *zb.UpsertAccountRequest) error {
	return battleground.CreateAccount(ctx, req)
}

var Contract plugin.Contract = contract.MakePluginContract(&ZombieBattleground{})

func main() {
	plugin.Serve(Contract)
}
