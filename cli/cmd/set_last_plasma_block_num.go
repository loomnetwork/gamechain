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

var setLastPlasmaBlockNumCmdArgs struct {
	blockNum uint64
}

var setLastPlasmaBlockNumCmd = &cobra.Command{
	Use:   "set_last_plasma_block_num",
	Short: "set last plasma block num",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := &loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := zb_calls.SetLastPlasmaBlockNumRequest{
			LastBlockNum: setLastPlasmaBlockNumCmdArgs.blockNum,
			Oracle:       callerAddr.MarshalPB(),
		}

		_, err := commonTxObjs.contract.Call("SetLastPlasmaBlockNum", &req, signer, nil)
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
	rootCmd.AddCommand(setLastPlasmaBlockNumCmd)

	setLastPlasmaBlockNumCmd.Flags().Uint64VarP(&setLastPlasmaBlockNumCmdArgs.blockNum, "blocknum", "n", 1, "block number")
}
