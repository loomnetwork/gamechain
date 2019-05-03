package cmd

import (
	"fmt"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/loomnetwork/gamechain/types/zb"
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

		req := &zb.ListDecksRequest{
			UserId: listDecksCmdArgs.userID,
			Version: listDecksCmdArgs.version,
		}
		var result zb.DeckList
		_, err := commonTxObjs.contract.Call("ListDecks", req, signer, &result)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			output, err := new(jsonpb.Marshaler).MarshalToString(&result)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
		default:
			fmt.Printf("deck size: %d\n", len(result.Decks))
			for _, deck := range result.Decks {
				fmt.Printf("id: %d\n", deck.Id)
				fmt.Printf("name: %s\n", deck.Name)
				for _, card := range deck.Cards {
					fmt.Printf("  mould id: %d, amount: %d\n", card.MouldId, card.Amount)
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
