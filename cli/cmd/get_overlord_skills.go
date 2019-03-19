package cmd

import (
	"fmt"
	"strings"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var getOverlordSkillsCmdArgs struct {
	userID string
	overlordID int64
}

var getOverlordSkillsCmd = &cobra.Command{
	Use:   "get_overlord_skills",
	Short: "get overlord skills",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := zb.GetOverlordSkillsRequest{
			UserId: getOverlordSkillsCmdArgs.userID,
			OverlordId: getOverlordSkillsCmdArgs.overlordID,
		}
		result := zb.GetOverlordSkillsResponse{}

		_, err := commonTxObjs.contract.StaticCall("GetOverlordSkills", &req, callerAddr, &result)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			return printProtoMessageAsJSONToStdout(&result)
		default:
			fmt.Printf("overlord_id: %d\n", result.OverlordId)
			for _, skill := range result.Skills {
				fmt.Printf("skill title: %s\n", skill.Title)
				fmt.Println(skill.SkillTargets)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getOverlordSkillsCmd)

	getOverlordSkillsCmd.Flags().StringVarP(&getOverlordSkillsCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	getOverlordSkillsCmd.Flags().Int64VarP(&getOverlordSkillsCmdArgs.overlordID, "overlordId", "", 1, "overlordID of overlord")
}
