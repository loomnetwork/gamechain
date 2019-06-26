package cmd

import (
	"fmt"
	"github.com/loomnetwork/gamechain/tools/battleground_utility"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"strings"

	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
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

		var req zb_calls.GetGameModeRequest
		var gameMode = zb_data.GameMode{}

		req.ID = getGameModeCmdArgs.ID

		_, err := commonTxObjs.contract.StaticCall("GetGameMode", &req, callerAddr, &gameMode)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			err := battleground_utility.PrintProtoMessageAsJsonToStdout(&gameMode)
			if err != nil {
				return err
			}
		default:
			fmt.Printf("found game mode: %+v", gameMode)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getGameModeCmd)
	getGameModeCmd.Flags().StringVar(&getGameModeCmdArgs.ID, "id", "", "id of the game mode")
}
