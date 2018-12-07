package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/spf13/viper"
)

func main() {
	cfg, err := parseConfig(nil)
	if err != nil {
		panic(err)
	}

	orc, err := CreateOracle(cfg)
	if err != nil {
		panic(err)
	}

	go orc.RunWithRecovery()

	http.HandleFunc("/status", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		//json.NewEncoder(w).Encode(orc.Status())
	})

	//http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(cfg.OracleQueryAddress, nil))
}

// Loads oracle.yml or equivalent from one of the usual location, or if overrideCfgDirs is provided
// from one of those config directories.
func parseConfig(overrideCfgDirs []string) (*OracleConfig, error) {
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
	conf := &OracleConfig{
		GameChainPrivateKeyPath: "priv",
		GameChainChainID:        "default",
		GameChainReadURI:        "http://localhost:46658/query",
		GameChainWriteURI:       "http://localhost:46658/rpc",
		GameChainEventsURI:      "",
	}
	err := v.Unmarshal(conf)
	if err != nil {
		return nil, err
	}
	return conf, err
}
