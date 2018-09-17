package cmd

import (
	"fmt"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var acceptMatchCmdArgs struct {
	userID string
}

var acceptMatchCmd = &cobra.Command{
	Use:   "accept_match",
	Short: "sample accept match",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var req = zb.AcceptMatchRequest{
			UserId: acceptMatchCmdArgs.userID,
		}
		var resp zb.AcceptMatchResponse

		_, err := commonTxObjs.contract.Call("AcceptMatch", &req, signer, &resp)
		if err != nil {
			return err
		}
		fmt.Printf("match id %d\n", resp.Match.Id)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(acceptMatchCmd)

	acceptMatchCmd.Flags().StringVarP(&acceptMatchCmdArgs.userID, "userId", "u", "loom", "UserId of account")
}
