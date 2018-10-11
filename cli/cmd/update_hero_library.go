package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var updateHeroLibraryCmdArgs struct {
	version string
	file    string
}

var updateHeroLibraryCmd = &cobra.Command{
	Use:   "update_hero_libary",
	Short: "updates the hero library",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		if updateHeroLibraryCmdArgs.version == "" {
			return fmt.Errorf("version not specified")
		}
		if updateHeroLibraryCmdArgs.file == "" {
			return fmt.Errorf("file name not provided")
		}

		f, err := ioutil.ReadFile(updateHeroLibraryCmdArgs.file)
		if err != nil {
			return fmt.Errorf("error reading file: %s", err.Error())
		}

		req := zb.UpdateHeroLibraryRequest{}
		if err := new(jsonpb.Unmarshaler).Unmarshal(bytes.NewReader(f), &req); err != nil {
			return fmt.Errorf("error parsing JSON file: %s", err.Error())
		}
		req.Version = updateHeroLibraryCmdArgs.version
		_, err = commonTxObjs.contract.Call("UpdateHeroLibrary", &req, signer, nil)
		if err != nil {
			return fmt.Errorf("error encountered while calling UpdateHeroLibrary: %s", err.Error())
		}
		fmt.Printf("Data updated successfully\n")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateHeroLibraryCmd)

	updateHeroLibraryCmd.Flags().StringVarP(&updateHeroLibraryCmdArgs.version, "version", "v", "v1", "Version")
	updateHeroLibraryCmd.Flags().StringVarP(&updateHeroLibraryCmdArgs.file, "file", "f", "", "File containing cards data to be updated in serialized json format")
}
