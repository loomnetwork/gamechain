package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var sendActiondrawCardCmdArgs struct {
	matchID            int64
	userID             string
	drawCardInstanceID int32
}

var sendActiondrawCardCmd = &cobra.Command{
	Use:   "send_action_drawcard",
	Short: "send_action_drawcard",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)

		var req = zb.PlayerActionRequest{
			MatchId: sendActiondrawCardCmdArgs.matchID,
			PlayerAction: &zb.PlayerAction{
				ActionType: zb.PlayerActionType_DrawCard,
				PlayerId:   sendActiondrawCardCmdArgs.userID,
				Action: &zb.PlayerAction_DrawCard{
					DrawCard: &zb.PlayerActionDrawCard{},
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
			fmt.Printf("sent action drawCard successfully")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sendActiondrawCardCmd)
	sendActiondrawCardCmd.Flags().Int64VarP(&sendActiondrawCardCmdArgs.matchID, "matchId", "m", 0, "Match Id")
	sendActiondrawCardCmd.Flags().StringVarP(&sendActiondrawCardCmdArgs.userID, "userId", "u", "loom", "UserId of account")
}
