package cmd

import (
	"fmt"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var sendActionCmdArgs struct {
	userID  string
	message string
}

var sendActionCmd = &cobra.Command{
	Use:   "send_action",
	Short: "send_action",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var req = zb.PlayerActionRequest{
			PlayerAction: &zb.PlayerAction{
				PlayerId: sendActionCmdArgs.userID,
			},
		}

		_, err := commonTxObjs.contract.Call("SendAction", &req, signer, nil)
		if err != nil {
			return err
		}
		fmt.Printf("sent action %v", req)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sendActionCmd)
	sendActionCmd.Flags().StringVarP(&sendActionCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	sendActionCmd.Flags().StringVarP(&sendActionCmdArgs.message, "message", "m", "hello loom", "Message")
}
