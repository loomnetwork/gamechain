package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"strings"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var setLastPlasmaBlockNumCmdArgs struct {
	blockNum uint64
}

var setLastPlasmaBlockNumCmd = &cobra.Command{
	Use:   "set_last_plasma_block_number",
	Short: "set last plasma block number",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		req := zb_calls.SetLastPlasmaBlockNumberRequest{
			LastPlasmachainBlockNumber: setLastPlasmaBlockNumCmdArgs.blockNum,
		}

		_, err := commonTxObjs.contract.Call("SetLastPlasmaBlockNumber", &req, signer, nil)
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

	setLastPlasmaBlockNumCmd.Flags().Uint64VarP(&setLastPlasmaBlockNumCmdArgs.blockNum, "blocknum", "n", 0, "block number")
	_ = setLastPlasmaBlockNumCmd.MarkFlagRequired("blocknum")
}
