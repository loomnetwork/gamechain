package cmd

import (
	"fmt"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var endMatchCmdArgs struct {
	userID  string
	matchID int64
}

var endMatchCmd = &cobra.Command{
	Use:   "end_match",
	Short: "end match for zombiebattleground",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var req = zb.EndMatchRequest{
			UserId:  endMatchCmdArgs.userID,
			MatchId: endMatchCmdArgs.matchID,
		}
		var resp zb.EndMatchResponse

		_, err := commonTxObjs.contract.Call("EndMatch", &req, signer, &resp)
		if err != nil {
			return err
		}
		fmt.Printf("left match: %v", req.MatchId)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(endMatchCmd)
	endMatchCmd.Flags().StringVarP(&endMatchCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	endMatchCmd.Flags().Int64VarP(&endMatchCmdArgs.matchID, "matchId", "m", 0, "Match ID")
}
