package cmd

import (
	"fmt"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var leaveMatchCmdArgs struct {
	userID  string
	matchID int64
}

var leaveMatchCmd = &cobra.Command{
	Use:   "leave_match",
	Short: "leave match for zombiebattleground",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var req = zb.LeaveMatchRequest{
			UserId:  leaveMatchCmdArgs.userID,
			MatchId: leaveMatchCmdArgs.matchID,
		}
		var resp zb.LeaveMatchResponse

		_, err := commonTxObjs.contract.Call("LeaveMatch", &req, signer, &resp)
		if err != nil {
			return err
		}
		fmt.Printf("left match: %v", req.MatchId)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(leaveMatchCmd)
	leaveMatchCmd.Flags().StringVarP(&leaveMatchCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	leaveMatchCmd.Flags().Int64VarP(&leaveMatchCmdArgs.matchID, "matchId", "m", 0, "Match ID")
}
