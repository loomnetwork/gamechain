package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"github.com/loomnetwork/gamechain/types/zb/zb_enums"
	"strings"

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

		var req = zb_calls.PlayerActionRequest{
			MatchId: sendActionLeaveMatchCmdArgs.matchID,
			PlayerAction: &zb_data.PlayerAction{
				ActionType: zb_enums.PlayerActionType_LeaveMatch,
				PlayerId:   sendActionLeaveMatchCmdArgs.userID,
				Action: &zb_data.PlayerAction_LeaveMatch{
					LeaveMatch: &zb_data.PlayerActionLeaveMatch{
						Reason: zb_data.PlayerActionLeaveMatch_PlayerLeave,
					},
				},
			},
		}

		_, err := commonTxObjs.contract.Call("SendPlayerAction", &req, signer, nil)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			output, err := json.Marshal(map[string]interface{}{"success": true})
			if err != nil {
				return err
			}
			fmt.Println(string(output))
		default:
			fmt.Println("sent action leavematch successfully")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sendActionLeaveMatchCmd)
	sendActionLeaveMatchCmd.Flags().Int64VarP(&sendActionLeaveMatchCmdArgs.matchID, "matchId", "m", 0, "Match Id")
	sendActionLeaveMatchCmd.Flags().StringVarP(&sendActionLeaveMatchCmdArgs.userID, "userId", "u", "loom", "UserId of account")
}
