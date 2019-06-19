package battleground

import (
	"testing"

	"github.com/loomnetwork/go-loom"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
)

func setupInitFromFile(c *ZombieBattleground, pubKeyHex string, addr *loom.Address, ctx *contract.Context, t *testing.T) {
	setup(c, pubKeyHex, addr, ctx, t)
}
