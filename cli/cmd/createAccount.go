package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var createAccCmdArgs struct {
	userName *string
	value    *string
}

var createAccountCmd = &cobra.Command{
	Use:   "createAccount",
	Short: "creates an account for zombiebattleground",
	Run: func(cmd *cobra.Command, args []string) {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var accountData zb.UpsertAccountRequest

		if err := json.Unmarshal([]byte(*updateAccCmdArgs.value), &accountData); err != nil {
			fmt.Printf("Invalid JSON passed in value field. Error: %s\n", err.Error())
			return
		}

		accountData.Username = *updateAccCmdArgs.userName

		_, err := commonTxObjs.contract.Call("CreateAccount", &accountData, signer, nil)
		if err != nil {
			fmt.Printf("Error encountered while calling CreateAccount: %s\n", err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(createAccountCmd)

	createAccCmdArgs.userName = createAccountCmd.Flags().StringP("username", "u", "", "Username of account")
	createAccCmdArgs.value = createAccountCmd.Flags().StringP("value", "v", "", "Account data in serialized json format")
}
