package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"

	"github.com/loomnetwork/zombie_battleground/types/zb"
)

var updateAccCmdArgs struct {
	userName *string
	value    *string
}

var updateAccountCmd = &cobra.Command{
	Use:   "updateAccount",
	Short: "creates an account for zombiebattleground",
	Run: func(cmd *cobra.Command, args []string) {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var result zb.Account
		var accountData zb.UpsertAccountRequest

		if err := json.Unmarshal([]byte(*updateAccCmdArgs.value), &accountData); err != nil {
			fmt.Printf("Invalid JSON passed in value field. Error: %s\n", err.Error())
			return
		}

		accountData.Username = *updateAccCmdArgs.userName

		_, err := commonTxObjs.contract.Call("UpdateAccount", &accountData, signer, &result)
		if err != nil {
			fmt.Printf("Error encountered while calling UpdateAccount: %s\n", err.Error())
		} else {
			fmt.Println(result)
		}
	},
}

func init() {
	rootCmd.AddCommand(updateAccountCmd)

	updateAccCmdArgs.userName = updateAccountCmd.Flags().StringP("username", "u", "", "Username of account")
	updateAccCmdArgs.value = updateAccountCmd.Flags().StringP("value", "v", "", "Account data in serialized json format")
}
