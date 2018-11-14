package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var findMatchCmdArgs struct {
	userID     string
	deckID     int64
	version    string
	randomSeed int64
	tags       []string
}

var findMatchCmd = &cobra.Command{
	Use:   "find_match",
	Short: "find match for zombiebattleground",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var req = zb.FindMatchRequest{
			UserId: findMatchCmdArgs.userID,
			DeckId: findMatchCmdArgs.deckID,
			Tags:   findMatchCmdArgs.tags,
		}
		var resp zb.FindMatchResponse

		req.UserId = findMatchCmdArgs.userID
		req.Version = findMatchCmdArgs.version
		req.RandomSeed = findMatchCmdArgs.randomSeed

		_, err := commonTxObjs.contract.Call("FindMatch", &req, signer, &resp)
		if err != nil {
			return err
		}
		match := resp.Match

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			output, err := new(jsonpb.Marshaler).MarshalToString(match)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
		default:
			fmt.Printf("MatchID: %d\n", match.Id)
			fmt.Printf("Status: %s\n", match.Status)
			fmt.Printf("Topic: %v\n", match.Topics)
			fmt.Printf("Players:\n")
			for _, player := range match.PlayerStates {
				fmt.Printf("\tPlayerID: %s\n", player.Id)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(findMatchCmd)

	findMatchCmd.Flags().StringVarP(&findMatchCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	findMatchCmd.Flags().Int64VarP(&findMatchCmdArgs.deckID, "deckId", "d", 1, "Deck Id")
	findMatchCmd.Flags().StringVarP(&findMatchCmdArgs.version, "version", "v", "", "version number like “0.10.0”")
	findMatchCmd.Flags().Int64VarP(&findMatchCmdArgs.randomSeed, "randomSeed", "s", time.Now().Unix(), "Random Seed")
	findMatchCmd.Flags().StringArrayVarP(&findMatchCmdArgs.tags, "tags", "t", nil, "tags")
}
