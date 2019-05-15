package cmd

import (
	"fmt"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var updateEloCmdArgs struct {
	userID string
	value  int64
}

var updateEloCmd = &cobra.Command{
	Use:   "update_elo",
	Short: "updates the user's elo score",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var requestData zb_calls.UpdateUserEloRequest

		requestData.UserId = updateEloCmdArgs.userID
		requestData.EloScore = updateEloCmdArgs.value

		_, err := commonTxObjs.contract.Call("UpdateUserElo", &requestData, signer, nil)
		if err != nil {
			return fmt.Errorf("error encountered while calling UpdateUserElo: %s", err.Error())
		}

		fmt.Println("Elo updated successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateEloCmd)

	updateEloCmd.Flags().StringVarP(&updateEloCmdArgs.userID, "userId", "u", "loom", "UserId of request")
	updateEloCmd.Flags().Int64VarP(&updateEloCmdArgs.value, "value", "v", 0, "new elo value to update")
}
