package cmd

import (
	"fmt"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var sendActionCmdArgs struct {
	matchID    int64
	userID     string
	actionType int32
}

var sendActionCmd = &cobra.Command{
	Use:   "send_action",
	Short: "send_action",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)

		var req = zb.PlayerActionRequest{
			MatchId: sendActionCmdArgs.matchID,
			PlayerAction: &zb.PlayerAction{
				ActionType: zb.PlayerActionType(sendActionCmdArgs.actionType),
				PlayerId:   sendActionCmdArgs.userID,
			},
		}

		switch zb.PlayerActionType(sendActionCmdArgs.actionType) {
		case zb.PlayerActionType_DrawCardPlayer:
			req.PlayerAction.Action = &zb.PlayerAction_DrawCard{
				DrawCard: &zb.PlayerActionDrawCard{
					CardInstance: &zb.CardInstance{
						InstanceId: 1,
					},
				},
			}
		case zb.PlayerActionType_EndTurn:
			req.PlayerAction.Action = &zb.PlayerAction_EndTurn{}
		default:
			return fmt.Errorf("not support action type: %v", zb.PlayerActionType(sendActionCmdArgs.actionType))
		}

		_, err := commonTxObjs.contract.Call("SendPlayerAction", &req, signer, nil)
		if err != nil {
			return err
		}
		fmt.Printf("sent action %v", req)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sendActionCmd)
	sendActionCmd.Flags().Int64VarP(&sendActionCmdArgs.matchID, "matchId", "m", 0, "Match Id")
	sendActionCmd.Flags().StringVarP(&sendActionCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	sendActionCmd.Flags().Int32VarP(&sendActionCmdArgs.actionType, "actionType", "t", 0, "Player Action Type")
}
