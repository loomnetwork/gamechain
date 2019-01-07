package oracle

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/loomnetwork/gamechain/oracle/ethcontract"
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
	// client
	client *client.DAppChainRPCClient
	// simle store TODO: replace me
	simpleStore *ethcontract.SimpleStoreContract
}

type EchoEvent struct {
	Name  string
	Value *big.Int `json:"_value"`
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
		ChainID: "default",
		Local:   addr,
	}

	backend := NewLoomchainBackend(loomClient)
	simpleStore, err := ethcontract.NewSimpleStoreContract(common.HexToAddress(contractAddressHex), backend)

	return &PlasmachainGateway{
		Address:          gatewayAddr,
		LastResponseTime: time.Now(),
		contract:         client.NewContract(loomClient, gatewayAddr.Local),
		caller:           caller,
		signer:           signer,
		logger:           logger,
		client:           loomClient,
		simpleStore:      simpleStore,
	}, nil
}

func (gw *PlasmachainGateway) LastBlockNumber() (uint64, error) {
	// gw.client.GetEvm
	return 7000, nil
}

func (gw *PlasmachainGateway) FetchTokenClaimed(filterOpts *bind.FilterOpts) ([]*plasmachainEventInfo, error) {
	var err error
	// var numEvents int

	filterOpts = &bind.FilterOpts{
		Start: 1,
	}
	it, err := gw.simpleStore.FilterEcho(filterOpts, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get logs for TokenClaimed")
	}
	events := []*plasmachainEventInfo{}
	for {
		ok := it.Next()
		if ok {
			ev := it.Event
			fmt.Printf("====> EV %#v\n", ev)
			// tokenAddr, err := loom.LocalAddressFromHexString(ev..Hex())
			// if err != nil {
			// 	return nil, errors.Wrap(err, "failed to parse ERC20Received token address")
			// }
			// fromAddr, err := loom.LocalAddressFromHexString(ev.From.Hex())
			// if err != nil {
			// 	return nil, errors.Wrap(err, "failed to parse ERC20Received from address")
			// }
			events = append(events, &plasmachainEventInfo{
				BlockNum: ev.Raw.BlockNumber,
				TxIdx:    ev.Raw.TxIndex,
				Event:    &PlasmachainEventData{
					// EthBlock: ev.Raw.BlockNumber,
					// Payload: &MainnetDepositEvent{
					// 	Deposit: &MainnetTokenDeposited{
					// 		TokenKind:     TokenKind_ERC20,
					// 		TokenContract: loom.Address{ChainID: "eth", Local: tokenAddr}.MarshalPB(),
					// 		TokenOwner:    loom.Address{ChainID: "eth", Local: fromAddr}.MarshalPB(),
					// 		TokenAmount:   &ltypes.BigUInt{Value: *loom.NewBigUInt(ev.Amount)},
					// 	},
					// },
				},
			})
		} else {
			err = it.Error()
			if err != nil {
				return nil, errors.Wrap(err, "Failed to get event data for TokenClaimed")
			}
			it.Close()
			break
		}
	}
	// numEvents = len(events)
	return events, nil
}

// func (gw *PlasmachainGateway) FetchOpenPack(filterOpts *bind.FilterOpts) ([]*plasmachainEventInfo, error) {
// 	fmt.Printf("--> FetchOpenPack\n")
// 	filter := `{"fromBlock": "1"}`
// 	logs, err := gw.client.GetEvmLogs(filter)
// 	if err != nil {
// 		return nil, err
// 	}
// 	eventSignature := []byte("Echo(string,uint256)")
// 	hash := crypto.Keccak256Hash(eventSignature)
// 	fmt.Printf("-->want topic hash: %s\n\n", hash.Hex())
// 	for _, log := range logs.EthBlockLogs {
// 		// find signature
// 		// fmt.Printf("-->topic2: %s\n\n", log.Topics[0])
// 		if strings.Compare(string(log.Topics[0]), hash.Hex()) != 0 {
// 			continue
// 		}
// 		fmt.Printf("-->num: %+v, %x\n\n", log.BlockNumber, log.TransactionHash)
// 		var event EchoEvent
// 		err := simpleStoreContractABI.Unpack(&event, "Echo", log.Data)
// 		if err == nil {
// 			fmt.Printf("-->topic: %s\n", log.Topics)
// 			fmt.Printf("-->event: %+v\n", event)
// 			fmt.Printf("-->topic: %s\n\n", log.Topics[1])
// 			data, _ := hexutil.Decode(string(log.Topics[1]))
// 			fmt.Printf("%s", data)
// 			// tx, err := gw.client.GetEvmTxReceipt(log.TransactionHash)
// 			// if err != nil {
// 			// 	return nil, err
// 			// }
// 			// // fmt.Printf("-->tx: %+v\n\n", tx.Value)
// 			// for _, data := range tx.Logs {
// 			// 	fmt.Printf("-->txhash: %s\n", data.Caller.Local.String())
// 			// 	fmt.Printf("-->data topic: %+v\n", data.Topics)
// 			// 	fmt.Printf("-->body: %x\n", data.EncodedBody)
// 			// 	fmt.Printf("-->data len topic: %+v\n\n", len(data.Topics))
// 			// }
// 		}

// 	}
// 	return nil, nil
// }
