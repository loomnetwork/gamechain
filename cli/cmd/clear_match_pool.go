package cmd

import (
	"fmt"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var clearMatchPoolCmdArgs struct {
}

var clearMatchPoolCmd = &cobra.Command{
	Use:   "clear_match_pool",
	Short: "clears out the matchmaking pool from the contract",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)

		req := &zb.ClearMatchPoolRequest{}

		_, err := commonTxObjs.contract.Call("ClearMatchPool", req, signer, nil)
		if err != nil {
			return err
		}
		fmt.Printf("deck deleted successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(clearMatchPoolCmd)
}
