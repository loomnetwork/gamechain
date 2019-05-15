package cmd

import (
	"fmt"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
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

		req := &zb_calls.GetCollectionRequest{
			UserId: getCollectionCmdArgs.userID,
		}
		var result zb_calls.GetCollectionResponse
		_, err := commonTxObjs.contract.StaticCall("GetCollection", req, callerAddr, &result)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			output, err := new(jsonpb.Marshaler).MarshalToString(&result)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
		default:
			fmt.Printf("collection:\n")
			for _, card := range result.Cards {
				fmt.Printf("mould id: %d, amount: %d\n", card.MouldId, card.Amount)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCollectionCmd)

	getCollectionCmd.Flags().StringVarP(&getCollectionCmdArgs.userID, "userId", "u", "loom", "UserId of account")
}
