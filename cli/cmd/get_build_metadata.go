package cmd

import (
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var getBuildMetadataCmd = &cobra.Command{
	Use:   "get_build_metadata",
	Short: "get contract build metadata",
	RunE: func(cmd *cobra.Command, args []string) error {

		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := zb_calls.GetContractBuildMetadataRequest{}
		result := zb_calls.GetContractBuildMetadataResponse{}

		_, err := commonTxObjs.contract.StaticCall("GetContractBuildMetadata", &req, callerAddr, &result)
		if err != nil {
			return err
		}

		return printProtoMessageAsJSONToStdout(&result)
	},
}

func init() {
	rootCmd.AddCommand(getBuildMetadataCmd)
}
