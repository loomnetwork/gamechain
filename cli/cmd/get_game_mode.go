package cmd

import (
	"fmt"

	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var getGameModeCmdArgs struct {
	ID string
}

var getGameModeCmd = &cobra.Command{
	Use:   "get_game_mode",
	Short: "get game mode by id",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		var req zb.GetGameModeRequest
		var gameMode = zb.GameMode{}

		req.ID = getGameModeCmdArgs.ID

		_, err := commonTxObjs.contract.StaticCall("GetGameMode", &req, callerAddr, &gameMode)
		if err != nil {
			return err
		}
		fmt.Printf("found game mode: %+v", gameMode)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getGameModeCmd)
	getGameModeCmd.Flags().StringVar(&getGameModeCmdArgs.ID, "id", "", "id of the game mode")
}
