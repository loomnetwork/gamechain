package cmd

import (
	"fmt"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var deleteGameModeCmdArgs struct {
	name string
}

var deleteGameModeCmd = &cobra.Command{
	Use:   "delete_game_mode",
	Short: "delete game mode by id",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var req zb.DeleteGameModeRequest

		req.Name = deleteGameModeCmdArgs.name

		_, err := commonTxObjs.contract.Call("DeleteGameMode", &req, signer, nil)
		if err != nil {
			return err
		}
		fmt.Printf("deleted game mode: %s", req.Name)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteGameModeCmd)
	deleteGameModeCmd.Flags().StringVarP(&deleteGameModeCmdArgs.name, "name", "n", "", "name of the game mode")
}
