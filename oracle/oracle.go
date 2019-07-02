package oracle

import (
	"encoding/base64"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/loomnetwork/gamechain/tools/battleground_utility"
	"github.com/loomnetwork/go-loom/common"
	"io/ioutil"
	"reflect"
	"runtime"
	"runtime/debug"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	orctype "github.com/loomnetwork/gamechain/types/oracle"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/go-loom/client"
	ptypes "github.com/loomnetwork/go-loom/plugin/types"
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
	NextPlasmachainBlockNumber uint64 `json:",string"`
	// Number of Plamachain events submitted to the DAppChain Gateway successfully
	PlasmachainEventsFetchedCount uint64 `json:",string"`
	// Total number of Plamachain events fetched
	PlasmachainEventsSubmittedCount uint64 `json:",string"`
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
	logger := loom.NewLoomLogger(cfg.OracleLogLevel, cfg.OracleLogDestination)

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

	logger.Info("Gamechain address " + gcAddress.String())
	logger.Info("Plasmachain address " + pcAddress.String())

	return &Oracle{
		cfg:               *cfg,
		gcAddress:         gcAddress,
		gcSigner:          gcSigner,
		pcAddress:         pcAddress,
		pcSigner:          pcSigner,
		pcContractAddress: cfg.PlasmachainContractHexAddress,
		metrics:           NewMetrics(metricSubsystem),
		logger:            logger,
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

	orc.status.NextPlasmachainBlockNumber = orc.startBlock
	orc.status.PlasmachainEventsFetchedCount = orc.numPlasmachainEventsFetched
	orc.status.PlasmachainEventsSubmittedCount = orc.numPlasmachainEventsSubmitted

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
			orc.logger.Error("recovered from panic in Oracle", "r", r, "stacktrace", string(debug.Stack()))
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
		err := orc.doCommunicationRound()
		if err != nil {
			orc.logger.Error(err.Error())
		}
	}
}

func (orc *Oracle) doCommunicationRound() error {
	latestPlasmaBlock, err := orc.pollPlasmachainForEvents()
	if err != nil {
		return errors.Wrap(err, "failed to poll Plasmachain for events")
	}

	if err := orc.executeGamechainCommands(latestPlasmaBlock); err != nil {
		return errors.Wrap(err, "failed to execute Gamechain commands")
	}

	return nil
}

func (orc *Oracle) executeGamechainCommands(latestPlasmaBlock uint64) error {
	orc.logger.Debug("Fetching Gamechain commands")

	commandRequests, err := orc.gcGateway.GetOracleCommandRequestList()
	if err != nil {
		return err
	}

	commandResponses := make([]*orctype.OracleCommandResponse, 0, len(commandRequests))

	orc.logger.Debug("Executing Gamechain commands", "len(commandRequests)", len(commandRequests))
	for _, commandRequestWrapper := range commandRequests {
		orc.logger.Info("Executing command", "commandId", commandRequestWrapper.CommandId, "commandType", reflect.TypeOf(commandRequestWrapper.GetCommand()).String())
		switch commandRequest := commandRequestWrapper.Command.(type) {
		case *orctype.OracleCommandRequest_GetUserFullCardCollection:
			userAddress := commandRequest.GetUserFullCardCollection.UserAddress
			tokensOwned, err := orc.pcGateway.GetTokensOwned(loom.UnmarshalAddressPB(userAddress).Local)
			if err != nil {
				return err
			}

			orc.logger.Info("GetUserFullCardCollection", "userAddress", loom.UnmarshalAddressPB(userAddress), "tokensOwned", len(tokensOwned))

			response := &orctype.OracleCommandResponse_GetUserFullCardCollectionCommandResponse{
				UserAddress: userAddress,
				BlockHeight: latestPlasmaBlock,
			}

			for _, tokensOwnedResponseItem := range tokensOwned {
				response.OwnedCards = append(response.OwnedCards, &orctype.RawCardCollectionCard{
					CardTokenId: battleground_utility.MarshalBigIntProto(tokensOwnedResponseItem.Index),
					Amount:      battleground_utility.MarshalBigIntProto(tokensOwnedResponseItem.Balance),
				})
			}

			responseWrapper := &orctype.OracleCommandResponse{
				CommandId: commandRequestWrapper.CommandId,
				Command: &orctype.OracleCommandResponse_GetUserFullCardCollection{
					GetUserFullCardCollection: response,
				},
			}
			commandResponses = append(commandResponses, responseWrapper)
		default:
			orc.logger.Warn("unknown command type", "commandType", reflect.TypeOf(commandRequestWrapper.GetCommand()).String())
		}
	}

	if len(commandRequests) > 0 {
		orc.logger.Debug("Sending executed command responses to Gamechain")
		err := orc.gcGateway.ProcessOracleCommandResponseBatch(commandResponses)
		if err != nil {
			return err
		}
	}

	orc.logger.Debug("Finished executing Gamechain commands")
	return nil
}

