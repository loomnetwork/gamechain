package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var initStateCmd = &cobra.Command{
	Use:   "init_state",
	Short: "initialize gamechain state",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		_, err := commonTxObjs.contract.Call("InitState", &zb.InitGamechainStateRequest{}, signer, nil)
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
			fmt.Println("state initialized")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initStateCmd)
}
