package oracle

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/client"
	loomcommon "github.com/loomnetwork/go-loom/common"
	"github.com/loomnetwork/loomchain/rpc/eth"
	"github.com/pkg/errors"
)

var ErrNotImplemented = errors.New("not implememented")

type LoomchainBackend struct {
	*client.DAppChainRPCClient
}

var _ bind.ContractBackend = &LoomchainBackend{}

func NewLoomchainBackend(cli *client.DAppChainRPCClient) bind.ContractBackend {
	return &LoomchainBackend{DAppChainRPCClient: cli}
}

func (l *LoomchainBackend) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	return nil, ErrNotImplemented
}

func (l *LoomchainBackend) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return l.CallContract(ctx, call, blockNumber)
}

func (l *LoomchainBackend) PendingCodeAt(ctx context.Context, contract common.Address) ([]byte, error) {
	return nil, ErrNotImplemented
}

func (l *LoomchainBackend) PendingCallContract(ctx context.Context, call ethereum.CallMsg) ([]byte, error) {
	return nil, ErrNotImplemented
}

func (l *LoomchainBackend) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return 0, ErrNotImplemented
}

func (l *LoomchainBackend) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return big.NewInt(0), ErrNotImplemented
}

func (l *LoomchainBackend) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	return 0, ErrNotImplemented
}

func (l *LoomchainBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return ErrNotImplemented
}

func (l *LoomchainBackend) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	var ethFilter eth.EthFilter
	if query.FromBlock != nil {
		ethFilter.FromBlock = eth.BlockHeight(query.FromBlock.String())
	}
	if query.ToBlock != nil {
		ethFilter.ToBlock = eth.BlockHeight(query.ToBlock.String())
	}
	if len(query.Addresses) > 0 {
		addrs, err := commonAddressToLoomAddress(query.Addresses...)
		if err != nil {
			return nil, err
		}
		ethFilter.Addresses = addrs
	}
	if len(query.Topics) > 0 {
		ethFilter.Topics = hashToStrings(query.Topics)
	}
	filter, err := json.Marshal(&ethFilter)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(filter))
	logs, err := l.GetEvmLogs(string(filter))
	if err != nil {
		return nil, err
	}
	var tlogs []types.Log
	for _, log := range logs.EthBlockLogs {
		tlogs = append(tlogs, types.Log{
			Address:     common.BytesToAddress(log.Address),
			Topics:      byteArrayToHashes(log.Topics...),
			Data:        log.Data,
			BlockNumber: uint64(log.BlockNumber),
			TxHash:      common.BytesToHash(log.TransactionHash),
			TxIndex:     uint(log.TransactionIndex),
			BlockHash:   common.BytesToHash(log.BlockHash),
			Index:       uint(log.LogIndex),
			Removed:     log.Removed,
		})
	}
	return tlogs, nil
}

func (l *LoomchainBackend) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	return nil, ErrNotImplemented
}

func (l *LoomchainBackend) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return nil, ErrNotImplemented
}

func byteArrayToHashes(bs ...[]byte) []common.Hash {
	var hashes []common.Hash
	for _, b := range bs {
		hashes = append(hashes, common.BytesToHash(b))
	}
	return hashes
}

func hashToStrings(hss [][]common.Hash) [][]string {
	var strs [][]string
	for _, hs := range hss {
		var newhs []string
		for _, h := range hs {
			newhs = append(newhs, h.String())
		}
		if len(newhs) > 0 {
			strs = append(strs, newhs)
		}
	}
	return strs
}

func commonAddressToLoomAddress(ca ...common.Address) ([]loomcommon.LocalAddress, error) {
	var addrs []loomcommon.LocalAddress
	for _, a := range ca {
		addr, err := loom.LocalAddressFromHexString(a.Hex())
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, addr)
	}
	return addrs, nil
}
