package cmd

import (
	"fmt"

	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/spf13/cobra"
)

var listCardCmd = &cobra.Command{
	Use:   "listcard",
	Short: "list card",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := zb.ListCardLibraryRequest{}
		resp := zb.ListCardLibraryResponse{}

		_, err := commonTxObjs.contract.StaticCall("ListCardLibrary", &req, callerAddr, &resp)
		if err != nil {
			return err
		}
		fmt.Println(resp)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCardCmd)
}
