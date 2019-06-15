package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"strings"

	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var updateOracleCmdArgs struct {
	newOracleAddressString    string
	newOraclePrivateKeyBase64 string
	chainID                   string
}

var updateOracleCmd = &cobra.Command{
	Use:   "update_oracle",
	Short: "change the oracle or set initial oracle",
	RunE: func(cmd *cobra.Command, args []string) error {
		if updateOracleCmdArgs.newOracleAddressString == "" && updateOracleCmdArgs.newOraclePrivateKeyBase64 == "" {
			return fmt.Errorf("newOracleAddress or ")
		}
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)

		var newOracle *loom.Address
		if updateOracleCmdArgs.newOraclePrivateKeyBase64 == "" {
			newOracleTmp, err := loom.ParseAddress(updateOracleCmdArgs.newOracleAddressString)
			if err != nil {
				return errors.Wrapf(err, "failed to parse new oracle address %s", updateOracleCmdArgs.newOracleAddressString)
			}

			newOracle = &newOracleTmp
		} else {
			newOraclePrivateKey, err := readKeyFile(rootCmdArgs.privateKeyFilePath)
			if err != nil {
				return errors.Wrapf(err, "error while reading new oracle private key file")
			}

			decodeBuffer := make([]byte, len(updateOracleCmdArgs.newOraclePrivateKeyBase64))
			bytesDecoded, err := base64.StdEncoding.Decode(decodeBuffer, []byte(updateOracleCmdArgs.newOraclePrivateKeyBase64))
			if err != nil {
				return errors.Wrapf(err, "invalid base64 content in new oracle private key: %s", updateOracleCmdArgs.newOraclePrivateKeyBase64)
			}

			newOraclePrivateKey = decodeBuffer[:bytesDecoded]
			newOracleSigner := auth.NewEd25519Signer(newOraclePrivateKey)
			newOracle = &loom.Address{
				ChainID: updateOracleCmdArgs.chainID,
				Local:   loom.LocalAddressFromPublicKey(newOracleSigner.PublicKey()),
			}
		}

		fmt.Printf("new oracle address: %s\n", newOracle.String())

		_, err := commonTxObjs.contract.Call("UpdateOracle", &zb_calls.UpdateOracleRequest{
			NewOracle: newOracle.MarshalPB(),
		}, signer, nil)
		if err != nil {
			return errors.Wrap(err, "error when calling UpdateOracle")
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			output, err := json.Marshal(map[string]interface{}{"success": true})
			if err != nil {
				return err
			}
			fmt.Println(string(output))
		default:
			fmt.Println("oracle changed")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateOracleCmd)

	updateOracleCmd.Flags().StringVarP(&updateOracleCmdArgs.newOracleAddressString, "newOracleAddress", "o", "", "Address of the new oracle")
	updateOracleCmd.Flags().StringVarP(&updateOracleCmdArgs.newOraclePrivateKeyBase64, "newOraclePrivateKey", "p", "", "Private key of the new oracle in base64")
	updateOracleCmd.Flags().StringVarP(&updateOracleCmdArgs.chainID, "chainID", "c", "default", "Chain ID")
}
