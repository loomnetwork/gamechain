package cmd

import (
	"fmt"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var updateVersionsCmdArgs struct {
	contentVersion string
	pvpVersion     string
}

var updateVersionsCmd = &cobra.Command{
	Use:   "update_versions",
	Short: "update the content and pvp versions",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		req := zb.UpdateVersionsRequest{
			ContentVersion: updateVersionsCmdArgs.contentVersion,
			PvpVersion:     updateVersionsCmdArgs.pvpVersion,
		}

		_, err := commonTxObjs.contract.Call("UpdateVersions", &req, signer, nil)
		if err != nil {
			return fmt.Errorf("error encountered while calling UpdateVersions: %s", err.Error())
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateVersionsCmd)

	updateVersionsCmd.Flags().StringVarP(&updateVersionsCmdArgs.contentVersion, "content", "c", "", "content version")
	updateVersionsCmd.Flags().StringVarP(&updateVersionsCmdArgs.pvpVersion, "pvp", "p", "", "pvp version")
}
