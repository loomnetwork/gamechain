package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var setTutorialRewardAmountCmdArgs struct {
	amount uint64
}

var setTutorialRewardAmountCmd = &cobra.Command{
	Use:   "set_tutorial_reward_amount",
	Short: "set tutorial reward amount",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := &loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := zb_calls.SetTutorialRewardAmountRequest{
			Amount: setTutorialRewardAmountCmdArgs.amount,
			Oracle: callerAddr.MarshalPB(),
		}

		_, err := commonTxObjs.contract.Call("SetTutorialRewardAmount", &req, signer, nil)
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
			fmt.Println("success")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setTutorialRewardAmountCmd)

	setTutorialRewardAmountCmd.Flags().Uint64VarP(&setTutorialRewardAmountCmdArgs.amount, "amount", "n", 1, "amount")
}
