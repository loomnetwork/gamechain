package cmd

import (
	"fmt"

	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var getHeroForUserCmdArgs struct {
	userID string
	heroID int64
}

var getHeroForUserCmd = &cobra.Command{
	Use:   "get_hero_user",
	Short: "get hero for user",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := zb.GetHeroForUserRequest{
			UserId: getHeroForUserCmdArgs.userID,
			HeroId: getHeroForUserCmdArgs.heroID,
		}
		result := zb.GetHeroForUserResponse{}

		_, err := commonTxObjs.contract.StaticCall("GetHeroForUser", &req, callerAddr, &result)
		if err != nil {
			return err
		}

		fmt.Printf("hero_id: %d\n", result.Hero.HeroId)
		fmt.Printf("experience: %d\n", result.Hero.Experience)
		fmt.Printf("level: %d\n", result.Hero.Level)
		for _, skill := range result.Hero.Skills {
			fmt.Printf("skill_title: %s\n", skill.Title)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getHeroForUserCmd)

	getHeroForUserCmd.Flags().StringVarP(&getHeroForUserCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	getHeroForUserCmd.Flags().Int64VarP(&getHeroForUserCmdArgs.heroID, "heroId", "", 1, "heroID of hero")
}