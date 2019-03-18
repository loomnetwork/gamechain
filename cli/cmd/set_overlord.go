package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var setOverlordCmdArgs struct {
	userID   string
	overlordID   int64
	filename string
}

var setOverlordCmd = &cobra.Command{
	Use:   "set_overlord",
	Short: "set overlord",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)

		if setOverlordCmdArgs.filename == "" {
			return fmt.Errorf("filename not provided")
		}
		f, err := os.Open(setOverlordCmdArgs.filename)
		if err != nil {
			return fmt.Errorf("error reading file: %s", err.Error())
		}
		defer f.Close()

		var overlord zb.Overlord
		if err := new(jsonpb.Unmarshaler).Unmarshal(f, &overlord); err != nil {
			return fmt.Errorf("error parsing JSON file: %s", err.Error())
		}
		req := zb.SetOverlordRequest{
			UserId: setOverlordCmdArgs.userID,
			OverlordId: setOverlordCmdArgs.overlordID,
			Overlord:   &overlord,
		}

		result := zb.SetOverlordResponse{}

		_, err = commonTxObjs.contract.Call("SetOverlord", &req, signer, &result)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			return printProtoMessageAsJSONToStdout(&result)
		default:
			if result.Overlord != nil {
				fmt.Printf("overlord_id: %d\n", result.Overlord.OverlordId)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setOverlordCmd)

	setOverlordCmd.Flags().StringVarP(&setOverlordCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	setOverlordCmd.Flags().Int64VarP(&setOverlordCmdArgs.overlordID, "overlordId", "i", 1, "overlordID of overlord")
	setOverlordCmd.Flags().StringVarP(&setOverlordCmdArgs.filename, "filename", "f", "overlord.json", "Overlord file name in JSON format")
}
