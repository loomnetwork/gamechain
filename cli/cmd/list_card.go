package cmd

import (
	"fmt"
	"os"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom"
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

		req := zb.ListCardLibraryRequest{}
		result := zb.ListCardLibraryResponse{}

		req.Version = listCardCmdArgs.version

		_, err := commonTxObjs.contract.StaticCall("ListCardLibrary", &req, callerAddr, &result)
		if err != nil {
			return err
		}
		m := jsonpb.Marshaler{OrigName: true}

		if err := m.Marshal(os.Stdout, &result); err != nil {
			return fmt.Errorf("error parsing JSON file: %s", err.Error())
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCardCmd)
	listCardCmd.Flags().StringVarP(&listCardCmdArgs.version, "version", "v", "", "Version")
}
