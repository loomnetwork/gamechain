// +build evm

package main

import (
	"encoding/base64"
	"io/ioutil"
	"runtime"
	"time"

	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/go-loom/client"
	"github.com/pkg/errors"
)

type Oracle struct {
	cfg *OracleConfig

	gcAddress loom.Address // gamechain address
	gcSigner  auth.Signer

	pcAddress loom.Address // plasmachain address
	pcSigner  auth.Signer

	/*
		chainID   string
		solGateway *ethcontract.MainnetGatewayContract
		goGateway  *DAppChainGateway
		startBlock uint64
		logger     *loom.Logger
		ethClient  *MainnetClient
		// Used to sign tx/data sent to the DAppChain Gateway contract
		signer auth.Signer
		// Private key that should be used to sign tx/data sent to Mainnet Gateway contract
		mainnetPrivateKey     *ecdsa.PrivateKey
		dAppChainPollInterval time.Duration
		mainnetPollInterval   time.Duration
		startupDelay          time.Duration
		reconnectInterval     time.Duration
		mainnetGatewayAddress loom.Address

		numMainnetEventsFetched   uint64
		numMainnetEventsSubmitted uint64

		statusMutex sync.RWMutex
		status      Status

		metrics *Metrics

		hashPool *recentHashPool

		isLoomCoinOracle bool
	*/
}

func CreateOracle(cfg *OracleConfig) (*Oracle, error) {
	return createOracle(cfg, chainID, "tg_oracle", false)
}

func createOracle(cfg *OracleConfig) (*Oracle, error) {
	privKey, err := LoadDappChainPrivateKey(cfg.GameChainPrivateKeyPath)
	if err != nil {
		return nil, err
	}
	gcSigner := auth.NewEd25519Signer(privKey)

	gcAddress := loom.Address{
		ChainID: chainID,
		Local:   loom.LocalAddressFromPublicKey(gcSigner.PublicKey()),
	}

	privKey, err := LoadDappChainPrivateKey(cfg.PlasmaChainPrivateKeyPath)
	if err != nil {
		return nil, err
	}
	pcSigner := auth.NewEd25519Signer(privKey)

	pcAddress := loom.Address{
		ChainID: chainID,
		Local:   loom.LocalAddressFromPublicKey(pcSigner.PublicKey()),
	}

	return &Oracle{
		cfg:       cfg,
		gcAddress: gcAddress,
		gcSigner:  gcSigner,

		pcAddress: pcAddress,
		pcSigner:  pcSigner,

		/*
			chainID:               chainID,
			logger:                loom.NewLoomLogger(cfg.OracleLogLevel, cfg.OracleLogDestination),
			gcAddress:             gcAddress,
			mainnetPrivateKey:     mainnetPrivateKey,
			dAppChainPollInterval: time.Duration(cfg.DAppChainPollInterval) * time.Second,
			mainnetPollInterval:   time.Duration(cfg.MainnetPollInterval) * time.Second,
			startupDelay:          time.Duration(cfg.OracleStartupDelay) * time.Second,
			reconnectInterval:     time.Duration(cfg.OracleReconnectInterval) * time.Second,
			mainnetGatewayAddress: loom.Address{
				ChainID: "eth",
				Local:   common.HexToAddress(cfg.MainnetContractHexAddress).Bytes(),
			},
			status: Status{
				Version:               loomchain.FullVersion(),
				OracleAddress:         address.String(),
				MainnetGatewayAddress: cfg.MainnetContractHexAddress,
			},
			metrics:  NewMetrics(metricSubsystem),
			hashPool: hashPool,

			isLoomCoinOracle: isLoomCoinOracle,
		*/
	}, nil
}

// Status returns some basic info about the current state of the Oracle.
func (orc *Oracle) connect() error {
	var err error

	if orc.gcGateway == nil {
		gcDappClient := client.NewDAppChainRPCClient(orc.gcChainID, orc.cfg.GameChainWriteURI, orc.cfg.GameChainReadURI)
		orc.gcGateway, err = ConnectToDAppChainGateway(gcDappClient, orc.gcAddress, orc.gcSigner, orc.logger)
		if err != nil {
			return errors.Wrap(err, "failed to create gc dappchain gateway")
		}
	}
	if orc.pcGateway == nil {
		pcDappClient := client.NewDAppChainRPCClient(orc.pcChainID, orc.cfg.PlasmaChainWriteURI, orc.cfg.PlasmaChainReadURI)
		orc.pcGateway, err = ConnectToDAppChainGateway(pcDappClient, orc.gpcAddress, orc.pcSigner, orc.logger)
		if err != nil {
			return errors.Wrap(err, "failed to create pc dappchain gateway")
		}

	}
	return nil
}

// RunWithRecovery should run in a goroutine, it will ensure the oracle keeps on running as long
// as it doesn't panic due to a runtime error.
func (orc *Oracle) RunWithRecovery() {
	defer func() {
		if r := recover(); r != nil {
			orc.logger.Error("recovered from panic in Oracle", "r", r)
			// Unless it's a runtime error restart the goroutine
			if _, ok := r.(runtime.Error); !ok {
				time.Sleep(30 * time.Second)
				orc.logger.Info("Restarting Oracle...")
				go orc.RunWithRecovery()
			}
		}
	}()

	// When running in-process give the node a bit of time to spin up.
	if orc.startupDelay > 0 {
		time.Sleep(orc.startupDelay)
	}

	orc.Run()
}

// TODO: Graceful shutdown
func (orc *Oracle) Run() {
	for {
		if err := orc.connect(); err == nil {
			break
		}
		orc.updateStatus()
		time.Sleep(orc.reconnectInterval)
	}

	skipSleep := true
	for {
		if !skipSleep {
			time.Sleep(orc.mainnetPollInterval)
		} else {
			skipSleep = false
		}
		// TODO: should be possible to poll DAppChain & Mainnet at different intervals
		orc.pollGameChain()
		orc.pollPlasmaChain()
	}
}

func (orc *Oracle) pollGameChain() error {

}

func LoadDappChainPrivateKey(path string) ([]byte, error) {
	privKeyB64, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	privKey, err := base64.StdEncoding.DecodeString(string(privKeyB64))
	if err != nil {
		return nil, err
	}

	return privKey, nil
}
