package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"strings"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var addSoloExperienceCmdArgs struct {
	version    string
	userId     string
	overlordId int64
	experience int64
}

var addSoloExperienceCmd = &cobra.Command{
	Use:   "add_solo_experience",
	Short: "add experience to an overlord (used for Solo mode)",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var req = zb_calls.AddSoloExperienceRequest{
			Version:    addSoloExperienceCmdArgs.version,
			UserId:     addSoloExperienceCmdArgs.userId,
			OverlordId: addSoloExperienceCmdArgs.overlordId,
			Experience: addSoloExperienceCmdArgs.experience,
		}
		var resp zb_calls.AddSoloExperienceResponse

		_, err := commonTxObjs.contract.Call("AddSoloExperience", &req, signer, &resp)
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
			fmt.Println("added experience successfully")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addSoloExperienceCmd)
	addSoloExperienceCmd.Flags().StringVarP(&addSoloExperienceCmdArgs.userId, "userId", "u", "loom", "UserId of account")
	addSoloExperienceCmd.Flags().StringVarP(&addSoloExperienceCmdArgs.version, "version", "v", "v1", "Version")
	addSoloExperienceCmd.Flags().Int64VarP(&addSoloExperienceCmdArgs.overlordId, "overlordId", "o", 0, "Overlord ID")
	addSoloExperienceCmd.Flags().Int64VarP(&addSoloExperienceCmdArgs.experience, "experience", "e", 10, "Experience number to add")

	_ = addSoloExperienceCmd.MarkFlagRequired("version")
}
