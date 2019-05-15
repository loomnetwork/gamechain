package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var updateAccCmdArgs struct {
	userID string
	value  string
}

var updateAccountCmd = &cobra.Command{
	Use:   "update_account",
	Short: "creates an account for zombiebattleground",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var result zb_data.Account
		var accountData zb_calls.UpsertAccountRequest

		if err := json.Unmarshal([]byte(updateAccCmdArgs.value), &accountData); err != nil {
			return fmt.Errorf("invalid JSON passed in value field. Error: %s", err.Error())
		}

		accountData.UserId = updateAccCmdArgs.userID

		_, err := commonTxObjs.contract.Call("UpdateAccount", &accountData, signer, &result)
		if err != nil {
			return fmt.Errorf("error encountered while calling UpdateAccount: %s", err.Error())
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			output, err := json.Marshal(map[string]interface{}{"success": true})
			if err != nil {
				return err
			}
			fmt.Println(string(output))
		default:
			fmt.Printf("Account updated successfully\n")
			fmt.Printf("User: %s\n", result.UserId)
			fmt.Printf("Image: %s\n", result.Image)
			fmt.Printf("Game Membership Tier: %d\n", result.GameMembershipTier)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateAccountCmd)

	updateAccountCmd.Flags().StringVarP(&updateAccCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	updateAccountCmd.Flags().StringVarP(&updateAccCmdArgs.value, "value", "v", "{\"image\":\"Image2\", \"game_membership_tier\": 2}", "Account data in serialized json format")
}
