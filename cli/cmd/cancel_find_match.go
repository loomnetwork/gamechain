package cmd

import (
	"fmt"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var cancelFindMatchCmdArgs struct {
	userID  string
	matchID int64
}

var cancelFindMatchCmd = &cobra.Command{
	Use:   "cancel_find_match",
	Short: "cancel find match for zombiebattleground",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var req = zb.CancelFindMatchRequest{
			UserId:  cancelFindMatchCmdArgs.userID,
			MatchId: cancelFindMatchCmdArgs.matchID,
		}

		_, err := commonTxObjs.contract.Call("CancelFindMatch", &req, signer, nil)
		if err != nil {
			return err
		}
		fmt.Printf("Cancel match %d successfully", req.MatchId)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cancelFindMatchCmd)

	cancelFindMatchCmd.Flags().StringVarP(&cancelFindMatchCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	cancelFindMatchCmd.Flags().Int64VarP(&cancelFindMatchCmdArgs.matchID, "matchId", "m", 0, "Match Id")
}
