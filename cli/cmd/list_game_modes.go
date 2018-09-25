package cmd

import (
	"fmt"

	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var listGameModesCmd = &cobra.Command{
	Use:   "list_game_modes",
	Short: "list game modes",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := &zb.ListGameModesRequest{}
		var result zb.GameModeList
		_, err := commonTxObjs.contract.StaticCall("ListGameModes", req, callerAddr, &result)
		if err != nil {
			return err
		}

		for _, gameMode := range result.GameModes {
			fmt.Printf("name: %s\n", gameMode.Name)
			fmt.Printf("name: %s\n", gameMode.Description)
			fmt.Printf("name: %s\n", gameMode.Version)
			fmt.Printf("name: %s\n", gameMode.GameModeType)
			fmt.Printf("name: %s\n", gameMode.Address)
			fmt.Printf("name: %s\n", gameMode.Owner)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listGameModesCmd)
}
