package cmd

import (
	"fmt"
	"strings"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var setHeroLevelCmdArgs struct {
	userID string
	heroID int64
	level  int64
}

var setHeroLevelCmd = &cobra.Command{
	Use:   "set_hero_level",
	Short: "set hero level",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)

		req := zb.SetHeroLevelRequest{
			UserId: setHeroLevelCmdArgs.userID,
			HeroId: setHeroLevelCmdArgs.heroID,
			Level:  setHeroLevelCmdArgs.level,
		}
		result := zb.SetHeroLevelResponse{}

		_, err := commonTxObjs.contract.Call("SetHeroLevel", &req, signer, &result)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			return printProtoMessageAsJSONToStdout(&result)
		default:
			fmt.Printf("hero_id: %d\n", result.HeroId)
			fmt.Printf("level: %d\n", result.Level)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setHeroLevelCmd)

	setHeroLevelCmd.Flags().StringVarP(&setHeroLevelCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	setHeroLevelCmd.Flags().Int64VarP(&setHeroLevelCmdArgs.heroID, "heroId", "i", 1, "heroID of hero")
	setHeroLevelCmd.Flags().Int64VarP(&setHeroLevelCmdArgs.level, "level", "l", 1, "level to be set")
}
