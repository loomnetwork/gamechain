package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var setHeroCmdArgs struct {
	userID   string
	heroID   int64
	filename string
}

var setHeroCmd = &cobra.Command{
	Use:   "set_hero",
	Short: "set hero",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)

		if setHeroCmdArgs.filename == "" {
			return fmt.Errorf("filename not provided")
		}
		f, err := os.Open(setHeroCmdArgs.filename)
		if err != nil {
			return fmt.Errorf("error reading file: %s", err.Error())
		}
		defer f.Close()

		var hero zb.Hero
		if err := new(jsonpb.Unmarshaler).Unmarshal(f, &hero); err != nil {
			return fmt.Errorf("error parsing JSON file: %s", err.Error())
		}
		req := zb.SetHeroRequest{
			UserId: setHeroCmdArgs.userID,
			HeroId: setHeroCmdArgs.heroID,
			Hero:   &hero,
		}

		result := zb.SetHeroResponse{}

		_, err = commonTxObjs.contract.Call("SetHero", &req, signer, &result)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			return printProtoMessageAsJSONToStdout(&result)
		default:
			if result.Hero != nil {
				fmt.Printf("hero_id: %d\n", result.Hero.HeroId)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setHeroCmd)

	setHeroCmd.Flags().StringVarP(&setHeroCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	setHeroCmd.Flags().Int64VarP(&setHeroCmdArgs.heroID, "heroId", "i", 1, "heroID of hero")
	setHeroCmd.Flags().StringVarP(&setHeroCmdArgs.filename, "filename", "f", "hero.json", "Hero file name in JSON format")
}
