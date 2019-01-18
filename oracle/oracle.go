package oracle

import (
	"encoding/base64"
	"io/ioutil"
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	orctype "github.com/loomnetwork/gamechain/types/oracle"
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/go-loom/client"
	ltypes "github.com/loomnetwork/go-loom/types"

	"github.com/pkg/errors"
)

type Status struct {
	Version                    string
	OracleAddress              string
	GamechainGatewayAddress    string
	GamechainGatewayLastSeen   time.Time
	PlasmachainGatewayAddress  string
	PlasmachainGatewayLastSeen time.Time
	NextPlasmachainBlockNum    uint64 `json:",string"`
	// Number of Plamachain events submitted to the DAppChain Gateway successfully
	NumPlamachainEventsFetched uint64 `json:",string"`
	// Total number of Plamachain events fetched
	NumPlamachainEventsSubmitted uint64 `json:",string"`
}

type Oracle struct {
	cfg     Config
	chainID string
	// Plasmachain address
	pcAddress         loom.Address
	pcSigner          auth.Signer
	pcChainID         string
	pcContractAddress string
	pcGateway         *PlasmachainGateway
	pcPollInterval    time.Duration
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
	privKey, err := LoadDappChainPrivateKey(cfg.GamechainPrivateKey)
	if err != nil {
		return nil, err
	}
	gcSigner := auth.NewEd25519Signer(privKey)
	gcAddress := loom.Address{
		ChainID: cfg.GamechainChainID,
		Local:   loom.LocalAddressFromPublicKey(gcSigner.PublicKey()),
	}

	privKey, err = LoadDappChainPrivateKey(cfg.PlasmachainPrivateKey)
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
		pcContractAddress: cfg.PlasmachainContractHexAddress,
		metrics:           NewMetrics(metricSubsystem),
		logger:            loom.NewLoomLogger(cfg.OracleLogLevel, cfg.OracleLogDestination),
		pcPollInterval:    time.Duration(cfg.PlasmachainPollInterval) * time.Second,
		startupDelay:      time.Duration(cfg.OracleStartupDelay) * time.Second,
		reconnectInterval: time.Duration(cfg.OracleReconnectInterval) * time.Second,
		status: Status{
			Version: "1.0.0",
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

	if orc.gcGateway != nil {
		orc.status.GamechainGatewayAddress = orc.gcGateway.Address.String()
		orc.status.GamechainGatewayLastSeen = orc.gcGateway.LastResponseTime
	}
	if orc.pcGateway != nil {
		orc.status.PlasmachainGatewayAddress = orc.pcGateway.Address.String()
		orc.status.PlasmachainGatewayLastSeen = orc.pcGateway.LastResponseTime
	}

	orc.statusMutex.Unlock()
}

// Status returns some basic info about the current state of the Oracle.
func (orc *Oracle) connect() error {
	var err error
	if orc.pcGateway == nil {
		dappClient := client.NewDAppChainRPCClient(orc.cfg.PlasmachainChainID, orc.cfg.PlasmachainWriteURI, orc.cfg.PlasmachainReadURI)
		orc.pcGateway, err = ConnectToPlasmachainGateway(dappClient, orc.pcAddress, orc.cfg.PlasmachainContractHexAddress, orc.pcSigner, orc.logger)
		if err != nil {
			return errors.Wrap(err, "failed to create plasmachain gateway")
		}
		orc.logger.Info("connected to Plasmachain")
	}

	if orc.gcGateway == nil {
		dappClient := client.NewDAppChainRPCClient(orc.cfg.GamechainChainID, orc.cfg.GamechainWriteURI, orc.cfg.GamechainReadURI)
		orc.gcGateway, err = ConnectToGamechainGateway(dappClient, orc.gcAddress, orc.cfg.GamechainContractName, orc.gcSigner, orc.logger, orc.cfg.GamechainCardVersion)
		if err != nil {
			return errors.Wrap(err, "failed to create gamechain gateway")
		}
		orc.logger.Info("connected to Gamechain")
	}

	return nil
}

// RunWithRecovery should run in a goroutine, it will ensure the oracle keeps on running as long
// as it doesn't panic due to a runtime error.
func (orc *Oracle) RunWithRecovery() {
	orc.logger.Info("Running Oracle...")
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
		err := orc.connect()
		if err == nil {
			break
		}
		orc.logger.Info(err.Error())
		orc.updateStatus()
		time.Sleep(orc.reconnectInterval)
	}

	skipSleep := true
	for {
		if !skipSleep {
			time.Sleep(orc.pcPollInterval)
		} else {
			skipSleep = false
		}
		err := orc.pollPlasmaChain()
		if err != nil {
			orc.logger.Error(err.Error())
		}
	}

}

