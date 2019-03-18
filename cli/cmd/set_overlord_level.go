package cmd

import (
	"fmt"
	"strings"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var setOverlordLevelCmdArgs struct {
	userID string
	overlordID int64
	level  int64
}

var setOverlordLevelCmd = &cobra.Command{
	Use:   "set_overlord_level",
	Short: "set overlord level",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)

		req := zb.SetOverlordLevelRequest{
			UserId: setOverlordLevelCmdArgs.userID,
			OverlordId: setOverlordLevelCmdArgs.overlordID,
			Level:  setOverlordLevelCmdArgs.level,
		}
		result := zb.SetOverlordLevelResponse{}

		_, err := commonTxObjs.contract.Call("SetOverlordLevel", &req, signer, &result)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			return printProtoMessageAsJSONToStdout(&result)
		default:
			fmt.Printf("overlord_id: %d\n", result.OverlordId)
			fmt.Printf("level: %d\n", result.Level)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setOverlordLevelCmd)

	setOverlordLevelCmd.Flags().StringVarP(&setOverlordLevelCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	setOverlordLevelCmd.Flags().Int64VarP(&setOverlordLevelCmdArgs.overlordID, "overlordId", "i", 1, "overlordID of overlord")
	setOverlordLevelCmd.Flags().Int64VarP(&setOverlordLevelCmdArgs.level, "level", "l", 1, "level to be set")
}
