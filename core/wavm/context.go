// Copyright 2019 The go-vnt Authors
// This file is part of the go-vnt library.
//
// The go-vnt library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-vnt library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-vnt library. If not, see <http://www.gnu.org/licenses/>.

package wavm

import (
	"github.com/vntchain/go-vnt/accounts/abi"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core/state"
	"github.com/vntchain/go-vnt/core/vm/interface"
	"github.com/vntchain/go-vnt/core/wavm/contract"
	"github.com/vntchain/go-vnt/core/wavm/gas"
	"github.com/vntchain/go-vnt/core/wavm/storage"
	"github.com/vntchain/go-vnt/params"
	"math/big"
)

type ChainContext struct {
	// CanTransfer returns whether the account contains
	// sufficient ether to transfer the value
	CanTransfer func(inter.StateDB, common.Address, *big.Int) bool
	// Transfer transfers ether from one account to the other
	Transfer func(inter.StateDB, common.Address, common.Address, *big.Int)
	// GetHash returns the hash corresponding to n
	GetHash func(uint64) common.Hash
	// Message information
	Origin   common.Address // Provides information for ORIGIN
	GasPrice *big.Int       // Provides information for GASPRICE

	// Block information
	Coinbase       common.Address // Provides information for COINBASE
	GasLimit       uint64         // Provides information for GASLIMIT
	BlockNumber    *big.Int       // Provides information for NUMBER
	Time           *big.Int       // Provides information for TIME
	Difficulty     *big.Int       // Provides information for DIFFICULTY
	StateDB        *state.StateDB
	Contract       *contract.WASMContract
	Code           []byte  //Wasm contract code
	Abi            abi.ABI //Wasm contract abi
	Wavm           *WAVM
	IsCreated      bool
	StorageMapping map[uint64]storage.StorageMapping
	GasRule        gas.Gas
	GasCounter     gas.GasCounter
	GasTable       params.GasTable
}
