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

// SimpleStoreContractABI is the input ABI used to generate the binding from.
const SimpleStoreContractABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"NewValueSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Echo\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"set\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"get\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// SimpleStoreContract is an auto generated Go binding around an Ethereum contract.
type SimpleStoreContract struct {
	SimpleStoreContractCaller     // Read-only binding to the contract
	SimpleStoreContractTransactor // Write-only binding to the contract
	SimpleStoreContractFilterer   // Log filterer for contract events
}

// SimpleStoreContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type SimpleStoreContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleStoreContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SimpleStoreContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleStoreContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SimpleStoreContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleStoreContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SimpleStoreContractSession struct {
	Contract     *SimpleStoreContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// SimpleStoreContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SimpleStoreContractCallerSession struct {
	Contract *SimpleStoreContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// SimpleStoreContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SimpleStoreContractTransactorSession struct {
	Contract     *SimpleStoreContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// SimpleStoreContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type SimpleStoreContractRaw struct {
	Contract *SimpleStoreContract // Generic contract binding to access the raw methods on
}

// SimpleStoreContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SimpleStoreContractCallerRaw struct {
	Contract *SimpleStoreContractCaller // Generic read-only contract binding to access the raw methods on
}

// SimpleStoreContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SimpleStoreContractTransactorRaw struct {
	Contract *SimpleStoreContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSimpleStoreContract creates a new instance of SimpleStoreContract, bound to a specific deployed contract.
func NewSimpleStoreContract(address common.Address, backend bind.ContractBackend) (*SimpleStoreContract, error) {
	contract, err := bindSimpleStoreContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SimpleStoreContract{SimpleStoreContractCaller: SimpleStoreContractCaller{contract: contract}, SimpleStoreContractTransactor: SimpleStoreContractTransactor{contract: contract}, SimpleStoreContractFilterer: SimpleStoreContractFilterer{contract: contract}}, nil
}

// NewSimpleStoreContractCaller creates a new read-only instance of SimpleStoreContract, bound to a specific deployed contract.
func NewSimpleStoreContractCaller(address common.Address, caller bind.ContractCaller) (*SimpleStoreContractCaller, error) {
	contract, err := bindSimpleStoreContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleStoreContractCaller{contract: contract}, nil
}

// NewSimpleStoreContractTransactor creates a new write-only instance of SimpleStoreContract, bound to a specific deployed contract.
func NewSimpleStoreContractTransactor(address common.Address, transactor bind.ContractTransactor) (*SimpleStoreContractTransactor, error) {
	contract, err := bindSimpleStoreContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleStoreContractTransactor{contract: contract}, nil
}

// NewSimpleStoreContractFilterer creates a new log filterer instance of SimpleStoreContract, bound to a specific deployed contract.
func NewSimpleStoreContractFilterer(address common.Address, filterer bind.ContractFilterer) (*SimpleStoreContractFilterer, error) {
	contract, err := bindSimpleStoreContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SimpleStoreContractFilterer{contract: contract}, nil
}

// bindSimpleStoreContract binds a generic wrapper to an already deployed contract.
func bindSimpleStoreContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SimpleStoreContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimpleStoreContract *SimpleStoreContractRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SimpleStoreContract.Contract.SimpleStoreContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimpleStoreContract *SimpleStoreContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleStoreContract.Contract.SimpleStoreContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimpleStoreContract *SimpleStoreContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleStoreContract.Contract.SimpleStoreContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimpleStoreContract *SimpleStoreContractCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SimpleStoreContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimpleStoreContract *SimpleStoreContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleStoreContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimpleStoreContract *SimpleStoreContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleStoreContract.Contract.contract.Transact(opts, method, params...)
}

// Get is a free data retrieval call binding the contract method 0x6d4ce63c.
//
// Solidity: function get() constant returns(uint256)
func (_SimpleStoreContract *SimpleStoreContractCaller) Get(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SimpleStoreContract.contract.Call(opts, out, "get")
	return *ret0, err
}

// Get is a free data retrieval call binding the contract method 0x6d4ce63c.
//
// Solidity: function get() constant returns(uint256)
func (_SimpleStoreContract *SimpleStoreContractSession) Get() (*big.Int, error) {
	return _SimpleStoreContract.Contract.Get(&_SimpleStoreContract.CallOpts)
}

// Get is a free data retrieval call binding the contract method 0x6d4ce63c.
//
// Solidity: function get() constant returns(uint256)
func (_SimpleStoreContract *SimpleStoreContractCallerSession) Get() (*big.Int, error) {
	return _SimpleStoreContract.Contract.Get(&_SimpleStoreContract.CallOpts)
}

// Set is a paid mutator transaction binding the contract method 0x60fe47b1.
//
// Solidity: function set(_value uint256) returns()
func (_SimpleStoreContract *SimpleStoreContractTransactor) Set(opts *bind.TransactOpts, _value *big.Int) (*types.Transaction, error) {
	return _SimpleStoreContract.contract.Transact(opts, "set", _value)
}

