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

var setDefaultPlayerDefenseCmdArgs struct {
	defense uint64
}

var setDefaultPlayerDefenseCmd = &cobra.Command{
	Use:   "set_default_player_defense",
	Short: "set default player defense",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := &loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := zb_calls.SetDefaultPlayerDefenseRequest{
			Defense: setDefaultPlayerDefenseCmdArgs.defense,
			Oracle:  callerAddr.MarshalPB(),
		}

		_, err := commonTxObjs.contract.Call("SetDefaultPlayerDefense", &req, signer, nil)
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
	rootCmd.AddCommand(setDefaultPlayerDefenseCmd)

	setDefaultPlayerDefenseCmd.Flags().Uint64VarP(&setDefaultPlayerDefenseCmdArgs.defense, "defense", "d", 20, "default defense")
}
