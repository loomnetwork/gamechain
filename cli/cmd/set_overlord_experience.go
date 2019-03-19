package cmd

import (
	"fmt"
	"strings"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var setOverlordExperienceCmdArgs struct {
	userID     string
	overlordID     int64
	experience int64
}

var setOverlordExperienceCmd = &cobra.Command{
	Use:   "set_overlord_experience",
	Short: "set overlord experience",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)

		req := zb.SetOverlordExperienceRequest{
			UserId:     setOverlordExperienceCmdArgs.userID,
			OverlordId:     setOverlordExperienceCmdArgs.overlordID,
			Experience: setOverlordExperienceCmdArgs.experience,
		}
		result := zb.SetOverlordExperienceResponse{}

		_, err := commonTxObjs.contract.Call("SetOverlordExperience", &req, signer, &result)
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
	rootCmd.AddCommand(setOverlordExperienceCmd)

	setOverlordExperienceCmd.Flags().StringVarP(&setOverlordExperienceCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	setOverlordExperienceCmd.Flags().Int64VarP(&setOverlordExperienceCmdArgs.overlordID, "overlordId", "i", 1, "overlordID of overlord")
	setOverlordExperienceCmd.Flags().Int64VarP(&setOverlordExperienceCmdArgs.experience, "experience", "e", 1, "experience to be set")
}
