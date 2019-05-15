package main

import (
	"encoding/hex"
	"fmt"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/battleground"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/plugin"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/pkg/errors"
	"io/ioutil"
)

var initRequest = zb_calls.InitRequest {
}

var updateInitRequest = zb_calls.UpdateInitRequest {
}

func readJsonFileToProtobuf(filename string, message proto.Message) error {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	json := string(bytes)
	if err := jsonpb.UnmarshalString(json, message); err != nil {
		return errors.Wrap(err, "error parsing JSON file " + filename)
	}

	return nil
}

func setup(c *battleground.ZombieBattleground, pubKeyHex string, addr *loom.Address, ctx *contract.Context) error {
	updateInitRequest.InitData = &zb_data.InitData{}
	err := readJsonFileToProtobuf("simple-init.json", updateInitRequest.InitData)
	if err != nil {
		return err
	}

	initRequest = zb_calls.InitRequest{
		DefaultDecks:         updateInitRequest.InitData.DefaultDecks,
		DefaultCollection:    updateInitRequest.InitData.DefaultCollection,
		Cards:                updateInitRequest.InitData.Cards,
		Overlords:            updateInitRequest.InitData.Overlords,
		AiDecks:              updateInitRequest.InitData.AiDecks,
		Version:              updateInitRequest.InitData.Version,
		Oracle:               updateInitRequest.InitData.Oracle,
		OverlordLeveling:     updateInitRequest.InitData.OverlordLeveling,
	}

	pubKey, _ := hex.DecodeString(pubKeyHex)

	addr = &loom.Address{
		Local: loom.LocalAddressFromPublicKey(pubKey),
	}

	*ctx = contract.WrapPluginContext(
		plugin.CreateFakeContext(*addr, *addr),
	)

	err = c.Init(*ctx, &initRequest)
	if err != nil {
		return err
	}

	return nil
}

func setupAccount(c *battleground.ZombieBattleground, ctx contract.Context, upsertAccountRequest *zb_calls.UpsertAccountRequest) {
	err := c.CreateAccount(ctx, upsertAccountRequest)
	if err != nil {
		panic(err)
	}
}
func setupZBContract() {

	var pubKeyHexString = "e4008e26428a9bca87465e8de3a8d0e9c37a56ca619d3d6202b0567528786618"
	var addr loom.Address

	setup(zvContract, pubKeyHexString, &addr, &ctx)
	setupAccount(zvContract, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "AccountUser",
		Image:   "PathToImage",
		Version: "v1",
	})

}
func listItemsForPlayer(playerId int) []string {
	res := []string{}

	cardCollection, err := zvContract.GetCollection(ctx, &zb_calls.GetCollectionRequest{
		UserId: "AccountUser",
	})
	if err != nil {
		panic(err)
	}
	for _, v := range cardCollection.Cards {
		res = append(res, fmt.Sprintf("Mould Id %d", v.MouldId))
	}

	return res
}

var zvContract *battleground.ZombieBattleground
var ctx contract.Context

func main() {
	zvContract = &battleground.ZombieBattleground{}
	setupZBContract()

	runGocui()
	return
}
