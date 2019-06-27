package oracle

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/go-loom/client"
	loomcommon "github.com/loomnetwork/go-loom/common"
	"github.com/loomnetwork/loomchain/rpc/eth"
	"github.com/pkg/errors"
	"math/big"
)

var ErrNotImplemented = errors.New("not implememented")

type LoomchainBackend struct {
	*client.DAppChainRPCClient
	signer auth.Signer
}

var _ bind.ContractBackend = &LoomchainBackend{}

func NewLoomchainBackend(cli *client.DAppChainRPCClient, signer auth.Signer) bind.ContractBackend {
	return &LoomchainBackend{
		DAppChainRPCClient: cli,
		signer: signer,
	}
}

func (l *LoomchainBackend) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	return nil, ErrNotImplemented
}

func (l *LoomchainBackend) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {

	contractAddress, err := commonAddressToLoomAddress(*call.To, l.GetChainID())
	if err != nil {
		return nil, err
	}
	callerAddres := loom.Address{
		ChainID: l.GetChainID(),
		Local:   loom.LocalAddressFromPublicKey(l.signer.PublicKey()),
	}

	evmContract := client.NewEvmContract(l.DAppChainRPCClient, contractAddress.Local)

	bytes, err := evmContract.StaticCall(call.Data, callerAddres)
	if err != nil {
		return nil, err
	}

	return bytes, nil
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
		addrs, err := commonAddressesToLoomAddresses(query.Addresses...)
		if err != nil {
			return nil, err
		}
		ethFilter.Addresses = addrs
	}
	if len(query.Topics) > 0 {
		ethFilter.Topics = hashToStrings(query.Topics)
	}

	// If eth.EthFilter is used, we get strange results with missing first topic
	jsonFilter := eth.JsonFilter{}
	jsonFilter.FromBlock = ethFilter.FromBlock
	jsonFilter.ToBlock = ethFilter.ToBlock
	jsonFilter.Address = query.Addresses[0]
	jsonFilter.Topics = []interface{}{query.Topics[0][0].String()}

	filter, err := json.Marshal(&jsonFilter)
	if err != nil {
		return nil, err
	}

	logs, err := l.GetEvmLogs(string(filter))
	if err != nil {
		return nil, err
	}

	var tlogs []types.Log
	for _, log := range logs.EthBlockLogs {
		topicBytes := make([][]byte, len(log.Topics))
		for i, topicHexStringBytes := range log.Topics {
			topicBytes[i], err = hexutil.Decode(string(topicHexStringBytes))
			if err != nil {
				return nil, errors.Wrap(err, "failed to decode topics")
			}
		}
		topicHashes := byteArrayToHashes(topicBytes...)

		tlogs = append(tlogs, types.Log{
			Address:     common.BytesToAddress(log.Address),
			Topics:      topicHashes,
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

func commonAddressToLoomAddress(ca common.Address, chainId string) (*loom.Address, error) {
	localAddress, err := loom.LocalAddressFromHexString(ca.Hex())
	if err != nil {
		return nil, err
	}

	return &loom.Address{
		ChainID: chainId,
		Local:   localAddress,
	}, nil
}

func commonAddressesToLoomAddresses(ca ...common.Address) ([]loomcommon.LocalAddress, error) {
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
