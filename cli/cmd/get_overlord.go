package cmd

import (
	"fmt"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var getOverlordForUserCmdArgs struct {
	userID string
	overlordID int64
}

var getOverlordForUserCmd = &cobra.Command{
	Use:   "get_overlord",
	Short: "get overlord for user",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := zb_calls.GetOverlordRequest{
			UserId: getOverlordForUserCmdArgs.userID,
			OverlordId: getOverlordForUserCmdArgs.overlordID,
		}
		result := zb_calls.GetOverlordResponse{}

		_, err := commonTxObjs.contract.StaticCall("GetOverlord", &req, callerAddr, &result)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			output, err := new(jsonpb.Marshaler).MarshalToString(&result)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
		default:
			fmt.Printf("overlord_id: %d\n", result.Overlord.OverlordId)
			fmt.Printf("experience: %d\n", result.Overlord.Experience)
			fmt.Printf("level: %d\n", result.Overlord.Level)
			for _, skill := range result.Overlord.Skills {
				fmt.Printf("skill title: %s\n", skill.Title)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getOverlordForUserCmd)

	getOverlordForUserCmd.Flags().StringVarP(&getOverlordForUserCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	getOverlordForUserCmd.Flags().Int64VarP(&getOverlordForUserCmdArgs.overlordID, "overlordId", "", 1, "overlordID of overlord")
}
