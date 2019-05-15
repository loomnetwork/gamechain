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
	tags                []string
	useBackendGameLogic bool

	enableDebugCheats    bool
	disableDeckShuffle   bool
	forceFirstTurnUserId string
	randomSeed           int64
}

var registerPlayerPoolCmd = &cobra.Command{
	Use:   "register_player_pool",
	Short: "register player to find_match pool",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var req = zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				UserId:              registerPlayerPoolCmdArgs.userID,
				DeckId:              registerPlayerPoolCmdArgs.deckID,
				Version:             registerPlayerPoolCmdArgs.version,
				Tags:                registerPlayerPoolCmdArgs.tags,
				UseBackendGameLogic: registerPlayerPoolCmdArgs.useBackendGameLogic,
			},
		}
		var resp zb_calls.RegisterPlayerPoolResponse

		if registerPlayerPoolCmdArgs.enableDebugCheats {
			req.RegistrationData.DebugCheats.Enabled = true
		}

		req.RegistrationData.DebugCheats.DisableDeckShuffle = registerPlayerPoolCmdArgs.disableDeckShuffle
		req.RegistrationData.DebugCheats.ForceFirstTurnUserId = registerPlayerPoolCmdArgs.forceFirstTurnUserId

		if registerPlayerPoolCmdArgs.randomSeed != 0 {
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
	registerPlayerPoolCmd.Flags().BoolVarP(&registerPlayerPoolCmdArgs.enableDebugCheats, "enableDebugCheats", "", false, "Enable Debug Cheats")
	registerPlayerPoolCmd.Flags().BoolVarP(&registerPlayerPoolCmdArgs.disableDeckShuffle, "disableDeckShuffle", "", false, "Cheat - Disable Deck Shuffle")
	registerPlayerPoolCmd.Flags().StringVarP(&registerPlayerPoolCmdArgs.forceFirstTurnUserId, "forceFirstTurnUserId", "", "", "Cheat - UserId of the player who will have the first turn")
	registerPlayerPoolCmd.Flags().Int64VarP(&registerPlayerPoolCmdArgs.randomSeed, "randomSeed", "", 0, "Cheat - Random Seed")
	registerPlayerPoolCmd.Flags().StringArrayVarP(&registerPlayerPoolCmdArgs.tags, "tags", "t", nil, "tags")
	registerPlayerPoolCmd.Flags().BoolVarP(&registerPlayerPoolCmdArgs.useBackendGameLogic, "useBackendGameLogic", "b", false, "useBackendGameLogic")
}
