package cmd

import (
	"fmt"
	"strings"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var listHeroLibraryCmdArgs struct {
	version string
}

var listHeroLibraryCmd = &cobra.Command{
	Use:   "list_hero_library",
	Short: "list hero library",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := zb.ListHeroLibraryRequest{
			Version: listHeroLibraryCmdArgs.version,
		}
		result := zb.ListHeroLibraryResponse{}

		_, err := commonTxObjs.contract.StaticCall("ListHeroLibrary", &req, callerAddr, &result)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			return printProtoMessageAsJSONToStdout(&result)
		default:
			for _, heroInfo := range result.Heroes {
				fmt.Printf("hero_id: %d\n", heroInfo.HeroId)
				for _, skill := range heroInfo.Skills {
					fmt.Printf("skill title: %s\n", skill.Title)
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listHeroLibraryCmd)

	listHeroLibraryCmd.Flags().StringVarP(&listHeroLibraryCmdArgs.version, "version", "v", "v1", "Version")
}
