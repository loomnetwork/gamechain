package cmd

import (
	"fmt"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var addGameModeCmdArgs struct {
	name         string
	description  string
	version      string
	gameModeType int
}

var addGameModeCmd = &cobra.Command{
	Use:   "add_game_mode",
	Short: "add a game mode for zombiebattleground",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var req zb.GameModeRequest

		req.Name = addGameModeCmdArgs.name
		req.Description = addGameModeCmdArgs.description
		req.Version = addGameModeCmdArgs.version
		req.GameModeType = zb.GameModeType(addGameModeCmdArgs.gameModeType)

		_, err := commonTxObjs.contract.Call("AddGameMode", &req, signer, nil)
		if err != nil {
			return err
		}
		fmt.Printf("added game mode")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addGameModeCmd)
	addGameModeCmd.Flags().StringVarP(&addGameModeCmdArgs.name, "name", "n", "", "name for the new game mode")
	addGameModeCmd.Flags().StringVarP(&addGameModeCmdArgs.description, "description", "d", "", "description")
	addGameModeCmd.Flags().StringVarP(&addGameModeCmdArgs.version, "version", "v", "", "version number like “0.10.0”")
	addGameModeCmd.Flags().IntVarP(&addGameModeCmdArgs.gameModeType, "gameModeType", "t", 0, "type of game mode")
}
