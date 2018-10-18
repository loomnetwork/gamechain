package cmd

import (
	"fmt"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var sendActionCmdArgs struct {
	matchID            int64
	userID             string
	actionType         int32
	cardPlayInstanceID int32
	attackerID         int32
	targetID           int32
	objectType         int32
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
		case zb.PlayerActionType_DrawCard:
			req.PlayerAction.Action = &zb.PlayerAction_DrawCard{
				DrawCard: &zb.PlayerActionDrawCard{
					CardInstance: &zb.CardInstance{},
				},
			}
		case zb.PlayerActionType_CardAttack:
			var otype zb.AffectObjectType
			switch sendActionCmdArgs.objectType {
			case 0:
				otype = zb.AffectObjectType_PLAYER
			case 1:
				otype = zb.AffectObjectType_CHARACTER
			}
			req.PlayerAction.Action = &zb.PlayerAction_CardAttack{
				CardAttack: &zb.PlayerActionCardAttack{
					Attacker: &zb.CardInstance{
						InstanceId: sendActionCmdArgs.attackerID,
					},
					AffectObjectType: otype,
					Target: &zb.Unit{
						InstanceId: sendActionCmdArgs.targetID,
					},
				},
			}
		case zb.PlayerActionType_CardPlay:
			req.PlayerAction.Action = &zb.PlayerAction_CardPlay{
				CardPlay: &zb.PlayerActionCardPlay{
					Card: &zb.CardInstance{
						InstanceId: sendActionCmdArgs.cardPlayInstanceID,
					},
				},
			}
		case zb.PlayerActionType_EndTurn:
			req.PlayerAction.Action = &zb.PlayerAction_EndTurn{
				EndTurn: &zb.PlayerActionEndTurn{},
			}
		case zb.PlayerActionType_LeaveMatch:
			req.PlayerAction.Action = &zb.PlayerAction_LeaveMatch{
				LeaveMatch: &zb.PlayerActionLeaveMatch{},
			}

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
	sendActionCmd.Flags().Int32VarP(&sendActionCmdArgs.cardPlayInstanceID, "instanceId", "i", 1, "card instance id for card play")
	sendActionCmd.Flags().Int32VarP(&sendActionCmdArgs.attackerID, "attackerID", "a", 0, "Attacker ID")
	sendActionCmd.Flags().Int32VarP(&sendActionCmdArgs.targetID, "targetID", "g", 0, "Target ID")
	sendActionCmd.Flags().Int32VarP(&sendActionCmdArgs.objectType, "objectType", "o", 0, "Object Type")
}
