package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var getCollectionCmdArgs struct {
	userID string
}

var getCollectionCmd = &cobra.Command{
	Use:   "get_collection",
	Short: "get collection",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := &zb.GetCollectionRequest{
			UserId: getCollectionCmdArgs.userID,
		}
		var result zb.GetCollectionResponse
		_, err := commonTxObjs.contract.StaticCall("GetCollection", req, callerAddr, &result)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			output, err := json.Marshal(result)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
		default:
			fmt.Printf("collection:\n")
			for _, card := range result.Cards {
				fmt.Printf("card_name: %s, amount: %d\n", card.CardName, card.Amount)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCollectionCmd)

	getCollectionCmd.Flags().StringVarP(&getCollectionCmdArgs.userID, "userId", "u", "loom", "UserId of account")
}
