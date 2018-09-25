package cmd

import (
	"fmt"

	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
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
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := &zb.GetDeckRequest{
			UserId: getDeckCmdArgs.userID,
			DeckId: getDeckCmdArgs.deckID,
		}
		var result zb.GetDeckResponse
		_, err := commonTxObjs.contract.StaticCall("GetDeck", req, callerAddr, &result)
		if err != nil {
			return err
		}
		fmt.Printf("deck name: %v\n", result.Deck.Name)
		fmt.Printf("deck id: %v\n", result.Deck.Id)
		fmt.Printf("hero id: %v\n", result.Deck.HeroId)
		for _, card := range result.Deck.Cards {
			fmt.Printf("card_name: %s, amount: %d\n", card.CardName, card.Amount)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getDeckCmd)

	getDeckCmd.Flags().StringVarP(&getDeckCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	getDeckCmd.Flags().Int64VarP(&getDeckCmdArgs.deckID, "deckId", "", 0, "DeckId of account")
}
