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
	log "github.com/sirupsen/logrus"
)

/*
* create rpc clients and DappchainGateway for gamechain (gc) and plasmachain (pc)
* 	done TODO: use correct contract names in `Resolve`
* open a websocket conn for events in `Run`
* websocket event should call the DappchainGateway's process events method
* http calls to loomauth:
* 	I think the endpoint is `/user/reward`, not sure and need to confirm.
* 	another gateway kind of thing for that?
* write contract method in ZB contract that stores the data coming from the pc events
* get config from file etc (have defaults too?)
 */
type Oracle struct {
	cfg *OracleConfig

	gcAddress      loom.Address // gamechain address
	gcSigner       auth.Signer
	gcChainID      string
	gcContractName string
	gcGateway      *DAppChainGateway

	pcAddress         loom.Address // plasmachain address
	pcSigner          auth.Signer
	pcChainID         string
	pcContractName    string
	pcGateway         *DAppChainGateway
	logger            *loom.Logger
	reconnectInterval time.Duration

	/*
		chainID   string
		solGateway *ethcontract.MainnetGatewayContract
		goGateway  *DAppChainGateway
		startBlock uint64
		ethClient  *MainnetClient
		// Used to sign tx/data sent to the DAppChain Gateway contract
		signer auth.Signer
		// Private key that should be used to sign tx/data sent to Mainnet Gateway contract
		mainnetPrivateKey     *ecdsa.PrivateKey
		dAppChainPollInterval time.Duration
		mainnetPollInterval   time.Duration
		startupDelay          time.Duration
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
	return createOracle(cfg)
}

func createOracle(cfg *OracleConfig) (*Oracle, error) {
	privKey, err := LoadDappChainPrivateKey(cfg.GameChainPrivateKeyPath)
	if err != nil {
		return nil, err
	}
	gcSigner := auth.NewEd25519Signer(privKey)

	gcAddress := loom.Address{
		ChainID: cfg.GameChainChainID,
		Local:   loom.LocalAddressFromPublicKey(gcSigner.PublicKey()),
	}

	/*
			privKey, err = LoadDappChainPrivateKey(cfg.PlasmaChainPrivateKeyPath)
			if err != nil {
				return nil, err
			}
			pcSigner := auth.NewEd25519Signer(privKey)

		pcAddress := loom.Address{
			ChainID: cfg.PlasmaChainChainID,
			Local:   loom.LocalAddressFromPublicKey(pcSigner.PublicKey()),
		}
	*/
	return &Oracle{
		cfg:            cfg,
		gcAddress:      gcAddress,
		gcSigner:       gcSigner,
		gcContractName: "ZombieBattleground",

		//pcAddress:      pcAddress,
		//pcSigner:       pcSigner,
		//pcContractName: "ZBGCard",

		logger: loom.NewLoomLogger("info", ""),
		/*
			chainID:               chainID,
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
		orc.gcGateway, err = ConnectToDAppChainGateway(gcDappClient, orc.gcAddress, orc.gcContractName, orc.gcSigner, orc.logger)
		if err != nil {
			return errors.Wrap(err, "failed to create gc dappchain gateway")
		}
	}

	if orc.pcGateway == nil {
		pcDappClient := client.NewDAppChainRPCClient(orc.pcChainID, orc.cfg.PlasmaChainWriteURI, orc.cfg.PlasmaChainReadURI)
		orc.pcGateway, err = ConnectToDAppChainGateway(pcDappClient, orc.pcAddress, orc.pcContractName, orc.pcSigner, orc.logger)
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
			log.Error("recovered from panic in Oracle", "r", r)
			// Unless it's a runtime error restart the goroutine
			if _, ok := r.(runtime.Error); !ok {
				time.Sleep(30 * time.Second)
				log.Info("Restarting Oracle...")
				go orc.RunWithRecovery()
			}
		}
	}()

	// When running in-process give the node a bit of time to spin up.
	//if orc.startupDelay > 0 {
	//	time.Sleep(orc.startupDelay)
	//}

	orc.Run()
}

// TODO: Graceful shutdown
func (orc *Oracle) Run() {
	for {
		if err := orc.connect(); err == nil {
			break
		}
		//	orc.updateStatus()
		time.Sleep(orc.reconnectInterval)
	}
	go orc.listenToGameChain()
	go orc.listenToPlasmaChain()
}

func (orc *Oracle) listenToGameChain() error {
	log.Info("Listening to GameChain")
	/*
		gcEventClient, err := NewDAppChainEventClient(orc.gcAddress, orc.cfg.GameChainEventsURI)
		if err != nil {
			return err
		}

		gcEventClient.WatchTopic()
	*/
	return nil
}

func (orc *Oracle) listenToPlasmaChain() error {
	log.Info("Listening to PlasmaChain")
	return nil
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

func ConnectToDAppChainGateway(
	loomClient *client.DAppChainRPCClient, caller loom.Address, contractName string, signer auth.Signer,
	logger *loom.Logger,
) (*DAppChainGateway, error) {
	gatewayAddr, err := loomClient.Resolve(contractName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve Gateway Go contract address")
	}

	return &DAppChainGateway{
		Address:          gatewayAddr,
		LastResponseTime: time.Now(),
		contract:         client.NewContract(loomClient, gatewayAddr.Local),
		caller:           caller,
		signer:           signer,
		//	logger:           logger,
	}, nil
}

type DAppChainGateway struct {
	Address loom.Address
	// Timestamp of the last successful response from the DAppChain
	LastResponseTime time.Time

	contract *client.Contract
	caller   loom.Address
	//logger   *loom.Logger
	signer auth.Signer
}
