package cmd

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var getGameStateCmdArgs struct {
	MatchID int64
}

var getGameStateCmd = &cobra.Command{
	Use:   "get_game_state",
	Short: "get gamestate",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}
		var req = zb.GetGameStateRequest{
			MatchId: getGameStateCmdArgs.MatchID,
		}
		var resp zb.GetGameStateResponse
		_, err := commonTxObjs.contract.StaticCall("GetGameState", &req, callerAddr, &resp)
		if err != nil {
			return err
		}

		state := resp.GameState

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			output, err := new(jsonpb.Marshaler).MarshalToString(state)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
		default:
			formatAbility := func(abilities []*zb.CardAbilityInstance) string {
				b := new(bytes.Buffer)
				for _, a := range abilities {
					b.WriteString(fmt.Sprintf("Abilities: [%+v trigger=%v active=%v]\n", a.AbilityType, a.Trigger, a.IsActive))
				}
				return b.String()
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
					fmt.Printf("\t\tId:%-2d Name:%-14s Dmg:%2d Def:%2d Goo:%2d, Zone:%0v, OwnerIndex:%d %s\n", card.InstanceId.Id, card.Prototype.Name, card.Instance.Damage, card.Instance.Defense, card.Instance.GooCost, card.Zone, card.OwnerIndex, formatAbility(card.AbilitiesInstances))
				}
				fmt.Printf("\tcard in play (%d):\n", len(player.CardsInPlay))
				for _, card := range player.CardsInPlay {
					fmt.Printf("\t\tId:%-2d Name:%-14s Dmg:%2d Def:%2d Goo:%2d, Zone:%0v, OwnerIndex:%d %s\n", card.InstanceId.Id, card.Prototype.Name, card.Instance.Damage, card.Instance.Defense, card.Instance.GooCost, card.Zone, card.OwnerIndex, formatAbility(card.AbilitiesInstances))
				}
				fmt.Printf("\tcard in deck (%d):\n", len(player.CardsInDeck))
				for _, card := range player.CardsInDeck {
					fmt.Printf("\t\tId:%-2d Name:%-14s Dmg:%2d Def:%2d Goo:%2d, Zone:%0v, OwnerIndex:%d %s\n", card.InstanceId.Id, card.Prototype.Name, card.Instance.Damage, card.Instance.Defense, card.Instance.GooCost, card.Zone, card.OwnerIndex, formatAbility(card.AbilitiesInstances))
				}
				fmt.Printf("\tcard in graveyard (%d):\n", len(player.CardsInGraveyard))
				for _, card := range player.CardsInGraveyard {
					fmt.Printf("\t\tId:%-2d Name:%-14s Dmg:%2d Def:%2d Goo:%2d, Zone:%0v, OwnerIndex:%d %s\n", card.InstanceId.Id, card.Prototype.Name, card.Instance.Damage, card.Instance.Defense, card.Instance.GooCost, card.Zone, card.OwnerIndex, formatAbility(card.AbilitiesInstances))
				}
				fmt.Printf("\n") // extra line
			}
			fmt.Printf("Actions: count %v\n", len(state.PlayerActions))
			for i, action := range state.PlayerActions {
				if int64(i) == state.CurrentActionIndex {
					fmt.Printf("   -->> [%d] %v\n", i, action)
				} else {
					fmt.Printf("\t[%d] %v\n", i, action)
				}
			}
			fmt.Printf("Current Action Index: %v\n", state.CurrentActionIndex)
			fmt.Printf("Is ended: %v, Winner: %s\n", state.IsEnded, state.Winner)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getGameStateCmd)

	getGameStateCmd.Flags().Int64VarP(&getGameStateCmdArgs.MatchID, "matchId", "m", 0, "Match ID")
}
