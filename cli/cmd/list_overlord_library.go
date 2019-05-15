package cmd

import (
	"fmt"
	"strings"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var listOverlordsLibraryCmdArgs struct {
	version string
}

var listOverlordsLibraryCmd = &cobra.Command{
	Use:   "list_overlord_library",
	Short: "list overlord library",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		req := zb_calls.ListOverlordLibraryRequest{
			Version: listOverlordsLibraryCmdArgs.version,
		}
		result := zb_calls.ListOverlordLibraryResponse{}

		_, err := commonTxObjs.contract.StaticCall("ListOverlordLibrary", &req, callerAddr, &result)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			return printProtoMessageAsJSONToStdout(&result)
		default:
			for _, overlordInfo := range result.Overlords {
				fmt.Printf("overlord_id: %d\n", overlordInfo.OverlordId)
				for _, skill := range overlordInfo.Skills {
					fmt.Printf("skill title: %s\n", skill.Title)
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listOverlordsLibraryCmd)

	listOverlordsLibraryCmd.Flags().StringVarP(&listOverlordsLibraryCmdArgs.version, "version", "v", "v1", "Version")

	_ = listOverlordsLibraryCmd.MarkFlagRequired("version")
}
