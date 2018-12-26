package oracle

type Config struct {
	// Plasmachain
	PlasmachainPrivateKeyPath string
	PlasmachainChainID        string
	PlasmachainReadURI        string
	PlasmachainWriteURI       string
	PlasmachainEventsURI      string
	PlasmachainContractName   string
	PlasmachainPollInterval   int // in second
	// Gamechain
	GamechainPrivateKeyPath string
	GamechainChainID        string
	GamechainReadURI        string
	GamechainWriteURI       string
	GamechainEventsURI      string
	GamechainContractName   string
	// Oracle log verbosity (debug, info, error, etc.)
	OracleLogLevel       string
	OracleLogDestination string
	// Number of seconds to wait before starting the Oracle.
	OracleStartupDelay int32
	// Number of seconds to wait between reconnection attempts.
	OracleReconnectInterval int32
	// Address on from which the out-of-process Oracle should expose the status & metrics endpoints.
	OracleQueryAddress string
}

func DefaultConfig() *Config {
	return &Config{
		PlasmachainChainID:      "default",
		PlasmachainReadURI:      "http://127.0.0.1:46658/query",
		PlasmachainWriteURI:     "http://127.0.0.1:46658/rpc",
		PlasmachainEventsURI:    "ws://127.0.0.1:%d/queryws",
		PlasmachainContractName: "ZBGCard",
		PlasmachainPollInterval: 10,
		GamechainChainID:        "default",
		GamechainReadURI:        "http://127.0.0.1:46658/query",
		GamechainWriteURI:       "http://127.0.0.1:46658/rpc",
		GamechainEventsURI:      "ws://127.0.0.1:%d/queryws",
		GamechainContractName:   "zombiebattleground",
	}
}
