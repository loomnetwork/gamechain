package cmd

import (
	"fmt"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var deleteDeckCmdArgs struct {
	userId string
	deckId string
}

var deleteDeckCmd = &cobra.Command{
	Use:   "deleteDeck",
	Short: "deletes deck for zombiebattleground by its name",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)

		req := &zb.DeleteDeckRequest{
			UserId: deleteDeckCmdArgs.userId,
			DeckId: deleteDeckCmdArgs.deckId,
		}

		_, err := commonTxObjs.contract.Call("DeleteDeck", req, signer, nil)
		if err != nil {
			return fmt.Errorf("Error encountered while calling DeleteDeck: %s\n", err.Error())
		} else {
			return nil
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteDeckCmd)

	deleteDeckCmd.Flags().StringVarP(&deleteDeckCmdArgs.userId, "userId", "u", "loom", "UserId of account")
	deleteDeckCmd.Flags().StringVarP(&deleteDeckCmdArgs.deckId, "deckId", "d", "NewDeck", "DeckId of account")
}
