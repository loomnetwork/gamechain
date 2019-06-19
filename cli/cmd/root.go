package cmd

import (
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
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

func readKeyFile(path string) ([]byte, error) {
	fileContents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read private key from file: %s", path)
	}

	decodeBuffer := make([]byte, len(fileContents))
	bytesDecoded, err := base64.StdEncoding.Decode(decodeBuffer, fileContents)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid base64 content in private key file: %s", path)
	}

	return decodeBuffer[:bytesDecoded], nil
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
		if command.Use == "merge_json_to_init" {
			return nil
		}

		var err error

		commonTxObjs.privateKey, err = readKeyFile(rootCmdArgs.privateKeyFilePath)
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
	rootCmd.PersistentFlags().StringVarP(&rootCmdArgs.outputFormat, "output", "O", "plaintext", "format of the output (json, plaintext)")

	return rootCmd.Execute()
}
