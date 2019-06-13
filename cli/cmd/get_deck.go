package cmd

import (
	"fmt"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var getDeckCmdArgs struct {
	userID string
	deckID int64
	version string
}

var getDeckCmd = &cobra.Command{
	Use:   "get_deck",
	Short: "gets deck for zombiebattleground by its id",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)

		req := &zb_calls.GetDeckRequest{
			UserId: getDeckCmdArgs.userID,
			DeckId: getDeckCmdArgs.deckID,
			Version: getDeckCmdArgs.version,
		}
		var result zb_calls.GetDeckResponse
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
				fmt.Printf("card key: [%v], amount: %d\n", card.CardKey.String(), card.Amount)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getDeckCmd)

	getDeckCmd.Flags().StringVarP(&getDeckCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	getDeckCmd.Flags().Int64VarP(&getDeckCmdArgs.deckID, "deckId", "", 0, "DeckId of account")
	getDeckCmd.Flags().StringVarP(&getDeckCmdArgs.version, "version", "v", "v1", "Version")

	_ = getDeckCmd.MarkFlagRequired("version")
}
