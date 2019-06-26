package cmd

import (
	"fmt"
	"github.com/loomnetwork/gamechain/tools/battleground_utility"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"strings"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var listDecksCmdArgs struct {
	userID string
	version string
}

var listDecksCmd = &cobra.Command{
	Use:   "list_decks",
	Short: "list decks",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)

		req := &zb_calls.ListDecksRequest{
			UserId: listDecksCmdArgs.userID,
			Version: listDecksCmdArgs.version,
		}
		var result zb_data.DeckList
		_, err := commonTxObjs.contract.Call("ListDecks", req, signer, &result)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			err := battleground_utility.PrintProtoMessageAsJsonToStdout(&result)
			if err != nil {
				return err
			}
		default:
			fmt.Printf("deck size: %d\n", len(result.Decks))
			for _, deck := range result.Decks {
				fmt.Printf("id: %d\n", deck.Id)
				fmt.Printf("name: %s\n", deck.Name)
				for _, card := range deck.Cards {
					fmt.Printf("  card key: [%v], amount: %d\n", card.CardKey.String(), card.Amount)
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listDecksCmd)

	listDecksCmd.Flags().StringVarP(&listDecksCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	listDecksCmd.Flags().StringVarP(&listDecksCmdArgs.version, "version", "v", "v1", "Version")

	_ = listDecksCmd.MarkFlagRequired("version")
}
