package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var setRewardContractVersionCmdArgs struct {
	version uint64
}

var setRewardContractVersionCmd = &cobra.Command{
	Use:   "set_reward_contract_version",
	Short: "set reward contract version",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := &loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := zb_calls.SetRewardContractVersionRequest{
			Version: setRewardContractVersionCmdArgs.version,
			Oracle:  callerAddr.MarshalPB(),
		}

		_, err := commonTxObjs.contract.Call("SetRewardContractVersion", &req, signer, nil)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			output, err := json.Marshal(map[string]interface{}{"success": true})
			if err != nil {
				return err
			}
			fmt.Println(string(output))
		default:
			fmt.Println("success")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setRewardContractVersionCmd)

	setRewardContractVersionCmd.Flags().Uint64VarP(&setRewardContractVersionCmdArgs.version, "version", "v", 1, "version")

	_ = setRewardContractVersionCmd.MarkFlagRequired("version")
}
