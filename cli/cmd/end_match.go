package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var endMatchCmdArgs struct {
	userID                  string
	matchID                 int64
	winnerID                string
	playerMatchExperience   int64
	opponentMatchExperience int64
}

var endMatchCmd = &cobra.Command{
	Use:   "end_match",
	Short: "end match for zombiebattleground",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var req = zb.EndMatchRequest{
			UserId:                  endMatchCmdArgs.userID,
			MatchId:                 endMatchCmdArgs.matchID,
			WinnerId:                endMatchCmdArgs.winnerID,
			PlayerMatchExperience:   endMatchCmdArgs.playerMatchExperience,
			OpponentMatchExperience: endMatchCmdArgs.opponentMatchExperience,
		}
		var resp zb.EndMatchResponse

		_, err := commonTxObjs.contract.Call("EndMatch", &req, signer, &resp)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			output, err := json.Marshal(map[string]interface{}{"success": true})
			if err != nil {
				return err
			}
			fmt.Println(string(output))
		default:
			fmt.Printf("end match %v successfully", req.MatchId)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(endMatchCmd)
	endMatchCmd.Flags().StringVarP(&endMatchCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	endMatchCmd.Flags().Int64VarP(&endMatchCmdArgs.matchID, "matchId", "m", 0, "Match ID")
	endMatchCmd.Flags().StringVar(&endMatchCmdArgs.winnerID, "winnerId", "loom", "Winner ID")
	endMatchCmd.Flags().Int64VarP(&endMatchCmdArgs.playerMatchExperience, "playerMatchExperience", "p", 0, "Player Match Experience")
	endMatchCmd.Flags().Int64VarP(&endMatchCmdArgs.opponentMatchExperience, "opponentMatchExperience", "o", 0, "Opponent Match Experience")
}
