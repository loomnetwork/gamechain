// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ethcontract

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// ZBGCardABI is the input ABI used to generate the binding from.
const ZBGCardABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"_interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"implementsERC721\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"toggleToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_tokenIds\",\"type\":\"uint256[]\"},{\"name\":\"_amounts\",\"type\":\"uint256[]\"}],\"name\":\"batchTransferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"InterfaceId_ERC165\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"tokensOwned\",\"outputs\":[{\"name\":\"indexes\",\"type\":\"uint256[]\"},{\"name\":\"balances\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_tokenIds\",\"type\":\"uint256[]\"},{\"name\":\"_amounts\",\"type\":\"uint256[]\"},{\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"safeBatchTransferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"tokenOfOwnerByIndex\",\"outputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"exists\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"tokenByIndex\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numValidators\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"tokenIdToDNA\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_address\",\"type\":\"address\"},{\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"balanceOfCoin\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"checkValidator\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"implementsERC721X\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_validator\",\"type\":\"address\"},{\"name\":\"_v\",\"type\":\"uint8[]\"},{\"name\":\"_r\",\"type\":\"bytes32[]\"},{\"name\":\"_s\",\"type\":\"bytes32[]\"}],\"name\":\"addValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_operator\",\"type\":\"address\"},{\"name\":\"_approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"nonce\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_validator\",\"type\":\"address\"},{\"name\":\"_v\",\"type\":\"uint8[]\"},{\"name\":\"_r\",\"type\":\"bytes32[]\"},{\"name\":\"_s\",\"type\":\"bytes32[]\"}],\"name\":\"removeValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowedTokens\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"name\":\"isOperator\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"name\":\"_amount\",\"type\":\"uint256\"},{\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_gateway\",\"type\":\"address\"},{\"name\":\"_validators\",\"type\":\"address[]\"},{\"name\":\"_threshold_num\",\"type\":\"uint8\"},{\"name\":\"_threshold_denom\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"supply\",\"type\":\"uint256\"}],\"name\":\"TokenMinted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"claimer\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TokenClaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"isApproved\",\"type\":\"bool\"}],\"name\":\"FaucetToggled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TransferWithQuantity\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"validator\",\"type\":\"address\"}],\"name\":\"AddedValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"validator\",\"type\":\"address\"}],\"name\":\"RemovedValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"}],\"name\":\"OwnershipRenounced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokenTypes\",\"type\":\"uint256[]\"},{\"indexed\":false,\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"name\":\"BatchTransfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_approved\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_operator\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"depositToGateway\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"individualSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"name\":\"_sig\",\"type\":\"bytes\"}],\"name\":\"claimTokenNFT\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_tokenIds\",\"type\":\"uint256[]\"},{\"name\":\"_tokenDNAs\",\"type\":\"uint256[]\"}],\"name\":\"batchMintToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"name\":\"_tokenDNA\",\"type\":\"uint256\"}],\"name\":\"mintToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"name\":\"_amount\",\"type\":\"uint256\"},{\"name\":\"_sig\",\"type\":\"bytes\"}],\"name\":\"claimToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_faucet\",\"type\":\"address\"}],\"name\":\"enableFaucet\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_faucet\",\"type\":\"address\"}],\"name\":\"disableFaucet\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"getTokenDetailsById\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"name\":\"tokenUri\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"isLimitedEdition\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom1\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// ZBGCard is an auto generated Go binding around an Ethereum contract.
type ZBGCard struct {
	ZBGCardCaller     // Read-only binding to the contract
	ZBGCardTransactor // Write-only binding to the contract
	ZBGCardFilterer   // Log filterer for contract events
}

// ZBGCardCaller is an auto generated read-only Go binding around an Ethereum contract.
type ZBGCardCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZBGCardTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ZBGCardTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZBGCardFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ZBGCardFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZBGCardSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ZBGCardSession struct {
	Contract     *ZBGCard          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ZBGCardCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ZBGCardCallerSession struct {
	Contract *ZBGCardCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// ZBGCardTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ZBGCardTransactorSession struct {
	Contract     *ZBGCardTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ZBGCardRaw is an auto generated low-level Go binding around an Ethereum contract.
type ZBGCardRaw struct {
	Contract *ZBGCard // Generic contract binding to access the raw methods on
}

// ZBGCardCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ZBGCardCallerRaw struct {
	Contract *ZBGCardCaller // Generic read-only contract binding to access the raw methods on
}

// ZBGCardTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ZBGCardTransactorRaw struct {
	Contract *ZBGCardTransactor // Generic write-only contract binding to access the raw methods on
}

