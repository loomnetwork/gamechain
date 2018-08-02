package cmd

import (
	"fmt"

	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var getHeroAbilitiesCmdArgs struct {
	userID string
	heroID int64
}

var getHeroAbilitiesCmd = &cobra.Command{
	Use:   "get_hero_abilities",
	Short: "get hero abilities",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := zb.GetHeroAbilitiesRequest{
			UserId: getHeroAbilitiesCmdArgs.userID,
			HeroId: getHeroAbilitiesCmdArgs.heroID,
		}
		result := zb.GetHeroAbilitiesResponse{}

		_, err := commonTxObjs.contract.StaticCall("GetHeroAbilities", &req, callerAddr, &result)
		if err != nil {
			return err
		}

		fmt.Printf("hero_id: %d", result.HeroId)
		for _, ability := range result.Abilities {
			fmt.Printf("ability_type: %s\n", ability.Type)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getHeroAbilitiesCmd)

	getHeroAbilitiesCmd.Flags().StringVarP(&getHeroAbilitiesCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	getHeroAbilitiesCmd.Flags().Int64VarP(&getHeroAbilitiesCmdArgs.heroID, "heroId", "", 1, "heroID of hero")
}
