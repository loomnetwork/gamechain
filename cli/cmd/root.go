package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/loomnetwork/go-loom/client"
	"github.com/spf13/cobra"
)

var rootCmdArgs struct {
	privateKeyFilePath *string
	readURI            *string
	writeURI           *string
}

var commonTxObjs struct {
	privateKey []byte
	contract   *client.Contract
	rpcClient  *client.DAppChainRPCClient
}

func readKeyFile() error {
	fileContents, err := ioutil.ReadFile(*rootCmdArgs.privateKeyFilePath)
	if err != nil {
		return fmt.Errorf("Unable to read private key from file: %s",
			*rootCmdArgs.privateKeyFilePath)
	}

	decodeBuffer := make([]byte, len(fileContents))
	bytesDecoded, err := base64.StdEncoding.Decode(decodeBuffer, fileContents)
	if err != nil {
		return fmt.Errorf("Invalid base64 content in private key file: %s",
			*rootCmdArgs.privateKeyFilePath)
	}

	commonTxObjs.privateKey = decodeBuffer[:bytesDecoded]
	return nil
}

func connectToRPC(readURI string, writeURI string) error {
	rpcClient := client.NewDAppChainRPCClient("default", writeURI, readURI)

	loomAddress, err := rpcClient.Resolve("ZombieBattleground")
	if err != nil {
		return fmt.Errorf("Unable to resolve RPC connection. RPC Error:%s", err.Error())
	}

	commonTxObjs.contract = client.NewContract(rpcClient, loomAddress.Local)
	commonTxObjs.rpcClient = rpcClient

	return nil
}

var rootCmd = &cobra.Command{
	Use:   "zb-cli",
	Short: "ZombieBattleGround cli tool",
	PersistentPreRunE: func(command *cobra.Command, args []string) error {
		var err error

		err = readKeyFile()
		if err != nil {
			return err
		}

		err = connectToRPC(*rootCmdArgs.readURI, *rootCmdArgs.writeURI)
		if err != nil {
			return err
		}

		return nil
	},
}

func Execute() {
	rootCmdArgs.privateKeyFilePath = rootCmd.PersistentFlags().StringP("key", "k", "", "Private key file path")
	rootCmdArgs.readURI = rootCmd.PersistentFlags().StringP("readURI", "r", "", "Read URI for rpc")
	rootCmdArgs.writeURI = rootCmd.PersistentFlags().StringP("writeURI", "w", "", "Write URI for rpc")

	rootCmd.Execute()
}
