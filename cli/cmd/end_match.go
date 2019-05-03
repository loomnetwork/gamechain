package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var endMatchCmdArgs struct {
	userID                 string
	matchID                int64
	winnerID               string
	playerMatchExperiences *[]int64
}

var endMatchCmd = &cobra.Command{
	Use:   "end_match",
	Short: "end match for zombiebattleground",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(*endMatchCmdArgs.playerMatchExperiences) != 2 {
			return errors.New("'playerMatchExperience' length must be 2")
		}

		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var req = zb.EndMatchRequest{
			UserId:           endMatchCmdArgs.userID,
			MatchId:          endMatchCmdArgs.matchID,
			WinnerId:         endMatchCmdArgs.winnerID,
			MatchExperiences: *endMatchCmdArgs.playerMatchExperiences,
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
	endMatchCmdArgs.playerMatchExperiences = &[]int64{0, 0}

	rootCmd.AddCommand(endMatchCmd)
	endMatchCmd.Flags().StringVarP(&endMatchCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	endMatchCmd.Flags().Int64VarP(&endMatchCmdArgs.matchID, "matchId", "m", 0, "Match ID")
	endMatchCmd.Flags().StringVar(&endMatchCmdArgs.winnerID, "winnerId", "loom", "Winner ID")
	endMatchCmd.Flags().Int64SliceVarP(endMatchCmdArgs.playerMatchExperiences, "playerMatchExperience", "p", []int64{0, 0},  "Players Match Experiences")
}
