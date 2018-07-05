package cmd

import (
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"

	"github.com/loomnetwork/zombie_battleground/types/zb"

	"fmt"
)

var createAccCmdArgs struct {
	userName *string
	image    *string
}

var createAccountCmd = &cobra.Command{
	Use:   "createAccount",
	Short: "creates an account for zombiechain",
	Run: func(cmd *cobra.Command, args []string) {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var result interface{}

		accountData := &zb.Account{
			Username: *createAccCmdArgs.userName,
			Image:    *createAccCmdArgs.image,
		}

		_, err := commonTxObjs.contract.Call("CreateAccount", accountData, signer, result)
		if err != nil {
			fmt.Printf("Error encountered while calling CreateAccount: %s", err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(createAccountCmd)

	createAccCmdArgs.userName = createAccountCmd.Flags().StringP("username", "u", "", "Username of account")
	createAccCmdArgs.image = createAccountCmd.Flags().StringP("image", "i", "", "Image of user account")
}
