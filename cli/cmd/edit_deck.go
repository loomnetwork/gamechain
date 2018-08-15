package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var editDeckCmdArgs struct {
	userID string
	value  string
}

var editDeckCmd = &cobra.Command{
	Use:   "edit_deck",
	Short: "edit deck in zombie battleground",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var deck zb.Deck

		if err := json.Unmarshal([]byte(editDeckCmdArgs.value), &deck); err != nil {
			return fmt.Errorf("invalid JSON passed in value field. Error: %s", err.Error())
		}

		req := &zb.EditDeckRequest{
			Deck:   &deck,
			UserId: editDeckCmdArgs.userID,
		}

		_, err := commonTxObjs.contract.Call("EditDeck", req, signer, nil)
		if err != nil {
			return fmt.Errorf("error encountered while calling EditDeck: %s", err.Error())
		}
		fmt.Printf("deck edited successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(editDeckCmd)

	editDeckCmd.Flags().StringVarP(&editDeckCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	editDeckCmd.Flags().StringVarP(&editDeckCmdArgs.value, "value", "v", "{\"id\": 0, \"hero_id\":1, \"name\": \"NewDefaultDeck\", \"cards\": [ {\"card_name\": \"Pyromaz\", \"amount\": 2}, {\"card_name\": \"Burrrnn\", \"amount\": 1} ]}", "Deck data in serialized json format")
}
