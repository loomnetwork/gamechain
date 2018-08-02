package cmd

import (
	"fmt"

	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var getHeroInfoCmdArgs struct {
	userID string
	heroID int64
}

var getHeroInfoCmd = &cobra.Command{
	Use:   "get_hero_info",
	Short: "get hero info",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := zb.GetHeroInfoRequest{
			UserId: getHeroInfoCmdArgs.userID,
			HeroId: getHeroInfoCmdArgs.heroID,
		}
		result := zb.GetHeroInfoResponse{}

		_, err := commonTxObjs.contract.StaticCall("GetHeroInfo", &req, callerAddr, &result)
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
	rootCmd.AddCommand(getHeroInfoCmd)

	getHeroInfoCmd.Flags().StringVarP(&getHeroInfoCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	getHeroInfoCmd.Flags().Int64VarP(&getHeroInfoCmdArgs.heroID, "heroId", "", 1, "heroID of hero")
}
