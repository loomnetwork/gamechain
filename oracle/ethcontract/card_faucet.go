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

// CardFaucetABI is the input ABI used to generate the binding from.
const CardFaucetABI = "[{\"inputs\":[{\"name\":\"creator\",\"type\":\"address\"},{\"name\":\"cardAddress\",\"type\":\"address\"},{\"name\":\"boosterPackAddr\",\"type\":\"address\"},{\"name\":\"superPackAddr\",\"type\":\"address\"},{\"name\":\"airPackAddr\",\"type\":\"address\"},{\"name\":\"earthPackAddr\",\"type\":\"address\"},{\"name\":\"firePackAddr\",\"type\":\"address\"},{\"name\":\"lifePackAddr\",\"type\":\"address\"},{\"name\":\"toxicPackAddr\",\"type\":\"address\"},{\"name\":\"waterPackAddr\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"cardId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"boosterType\",\"type\":\"uint8\"}],\"name\":\"GeneratedCard\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"cardId\",\"type\":\"uint256\"}],\"name\":\"UpgradedCardToBE\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"cardId\",\"type\":\"uint256\"}],\"name\":\"UpgradedCardToLE\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"target\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"total\",\"type\":\"uint256\"}],\"name\":\"LimitedDifficultyChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"target\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"total\",\"type\":\"uint256\"}],\"name\":\"BackerDifficultyChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"boosterTypes\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"rarity\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"cardList\",\"type\":\"uint256[]\"}],\"name\":\"DropratesSet\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"boosterType\",\"type\":\"uint256\"},{\"name\":\"rarity\",\"type\":\"uint256\"},{\"name\":\"cardList\",\"type\":\"uint256[]\"}],\"name\":\"setupDroprates\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"boosterType\",\"type\":\"uint256\"},{\"name\":\"rarity\",\"type\":\"uint256\"}],\"name\":\"getDroprates\",\"outputs\":[{\"name\":\"cardList\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"boosterType\",\"type\":\"uint8\"}],\"name\":\"openBoosterPack\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_target\",\"type\":\"uint256\"},{\"name\":\"_total\",\"type\":\"uint256\"}],\"name\":\"updateDifficultyBE\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_target\",\"type\":\"uint256\"},{\"name\":\"_total\",\"type\":\"uint256\"}],\"name\":\"updateDifficultyLE\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"endBackersPeriod\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newValidator\",\"type\":\"address\"}],\"name\":\"addValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// CardFaucet is an auto generated Go binding around an Ethereum contract.
type CardFaucet struct {
	CardFaucetCaller     // Read-only binding to the contract
	CardFaucetTransactor // Write-only binding to the contract
	CardFaucetFilterer   // Log filterer for contract events
}

