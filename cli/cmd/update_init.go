package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"

	"github.com/loomnetwork/zombie_battleground/types/zb"
)

var updateInitCmdArgs struct {
	version string
	file    string
}

var updateInitCmd = &cobra.Command{
	Use:   "update_init",
	Short: "updates the init data for zombiebattleground",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var updateInitData zb.UpdateInitRequest

		if updateInitCmdArgs.file == "" {
			return fmt.Errorf("file name not provided")
		}

		f, err := ioutil.ReadFile(updateInitCmdArgs.file)
		if err != nil {
			return fmt.Errorf("error reading file: %s", err.Error())
		}

		if err := json.Unmarshal(f, &updateInitData); err != nil {
			return fmt.Errorf("invalid JSON passed in data field. Error: %s", err.Error())
		}

		if updateInitCmdArgs.version == "" {
			return fmt.Errorf("version not specified")
		}

		updateInitData.Version = updateInitCmdArgs.version
		_, err = commonTxObjs.contract.Call("UpdateInit", &updateInitData, signer, nil)
		if err != nil {
			return fmt.Errorf("error encountered while calling UpdateInit: %s", err.Error())
		}
		fmt.Printf("Data updated successfully\n")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateInitCmd)

	updateInitCmd.Flags().StringVarP(&updateInitCmdArgs.version, "version", "v", "", "UserId of account")
	updateInitCmd.Flags().StringVarP(&updateInitCmdArgs.file, "file", "f", "", "File of init data to be updated in serialized json format")
}
