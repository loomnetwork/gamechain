package cmd

import (
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var addHeroExperienceCmdArgs struct {
	userID     string
	heroID     int64
	experience int64
}

var addHeroExperienceCmd = &cobra.Command{
	Use:   "add_hero_experience",
	Short: "add hero experience",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)

		req := zb.AddHeroExperienceRequest{
			UserId:     addHeroExperienceCmdArgs.userID,
			HeroId:     addHeroExperienceCmdArgs.heroID,
			Experience: addHeroExperienceCmdArgs.experience,
		}
		result := zb.AddHeroExperienceResponse{}

		_, err := commonTxObjs.contract.Call("AddHeroExperience", &req, signer, &result)
		if err != nil {
			return err
		}

		printProtoMessageAsJsonToStdout(&result)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addHeroExperienceCmd)

	addHeroExperienceCmd.Flags().StringVarP(&addHeroExperienceCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	addHeroExperienceCmd.Flags().Int64VarP(&addHeroExperienceCmdArgs.heroID, "heroId", "", 1, "heroID of hero")
	addHeroExperienceCmd.Flags().Int64VarP(&addHeroExperienceCmdArgs.experience, "experience", "e", 1, "experience to be added")
}
