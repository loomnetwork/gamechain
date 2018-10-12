package cmd

import (
	"fmt"
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var staticCallGameModeCustomGameModeFunctionArgs struct {
	ID       string
	Function string
}

var staticCallGameModeCustomGameModeFunctionCmd = &cobra.Command{
	Use:   "static_call_game_mode_custom_function",
	Short: "calls a custom function on a game mode",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		var req zb.GetGameModeRequest
		var gameMode = zb.GameMode{}

		req.ID = staticCallGameModeCustomGameModeFunctionArgs.ID

		_, err := commonTxObjs.contract.StaticCall("GetGameMode", &req, callerAddr, &gameMode)
		if err != nil {
			return err
		}

		var reqUi zb.CallCustomGameModeFunctionRequest

		reqUi.Address = gameMode.Address
		reqUi.FunctionName = staticCallGameModeCustomGameModeFunctionArgs.Function

		var resp zb.StaticCallCustomGameModeFunctionResponse
		_, err = commonTxObjs.contract.StaticCall("StaticCallCustomGameModeFunction", &reqUi, callerAddr, &resp)
		if err != nil {
			return err
		}

		fmt.Println(resp)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(staticCallGameModeCustomGameModeFunctionCmd)
	staticCallGameModeCustomGameModeFunctionCmd.Flags().StringVar(&staticCallGameModeCustomGameModeFunctionArgs.ID, "id", "", "id of the game mode")
	staticCallGameModeCustomGameModeFunctionCmd.Flags().StringVar(&staticCallGameModeCustomGameModeFunctionArgs.Function, "function", "", "function name to call")
}
