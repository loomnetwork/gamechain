package cmd

import (
	"fmt"

	"github.com/gogo/protobuf/jsonpb"

	"github.com/loomnetwork/gamechain/types/zb"
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
		req := &zb.GetAIDecksRequest{
			Version: getAIDecksCmdArgs.version,
		}
		var result zb.GetAIDecksResponse
		_, err := commonTxObjs.contract.Call("GetAIDecks", req, signer, &result)
		if err != nil {
			return err
		}

		jsonMarshaler := jsonpb.Marshaler{
			OrigName: true,
			Indent:   "  ",
		}

		j, err := jsonMarshaler.MarshalToString(&result)
		if err != nil {
			return err
		}
		fmt.Println(j)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getAIDecksCmd)

	getAIDecksCmd.Flags().StringVarP(&getAIDecksCmdArgs.version, "version", "v", "v1", "version")

	_ = getAIDecksCmd.MarkFlagRequired("version")
}
