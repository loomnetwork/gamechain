package cmd

import (
	"fmt"
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
	"io/ioutil"
)

var callGameModeCustomGameModeFunctionArgs struct {
	ID       string
	abiInputFile string
}

var callGameModeCustomGameModeFunctionCmd = &cobra.Command{
	Use:   "call_game_mode_custom_function",
	Short: "calls a custom function on a game mode",
	RunE: func(cmd *cobra.Command, args []string) error {
		abiInputFileContents, err := ioutil.ReadFile(callGameModeCustomGameModeFunctionArgs.abiInputFile)
		if err != nil {
			return fmt.Errorf("unable to ABI-encoded call data from file: %s",
				callGameModeCustomGameModeFunctionArgs.abiInputFile)
		}

		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		var req zb.GetGameModeRequest
		var gameMode = zb.GameMode{}

		req.ID = callGameModeCustomGameModeFunctionArgs.ID

		_, err = commonTxObjs.contract.StaticCall("GetGameMode", &req, callerAddr, &gameMode)
		if err != nil {
			return err
		}

		var reqUi zb.CallCustomGameModeFunctionRequest

		reqUi.Address = gameMode.Address
		reqUi.CallData = abiInputFileContents

		_, err = commonTxObjs.contract.Call("CallCustomGameModeFunction", &reqUi, signer, nil)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(callGameModeCustomGameModeFunctionCmd)
	callGameModeCustomGameModeFunctionCmd.Flags().StringVar(&callGameModeCustomGameModeFunctionArgs.ID, "id", "", "id of the game mode")
	callGameModeCustomGameModeFunctionCmd.Flags().StringVar(&callGameModeCustomGameModeFunctionArgs.abiInputFile, "abiInputFile", "call.bin", "Binary ABI-encoded function call data file")
}
