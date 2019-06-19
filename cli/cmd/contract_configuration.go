package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/go-loom/common"
	"github.com/loomnetwork/go-loom/types"
	"github.com/spf13/cobra"
	"math/big"
	"strings"
)

var configurationCmdArgs struct {
	fiatPurchaseContractVersion    uint64
	initialFiatPurchaseTxId        string
	useCardLibraryAsUserCollection bool
	debugMode                      bool
}

var configuration_setDataWipeConfigurationCmdArgs struct {
	version   string
	wipeDecks bool
}

var configuration_get = &cobra.Command{
	Use:   "get",
	Short: "get contract configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
		callerAddr := loom.Address{
			ChainID: commonTxObjs.rpcClient.GetChainID(),
			Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
		}

		var resp zb_calls.GetContractConfigurationResponse
		_, err := commonTxObjs.contract.StaticCall("GetContractConfiguration", &zb_calls.EmptyRequest{}, callerAddr, &resp)
		if err != nil {
			return err
		}

		switch strings.ToLower(rootCmdArgs.outputFormat) {
		case "json":
			err := printProtoMessageAsJSONToStdout(resp.Configuration)
			if err != nil {
				return err
			}
		default:
			fmt.Printf("%+v\n", resp.Configuration)
		}

		return nil
	},
}

var configurationCmd = &cobra.Command{
	Use:   "contract_configuration",
	Short: "manipulates contract configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var configuration_setFiatPurchaseContractVersionCmd = &cobra.Command{
	Use:   "set_fiat_purchase_contract_version",
	Short: "sets fiatPurchaseContractVersion",
	RunE: func(cmd *cobra.Command, args []string) error {
		request := &zb_calls.UpdateContractConfigurationRequest{
			SetFiatPurchaseContractVersion: true,
			FiatPurchaseContractVersion:    configurationCmdArgs.fiatPurchaseContractVersion,
		}
		return configurationSetMain(request)
	},
}

var configuration_setInitialFiatPurchaseTxIdCmd = &cobra.Command{
	Use:   "set_initial_fiat_purchase_txid",
	Short: "sets initialFiatPurchaseTxId",
	RunE: func(cmd *cobra.Command, args []string) error {
		initialFiatPurchaseTxId, ok := new(big.Int).SetString(configurationCmdArgs.initialFiatPurchaseTxId, 10)
		if !ok {
			return fmt.Errorf("failed to parse %s as a number", configurationCmdArgs.initialFiatPurchaseTxId)
		}

		request := &zb_calls.UpdateContractConfigurationRequest{
			SetInitialFiatPurchaseTxId: true,
			InitialFiatPurchaseTxId:    &types.BigUInt{Value: common.BigUInt{Int: initialFiatPurchaseTxId}},
		}
		return configurationSetMain(request)
	},
}

var configuration_useCardLibraryAsUserCollectionCmd = &cobra.Command{
	Use:   "set_use_card_library_as_user_collection",
	Short: "sets useCardLibraryAsUserCollection",
	RunE: func(cmd *cobra.Command, args []string) error {
		request := &zb_calls.UpdateContractConfigurationRequest{
			SetUseCardLibraryAsUserCollection: true,
			UseCardLibraryAsUserCollection:    configurationCmdArgs.useCardLibraryAsUserCollection,
		}
		return configurationSetMain(request)
	},
}

var configuration_setDataWipeConfigurationCmd = &cobra.Command{
	Use:   "set_data_wipe_configuration",
	Short: "sets data wipe configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		request := &zb_calls.UpdateContractConfigurationRequest{
			SetDataWipeConfiguration: true,
			DataWipeConfiguration: &zb_data.DataWipeConfiguration{
				Version:   configuration_setDataWipeConfigurationCmdArgs.version,
				WipeDecks: configuration_setDataWipeConfigurationCmdArgs.wipeDecks,
			},
		}
		return configurationSetMain(request)
	},
}

var configuration_setDebugModeCmd = &cobra.Command{
	Use:   "set_debug_mode",
	Short: "sets debug mode state",
	RunE: func(cmd *cobra.Command, args []string) error {
		request := &zb_calls.UpdateContractConfigurationRequest{
			SetDebugMode: true,
			DebugMode:    configurationCmdArgs.debugMode,
		}
		return configurationSetMain(request)
	},
}

func configurationSetMain(configurationRequest *zb_calls.UpdateContractConfigurationRequest) error {
	signer := auth.NewEd25519Signer(commonTxObjs.privateKey)
	var result zb_calls.EmptyResponse

	_, err := commonTxObjs.contract.Call("UpdateContractConfiguration", configurationRequest, signer, &result)
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
}

func init() {
	configuration_setFiatPurchaseContractVersionCmd.Flags().Uint64VarP(&configurationCmdArgs.fiatPurchaseContractVersion, "value", "v", 3, "")
	configuration_setInitialFiatPurchaseTxIdCmd.Flags().StringVarP(&configurationCmdArgs.initialFiatPurchaseTxId, "value", "v", "0", "Starting txId used for transaction receipt created by the contract")
	configuration_useCardLibraryAsUserCollectionCmd.Flags().BoolVarP(&configurationCmdArgs.useCardLibraryAsUserCollection, "value", "v", false, "If false, user personal collection is used, if true, card library is used to make a full fake collection")
	configuration_setDebugModeCmd.Flags().BoolVarP(&configurationCmdArgs.debugMode, "value", "v", false, "Debug mode state")

	configuration_setDataWipeConfigurationCmd.Flags().StringVarP(&configuration_setDataWipeConfigurationCmdArgs.version, "version", "v", "v1", "Data version to wipe on")
	configuration_setDataWipeConfigurationCmd.Flags().BoolVarP(&configuration_setDataWipeConfigurationCmdArgs.wipeDecks, "wipeDecks", "d", false, "Whether to wipe user decks")
	_ = configuration_setFiatPurchaseContractVersionCmd.MarkFlagRequired("value")
	_ = configuration_setInitialFiatPurchaseTxIdCmd.MarkFlagRequired("value")
	_ = configuration_setDebugModeCmd.MarkFlagRequired("value")
	_ = configuration_setDataWipeConfigurationCmd.MarkFlagRequired("version")

	configurationCmd.AddCommand(
		configuration_get,
		configuration_setFiatPurchaseContractVersionCmd,
		configuration_setInitialFiatPurchaseTxIdCmd,
		configuration_useCardLibraryAsUserCollectionCmd,
		configuration_setDataWipeConfigurationCmd,
		configuration_setDebugModeCmd,
	)

	rootCmd.AddCommand(configurationCmd)
}
