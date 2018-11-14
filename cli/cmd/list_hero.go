package cmd

import (
	"fmt"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom"
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
			output, err := new(jsonpb.Marshaler).MarshalToString(&result)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
		default:
			for _, heroInfo := range result.Heroes {
				fmt.Printf("hero_id: %d\n", heroInfo.HeroId)
				fmt.Printf("experience: %d\n", heroInfo.Experience)
				for _, skill := range heroInfo.Skills {
					fmt.Printf("skill_title: %s\n", skill.Title)
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
