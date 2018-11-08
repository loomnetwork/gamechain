package cmd

import (
	"fmt"

	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var getAIDecksCmdArgs struct {
	version string
}

var getAIDecksCmd = &cobra.Command{
	Use:   "get_ai_decks",
	Short: "get AI decks",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := &zb.GetAIDecksRequest{
			Version: getAIDecksCmdArgs.version,
		}
		var result zb.GetAIDecksResponse
		_, err := commonTxObjs.contract.StaticCall("GetAIDecks", req, callerAddr, &result)
		if err != nil {
			return err
		}
		fmt.Printf("deck size: %d\n", len(result.Decks))
		for _, deck := range result.Decks {
			fmt.Printf("id: %d\n", deck.Id)
			fmt.Printf("name: %s\n", deck.Name)
			for _, card := range deck.Cards {
				fmt.Printf("  card_name: %s, amount: %d\n", card.CardName, card.Amount)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getAIDecksCmd)

	getAIDecksCmd.Flags().StringVarP(&getAIDecksCmdArgs.version, "version", "v", "v1", "version")
}
