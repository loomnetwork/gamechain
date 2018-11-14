package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var getHeroSkillsCmdArgs struct {
	userID string
	heroID int64
}

var getHeroSkillsCmd = &cobra.Command{
	Use:   "get_hero_skills",
	Short: "get hero skills",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := zb.GetHeroSkillsRequest{
			UserId: getHeroSkillsCmdArgs.userID,
			HeroId: getHeroSkillsCmdArgs.heroID,
		}
		result := zb.GetHeroSkillsResponse{}

		_, err := commonTxObjs.contract.StaticCall("GetHeroSkills", &req, callerAddr, &result)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			output, err := json.Marshal(result)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
		default:
			fmt.Printf("hero_id: %d\n", result.HeroId)
			for _, skill := range result.Skills {
				fmt.Printf("skill_title: %s\n", skill.Title)
				fmt.Println(skill.SkillTargets)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getHeroSkillsCmd)

	getHeroSkillsCmd.Flags().StringVarP(&getHeroSkillsCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	getHeroSkillsCmd.Flags().Int64VarP(&getHeroSkillsCmdArgs.heroID, "heroId", "", 1, "heroID of hero")
}
