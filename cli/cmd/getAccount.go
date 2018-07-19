package cmd

import (
	"fmt"

	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var getAccCmdArgs struct {
	userId string
}

var getAccountCmd = &cobra.Command{
	Use:   "getAccount",
	Short: "gets account data for zombiebattleground",
	RunE: func(cmd *cobra.Command, args []string) error {
		var result zb.Account

		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := &zb.GetAccountRequest{
			UserId: getAccCmdArgs.userId,
		}

		_, err := commonTxObjs.contract.StaticCall("GetAccount", req, callerAddr, &result)
		if err != nil {
			return fmt.Errorf("error encountered while calling GetAccount: %s\n", err.Error())
		} else {
			fmt.Println(result)
			return nil
		}
	},
}

func init() {
	rootCmd.AddCommand(getAccountCmd)

	getAccountCmd.Flags().StringVarP(&getAccCmdArgs.userId, "userId", "u", "loom", "UserId of account")
}
