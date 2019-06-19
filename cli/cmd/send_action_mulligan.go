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

var sendActionMulliganCmdArgs struct {
	matchID         int64
	userID          string
	mulliganedCards []int32
}

var sendActionMulliganCmd = &cobra.Command{
	Use:   "send_action_mulligan",
	Short: "send_action_mulligan",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)

		var ids []*zb_data.InstanceId
		for _, id := range sendActionMulliganCmdArgs.mulliganedCards {
			ids = append(ids, &zb_data.InstanceId{Id: id})
		}

		var req = zb_calls.PlayerActionRequest{
			MatchId: sendActionMulliganCmdArgs.matchID,
			PlayerAction: &zb_data.PlayerAction{
				ActionType: zb_enums.PlayerActionType_Mulligan,
				PlayerId:   sendActionMulliganCmdArgs.userID,
				Action: &zb_data.PlayerAction_Mulligan{
					Mulligan: &zb_data.PlayerActionMulligan{
						MulliganedCards: ids,
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
			fmt.Println("sent action mulligan successfully")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sendActionMulliganCmd)
	sendActionMulliganCmd.Flags().Int64VarP(&sendActionMulliganCmdArgs.matchID, "matchId", "m", 0, "Match Id")
	sendActionMulliganCmd.Flags().StringVarP(&sendActionMulliganCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	sendActionMulliganCmd.Flags().Int32SliceVarP(&sendActionMulliganCmdArgs.mulliganedCards, "mulliganedCards", "i", nil, "comma-separated card instance ids to mulligan")
}
