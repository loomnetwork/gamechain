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

var sendActionCardPlayCmdArgs struct {
	matchID            int64
	userID             string
	cardPlayInstanceID int32
	attackerID         int32
	targetID           int32
	objectType         int32
}

var sendActionCardPlayCmd = &cobra.Command{
	Use:   "send_action_cardplay",
	Short: "send_action_cardplay",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)

		var req = zb_calls.PlayerActionRequest{
			MatchId: sendActionCardPlayCmdArgs.matchID,
			PlayerAction: &zb_data.PlayerAction{
				ActionType: zb_enums.PlayerActionType_CardPlay,
				PlayerId:   sendActionCardPlayCmdArgs.userID,
				Action: &zb_data.PlayerAction_CardPlay{
					CardPlay: &zb_data.PlayerActionCardPlay{
						Card: &zb_data.InstanceId{Id: sendActionCardPlayCmdArgs.cardPlayInstanceID},
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
			fmt.Println("sent action cardplay successfully")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sendActionCardPlayCmd)
	sendActionCardPlayCmd.Flags().Int64VarP(&sendActionCardPlayCmdArgs.matchID, "matchId", "m", 0, "Match Id")
	sendActionCardPlayCmd.Flags().StringVarP(&sendActionCardPlayCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	sendActionCardPlayCmd.Flags().Int32VarP(&sendActionCardPlayCmdArgs.cardPlayInstanceID, "instanceId", "i", 1, "card instance id for card play")
	sendActionCardPlayCmd.Flags().Int32VarP(&sendActionCardPlayCmdArgs.attackerID, "attackerID", "a", 0, "Attacker ID")
	sendActionCardPlayCmd.Flags().Int32VarP(&sendActionCardPlayCmdArgs.targetID, "targetID", "g", 0, "Target ID")
	sendActionCardPlayCmd.Flags().Int32VarP(&sendActionCardPlayCmdArgs.objectType, "objectType", "o", 0, "Object Type")
}
