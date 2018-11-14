package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/loomnetwork/go-loom/client"
	"github.com/spf13/cobra"
)

var rootCmdArgs struct {
	privateKeyFilePath string
	readURI            string
	writeURI           string
	chainID            string
	outputFormat       string
}

var commonTxObjs struct {
	privateKey []byte
	contract   *client.Contract
	rpcClient  *client.DAppChainRPCClient
}

func readKeyFile() error {
	fileContents, err := ioutil.ReadFile(rootCmdArgs.privateKeyFilePath)
	if err != nil {
		return fmt.Errorf("unable to read private key from file: %s",
			rootCmdArgs.privateKeyFilePath)
	}

	decodeBuffer := make([]byte, len(fileContents))
	bytesDecoded, err := base64.StdEncoding.Decode(decodeBuffer, fileContents)
	if err != nil {
		return fmt.Errorf("invalid base64 content in private key file: %s",
			rootCmdArgs.privateKeyFilePath)
	}

	commonTxObjs.privateKey = decodeBuffer[:bytesDecoded]
	return nil
}

func connectToRPC(chainID, readURI, writeURI string) error {
	rpcClient := client.NewDAppChainRPCClient(chainID, writeURI, readURI)

	loomAddress, err := rpcClient.Resolve("ZombieBattleground")
	if err != nil {
		return fmt.Errorf("unable to resolve RPC connection. RPC Error:%s", err.Error())
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
			return fmt.Errorf("error while reading private key file: %s", err.Error())
		}

		err = connectToRPC(rootCmdArgs.chainID, rootCmdArgs.readURI, rootCmdArgs.writeURI)
		if err != nil {
			return fmt.Errorf("error while establishing RPC connection: %s", err.Error())
		}

		return nil
	},
}

func Execute() error {
	rootCmd.PersistentFlags().StringVarP(&rootCmdArgs.privateKeyFilePath, "key", "k", "priv.key", "Private key file path")
	rootCmd.PersistentFlags().StringVarP(&rootCmdArgs.readURI, "readURI", "r", "http://localhost:46658/query", "Read URI for rpc")
	rootCmd.PersistentFlags().StringVarP(&rootCmdArgs.writeURI, "writeURI", "w", "http://localhost:46658/rpc", "Write URI for rpc")
	rootCmd.PersistentFlags().StringVarP(&rootCmdArgs.chainID, "chainID", "c", "default", "Chain ID")
	rootCmd.PersistentFlags().StringVar(&rootCmdArgs.outputFormat, "output", "plaintext", "format of the output (json, plaintext)")

	return rootCmd.Execute()
}
