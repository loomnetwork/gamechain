package cmd

import (
	"fmt"
	"strings"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var listOverlordsForUserCmdArgs struct {
	userID string
}

var listOverlordsForUserCmd = &cobra.Command{
	Use:   "list_overlords",
	Short: "list overlords for user",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := zb.ListOverlordsRequest{
			UserId: listOverlordsForUserCmdArgs.userID,
		}
		result := zb.ListOverlordsResponse{}

		_, err := commonTxObjs.contract.StaticCall("ListOverlords", &req, callerAddr, &result)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			return printProtoMessageAsJSONToStdout(&result)
		default:
			for _, overlordInfo := range result.Overlords {
				fmt.Printf("overlord_id: %d\n", overlordInfo.OverlordId)
				fmt.Printf("experience: %d\n", overlordInfo.Experience)
				fmt.Printf("level: %d\n", overlordInfo.Level)
				for _, skill := range overlordInfo.Skills {
					fmt.Printf("skill title: %s\n", skill.Title)
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listOverlordsForUserCmd)

	listOverlordsForUserCmd.Flags().StringVarP(&listOverlordsForUserCmdArgs.userID, "userId", "u", "loom", "UserId of account")
}
