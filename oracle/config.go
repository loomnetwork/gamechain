package main

type OracleConfig struct {
	GameChainPrivateKeyPath string
	GameChainReadURI        string
	GameChainWriteURI       string
	GameChainEventsURI      string
	GameChainPollInterval   int

	PlasmaChainPrivateKeyPath string
	PlasmaChainReadURI        string
	PlasmaChainWriteURI       string
	PlasmaChainEventsURI      string
	PlasmaChainPollInterval   int
}
