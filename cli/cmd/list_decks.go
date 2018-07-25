package cmd

import (
	"fmt"

	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var listDecksCmdArgs struct {
	userID string
}

var listDecksCmd = &cobra.Command{
	Use:   "list_decks",
	Short: "list decks",
	RunE: func(cmd *cobra.Command, args []string) error {

		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := &zb.ListDecksRequest{
			UserId: listDecksCmdArgs.userID,
		}
		var result zb.DeckList
		_, err := commonTxObjs.contract.StaticCall("ListDecks", req, callerAddr, &result)
		if err != nil {
			return err
		}
		for _, deck := range result.Decks {
			fmt.Printf("name: %s\n", deck.Name)
			for _, card := range deck.Cards {
				fmt.Printf("  card_id: %d, amount: %d\n", card.CardId, card.Amount)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listDecksCmd)

	listDecksCmd.Flags().StringVarP(&listDecksCmdArgs.userID, "userId", "u", "loom", "UserId of account")
}
