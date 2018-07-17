package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var addDeckCmdArgs struct {
	userId string
	value  string
}

var addDeckCmd = &cobra.Command{
	Use:   "addDeck",
	Short: "add deck in zombie battleground",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var deck zb.ZBDeck

		if err := json.Unmarshal([]byte(addDeckCmdArgs.value), &deck); err != nil {
			return fmt.Errorf("invalid JSON passed in value field. Error: %s\n", err.Error())
		}

		req := &zb.AddDeckRequest{
			Deck:   &deck,
			UserId: addDeckCmdArgs.userId,
		}

		_, err := commonTxObjs.contract.Call("AddDeck", req, signer, nil)
		if err != nil {
			return fmt.Errorf("error encountered while calling AddDeck: %s\n", err.Error())
		} else {
			return nil
		}
	},
}

func init() {
	rootCmd.AddCommand(addDeckCmd)

	addDeckCmd.Flags().StringVarP(&addDeckCmdArgs.userId, "userId", "u", "loom", "UserId of account")
	addDeckCmd.Flags().StringVarP(&addDeckCmdArgs.value, "value", "v", "{\"heroId\":\"1\", \"name\": \"NewDeck\", \"cards\": [ {\"card_id\": 1, \"amount\": 2}, {\"card_id\": 2, \"amount\": 1} ]}", "Deck data in serialized json format")
}
