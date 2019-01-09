package oracle

import (
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/loomnetwork/gamechain/oracle/ethcontract"
	orctype "github.com/loomnetwork/gamechain/types/oracle"
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/go-loom/client"
	ltypes "github.com/loomnetwork/go-loom/types"
	"github.com/pkg/errors"
)

type (
	PlasmachainEvent              = orctype.PlasmachainEvent
	PlasmachainGeneratedCardEvent = orctype.PlasmachainEvent_Card
	PlasmachainGeneratedCard      = orctype.PlasmachainGeneratedCard
)

type plasmachainEventInfo struct {
	BlockNum uint64
	TxIdx    uint
	Event    *PlasmachainEvent
}

type PlasmachainGateway struct {
	Address loom.Address
	// Timestamp of the last successful response from the DAppChain
	LastResponseTime time.Time

	contract *client.Contract
	caller   loom.Address
	logger   *loom.Logger
	signer   auth.Signer
	// client
	client *client.DAppChainRPCClient
	// cardfaucet
	cardFaucet *ethcontract.CardFaucet
}

func ConnectToPlasmachainGateway(
	loomClient *client.DAppChainRPCClient, caller loom.Address, contractAddressHex string, signer auth.Signer,
	logger *loom.Logger,
) (*PlasmachainGateway, error) {
	addr, err := loom.LocalAddressFromHexString(contractAddressHex)
	if err != nil {
		return nil, err
	}
	gatewayAddr := loom.Address{
		ChainID: loomClient.GetChainID(),
		Local:   addr,
	}

	backend := NewLoomchainBackend(loomClient)
	cardFaucet, err := ethcontract.NewCardFaucet(common.HexToAddress(contractAddressHex), backend)
	if err != nil {
		return nil, err
	}

	return &PlasmachainGateway{
		Address:          gatewayAddr,
		LastResponseTime: time.Now(),
		contract:         client.NewContract(loomClient, gatewayAddr.Local),
		caller:           caller,
		signer:           signer,
		logger:           logger,
		client:           loomClient,
		cardFaucet:       cardFaucet,
	}, nil
}

func (gw *PlasmachainGateway) LastBlockNumber() (uint64, error) {
	block, err := gw.client.GetEvmBlockByNumber("latest", false)
	if err != nil {
		return 0, err
	}
	return uint64(block.Number), nil
}

func (gw *PlasmachainGateway) FetchGeneratedCard(filterOpts *bind.FilterOpts) ([]*plasmachainEventInfo, error) {
	var err error
	// var numEvents int
	it, err := gw.cardFaucet.FilterGeneratedCard(filterOpts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get logs for Transfer")
	}
	var chainID = gw.client.GetChainID()
	events := []*plasmachainEventInfo{}
	for {
		ok := it.Next()
		if ok {
			ev := it.Event
			receipt, err := gw.client.GetEvmTxReceipt(ev.Raw.TxHash.Bytes())
			if err != nil {
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
							CardId:   &ltypes.BigUInt{Value: *loom.NewBigUInt(ev.CardId)},
							Amount:   &ltypes.BigUInt{Value: *loom.NewBigUIntFromInt(1)},
							Contract: contractAddr,
						},
					},
				},
			})
		} else {
			err = it.Error()
			if err != nil {
				return nil, errors.Wrap(err, "Failed to get event data for Transfer")
			}
			it.Close()
			break
		}
	}
	// numEvents = len(events)
	return events, nil
}
