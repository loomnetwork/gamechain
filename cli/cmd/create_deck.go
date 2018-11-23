package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var createDeckCmdArgs struct {
	userID  string
	data    string
	version string
}

var createDeckCmd = &cobra.Command{
	Use:   "create_deck",
	Short: "create a deck",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var deck zb.Deck

		if err := json.Unmarshal([]byte(createDeckCmdArgs.data), &deck); err != nil {
			return err
		}

		req := &zb.CreateDeckRequest{
			Deck:    &deck,
			UserId:  createDeckCmdArgs.userID,
			Version: createDeckCmdArgs.version,
		}

		var result zb.CreateDeckResponse
		_, err := commonTxObjs.contract.Call("CreateDeck", req, signer, &result)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			output, err := json.Marshal(map[string]interface{}{"success": true})
			if err != nil {
				return err
			}
			fmt.Println(string(output))
		default:
			fmt.Printf("deck created successfully with id %d", result.DeckId)
		}

		return nil

	},
}

func init() {
	rootCmd.AddCommand(createDeckCmd)

	createDeckCmd.Flags().StringVarP(&createDeckCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	createDeckCmd.Flags().StringVarP(&createDeckCmdArgs.data, "data", "d", "{\"hero_id\":1, \"name\": \"NewDeck\", \"cards\": [ {\"card_name\": \"Pyromaz\", \"amount\": 2}, {\"card_name\": \"Burrrnn\", \"amount\": 1} ]}", "Deck data in serialized json format")
	createDeckCmd.Flags().StringVarP(&createDeckCmdArgs.version, "version", "v", "v1", "Version")
}
