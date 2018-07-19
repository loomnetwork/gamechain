package cmd

import (
	"fmt"

	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var getDecksCmdArgs struct {
	userId string
}

var getDecksCmd = &cobra.Command{
	Use:   "getDecks",
	Short: "gets deck data for zombiebattleground",
	RunE: func(cmd *cobra.Command, args []string) error {
		var result zb.DeckList

		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := &zb.GetDecksRequest{
			UserId: getDecksCmdArgs.userId,
		}

		_, err := commonTxObjs.contract.StaticCall("GetDecks", req, callerAddr, &result)
		if err != nil {
			return fmt.Errorf("error encountered while calling GetDecks: %s\n", err.Error())
		} else {
			fmt.Println(result)
			return nil
		}
	},
}

func init() {
	rootCmd.AddCommand(getDecksCmd)

	getDecksCmd.Flags().StringVarP(&getDecksCmdArgs.userId, "userId", "u", "loom", "UserId of account")
}
