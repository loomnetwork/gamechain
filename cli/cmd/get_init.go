
package cmd

import (
	"fmt"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"

	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var getInitCmdArgs struct {
	version string
}

var getInitCmd = &cobra.Command{
	Use:   "get_init",
	Short: "get init card collections",
	RunE: func(cmd *cobra.Command, args []string) error {

		if getInitCmdArgs.version == "" {
			return fmt.Errorf("version not specified")
		}

		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := zb_calls.GetInitRequest{
			Version: getInitCmdArgs.version,
		}
		result := zb_calls.GetInitResponse{}

		_, err := commonTxObjs.contract.StaticCall("GetInit", &req, callerAddr, &result)
		if err != nil {
			return err
		}

		return printProtoMessageAsJSONToStdout(result.InitData)
	},
}

func init() {
	rootCmd.AddCommand(getInitCmd)
	getInitCmd.Flags().StringVarP(&getInitCmdArgs.version, "version", "v", "", "Version")

	_ = getInitCmd.MarkFlagRequired("version")
}
