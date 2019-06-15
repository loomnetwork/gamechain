package cmd

import (
	"fmt"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"strings"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var getPlayerPoolCmdArgs struct {
	MatchID            int64
	isTaggedPlayerPool bool
}

var getPlayerPoolCmd = &cobra.Command{
	Use:   "get_player_pool",
	Short: "get match",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)

		var req = zb_calls.PlayerPoolRequest{}
		var resp zb_calls.PlayerPoolResponse

		if getPlayerPoolCmdArgs.isTaggedPlayerPool {
			_, err := commonTxObjs.contract.Call("GetTaggedPlayerPool", &req, signer, &resp)
			if err != nil {
				return err
			}
		} else {
			_, err := commonTxObjs.contract.Call("GetPlayerPool", &req, signer, &resp)
			if err != nil {
				return err
			}
		}

		pool := resp.Pool

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			err := printProtoMessageAsJSONToStdout(pool)
			if err != nil {
				return err
			}
		default:
			fmt.Printf("Pool: %+v\n", pool)
			fmt.Printf("Players:\n")
			for _, player := range pool.PlayerProfiles {
				fmt.Printf("\t%+v\n", player)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getPlayerPoolCmd)

	getPlayerPoolCmd.Flags().BoolVarP(&getPlayerPoolCmdArgs.isTaggedPlayerPool, "tagged", "t", false, "Tagged Player Pool")
}
