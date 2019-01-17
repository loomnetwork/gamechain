package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/loomnetwork/gamechain/oracle"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:          "gcoracle",
	Short:        "Gamechain Oracle",
	Long:         `The oracle that connect plasmachain and gamechain`,
	Example:      `  gcoracle`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	// plasmchain
	rootCmd.PersistentFlags().String("plasmachain-private-key", "", "Plasmachain Private Key")
	rootCmd.PersistentFlags().String("plasmachain-chain-id", "default", "Plasmachain Chain ID")
	rootCmd.PersistentFlags().String("plasmachain-read-uri", "http://localhost:46658/query", "Plasmachain Read URI")
	rootCmd.PersistentFlags().String("plasmachain-write-uri", "http://localhost:46658/rpc", "Plasmachain Write URI")
	rootCmd.PersistentFlags().String("plasmachain-event-uri", "ws://localhost:9999/queryws", "Plasmachain Events URI")
	rootCmd.PersistentFlags().String("plasmachain-contract-hex-address", "0xea59a949651ffc6d3e039db2d89f4e047301718d", "Plasmachain Contract Hex Address")
	rootCmd.PersistentFlags().Int("plasmachain-poll-interval", 10, "Plasmachain Pool Interval in seconds")
	// gamechain
	rootCmd.PersistentFlags().String("gamechain-private-key", "", "Gamechain Private Key")
	rootCmd.PersistentFlags().String("gamechain-chain-id", "default", "Gamechain Chain ID")
	rootCmd.PersistentFlags().String("gamechain-read-uri", "http://localhost:46658/query", "Gamechain Read URI")
	rootCmd.PersistentFlags().String("gamechain-write-uri", "http://localhost:46658/rpc", "Gamechain Write URI")
	rootCmd.PersistentFlags().String("gamechain-event-uri", "ws://localhost:9999/queryws", "Gamechain Events URI")
	rootCmd.PersistentFlags().String("gamechain-contract-name", "ZombieBattleground", "Gamechain Contract Name")
	rootCmd.PersistentFlags().String("gamechain-card-version", "v3", "Gamechain Card Version")
	// oracle
	rootCmd.PersistentFlags().String("oracle-query-address", ":8888", "Oracle Query Address")
	rootCmd.PersistentFlags().String("oracle-log-level", "debug", "Oracle Log Level")
	rootCmd.PersistentFlags().String("oracle-log-destination", "file://oracle.log", "Oracle Log Destination")
	rootCmd.PersistentFlags().Int("oracle-reconnect-interval", 5, "Oracle Startup Delay in second")
	rootCmd.PersistentFlags().Int("oracle-startup-delay", 10, "Oracle Reconnect Interval in second")

	viper.BindPFlag("plasmachain-private-key", rootCmd.PersistentFlags().Lookup("plasmachain-private-key"))
	viper.BindPFlag("plasmachain-chain-id", rootCmd.PersistentFlags().Lookup("plasmachain-chain-id"))
	viper.BindPFlag("plasmachain-read-uri", rootCmd.PersistentFlags().Lookup("plasmachain-read-uri"))
	viper.BindPFlag("plasmachain-write-uri", rootCmd.PersistentFlags().Lookup("plasmachain-write-uri"))
	viper.BindPFlag("plasmachain-event-uri", rootCmd.PersistentFlags().Lookup("plasmachain-event-uri"))
	viper.BindPFlag("plasmachain-contract-hex-address", rootCmd.PersistentFlags().Lookup("plasmachain-contract-hex-address"))
	viper.BindPFlag("plasmachain-poll-interval", rootCmd.PersistentFlags().Lookup("plasmachain-poll-interval"))

	viper.BindPFlag("gamechain-private-key", rootCmd.PersistentFlags().Lookup("gamechain-private-key"))
	viper.BindPFlag("gamechain-chain-id", rootCmd.PersistentFlags().Lookup("gamechain-chain-id"))
	viper.BindPFlag("gamechain-read-uri", rootCmd.PersistentFlags().Lookup("gamechain-read-uri"))
	viper.BindPFlag("gamechain-write-uri", rootCmd.PersistentFlags().Lookup("gamechain-write-uri"))
	viper.BindPFlag("gamechain-event-uri", rootCmd.PersistentFlags().Lookup("gamechain-event-uri"))
	viper.BindPFlag("gamechain-contract-name", rootCmd.PersistentFlags().Lookup("gamechain-contract-name"))
	viper.BindPFlag("gamechain-card-version", rootCmd.PersistentFlags().Lookup("gamechain-card-version"))

	viper.BindPFlag("oracle-query-address", rootCmd.PersistentFlags().Lookup("oracle-query-address"))
	viper.BindPFlag("oracle-log-level", rootCmd.PersistentFlags().Lookup("oracle-log-level"))
	viper.BindPFlag("oracle-log-destination", rootCmd.PersistentFlags().Lookup("oracle-log-destination"))
	viper.BindPFlag("oracle-reconnect-interval", rootCmd.PersistentFlags().Lookup("oracle-reconnect-interval"))
	viper.BindPFlag("oracle-startup-delay", rootCmd.PersistentFlags().Lookup("oracle-startup-delay"))
}

func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	cfg := &oracle.Config{
		PlasmachainPrivateKey:         viper.GetString("plasmachain-private-key"),
		PlasmachainChainID:            viper.GetString("plasmachain-chain-id"),
		PlasmachainReadURI:            viper.GetString("plasmachain-read-uri"),
		PlasmachainWriteURI:           viper.GetString("plasmachain-write-uri"),
		PlasmachainEventsURI:          viper.GetString("plasmachain-event-uri"),
		PlasmachainContractHexAddress: viper.GetString("plasmachain-contract-hex-address"),
		PlasmachainPollInterval:       viper.GetInt("plasmachain-poll-interval"),
		GamechainPrivateKey:           viper.GetString("gamechain-private-key"),
		GamechainChainID:              viper.GetString("gamechain-chain-id"),
		GamechainReadURI:              viper.GetString("gamechain-read-uri"),
		GamechainWriteURI:             viper.GetString("gamechain-write-uri"),
		GamechainEventsURI:            viper.GetString("gamechain-event-uri"),
		GamechainContractName:         viper.GetString("gamechain-contract-name"),
		GamechainCardVersion:          viper.GetString("gamechain-card-version"),
		OracleQueryAddress:            viper.GetString("oracle-query-address"),
		OracleLogLevel:                viper.GetString("oracle-log-level"),
		OracleLogDestination:          viper.GetString("oracle-log-destination"),
		OracleReconnectInterval:       int32(viper.GetInt("oracle-reconnect-interval")),
		OracleStartupDelay:            int32(viper.GetInt("oracle-startup-delay")),
	}

	fmt.Printf("config %#v\n", cfg)

	if cfg.PlasmachainPrivateKey == "" {
		return errors.New("PlasmachainPrivateKey [--plasmachain-private-key] is required")
	}
	if cfg.GamechainPrivateKey == "" {
		return errors.New("GamechainPrivateKey [--gamechain-private-key] is required")
	}

	orc, err := oracle.CreateOracle(cfg, "gcoracle")
	if err != nil {
		panic(err)
	}

	go orc.RunWithRecovery()

	http.HandleFunc("/status", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(orc.Status())
	})

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(cfg.OracleQueryAddress, nil))
	return nil
}
