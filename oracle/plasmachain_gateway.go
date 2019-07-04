package oracle

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"time"

	"github.com/loomnetwork/gamechain/oracle/ethcontract"
	orctype "github.com/loomnetwork/gamechain/types/oracle"
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/go-loom/client"
)

type plasmachainEventInfo struct {
	BlockNum uint64
	TxIdx    uint
	Event    *orctype.PlasmachainEvent
}

type tokensOwnedResponseItem struct{
	Index  *big.Int
	Balance *big.Int
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
	// zbgCard
	zbgCard *ethcontract.ZBGCard
}

func ConnectToPlasmachainGateway(
	loomClient *client.DAppChainRPCClient,
	caller loom.Address,
	zbgCardContractAddress loom.Address,
	signer auth.Signer,
	logger *loom.Logger,
) (*PlasmachainGateway, error) {
	gatewayAddr := zbgCardContractAddress

	backend := NewLoomchainBackend(loomClient, signer)
	zbgCard, err := ethcontract.NewZBGCard(common.HexToAddress(zbgCardContractAddress.Local.Hex()), backend)
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
		zbgCard:          zbgCard,
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

func (gw *PlasmachainGateway) GetTokensOwned(owner loom.LocalAddress) ([]tokensOwnedResponseItem, error) {
	ownerCommonAddress := common.BytesToAddress(owner)
	opts := &bind.CallOpts{
	}

	tokensOwnedRaw, err := gw.zbgCard.TokensOwned(opts, ownerCommonAddress)
	if err != nil {
		return nil, err
	}

	tokensOwned := make([]tokensOwnedResponseItem, len(tokensOwnedRaw.Balances))
	for i := range tokensOwnedRaw.Balances {
		tokensOwned[i] = tokensOwnedResponseItem {
			Index: tokensOwnedRaw.Indexes[i],
			Balance: tokensOwnedRaw.Balances[i],
		}
	}

	gw.LastResponseTime = time.Now()
	return tokensOwned, nil
}
