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

var updateOverlordLibraryCmdArgs struct {
	version string
	file    string
}

var updateOverlordLibraryCmd = &cobra.Command{
	Use:   "update_overlord_library",
	Short: "updates the overlord library",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		if updateOverlordLibraryCmdArgs.version == "" {
			return fmt.Errorf("version not specified")
		}
		if updateOverlordLibraryCmdArgs.file == "" {
			return fmt.Errorf("file name not provided")
		}

		f, err := os.Open(updateOverlordLibraryCmdArgs.file)
		if err != nil {
			return fmt.Errorf("error reading file: %s", err.Error())
		}
		defer f.Close()

		req := zb.UpdateOverlordLibraryRequest{}
		if err := new(jsonpb.Unmarshaler).Unmarshal(f, &req); err != nil {
			return fmt.Errorf("error parsing JSON file: %s", err.Error())
		}
		req.Version = updateOverlordLibraryCmdArgs.version
		_, err = commonTxObjs.contract.Call("UpdateOverlordLibrary", &req, signer, nil)
		if err != nil {
			return fmt.Errorf("error encountered while calling UpdateOverlordLibrary: %s", err.Error())
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
	rootCmd.AddCommand(updateOverlordLibraryCmd)

	updateOverlordLibraryCmd.Flags().StringVarP(&updateOverlordLibraryCmdArgs.version, "version", "v", "v1", "Version")
	updateOverlordLibraryCmd.Flags().StringVarP(&updateOverlordLibraryCmdArgs.file, "file", "f", "", "File containing cards data to be updated in serialized json format")

	_ = updateOverlordLibraryCmd.MarkFlagRequired("version")
	_ = updateOverlordLibraryCmd.MarkFlagRequired("file")
}
