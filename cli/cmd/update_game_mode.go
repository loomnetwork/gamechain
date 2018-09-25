package cmd

import (
	"fmt"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var updateGameModeCmdArgs struct {
	name         string
	description  string
	version      string
	gameModeType int
}

var updateGameModeCmd = &cobra.Command{
	Use:   "update_game_mode",
	Short: "update a game mode for zombiebattleground",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var req zb.GameModeRequest

		req.Name = updateGameModeCmdArgs.name
		req.Description = updateGameModeCmdArgs.description
		req.Version = updateGameModeCmdArgs.version
		req.GameModeType = zb.GameModeType(updateGameModeCmdArgs.gameModeType)

		_, err := commonTxObjs.contract.Call("UpdateGameMode", &req, signer, nil)
		if err != nil {
			return err
		}
		fmt.Printf("updated game mode")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateGameModeCmd)
	updateGameModeCmd.Flags().StringVarP(&updateGameModeCmdArgs.name, "name", "n", "", "name for the game mode")
	updateGameModeCmd.Flags().StringVarP(&updateGameModeCmdArgs.description, "description", "d", "", "description")
	updateGameModeCmd.Flags().StringVarP(&updateGameModeCmdArgs.version, "version", "v", "", "version number like “0.10.0”")
	updateGameModeCmd.Flags().IntVarP(&updateGameModeCmdArgs.gameModeType, "gameModeType", "t", 0, "type of game mode")
}
