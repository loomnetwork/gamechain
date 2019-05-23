package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"strings"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var createAccCmdArgs struct {
	userID  string
	data    string
	version string
}

var createAccountCmd = &cobra.Command{
	Use:   "create_account",
	Short: "creates an account for zombiebattleground",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var accountData zb_calls.UpsertAccountRequest

		if createAccCmdArgs.version == "" {
			return fmt.Errorf("version not specified")
		}

		if err := json.Unmarshal([]byte(createAccCmdArgs.data), &accountData); err != nil {
			return fmt.Errorf("invalid JSON passed in data field. Error: %s", err.Error())
		}

		accountData.UserId = createAccCmdArgs.userID
		accountData.Version = createAccCmdArgs.version

		_, err := commonTxObjs.contract.Call("CreateAccount", &accountData, signer, nil)
		if err != nil {
			return fmt.Errorf("error encountered while calling CreateAccount. Error: %s", err.Error())
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			output, err := json.Marshal(map[string]interface{}{"success": true})
			if err != nil {
				return err
			}
			fmt.Println(string(output))
		default:
			fmt.Printf("account %s created successfully", createAccCmdArgs.userID)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createAccountCmd)

	createAccountCmd.Flags().StringVarP(&createAccCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	createAccountCmd.Flags().StringVarP(&createAccCmdArgs.data, "data", "d", "{\"image\":\"Image\", \"game_membership_tier\": 1}", "Account data in serialized json format")
	createAccountCmd.Flags().StringVarP(&createAccCmdArgs.version, "version", "v", "v1", "Version")

	_ = createAccountCmd.MarkFlagRequired("version")
}
