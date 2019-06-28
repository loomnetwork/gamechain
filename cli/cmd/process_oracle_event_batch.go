package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/tools/battleground_utility"
	"github.com/loomnetwork/gamechain/types/oracle"
	"io/ioutil"
	"strings"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var processOracleEventBatchCmdArgs struct {
	processOracleEventBatchJsonFile string
	processOracleEventBatchProtobufFile string
}

var processOracleEventBatchCmd = &cobra.Command{
	Use:   "process_oracle_event_batch",
	Short: "calls ProcessOracleEventBatch method of the contract with a request from a file",
	RunE: func(cmd *cobra.Command, args []string) error {
		var req oracle.ProcessOracleEventBatchRequest

		var err error
		if processOracleEventBatchCmdArgs.processOracleEventBatchJsonFile != "" {
			err = battleground_utility.ReadJsonStringToProtoMessage(processOracleEventBatchCmdArgs.processOracleEventBatchJsonFile, &req)
			if err != nil {
				return err
			}
		} else {
			bytes, err := ioutil.ReadFile(processOracleEventBatchCmdArgs.processOracleEventBatchProtobufFile)
			if err != nil {
				return err
			}

			err = proto.Unmarshal(bytes, &req)
			if err != nil {
				return err
			}
		}

		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		_, err = commonTxObjs.contract.Call("ProcessOracleEventBatch", &req, signer, nil)
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
			fmt.Println("simulated event batch sent")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(processOracleEventBatchCmd)
	processOracleEventBatchCmd.Flags().StringVarP(&processOracleEventBatchCmdArgs.processOracleEventBatchJsonFile, "requestJson", "", "", "file with JSON of ProcessOracleEventBatchRequest protobuf")
	processOracleEventBatchCmd.Flags().StringVarP(&processOracleEventBatchCmdArgs.processOracleEventBatchProtobufFile, "requestPb", "", "", "binary file with ProcessOracleEventBatchRequest protobuf")
}
