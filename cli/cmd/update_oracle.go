package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/go-loom/cli"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var updateOracleCmdArgs struct {
	URI     string
	chainID string
}

var updateOracleCmd = &cobra.Command{
	Use:   "update_oracle (new oracle) [old oracle]",
	Short: "change the oracle or set initial oracle",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		newOracle, err := cli.ResolveAddress(args[0], updateOracleCmdArgs.chainID, updateOracleCmdArgs.URI)
		if err != nil {
			return errors.Wrap(err, "resolve new oracle address arg")
		}
		var oldOracle loom.Address
		if len(args) > 1 {
			oldOracle, err = cli.ResolveAddress(args[1], updateOracleCmdArgs.chainID, updateOracleCmdArgs.URI)
			if err != nil {
				return errors.Wrap(err, "resolve old oracle address arg")
			}
		}
		_, err = commonTxObjs.contract.Call("UpdateOracle", &zb_calls.UpdateOracle{
			NewOracle: newOracle.MarshalPB(),
			OldOracle: oldOracle.MarshalPB(),
		}, signer, nil)
		if err != nil {
			return errors.Wrap(err, "call contract")
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			output, err := json.Marshal(map[string]interface{}{"success": true})
			if err != nil {
				return err
			}
			fmt.Println(string(output))
		default:
			fmt.Println("oracle changed")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateOracleCmd)

	updateOracleCmd.Flags().StringVarP(&updateOracleCmdArgs.URI, "URI", "u", "http://localhost:46658", "Root URI for rpc and query")
	updateOracleCmd.Flags().StringVarP(&updateOracleCmdArgs.chainID, "chainID", "c", "default", "Chain ID")

	_ = updateOracleCmd.MarkFlagRequired("URI")
}
