package cmd

import (
	"fmt"
	"strings"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var addOverlordExperienceCmdArgs struct {
	userID     string
	overlordID int64
	experience int64
}

var addOverlordExperienceCmd = &cobra.Command{
	Use:   "add_overlord_experience",
	Short: "add overlord experience",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)

		req := zb.AddOverlordExperienceRequest{
			UserId:     addOverlordExperienceCmdArgs.userID,
			OverlordId: addOverlordExperienceCmdArgs.overlordID,
			Experience: addOverlordExperienceCmdArgs.experience,
		}
		result := zb.AddOverlordExperienceResponse{}

		_, err := commonTxObjs.contract.Call("AddOverlordExperience", &req, signer, &result)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			return printProtoMessageAsJSONToStdout(&result)
		default:
			fmt.Printf("overlord_id: %d\n", result.OverlordId)
			fmt.Printf("experience: %d\n", result.Experience)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addOverlordExperienceCmd)

	addOverlordExperienceCmd.Flags().StringVarP(&addOverlordExperienceCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	addOverlordExperienceCmd.Flags().Int64VarP(&addOverlordExperienceCmdArgs.overlordID, "overlordId", "i", 1, "overlordID of overlord")
	addOverlordExperienceCmd.Flags().Int64VarP(&addOverlordExperienceCmdArgs.experience, "experience", "e", 1, "experience to be added")
}
