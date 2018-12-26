package oracle

import (
	"time"

	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/go-loom/client"
	"github.com/pkg/errors"
)

type GamechainGateway struct {
	Address loom.Address
	// Timestamp of the last successful response from the DAppChain
	LastResponseTime time.Time

	contract *client.Contract
	caller   loom.Address
	logger   *loom.Logger
	signer   auth.Signer
}

func ConnectToGamechainGateway(
	loomClient *client.DAppChainRPCClient, caller loom.Address, contractName string, signer auth.Signer,
	logger *loom.Logger,
) (*GamechainGateway, error) {
	gatewayAddr, err := loomClient.Resolve(contractName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve Gateway Go contract address")
	}

	return &GamechainGateway{
		Address:          gatewayAddr,
		LastResponseTime: time.Now(),
		contract:         client.NewContract(loomClient, gatewayAddr.Local),
		caller:           caller,
		signer:           signer,
		logger:           logger,
	}, nil
}

func (gw *GamechainGateway) LastPlasmaBlockNumber() (uint64, error) {
	var req zb.GetGamechainStateRequest
	var resp zb.GetGamechainStateResponse
	if _, err := gw.contract.StaticCall("GetState", &req, gw.Address, &resp); err != nil {
		gw.logger.Error("fail to get state from plasmachain")
		return 0, err
	}
	gw.LastResponseTime = time.Now()
	return resp.State.LastPlasmachainBlockNum, nil
}

func (gw *GamechainGateway) ProcessEventBatch(events []*plasmachainEventInfo) error {
	panic("need to have ProcessEventBatch")
}
