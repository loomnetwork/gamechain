package oracle

import (
	"context"
	"fmt"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/loomnetwork/go-loom/client"
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

	filter := `{"fromBlock": "1", "topics": ["0xcaf432cb38a3a6f6c9bdd5b57f1a5388e0f452215b40290524727e0a7a523da1"]}`
	logs, err := l.GetEvmLogs(filter)
	if err != nil {
		return nil, err
	}
	var tlogs []types.Log
	for _, log := range logs.EthBlockLogs {
		fmt.Printf("----> Log: %s\n", log.Topics)
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

func byteArrayToHashes(bs ...[]byte) []common.Hash {
	var hashes []common.Hash
	for _, b := range bs {
		hashes = append(hashes, common.BytesToHash(b))
	}
	return hashes
}

func (l *LoomchainBackend) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	return nil, ErrNotImplemented
}

func (l *LoomchainBackend) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return nil, ErrNotImplemented
}