func (orc *Oracle) pollPlasmachainForEvents() (latestPlasmaBlock uint64, err error) {
	orc.logger.Info("Start polling Plasmachain")
	lastPlasmachainBlockNumber, err := orc.gcGateway.GetLastPlasmaBlockNumber()
	if err != nil {
		orc.logger.Error("failed to obtain last Plasmachain block number from Gamechain", "err", err)
		return 0, err
	}

	orc.logger.Debug("got last processed Plasmachain block number from Gamechain", "lastPlasmachainBlockNumber", lastPlasmachainBlockNumber)
	if lastPlasmachainBlockNumber == 0 {
		err = errors.New("last processed Plasmachain block number from Gamechain == 0, unable to proceed, will retry")
		return 0, err
	}

	startBlock := lastPlasmachainBlockNumber + 1
	if orc.startBlock > startBlock {
		startBlock = orc.startBlock
	}

	// TODO: limit max block range per batch
	latestBlock, err := orc.getLatestEthBlockNumber()
	if err != nil {
		orc.logger.Error("failed to obtain latest Plasmachain block number", "err", err)
		return 0, err
	}

	orc.logger.Debug("current latest Plasmachain block number", "latestBlock", latestBlock)

	if latestBlock < startBlock {
		// Wait for Plasmachain to produce a new block...
		return 0, nil
	}

	orc.logger.Info("fetching events", "startBlock", startBlock, "latestBlock", latestBlock)
	events, err := orc.fetchEvents(startBlock, latestBlock)
	if err != nil {
		orc.logger.Error("failed to fetch events from Plasmachain", "err", err)
		return 0, err
	}

	orc.logger.Debug("finished fetching events", "len(events)", len(events))

	if len(events) > 0 {
		orc.numPlasmachainEventsFetched = orc.numPlasmachainEventsFetched + uint64(len(events))
		orc.updateStatus()

		orc.logger.Debug("calling ProcessOracleEventBatch")
		if err := orc.gcGateway.ProcessOracleEventBatch(events, latestBlock); err != nil {
			return 0, err
		}
		orc.logger.Debug("finished calling ProcessOracleEventBatch")

		orc.numPlasmachainEventsSubmitted = orc.numPlasmachainEventsSubmitted + uint64(len(events))
		orc.metrics.SubmittedPlasmachainEvents(len(events))
		orc.updateStatus()
	} else {
		// If there were no events, just update the latest Plasmachain block number
		// so that we won't process same events again.
		orc.logger.Info("calling SetLastPlasmaBlockNumber")
		if err := orc.gcGateway.SetLastPlasmaBlockNumber(latestBlock); err != nil {
			orc.logger.Warn(err.Error())
			return 0, err
		}
	}

	orc.startBlock = latestBlock + 1
	return latestBlock, nil
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

	var rawEvents []*plasmachainEventInfo
	var err error

	rawEvents, err = orc.fetchTransferEvents(filterOpts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch transfer events")
	}

	sortPlasmachainEvents(rawEvents)
	events := make([]*orctype.PlasmachainEvent, len(rawEvents))
	for i, event := range rawEvents {
		events[i] = event.Event
	}

	if len(rawEvents) > 0 {
		orc.logger.Debug("fetched Plasmachain events",
			"startBlock", startBlock,
			"endBlock", endBlock,
			"eventCount", len(rawEvents),
		)
	}

	return events, nil
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

func (orc *Oracle) processSingleRawEvent(rawEvent ethtypes.Log) (eventInfo *plasmachainEventInfo, receipt *ptypes.EvmTxReceipt, err error) {
	receiptRaw, err := orc.pcGateway.client.GetEvmTxReceipt(rawEvent.TxHash.Bytes())
	if err != nil {
		orc.logger.Error(err.Error(), "txHash", rawEvent.TxHash.Hex())
		return nil, nil, err
	}

	receipt = &receiptRaw

	return &plasmachainEventInfo{
		BlockNum: rawEvent.BlockNumber,
		TxIdx:    rawEvent.TxIndex,
		Event: &orctype.PlasmachainEvent{
			EthBlock: rawEvent.BlockNumber,
		},
	}, receipt, nil
}

func (orc *Oracle) fetchTransferEvents(filterOpts *bind.FilterOpts) ([]*plasmachainEventInfo, error) {
	var err error
	numTransferEvents := 0
	numTransferWithQuantityEvents := 0
	numBatchTransferEvents := 0
	defer func(begin time.Time) {
		orc.metrics.MethodCalled(begin, "fetchTransferEvents", err)
		orc.metrics.FetchedPlasmachainEvents(numTransferEvents, "Transfer")
		orc.metrics.FetchedPlasmachainEvents(numTransferWithQuantityEvents, "TransferWithQuantity")
		orc.metrics.FetchedPlasmachainEvents(numBatchTransferEvents, "BatchTransfer")
		orc.updateStatus()
	}(time.Now())

	var chainID = orc.pcGateway.client.GetChainID()
	events := make([]*plasmachainEventInfo, 0)

	// Transfer
	transferIterator, err := orc.pcGateway.zbgCard.FilterTransfer(filterOpts, nil, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get logs for Transfer")
	}
	for {
		ok := transferIterator.Next()
		if ok {
			event := transferIterator.Event
			eventInfo, _, err := orc.processSingleRawEvent(event.Raw)
			if err != nil {
				return nil, err
			}

			fromLocal, err := loom.LocalAddressFromHexString(event.From.Hex())
			if err != nil {
				return nil, errors.Wrapf(err, "error parsing address %s", event.From.Hex())
			}

			toLocal, err := loom.LocalAddressFromHexString(event.To.Hex())
			if err != nil {
				return nil, errors.Wrapf(err, "error parsing address %s", event.To.Hex())
			}

			eventInfo.Event.Payload = &orctype.PlasmachainEvent_Transfer{
				Transfer: &orctype.PlasmachainEventTransfer{
					From: &ltypes.Address{
						ChainId: chainID,
						Local:   fromLocal,
					},
					To: &ltypes.Address{
						ChainId: chainID,
						Local:   toLocal,
					},
					TokenId: &ltypes.BigUInt{Value: common.BigUInt{Int: event.TokenId}},
				},
			}

			events = append(events, eventInfo)
			numTransferEvents++
		} else {
			err = transferIterator.Error()
			if err != nil {
				return nil, errors.Wrap(err, "Failed to get event data for Transfer")
			}
			transferIterator.Close()
			break
		}
	}

	// TransferWithQuantity
	transferWithQuantityIterator, err := orc.pcGateway.zbgCard.FilterTransferWithQuantity(filterOpts, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get logs for TransferWithQuantity")
	}
	for {
		ok := transferWithQuantityIterator.Next()
		if ok {
			event := transferWithQuantityIterator.Event
			eventInfo, _, err := orc.processSingleRawEvent(event.Raw)
			if err != nil {
				return nil, err
			}

			fromLocal, err := loom.LocalAddressFromHexString(event.From.Hex())
			if err != nil {
				return nil, errors.Wrapf(err, "error parsing address %s", event.From.Hex())
			}

			toLocal, err := loom.LocalAddressFromHexString(event.To.Hex())
			if err != nil {
				return nil, errors.Wrapf(err, "error parsing address %s", event.To.Hex())
			}

			eventInfo.Event.Payload = &orctype.PlasmachainEvent_TransferWithQuantity{
				TransferWithQuantity: &orctype.PlasmachainEventTransferWithQuantity{
					From: &ltypes.Address{
						ChainId: chainID,
						Local:   fromLocal,
					},
					To: &ltypes.Address{
						ChainId: chainID,
						Local:   toLocal,
					},
					TokenId: &ltypes.BigUInt{Value: common.BigUInt{Int: event.TokenId}},
					Amount:  &ltypes.BigUInt{Value: common.BigUInt{Int: event.Amount}},
				},
			}

			events = append(events, eventInfo)
			numTransferWithQuantityEvents++
		} else {
			err = transferWithQuantityIterator.Error()
			if err != nil {
				return nil, errors.Wrap(err, "Failed to get event data for TransferWithQuantity")
			}
			transferWithQuantityIterator.Close()
			break
		}
	}

	// BatchTransfer
	batchTransferIterator, err := orc.pcGateway.zbgCard.FilterBatchTransfer(filterOpts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get logs for BatchTransfer")
	}
	for {
		ok := batchTransferIterator.Next()
		if ok {
			event := batchTransferIterator.Event
			eventInfo, _, err := orc.processSingleRawEvent(event.Raw)
			if err != nil {
				return nil, err
			}

			fromLocal, err := loom.LocalAddressFromHexString(event.From.Hex())
			if err != nil {
				return nil, errors.Wrapf(err, "error parsing address %s", event.From.Hex())
			}

			toLocal, err := loom.LocalAddressFromHexString(event.To.Hex())
			if err != nil {
				return nil, errors.Wrapf(err, "error parsing address %s", event.To.Hex())
			}

			tokenIds := make([]*ltypes.BigUInt, len(event.TokenTypes))
			amounts := make([]*ltypes.BigUInt, len(event.Amounts))

			for index, tokenType := range event.TokenTypes {
				amount := event.Amounts[index]

				tokenIds[index] = &ltypes.BigUInt{Value: common.BigUInt{Int: tokenType}}
				amounts[index] = &ltypes.BigUInt{Value: common.BigUInt{Int: amount}}
			}

			eventInfo.Event.Payload = &orctype.PlasmachainEvent_BatchTransfer{
				BatchTransfer: &orctype.PlasmachainEventBatchTransfer{
					From: &ltypes.Address{
						ChainId: chainID,
						Local:   fromLocal,
					},
					To: &ltypes.Address{
						ChainId: chainID,
						Local:   toLocal,
					},
					TokenIds: tokenIds,
					Amounts:  amounts,
				},
			}

			events = append(events, eventInfo)
			numBatchTransferEvents++
		} else {
			err = batchTransferIterator.Error()
			if err != nil {
				return nil, errors.Wrap(err, "Failed to get event data for BatchTransfer")
			}
			batchTransferIterator.Close()
			break
		}
	}

	return events, nil
}