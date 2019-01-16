package cmd

import (
	"fmt"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var registerPlayerPoolCmdArgs struct {
	userID              string
	deckID              int64
	version             string
	randomSeed          int64
	tags                []string
	useBackendGameLogic bool
}

var registerPlayerPoolCmd = &cobra.Command{
	Use:   "register_player_pool",
	Short: "register player to find_match pool",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var req = zb.RegisterPlayerPoolRequest{
			RegistrationData: &zb.PlayerProfileRegistrationData{
				UserId:              registerPlayerPoolCmdArgs.userID,
				DeckId:              registerPlayerPoolCmdArgs.deckID,
				Version:             registerPlayerPoolCmdArgs.version,
				Tags:                registerPlayerPoolCmdArgs.tags,
				UseBackendGameLogic: registerPlayerPoolCmdArgs.useBackendGameLogic,
			},
		}
		var resp zb.RegisterPlayerPoolResponse

		if registerPlayerPoolCmdArgs.randomSeed != 0 {
			req.RegistrationData.DebugCheats.Enabled = true
			req.RegistrationData.DebugCheats.UseCustomRandomSeed = true
			req.RegistrationData.DebugCheats.CustomRandomSeed = registerPlayerPoolCmdArgs.randomSeed
		}

		_, err := commonTxObjs.contract.Call("RegisterPlayerPool", &req, signer, &resp)
		if err != nil {
			return err
		}

		fmt.Printf("Registered player %s to pool", req.RegistrationData.UserId)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(registerPlayerPoolCmd)

	registerPlayerPoolCmd.Flags().StringVarP(&registerPlayerPoolCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	registerPlayerPoolCmd.Flags().Int64VarP(&registerPlayerPoolCmdArgs.deckID, "deckId", "d", 1, "Deck Id")
	registerPlayerPoolCmd.Flags().StringVarP(&registerPlayerPoolCmdArgs.version, "version", "v", "", "version number like “0.10.0”")
	registerPlayerPoolCmd.Flags().Int64VarP(&registerPlayerPoolCmdArgs.randomSeed, "randomSeed", "s", 0, "Random Seed")
	registerPlayerPoolCmd.Flags().StringArrayVarP(&registerPlayerPoolCmdArgs.tags, "tags", "t", nil, "tags")
	registerPlayerPoolCmd.Flags().BoolVarP(&registerPlayerPoolCmdArgs.useBackendGameLogic, "useBackendGameLogic", "b", false, "useBackendGameLogic")
}
