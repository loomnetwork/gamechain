package cmd

import (
	"fmt"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var getDeckCmdArgs struct {
	userID string
	deckID int64
}

var getDeckCmd = &cobra.Command{
	Use:   "get_deck",
	Short: "gets deck for zombiebattleground by its id",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)

		req := &zb.GetDeckRequest{
			UserId: getDeckCmdArgs.userID,
			DeckId: getDeckCmdArgs.deckID,
		}
		var result zb.GetDeckResponse
		_, err := commonTxObjs.contract.Call("GetDeck", req, signer, &result)
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
			fmt.Printf("deck name: %v\n", result.Deck.Name)
			fmt.Printf("deck id: %v\n", result.Deck.Id)
			fmt.Printf("overlord id: %v\n", result.Deck.OverlordId)
			for _, card := range result.Deck.Cards {
				fmt.Printf("mould id: %d, amount: %d\n", card.MouldId, card.Amount)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getDeckCmd)

	getDeckCmd.Flags().StringVarP(&getDeckCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	getDeckCmd.Flags().Int64VarP(&getDeckCmdArgs.deckID, "deckId", "", 0, "DeckId of account")
}
