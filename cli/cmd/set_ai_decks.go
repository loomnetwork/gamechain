package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var setAIDecksCmdArgs struct {
	data    string
	version string
}

var setAIDecksCmd = &cobra.Command{
	Use:   "set_ai_decks",
	Short: "set AI decks",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var deckList zb.DeckList

		if err := json.Unmarshal([]byte(setAIDecksCmdArgs.data), &deckList); err != nil {
			return err
		}

		req := &zb.SetAIDecksRequest{
			Decks:   deckList.Decks,
			Version: setAIDecksCmdArgs.version,
		}

		_, err := commonTxObjs.contract.Call("SetAIDecks", req, signer, nil)
		if err != nil {
			return err
		}
		fmt.Printf("decks set successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setAIDecksCmd)

	setAIDecksCmd.Flags().StringVarP(&setAIDecksCmdArgs.data, "data", "d", "{\"hero_id\":1, \"name\": \"NewDeck\", \"cards\": [ {\"card_name\": \"Pyromaz\", \"amount\": 2}, {\"card_name\": \"Burrrnn\", \"amount\": 1} ]}", "Deck data in serialized json format")
	setAIDecksCmd.Flags().StringVarP(&setAIDecksCmdArgs.version, "version", "v", "v1", "Version")
}
