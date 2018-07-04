package battleground

import (
	"context"

	"github.com/loomnetwork/go-loom/plugin"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/loomnetwork/zombie_battleground/types/zb"
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

func (z *ZombieBattleground) CreateAccount(ctx context.Context, req *zb.CreateAccountRequest) error {
	return nil
}

// TODO add more methods to support functionality
// - CreateAccount
// - ...

var Contract plugin.Contract = contract.MakePluginContract(&ZombieBattleground{})
