package cmd

import (
	"fmt"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"strings"

	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var getContractStateCmdArgs struct {
}

var getContractStateCmd = &cobra.Command{
	Use:   "get_contract_state",
	Short: "get contract state",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		var resp zb_calls.GetContractStateResponse
		_, err := commonTxObjs.contract.StaticCall("GetContractState", &zb_calls.EmptyRequest{}, callerAddr, &resp)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			err := printProtoMessageAsJSONToStdout(resp.State)
			if err != nil {
				return err
			}
		default:
			fmt.Printf("%+v\n", resp.State)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getContractStateCmd)
}
