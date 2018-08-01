package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var createAccCmdArgs struct {
	userID string
	value  string
}

var createAccountCmd = &cobra.Command{
	Use:   "create_account",
	Short: "creates an account for zombiebattleground",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var accountData zb.UpsertAccountRequest

		if err := json.Unmarshal([]byte(createAccCmdArgs.value), &accountData); err != nil {
			return fmt.Errorf("invalid JSON passed in value field. Error: %s", err.Error())
		}

		accountData.UserId = createAccCmdArgs.userID

		_, err := commonTxObjs.contract.Call("CreateAccount", &accountData, signer, nil)
		if err != nil {
			return fmt.Errorf("error encountered while calling CreateAccount. Error: %s", err.Error())
		}
		fmt.Printf("account %s created successfully", createAccCmdArgs.userID)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createAccountCmd)

	createAccountCmd.Flags().StringVarP(&createAccCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	createAccountCmd.Flags().StringVarP(&createAccCmdArgs.value, "value", "v", "{\"image\":\"Image\", \"game_membership_tier\": 1}", "Account data in serialized json format")
}
