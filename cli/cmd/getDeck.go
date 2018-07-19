package cmd

import (
	"fmt"

	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var getDeckCmdArgs struct {
	userId string
	deckId string
}

var getDeckCmd = &cobra.Command{
	Use:   "getDeck",
	Short: "gets deck for zombiebattleground by its name",
	RunE: func(cmd *cobra.Command, args []string) error {
		var result zb.Deck

		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := &zb.GetDeckRequest{
			UserId: getDeckCmdArgs.userId,
			DeckId: getDeckCmdArgs.deckId,
		}

		_, err := commonTxObjs.contract.StaticCall("GetDeck", req, callerAddr, &result)
		if err != nil {
			return fmt.Errorf("error encountered while calling GetDeck: %s\n", err.Error())
		} else {
			fmt.Println(result)
			return nil
		}
	},
}

func init() {
	rootCmd.AddCommand(getDeckCmd)

	getDeckCmd.Flags().StringVarP(&getDeckCmdArgs.userId, "userId", "u", "loom", "UserId of account")
	getDeckCmd.Flags().StringVarP(&getDeckCmdArgs.deckId, "deckId", "d", "Default", "DeckId of account")
}
