package cmd

import (
	"fmt"

	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/spf13/cobra"
)

var listCardCmdArgs struct {
	version string
}

var listCardCmd = &cobra.Command{
	Use:   "list_card",
	Short: "list card",
	RunE: func(cmd *cobra.Command, args []string) error {

		if listCardCmdArgs.version == "" {
			return fmt.Errorf("version not specified")
		}

		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := zb.ListCardLibraryRequest{}
		result := zb.ListCardLibraryResponse{}

		req.Version = listCardCmdArgs.version

		_, err := commonTxObjs.contract.StaticCall("ListCardLibrary", &req, callerAddr, &result)
		if err != nil {
			return err
		}
		fmt.Printf("card library size: %d\n", len(result.Sets))
		for _, set := range result.Sets {
			for _, card := range set.Cards {
				fmt.Printf("card mould ID: %d, name: %s\n", card.MouldId, card.Name)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCardCmd)
	listCardCmd.Flags().StringVarP(&listCardCmdArgs.version, "version", "v", "", "Version")
}
