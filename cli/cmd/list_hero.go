package cmd

import (
	"fmt"
	"strings"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var listHeroForUserCmdArgs struct {
	userID string
}

var listHeroForUserCmd = &cobra.Command{
	Use:   "list_hero",
	Short: "list hero for user",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := zb.ListHeroesRequest{
			UserId: listHeroForUserCmdArgs.userID,
		}
		result := zb.ListHeroesResponse{}

		_, err := commonTxObjs.contract.StaticCall("ListHeroes", &req, callerAddr, &result)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			return printProtoMessageAsJSONToStdout(&result)
		default:
			for _, heroInfo := range result.Heroes {
				fmt.Printf("hero_id: %d\n", heroInfo.HeroId)
				fmt.Printf("experience: %d\n", heroInfo.Experience)
				fmt.Printf("level: %d\n", heroInfo.Level)
				for _, skill := range heroInfo.Skills {
					fmt.Printf("skill title: %s\n", skill.Title)
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listHeroForUserCmd)

	listHeroForUserCmd.Flags().StringVarP(&listHeroForUserCmdArgs.userID, "userId", "u", "loom", "UserId of account")
}
