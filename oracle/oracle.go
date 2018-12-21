package oracle

import (
	"encoding/base64"
	"io/ioutil"
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/go-loom/client"
	lmtype "github.com/loomnetwork/go-loom/plugin/types"
	"github.com/loomnetwork/loomchain"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type (
	PlasmachainEventData = lmtype.EventData
)

type plasmachainEventInfo struct {
	BlockNum uint64
	TxIdx    uint
	Event    *PlasmachainEventData
}

type Status struct {
	Version                  string
	OracleAddress            string
	DAppChainGatewayAddress  string
	MainnetGatewayAddress    string
	NextPlasmachainBlockNum  uint64    `json:",string"`
	MainnetGatewayLastSeen   time.Time // TODO: hook this up
	DAppChainGatewayLastSeen time.Time
	// Number of Plamachain events submitted to the DAppChain Gateway successfully
	NumPlamachainEventsFetched uint64 `json:",string"`
	// Total number of Plamachain events fetched
	NumPlamachainEventsSubmitted uint64 `json:",string"`
}

type Oracle struct {
	cfg     Config
	chainID string
	// Plasmachain address
	pcAddress      loom.Address
	pcSigner       auth.Signer
	pcChainID      string
	pcContractName string
	pcGateway      *PlasmachainGateway
	pcPollInterval time.Duration
	// Gamechain
	gcAddress      loom.Address
	gcSigner       auth.Signer
	gcChainID      string
	gcContractName string
	gcGateway      *GamechainGateway
	// oracle
	logger            *loom.Logger
	reconnectInterval time.Duration
	statusMutex       sync.RWMutex
	status            Status
	metrics           *Metrics
	// Used to sign tx/data sent to the DAppChain Gateway contract
	signer                        auth.Signer
	startupDelay                  time.Duration
	startBlock                    uint64
	numPlasmachainEventsFetched   uint64
	numPlasmachainEventsSubmitted uint64

	hashPool *recentHashPool
}

func CreateOracle(cfg *Config, metricSubsystem string) (*Oracle, error) {
	privKey, err := LoadDappChainPrivateKey(cfg.GamechainPrivateKeyPath)
	if err != nil {
		return nil, err
	}
	gcSigner := auth.NewEd25519Signer(privKey)
	gcAddress := loom.Address{
		ChainID: cfg.GamechainChainID,
		Local:   loom.LocalAddressFromPublicKey(gcSigner.PublicKey()),
	}

	privKey, err = LoadDappChainPrivateKey(cfg.PlasmachainPrivateKeyPath)
	if err != nil {
		return nil, err
	}
	pcSigner := auth.NewEd25519Signer(privKey)
	pcAddress := loom.Address{
		ChainID: cfg.PlasmachainChainID,
		Local:   loom.LocalAddressFromPublicKey(pcSigner.PublicKey()),
	}

	hashPool := newRecentHashPool(time.Duration(cfg.PlasmachainPollInterval) * time.Second * 4)
	hashPool.startCleanupRoutine()

	return &Oracle{
		cfg:               *cfg,
		gcAddress:         gcAddress,
		gcSigner:          gcSigner,
		pcAddress:         pcAddress,
		pcSigner:          pcSigner,
		pcContractName:    cfg.PlasmachainContractName,
		metrics:           NewMetrics(metricSubsystem),
		logger:            loom.NewLoomLogger("info", ""),
		pcPollInterval:    time.Duration(cfg.PlasmachainPollInterval) * time.Second,
		startupDelay:      time.Duration(cfg.OracleStartupDelay) * time.Second,
		reconnectInterval: time.Duration(cfg.OracleReconnectInterval) * time.Second,
		status: Status{
			Version: loomchain.FullVersion(),
		},
		hashPool: hashPool,
	}, nil
}

// Status returns some basic info about the current state of the Oracle.
func (orc *Oracle) Status() *Status {
	orc.statusMutex.RLock()

	s := orc.status

	orc.statusMutex.RUnlock()
	return &s
}

func (orc *Oracle) updateStatus() {
	orc.statusMutex.Lock()

	orc.status.NextPlasmachainBlockNum = orc.startBlock
	orc.status.NumPlamachainEventsFetched = orc.numPlasmachainEventsFetched
	orc.status.NumPlamachainEventsSubmitted = orc.numPlasmachainEventsSubmitted

	// if orc.goGateway != nil {
	// 	orc.status.DAppChainGatewayAddress = orc.goGateway.Address.String()
	// 	orc.status.DAppChainGatewayLastSeen = orc.goGateway.LastResponseTime
	// }

	orc.statusMutex.Unlock()
}

// Status returns some basic info about the current state of the Oracle.
func (orc *Oracle) connect() error {
	var err error

	if orc.gcGateway == nil {
		dappClient := client.NewDAppChainRPCClient(orc.cfg.GamechainChainID, orc.cfg.GamechainWriteURI, orc.cfg.GamechainReadURI)
		orc.gcGateway, err = ConnectToGamechainGateway(dappClient, orc.gcAddress, orc.cfg.GamechainContractName, orc.gcSigner, orc.logger)
		if err != nil {
			return errors.Wrap(err, "failed to create gamechain gateway")
		}
		orc.logger.Info("connected to Gamechain")
	}

	if orc.pcGateway == nil {
		dappClient := client.NewDAppChainRPCClient(orc.cfg.PlasmachainChainID, orc.cfg.PlasmachainWriteURI, orc.cfg.PlasmachainReadURI)
		orc.pcGateway, err = ConnectToPlasmachainGateway(dappClient, orc.gcAddress, orc.cfg.PlasmachainContractName, orc.pcSigner, orc.logger)
		if err != nil {
			return errors.Wrap(err, "failed to create plasmachain gateway")
		}
		orc.logger.Info("connected to Plasmachain")
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
	if orc.startupDelay > 0 {
		time.Sleep(orc.startupDelay)
	}

	orc.Run()
}

// TODO: Graceful shutdown
func (orc *Oracle) Run() {
	for {
		if err := orc.connect(); err == nil {
			log.Error(err)
			break
		}
		orc.updateStatus()
		time.Sleep(orc.reconnectInterval)
	}
	go orc.pollPlasmaChain()
}

func (orc *Oracle) pollPlasmaChain() error {
	lastPlasmachainBlockNum, err := orc.pcGateway.LastBlockNumber()
	if err != nil {
		return err
	}

	startBlock := lastPlasmachainBlockNum + 1
	if orc.startBlock > startBlock {
		startBlock = orc.startBlock
	}

	// TODO: limit max block range per batch
	latestBlock, err := orc.gcGateway.LastPlasmaBlockNumber()
	if err != nil {
		orc.logger.Error("failed to obtain latest Plasmachain block number", "err", err)
		return err
	}

	if latestBlock < startBlock {
		// Wait for Plasmachain to produce a new block...
		return nil
	}

	events, err := orc.fetchEvents(startBlock, latestBlock)
	if err != nil {
		orc.logger.Error("failed to fetch events from Ethereum", "err", err)
		return err
	}

	if len(events) > 0 {
		orc.numPlasmachainEventsFetched = orc.numPlasmachainEventsFetched + uint64(len(events))
		orc.updateStatus()

		if err := orc.gcGateway.ProcessEventBatch(events); err != nil {
			return err
		}

		orc.numPlasmachainEventsSubmitted = orc.numPlasmachainEventsSubmitted + uint64(len(events))
		orc.metrics.SubmittedPlasmachainEvents(len(events))
		orc.updateStatus()
	}

	orc.startBlock = latestBlock + 1
	return nil
}

// Fetches all relevent events from an Plasmachain node from startBlock to endBlock (inclusive)
func (orc *Oracle) fetchEvents(startBlock, endBlock uint64) ([]*plasmachainEventInfo, error) {
	// NOTE: Currently either all blocks from w.StartBlock are processed successfully or none are.
	filterOpts := &bind.FilterOpts{
		Start: startBlock,
		End:   &endBlock,
	}

	var openPacks []*plasmachainEventInfo
	var err error

	openPacks, err = orc.fetchOpenPacks(filterOpts)
	if err != nil {
		return nil, err
	}

	events := make(
		[]*plasmachainEventInfo, 0,
		len(openPacks),
	)
	events = append(events, openPacks...)

	sortPlasmachainEvents(events)
	sortedEvents := make([]*plasmachainEventInfo, len(events))
	for i, event := range events {
		sortedEvents[i] = event
	}

	if len(events) > 0 {
		orc.logger.Debug("fetched Plasmachain events",
			"startBlock", startBlock,
			"endBlock", endBlock,
			"open-packs", len(openPacks),
		)
	}

	return sortedEvents, nil
}

func sortPlasmachainEvents(events []*plasmachainEventInfo) {
	// Sort events by block & tx index (within the block)?
	// Need to check if plasmachain event contains TxIdx
	sort.Slice(events, func(i, j int) bool {
		if events[i].BlockNum == events[j].BlockNum {
			return events[i].TxIdx < events[j].TxIdx
		}
		return events[i].BlockNum < events[j].BlockNum
	})
}

func (orc *Oracle) fetchOpenPacks(filterOpts *bind.FilterOpts) ([]*plasmachainEventInfo, error) {
	panic("not implemented")
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
