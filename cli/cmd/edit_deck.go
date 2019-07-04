package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"strings"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var editDeckCmdArgs struct {
	userID  string
	data    string
	version string
}

var editDeckCmd = &cobra.Command{
	Use:   "edit_deck",
	Short: "edit deck in zombie battleground",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var deck zb_data.Deck

		if err := json.Unmarshal([]byte(editDeckCmdArgs.data), &deck); err != nil {
			return fmt.Errorf("invalid JSON passed in data field. Error: %s", err.Error())
		}

		req := &zb_calls.EditDeckRequest{
			Deck:    &deck,
			UserId:  editDeckCmdArgs.userID,
			Version: editDeckCmdArgs.version,
		}

		_, err := commonTxObjs.contract.Call("EditDeck", req, signer, nil)
		if err != nil {
			return fmt.Errorf("error encountered while calling EditDeck: %s", err.Error())
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			output, err := json.Marshal(map[string]interface{}{"success": true})
			if err != nil {
				return err
			}
			fmt.Println(string(output))
		default:
			fmt.Printf("deck edited successfully\n")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(editDeckCmd)

	editDeckCmd.Flags().StringVarP(&editDeckCmdArgs.userID, "userId", "u", "loom", "UserId of account")
	editDeckCmd.Flags().StringVarP(&editDeckCmdArgs.data, "data", "d", "{\"id\": 0, \"overlord_id\":1, \"name\": \"NewDefaultDeck\", \"cards\": [ {\"cardKey\": {\"mouldId\": 4}, \"amount\": 2}, {\"cardKey\": {\"mouldId\": 43}, \"amount\": 2} ]}", "Deck data in serialized json format")
	editDeckCmd.Flags().StringVarP(&editDeckCmdArgs.version, "version", "v", "v1", "Version")

	_ = editDeckCmd.MarkFlagRequired("version")
}
