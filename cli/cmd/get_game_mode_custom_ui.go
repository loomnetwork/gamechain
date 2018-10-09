package cmd

import (
	"fmt"
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var getGameModeCustomUiCmdArgs struct {
	ID string
}

var getGameModeCustomUiCmd = &cobra.Command{
	Use:   "get_game_mode_custom_ui",
	Short: "get game mode custom ui by id",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		var req zb.GetGameModeRequest
		var gameMode = zb.GameMode{}

		req.ID = getGameModeCustomUiCmdArgs.ID

		_, err := commonTxObjs.contract.StaticCall("GetGameMode", &req, callerAddr, &gameMode)
		if err != nil {
			return err
		}

		var reqUi zb.GetCustomGameModeCustomUiRequest

		reqUi.Address = gameMode.Address

		result := zb.GetCustomGameModeCustomUiResponse{}
		_, err = commonTxObjs.contract.StaticCall("GetGameModeCustomUi", &reqUi, callerAddr, &result)
		if err != nil {
			return err
		}

		fmt.Println(result.UiElements)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getGameModeCustomUiCmd)
	getGameModeCustomUiCmd.Flags().StringVar(&getGameModeCustomUiCmdArgs.ID, "id", "", "id of the game mode")
}
