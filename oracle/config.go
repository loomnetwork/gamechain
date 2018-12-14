package main

type OracleConfig struct {
	GameChainPrivateKeyPath string
	GameChainChainID        string
	GameChainReadURI        string
	GameChainWriteURI       string
	GameChainEventsURI      string
	GameChainPollInterval   int
	GameChainContractName   string

	LoomAuthEndpoint string

	PlasmaChainPrivateKeyPath string
	PlasmaChainChainID        string
	PlasmaChainReadURI        string
	PlasmaChainWriteURI       string
	PlasmaChainEventsURI      string
	PlasmaChainPollInterval   int

	OracleQueryAddress string
}
