package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/go-loom/client"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type DAppChainGateway struct {
	Address loom.Address
	// Timestamp of the last successful response from the DAppChain
	LastResponseTime time.Time

	contract *client.Contract
	caller   loom.Address
	//logger   *loom.Logger
	signer auth.Signer
}

type LoomAuthGateway struct {
	EndPoint string
}

type LARequest struct {
	Elo int64 `json:"elo"`
}

func (g *LoomAuthGateway) ProcessEvent(eventBody []byte) error {
	var payload zb.Account
	err := proto.Unmarshal(eventBody, &payload)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("Sending event to loomauth")
	// use some other library such as gorequest?
	j, _ := json.Marshal(LARequest{
		Elo: payload.EloScore,
	})
	resp, err := http.Post(g.EndPoint, "application/json", bytes.NewBuffer(j))
	if err != nil {
		return err
	}
	log.Info("Response: ", resp.StatusCode)
	return nil
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
