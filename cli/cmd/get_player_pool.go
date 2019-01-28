package cmd

import (
	"fmt"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom"
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
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}
		var req = zb.PlayerPoolRequest{}
		var resp zb.PlayerPoolResponse

		if getPlayerPoolCmdArgs.isTaggedPlayerPool {
			_, err := commonTxObjs.contract.StaticCall("GetTaggedPlayerPool", &req, callerAddr, &resp)
			if err != nil {
				return err
			}
		} else {
			_, err := commonTxObjs.contract.StaticCall("GetPlayerPool", &req, callerAddr, &resp)
			if err != nil {
				return err
			}
		}

		pool := resp.Pool

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			output, err := new(jsonpb.Marshaler).MarshalToString(pool)
			if err != nil {
				return err
			}
			fmt.Println(output)
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
