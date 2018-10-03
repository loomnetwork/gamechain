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
		match := resp.Match
		fmt.Printf("MatchID: %d\n", match.Id)
		fmt.Printf("Status: %s\n", match.Status)
		fmt.Printf("Topic: %v\n", match.Topics)
		fmt.Printf("Players:\n")
		for i, player := range match.PlayerStates {
			fmt.Printf("\tPlayer%d: %s\n", i+1, player.Id)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getMatchCmd)

	getMatchCmd.Flags().Int64VarP(&getMatchCmdArgs.MatchID, "matchId", "m", 0, "Match ID")
}
