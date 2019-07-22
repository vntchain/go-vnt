// Copyright 2014 The go-ethereum Authors
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

package vm

import (
	"math/big"

	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core/vm/election"
	inter "github.com/vntchain/go-vnt/core/vm/interface"
)

// PrecompiledContract is the basic interface for native Go contracts. The implementation
// requires a deterministic gas count based on the input size of the Run method of the
// contract.
type PrecompiledContract interface {
	RequiredGas(input []byte) uint64                                              // RequiredPrice calculates the contract gas use
	Run(context inter.ChainContext, input []byte, value *big.Int) ([]byte, error) // Run runs the precompiled contract
}

// PrecompiledContractsHubble contains the default set of pre-compiled VNT
// contracts used in the Hubble release.
var PrecompiledContractsHubble = map[common.Address]PrecompiledContract{
	common.BytesToAddress([]byte{9}): &election.Election{},
}

// RunPrecompiledContract runs and evaluates the output of a precompiled contract.
func RunPrecompiledContract(context inter.ChainContext, p PrecompiledContract, input []byte, contract inter.Contract) (ret []byte, err error) {
	gas := p.RequiredGas(input)
	if contract.UseGas(gas) {
		return p.Run(context, input, contract.Value())
	}
	return nil, ErrOutOfGas
}
