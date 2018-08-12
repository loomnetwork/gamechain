package cmd

import (
	"fmt"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var deleteDeckCmdArgs struct {
	userID string
	deckId int64
}

var deleteDeckCmd = &cobra.Command{
	Use:   "delete_deck",
	Short: "deletes deck for zombiebattleground by its id",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)

		req := &zb.DeleteDeckRequest{
			UserId: deleteDeckCmdArgs.userID,
			DeckId: deleteDeckCmdArgs.deckId,
		}

		_, err := commonTxObjs.contract.Call("DeleteDeck", req, signer, nil)
		if err != nil {
			return err
		}
		fmt.Printf("deck deleted successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteDeckCmd)

	deleteDeckCmd.Flags().StringVarP(&deleteDeckCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	deleteDeckCmd.Flags().Int64VarP(&deleteDeckCmdArgs.deckId, "deckId", "", 0, "DeckId of account")
}
