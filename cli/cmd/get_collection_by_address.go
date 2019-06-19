package cmd

import (
	"fmt"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"strings"

	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var getCollectionByAddressCmd = &cobra.Command{
	Use:   "get_collection_by_address",
	Short: "get collection by user address",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		var req zb_calls.GetCollectionByAddressRequest
		var result zb_calls.GetCollectionByAddressResponse
		_, err := commonTxObjs.contract.StaticCall("GetCollectionByAddress", &req, callerAddr, &result)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			err := printProtoMessageAsJSONToStdout(&result)
			if err != nil {
				return err
			}
		default:
			fmt.Printf("collection:\n")
			for _, card := range result.Cards {
				fmt.Printf("card key: [%v], amount: %d\n", card.CardKey.String(), card.Amount)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCollectionByAddressCmd)
}
