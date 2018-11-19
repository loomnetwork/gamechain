package cmd

import (
	"fmt"

	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var getMatchCmdArgs struct {
	MatchID int64
}

var getMatchCmd = &cobra.Command{
	Use:   "get_match",
	Short: "get match",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}
		var req = zb.GetMatchRequest{
			MatchId: getMatchCmdArgs.MatchID,
		}
		var resp zb.GetMatchResponse

		_, err := commonTxObjs.contract.StaticCall("GetMatch", &req, callerAddr, &resp)
		if err != nil {
			return err
		}
		match := resp.Match
		fmt.Printf("MatchID: %d\n", match.Id)
		fmt.Printf("Status: %s\n", match.Status)
		fmt.Printf("Topic: %v\n", match.Topics)
		fmt.Printf("Version: %s\n", match.Version)
		fmt.Printf("Players:\n")
		for i, player := range match.PlayerStates {
			fmt.Printf("\tPlayer%d: %s - accepted: %v\n", i+1, player.Id, player.MatchAccepted)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getMatchCmd)

	getMatchCmd.Flags().Int64VarP(&getMatchCmdArgs.MatchID, "matchId", "m", 0, "Match ID")
}
