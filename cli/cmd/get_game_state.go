package cmd

import (
	"fmt"

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
		fmt.Printf("============StateInfo=============\n")
		fmt.Printf("Is ended: %v, Winner: %s\n", state.IsEnded, state.Winner)
		fmt.Printf("Current Player Index: %v\n", state.CurrentPlayerIndex)

		for i, player := range state.PlayerStates {
			if state.CurrentPlayerIndex == int32(i) {
				fmt.Printf("Player%d: %s 🧟\n", i+1, player.Id)
			} else {
				fmt.Printf("Player%d: %s\n", i+1, player.Id)
			}
			fmt.Printf("\tdefense: %v\n", player.Defense)
			fmt.Printf("\tcurrent goo: %v\n", player.CurrentGoo)
			fmt.Printf("\tgoo vials: %v\n", player.GooVials)
			fmt.Printf("\thas drawn card: %v\n", player.HasDrawnCard)
			fmt.Printf("\tcard in hand (%d): %v\n", len(player.CardsInHand), player.CardsInHand)
			fmt.Printf("\tcard in play (%d): %v\n", len(player.CardsInPlay), player.CardsInPlay)
			fmt.Printf("\tcard in deck (%d): %v\n", len(player.CardsInDeck), player.CardsInDeck)
			fmt.Printf("\tcard in graveyard (%d): %v\n", len(player.CardsInGraveyard), player.CardsInGraveyard)
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
		fmt.Printf("==================================\n")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getGameStateCmd)

	getGameStateCmd.Flags().Int64VarP(&getGameStateCmdArgs.MatchID, "matchId", "m", 0, "Match ID")
}
