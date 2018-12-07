package main

type OracleConfig struct {
	GameChainPrivateKeyPath string
	GameChainChainID        string
	GameChainReadURI        string
	GameChainWriteURI       string
	GameChainEventsURI      string
	GameChainPollInterval   int

	PlasmaChainPrivateKeyPath string
	PlasmaChainChainID        string
	PlasmaChainReadURI        string
	PlasmaChainWriteURI       string
	PlasmaChainEventsURI      string
	PlasmaChainPollInterval   int
}
