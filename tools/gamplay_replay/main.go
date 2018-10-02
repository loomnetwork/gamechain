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
	"github.com/loomnetwork/zombie_battleground/battleground"
	"github.com/loomnetwork/zombie_battleground/types/zb"
)

var c *battleground.ZombieBattleground
var pubKeyHexString = "e4008e26428a9bca87465e8de3a8d0e9c37a56ca619d3d6202b0567528786618"
var addr loom.Address
var ctx contract.Context

func main() {
	f, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Error: ", err)
	}

	var replayList []zb.PlayerActionEvent
	json.Unmarshal(f, &replayList)

	setup(c, pubKeyHexString, &addr, &ctx)
	// TODO: accounts

	startValidation(c, replayList)
}

func setup(c *battleground.ZombieBattleground, pubKeyHex string, addr *loom.Address, ctx *contract.Context) {

	c = &battleground.ZombieBattleground{}
	pubKey, _ := hex.DecodeString(pubKeyHex)

	addr = &loom.Address{
		Local: loom.LocalAddressFromPublicKey(pubKey),
	}

	*ctx = contract.WrapPluginContext(
		plugin.CreateFakeContext(*addr, *addr),
	)

	//err := c.Init(*ctx, &initRequest)
}

func startValidation(c *battleground.ZombieBattleground, replayActionList []zb.PlayerActionEvent) {
	for _, replayAction := range replayActionList {
		actionReq := zb.PlayerActionRequest{
			MatchId:      1, // TODO: handle better
			PlayerAction: replayAction.PlayerAction,
		}
		actionResp, err := c.SendPlayerAction(ctx, &actionReq)
		if err != nil {
			fmt.Println("error: ", err)
		}
		newGameState := actionResp.GameState
		newPlayerStates := newGameState.PlayerStates

		logPlayerStates := replayAction.Match.PlayerStates

		err = comparePlayerStates(newPlayerStates, logPlayerStates)
		if err != nil {
			fmt.Println("player states do not match: ", err)
		}

	}
}

func comparePlayerStates(newPlayerStates, logPlayerStates []*zb.PlayerState) error {
	for _, newPlayerState := range newPlayerStates {
		for _, logPlayerState := range logPlayerStates {
			if newPlayerState.Id == logPlayerState.Id {
				fmt.Println("comparing state for user ", newPlayerState.Id)
				// TODO: compare using some library??
				// hp
				if newPlayerState.Hp != logPlayerState.Hp {
					return fmt.Errorf("hp doesn't match")
				}

				// mana
				if newPlayerState.Mana != logPlayerState.Mana {
					return fmt.Errorf("mana doesn't match")
				}

				// current action

				// overlord instance

				// cardsinhand

				// cards in deck

				// deck
			}
		}
	}
	return nil
}
