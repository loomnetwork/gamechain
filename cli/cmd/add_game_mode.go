package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var addGameModeCmdArgs struct {
	name         string
	description  string
	version      string
	address      string
	gameModeType int
	oracle       string
}

var addGameModeCmd = &cobra.Command{
	Use:   "add_game_mode",
	Short: "add a game mode for zombiebattleground",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var req zb_calls.GameModeRequest

		req.Name = addGameModeCmdArgs.name
		req.Description = addGameModeCmdArgs.description
		req.Version = addGameModeCmdArgs.version
		req.Address = addGameModeCmdArgs.address
		req.GameModeType = zb_data.GameModeType(addGameModeCmdArgs.gameModeType)
		req.Oracle = addGameModeCmdArgs.oracle

		result := zb_data.GameMode{}

		_, err := commonTxObjs.contract.Call("AddGameMode", &req, signer, &result)
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
			fmt.Printf("added game mode: %+v", result)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addGameModeCmd)
	addGameModeCmd.Flags().StringVarP(&addGameModeCmdArgs.name, "name", "n", "", "name for the new game mode")
	addGameModeCmd.Flags().StringVarP(&addGameModeCmdArgs.description, "description", "d", "", "description")
	addGameModeCmd.Flags().StringVarP(&addGameModeCmdArgs.version, "version", "v", "", "version number like “0.10.0”")
	addGameModeCmd.Flags().StringVarP(&addGameModeCmdArgs.address, "address", "a", "", "address of game mode")
	addGameModeCmd.Flags().IntVarP(&addGameModeCmdArgs.gameModeType, "gameModeType", "t", 0, "type of game mode")
	addGameModeCmd.Flags().StringVarP(&addGameModeCmdArgs.oracle, "oracle", "o", "", "oracle address")
}
