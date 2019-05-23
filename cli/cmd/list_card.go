package cmd

import (
	"fmt"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"

	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var listCardCmdArgs struct {
	version string
}

var listCardCmd = &cobra.Command{
	Use:   "list_card_library",
	Short: "list card_library",
	RunE: func(cmd *cobra.Command, args []string) error {

		if listCardCmdArgs.version == "" {
			return fmt.Errorf("version not specified")
		}

		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := zb_calls.ListCardLibraryRequest{
			Version: listCardCmdArgs.version,
		}
		result := zb_calls.ListCardLibraryResponse{}

		_, err := commonTxObjs.contract.StaticCall("ListCardLibrary", &req, callerAddr, &result)
		if err != nil {
			return err
		}

		return printProtoMessageAsJSONToStdout(&result)
	},
}

func init() {
	rootCmd.AddCommand(listCardCmd)
	listCardCmd.Flags().StringVarP(&listCardCmdArgs.version, "version", "v", "", "Version")

	_ = listCardCmd.MarkFlagRequired("version")
}
