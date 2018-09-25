package cmd

import (
	"fmt"

	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/go-loom/cli"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var updateOracleCmd = &cobra.Command{
	Use:   "update_oracle (new oracle) [old oracle]",
	Short: "change the oracle or set initial oracle",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		newOracle, err := cli.ResolveAddress(args[0])
		if err != nil {
			return errors.Wrap(err, "resolve new oracle address arg")
		}
		var oldOracle loom.Address
		if len(args) > 1 {
			oldOracle, err = cli.ResolveAddress(args[1])
			if err != nil {
				return errors.Wrap(err, "resolve old oracle address arg")
			}
		}
		_, err = commonTxObjs.contract.Call("UpdateOracle", &zb.NewOracleValidator{
			NewOracle: newOracle.MarshalPB(),
			OldOracle: oldOracle.MarshalPB(),
		}, signer, nil)
		if err != nil {
			return errors.Wrap(err, "call contract")
		}
		fmt.Println("oracle changed")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateOracleCmd)
}
