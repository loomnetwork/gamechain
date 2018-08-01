package cmd

import (
	"fmt"

	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var listCardCmd = &cobra.Command{
	Use:   "list_card",
	Short: "list card",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := zb.ListCardLibraryRequest{}
		result := zb.ListCardLibraryResponse{}

		_, err := commonTxObjs.contract.StaticCall("ListCardLibrary", &req, callerAddr, &result)
		if err != nil {
			return err
		}
		fmt.Printf("card library size: %d\n", len(result.Sets))
		for _, set := range result.Sets {
			for _, card := range set.Cards {
				fmt.Printf("card_id: %d, name: %s\n", card.Id, card.Name)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCardCmd)
}
