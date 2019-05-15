package cmd

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var replayGameCmdArgs struct {
	matchID           int64
	stopAtActionIndex int32
}

var replayGameCmd = &cobra.Command{
	Use:   "replay_game",
	Short: "replay_game",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)

		var req = zb_calls.ReplayGameRequest{
			MatchId:           replayGameCmdArgs.matchID,
			StopAtActionIndex: replayGameCmdArgs.stopAtActionIndex,
		}
		var resp zb_calls.ReplayGameResponse

		_, err := commonTxObjs.contract.Call("ReplayGame", &req, signer, &resp)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			output, err := new(jsonpb.Marshaler).MarshalToString(resp.GameState)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
		default:
			state := resp.GameState
			actionOutcomes := resp.ActionOutcomes
			formatAbility := func(abilities []*zb_data.CardAbilityInstance) string {
				b := new(bytes.Buffer)
				for _, a := range abilities {
					b.WriteString(fmt.Sprintf("Abilities: [%+v trigger=%v active=%v]\n", a.AbilityType, a.Trigger, a.IsActive))
				}
				return b.String()
			}

			formatAction := func(action *zb_data.PlayerAction) string {
				return fmt.Sprintf("%s: %s, %+v", action.ActionType, action.PlayerId, action.Action)
			}

			for i, player := range state.PlayerStates {
				if state.CurrentPlayerIndex == int32(i) {
					fmt.Printf("Player%d: %s 🧟\n", i+1, player.Id)
				} else {
					fmt.Printf("Player%d: %s\n", i+1, player.Id)
				}
				fmt.Printf("\tstats:\n")
				fmt.Printf("\t\tdefense: %v\n", player.Defense)
				fmt.Printf("\t\tcurrent goo: %v\n", player.CurrentGoo)
				fmt.Printf("\t\tgoo vials: %v\n", player.GooVials)
				fmt.Printf("\t\thas drawn card: %v\n", player.HasDrawnCard)
				fmt.Printf("\tmulligan (%d):\n", len(player.MulliganCards))
				for _, card := range player.MulliganCards {
					fmt.Printf("\t\tName:%s\n", card.Prototype.Name)
				}
				fmt.Printf("\tcard in hand (%d):\n", len(player.CardsInHand))
				for _, card := range player.CardsInHand {
					fmt.Printf("\t\tId:%-2d Name:%-14s Dmg:%2d Def:%2d Goo:%2d, Zone:%0v, OwnerIndex:%d %s\n", card.InstanceId.Id, card.Prototype.Name, card.Instance.Damage, card.Instance.Defense, card.Instance.Cost, card.Zone, card.OwnerIndex, formatAbility(card.AbilitiesInstances))
				}
				fmt.Printf("\tcard in play (%d):\n", len(player.CardsInPlay))
				for _, card := range player.CardsInPlay {
					fmt.Printf("\t\tId:%-2d Name:%-14s Dmg:%2d Def:%2d Goo:%2d, Zone:%0v, OwnerIndex:%d %s\n", card.InstanceId.Id, card.Prototype.Name, card.Instance.Damage, card.Instance.Defense, card.Instance.Cost, card.Zone, card.OwnerIndex, formatAbility(card.AbilitiesInstances))
				}
				fmt.Printf("\tcard in deck (%d):\n", len(player.CardsInDeck))
				for _, card := range player.CardsInDeck {
					fmt.Printf("\t\tId:%-2d Name:%-14s Dmg:%2d Def:%2d Goo:%2d, Zone:%0v, OwnerIndex:%d %s\n", card.InstanceId.Id, card.Prototype.Name, card.Instance.Damage, card.Instance.Defense, card.Instance.Cost, card.Zone, card.OwnerIndex, formatAbility(card.AbilitiesInstances))
				}
				fmt.Printf("\tcard in graveyard (%d):\n", len(player.CardsInGraveyard))
				for _, card := range player.CardsInGraveyard {
					fmt.Printf("\t\tId:%-2d Name:%-14s Dmg:%2d Def:%2d Goo:%2d, Zone:%0v, OwnerIndex:%d %s\n", card.InstanceId.Id, card.Prototype.Name, card.Instance.Damage, card.Instance.Defense, card.Instance.Cost, card.Zone, card.OwnerIndex, formatAbility(card.AbilitiesInstances))
				}
				fmt.Printf("\n") // extra line
			}
			fmt.Printf("Actions: count %v\n", len(state.PlayerActions))
			for i, action := range state.PlayerActions {
				if int64(i) == state.CurrentActionIndex {
					fmt.Printf("   -->> [%-2d] %v\n", i, formatAction(action))
				} else {
					fmt.Printf("\t[%2d] %v\n", i, formatAction(action))
				}
			}
			fmt.Printf("Current Action Index: %v\n", state.CurrentActionIndex)

			fmt.Printf("Ability Outcomes:\n")
			for i, outcome := range actionOutcomes {
				fmt.Printf("\t[%d] %v\n", i, outcome)
			}

			fmt.Printf("==================================\n")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(replayGameCmd)
	replayGameCmd.Flags().Int64VarP(&replayGameCmdArgs.matchID, "matchId", "m", 0, "Match Id")
	replayGameCmd.Flags().Int32VarP(&replayGameCmdArgs.stopAtActionIndex, "stopAt", "i", -1, "stop at action index")
}
