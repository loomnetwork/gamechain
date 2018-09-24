package cmd

import (
	"fmt"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var getMatchCmdArgs struct {
	MatchID int64
}

var getMatchCmd = &cobra.Command{
	Use:   "get_match",
	Short: "sample get match",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var req = zb.GetMatchRequest{
			MatchId: getMatchCmdArgs.MatchID,
		}
		var resp zb.GetMatchResponse

		_, err := commonTxObjs.contract.Call("GetMatch", &req, signer, &resp)
		if err != nil {
			return err
		}
		fmt.Printf("ID: %d, %v\n", resp.Match.Id, resp.Match.Status)
		for _, player := range resp.Match.PlayerStates {
			fmt.Printf("\tplayer: %#v\n", player)
		}
		if gamestate := resp.GameState; gamestate != nil {
			fmt.Printf("gameStateID: %d\n", gamestate.Id)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getMatchCmd)

	getMatchCmd.Flags().Int64VarP(&getMatchCmdArgs.MatchID, "matchId", "m", 0, "Match ID")
}
