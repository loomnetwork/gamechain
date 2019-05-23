package cmd

import (
	"fmt"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var findMatchCmdArgs struct {
	userID string
	tags   []string
}

var findMatchCmd = &cobra.Command{
	Use:   "find_match",
	Short: "find match for zombiebattleground",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var req = zb_calls.FindMatchRequest{
			UserId: findMatchCmdArgs.userID,
			Tags:   findMatchCmdArgs.tags,
		}
		var resp zb_calls.FindMatchResponse

		req.UserId = findMatchCmdArgs.userID

		_, err := commonTxObjs.contract.Call("FindMatch", &req, signer, &resp)
		if err != nil {
			return err
		}
		if resp.Match != nil {
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
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(findMatchCmd)

	findMatchCmd.Flags().StringVarP(&findMatchCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	findMatchCmd.Flags().StringArrayVarP(&findMatchCmdArgs.tags, "tags", "t", nil, "tags")
}
