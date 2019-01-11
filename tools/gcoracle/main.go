package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"

	"github.com/loomnetwork/gamechain/oracle"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)

func main() {
	cfg, err := parseConfig(nil)
	if err != nil {
		panic(err)
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
}

// Loads oracle.yml or equivalent from one of the usual location, or if overrideCfgDirs is provided
// from one of those config directories.
func parseConfig(overrideCfgDirs []string) (*oracle.Config, error) {
	v := viper.New()
	v.SetConfigName("oracle")
	if len(overrideCfgDirs) == 0 {
		// look for the loom config file in all the places loom itself does
		v.AddConfigPath(".")
		v.AddConfigPath(filepath.Join(".", "config"))
	} else {
		for _, dir := range overrideCfgDirs {
			v.AddConfigPath(dir)
		}
	}
	v.ReadInConfig()
	conf := &oracle.Config{
		// plasmachain
		PlasmachainPrivateKeyPath:     "plasmachain.priv",
		PlasmachainChainID:            "default",
		PlasmachainReadURI:            "http://localhost:46658/query",
		PlasmachainWriteURI:           "http://localhost:46658/rpc",
		PlasmachainEventsURI:          "ws://localhost:9999/queryws",
		PlasmachainContractHexAddress: "0xC5dFc9282BF68DFAd041a04a0c09bE927b093992", // TODO: change me
		// gamechain
		GamechainPrivateKeyPath: "gamechain.priv",
		GamechainChainID:        "default",
		GamechainReadURI:        "http://localhost:46658/query",
		GamechainWriteURI:       "http://localhost:46658/rpc",
		GamechainEventsURI:      "ws://localhost:9999/queryws",
		GamechainContractName:   "ZombieBattleground",
		// oracle
		OracleQueryAddress:   ":8888",
		OracleLogLevel:       "debug",
		OracleLogDestination: "file://oracle.log",
	}
	err := v.Unmarshal(conf)
	if err != nil {
		return nil, err
	}
	return conf, err
}
