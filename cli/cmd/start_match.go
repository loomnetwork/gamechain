package cmd

import (
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var startMatchCmdArgs struct {
	userID string
}

var startMatchCmd = &cobra.Command{
	Use:   "start_match",
	Short: "sample start match",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var req zb.StartMatchRequest
		req.UserId = startMatchCmdArgs.userID

		_, err := commonTxObjs.contract.Call("StartMatch", &req, signer, nil)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(startMatchCmd)

	startMatchCmd.Flags().StringVarP(&startMatchCmdArgs.userID, "userId", "u", "loom", "UserId of account")
}
