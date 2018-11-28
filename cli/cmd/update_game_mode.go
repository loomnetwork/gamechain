package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var updateGameModeCmdArgs struct {
	ID           string
	name         string
	description  string
	version      string
	gameModeType int
	oracle       string
}

var updateGameModeCmd = &cobra.Command{
	Use:   "update_game_mode",
	Short: "update a game mode for zombiebattleground",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var req zb.UpdateGameModeRequest

		req.ID = updateGameModeCmdArgs.ID
		req.Name = updateGameModeCmdArgs.name
		req.Description = updateGameModeCmdArgs.description
		req.Version = updateGameModeCmdArgs.version
		req.GameModeType = zb.GameModeType(updateGameModeCmdArgs.gameModeType)
		req.Oracle = updateGameModeCmdArgs.oracle

		_, err := commonTxObjs.contract.Call("UpdateGameMode", &req, signer, nil)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			output, err := json.Marshal(map[string]interface{}{"success": true})
			if err != nil {
				return err
			}
			fmt.Println(string(output))
		default:
			fmt.Printf("updated game mode")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateGameModeCmd)
	updateGameModeCmd.Flags().StringVar(&updateGameModeCmdArgs.ID, "id", "", "id")
	updateGameModeCmd.Flags().StringVarP(&updateGameModeCmdArgs.name, "name", "n", "", "name for the game mode")
	updateGameModeCmd.Flags().StringVarP(&updateGameModeCmdArgs.description, "description", "d", "", "description")
	updateGameModeCmd.Flags().StringVarP(&updateGameModeCmdArgs.version, "version", "v", "", "version number like “0.10.0”")
	updateGameModeCmd.Flags().IntVarP(&updateGameModeCmdArgs.gameModeType, "gameModeType", "t", 0, "type of game mode")
	updateGameModeCmd.Flags().StringVarP(&updateGameModeCmdArgs.oracle, "oracle", "o", "", "oracle address")
}
