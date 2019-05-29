// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"math/big"
	"strings"

	hubble "github.com/vntchain/go-vnt"
	"github.com/vntchain/go-vnt/accounts/abi"
	"github.com/vntchain/go-vnt/accounts/abi/bind"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core/types"
	"github.com/vntchain/go-vnt/event"
)

// ChequebookABI is the input ABI used to generate the binding from.
const ChequebookABI = "[{\"constant\":false,\"inputs\":[],\"name\":\"kill\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"sent\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"beneficiary\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"sig_v\",\"type\":\"uint8\"},{\"name\":\"sig_r\",\"type\":\"bytes32\"},{\"name\":\"sig_s\",\"type\":\"bytes32\"}],\"name\":\"cash\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"deadbeat\",\"type\":\"address\"}],\"name\":\"Overdraft\",\"type\":\"event\"}]"

// ChequebookBin is the compiled bytecode used for deploying new contracts.
const ChequebookBin = `0x608060405260008054600160a060020a031916331790556102d0806100256000396000f3006080604052600436106100565763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166341c0e1b581146100585780637bf786f81461006d578063fbf788d6146100a0575b005b34801561006457600080fd5b506100566100d0565b34801561007957600080fd5b5061008e600160a060020a03600435166100f3565b60408051918252519081900360200190f35b3480156100ac57600080fd5b50610056600160a060020a036004351660243560ff60443516606435608435610105565b600054600160a060020a03163314156100f157600054600160a060020a0316ff5b565b60016020526000908152604090205481565b600160a060020a0385166000908152600160205260408120548190861161012b57600080fd5b604080516c010000000000000000000000003081028252600160a060020a038a160260148201526028810188905281519081900360480181206000808352602083810180865283905260ff8a16848601526060840189905260808401889052935191955060019360a0808501949193601f198101939281900390910191865af11580156101bc573d6000803e3d6000fd5b5050604051601f190151600054600160a060020a0390811691161490506101e257600080fd5b50600160a060020a03861660009081526001602052604090205485033031811161025057600160a060020a0387166000818152600160205260408082208990555183156108fc0291849190818181858888f1935050505015801561024a573d6000803e3d6000fd5b5061029b565b60005460408051600160a060020a039092168252517f2250e2993c15843b32621c89447cc589ee7a9f049c026986e545d3c2c0c6f9789181900360200190a186600160a060020a0316ff5b505050505050505600a165627a7a72305820b2c4748d2140dcc903f965abe3187ac942f8c04dacb9887c47bdb86a6c84dfc70029`

