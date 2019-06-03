package cmd

import (
	"fmt"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"strings"

	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var listOverlordsForUserCmdArgs struct {
	userID  string
	version string
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

		req := zb_calls.ListOverlordUserInstancesRequest{
			UserId:  listOverlordsForUserCmdArgs.userID,
			Version: listOverlordsForUserCmdArgs.version,
		}
		result := zb_calls.ListOverlordUserInstancesResponse{}

		_, err := commonTxObjs.contract.StaticCall("ListOverlordUserInstances", &req, callerAddr, &result)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			return printProtoMessageAsJSONToStdout(&result)
		default:
			for _, overlordInfo := range result.Overlords {
				fmt.Printf("overlord_id: %d\n", overlordInfo.Prototype.Id)
				fmt.Printf("experience: %d\n", overlordInfo.UserData.Experience)
				fmt.Printf("level: %d\n", overlordInfo.UserData.Level)
				for _, skill := range overlordInfo.Prototype.Skills {
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
	listOverlordsForUserCmd.Flags().StringVarP(&listOverlordsForUserCmdArgs.version, "version", "v", "v1", "Version")

	_ = listOverlordsForUserCmd.MarkFlagRequired("version")
}
