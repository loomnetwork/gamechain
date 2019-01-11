package oracle

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/loomnetwork/gamechain/oracle/ethcontract"
	orctype "github.com/loomnetwork/gamechain/types/oracle"
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/go-loom/client"
)

type (
	PlasmachainEvent              = orctype.PlasmachainEvent
	PlasmachainGeneratedCardEvent = orctype.PlasmachainEvent_Card
	PlasmachainGeneratedCard      = orctype.PlasmachainGeneratedCard
	ProcessEventBatchRequest      = orctype.ProcessEventBatchRequest
)

type plasmachainEventInfo struct {
	BlockNum uint64
	TxIdx    uint
	Event    *orctype.PlasmachainEvent
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
	gw.LastResponseTime = time.Now()
	return uint64(block.Number), nil
}