// NewZBGCard creates a new instance of ZBGCard, bound to a specific deployed contract.
func NewZBGCard(address common.Address, backend bind.ContractBackend) (*ZBGCard, error) {
	contract, err := bindZBGCard(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ZBGCard{ZBGCardCaller: ZBGCardCaller{contract: contract}, ZBGCardTransactor: ZBGCardTransactor{contract: contract}, ZBGCardFilterer: ZBGCardFilterer{contract: contract}}, nil
}

// NewZBGCardCaller creates a new read-only instance of ZBGCard, bound to a specific deployed contract.
func NewZBGCardCaller(address common.Address, caller bind.ContractCaller) (*ZBGCardCaller, error) {
	contract, err := bindZBGCard(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ZBGCardCaller{contract: contract}, nil
}

// NewZBGCardTransactor creates a new write-only instance of ZBGCard, bound to a specific deployed contract.
func NewZBGCardTransactor(address common.Address, transactor bind.ContractTransactor) (*ZBGCardTransactor, error) {
	contract, err := bindZBGCard(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ZBGCardTransactor{contract: contract}, nil
}

// NewZBGCardFilterer creates a new log filterer instance of ZBGCard, bound to a specific deployed contract.
func NewZBGCardFilterer(address common.Address, filterer bind.ContractFilterer) (*ZBGCardFilterer, error) {
	contract, err := bindZBGCard(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ZBGCardFilterer{contract: contract}, nil
}

// bindZBGCard binds a generic wrapper to an already deployed contract.
func bindZBGCard(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ZBGCardABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZBGCard *ZBGCardRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ZBGCard.Contract.ZBGCardCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZBGCard *ZBGCardRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZBGCard.Contract.ZBGCardTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZBGCard *ZBGCardRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZBGCard.Contract.ZBGCardTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZBGCard *ZBGCardCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ZBGCard.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZBGCard *ZBGCardTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZBGCard.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZBGCard *ZBGCardTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZBGCard.Contract.contract.Transact(opts, method, params...)
}

// InterfaceIdERC165 is a free data retrieval call binding the contract method 0x19fa8f50.
//
// Solidity: function InterfaceId_ERC165() constant returns(bytes4)
func (_ZBGCard *ZBGCardCaller) InterfaceIdERC165(opts *bind.CallOpts) ([4]byte, error) {
	var (
		ret0 = new([4]byte)
	)
	out := ret0
	err := _ZBGCard.contract.Call(opts, out, "InterfaceId_ERC165")
	return *ret0, err
}

// InterfaceIdERC165 is a free data retrieval call binding the contract method 0x19fa8f50.
//
// Solidity: function InterfaceId_ERC165() constant returns(bytes4)
func (_ZBGCard *ZBGCardSession) InterfaceIdERC165() ([4]byte, error) {
	return _ZBGCard.Contract.InterfaceIdERC165(&_ZBGCard.CallOpts)
}

// InterfaceIdERC165 is a free data retrieval call binding the contract method 0x19fa8f50.
//
// Solidity: function InterfaceId_ERC165() constant returns(bytes4)
func (_ZBGCard *ZBGCardCallerSession) InterfaceIdERC165() ([4]byte, error) {
	return _ZBGCard.Contract.InterfaceIdERC165(&_ZBGCard.CallOpts)
}

// AllowedTokens is a free data retrieval call binding the contract method 0xe744092e.
//
// Solidity: function allowedTokens( address) constant returns(bool)
func (_ZBGCard *ZBGCardCaller) AllowedTokens(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ZBGCard.contract.Call(opts, out, "allowedTokens", arg0)
	return *ret0, err
}

// AllowedTokens is a free data retrieval call binding the contract method 0xe744092e.
//
// Solidity: function allowedTokens( address) constant returns(bool)
func (_ZBGCard *ZBGCardSession) AllowedTokens(arg0 common.Address) (bool, error) {
	return _ZBGCard.Contract.AllowedTokens(&_ZBGCard.CallOpts, arg0)
}

// AllowedTokens is a free data retrieval call binding the contract method 0xe744092e.
//
// Solidity: function allowedTokens( address) constant returns(bool)
func (_ZBGCard *ZBGCardCallerSession) AllowedTokens(arg0 common.Address) (bool, error) {
	return _ZBGCard.Contract.AllowedTokens(&_ZBGCard.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_ZBGCard *ZBGCardCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ZBGCard.contract.Call(opts, out, "balanceOf", _owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_ZBGCard *ZBGCardSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _ZBGCard.Contract.BalanceOf(&_ZBGCard.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_ZBGCard *ZBGCardCallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _ZBGCard.Contract.BalanceOf(&_ZBGCard.CallOpts, _owner)
}

// BalanceOfCoin is a free data retrieval call binding the contract method 0x70af0e54.
//
// Solidity: function balanceOfCoin(_address address, _tokenId uint256) constant returns(uint256)
func (_ZBGCard *ZBGCardCaller) BalanceOfCoin(opts *bind.CallOpts, _address common.Address, _tokenId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ZBGCard.contract.Call(opts, out, "balanceOfCoin", _address, _tokenId)
	return *ret0, err
}

// BalanceOfCoin is a free data retrieval call binding the contract method 0x70af0e54.
//
// Solidity: function balanceOfCoin(_address address, _tokenId uint256) constant returns(uint256)
func (_ZBGCard *ZBGCardSession) BalanceOfCoin(_address common.Address, _tokenId *big.Int) (*big.Int, error) {
	return _ZBGCard.Contract.BalanceOfCoin(&_ZBGCard.CallOpts, _address, _tokenId)
}

// BalanceOfCoin is a free data retrieval call binding the contract method 0x70af0e54.
//
// Solidity: function balanceOfCoin(_address address, _tokenId uint256) constant returns(uint256)
func (_ZBGCard *ZBGCardCallerSession) BalanceOfCoin(_address common.Address, _tokenId *big.Int) (*big.Int, error) {
	return _ZBGCard.Contract.BalanceOfCoin(&_ZBGCard.CallOpts, _address, _tokenId)
}

// CheckValidator is a free data retrieval call binding the contract method 0x797327ae.
//
// Solidity: function checkValidator(_address address) constant returns(bool)
func (_ZBGCard *ZBGCardCaller) CheckValidator(opts *bind.CallOpts, _address common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ZBGCard.contract.Call(opts, out, "checkValidator", _address)
	return *ret0, err
}

// CheckValidator is a free data retrieval call binding the contract method 0x797327ae.
//
// Solidity: function checkValidator(_address address) constant returns(bool)
func (_ZBGCard *ZBGCardSession) CheckValidator(_address common.Address) (bool, error) {
	return _ZBGCard.Contract.CheckValidator(&_ZBGCard.CallOpts, _address)
}

// CheckValidator is a free data retrieval call binding the contract method 0x797327ae.
//
// Solidity: function checkValidator(_address address) constant returns(bool)
func (_ZBGCard *ZBGCardCallerSession) CheckValidator(_address common.Address) (bool, error) {
	return _ZBGCard.Contract.CheckValidator(&_ZBGCard.CallOpts, _address)
}

// Exists is a free data retrieval call binding the contract method 0x4f558e79.
//
// Solidity: function exists(_tokenId uint256) constant returns(bool)
func (_ZBGCard *ZBGCardCaller) Exists(opts *bind.CallOpts, _tokenId *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ZBGCard.contract.Call(opts, out, "exists", _tokenId)
	return *ret0, err
}

// Exists is a free data retrieval call binding the contract method 0x4f558e79.
//
// Solidity: function exists(_tokenId uint256) constant returns(bool)
func (_ZBGCard *ZBGCardSession) Exists(_tokenId *big.Int) (bool, error) {
	return _ZBGCard.Contract.Exists(&_ZBGCard.CallOpts, _tokenId)
}

// Exists is a free data retrieval call binding the contract method 0x4f558e79.
//
// Solidity: function exists(_tokenId uint256) constant returns(bool)
func (_ZBGCard *ZBGCardCallerSession) Exists(_tokenId *big.Int) (bool, error) {
	return _ZBGCard.Contract.Exists(&_ZBGCard.CallOpts, _tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(_tokenId uint256) constant returns(address)
func (_ZBGCard *ZBGCardCaller) GetApproved(opts *bind.CallOpts, _tokenId *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ZBGCard.contract.Call(opts, out, "getApproved", _tokenId)
	return *ret0, err
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(_tokenId uint256) constant returns(address)
func (_ZBGCard *ZBGCardSession) GetApproved(_tokenId *big.Int) (common.Address, error) {
	return _ZBGCard.Contract.GetApproved(&_ZBGCard.CallOpts, _tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(_tokenId uint256) constant returns(address)
func (_ZBGCard *ZBGCardCallerSession) GetApproved(_tokenId *big.Int) (common.Address, error) {
	return _ZBGCard.Contract.GetApproved(&_ZBGCard.CallOpts, _tokenId)
}

// GetTokenDetailsById is a free data retrieval call binding the contract method 0x8d43f0bf.
//
// Solidity: function getTokenDetailsById(_tokenId uint256) constant returns(uint256, uint256, uint256, uint256)
func (_ZBGCard *ZBGCardCaller) GetTokenDetailsById(opts *bind.CallOpts, _tokenId *big.Int) (*big.Int, *big.Int, *big.Int, *big.Int, error) {
	var (
		ret0 = new(*big.Int)
		ret1 = new(*big.Int)
		ret2 = new(*big.Int)
		ret3 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
	}
	err := _ZBGCard.contract.Call(opts, out, "getTokenDetailsById", _tokenId)
	return *ret0, *ret1, *ret2, *ret3, err
}

// GetTokenDetailsById is a free data retrieval call binding the contract method 0x8d43f0bf.
//
// Solidity: function getTokenDetailsById(_tokenId uint256) constant returns(uint256, uint256, uint256, uint256)
func (_ZBGCard *ZBGCardSession) GetTokenDetailsById(_tokenId *big.Int) (*big.Int, *big.Int, *big.Int, *big.Int, error) {
	return _ZBGCard.Contract.GetTokenDetailsById(&_ZBGCard.CallOpts, _tokenId)
}

// GetTokenDetailsById is a free data retrieval call binding the contract method 0x8d43f0bf.
//
// Solidity: function getTokenDetailsById(_tokenId uint256) constant returns(uint256, uint256, uint256, uint256)
func (_ZBGCard *ZBGCardCallerSession) GetTokenDetailsById(_tokenId *big.Int) (*big.Int, *big.Int, *big.Int, *big.Int, error) {
	return _ZBGCard.Contract.GetTokenDetailsById(&_ZBGCard.CallOpts, _tokenId)
}

// ImplementsERC721 is a free data retrieval call binding the contract method 0x1051db34.
//
// Solidity: function implementsERC721() constant returns(bool)
func (_ZBGCard *ZBGCardCaller) ImplementsERC721(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ZBGCard.contract.Call(opts, out, "implementsERC721")
	return *ret0, err
}

// ImplementsERC721 is a free data retrieval call binding the contract method 0x1051db34.
//
// Solidity: function implementsERC721() constant returns(bool)
func (_ZBGCard *ZBGCardSession) ImplementsERC721() (bool, error) {
	return _ZBGCard.Contract.ImplementsERC721(&_ZBGCard.CallOpts)
}

// ImplementsERC721 is a free data retrieval call binding the contract method 0x1051db34.
//
// Solidity: function implementsERC721() constant returns(bool)
func (_ZBGCard *ZBGCardCallerSession) ImplementsERC721() (bool, error) {
	return _ZBGCard.Contract.ImplementsERC721(&_ZBGCard.CallOpts)
}

// ImplementsERC721X is a free data retrieval call binding the contract method 0x7fb42a36.
//
// Solidity: function implementsERC721X() constant returns(bool)
func (_ZBGCard *ZBGCardCaller) ImplementsERC721X(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ZBGCard.contract.Call(opts, out, "implementsERC721X")
	return *ret0, err
}

// ImplementsERC721X is a free data retrieval call binding the contract method 0x7fb42a36.
//
// Solidity: function implementsERC721X() constant returns(bool)
func (_ZBGCard *ZBGCardSession) ImplementsERC721X() (bool, error) {
	return _ZBGCard.Contract.ImplementsERC721X(&_ZBGCard.CallOpts)
}

// ImplementsERC721X is a free data retrieval call binding the contract method 0x7fb42a36.
//
// Solidity: function implementsERC721X() constant returns(bool)
func (_ZBGCard *ZBGCardCallerSession) ImplementsERC721X() (bool, error) {
	return _ZBGCard.Contract.ImplementsERC721X(&_ZBGCard.CallOpts)
}

// IndividualSupply is a free data retrieval call binding the contract method 0x17be1933.
//
// Solidity: function individualSupply(_tokenId uint256) constant returns(uint256)
func (_ZBGCard *ZBGCardCaller) IndividualSupply(opts *bind.CallOpts, _tokenId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ZBGCard.contract.Call(opts, out, "individualSupply", _tokenId)
	return *ret0, err
}

// IndividualSupply is a free data retrieval call binding the contract method 0x17be1933.
//
// Solidity: function individualSupply(_tokenId uint256) constant returns(uint256)
func (_ZBGCard *ZBGCardSession) IndividualSupply(_tokenId *big.Int) (*big.Int, error) {
	return _ZBGCard.Contract.IndividualSupply(&_ZBGCard.CallOpts, _tokenId)
}

// IndividualSupply is a free data retrieval call binding the contract method 0x17be1933.
//
// Solidity: function individualSupply(_tokenId uint256) constant returns(uint256)
func (_ZBGCard *ZBGCardCallerSession) IndividualSupply(_tokenId *big.Int) (*big.Int, error) {
	return _ZBGCard.Contract.IndividualSupply(&_ZBGCard.CallOpts, _tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(_owner address, _operator address) constant returns(isOperator bool)
func (_ZBGCard *ZBGCardCaller) IsApprovedForAll(opts *bind.CallOpts, _owner common.Address, _operator common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ZBGCard.contract.Call(opts, out, "isApprovedForAll", _owner, _operator)
	return *ret0, err
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(_owner address, _operator address) constant returns(isOperator bool)
func (_ZBGCard *ZBGCardSession) IsApprovedForAll(_owner common.Address, _operator common.Address) (bool, error) {
	return _ZBGCard.Contract.IsApprovedForAll(&_ZBGCard.CallOpts, _owner, _operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(_owner address, _operator address) constant returns(isOperator bool)
func (_ZBGCard *ZBGCardCallerSession) IsApprovedForAll(_owner common.Address, _operator common.Address) (bool, error) {
	return _ZBGCard.Contract.IsApprovedForAll(&_ZBGCard.CallOpts, _owner, _operator)
}

// IsLimitedEdition is a free data retrieval call binding the contract method 0xe3981429.
//
// Solidity: function isLimitedEdition(_tokenId uint256) constant returns(bool)
func (_ZBGCard *ZBGCardCaller) IsLimitedEdition(opts *bind.CallOpts, _tokenId *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ZBGCard.contract.Call(opts, out, "isLimitedEdition", _tokenId)
	return *ret0, err
}

// IsLimitedEdition is a free data retrieval call binding the contract method 0xe3981429.
//
// Solidity: function isLimitedEdition(_tokenId uint256) constant returns(bool)
func (_ZBGCard *ZBGCardSession) IsLimitedEdition(_tokenId *big.Int) (bool, error) {
	return _ZBGCard.Contract.IsLimitedEdition(&_ZBGCard.CallOpts, _tokenId)
}

// IsLimitedEdition is a free data retrieval call binding the contract method 0xe3981429.
//
// Solidity: function isLimitedEdition(_tokenId uint256) constant returns(bool)
func (_ZBGCard *ZBGCardCallerSession) IsLimitedEdition(_tokenId *big.Int) (bool, error) {
	return _ZBGCard.Contract.IsLimitedEdition(&_ZBGCard.CallOpts, _tokenId)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_ZBGCard *ZBGCardCaller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _ZBGCard.contract.Call(opts, out, "name")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_ZBGCard *ZBGCardSession) Name() (string, error) {
	return _ZBGCard.Contract.Name(&_ZBGCard.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_ZBGCard *ZBGCardCallerSession) Name() (string, error) {
	return _ZBGCard.Contract.Name(&_ZBGCard.CallOpts)
}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() constant returns(uint256)
func (_ZBGCard *ZBGCardCaller) Nonce(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ZBGCard.contract.Call(opts, out, "nonce")
	return *ret0, err
}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() constant returns(uint256)
func (_ZBGCard *ZBGCardSession) Nonce() (*big.Int, error) {
	return _ZBGCard.Contract.Nonce(&_ZBGCard.CallOpts)
}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() constant returns(uint256)
func (_ZBGCard *ZBGCardCallerSession) Nonce() (*big.Int, error) {
	return _ZBGCard.Contract.Nonce(&_ZBGCard.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces( address) constant returns(uint256)
func (_ZBGCard *ZBGCardCaller) Nonces(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ZBGCard.contract.Call(opts, out, "nonces", arg0)
	return *ret0, err
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces( address) constant returns(uint256)
func (_ZBGCard *ZBGCardSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _ZBGCard.Contract.Nonces(&_ZBGCard.CallOpts, arg0)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces( address) constant returns(uint256)
func (_ZBGCard *ZBGCardCallerSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _ZBGCard.Contract.Nonces(&_ZBGCard.CallOpts, arg0)
}

// NumValidators is a free data retrieval call binding the contract method 0x5d593f8d.
//
// Solidity: function numValidators() constant returns(uint256)
func (_ZBGCard *ZBGCardCaller) NumValidators(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ZBGCard.contract.Call(opts, out, "numValidators")
	return *ret0, err
}

// NumValidators is a free data retrieval call binding the contract method 0x5d593f8d.
//
// Solidity: function numValidators() constant returns(uint256)
func (_ZBGCard *ZBGCardSession) NumValidators() (*big.Int, error) {
	return _ZBGCard.Contract.NumValidators(&_ZBGCard.CallOpts)
}

// NumValidators is a free data retrieval call binding the contract method 0x5d593f8d.
//
// Solidity: function numValidators() constant returns(uint256)
func (_ZBGCard *ZBGCardCallerSession) NumValidators() (*big.Int, error) {
	return _ZBGCard.Contract.NumValidators(&_ZBGCard.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_ZBGCard *ZBGCardCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ZBGCard.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_ZBGCard *ZBGCardSession) Owner() (common.Address, error) {
	return _ZBGCard.Contract.Owner(&_ZBGCard.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_ZBGCard *ZBGCardCallerSession) Owner() (common.Address, error) {
	return _ZBGCard.Contract.Owner(&_ZBGCard.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(_tokenId uint256) constant returns(address)
func (_ZBGCard *ZBGCardCaller) OwnerOf(opts *bind.CallOpts, _tokenId *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ZBGCard.contract.Call(opts, out, "ownerOf", _tokenId)
	return *ret0, err
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(_tokenId uint256) constant returns(address)
func (_ZBGCard *ZBGCardSession) OwnerOf(_tokenId *big.Int) (common.Address, error) {
	return _ZBGCard.Contract.OwnerOf(&_ZBGCard.CallOpts, _tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(_tokenId uint256) constant returns(address)
func (_ZBGCard *ZBGCardCallerSession) OwnerOf(_tokenId *big.Int) (common.Address, error) {
	return _ZBGCard.Contract.OwnerOf(&_ZBGCard.CallOpts, _tokenId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(_interfaceId bytes4) constant returns(bool)
func (_ZBGCard *ZBGCardCaller) SupportsInterface(opts *bind.CallOpts, _interfaceId [4]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ZBGCard.contract.Call(opts, out, "supportsInterface", _interfaceId)
	return *ret0, err
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(_interfaceId bytes4) constant returns(bool)
func (_ZBGCard *ZBGCardSession) SupportsInterface(_interfaceId [4]byte) (bool, error) {
	return _ZBGCard.Contract.SupportsInterface(&_ZBGCard.CallOpts, _interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(_interfaceId bytes4) constant returns(bool)
func (_ZBGCard *ZBGCardCallerSession) SupportsInterface(_interfaceId [4]byte) (bool, error) {
	return _ZBGCard.Contract.SupportsInterface(&_ZBGCard.CallOpts, _interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_ZBGCard *ZBGCardCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _ZBGCard.contract.Call(opts, out, "symbol")
	return *ret0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_ZBGCard *ZBGCardSession) Symbol() (string, error) {
	return _ZBGCard.Contract.Symbol(&_ZBGCard.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_ZBGCard *ZBGCardCallerSession) Symbol() (string, error) {
	return _ZBGCard.Contract.Symbol(&_ZBGCard.CallOpts)
}

// TokenByIndex is a free data retrieval call binding the contract method 0x4f6ccce7.
//
// Solidity: function tokenByIndex(_index uint256) constant returns(uint256)
func (_ZBGCard *ZBGCardCaller) TokenByIndex(opts *bind.CallOpts, _index *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ZBGCard.contract.Call(opts, out, "tokenByIndex", _index)
	return *ret0, err
}

// TokenByIndex is a free data retrieval call binding the contract method 0x4f6ccce7.
//
// Solidity: function tokenByIndex(_index uint256) constant returns(uint256)
func (_ZBGCard *ZBGCardSession) TokenByIndex(_index *big.Int) (*big.Int, error) {
	return _ZBGCard.Contract.TokenByIndex(&_ZBGCard.CallOpts, _index)
}

// TokenByIndex is a free data retrieval call binding the contract method 0x4f6ccce7.
//
// Solidity: function tokenByIndex(_index uint256) constant returns(uint256)
func (_ZBGCard *ZBGCardCallerSession) TokenByIndex(_index *big.Int) (*big.Int, error) {
	return _ZBGCard.Contract.TokenByIndex(&_ZBGCard.CallOpts, _index)
}

// TokenIdToDNA is a free data retrieval call binding the contract method 0x6e5ed979.
//
// Solidity: function tokenIdToDNA( uint256) constant returns(uint256)
func (_ZBGCard *ZBGCardCaller) TokenIdToDNA(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ZBGCard.contract.Call(opts, out, "tokenIdToDNA", arg0)
	return *ret0, err
}

// TokenIdToDNA is a free data retrieval call binding the contract method 0x6e5ed979.
//
// Solidity: function tokenIdToDNA( uint256) constant returns(uint256)
func (_ZBGCard *ZBGCardSession) TokenIdToDNA(arg0 *big.Int) (*big.Int, error) {
	return _ZBGCard.Contract.TokenIdToDNA(&_ZBGCard.CallOpts, arg0)
}

// TokenIdToDNA is a free data retrieval call binding the contract method 0x6e5ed979.
//
// Solidity: function tokenIdToDNA( uint256) constant returns(uint256)
func (_ZBGCard *ZBGCardCallerSession) TokenIdToDNA(arg0 *big.Int) (*big.Int, error) {
	return _ZBGCard.Contract.TokenIdToDNA(&_ZBGCard.CallOpts, arg0)
}

// TokenOfOwnerByIndex is a free data retrieval call binding the contract method 0x2f745c59.
//
// Solidity: function tokenOfOwnerByIndex(_owner address, _index uint256) constant returns(_tokenId uint256)
func (_ZBGCard *ZBGCardCaller) TokenOfOwnerByIndex(opts *bind.CallOpts, _owner common.Address, _index *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ZBGCard.contract.Call(opts, out, "tokenOfOwnerByIndex", _owner, _index)
	return *ret0, err
}

// TokenOfOwnerByIndex is a free data retrieval call binding the contract method 0x2f745c59.
//
// Solidity: function tokenOfOwnerByIndex(_owner address, _index uint256) constant returns(_tokenId uint256)
func (_ZBGCard *ZBGCardSession) TokenOfOwnerByIndex(_owner common.Address, _index *big.Int) (*big.Int, error) {
	return _ZBGCard.Contract.TokenOfOwnerByIndex(&_ZBGCard.CallOpts, _owner, _index)
}

// TokenOfOwnerByIndex is a free data retrieval call binding the contract method 0x2f745c59.
//
// Solidity: function tokenOfOwnerByIndex(_owner address, _index uint256) constant returns(_tokenId uint256)
func (_ZBGCard *ZBGCardCallerSession) TokenOfOwnerByIndex(_owner common.Address, _index *big.Int) (*big.Int, error) {
	return _ZBGCard.Contract.TokenOfOwnerByIndex(&_ZBGCard.CallOpts, _owner, _index)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(_tokenId uint256) constant returns(tokenUri string)
func (_ZBGCard *ZBGCardCaller) TokenURI(opts *bind.CallOpts, _tokenId *big.Int) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _ZBGCard.contract.Call(opts, out, "tokenURI", _tokenId)
	return *ret0, err
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(_tokenId uint256) constant returns(tokenUri string)
func (_ZBGCard *ZBGCardSession) TokenURI(_tokenId *big.Int) (string, error) {
	return _ZBGCard.Contract.TokenURI(&_ZBGCard.CallOpts, _tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(_tokenId uint256) constant returns(tokenUri string)
func (_ZBGCard *ZBGCardCallerSession) TokenURI(_tokenId *big.Int) (string, error) {
	return _ZBGCard.Contract.TokenURI(&_ZBGCard.CallOpts, _tokenId)
}

// TokensOwned is a free data retrieval call binding the contract method 0x21cda790.
//
// Solidity: function tokensOwned(_owner address) constant returns(indexes uint256[], balances uint256[])
func (_ZBGCard *ZBGCardCaller) TokensOwned(opts *bind.CallOpts, _owner common.Address) (struct {
	Indexes  []*big.Int
	Balances []*big.Int
}, error) {
	ret := new(struct {
		Indexes  []*big.Int
		Balances []*big.Int
	})
	out := ret
	err := _ZBGCard.contract.Call(opts, out, "tokensOwned", _owner)
	return *ret, err
}

// TokensOwned is a free data retrieval call binding the contract method 0x21cda790.
//
// Solidity: function tokensOwned(_owner address) constant returns(indexes uint256[], balances uint256[])
func (_ZBGCard *ZBGCardSession) TokensOwned(_owner common.Address) (struct {
	Indexes  []*big.Int
	Balances []*big.Int
}, error) {
	return _ZBGCard.Contract.TokensOwned(&_ZBGCard.CallOpts, _owner)
}

// TokensOwned is a free data retrieval call binding the contract method 0x21cda790.
//
// Solidity: function tokensOwned(_owner address) constant returns(indexes uint256[], balances uint256[])
func (_ZBGCard *ZBGCardCallerSession) TokensOwned(_owner common.Address) (struct {
	Indexes  []*big.Int
	Balances []*big.Int
}, error) {
	return _ZBGCard.Contract.TokensOwned(&_ZBGCard.CallOpts, _owner)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ZBGCard *ZBGCardCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ZBGCard.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ZBGCard *ZBGCardSession) TotalSupply() (*big.Int, error) {
	return _ZBGCard.Contract.TotalSupply(&_ZBGCard.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ZBGCard *ZBGCardCallerSession) TotalSupply() (*big.Int, error) {
	return _ZBGCard.Contract.TotalSupply(&_ZBGCard.CallOpts)
}

// AddValidator is a paid mutator transaction binding the contract method 0x90b616c8.
//
// Solidity: function addValidator(_validator address, _v uint8[], _r bytes32[], _s bytes32[]) returns()
func (_ZBGCard *ZBGCardTransactor) AddValidator(opts *bind.TransactOpts, _validator common.Address, _v []uint8, _r [][32]byte, _s [][32]byte) (*types.Transaction, error) {
	return _ZBGCard.contract.Transact(opts, "addValidator", _validator, _v, _r, _s)
}

// AddValidator is a paid mutator transaction binding the contract method 0x90b616c8.
//
// Solidity: function addValidator(_validator address, _v uint8[], _r bytes32[], _s bytes32[]) returns()
func (_ZBGCard *ZBGCardSession) AddValidator(_validator common.Address, _v []uint8, _r [][32]byte, _s [][32]byte) (*types.Transaction, error) {
	return _ZBGCard.Contract.AddValidator(&_ZBGCard.TransactOpts, _validator, _v, _r, _s)
}

// AddValidator is a paid mutator transaction binding the contract method 0x90b616c8.
//
// Solidity: function addValidator(_validator address, _v uint8[], _r bytes32[], _s bytes32[]) returns()
func (_ZBGCard *ZBGCardTransactorSession) AddValidator(_validator common.Address, _v []uint8, _r [][32]byte, _s [][32]byte) (*types.Transaction, error) {
	return _ZBGCard.Contract.AddValidator(&_ZBGCard.TransactOpts, _validator, _v, _r, _s)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_to address, _tokenId uint256) returns()
func (_ZBGCard *ZBGCardTransactor) Approve(opts *bind.TransactOpts, _to common.Address, _tokenId *big.Int) (*types.Transaction, error) {
	return _ZBGCard.contract.Transact(opts, "approve", _to, _tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_to address, _tokenId uint256) returns()
func (_ZBGCard *ZBGCardSession) Approve(_to common.Address, _tokenId *big.Int) (*types.Transaction, error) {
	return _ZBGCard.Contract.Approve(&_ZBGCard.TransactOpts, _to, _tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_to address, _tokenId uint256) returns()
func (_ZBGCard *ZBGCardTransactorSession) Approve(_to common.Address, _tokenId *big.Int) (*types.Transaction, error) {
	return _ZBGCard.Contract.Approve(&_ZBGCard.TransactOpts, _to, _tokenId)
}

// BatchMintToken is a paid mutator transaction binding the contract method 0x480cd004.
//
// Solidity: function batchMintToken(_tokenIds uint256[], _tokenDNAs uint256[]) returns()
func (_ZBGCard *ZBGCardTransactor) BatchMintToken(opts *bind.TransactOpts, _tokenIds []*big.Int, _tokenDNAs []*big.Int) (*types.Transaction, error) {
	return _ZBGCard.contract.Transact(opts, "batchMintToken", _tokenIds, _tokenDNAs)
}

// BatchMintToken is a paid mutator transaction binding the contract method 0x480cd004.
//
// Solidity: function batchMintToken(_tokenIds uint256[], _tokenDNAs uint256[]) returns()
func (_ZBGCard *ZBGCardSession) BatchMintToken(_tokenIds []*big.Int, _tokenDNAs []*big.Int) (*types.Transaction, error) {
	return _ZBGCard.Contract.BatchMintToken(&_ZBGCard.TransactOpts, _tokenIds, _tokenDNAs)
}

// BatchMintToken is a paid mutator transaction binding the contract method 0x480cd004.
//
// Solidity: function batchMintToken(_tokenIds uint256[], _tokenDNAs uint256[]) returns()
func (_ZBGCard *ZBGCardTransactorSession) BatchMintToken(_tokenIds []*big.Int, _tokenDNAs []*big.Int) (*types.Transaction, error) {
	return _ZBGCard.Contract.BatchMintToken(&_ZBGCard.TransactOpts, _tokenIds, _tokenDNAs)
}

// BatchTransferFrom is a paid mutator transaction binding the contract method 0x17fad7fc.
//
// Solidity: function batchTransferFrom(_from address, _to address, _tokenIds uint256[], _amounts uint256[]) returns()
func (_ZBGCard *ZBGCardTransactor) BatchTransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _tokenIds []*big.Int, _amounts []*big.Int) (*types.Transaction, error) {
	return _ZBGCard.contract.Transact(opts, "batchTransferFrom", _from, _to, _tokenIds, _amounts)
}

// BatchTransferFrom is a paid mutator transaction binding the contract method 0x17fad7fc.
//
// Solidity: function batchTransferFrom(_from address, _to address, _tokenIds uint256[], _amounts uint256[]) returns()
func (_ZBGCard *ZBGCardSession) BatchTransferFrom(_from common.Address, _to common.Address, _tokenIds []*big.Int, _amounts []*big.Int) (*types.Transaction, error) {
	return _ZBGCard.Contract.BatchTransferFrom(&_ZBGCard.TransactOpts, _from, _to, _tokenIds, _amounts)
}

// BatchTransferFrom is a paid mutator transaction binding the contract method 0x17fad7fc.
//
// Solidity: function batchTransferFrom(_from address, _to address, _tokenIds uint256[], _amounts uint256[]) returns()
func (_ZBGCard *ZBGCardTransactorSession) BatchTransferFrom(_from common.Address, _to common.Address, _tokenIds []*big.Int, _amounts []*big.Int) (*types.Transaction, error) {
	return _ZBGCard.Contract.BatchTransferFrom(&_ZBGCard.TransactOpts, _from, _to, _tokenIds, _amounts)
}

// ClaimToken is a paid mutator transaction binding the contract method 0xb68e6fea.
//
// Solidity: function claimToken(_tokenId uint256, _amount uint256, _sig bytes) returns()
func (_ZBGCard *ZBGCardTransactor) ClaimToken(opts *bind.TransactOpts, _tokenId *big.Int, _amount *big.Int, _sig []byte) (*types.Transaction, error) {
	return _ZBGCard.contract.Transact(opts, "claimToken", _tokenId, _amount, _sig)
}

// ClaimToken is a paid mutator transaction binding the contract method 0xb68e6fea.
//
// Solidity: function claimToken(_tokenId uint256, _amount uint256, _sig bytes) returns()
func (_ZBGCard *ZBGCardSession) ClaimToken(_tokenId *big.Int, _amount *big.Int, _sig []byte) (*types.Transaction, error) {
	return _ZBGCard.Contract.ClaimToken(&_ZBGCard.TransactOpts, _tokenId, _amount, _sig)
}

// ClaimToken is a paid mutator transaction binding the contract method 0xb68e6fea.
//
// Solidity: function claimToken(_tokenId uint256, _amount uint256, _sig bytes) returns()
func (_ZBGCard *ZBGCardTransactorSession) ClaimToken(_tokenId *big.Int, _amount *big.Int, _sig []byte) (*types.Transaction, error) {
	return _ZBGCard.Contract.ClaimToken(&_ZBGCard.TransactOpts, _tokenId, _amount, _sig)
}

// ClaimTokenNFT is a paid mutator transaction binding the contract method 0xed186dd3.
//
// Solidity: function claimTokenNFT(_tokenId uint256, _sig bytes) returns()
func (_ZBGCard *ZBGCardTransactor) ClaimTokenNFT(opts *bind.TransactOpts, _tokenId *big.Int, _sig []byte) (*types.Transaction, error) {
	return _ZBGCard.contract.Transact(opts, "claimTokenNFT", _tokenId, _sig)
}

// ClaimTokenNFT is a paid mutator transaction binding the contract method 0xed186dd3.
//
// Solidity: function claimTokenNFT(_tokenId uint256, _sig bytes) returns()
func (_ZBGCard *ZBGCardSession) ClaimTokenNFT(_tokenId *big.Int, _sig []byte) (*types.Transaction, error) {
	return _ZBGCard.Contract.ClaimTokenNFT(&_ZBGCard.TransactOpts, _tokenId, _sig)
}

// ClaimTokenNFT is a paid mutator transaction binding the contract method 0xed186dd3.
//
// Solidity: function claimTokenNFT(_tokenId uint256, _sig bytes) returns()
func (_ZBGCard *ZBGCardTransactorSession) ClaimTokenNFT(_tokenId *big.Int, _sig []byte) (*types.Transaction, error) {
	return _ZBGCard.Contract.ClaimTokenNFT(&_ZBGCard.TransactOpts, _tokenId, _sig)
}

// DepositToGateway is a paid mutator transaction binding the contract method 0x45f0edb7.
//
// Solidity: function depositToGateway(_tokenId uint256, amount uint256) returns()
func (_ZBGCard *ZBGCardTransactor) DepositToGateway(opts *bind.TransactOpts, _tokenId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _ZBGCard.contract.Transact(opts, "depositToGateway", _tokenId, amount)
}

// DepositToGateway is a paid mutator transaction binding the contract method 0x45f0edb7.
//
// Solidity: function depositToGateway(_tokenId uint256, amount uint256) returns()
func (_ZBGCard *ZBGCardSession) DepositToGateway(_tokenId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _ZBGCard.Contract.DepositToGateway(&_ZBGCard.TransactOpts, _tokenId, amount)
}

// DepositToGateway is a paid mutator transaction binding the contract method 0x45f0edb7.
//
// Solidity: function depositToGateway(_tokenId uint256, amount uint256) returns()
func (_ZBGCard *ZBGCardTransactorSession) DepositToGateway(_tokenId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _ZBGCard.Contract.DepositToGateway(&_ZBGCard.TransactOpts, _tokenId, amount)
}

// DisableFaucet is a paid mutator transaction binding the contract method 0x87a8af4e.
//
// Solidity: function disableFaucet(_faucet address) returns()
func (_ZBGCard *ZBGCardTransactor) DisableFaucet(opts *bind.TransactOpts, _faucet common.Address) (*types.Transaction, error) {
	return _ZBGCard.contract.Transact(opts, "disableFaucet", _faucet)
}

// DisableFaucet is a paid mutator transaction binding the contract method 0x87a8af4e.
//
// Solidity: function disableFaucet(_faucet address) returns()
func (_ZBGCard *ZBGCardSession) DisableFaucet(_faucet common.Address) (*types.Transaction, error) {
	return _ZBGCard.Contract.DisableFaucet(&_ZBGCard.TransactOpts, _faucet)
}

// DisableFaucet is a paid mutator transaction binding the contract method 0x87a8af4e.
//
// Solidity: function disableFaucet(_faucet address) returns()
func (_ZBGCard *ZBGCardTransactorSession) DisableFaucet(_faucet common.Address) (*types.Transaction, error) {
	return _ZBGCard.Contract.DisableFaucet(&_ZBGCard.TransactOpts, _faucet)
}

// EnableFaucet is a paid mutator transaction binding the contract method 0xe4596dc4.
//
// Solidity: function enableFaucet(_faucet address) returns()
func (_ZBGCard *ZBGCardTransactor) EnableFaucet(opts *bind.TransactOpts, _faucet common.Address) (*types.Transaction, error) {
	return _ZBGCard.contract.Transact(opts, "enableFaucet", _faucet)
}

// EnableFaucet is a paid mutator transaction binding the contract method 0xe4596dc4.
//
// Solidity: function enableFaucet(_faucet address) returns()
func (_ZBGCard *ZBGCardSession) EnableFaucet(_faucet common.Address) (*types.Transaction, error) {
	return _ZBGCard.Contract.EnableFaucet(&_ZBGCard.TransactOpts, _faucet)
}

// EnableFaucet is a paid mutator transaction binding the contract method 0xe4596dc4.
//
// Solidity: function enableFaucet(_faucet address) returns()
func (_ZBGCard *ZBGCardTransactorSession) EnableFaucet(_faucet common.Address) (*types.Transaction, error) {
	return _ZBGCard.Contract.EnableFaucet(&_ZBGCard.TransactOpts, _faucet)
}

// MintToken is a paid mutator transaction binding the contract method 0x20cbf5f9.
//
// Solidity: function mintToken(_tokenId uint256, _tokenDNA uint256) returns()
func (_ZBGCard *ZBGCardTransactor) MintToken(opts *bind.TransactOpts, _tokenId *big.Int, _tokenDNA *big.Int) (*types.Transaction, error) {
	return _ZBGCard.contract.Transact(opts, "mintToken", _tokenId, _tokenDNA)
}

// MintToken is a paid mutator transaction binding the contract method 0x20cbf5f9.
//
// Solidity: function mintToken(_tokenId uint256, _tokenDNA uint256) returns()
func (_ZBGCard *ZBGCardSession) MintToken(_tokenId *big.Int, _tokenDNA *big.Int) (*types.Transaction, error) {
	return _ZBGCard.Contract.MintToken(&_ZBGCard.TransactOpts, _tokenId, _tokenDNA)
}

// MintToken is a paid mutator transaction binding the contract method 0x20cbf5f9.
//
// Solidity: function mintToken(_tokenId uint256, _tokenDNA uint256) returns()
func (_ZBGCard *ZBGCardTransactorSession) MintToken(_tokenId *big.Int, _tokenDNA *big.Int) (*types.Transaction, error) {
	return _ZBGCard.Contract.MintToken(&_ZBGCard.TransactOpts, _tokenId, _tokenDNA)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0xc7e7f6f6.
//
// Solidity: function removeValidator(_validator address, _v uint8[], _r bytes32[], _s bytes32[]) returns()
func (_ZBGCard *ZBGCardTransactor) RemoveValidator(opts *bind.TransactOpts, _validator common.Address, _v []uint8, _r [][32]byte, _s [][32]byte) (*types.Transaction, error) {
	return _ZBGCard.contract.Transact(opts, "removeValidator", _validator, _v, _r, _s)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0xc7e7f6f6.
//
// Solidity: function removeValidator(_validator address, _v uint8[], _r bytes32[], _s bytes32[]) returns()
func (_ZBGCard *ZBGCardSession) RemoveValidator(_validator common.Address, _v []uint8, _r [][32]byte, _s [][32]byte) (*types.Transaction, error) {
	return _ZBGCard.Contract.RemoveValidator(&_ZBGCard.TransactOpts, _validator, _v, _r, _s)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0xc7e7f6f6.
//
// Solidity: function removeValidator(_validator address, _v uint8[], _r bytes32[], _s bytes32[]) returns()
func (_ZBGCard *ZBGCardTransactorSession) RemoveValidator(_validator common.Address, _v []uint8, _r [][32]byte, _s [][32]byte) (*types.Transaction, error) {
	return _ZBGCard.Contract.RemoveValidator(&_ZBGCard.TransactOpts, _validator, _v, _r, _s)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ZBGCard *ZBGCardTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZBGCard.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ZBGCard *ZBGCardSession) RenounceOwnership() (*types.Transaction, error) {
	return _ZBGCard.Contract.RenounceOwnership(&_ZBGCard.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ZBGCard *ZBGCardTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ZBGCard.Contract.RenounceOwnership(&_ZBGCard.TransactOpts)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(_from address, _to address, _tokenIds uint256[], _amounts uint256[], _data bytes) returns()
func (_ZBGCard *ZBGCardTransactor) SafeBatchTransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _tokenIds []*big.Int, _amounts []*big.Int, _data []byte) (*types.Transaction, error) {
	return _ZBGCard.contract.Transact(opts, "safeBatchTransferFrom", _from, _to, _tokenIds, _amounts, _data)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(_from address, _to address, _tokenIds uint256[], _amounts uint256[], _data bytes) returns()
func (_ZBGCard *ZBGCardSession) SafeBatchTransferFrom(_from common.Address, _to common.Address, _tokenIds []*big.Int, _amounts []*big.Int, _data []byte) (*types.Transaction, error) {
	return _ZBGCard.Contract.SafeBatchTransferFrom(&_ZBGCard.TransactOpts, _from, _to, _tokenIds, _amounts, _data)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(_from address, _to address, _tokenIds uint256[], _amounts uint256[], _data bytes) returns()
func (_ZBGCard *ZBGCardTransactorSession) SafeBatchTransferFrom(_from common.Address, _to common.Address, _tokenIds []*big.Int, _amounts []*big.Int, _data []byte) (*types.Transaction, error) {
	return _ZBGCard.Contract.SafeBatchTransferFrom(&_ZBGCard.TransactOpts, _from, _to, _tokenIds, _amounts, _data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(_from address, _to address, _tokenId uint256, _amount uint256, _data bytes) returns()
func (_ZBGCard *ZBGCardTransactor) SafeTransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _tokenId *big.Int, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _ZBGCard.contract.Transact(opts, "safeTransferFrom", _from, _to, _tokenId, _amount, _data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(_from address, _to address, _tokenId uint256, _amount uint256, _data bytes) returns()
func (_ZBGCard *ZBGCardSession) SafeTransferFrom(_from common.Address, _to common.Address, _tokenId *big.Int, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _ZBGCard.Contract.SafeTransferFrom(&_ZBGCard.TransactOpts, _from, _to, _tokenId, _amount, _data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(_from address, _to address, _tokenId uint256, _amount uint256, _data bytes) returns()
func (_ZBGCard *ZBGCardTransactorSession) SafeTransferFrom(_from common.Address, _to common.Address, _tokenId *big.Int, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _ZBGCard.Contract.SafeTransferFrom(&_ZBGCard.TransactOpts, _from, _to, _tokenId, _amount, _data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(_operator address, _approved bool) returns()
func (_ZBGCard *ZBGCardTransactor) SetApprovalForAll(opts *bind.TransactOpts, _operator common.Address, _approved bool) (*types.Transaction, error) {
	return _ZBGCard.contract.Transact(opts, "setApprovalForAll", _operator, _approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(_operator address, _approved bool) returns()
func (_ZBGCard *ZBGCardSession) SetApprovalForAll(_operator common.Address, _approved bool) (*types.Transaction, error) {
	return _ZBGCard.Contract.SetApprovalForAll(&_ZBGCard.TransactOpts, _operator, _approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(_operator address, _approved bool) returns()
func (_ZBGCard *ZBGCardTransactorSession) SetApprovalForAll(_operator common.Address, _approved bool) (*types.Transaction, error) {
	return _ZBGCard.Contract.SetApprovalForAll(&_ZBGCard.TransactOpts, _operator, _approved)
}

// ToggleToken is a paid mutator transaction binding the contract method 0x15c75f89.
//
// Solidity: function toggleToken(_token address) returns()
func (_ZBGCard *ZBGCardTransactor) ToggleToken(opts *bind.TransactOpts, _token common.Address) (*types.Transaction, error) {
	return _ZBGCard.contract.Transact(opts, "toggleToken", _token)
}

// ToggleToken is a paid mutator transaction binding the contract method 0x15c75f89.
//
// Solidity: function toggleToken(_token address) returns()
func (_ZBGCard *ZBGCardSession) ToggleToken(_token common.Address) (*types.Transaction, error) {
	return _ZBGCard.Contract.ToggleToken(&_ZBGCard.TransactOpts, _token)
}

// ToggleToken is a paid mutator transaction binding the contract method 0x15c75f89.
//
// Solidity: function toggleToken(_token address) returns()
func (_ZBGCard *ZBGCardTransactorSession) ToggleToken(_token common.Address) (*types.Transaction, error) {
	return _ZBGCard.Contract.ToggleToken(&_ZBGCard.TransactOpts, _token)
}

// Transfer is a paid mutator transaction binding the contract method 0x095bcdb6.
//
// Solidity: function transfer(_to address, _tokenId uint256, _amount uint256) returns()
func (_ZBGCard *ZBGCardTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _tokenId *big.Int, _amount *big.Int) (*types.Transaction, error) {
	return _ZBGCard.contract.Transact(opts, "transfer", _to, _tokenId, _amount)
}

// Transfer is a paid mutator transaction binding the contract method 0x095bcdb6.
//
// Solidity: function transfer(_to address, _tokenId uint256, _amount uint256) returns()
func (_ZBGCard *ZBGCardSession) Transfer(_to common.Address, _tokenId *big.Int, _amount *big.Int) (*types.Transaction, error) {
	return _ZBGCard.Contract.Transfer(&_ZBGCard.TransactOpts, _to, _tokenId, _amount)
}

// Transfer is a paid mutator transaction binding the contract method 0x095bcdb6.
//
// Solidity: function transfer(_to address, _tokenId uint256, _amount uint256) returns()
func (_ZBGCard *ZBGCardTransactorSession) Transfer(_to common.Address, _tokenId *big.Int, _amount *big.Int) (*types.Transaction, error) {
	return _ZBGCard.Contract.Transfer(&_ZBGCard.TransactOpts, _to, _tokenId, _amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0xfe99049a.
//
// Solidity: function transferFrom(_from address, _to address, _tokenId uint256, _amount uint256) returns()
func (_ZBGCard *ZBGCardTransactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _tokenId *big.Int, _amount *big.Int) (*types.Transaction, error) {
	return _ZBGCard.contract.Transact(opts, "transferFrom", _from, _to, _tokenId, _amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0xfe99049a.
//
// Solidity: function transferFrom(_from address, _to address, _tokenId uint256, _amount uint256) returns()
func (_ZBGCard *ZBGCardSession) TransferFrom(_from common.Address, _to common.Address, _tokenId *big.Int, _amount *big.Int) (*types.Transaction, error) {
	return _ZBGCard.Contract.TransferFrom(&_ZBGCard.TransactOpts, _from, _to, _tokenId, _amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0xfe99049a.
//
// Solidity: function transferFrom(_from address, _to address, _tokenId uint256, _amount uint256) returns()
func (_ZBGCard *ZBGCardTransactorSession) TransferFrom(_from common.Address, _to common.Address, _tokenId *big.Int, _amount *big.Int) (*types.Transaction, error) {
	return _ZBGCard.Contract.TransferFrom(&_ZBGCard.TransactOpts, _from, _to, _tokenId, _amount)
}

// TransferFrom1 is a paid mutator transaction binding the contract method 0x919a96cc.
//
// Solidity: function transferFrom1(_from address, _to address, _tokenId uint256, _amount uint256) returns()
func (_ZBGCard *ZBGCardTransactor) TransferFrom1(opts *bind.TransactOpts, _from common.Address, _to common.Address, _tokenId *big.Int, _amount *big.Int) (*types.Transaction, error) {
	return _ZBGCard.contract.Transact(opts, "transferFrom1", _from, _to, _tokenId, _amount)
}

// TransferFrom1 is a paid mutator transaction binding the contract method 0x919a96cc.
//
// Solidity: function transferFrom1(_from address, _to address, _tokenId uint256, _amount uint256) returns()
func (_ZBGCard *ZBGCardSession) TransferFrom1(_from common.Address, _to common.Address, _tokenId *big.Int, _amount *big.Int) (*types.Transaction, error) {
	return _ZBGCard.Contract.TransferFrom1(&_ZBGCard.TransactOpts, _from, _to, _tokenId, _amount)
}

// TransferFrom1 is a paid mutator transaction binding the contract method 0x919a96cc.
//
// Solidity: function transferFrom1(_from address, _to address, _tokenId uint256, _amount uint256) returns()
func (_ZBGCard *ZBGCardTransactorSession) TransferFrom1(_from common.Address, _to common.Address, _tokenId *big.Int, _amount *big.Int) (*types.Transaction, error) {
	return _ZBGCard.Contract.TransferFrom1(&_ZBGCard.TransactOpts, _from, _to, _tokenId, _amount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_ZBGCard *ZBGCardTransactor) TransferOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _ZBGCard.contract.Transact(opts, "transferOwnership", _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_ZBGCard *ZBGCardSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _ZBGCard.Contract.TransferOwnership(&_ZBGCard.TransactOpts, _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_ZBGCard *ZBGCardTransactorSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _ZBGCard.Contract.TransferOwnership(&_ZBGCard.TransactOpts, _newOwner)
}

// ZBGCardAddedValidatorIterator is returned from FilterAddedValidator and is used to iterate over the raw logs and unpacked data for AddedValidator events raised by the ZBGCard contract.
type ZBGCardAddedValidatorIterator struct {
	Event *ZBGCardAddedValidator // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ZBGCardAddedValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZBGCardAddedValidator)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ZBGCardAddedValidator)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ZBGCardAddedValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZBGCardAddedValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZBGCardAddedValidator represents a AddedValidator event raised by the ZBGCard contract.
type ZBGCardAddedValidator struct {
	Validator common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAddedValidator is a free log retrieval operation binding the contract event 0x8e15bf46bd11add443414ada75aa9592a4af68f3f2ec02ae3d49572f9843c2a8.
//
// Solidity: e AddedValidator(validator address)
func (_ZBGCard *ZBGCardFilterer) FilterAddedValidator(opts *bind.FilterOpts) (*ZBGCardAddedValidatorIterator, error) {

	logs, sub, err := _ZBGCard.contract.FilterLogs(opts, "AddedValidator")
	if err != nil {
		return nil, err
	}
	return &ZBGCardAddedValidatorIterator{contract: _ZBGCard.contract, event: "AddedValidator", logs: logs, sub: sub}, nil
}

// WatchAddedValidator is a free log subscription operation binding the contract event 0x8e15bf46bd11add443414ada75aa9592a4af68f3f2ec02ae3d49572f9843c2a8.
//
// Solidity: e AddedValidator(validator address)
func (_ZBGCard *ZBGCardFilterer) WatchAddedValidator(opts *bind.WatchOpts, sink chan<- *ZBGCardAddedValidator) (event.Subscription, error) {

	logs, sub, err := _ZBGCard.contract.WatchLogs(opts, "AddedValidator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZBGCardAddedValidator)
				if err := _ZBGCard.contract.UnpackLog(event, "AddedValidator", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ZBGCardApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ZBGCard contract.
type ZBGCardApprovalIterator struct {
	Event *ZBGCardApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ZBGCardApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZBGCardApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ZBGCardApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ZBGCardApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZBGCardApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZBGCardApproval represents a Approval event raised by the ZBGCard contract.
type ZBGCardApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(_owner indexed address, _approved indexed address, _tokenId indexed uint256)
func (_ZBGCard *ZBGCardFilterer) FilterApproval(opts *bind.FilterOpts, _owner []common.Address, _approved []common.Address, _tokenId []*big.Int) (*ZBGCardApprovalIterator, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _approvedRule []interface{}
	for _, _approvedItem := range _approved {
		_approvedRule = append(_approvedRule, _approvedItem)
	}
	var _tokenIdRule []interface{}
	for _, _tokenIdItem := range _tokenId {
		_tokenIdRule = append(_tokenIdRule, _tokenIdItem)
	}

	logs, sub, err := _ZBGCard.contract.FilterLogs(opts, "Approval", _ownerRule, _approvedRule, _tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &ZBGCardApprovalIterator{contract: _ZBGCard.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(_owner indexed address, _approved indexed address, _tokenId indexed uint256)
func (_ZBGCard *ZBGCardFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ZBGCardApproval, _owner []common.Address, _approved []common.Address, _tokenId []*big.Int) (event.Subscription, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _approvedRule []interface{}
	for _, _approvedItem := range _approved {
		_approvedRule = append(_approvedRule, _approvedItem)
	}
	var _tokenIdRule []interface{}
	for _, _tokenIdItem := range _tokenId {
		_tokenIdRule = append(_tokenIdRule, _tokenIdItem)
	}

	logs, sub, err := _ZBGCard.contract.WatchLogs(opts, "Approval", _ownerRule, _approvedRule, _tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZBGCardApproval)
				if err := _ZBGCard.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ZBGCardApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the ZBGCard contract.
type ZBGCardApprovalForAllIterator struct {
	Event *ZBGCardApprovalForAll // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ZBGCardApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZBGCardApprovalForAll)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ZBGCardApprovalForAll)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ZBGCardApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZBGCardApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZBGCardApprovalForAll represents a ApprovalForAll event raised by the ZBGCard contract.
type ZBGCardApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: e ApprovalForAll(_owner indexed address, _operator indexed address, _approved bool)
func (_ZBGCard *ZBGCardFilterer) FilterApprovalForAll(opts *bind.FilterOpts, _owner []common.Address, _operator []common.Address) (*ZBGCardApprovalForAllIterator, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _operatorRule []interface{}
	for _, _operatorItem := range _operator {
		_operatorRule = append(_operatorRule, _operatorItem)
	}

	logs, sub, err := _ZBGCard.contract.FilterLogs(opts, "ApprovalForAll", _ownerRule, _operatorRule)
	if err != nil {
		return nil, err
	}
	return &ZBGCardApprovalForAllIterator{contract: _ZBGCard.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: e ApprovalForAll(_owner indexed address, _operator indexed address, _approved bool)
func (_ZBGCard *ZBGCardFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *ZBGCardApprovalForAll, _owner []common.Address, _operator []common.Address) (event.Subscription, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _operatorRule []interface{}
	for _, _operatorItem := range _operator {
		_operatorRule = append(_operatorRule, _operatorItem)
	}

	logs, sub, err := _ZBGCard.contract.WatchLogs(opts, "ApprovalForAll", _ownerRule, _operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZBGCardApprovalForAll)
				if err := _ZBGCard.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ZBGCardBatchTransferIterator is returned from FilterBatchTransfer and is used to iterate over the raw logs and unpacked data for BatchTransfer events raised by the ZBGCard contract.
type ZBGCardBatchTransferIterator struct {
	Event *ZBGCardBatchTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ZBGCardBatchTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZBGCardBatchTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ZBGCardBatchTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ZBGCardBatchTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZBGCardBatchTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZBGCardBatchTransfer represents a BatchTransfer event raised by the ZBGCard contract.
type ZBGCardBatchTransfer struct {
	From       common.Address
	To         common.Address
	TokenTypes []*big.Int
	Amounts    []*big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterBatchTransfer is a free log retrieval operation binding the contract event 0xf59807b2c31ca3ba212e90599175c120c556422950bac5be656274483e8581df.
//
// Solidity: e BatchTransfer(from address, to address, tokenTypes uint256[], amounts uint256[])
func (_ZBGCard *ZBGCardFilterer) FilterBatchTransfer(opts *bind.FilterOpts) (*ZBGCardBatchTransferIterator, error) {

	logs, sub, err := _ZBGCard.contract.FilterLogs(opts, "BatchTransfer")
	if err != nil {
		return nil, err
	}
	return &ZBGCardBatchTransferIterator{contract: _ZBGCard.contract, event: "BatchTransfer", logs: logs, sub: sub}, nil
}

// WatchBatchTransfer is a free log subscription operation binding the contract event 0xf59807b2c31ca3ba212e90599175c120c556422950bac5be656274483e8581df.
//
// Solidity: e BatchTransfer(from address, to address, tokenTypes uint256[], amounts uint256[])
func (_ZBGCard *ZBGCardFilterer) WatchBatchTransfer(opts *bind.WatchOpts, sink chan<- *ZBGCardBatchTransfer) (event.Subscription, error) {

	logs, sub, err := _ZBGCard.contract.WatchLogs(opts, "BatchTransfer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZBGCardBatchTransfer)
				if err := _ZBGCard.contract.UnpackLog(event, "BatchTransfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ZBGCardFaucetToggledIterator is returned from FilterFaucetToggled and is used to iterate over the raw logs and unpacked data for FaucetToggled events raised by the ZBGCard contract.
type ZBGCardFaucetToggledIterator struct {
	Event *ZBGCardFaucetToggled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ZBGCardFaucetToggledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZBGCardFaucetToggled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ZBGCardFaucetToggled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ZBGCardFaucetToggledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZBGCardFaucetToggledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZBGCardFaucetToggled represents a FaucetToggled event raised by the ZBGCard contract.
type ZBGCardFaucetToggled struct {
	IsApproved bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterFaucetToggled is a free log retrieval operation binding the contract event 0x87e7975983dfa9eadba810a775cce04da400ce8dcc3f201b0e2def605d58911d.
//
// Solidity: e FaucetToggled(isApproved bool)
func (_ZBGCard *ZBGCardFilterer) FilterFaucetToggled(opts *bind.FilterOpts) (*ZBGCardFaucetToggledIterator, error) {

	logs, sub, err := _ZBGCard.contract.FilterLogs(opts, "FaucetToggled")
	if err != nil {
		return nil, err
	}
	return &ZBGCardFaucetToggledIterator{contract: _ZBGCard.contract, event: "FaucetToggled", logs: logs, sub: sub}, nil
}

// WatchFaucetToggled is a free log subscription operation binding the contract event 0x87e7975983dfa9eadba810a775cce04da400ce8dcc3f201b0e2def605d58911d.
//
// Solidity: e FaucetToggled(isApproved bool)
func (_ZBGCard *ZBGCardFilterer) WatchFaucetToggled(opts *bind.WatchOpts, sink chan<- *ZBGCardFaucetToggled) (event.Subscription, error) {

	logs, sub, err := _ZBGCard.contract.WatchLogs(opts, "FaucetToggled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZBGCardFaucetToggled)
				if err := _ZBGCard.contract.UnpackLog(event, "FaucetToggled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ZBGCardOwnershipRenouncedIterator is returned from FilterOwnershipRenounced and is used to iterate over the raw logs and unpacked data for OwnershipRenounced events raised by the ZBGCard contract.
type ZBGCardOwnershipRenouncedIterator struct {
	Event *ZBGCardOwnershipRenounced // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ZBGCardOwnershipRenouncedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZBGCardOwnershipRenounced)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ZBGCardOwnershipRenounced)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ZBGCardOwnershipRenouncedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZBGCardOwnershipRenouncedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZBGCardOwnershipRenounced represents a OwnershipRenounced event raised by the ZBGCard contract.
type ZBGCardOwnershipRenounced struct {
	PreviousOwner common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipRenounced is a free log retrieval operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_ZBGCard *ZBGCardFilterer) FilterOwnershipRenounced(opts *bind.FilterOpts, previousOwner []common.Address) (*ZBGCardOwnershipRenouncedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _ZBGCard.contract.FilterLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ZBGCardOwnershipRenouncedIterator{contract: _ZBGCard.contract, event: "OwnershipRenounced", logs: logs, sub: sub}, nil
}

// WatchOwnershipRenounced is a free log subscription operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_ZBGCard *ZBGCardFilterer) WatchOwnershipRenounced(opts *bind.WatchOpts, sink chan<- *ZBGCardOwnershipRenounced, previousOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _ZBGCard.contract.WatchLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZBGCardOwnershipRenounced)
				if err := _ZBGCard.contract.UnpackLog(event, "OwnershipRenounced", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ZBGCardOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ZBGCard contract.
type ZBGCardOwnershipTransferredIterator struct {
	Event *ZBGCardOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ZBGCardOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZBGCardOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ZBGCardOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ZBGCardOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZBGCardOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZBGCardOwnershipTransferred represents a OwnershipTransferred event raised by the ZBGCard contract.
type ZBGCardOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_ZBGCard *ZBGCardFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ZBGCardOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ZBGCard.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ZBGCardOwnershipTransferredIterator{contract: _ZBGCard.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_ZBGCard *ZBGCardFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ZBGCardOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ZBGCard.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZBGCardOwnershipTransferred)
				if err := _ZBGCard.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ZBGCardRemovedValidatorIterator is returned from FilterRemovedValidator and is used to iterate over the raw logs and unpacked data for RemovedValidator events raised by the ZBGCard contract.
type ZBGCardRemovedValidatorIterator struct {
	Event *ZBGCardRemovedValidator // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ZBGCardRemovedValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZBGCardRemovedValidator)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ZBGCardRemovedValidator)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ZBGCardRemovedValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZBGCardRemovedValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZBGCardRemovedValidator represents a RemovedValidator event raised by the ZBGCard contract.
type ZBGCardRemovedValidator struct {
	Validator common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRemovedValidator is a free log retrieval operation binding the contract event 0xb625c55cf7e37b54fcd18bc4edafdf3f4f9acd59a5ec824c77c795dcb2d65070.
//
// Solidity: e RemovedValidator(validator address)
func (_ZBGCard *ZBGCardFilterer) FilterRemovedValidator(opts *bind.FilterOpts) (*ZBGCardRemovedValidatorIterator, error) {

	logs, sub, err := _ZBGCard.contract.FilterLogs(opts, "RemovedValidator")
	if err != nil {
		return nil, err
	}
	return &ZBGCardRemovedValidatorIterator{contract: _ZBGCard.contract, event: "RemovedValidator", logs: logs, sub: sub}, nil
}

// WatchRemovedValidator is a free log subscription operation binding the contract event 0xb625c55cf7e37b54fcd18bc4edafdf3f4f9acd59a5ec824c77c795dcb2d65070.
//
// Solidity: e RemovedValidator(validator address)
func (_ZBGCard *ZBGCardFilterer) WatchRemovedValidator(opts *bind.WatchOpts, sink chan<- *ZBGCardRemovedValidator) (event.Subscription, error) {

	logs, sub, err := _ZBGCard.contract.WatchLogs(opts, "RemovedValidator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZBGCardRemovedValidator)
				if err := _ZBGCard.contract.UnpackLog(event, "RemovedValidator", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ZBGCardTokenClaimedIterator is returned from FilterTokenClaimed and is used to iterate over the raw logs and unpacked data for TokenClaimed events raised by the ZBGCard contract.
type ZBGCardTokenClaimedIterator struct {
	Event *ZBGCardTokenClaimed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ZBGCardTokenClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZBGCardTokenClaimed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ZBGCardTokenClaimed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ZBGCardTokenClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZBGCardTokenClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZBGCardTokenClaimed represents a TokenClaimed event raised by the ZBGCard contract.
type ZBGCardTokenClaimed struct {
	TokenId *big.Int
	Claimer common.Address
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTokenClaimed is a free log retrieval operation binding the contract event 0x7ebec69a24f70f6adc732c5d495ea40faf1248f959f11feafe0f9bbdc4e07b5a.
//
// Solidity: e TokenClaimed(tokenId indexed uint256, claimer address, amount uint256)
func (_ZBGCard *ZBGCardFilterer) FilterTokenClaimed(opts *bind.FilterOpts, tokenId []*big.Int) (*ZBGCardTokenClaimedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ZBGCard.contract.FilterLogs(opts, "TokenClaimed", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &ZBGCardTokenClaimedIterator{contract: _ZBGCard.contract, event: "TokenClaimed", logs: logs, sub: sub}, nil
}

// WatchTokenClaimed is a free log subscription operation binding the contract event 0x7ebec69a24f70f6adc732c5d495ea40faf1248f959f11feafe0f9bbdc4e07b5a.
//
// Solidity: e TokenClaimed(tokenId indexed uint256, claimer address, amount uint256)
func (_ZBGCard *ZBGCardFilterer) WatchTokenClaimed(opts *bind.WatchOpts, sink chan<- *ZBGCardTokenClaimed, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ZBGCard.contract.WatchLogs(opts, "TokenClaimed", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZBGCardTokenClaimed)
				if err := _ZBGCard.contract.UnpackLog(event, "TokenClaimed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ZBGCardTokenMintedIterator is returned from FilterTokenMinted and is used to iterate over the raw logs and unpacked data for TokenMinted events raised by the ZBGCard contract.
type ZBGCardTokenMintedIterator struct {
	Event *ZBGCardTokenMinted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ZBGCardTokenMintedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZBGCardTokenMinted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ZBGCardTokenMinted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ZBGCardTokenMintedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZBGCardTokenMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZBGCardTokenMinted represents a TokenMinted event raised by the ZBGCard contract.
type ZBGCardTokenMinted struct {
	TokenId *big.Int
	Supply  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTokenMinted is a free log retrieval operation binding the contract event 0x5f7666687319b40936f33c188908d86aea154abd3f4127b4fa0a3f04f303c7da.
//
// Solidity: e TokenMinted(tokenId uint256, supply uint256)
func (_ZBGCard *ZBGCardFilterer) FilterTokenMinted(opts *bind.FilterOpts) (*ZBGCardTokenMintedIterator, error) {

	logs, sub, err := _ZBGCard.contract.FilterLogs(opts, "TokenMinted")
	if err != nil {
		return nil, err
	}
	return &ZBGCardTokenMintedIterator{contract: _ZBGCard.contract, event: "TokenMinted", logs: logs, sub: sub}, nil
}

// WatchTokenMinted is a free log subscription operation binding the contract event 0x5f7666687319b40936f33c188908d86aea154abd3f4127b4fa0a3f04f303c7da.
//
// Solidity: e TokenMinted(tokenId uint256, supply uint256)
func (_ZBGCard *ZBGCardFilterer) WatchTokenMinted(opts *bind.WatchOpts, sink chan<- *ZBGCardTokenMinted) (event.Subscription, error) {

	logs, sub, err := _ZBGCard.contract.WatchLogs(opts, "TokenMinted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZBGCardTokenMinted)
				if err := _ZBGCard.contract.UnpackLog(event, "TokenMinted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ZBGCardTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ZBGCard contract.
type ZBGCardTransferIterator struct {
	Event *ZBGCardTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ZBGCardTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZBGCardTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ZBGCardTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ZBGCardTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZBGCardTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZBGCardTransfer represents a Transfer event raised by the ZBGCard contract.
type ZBGCardTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(_from indexed address, _to indexed address, _tokenId indexed uint256)
func (_ZBGCard *ZBGCardFilterer) FilterTransfer(opts *bind.FilterOpts, _from []common.Address, _to []common.Address, _tokenId []*big.Int) (*ZBGCardTransferIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}
	var _tokenIdRule []interface{}
	for _, _tokenIdItem := range _tokenId {
		_tokenIdRule = append(_tokenIdRule, _tokenIdItem)
	}

	logs, sub, err := _ZBGCard.contract.FilterLogs(opts, "Transfer", _fromRule, _toRule, _tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &ZBGCardTransferIterator{contract: _ZBGCard.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(_from indexed address, _to indexed address, _tokenId indexed uint256)
func (_ZBGCard *ZBGCardFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ZBGCardTransfer, _from []common.Address, _to []common.Address, _tokenId []*big.Int) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}
	var _tokenIdRule []interface{}
	for _, _tokenIdItem := range _tokenId {
		_tokenIdRule = append(_tokenIdRule, _tokenIdItem)
	}

	logs, sub, err := _ZBGCard.contract.WatchLogs(opts, "Transfer", _fromRule, _toRule, _tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZBGCardTransfer)
				if err := _ZBGCard.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ZBGCardTransferWithQuantityIterator is returned from FilterTransferWithQuantity and is used to iterate over the raw logs and unpacked data for TransferWithQuantity events raised by the ZBGCard contract.
type ZBGCardTransferWithQuantityIterator struct {
	Event *ZBGCardTransferWithQuantity // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ZBGCardTransferWithQuantityIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZBGCardTransferWithQuantity)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ZBGCardTransferWithQuantity)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ZBGCardTransferWithQuantityIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZBGCardTransferWithQuantityIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZBGCardTransferWithQuantity represents a TransferWithQuantity event raised by the ZBGCard contract.
type ZBGCardTransferWithQuantity struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransferWithQuantity is a free log retrieval operation binding the contract event 0x2114851a3e2a54429989f46c1ab0743e37ded205d9bbdfd85635aed5bd595a06.
//
// Solidity: e TransferWithQuantity(from address, to address, tokenId indexed uint256, amount uint256)
func (_ZBGCard *ZBGCardFilterer) FilterTransferWithQuantity(opts *bind.FilterOpts, tokenId []*big.Int) (*ZBGCardTransferWithQuantityIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ZBGCard.contract.FilterLogs(opts, "TransferWithQuantity", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &ZBGCardTransferWithQuantityIterator{contract: _ZBGCard.contract, event: "TransferWithQuantity", logs: logs, sub: sub}, nil
}

// WatchTransferWithQuantity is a free log subscription operation binding the contract event 0x2114851a3e2a54429989f46c1ab0743e37ded205d9bbdfd85635aed5bd595a06.
//
// Solidity: e TransferWithQuantity(from address, to address, tokenId indexed uint256, amount uint256)
func (_ZBGCard *ZBGCardFilterer) WatchTransferWithQuantity(opts *bind.WatchOpts, sink chan<- *ZBGCardTransferWithQuantity, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ZBGCard.contract.WatchLogs(opts, "TransferWithQuantity", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZBGCardTransferWithQuantity)
				if err := _ZBGCard.contract.UnpackLog(event, "TransferWithQuantity", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
