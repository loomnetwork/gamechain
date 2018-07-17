package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var editDeckCmdArgs struct {
	userId string
	value  string
}

var editDeckCmd = &cobra.Command{
	Use:   "editDeck",
	Short: "edit deck in zombie battleground",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var deck zb.ZBDeck

		if err := json.Unmarshal([]byte(editDeckCmdArgs.value), &deck); err != nil {
			return fmt.Errorf("invalid JSON passed in value field. Error: %s\n", err.Error())
		}

		req := &zb.EditDeckRequest{
			Deck:   &deck,
			UserId: editDeckCmdArgs.userId,
		}

		_, err := commonTxObjs.contract.Call("EditDeck", req, signer, nil)
		if err != nil {
			return fmt.Errorf("error encountered while calling EditDeck: %s\n", err.Error())
		} else {
			return nil
		}
	},
}

func init() {
	rootCmd.AddCommand(editDeckCmd)

	editDeckCmd.Flags().StringVarP(&editDeckCmdArgs.userId, "userId", "u", "loom", "UserId of account")
	editDeckCmd.Flags().StringVarP(&editDeckCmdArgs.value, "value", "v", "{\"heroId\":\"1\", \"name\": \"NewDeck\", \"cards\": [ {\"card_id\": 1, \"amount\": 2}, {\"card_id\": 5, \"amount\": 2} ]}", "Deck data in serialized json format")
}
