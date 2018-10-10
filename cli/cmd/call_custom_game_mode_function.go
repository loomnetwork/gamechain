package cmd

import (
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var callGameModeCustomGameModeFunctionArgs struct {
	ID       string
	Function string
}

var allGameModeCustomGameModeFunctionCmd = &cobra.Command{
	Use:   "call_game_mode_custom_function",
	Short: "calls a custom function on a game mode",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		var req zb.GetGameModeRequest
		var gameMode = zb.GameMode{}

		req.ID = callGameModeCustomGameModeFunctionArgs.ID

		_, err := commonTxObjs.contract.StaticCall("GetGameMode", &req, callerAddr, &gameMode)
		if err != nil {
			return err
		}

		var reqUi zb.CallCustomGameModeFunctionRequest

		reqUi.Address = gameMode.Address
		reqUi.FunctionName = callGameModeCustomGameModeFunctionArgs.Function

		_, err = commonTxObjs.contract.Call("CallCustomGameModeFunction", &reqUi, signer, nil)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(allGameModeCustomGameModeFunctionCmd)
	allGameModeCustomGameModeFunctionCmd.Flags().StringVar(&callGameModeCustomGameModeFunctionArgs.ID, "id", "", "id of the game mode")
	allGameModeCustomGameModeFunctionCmd.Flags().StringVar(&callGameModeCustomGameModeFunctionArgs.Function, "function", "", "function name to call")
}
