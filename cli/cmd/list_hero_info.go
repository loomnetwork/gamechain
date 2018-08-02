package cmd

import (
	"fmt"

	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var listHeroInfoCmdArgs struct {
	userID string
}

var listHeroInfoCmd = &cobra.Command{
	Use:   "list_hero_info",
	Short: "list hero info",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := zb.GetHeroInfoListRequest{
			UserId: listHeroInfoCmdArgs.userID,
		}
		result := zb.GetHeroInfoListResponse{}

		_, err := commonTxObjs.contract.StaticCall("ListHeroInfo", &req, callerAddr, &result)
		if err != nil {
			return err
		}

		for _, heroInfo := range result.Heroes {
			fmt.Printf("hero_id: %d\n", heroInfo.HeroId)
			fmt.Printf("experience: %d\n", heroInfo.Experience)
			for _, skill := range heroInfo.Skills {
				fmt.Printf("skill_title: %s\n", skill.Title)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listHeroInfoCmd)

	listHeroInfoCmd.Flags().StringVarP(&listHeroInfoCmdArgs.userID, "userId", "u", "loom", "UserId of account")
}