// DeployChequebook deploys a new VNT contract, binding an instance of Chequebook to it.
func DeployChequebook(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Chequebook, error) {
	parsed, err := abi.JSON(strings.NewReader(ChequebookABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ChequebookBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Chequebook{ChequebookCaller: ChequebookCaller{contract: contract}, ChequebookTransactor: ChequebookTransactor{contract: contract}, ChequebookFilterer: ChequebookFilterer{contract: contract}}, nil
}

// Chequebook is an auto generated Go binding around an VNT contract.
type Chequebook struct {
	ChequebookCaller     // Read-only binding to the contract
	ChequebookTransactor // Write-only binding to the contract
	ChequebookFilterer   // Log filterer for contract events
}

// ChequebookCaller is an auto generated read-only Go binding around an VNT contract.
type ChequebookCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChequebookTransactor is an auto generated write-only Go binding around an VNT contract.
type ChequebookTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChequebookFilterer is an auto generated log filtering Go binding around an VNT contract events.
type ChequebookFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChequebookSession is an auto generated Go binding around an VNT contract,
// with pre-set call and transact options.
type ChequebookSession struct {
	Contract     *Chequebook       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ChequebookCallerSession is an auto generated read-only Go binding around an VNT contract,
// with pre-set call options.
type ChequebookCallerSession struct {
	Contract *ChequebookCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// ChequebookTransactorSession is an auto generated write-only Go binding around an VNT contract,
// with pre-set transact options.
type ChequebookTransactorSession struct {
	Contract     *ChequebookTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ChequebookRaw is an auto generated low-level Go binding around an VNT contract.
type ChequebookRaw struct {
	Contract *Chequebook // Generic contract binding to access the raw methods on
}

// ChequebookCallerRaw is an auto generated low-level read-only Go binding around an VNT contract.
type ChequebookCallerRaw struct {
	Contract *ChequebookCaller // Generic read-only contract binding to access the raw methods on
}

// ChequebookTransactorRaw is an auto generated low-level write-only Go binding around an VNT contract.
type ChequebookTransactorRaw struct {
	Contract *ChequebookTransactor // Generic write-only contract binding to access the raw methods on
}

// NewChequebook creates a new instance of Chequebook, bound to a specific deployed contract.
func NewChequebook(address common.Address, backend bind.ContractBackend) (*Chequebook, error) {
	contract, err := bindChequebook(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Chequebook{ChequebookCaller: ChequebookCaller{contract: contract}, ChequebookTransactor: ChequebookTransactor{contract: contract}, ChequebookFilterer: ChequebookFilterer{contract: contract}}, nil
}

// NewChequebookCaller creates a new read-only instance of Chequebook, bound to a specific deployed contract.
func NewChequebookCaller(address common.Address, caller bind.ContractCaller) (*ChequebookCaller, error) {
	contract, err := bindChequebook(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChequebookCaller{contract: contract}, nil
}

// NewChequebookTransactor creates a new write-only instance of Chequebook, bound to a specific deployed contract.
func NewChequebookTransactor(address common.Address, transactor bind.ContractTransactor) (*ChequebookTransactor, error) {
	contract, err := bindChequebook(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChequebookTransactor{contract: contract}, nil
}

// NewChequebookFilterer creates a new log filterer instance of Chequebook, bound to a specific deployed contract.
func NewChequebookFilterer(address common.Address, filterer bind.ContractFilterer) (*ChequebookFilterer, error) {
	contract, err := bindChequebook(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChequebookFilterer{contract: contract}, nil
}

// bindChequebook binds a generic wrapper to an already deployed contract.
func bindChequebook(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ChequebookABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Chequebook *ChequebookRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Chequebook.Contract.ChequebookCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Chequebook *ChequebookRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Chequebook.Contract.ChequebookTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Chequebook *ChequebookRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Chequebook.Contract.ChequebookTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Chequebook *ChequebookCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Chequebook.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Chequebook *ChequebookTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Chequebook.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Chequebook *ChequebookTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Chequebook.Contract.contract.Transact(opts, method, params...)
}

// Sent is a free data retrieval call binding the contract method 0x7bf786f8.
//
// Solidity: function sent( address) constant returns(uint256)
func (_Chequebook *ChequebookCaller) Sent(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Chequebook.contract.Call(opts, out, "sent", arg0)
	return *ret0, err
}

// Sent is a free data retrieval call binding the contract method 0x7bf786f8.
//
// Solidity: function sent( address) constant returns(uint256)
func (_Chequebook *ChequebookSession) Sent(arg0 common.Address) (*big.Int, error) {
	return _Chequebook.Contract.Sent(&_Chequebook.CallOpts, arg0)
}

// Sent is a free data retrieval call binding the contract method 0x7bf786f8.
//
// Solidity: function sent( address) constant returns(uint256)
func (_Chequebook *ChequebookCallerSession) Sent(arg0 common.Address) (*big.Int, error) {
	return _Chequebook.Contract.Sent(&_Chequebook.CallOpts, arg0)
}

// Cash is a paid mutator transaction binding the contract method 0xfbf788d6.
//
// Solidity: function cash(beneficiary address, amount uint256, sig_v uint8, sig_r bytes32, sig_s bytes32) returns()
func (_Chequebook *ChequebookTransactor) Cash(opts *bind.TransactOpts, beneficiary common.Address, amount *big.Int, sig_v uint8, sig_r [32]byte, sig_s [32]byte) (*types.Transaction, error) {
	return _Chequebook.contract.Transact(opts, "cash", beneficiary, amount, sig_v, sig_r, sig_s)
}

// Cash is a paid mutator transaction binding the contract method 0xfbf788d6.
//
// Solidity: function cash(beneficiary address, amount uint256, sig_v uint8, sig_r bytes32, sig_s bytes32) returns()
func (_Chequebook *ChequebookSession) Cash(beneficiary common.Address, amount *big.Int, sig_v uint8, sig_r [32]byte, sig_s [32]byte) (*types.Transaction, error) {
	return _Chequebook.Contract.Cash(&_Chequebook.TransactOpts, beneficiary, amount, sig_v, sig_r, sig_s)
}

// Cash is a paid mutator transaction binding the contract method 0xfbf788d6.
//
// Solidity: function cash(beneficiary address, amount uint256, sig_v uint8, sig_r bytes32, sig_s bytes32) returns()
func (_Chequebook *ChequebookTransactorSession) Cash(beneficiary common.Address, amount *big.Int, sig_v uint8, sig_r [32]byte, sig_s [32]byte) (*types.Transaction, error) {
	return _Chequebook.Contract.Cash(&_Chequebook.TransactOpts, beneficiary, amount, sig_v, sig_r, sig_s)
}

// Kill is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: function kill() returns()
func (_Chequebook *ChequebookTransactor) Kill(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Chequebook.contract.Transact(opts, "kill")
}

// Kill is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: function kill() returns()
func (_Chequebook *ChequebookSession) Kill() (*types.Transaction, error) {
	return _Chequebook.Contract.Kill(&_Chequebook.TransactOpts)
}

// Kill is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: function kill() returns()
func (_Chequebook *ChequebookTransactorSession) Kill() (*types.Transaction, error) {
	return _Chequebook.Contract.Kill(&_Chequebook.TransactOpts)
}

// ChequebookOverdraftIterator is returned from FilterOverdraft and is used to iterate over the raw logs and unpacked data for Overdraft events raised by the Chequebook contract.
type ChequebookOverdraftIterator struct {
	Event *ChequebookOverdraft // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  hubble.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ChequebookOverdraftIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChequebookOverdraft)
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
		it.Event = new(ChequebookOverdraft)
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
func (it *ChequebookOverdraftIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChequebookOverdraftIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChequebookOverdraft represents a Overdraft event raised by the Chequebook contract.
type ChequebookOverdraft struct {
	Deadbeat common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOverdraft is a free log retrieval operation binding the contract event 0x2250e2993c15843b32621c89447cc589ee7a9f049c026986e545d3c2c0c6f978.
//
// Solidity: e Overdraft(deadbeat address)
func (_Chequebook *ChequebookFilterer) FilterOverdraft(opts *bind.FilterOpts) (*ChequebookOverdraftIterator, error) {

	logs, sub, err := _Chequebook.contract.FilterLogs(opts, "Overdraft")
	if err != nil {
		return nil, err
	}
	return &ChequebookOverdraftIterator{contract: _Chequebook.contract, event: "Overdraft", logs: logs, sub: sub}, nil
}

// WatchOverdraft is a free log subscription operation binding the contract event 0x2250e2993c15843b32621c89447cc589ee7a9f049c026986e545d3c2c0c6f978.
//
// Solidity: e Overdraft(deadbeat address)
func (_Chequebook *ChequebookFilterer) WatchOverdraft(opts *bind.WatchOpts, sink chan<- *ChequebookOverdraft) (event.Subscription, error) {

	logs, sub, err := _Chequebook.contract.WatchLogs(opts, "Overdraft")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChequebookOverdraft)
				if err := _Chequebook.contract.UnpackLog(event, "Overdraft", log); err != nil {
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
