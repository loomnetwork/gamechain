package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"strings"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var deleteGameModeCmdArgs struct {
	ID     string
	oracle string
}

var deleteGameModeCmd = &cobra.Command{
	Use:   "delete_game_mode",
	Short: "delete game mode by id",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var req zb_calls.DeleteGameModeRequest

		req.ID = deleteGameModeCmdArgs.ID
		req.Oracle = deleteGameModeCmdArgs.oracle

		_, err := commonTxObjs.contract.Call("DeleteGameMode", &req, signer, nil)
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
			fmt.Printf("deleted game mode: %s", req.ID)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteGameModeCmd)
	deleteGameModeCmd.Flags().StringVar(&deleteGameModeCmdArgs.ID, "id", "", "id of the game mode")
	deleteGameModeCmd.Flags().StringVarP(&deleteGameModeCmdArgs.oracle, "oracle", "o", "", "oracle address")
}