func (orc *Oracle) pollPlasmaChain() error {
	orc.logger.Info("polling Plasmachain")
	lastPlasmachainBlockNum, err := orc.gcGateway.LastPlasmaBlockNumber()
	if err != nil {
		return err
	}
	startBlock := lastPlasmachainBlockNum + 1
	if orc.startBlock > startBlock {
		startBlock = orc.startBlock
	}

	// TODO: limit max block range per batch
	latestBlock, err := orc.getLatestEthBlockNumber()
	if err != nil {
		orc.logger.Error("failed to obtain latest Plasmachain block number from Gamechain", "err", err)
		return err
	}

	if latestBlock < startBlock {
		// Wait for Plasmachain to produce a new block...
		return nil
	}

	orc.logger.Info("fetching blocks", "latestBlock", startBlock, "latestBlock", latestBlock)
	events, err := orc.fetchEvents(startBlock, latestBlock)
	if err != nil {
		orc.logger.Error("failed to fetch events from Plasmachain", "err", err)
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

func (orc *Oracle) getLatestEthBlockNumber() (uint64, error) {
	return orc.pcGateway.LastBlockNumber()
}

// Fetches all relevent events from an Plasmachain node from startBlock to endBlock (inclusive)
func (orc *Oracle) fetchEvents(startBlock, endBlock uint64) ([]*orctype.PlasmachainEvent, error) {
	// NOTE: Currently either all blocks from w.StartBlock are processed successfully or none are.
	filterOpts := &bind.FilterOpts{
		Start: startBlock,
		End:   &endBlock,
	}

	var generatedCards []*plasmachainEventInfo
	var err error

	generatedCards, err = orc.fetchGeneratedCard(filterOpts)
	if err != nil {
		return nil, err
	}

	events := make(
		[]*plasmachainEventInfo, 0,
		len(generatedCards),
	)
	events = append(events, generatedCards...)

	sortPlasmachainEvents(events)
	sortedEvents := make([]*orctype.PlasmachainEvent, len(events))
	for i, event := range events {
		sortedEvents[i] = event.Event
	}

	if len(events) > 0 {
		orc.logger.Debug("fetched Plasmachain events",
			"startBlock", startBlock,
			"endBlock", endBlock,
			"generatedCard", len(generatedCards),
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

func LoadDappChainPrivateKeyFile(path string) ([]byte, error) {
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

func LoadDappChainPrivateKey(privKeyB64 string) ([]byte, error) {
	privKey, err := base64.StdEncoding.DecodeString(string(privKeyB64))
	if err != nil {
		return nil, err
	}

	return privKey, nil
}

func (orc *Oracle) fetchGeneratedCard(filterOpts *bind.FilterOpts) ([]*plasmachainEventInfo, error) {
	var err error
	var numEvents int
	defer func(begin time.Time) {
		orc.metrics.MethodCalled(begin, "fetchGeneratedCard", err)
		orc.metrics.FetchedPlasmachainEvents(numEvents, "GeneratedCard")
		orc.updateStatus()
	}(time.Now())

	it, err := orc.pcGateway.cardFaucet.FilterGeneratedCard(filterOpts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get logs for GeneratedCard")
	}
	var chainID = orc.pcGateway.client.GetChainID()
	events := []*plasmachainEventInfo{}
	for {
		ok := it.Next()
		if ok {
			ev := it.Event
			receipt, err := orc.pcGateway.client.GetEvmTxReceipt(ev.Raw.TxHash.Bytes())
			if err != nil {
				orc.logger.Error(err.Error(), "txHash", ev.Raw.TxHash.Hex())
				return nil, err
			}
			contractAddr := loom.Address{ChainID: chainID, Local: receipt.ContractAddress}.MarshalPB()
			events = append(events, &plasmachainEventInfo{
				BlockNum: ev.Raw.BlockNumber,
				TxIdx:    ev.Raw.TxIndex,
				Event: &PlasmachainEvent{
					EthBlock: ev.Raw.BlockNumber,
					Payload: &PlasmachainGeneratedCardEvent{
						Card: &PlasmachainGeneratedCard{
							Owner:    receipt.CallerAddress,
							CardID:   &ltypes.BigUInt{Value: *loom.NewBigUInt(ev.CardId)},
							Amount:   &ltypes.BigUInt{Value: *loom.NewBigUIntFromInt(1)},
							Contract: contractAddr,
						},
					},
				},
			})
		} else {
			err = it.Error()
			if err != nil {
				return nil, errors.Wrap(err, "Failed to get event data for GeneratedCard")
			}
			it.Close()
			break
		}
	}
	numEvents = len(events)
	return events, nil
}
