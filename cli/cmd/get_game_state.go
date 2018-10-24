package cmd

import (
	"github.com/loomnetwork/gamechain/battleground"
	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var getGameStateCmdArgs struct {
	MatchID int64
}

var getGameStateCmd = &cobra.Command{
	Use:   "get_game_state",
	Short: "get gamestate",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}
		var req = zb.GetGameStateRequest{
			MatchId: getGameStateCmdArgs.MatchID,
		}
		var resp zb.GetGameStateResponse
		_, err := commonTxObjs.contract.StaticCall("GetGameState", &req, callerAddr, &resp)
		if err != nil {
			return err
		}

		gp := battleground.Gameplay{State: resp.GameState}
		gp.PrintState()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getGameStateCmd)

	getGameStateCmd.Flags().Int64VarP(&getGameStateCmdArgs.MatchID, "matchId", "m", 0, "Match ID")
}