// Set is a paid mutator transaction binding the contract method 0x60fe47b1.
//
// Solidity: function set(_value uint256) returns()
func (_SimpleStoreContract *SimpleStoreContractSession) Set(_value *big.Int) (*types.Transaction, error) {
	return _SimpleStoreContract.Contract.Set(&_SimpleStoreContract.TransactOpts, _value)
}

// Set is a paid mutator transaction binding the contract method 0x60fe47b1.
//
// Solidity: function set(_value uint256) returns()
func (_SimpleStoreContract *SimpleStoreContractTransactorSession) Set(_value *big.Int) (*types.Transaction, error) {
	return _SimpleStoreContract.Contract.Set(&_SimpleStoreContract.TransactOpts, _value)
}

// SimpleStoreContractEchoIterator is returned from FilterEcho and is used to iterate over the raw logs and unpacked data for Echo events raised by the SimpleStoreContract contract.
type SimpleStoreContractEchoIterator struct {
	Event *SimpleStoreContractEcho // Event containing the contract specifics and raw log

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
func (it *SimpleStoreContractEchoIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleStoreContractEcho)
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
		it.Event = new(SimpleStoreContractEcho)
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
func (it *SimpleStoreContractEchoIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleStoreContractEchoIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleStoreContractEcho represents a Echo event raised by the SimpleStoreContract contract.
type SimpleStoreContractEcho struct {
	Name  common.Hash
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterEcho is a free log retrieval operation binding the contract event 0xcaf432cb38a3a6f6c9bdd5b57f1a5388e0f452215b40290524727e0a7a523da1.
//
// Solidity: e Echo(name indexed string, _value uint256)
func (_SimpleStoreContract *SimpleStoreContractFilterer) FilterEcho(opts *bind.FilterOpts, name []string) (*SimpleStoreContractEchoIterator, error) {

	var nameRule []interface{}
	for _, nameItem := range name {
		nameRule = append(nameRule, nameItem)
	}

	logs, sub, err := _SimpleStoreContract.contract.FilterLogs(opts, "Echo", nameRule)
	if err != nil {
		return nil, err
	}
	return &SimpleStoreContractEchoIterator{contract: _SimpleStoreContract.contract, event: "Echo", logs: logs, sub: sub}, nil
}

// WatchEcho is a free log subscription operation binding the contract event 0xcaf432cb38a3a6f6c9bdd5b57f1a5388e0f452215b40290524727e0a7a523da1.
//
// Solidity: e Echo(name indexed string, _value uint256)
func (_SimpleStoreContract *SimpleStoreContractFilterer) WatchEcho(opts *bind.WatchOpts, sink chan<- *SimpleStoreContractEcho, name []string) (event.Subscription, error) {

	var nameRule []interface{}
	for _, nameItem := range name {
		nameRule = append(nameRule, nameItem)
	}

	logs, sub, err := _SimpleStoreContract.contract.WatchLogs(opts, "Echo", nameRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimpleStoreContractEcho)
				if err := _SimpleStoreContract.contract.UnpackLog(event, "Echo", log); err != nil {
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

// SimpleStoreContractNewValueSetIterator is returned from FilterNewValueSet and is used to iterate over the raw logs and unpacked data for NewValueSet events raised by the SimpleStoreContract contract.
type SimpleStoreContractNewValueSetIterator struct {
	Event *SimpleStoreContractNewValueSet // Event containing the contract specifics and raw log

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
func (it *SimpleStoreContractNewValueSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleStoreContractNewValueSet)
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
		it.Event = new(SimpleStoreContractNewValueSet)
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
func (it *SimpleStoreContractNewValueSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleStoreContractNewValueSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleStoreContractNewValueSet represents a NewValueSet event raised by the SimpleStoreContract contract.
type SimpleStoreContractNewValueSet struct {
	Name  string
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterNewValueSet is a free log retrieval operation binding the contract event 0x7e8ee33c8615178a01c0dbd6263ef0af255dcaba4dde4f387384567abfab718f.
//
// Solidity: e NewValueSet(name string, _value uint256)
func (_SimpleStoreContract *SimpleStoreContractFilterer) FilterNewValueSet(opts *bind.FilterOpts) (*SimpleStoreContractNewValueSetIterator, error) {

	logs, sub, err := _SimpleStoreContract.contract.FilterLogs(opts, "NewValueSet")
	if err != nil {
		return nil, err
	}
	return &SimpleStoreContractNewValueSetIterator{contract: _SimpleStoreContract.contract, event: "NewValueSet", logs: logs, sub: sub}, nil
}

// WatchNewValueSet is a free log subscription operation binding the contract event 0x7e8ee33c8615178a01c0dbd6263ef0af255dcaba4dde4f387384567abfab718f.
//
// Solidity: e NewValueSet(name string, _value uint256)
func (_SimpleStoreContract *SimpleStoreContractFilterer) WatchNewValueSet(opts *bind.WatchOpts, sink chan<- *SimpleStoreContractNewValueSet) (event.Subscription, error) {

	logs, sub, err := _SimpleStoreContract.contract.WatchLogs(opts, "NewValueSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimpleStoreContractNewValueSet)
				if err := _SimpleStoreContract.contract.UnpackLog(event, "NewValueSet", log); err != nil {
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