// CardFaucetCaller is an auto generated read-only Go binding around an Ethereum contract.
type CardFaucetCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CardFaucetTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CardFaucetTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CardFaucetFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CardFaucetFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CardFaucetSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CardFaucetSession struct {
	Contract     *CardFaucet       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CardFaucetCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CardFaucetCallerSession struct {
	Contract *CardFaucetCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// CardFaucetTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CardFaucetTransactorSession struct {
	Contract     *CardFaucetTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// CardFaucetRaw is an auto generated low-level Go binding around an Ethereum contract.
type CardFaucetRaw struct {
	Contract *CardFaucet // Generic contract binding to access the raw methods on
}

// CardFaucetCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CardFaucetCallerRaw struct {
	Contract *CardFaucetCaller // Generic read-only contract binding to access the raw methods on
}

// CardFaucetTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CardFaucetTransactorRaw struct {
	Contract *CardFaucetTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCardFaucet creates a new instance of CardFaucet, bound to a specific deployed contract.
func NewCardFaucet(address common.Address, backend bind.ContractBackend) (*CardFaucet, error) {
	contract, err := bindCardFaucet(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CardFaucet{CardFaucetCaller: CardFaucetCaller{contract: contract}, CardFaucetTransactor: CardFaucetTransactor{contract: contract}, CardFaucetFilterer: CardFaucetFilterer{contract: contract}}, nil
}

// NewCardFaucetCaller creates a new read-only instance of CardFaucet, bound to a specific deployed contract.
func NewCardFaucetCaller(address common.Address, caller bind.ContractCaller) (*CardFaucetCaller, error) {
	contract, err := bindCardFaucet(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CardFaucetCaller{contract: contract}, nil
}

// NewCardFaucetTransactor creates a new write-only instance of CardFaucet, bound to a specific deployed contract.
func NewCardFaucetTransactor(address common.Address, transactor bind.ContractTransactor) (*CardFaucetTransactor, error) {
	contract, err := bindCardFaucet(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CardFaucetTransactor{contract: contract}, nil
}

// NewCardFaucetFilterer creates a new log filterer instance of CardFaucet, bound to a specific deployed contract.
func NewCardFaucetFilterer(address common.Address, filterer bind.ContractFilterer) (*CardFaucetFilterer, error) {
	contract, err := bindCardFaucet(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CardFaucetFilterer{contract: contract}, nil
}

// bindCardFaucet binds a generic wrapper to an already deployed contract.
func bindCardFaucet(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CardFaucetABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CardFaucet *CardFaucetRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _CardFaucet.Contract.CardFaucetCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CardFaucet *CardFaucetRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CardFaucet.Contract.CardFaucetTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CardFaucet *CardFaucetRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CardFaucet.Contract.CardFaucetTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CardFaucet *CardFaucetCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _CardFaucet.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CardFaucet *CardFaucetTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CardFaucet.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CardFaucet *CardFaucetTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CardFaucet.Contract.contract.Transact(opts, method, params...)
}

// GetDroprates is a free data retrieval call binding the contract method 0x54b156fb.
//
// Solidity: function getDroprates(boosterType uint256, rarity uint256) constant returns(cardList uint256[])
func (_CardFaucet *CardFaucetCaller) GetDroprates(opts *bind.CallOpts, boosterType *big.Int, rarity *big.Int) ([]*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
	)
	out := ret0
	err := _CardFaucet.contract.Call(opts, out, "getDroprates", boosterType, rarity)
	return *ret0, err
}

// GetDroprates is a free data retrieval call binding the contract method 0x54b156fb.
//
// Solidity: function getDroprates(boosterType uint256, rarity uint256) constant returns(cardList uint256[])
func (_CardFaucet *CardFaucetSession) GetDroprates(boosterType *big.Int, rarity *big.Int) ([]*big.Int, error) {
	return _CardFaucet.Contract.GetDroprates(&_CardFaucet.CallOpts, boosterType, rarity)
}

// GetDroprates is a free data retrieval call binding the contract method 0x54b156fb.
//
// Solidity: function getDroprates(boosterType uint256, rarity uint256) constant returns(cardList uint256[])
func (_CardFaucet *CardFaucetCallerSession) GetDroprates(boosterType *big.Int, rarity *big.Int) ([]*big.Int, error) {
	return _CardFaucet.Contract.GetDroprates(&_CardFaucet.CallOpts, boosterType, rarity)
}

// AddValidator is a paid mutator transaction binding the contract method 0x4d238c8e.
//
// Solidity: function addValidator(newValidator address) returns()
func (_CardFaucet *CardFaucetTransactor) AddValidator(opts *bind.TransactOpts, newValidator common.Address) (*types.Transaction, error) {
	return _CardFaucet.contract.Transact(opts, "addValidator", newValidator)
}

// AddValidator is a paid mutator transaction binding the contract method 0x4d238c8e.
//
// Solidity: function addValidator(newValidator address) returns()
func (_CardFaucet *CardFaucetSession) AddValidator(newValidator common.Address) (*types.Transaction, error) {
	return _CardFaucet.Contract.AddValidator(&_CardFaucet.TransactOpts, newValidator)
}

// AddValidator is a paid mutator transaction binding the contract method 0x4d238c8e.
//
// Solidity: function addValidator(newValidator address) returns()
func (_CardFaucet *CardFaucetTransactorSession) AddValidator(newValidator common.Address) (*types.Transaction, error) {
	return _CardFaucet.Contract.AddValidator(&_CardFaucet.TransactOpts, newValidator)
}

// EndBackersPeriod is a paid mutator transaction binding the contract method 0x39210452.
//
// Solidity: function endBackersPeriod() returns()
func (_CardFaucet *CardFaucetTransactor) EndBackersPeriod(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CardFaucet.contract.Transact(opts, "endBackersPeriod")
}

// EndBackersPeriod is a paid mutator transaction binding the contract method 0x39210452.
//
// Solidity: function endBackersPeriod() returns()
func (_CardFaucet *CardFaucetSession) EndBackersPeriod() (*types.Transaction, error) {
	return _CardFaucet.Contract.EndBackersPeriod(&_CardFaucet.TransactOpts)
}

// EndBackersPeriod is a paid mutator transaction binding the contract method 0x39210452.
//
// Solidity: function endBackersPeriod() returns()
func (_CardFaucet *CardFaucetTransactorSession) EndBackersPeriod() (*types.Transaction, error) {
	return _CardFaucet.Contract.EndBackersPeriod(&_CardFaucet.TransactOpts)
}

// OpenBoosterPack is a paid mutator transaction binding the contract method 0x576de247.
//
// Solidity: function openBoosterPack(boosterType uint8) returns()
func (_CardFaucet *CardFaucetTransactor) OpenBoosterPack(opts *bind.TransactOpts, boosterType uint8) (*types.Transaction, error) {
	return _CardFaucet.contract.Transact(opts, "openBoosterPack", boosterType)
}

// OpenBoosterPack is a paid mutator transaction binding the contract method 0x576de247.
//
// Solidity: function openBoosterPack(boosterType uint8) returns()
func (_CardFaucet *CardFaucetSession) OpenBoosterPack(boosterType uint8) (*types.Transaction, error) {
	return _CardFaucet.Contract.OpenBoosterPack(&_CardFaucet.TransactOpts, boosterType)
}

// OpenBoosterPack is a paid mutator transaction binding the contract method 0x576de247.
//
// Solidity: function openBoosterPack(boosterType uint8) returns()
func (_CardFaucet *CardFaucetTransactorSession) OpenBoosterPack(boosterType uint8) (*types.Transaction, error) {
	return _CardFaucet.Contract.OpenBoosterPack(&_CardFaucet.TransactOpts, boosterType)
}

// SetupDroprates is a paid mutator transaction binding the contract method 0xa70afe6c.
//
// Solidity: function setupDroprates(boosterType uint256, rarity uint256, cardList uint256[]) returns()
func (_CardFaucet *CardFaucetTransactor) SetupDroprates(opts *bind.TransactOpts, boosterType *big.Int, rarity *big.Int, cardList []*big.Int) (*types.Transaction, error) {
	return _CardFaucet.contract.Transact(opts, "setupDroprates", boosterType, rarity, cardList)
}

// SetupDroprates is a paid mutator transaction binding the contract method 0xa70afe6c.
//
// Solidity: function setupDroprates(boosterType uint256, rarity uint256, cardList uint256[]) returns()
func (_CardFaucet *CardFaucetSession) SetupDroprates(boosterType *big.Int, rarity *big.Int, cardList []*big.Int) (*types.Transaction, error) {
	return _CardFaucet.Contract.SetupDroprates(&_CardFaucet.TransactOpts, boosterType, rarity, cardList)
}

// SetupDroprates is a paid mutator transaction binding the contract method 0xa70afe6c.
//
// Solidity: function setupDroprates(boosterType uint256, rarity uint256, cardList uint256[]) returns()
func (_CardFaucet *CardFaucetTransactorSession) SetupDroprates(boosterType *big.Int, rarity *big.Int, cardList []*big.Int) (*types.Transaction, error) {
	return _CardFaucet.Contract.SetupDroprates(&_CardFaucet.TransactOpts, boosterType, rarity, cardList)
}

// UpdateDifficultyBE is a paid mutator transaction binding the contract method 0x9770993f.
//
// Solidity: function updateDifficultyBE(_target uint256, _total uint256) returns()
func (_CardFaucet *CardFaucetTransactor) UpdateDifficultyBE(opts *bind.TransactOpts, _target *big.Int, _total *big.Int) (*types.Transaction, error) {
	return _CardFaucet.contract.Transact(opts, "updateDifficultyBE", _target, _total)
}

// UpdateDifficultyBE is a paid mutator transaction binding the contract method 0x9770993f.
//
// Solidity: function updateDifficultyBE(_target uint256, _total uint256) returns()
func (_CardFaucet *CardFaucetSession) UpdateDifficultyBE(_target *big.Int, _total *big.Int) (*types.Transaction, error) {
	return _CardFaucet.Contract.UpdateDifficultyBE(&_CardFaucet.TransactOpts, _target, _total)
}

// UpdateDifficultyBE is a paid mutator transaction binding the contract method 0x9770993f.
//
// Solidity: function updateDifficultyBE(_target uint256, _total uint256) returns()
func (_CardFaucet *CardFaucetTransactorSession) UpdateDifficultyBE(_target *big.Int, _total *big.Int) (*types.Transaction, error) {
	return _CardFaucet.Contract.UpdateDifficultyBE(&_CardFaucet.TransactOpts, _target, _total)
}

// UpdateDifficultyLE is a paid mutator transaction binding the contract method 0xe6fa7717.
//
// Solidity: function updateDifficultyLE(_target uint256, _total uint256) returns()
func (_CardFaucet *CardFaucetTransactor) UpdateDifficultyLE(opts *bind.TransactOpts, _target *big.Int, _total *big.Int) (*types.Transaction, error) {
	return _CardFaucet.contract.Transact(opts, "updateDifficultyLE", _target, _total)
}

// UpdateDifficultyLE is a paid mutator transaction binding the contract method 0xe6fa7717.
//
// Solidity: function updateDifficultyLE(_target uint256, _total uint256) returns()
func (_CardFaucet *CardFaucetSession) UpdateDifficultyLE(_target *big.Int, _total *big.Int) (*types.Transaction, error) {
	return _CardFaucet.Contract.UpdateDifficultyLE(&_CardFaucet.TransactOpts, _target, _total)
}

// UpdateDifficultyLE is a paid mutator transaction binding the contract method 0xe6fa7717.
//
// Solidity: function updateDifficultyLE(_target uint256, _total uint256) returns()
func (_CardFaucet *CardFaucetTransactorSession) UpdateDifficultyLE(_target *big.Int, _total *big.Int) (*types.Transaction, error) {
	return _CardFaucet.Contract.UpdateDifficultyLE(&_CardFaucet.TransactOpts, _target, _total)
}

// CardFaucetBackerDifficultyChangedIterator is returned from FilterBackerDifficultyChanged and is used to iterate over the raw logs and unpacked data for BackerDifficultyChanged events raised by the CardFaucet contract.
type CardFaucetBackerDifficultyChangedIterator struct {
	Event *CardFaucetBackerDifficultyChanged // Event containing the contract specifics and raw log

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
func (it *CardFaucetBackerDifficultyChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CardFaucetBackerDifficultyChanged)
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
		it.Event = new(CardFaucetBackerDifficultyChanged)
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
func (it *CardFaucetBackerDifficultyChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CardFaucetBackerDifficultyChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CardFaucetBackerDifficultyChanged represents a BackerDifficultyChanged event raised by the CardFaucet contract.
type CardFaucetBackerDifficultyChanged struct {
	Target *big.Int
	Total  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBackerDifficultyChanged is a free log retrieval operation binding the contract event 0xc39d86a0778a5e4422276c5eb2b4e7879a2bf8f5d90c5df2486e3fe3a2dc5695.
//
// Solidity: e BackerDifficultyChanged(target uint256, total uint256)
func (_CardFaucet *CardFaucetFilterer) FilterBackerDifficultyChanged(opts *bind.FilterOpts) (*CardFaucetBackerDifficultyChangedIterator, error) {

	logs, sub, err := _CardFaucet.contract.FilterLogs(opts, "BackerDifficultyChanged")
	if err != nil {
		return nil, err
	}
	return &CardFaucetBackerDifficultyChangedIterator{contract: _CardFaucet.contract, event: "BackerDifficultyChanged", logs: logs, sub: sub}, nil
}

// WatchBackerDifficultyChanged is a free log subscription operation binding the contract event 0xc39d86a0778a5e4422276c5eb2b4e7879a2bf8f5d90c5df2486e3fe3a2dc5695.
//
// Solidity: e BackerDifficultyChanged(target uint256, total uint256)
func (_CardFaucet *CardFaucetFilterer) WatchBackerDifficultyChanged(opts *bind.WatchOpts, sink chan<- *CardFaucetBackerDifficultyChanged) (event.Subscription, error) {

	logs, sub, err := _CardFaucet.contract.WatchLogs(opts, "BackerDifficultyChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CardFaucetBackerDifficultyChanged)
				if err := _CardFaucet.contract.UnpackLog(event, "BackerDifficultyChanged", log); err != nil {
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

// CardFaucetDropratesSetIterator is returned from FilterDropratesSet and is used to iterate over the raw logs and unpacked data for DropratesSet events raised by the CardFaucet contract.
type CardFaucetDropratesSetIterator struct {
	Event *CardFaucetDropratesSet // Event containing the contract specifics and raw log

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
func (it *CardFaucetDropratesSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CardFaucetDropratesSet)
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
		it.Event = new(CardFaucetDropratesSet)
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
func (it *CardFaucetDropratesSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CardFaucetDropratesSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CardFaucetDropratesSet represents a DropratesSet event raised by the CardFaucet contract.
type CardFaucetDropratesSet struct {
	BoosterTypes *big.Int
	Rarity       *big.Int
	CardList     []*big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterDropratesSet is a free log retrieval operation binding the contract event 0x14a183efe3b4ac142b43a017dc8b0436c6ca117f708a267c5fd178ab52e7d47a.
//
// Solidity: e DropratesSet(boosterTypes uint256, rarity uint256, cardList uint256[])
func (_CardFaucet *CardFaucetFilterer) FilterDropratesSet(opts *bind.FilterOpts) (*CardFaucetDropratesSetIterator, error) {

	logs, sub, err := _CardFaucet.contract.FilterLogs(opts, "DropratesSet")
	if err != nil {
		return nil, err
	}
	return &CardFaucetDropratesSetIterator{contract: _CardFaucet.contract, event: "DropratesSet", logs: logs, sub: sub}, nil
}

// WatchDropratesSet is a free log subscription operation binding the contract event 0x14a183efe3b4ac142b43a017dc8b0436c6ca117f708a267c5fd178ab52e7d47a.
//
// Solidity: e DropratesSet(boosterTypes uint256, rarity uint256, cardList uint256[])
func (_CardFaucet *CardFaucetFilterer) WatchDropratesSet(opts *bind.WatchOpts, sink chan<- *CardFaucetDropratesSet) (event.Subscription, error) {

	logs, sub, err := _CardFaucet.contract.WatchLogs(opts, "DropratesSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CardFaucetDropratesSet)
				if err := _CardFaucet.contract.UnpackLog(event, "DropratesSet", log); err != nil {
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

// CardFaucetGeneratedCardIterator is returned from FilterGeneratedCard and is used to iterate over the raw logs and unpacked data for GeneratedCard events raised by the CardFaucet contract.
type CardFaucetGeneratedCardIterator struct {
	Event *CardFaucetGeneratedCard // Event containing the contract specifics and raw log

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
func (it *CardFaucetGeneratedCardIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CardFaucetGeneratedCard)
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
		it.Event = new(CardFaucetGeneratedCard)
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
func (it *CardFaucetGeneratedCardIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CardFaucetGeneratedCardIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CardFaucetGeneratedCard represents a GeneratedCard event raised by the CardFaucet contract.
type CardFaucetGeneratedCard struct {
	CardId      *big.Int
	BoosterType uint8
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterGeneratedCard is a free log retrieval operation binding the contract event 0x990c7cbd32f45065c898d3b6c2f238fea70229c28ddeb9ddce2f402eeb60a10e.
//
// Solidity: e GeneratedCard(cardId uint256, boosterType uint8)
func (_CardFaucet *CardFaucetFilterer) FilterGeneratedCard(opts *bind.FilterOpts) (*CardFaucetGeneratedCardIterator, error) {

	logs, sub, err := _CardFaucet.contract.FilterLogs(opts, "GeneratedCard")
	if err != nil {
		return nil, err
	}
	return &CardFaucetGeneratedCardIterator{contract: _CardFaucet.contract, event: "GeneratedCard", logs: logs, sub: sub}, nil
}

// WatchGeneratedCard is a free log subscription operation binding the contract event 0x990c7cbd32f45065c898d3b6c2f238fea70229c28ddeb9ddce2f402eeb60a10e.
//
// Solidity: e GeneratedCard(cardId uint256, boosterType uint8)
func (_CardFaucet *CardFaucetFilterer) WatchGeneratedCard(opts *bind.WatchOpts, sink chan<- *CardFaucetGeneratedCard) (event.Subscription, error) {

	logs, sub, err := _CardFaucet.contract.WatchLogs(opts, "GeneratedCard")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CardFaucetGeneratedCard)
				if err := _CardFaucet.contract.UnpackLog(event, "GeneratedCard", log); err != nil {
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

// CardFaucetLimitedDifficultyChangedIterator is returned from FilterLimitedDifficultyChanged and is used to iterate over the raw logs and unpacked data for LimitedDifficultyChanged events raised by the CardFaucet contract.
type CardFaucetLimitedDifficultyChangedIterator struct {
	Event *CardFaucetLimitedDifficultyChanged // Event containing the contract specifics and raw log

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
func (it *CardFaucetLimitedDifficultyChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CardFaucetLimitedDifficultyChanged)
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
		it.Event = new(CardFaucetLimitedDifficultyChanged)
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
func (it *CardFaucetLimitedDifficultyChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CardFaucetLimitedDifficultyChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CardFaucetLimitedDifficultyChanged represents a LimitedDifficultyChanged event raised by the CardFaucet contract.
type CardFaucetLimitedDifficultyChanged struct {
	Target *big.Int
	Total  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterLimitedDifficultyChanged is a free log retrieval operation binding the contract event 0xbd9657322623c7e27f36dd8e37a8f3f3af601a8220da154bc738501de723a345.
//
// Solidity: e LimitedDifficultyChanged(target uint256, total uint256)
func (_CardFaucet *CardFaucetFilterer) FilterLimitedDifficultyChanged(opts *bind.FilterOpts) (*CardFaucetLimitedDifficultyChangedIterator, error) {

	logs, sub, err := _CardFaucet.contract.FilterLogs(opts, "LimitedDifficultyChanged")
	if err != nil {
		return nil, err
	}
	return &CardFaucetLimitedDifficultyChangedIterator{contract: _CardFaucet.contract, event: "LimitedDifficultyChanged", logs: logs, sub: sub}, nil
}

// WatchLimitedDifficultyChanged is a free log subscription operation binding the contract event 0xbd9657322623c7e27f36dd8e37a8f3f3af601a8220da154bc738501de723a345.
//
// Solidity: e LimitedDifficultyChanged(target uint256, total uint256)
func (_CardFaucet *CardFaucetFilterer) WatchLimitedDifficultyChanged(opts *bind.WatchOpts, sink chan<- *CardFaucetLimitedDifficultyChanged) (event.Subscription, error) {

	logs, sub, err := _CardFaucet.contract.WatchLogs(opts, "LimitedDifficultyChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CardFaucetLimitedDifficultyChanged)
				if err := _CardFaucet.contract.UnpackLog(event, "LimitedDifficultyChanged", log); err != nil {
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

// CardFaucetUpgradedCardToBEIterator is returned from FilterUpgradedCardToBE and is used to iterate over the raw logs and unpacked data for UpgradedCardToBE events raised by the CardFaucet contract.
type CardFaucetUpgradedCardToBEIterator struct {
	Event *CardFaucetUpgradedCardToBE // Event containing the contract specifics and raw log

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
func (it *CardFaucetUpgradedCardToBEIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CardFaucetUpgradedCardToBE)
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
		it.Event = new(CardFaucetUpgradedCardToBE)
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
func (it *CardFaucetUpgradedCardToBEIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CardFaucetUpgradedCardToBEIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CardFaucetUpgradedCardToBE represents a UpgradedCardToBE event raised by the CardFaucet contract.
type CardFaucetUpgradedCardToBE struct {
	CardId *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterUpgradedCardToBE is a free log retrieval operation binding the contract event 0x7950d9168421908781ca1e29d338b03d59d63e4afeebae90f60b4adfa14cda5b.
//
// Solidity: e UpgradedCardToBE(cardId uint256)
func (_CardFaucet *CardFaucetFilterer) FilterUpgradedCardToBE(opts *bind.FilterOpts) (*CardFaucetUpgradedCardToBEIterator, error) {

	logs, sub, err := _CardFaucet.contract.FilterLogs(opts, "UpgradedCardToBE")
	if err != nil {
		return nil, err
	}
	return &CardFaucetUpgradedCardToBEIterator{contract: _CardFaucet.contract, event: "UpgradedCardToBE", logs: logs, sub: sub}, nil
}

// WatchUpgradedCardToBE is a free log subscription operation binding the contract event 0x7950d9168421908781ca1e29d338b03d59d63e4afeebae90f60b4adfa14cda5b.
//
// Solidity: e UpgradedCardToBE(cardId uint256)
func (_CardFaucet *CardFaucetFilterer) WatchUpgradedCardToBE(opts *bind.WatchOpts, sink chan<- *CardFaucetUpgradedCardToBE) (event.Subscription, error) {

	logs, sub, err := _CardFaucet.contract.WatchLogs(opts, "UpgradedCardToBE")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CardFaucetUpgradedCardToBE)
				if err := _CardFaucet.contract.UnpackLog(event, "UpgradedCardToBE", log); err != nil {
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

// CardFaucetUpgradedCardToLEIterator is returned from FilterUpgradedCardToLE and is used to iterate over the raw logs and unpacked data for UpgradedCardToLE events raised by the CardFaucet contract.
type CardFaucetUpgradedCardToLEIterator struct {
	Event *CardFaucetUpgradedCardToLE // Event containing the contract specifics and raw log

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
func (it *CardFaucetUpgradedCardToLEIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CardFaucetUpgradedCardToLE)
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
		it.Event = new(CardFaucetUpgradedCardToLE)
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
func (it *CardFaucetUpgradedCardToLEIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CardFaucetUpgradedCardToLEIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CardFaucetUpgradedCardToLE represents a UpgradedCardToLE event raised by the CardFaucet contract.
type CardFaucetUpgradedCardToLE struct {
	CardId *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterUpgradedCardToLE is a free log retrieval operation binding the contract event 0xed0686f5b22006b8759c636c5c54007c681f80440f137b685adacb2907d43846.
//
// Solidity: e UpgradedCardToLE(cardId uint256)
func (_CardFaucet *CardFaucetFilterer) FilterUpgradedCardToLE(opts *bind.FilterOpts) (*CardFaucetUpgradedCardToLEIterator, error) {

	logs, sub, err := _CardFaucet.contract.FilterLogs(opts, "UpgradedCardToLE")
	if err != nil {
		return nil, err
	}
	return &CardFaucetUpgradedCardToLEIterator{contract: _CardFaucet.contract, event: "UpgradedCardToLE", logs: logs, sub: sub}, nil
}

// WatchUpgradedCardToLE is a free log subscription operation binding the contract event 0xed0686f5b22006b8759c636c5c54007c681f80440f137b685adacb2907d43846.
//
// Solidity: e UpgradedCardToLE(cardId uint256)
func (_CardFaucet *CardFaucetFilterer) WatchUpgradedCardToLE(opts *bind.WatchOpts, sink chan<- *CardFaucetUpgradedCardToLE) (event.Subscription, error) {

	logs, sub, err := _CardFaucet.contract.WatchLogs(opts, "UpgradedCardToLE")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CardFaucetUpgradedCardToLE)
				if err := _CardFaucet.contract.UnpackLog(event, "UpgradedCardToLE", log); err != nil {
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
