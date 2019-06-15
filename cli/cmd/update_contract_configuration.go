package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"github.com/loomnetwork/go-loom/common"
	"github.com/loomnetwork/go-loom/types"
	"math/big"
	"strings"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/spf13/cobra"
)

var updateConfigurationCmdArgs struct {
	setFiatPurchaseContractVersion bool
	fiatPurchaseContractVersion    uint64
	setInitialFiatPurchaseTxId     bool
	initialFiatPurchaseTxId        string
}

var updateConfigurationCmd = &cobra.Command{
	Use:   "update_contract_configuration",
	Short: "changes contract configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		var result zb_calls.EmptyResponse

		initialFiatPurchaseTxId, ok := new(big.Int).SetString(updateConfigurationCmdArgs.initialFiatPurchaseTxId, 10)
		if !ok {
			return fmt.Errorf("failed to parse %s as a number", updateConfigurationCmdArgs.initialFiatPurchaseTxId)
		}

		configurationRequest := zb_calls.UpdateContractConfigurationRequest{
			SetInitialFiatPurchaseTxId:     updateConfigurationCmdArgs.setInitialFiatPurchaseTxId,
			InitialFiatPurchaseTxId:        &types.BigUInt{Value: common.BigUInt{Int: initialFiatPurchaseTxId}},
			SetFiatPurchaseContractVersion: updateConfigurationCmdArgs.setFiatPurchaseContractVersion,
			FiatPurchaseContractVersion:    updateConfigurationCmdArgs.fiatPurchaseContractVersion,
		}

		_, err := commonTxObjs.contract.Call("UpdateContractConfiguration", &configurationRequest, signer, &result)
		if err != nil {
			return fmt.Errorf("error encountered while calling UpdateContractConfiguration: %s", err.Error())
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			output, err := json.Marshal(map[string]interface{}{"success": true})
			if err != nil {
				return err
			}
			fmt.Println(string(output))
		default:
			fmt.Println("success")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateConfigurationCmd)

	updateConfigurationCmd.Flags().BoolVarP(&updateConfigurationCmdArgs.setFiatPurchaseContractVersion, "setFiatPurchaseContractVersion", "", false, "Whether to set fiatPurchaseContractVersion")
	updateConfigurationCmd.Flags().Uint64VarP(&updateConfigurationCmdArgs.fiatPurchaseContractVersion, "fiatPurchaseContractVersion", "", 0, "")
	updateConfigurationCmd.Flags().BoolVarP(&updateConfigurationCmdArgs.setInitialFiatPurchaseTxId, "setInitialFiatPurchaseTxId", "", false, "Whether to set initialFiatPurchaseTxId")
	updateConfigurationCmd.Flags().StringVarP(&updateConfigurationCmdArgs.initialFiatPurchaseTxId, "initialFiatPurchaseTxId", "", "0", "Starting txId used for transaction receipt created by the contract")
}
