package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var updateInitCmdArgs struct {
	file       string
	oldVersion string
}

var updateInitCmd = &cobra.Command{
	Use:   "update_init",
	Short: "updates the init data for zombiebattleground",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var initData zb.InitData

		if updateInitCmdArgs.file == "" {
			return fmt.Errorf("file name not provided")
		}

		f, err := os.Open(updateInitCmdArgs.file)
		if err != nil {
			return fmt.Errorf("error reading file: %s", err.Error())
		}
		defer f.Close()

		if err := new(jsonpb.Unmarshaler).Unmarshal(f, &initData); err != nil {
			return fmt.Errorf("error parsing JSON file: %s", err.Error())
		}

		if initData.Version == "" {
			return fmt.Errorf("version not specified")
		}

		updateInitData := zb.UpdateInitRequest{
			InitData: &initData,
		}

		_, err = commonTxObjs.contract.Call("UpdateInit", &updateInitData, signer, nil)
		if err != nil {
			return fmt.Errorf("error encountered while calling UpdateInit: %s", err.Error())
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			output, err := json.Marshal(map[string]interface{}{"success": true})
			if err != nil {
				return err
			}
			fmt.Println(string(output))
		default:
			fmt.Printf("Data updated successfully\n")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateInitCmd)

	updateInitCmd.Flags().StringVarP(&updateInitCmdArgs.file, "file", "f", "", "File of init data to be updated in serialized json format")
	
	_ = updateInitCmd.MarkFlagRequired("file")
}
