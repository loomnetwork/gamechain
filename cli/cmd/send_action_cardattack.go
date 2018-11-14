package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var sendActionCardAttackCmdArgs struct {
	matchID            int64
	userID             string
	actionType         int32
	cardPlayInstanceID int32
	attackerID         int32
	targetID           int32
	objectType         int32
}

var sendActionCardAttackCmd = &cobra.Command{
	Use:   "send_action_cardattack",
	Short: "send_action_cardattack",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)

		var otype zb.AffectObjectType
		switch sendActionCardAttackCmdArgs.objectType {
		case 0:
			otype = zb.AffectObjectType_PLAYER
		case 1:
			otype = zb.AffectObjectType_CHARACTER
		}

		var req = zb.PlayerActionRequest{
			MatchId: sendActionCardAttackCmdArgs.matchID,
			PlayerAction: &zb.PlayerAction{
				ActionType: zb.PlayerActionType_CardAttack,
				PlayerId:   sendActionCardAttackCmdArgs.userID,
				Action: &zb.PlayerAction_CardAttack{
					CardAttack: &zb.PlayerActionCardAttack{
						Attacker: &zb.CardInstance{
							InstanceId: sendActionCardAttackCmdArgs.attackerID,
						},
						AffectObjectType: otype,
						Target: &zb.Unit{
							InstanceId: sendActionCardAttackCmdArgs.targetID,
						},
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
			fmt.Printf("sent action cardattack successfully")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sendActionCardAttackCmd)
	sendActionCardAttackCmd.Flags().Int64VarP(&sendActionCardAttackCmdArgs.matchID, "matchId", "m", 0, "Match Id")
	sendActionCardAttackCmd.Flags().StringVarP(&sendActionCardAttackCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	sendActionCardAttackCmd.Flags().Int32VarP(&sendActionCardAttackCmdArgs.attackerID, "attackerID", "a", 0, "Attacker ID")
	sendActionCardAttackCmd.Flags().Int32VarP(&sendActionCardAttackCmdArgs.targetID, "targetID", "g", 0, "Target ID")
	sendActionCardAttackCmd.Flags().Int32VarP(&sendActionCardAttackCmdArgs.objectType, "objectType", "o", 0, "Object Type")
}
