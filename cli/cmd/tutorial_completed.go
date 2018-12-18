package cmd

import (
	"fmt"

	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var tutorialCompletedCmdArgs struct {
	userID string
}

var tutorialCompletedCmd = &cobra.Command{
	Use:   "tutorial_completed",
	Short: "complete the tutorial for a user and get reward",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := &zb.RewardTutorialCompletedRequest{
			UserId: tutorialCompletedCmdArgs.userID,
		}
		var resp zb.RewardTutorialCompletedResponse
		_, err := commonTxObjs.contract.Call("RewardTutorialCompleted", req, callerAddr, &resp)
		if err != nil {
			return err
		}

		fmt.Println(resp)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tutorialCompletedCmd)

	tutorialCompletedCmd.Flags().StringVarP(&tutorialCompletedCmdArgs.userID, "userId", "u", "loom", "UserId of account")
}
