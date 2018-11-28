package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var setAIDecksCmdArgs struct {
	file    string
	version string
}

var setAIDecksCmd = &cobra.Command{
	Use:   "set_ai_decks",
	Short: "set AI decks",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var deckList zb.AIDeckList

		if setAIDecksCmdArgs.file == "" {
			return fmt.Errorf("filename not provided")
		}

		f, err := os.Open(setAIDecksCmdArgs.file)
		if err != nil {
			return fmt.Errorf("error reading file: %s", err.Error())
		}
		defer f.Close()

		if err := new(jsonpb.Unmarshaler).Unmarshal(f, &deckList); err != nil {
			return fmt.Errorf("error parsing JSON file: %s", err.Error())
		}
		req := &zb.SetAIDecksRequest{
			Decks:   deckList.Decks,
			Version: setAIDecksCmdArgs.version,
		}

		_, err = commonTxObjs.contract.Call("SetAIDecks", req, signer, nil)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			output, err := json.Marshal(map[string]interface{}{"success": true})
			if err != nil {
				return err
			}
			fmt.Println(string(output))
		default:
			fmt.Printf("decks set successfully")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setAIDecksCmd)

	setAIDecksCmd.Flags().StringVarP(&setAIDecksCmdArgs.file, "file", "f", "", "json file containing decks data")
	setAIDecksCmd.Flags().StringVarP(&setAIDecksCmdArgs.version, "version", "v", "v1", "Version")
}
