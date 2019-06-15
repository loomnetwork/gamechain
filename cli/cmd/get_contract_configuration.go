package cmd

import (
	"fmt"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"strings"

	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var getContractConfigurationCmdArgs struct {
}

var getContractConfigurationCmd = &cobra.Command{
	Use:   "get_contract_configuration",
	Short: "get contract configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		var resp zb_calls.GetContractConfigurationResponse
		_, err := commonTxObjs.contract.StaticCall("GetContractConfiguration", &zb_calls.EmptyRequest{}, callerAddr, &resp)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			err := printProtoMessageAsJSONToStdout(resp.Configuration)
			if err != nil {
				return err
			}
		default:
			fmt.Printf("%+v\n", resp.Configuration)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getContractConfigurationCmd)
}
