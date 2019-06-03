package cmd

import (
	"fmt"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var getStateCmdArgs struct {
	MatchID int64
}

var getStateCmd = &cobra.Command{
	Use:   "get_state",
	Short: "get state",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}
		var req zb_calls.GetGamechainStateRequest
		var resp zb_calls.GetGamechainStateResponse
		_, err := commonTxObjs.contract.StaticCall("GetState", &req, callerAddr, &resp)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			output, err := new(jsonpb.Marshaler).MarshalToString(resp.State)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
		default:
			fmt.Printf("%+v", resp.State)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getStateCmd)
}
