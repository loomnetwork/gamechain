package battleground

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/plugin"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/loomnetwork/zombie_battleground/types/zb"
)

type ReplayEntry struct {
	PlayerActionType string      `json:"playerActionType"`
	UserID           string      `json:"userId"`
	Match            interface{} `json:"match"`
	PlayerAction     interface{} `json:"playerAction"`
}

func main() {
	f, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Error: ", err)
	}

	var replayList []ReplayEntry
	json.Unmarshal(f, &replayList)

	var c *ZombieBattleground
	var pubKeyHexString = "e4008e26428a9bca87465e8de3a8d0e9c37a56ca619d3d6202b0567528786618"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx)

	startValidation(replayList)
}

func setup(c *ZombieBattleground, pubKeyHex string, addr *loom.Address, ctx *contract.Context) {

	c = &ZombieBattleground{}
	pubKey, _ := hex.DecodeString(pubKeyHex)

	addr = &loom.Address{
		Local: loom.LocalAddressFromPublicKey(pubKey),
	}

	*ctx = contract.WrapPluginContext(
		plugin.CreateFakeContext(*addr, *addr),
	)

	//err := c.Init(*ctx, &initRequest)
}

func startValidation(replayList []ReplayEntry) {
	var err error
	for _, replay := range replayList {
		var action zb.PlayerAction

		// add action
		err = gp.AddAction(&action)
		if err != nil {
			fmt.Println("error: ", err)
		}
		var logState zb.GameState
		newState := gp.State
		compareState(newState, zb.logState)
	}

}
