package cmd

import (
	"fmt"
	"strings"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var setHeroExperienceCmdArgs struct {
	userID     string
	heroID     int64
	experience int64
}

var setHeroExperienceCmd = &cobra.Command{
	Use:   "set_hero_experience",
	Short: "set hero experience",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)

		req := zb.SetHeroExperienceRequest{
			UserId:     setHeroExperienceCmdArgs.userID,
			HeroId:     setHeroExperienceCmdArgs.heroID,
			Experience: setHeroExperienceCmdArgs.experience,
		}
		result := zb.SetHeroExperienceResponse{}

		_, err := commonTxObjs.contract.Call("SetHeroExperience", &req, signer, &result)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			return printProtoMessageAsJSONToStdout(&result)
		default:
			fmt.Printf("hero_id: %d\n", result.HeroId)
			fmt.Printf("experience: %d\n", result.Experience)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setHeroExperienceCmd)

	setHeroExperienceCmd.Flags().StringVarP(&setHeroExperienceCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	setHeroExperienceCmd.Flags().Int64VarP(&setHeroExperienceCmdArgs.heroID, "heroId", "i", 1, "heroID of hero")
	setHeroExperienceCmd.Flags().Int64VarP(&setHeroExperienceCmdArgs.experience, "experience", "e", 1, "experience to be set")
}
