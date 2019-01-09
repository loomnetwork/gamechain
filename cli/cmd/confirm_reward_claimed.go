package cmd

import (
	"fmt"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var confirmRewardClaimedCmdArgs struct {
	userID        string
	rewardTypeInt int
}

var confirmRewardClaimedCmd = &cobra.Command{
	Use:   "confirm_reward_claimed",
	Short: "confirm that a reward has been claimed from faucet",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)

		var rewardType string
		switch confirmRewardClaimedCmdArgs.rewardTypeInt {
		case 0:
			rewardType = "tutorial-completed"
		default:
			return fmt.Errorf("invalid reward type")
		}

		req := &zb.ConfirmRewardClaimedRequest{
			UserId:     confirmRewardClaimedCmdArgs.userID,
			RewardType: rewardType,
		}

		_, err := commonTxObjs.contract.Call("ConfirmRewardClaimed", req, signer, nil)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(confirmRewardClaimedCmd)

	confirmRewardClaimedCmd.Flags().StringVarP(&confirmRewardClaimedCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	confirmRewardClaimedCmd.Flags().IntVarP(&confirmRewardClaimedCmdArgs.rewardTypeInt, "type", "t", 0, "reward type claimed")
}
