package cmd

import (
	"fmt"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var sendActionLeaveMatchCmdArgs struct {
	matchID int64
	userID  string
}

var sendActionLeaveMatchCmd = &cobra.Command{
	Use:   "send_action_leavematch",
	Short: "send_action_leavematch",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)

		var req = zb.PlayerActionRequest{
			MatchId: sendActionLeaveMatchCmdArgs.matchID,
			PlayerAction: &zb.PlayerAction{
				ActionType: zb.PlayerActionType_LeaveMatch,
				PlayerId:   sendActionLeaveMatchCmdArgs.userID,
				Action: &zb.PlayerAction_LeaveMatch{
					LeaveMatch: &zb.PlayerActionLeaveMatch{},
				},
			},
		}

		_, err := commonTxObjs.contract.Call("SendPlayerAction", &req, signer, nil)
		if err != nil {
			return err
		}
		fmt.Printf("sent action leavematch successfully")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sendActionLeaveMatchCmd)
	sendActionLeaveMatchCmd.Flags().Int64VarP(&sendActionLeaveMatchCmdArgs.matchID, "matchId", "m", 0, "Match Id")
	sendActionLeaveMatchCmd.Flags().StringVarP(&sendActionLeaveMatchCmdArgs.userID, "userId", "u", "loom", "UserId of account")
}
