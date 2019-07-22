// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package contract

import (
	"math/big"

	"github.com/vntchain/go-vnt/core/vm/interface"

	"github.com/vntchain/go-vnt/common"
)

type WasmCode struct {
	Code     []byte
	Abi      []byte
	Compiled []byte
}

type WASMContractRef interface {
	Address() common.Address
}

// Contract represents an VNTChain contract in the state database. It contains
// the the contract code, calling arguments. Contract implements ContractRef
type WASMContract struct {
	// CallerAddress is the result of the caller which initialised this
	// contract. However when the "call method" is delegated this value
	// needs to be initialised to that of the caller's caller.
	CallerAddress common.Address
	caller        WASMContractRef
	self          WASMContractRef

	Code     []byte
	CodeHash common.Hash
	CodeAddr *common.Address
	Input    []byte

	GasLimit uint64
	Gas      uint64
	value    *big.Int

	Args []byte

	DelegateCall bool
}

// NewWASMContract returns a new contract environment for the execution of WAVM.
func NewWASMContract(caller WASMContractRef, object WASMContractRef, value *big.Int, gas uint64) *WASMContract {
	c := &WASMContract{CallerAddress: caller.Address(), caller: caller, self: object, Args: nil}

	// Gas should be a pointer so it can safely be reduced through the run
	// This pointer will be off the state transition
	c.Gas = gas
	c.GasLimit = gas
	// ensures a value is set
	c.value = value

	return c
}

// AsDelegate sets the contract to be a delegate call and returns the current
// contract (for chaining calls)
func (c *WASMContract) AsDelegate() inter.Contract {
	c.DelegateCall = true
	// NOTE: caller must, at all times be a contract. It should never happen
	// that caller is something other than a Contract.
	parent := c.caller.(*WASMContract)
	c.CallerAddress = parent.CallerAddress
	c.value = parent.value

	return c
}

// GetByte returns the n'th byte in the contract's byte array
func (c *WASMContract) GetByte(n uint64) byte {
	if n < uint64(len(c.Code)) {
		return c.Code[n]
	}

	return 0
}

// Caller returns the caller of the contract.
//
// Caller will recursively call caller when the contract is a delegate
// call, including that of caller's caller.
func (c *WASMContract) Caller() common.Address {
	return c.CallerAddress
}

// UseGas attempts the use gas and subtracts it and returns true on success
func (c *WASMContract) UseGas(gas uint64) (ok bool) {
	if c.Gas < gas {
		return false
	}
	c.Gas -= gas
	return true
}

// Address returns the contracts address
func (c *WASMContract) Address() common.Address {
	return c.self.Address()
}

// Value returns the contracts value (sent to it from it's caller)
func (c *WASMContract) Value() *big.Int {
	return c.value
}

// SetCode sets the code to the contract
func (c *WASMContract) SetCode(hash common.Hash, code []byte) {
	c.Code = code
	c.CodeHash = hash
}

// SetCallCode sets the code of the contract and address of the backing data
// object
func (c *WASMContract) SetCallCode(addr *common.Address, hash common.Hash, code []byte) {
	c.Code = code
	c.CodeHash = hash
	c.CodeAddr = addr
}
