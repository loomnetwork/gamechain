package cmd

import (
	"time"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var registerPlayerPoolCmdArgs struct {
	userID     string
	deckID     int64
	version    string
	randomSeed int64
	tags       []string
}

var registerPlayerPoolCmd = &cobra.Command{
	Use:   "register_player_pool",
	Short: "register player to find_match pool",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var req = zb.RegisterPlayerPoolRequest{
			UserId: registerPlayerPoolCmdArgs.userID,
			DeckId: registerPlayerPoolCmdArgs.deckID,
			Tags:   registerPlayerPoolCmdArgs.tags,
		}
		var resp zb.RegisterPlayerPoolResponse

		req.UserId = registerPlayerPoolCmdArgs.userID
		req.Version = registerPlayerPoolCmdArgs.version
		req.RandomSeed = registerPlayerPoolCmdArgs.randomSeed

		_, err := commonTxObjs.contract.Call("RegisterPlayerPool", &req, signer, &resp)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(registerPlayerPoolCmd)

	registerPlayerPoolCmd.Flags().StringVarP(&registerPlayerPoolCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	registerPlayerPoolCmd.Flags().Int64VarP(&registerPlayerPoolCmdArgs.deckID, "deckId", "d", 1, "Deck Id")
	registerPlayerPoolCmd.Flags().StringVarP(&registerPlayerPoolCmdArgs.version, "version", "v", "", "version number like “0.10.0”")
	registerPlayerPoolCmd.Flags().Int64VarP(&registerPlayerPoolCmdArgs.randomSeed, "randomSeed", "s", time.Now().Unix(), "Random Seed")
	registerPlayerPoolCmd.Flags().StringArrayVarP(&registerPlayerPoolCmdArgs.tags, "tags", "t", nil, "tags")
}
