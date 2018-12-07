package cmd

import (
	"fmt"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var getHeroForUserCmdArgs struct {
	userID string
	heroID int64
}

var getHeroForUserCmd = &cobra.Command{
	Use:   "get_hero",
	Short: "get hero for user",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := zb.GetHeroRequest{
			UserId: getHeroForUserCmdArgs.userID,
			HeroId: getHeroForUserCmdArgs.heroID,
		}
		result := zb.GetHeroResponse{}

		_, err := commonTxObjs.contract.StaticCall("GetHero", &req, callerAddr, &result)
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
			fmt.Printf("hero_id: %d\n", result.Hero.HeroId)
			fmt.Printf("experience: %d\n", result.Hero.Experience)
			fmt.Printf("level: %d\n", result.Hero.Level)
			for _, skill := range result.Hero.Skills {
				fmt.Printf("skill title: %s\n", skill.Title)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getHeroForUserCmd)

	getHeroForUserCmd.Flags().StringVarP(&getHeroForUserCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	getHeroForUserCmd.Flags().Int64VarP(&getHeroForUserCmdArgs.heroID, "heroId", "", 1, "heroID of hero")
}
