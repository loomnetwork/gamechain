package cmd

import (
	"fmt"

	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/gamechain/types/zb"
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
			fmt.Printf("ID: %s\n", gameMode.ID)
			fmt.Printf("Name: %s\n", gameMode.Name)
			fmt.Printf("Description: %s\n", gameMode.Description)
			fmt.Printf("Version: %s\n", gameMode.Version)
			fmt.Printf("GameModeType: %s\n", gameMode.GameModeType)
			fmt.Printf("Address: %s\n", gameMode.Address.String())
			fmt.Printf("Owner: %s\n", gameMode.Owner.String())
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listGameModesCmd)
}
