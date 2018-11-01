package cmd

import (
	"fmt"

	"github.com/loomnetwork/gamechain/types/zb"
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

		var req = zb.PlayerActionRequest{
			MatchId: sendActionEndturnCmdArgs.matchID,
			PlayerAction: &zb.PlayerAction{
				ActionType: zb.PlayerActionType_EndTurn,
				PlayerId:   sendActionEndturnCmdArgs.userID,
				Action: &zb.PlayerAction_EndTurn{
					EndTurn: &zb.PlayerActionEndTurn{},
				},
			},
		}

		_, err := commonTxObjs.contract.Call("SendPlayerAction", &req, signer, nil)
		if err != nil {
			return err
		}
		fmt.Printf("sent action endturn successfully")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sendActionEndturnCmd)
	sendActionEndturnCmd.Flags().Int64VarP(&sendActionEndturnCmdArgs.matchID, "matchId", "m", 0, "Match Id")
	sendActionEndturnCmd.Flags().StringVarP(&sendActionEndturnCmdArgs.userID, "userId", "u", "loom", "UserId of account")
}
