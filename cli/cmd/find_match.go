package cmd

import (
	"fmt"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var findMatchCmdArgs struct {
	userID string
	deckID int64
}

var findMatchCmd = &cobra.Command{
	Use:   "find_match",
	Short: "find match for zombiebattleground",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var req = zb.FindMatchRequest{
			UserId: findMatchCmdArgs.userID,
			DeckId: findMatchCmdArgs.deckID,
		}
		var resp zb.FindMatchResponse

		req.UserId = findMatchCmdArgs.userID

		_, err := commonTxObjs.contract.Call("FindMatch", &req, signer, &resp)
		if err != nil {
			return err
		}
		fmt.Printf("find match: %v", resp.Match)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(findMatchCmd)

	findMatchCmd.Flags().StringVarP(&findMatchCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	findMatchCmd.Flags().Int64VarP(&findMatchCmdArgs.deckID, "deckId", "d", 1, "Deck Id")
}
