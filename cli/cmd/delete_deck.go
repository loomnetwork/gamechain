package cmd

import (
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var deleteDeckCmdArgs struct {
	userID   string
	deckName string
}

var deleteDeckCmd = &cobra.Command{
	Use:   "delete_deck",
	Short: "deletes deck for zombiebattleground by its name",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)

		req := &zb.DeleteDeckRequest{
			UserId:   deleteDeckCmdArgs.userID,
			DeckName: deleteDeckCmdArgs.deckName,
		}

		_, err := commonTxObjs.contract.Call("DeleteDeck", req, signer, nil)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteDeckCmd)

	deleteDeckCmd.Flags().StringVarP(&deleteDeckCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	deleteDeckCmd.Flags().StringVarP(&deleteDeckCmdArgs.deckName, "deckName", "d", "NewDeck", "DeckName of account")
}
