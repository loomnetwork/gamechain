package cmd

import (
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var getDeckCmdArgs struct {
	userID   string
	deckName string
}

var getDeckCmd = &cobra.Command{
	Use:   "get_deck",
	Short: "gets deck for zombiebattleground by its name",
	RunE: func(cmd *cobra.Command, args []string) error {

		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := &zb.GetDeckRequest{
			UserId:   getDeckCmdArgs.userID,
			DeckName: getDeckCmdArgs.deckName,
		}
		var result zb.GetDeckResponse
		_, err := commonTxObjs.contract.StaticCall("GetDeck", req, callerAddr, &result)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getDeckCmd)

	getDeckCmd.Flags().StringVarP(&getDeckCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	getDeckCmd.Flags().StringVarP(&getDeckCmdArgs.deckName, "deckName", "d", "Default", "DeckId of account")
}
