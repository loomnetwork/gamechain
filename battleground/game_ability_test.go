package battleground

import (
	"encoding/hex"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"os"
	"testing"

	"github.com/gogo/protobuf/jsonpb"
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/plugin"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/stretchr/testify/assert"
)

func setupInitFromFile(c *ZombieBattleground, pubKeyHex string, addr *loom.Address, ctx *contract.Context, t *testing.T) {
	c = &ZombieBattleground{}
	pubKey, _ := hex.DecodeString(pubKeyHex)

	addr = &loom.Address{
		Local: loom.LocalAddressFromPublicKey(pubKey),
	}

	*ctx = contract.WrapPluginContext(
		plugin.CreateFakeContext(*addr, *addr),
	)

	// read from update-init-test file
	f, err := os.Open("./update-init-test.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var updateInitData zb_calls.InitRequest

	if err := new(jsonpb.Unmarshaler).Unmarshal(f, &updateInitData); err != nil {
		panic(err)
	}

	err = c.Init(*ctx, &updateInitData)
	assert.Nil(t, err)
}
