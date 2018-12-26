package oracle

import (
	"time"

	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/go-loom/client"
	"github.com/pkg/errors"
)

type PlasmachainGateway struct {
	Address loom.Address
	// Timestamp of the last successful response from the DAppChain
	LastResponseTime time.Time

	contract *client.Contract
	caller   loom.Address
	logger   *loom.Logger
	signer   auth.Signer
}

func ConnectToPlasmachainGateway(
	loomClient *client.DAppChainRPCClient, caller loom.Address, contractName string, signer auth.Signer,
	logger *loom.Logger,
) (*PlasmachainGateway, error) {
	gatewayAddr, err := loomClient.Resolve(contractName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve Gateway Go contract address")
	}

	return &PlasmachainGateway{
		Address:          gatewayAddr,
		LastResponseTime: time.Now(),
		contract:         client.NewContract(loomClient, gatewayAddr.Local),
		caller:           caller,
		signer:           signer,
		logger:           logger,
	}, nil
}

func (gw *PlasmachainGateway) LastBlockNumber() (uint64, error) {
	// panic("need to have GetState")
	// FIX ME:
	return 1, nil
}
