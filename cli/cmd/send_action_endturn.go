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

var sendActionEndturnCmdArgs struct {
	matchID int64
	userID  string
}

var sendActionEndturnCmd = &cobra.Command{
	Use:   "send_action_endturn",
	Short: "send_action_endturn",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)

		var req = zb_calls.PlayerActionRequest{
			MatchId: sendActionEndturnCmdArgs.matchID,
			PlayerAction: &zb_data.PlayerAction{
				ActionType: zb_enums.PlayerActionType_EndTurn,
				PlayerId:   sendActionEndturnCmdArgs.userID,
				Action: &zb_data.PlayerAction_EndTurn{
					EndTurn: &zb_data.PlayerActionEndTurn{},
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
			fmt.Println("sent action endturn successfully")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sendActionEndturnCmd)
	sendActionEndturnCmd.Flags().Int64VarP(&sendActionEndturnCmdArgs.matchID, "matchId", "m", 0, "Match Id")
	sendActionEndturnCmd.Flags().StringVarP(&sendActionEndturnCmdArgs.userID, "userId", "u", "loom", "UserId of account")
}
