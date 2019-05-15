package cmd

import (
	"fmt"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var acceptMatchCmdArgs struct {
	userID  string
	matchID int64
}

var acceptMatchCmd = &cobra.Command{
	Use:   "accept_match",
	Short: "accept match",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var req = zb_calls.AcceptMatchRequest{
			UserId:  acceptMatchCmdArgs.userID,
			MatchId: acceptMatchCmdArgs.matchID,
		}
		var resp zb_calls.AcceptMatchResponse

		_, err := commonTxObjs.contract.Call("AcceptMatch", &req, signer, &resp)
		if err != nil {
			return err
		}
		match := resp.Match
		fmt.Printf("MatchID: %d\n", match.Id)
		fmt.Printf("Status: %s\n", match.Status)
		fmt.Printf("Topic: %v\n", match.Topics)
		fmt.Printf("Players:\n")
		for _, player := range match.PlayerStates {
			fmt.Printf("\tPlayerID: %s\n", player.Id)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(acceptMatchCmd)

	acceptMatchCmd.Flags().StringVarP(&acceptMatchCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	acceptMatchCmd.Flags().Int64VarP(&acceptMatchCmdArgs.matchID, "matchId", "m", 0, "matchId")
}
