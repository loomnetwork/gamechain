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

var updateCardsCmdArgs struct {
	version string
	file    string
}

var updateCardsCmd = &cobra.Command{
	Use:   "update_cards",
	Short: "updates the card list for zombiebattleground",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var updateCardsData zb.UpdateCardListRequest

		if updateCardsCmdArgs.file == "" {
			return fmt.Errorf("file name not provided")
		}

		f, err := os.Open(updateCardsCmdArgs.file)
		if err != nil {
			return fmt.Errorf("error reading file: %s", err.Error())
		}
		defer f.Close()

		if err := new(jsonpb.Unmarshaler).Unmarshal(f, &updateCardsData); err != nil {
			return fmt.Errorf("error parsing JSON file: %s", err.Error())
		}

		if updateCardsCmdArgs.version == "" {
			return fmt.Errorf("version not specified")
		}
		fmt.Printf("Updating %d cards with version %s\n", len(updateCardsData.Cards), updateCardsCmdArgs.version)

		updateCardsData.Version = updateCardsCmdArgs.version
		_, err = commonTxObjs.contract.Call("UpdateCardList", &updateCardsData, signer, nil)
		if err != nil {
			return fmt.Errorf("error encountered while calling UpdateCardList: %s", err.Error())
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			output, err := json.Marshal(map[string]interface{}{"success": true})
			if err != nil {
				return err
			}
			fmt.Println(string(output))
		default:
			fmt.Printf("Cards updated successfully\n")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCardsCmd)

	updateCardsCmd.Flags().StringVarP(&updateCardsCmdArgs.version, "version", "v", "v1", "Version")
	updateCardsCmd.Flags().StringVarP(&updateCardsCmdArgs.file, "file", "f", "", "File containing cards data to be updated in serialized json format")

	_ = updateCardsCmd.MarkFlagRequired("version")
	_ = updateCardsCmd.MarkFlagRequired("file")
}
