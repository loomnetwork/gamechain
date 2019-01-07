package cmd

import (
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var getVersionsCmd = &cobra.Command{
	Use:   "get_versions",
	Short: "get content and pvp versions",
	RunE: func(cmd *cobra.Command, args []string) error {

		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := zb.GetVersionsRequest{}
		result := zb.GetVersionsResponse{}

		_, err := commonTxObjs.contract.StaticCall("GetVersions", &req, callerAddr, &result)
		if err != nil {
			return err
		}

		return printProtoMessageAsJSONToStdout(&result)
	},
}

func init() {
	rootCmd.AddCommand(getVersionsCmd)
}
