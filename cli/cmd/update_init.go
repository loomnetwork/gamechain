package cmd

import (
	"fmt"
	"os"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var updateInitCmdArgs struct {
	version    string
	file       string
	oldVersion string
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

		f, err := os.Open(updateInitCmdArgs.file)
		if err != nil {
			return fmt.Errorf("error reading file: %s", err.Error())
		}
		defer f.Close()

		if err := new(jsonpb.Unmarshaler).Unmarshal(f, &updateInitData); err != nil {
			return fmt.Errorf("error parsing JSON file: %s", err.Error())
		}

		if updateInitCmdArgs.version == "" {
			return fmt.Errorf("version not specified")
		}

		updateInitData.Version = updateInitCmdArgs.version
		updateInitData.OldVersion = updateInitCmdArgs.oldVersion
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

	updateInitCmd.Flags().StringVarP(&updateInitCmdArgs.version, "version", "v", "", "Version to update")
	updateInitCmd.Flags().StringVarP(&updateInitCmdArgs.file, "file", "f", "", "File of init data to be updated in serialized json format")
	updateInitCmd.Flags().StringVarP(&updateInitCmdArgs.oldVersion, "old_version", "o", "", "Old version to copy missing keys from")
}
