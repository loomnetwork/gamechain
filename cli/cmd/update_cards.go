package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"

	"github.com/loomnetwork/zombie_battleground/types/zb"
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

		f, err := ioutil.ReadFile(updateCardsCmdArgs.file)
		if err != nil {
			return fmt.Errorf("error reading file: %s", err.Error())
		}

		if err := json.Unmarshal(f, &updateCardsData); err != nil {
			return fmt.Errorf("error parsing JSON file: %s", err.Error())
		}

		if updateCardsCmdArgs.version == "" {
			return fmt.Errorf("version not specified")
		}

		updateCardsData.Version = updateCardsCmdArgs.version
		_, err = commonTxObjs.contract.Call("UpdateCardList", &updateCardsData, signer, nil)
		if err != nil {
			return fmt.Errorf("error encountered while calling UpdateCardList: %s", err.Error())
		}
		fmt.Printf("Data updated successfully\n")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCardsCmd)

	updateCardsCmd.Flags().StringVarP(&updateCardsCmdArgs.version, "version", "v", "", "Version")
	updateCardsCmd.Flags().StringVarP(&updateCardsCmdArgs.file, "file", "f", "", "File containing cards data to be updated in serialized json format")
}
