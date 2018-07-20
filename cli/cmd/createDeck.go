package cmd

import (
	"encoding/json"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var createDeckCmdArgs struct {
	userID string
	value  string
}

var createDeckCmd = &cobra.Command{
	Use:   "createDeck",
	Short: "create a deck",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var deck zb.Deck

		if err := json.Unmarshal([]byte(createDeckCmdArgs.value), &deck); err != nil {
			return err
		}

		req := &zb.CreateDeckRequest{
			Deck:   &deck,
			UserId: createDeckCmdArgs.userID,
		}

		_, err := commonTxObjs.contract.Call("CreateDeck", req, signer, nil)
		if err != nil {
			return err
		}
		return nil

	},
}

func init() {
	rootCmd.AddCommand(createDeckCmd)

	createDeckCmd.Flags().StringVarP(&createDeckCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	createDeckCmd.Flags().StringVarP(&createDeckCmdArgs.value, "value", "v", "{\"heroId\":\"1\", \"name\": \"NewDeck\", \"cards\": [ {\"card_id\": 1, \"amount\": 2}, {\"card_id\": 2, \"amount\": 1} ]}", "Deck data in serialized json format")
}
