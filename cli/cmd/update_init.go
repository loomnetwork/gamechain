package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"

	"github.com/loomnetwork/zombie_battleground/types/zb"
)

var updateInitCmdArgs struct {
	version string
	data    string
}

var updateInitCmd = &cobra.Command{
	Use:   "update_init",
	Short: "updates the init data for zombiebattleground",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var updateInitData zb.UpdateInitRequest

		if err := json.Unmarshal([]byte(updateInitCmdArgs.data), &updateInitData); err != nil {
			return fmt.Errorf("invalid JSON passed in data field. Error: %s", err.Error())
		}

		updateInitData.Version = updateInitCmdArgs.version
		_, err := commonTxObjs.contract.Call("UpdateInit", &updateInitData, signer, nil)
		if err != nil {
			return fmt.Errorf("error encountered while calling UpdateInit: %s", err.Error())
		}
		fmt.Printf("Data updated successfully\n")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateInitCmd)

	updateInitCmd.Flags().StringVarP(&updateInitCmdArgs.version, "version", "v", "1", "UserId of account")
	updateInitCmd.Flags().StringVarP(&updateInitCmdArgs.data, "data", "d", "{}", "Init data to be updated in serialized json format")
}
