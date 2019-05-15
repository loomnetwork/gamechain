package cmd

import (
	"fmt"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom"
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
			output, err := new(jsonpb.Marshaler).MarshalToString(&gameMode)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
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
