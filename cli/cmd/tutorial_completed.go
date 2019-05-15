package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var tutorialCompletedCmdArgs struct {
}

var tutorialCompletedCmd = &cobra.Command{
	Use:   "tutorial_completed",
	Short: "complete the tutorial for a user and get reward",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)

		req := &zb_calls.RewardTutorialCompletedRequest{}
		var resp zb_calls.RewardTutorialCompletedResponse
		_, err := commonTxObjs.contract.Call("RewardTutorialCompleted", req, signer, &resp)
		if err != nil {
			return err
		}

		j, _ := json.Marshal(resp)
		fmt.Println(string(j))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tutorialCompletedCmd)
}
