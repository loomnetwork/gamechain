package oracle

import (
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"time"

	orctype "github.com/loomnetwork/gamechain/types/oracle"
	"github.com/loomnetwork/go-loom"
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

func (gw *GamechainGateway) GetLastPlasmaBlockNumber() (uint64, error) {
	var req zb_calls.EmptyRequest
	var resp zb_calls.GetContractStateResponse

	if _, err := gw.contract.StaticCall("GetContractState", &req, gw.caller, &resp); err != nil {
		err = errors.Wrap(err, "failed to call GetContractState")
		gw.logger.Error(err.Error())
		return 0, err
	}
	gw.LastResponseTime = time.Now()
	return resp.State.LastPlasmachainBlockNumber, nil
}

func (gw *GamechainGateway) SetLastPlasmaBlockNumber(lastBlock uint64) error {
	req := zb_calls.SetLastPlasmaBlockNumberRequest{
		LastPlasmachainBlockNumber: lastBlock,
	}

	if _, err := gw.contract.Call("SetLastPlasmaBlockNumber", &req, gw.signer, nil); err != nil {
		err = errors.Wrap(err, "failed to call SetLastPlasmaBlockNumber")
		gw.logger.Error(err.Error())
		return err
	}
	gw.LastResponseTime = time.Now()
	return nil
}

func (gw *GamechainGateway) GetOracleCommandRequestList() ([]*orctype.OracleCommandRequest, error) {
	var req orctype.GetOracleCommandRequestListRequest
	var resp orctype.GetOracleCommandRequestListResponse

	if _, err := gw.contract.StaticCall("GetOracleCommandRequestList", &req, gw.caller, &resp); err != nil {
		err = errors.Wrap(err, "failed to call GetOracleCommandRequestList")
		gw.logger.Error(err.Error())
		return nil, err
	}
	gw.LastResponseTime = time.Now()
	return resp.CommandRequests, nil
}

func (gw *GamechainGateway) ProcessOracleCommandResponseBatch(commandResponses []*orctype.OracleCommandResponse) error {
	req := orctype.ProcessOracleCommandResponseBatchRequest{
		CommandResponses: commandResponses,
	}

	if _, err := gw.contract.Call("ProcessOracleCommandResponseBatch", &req, gw.signer, nil); err != nil {
		err = errors.Wrap(err, "failed to call ProcessOracleCommandResponseBatch")
		gw.logger.Error(err.Error())
		return err
	}
	gw.LastResponseTime = time.Now()
	return nil
}

func (gw *GamechainGateway) ProcessOracleEventBatch(events []*orctype.PlasmachainEvent, endBlock uint64, zbgCardContractAddress loom.Address) error {
	req := orctype.ProcessOracleEventBatchRequest{
		Events:                     events,
		LastPlasmachainBlockNumber: endBlock,
		ZbgCardContractAddress:     zbgCardContractAddress.MarshalPB(),
	}

	if _, err := gw.contract.Call("ProcessOracleEventBatch", &req, gw.signer, nil); err != nil {
		err = errors.Wrap(err, "failed to call ProcessOracleEventBatch")
		gw.logger.Error(err.Error())
		return err
	}
	gw.LastResponseTime = time.Now()
	return nil
}
