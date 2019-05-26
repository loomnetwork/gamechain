package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var deleteDeckCmdArgs struct {
	userID string
	deckID int64
	version string
}

var deleteDeckCmd = &cobra.Command{
	Use:   "delete_deck",
	Short: "deletes deck for zombiebattleground by its id",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)

		req := &zb.DeleteDeckRequest{
			UserId: deleteDeckCmdArgs.userID,
			DeckId: deleteDeckCmdArgs.deckID,
			Version: deleteDeckCmdArgs.version,
		}

		_, err := commonTxObjs.contract.Call("DeleteDeck", req, signer, nil)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			output, err := json.Marshal(map[string]interface{}{"success": true})
			if err != nil {
				return err
			}
			fmt.Println(string(output))
		default:
			fmt.Printf("deck deleted successfully")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteDeckCmd)

	deleteDeckCmd.Flags().StringVarP(&deleteDeckCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	deleteDeckCmd.Flags().Int64VarP(&deleteDeckCmdArgs.deckID, "deckId", "", 0, "DeckId of account")
	deleteDeckCmd.Flags().StringVarP(&deleteDeckCmdArgs.version, "version", "v", "v1", "Version")

	_ = deleteDeckCmd.MarkFlagRequired("version")
}
