package cmd

import (
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var confirmRewardTutorialClaimedCmdArgs struct {
}

var confirmRewardTutorialClaimedCmd = &cobra.Command{
	Use:   "confirm_reward_tutorial_completed",
	Short: "confirm that a reward has been claimed from faucet",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)

		req := &zb.ConfirmRewardTutorialClaimedRequest{}

		_, err := commonTxObjs.contract.Call("ConfirmRewardTutorialClaimed", req, signer, nil)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(confirmRewardTutorialClaimedCmd)
}
