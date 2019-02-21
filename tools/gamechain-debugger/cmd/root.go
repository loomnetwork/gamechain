package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/loomnetwork/gamechain/tools/gamechain-debugger/controller"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmdArgs struct {
	privateKeyFilePath string
	port               string
	cliFilePath        string
}

func verifyFlag() error {
	if err := readKeyFile(); err != nil {
		return err
	}
	if err := checkCli(); err != nil {
		return err
	}
	return nil
}

func checkCli() error {
	if _, err := os.Stat(rootCmdArgs.cliFilePath); os.IsNotExist(err) {
		fmt.Println("Cli file doesn't exist, please use flag -b [file path]")
		return err
	}
	return nil
}

func readKeyFile() error {
	fileContents, err := ioutil.ReadFile(rootCmdArgs.privateKeyFilePath)
	if err != nil {
		fmt.Printf("Cannot read file, please use flag -k [key path]")
		return err
	}

	decodeBuffer := make([]byte, len(fileContents))
	_, err = base64.StdEncoding.Decode(decodeBuffer, fileContents)
	if err != nil {
		return fmt.Errorf("invalid base64 content in private key file: %s",
			rootCmdArgs.privateKeyFilePath)
	}

	return nil
}

var rootCmd = &cobra.Command{
	Use:   "gamechain-debugger",
	Short: "A tool for gamestate debugging",
	Run: func(cmd *cobra.Command, args []string) {
		if verifyFlag() != nil {
			return
		}
		Serve()
	},
}

func Execute() error {
	rootCmd.PersistentFlags().StringVarP(&rootCmdArgs.privateKeyFilePath, "key", "k", "priv", "Private key file path")
	rootCmd.PersistentFlags().StringVarP(&rootCmdArgs.port, "port", "p", "3000", "Running port for web server")
	rootCmd.PersistentFlags().StringVarP(&rootCmdArgs.cliFilePath, "cli", "b", "./zb-cli", "path to zb-cli")
	viper.BindPFlag("key", rootCmd.PersistentFlags().Lookup("key"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("cli", rootCmd.PersistentFlags().Lookup("cli"))
	return rootCmd.Execute()
}

func Serve() {
	fmt.Println("Gamechain Debugger is running on port " + rootCmdArgs.port)
	r := mux.NewRouter()
	mc := controller.NewMainController(rootCmdArgs.cliFilePath, rootCmdArgs.privateKeyFilePath)
	r.HandleFunc("/", mc.GamechainDebugger)
	r.HandleFunc("/client_state", mc.ClientStateDebugger)
	r.HandleFunc("/get_state/{MatchId}", mc.GetState)
	r.HandleFunc("/save_state/{MatchId}", mc.SaveState)
	http.ListenAndServe(":"+rootCmdArgs.port, r)
}
