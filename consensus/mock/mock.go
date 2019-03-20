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

package mock

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/consensus"
	"github.com/vntchain/go-vnt/core/state"
	"github.com/vntchain/go-vnt/core/types"
	"github.com/vntchain/go-vnt/log"
	"github.com/vntchain/go-vnt/rpc"
)

var (
	// errUnknownBlock is returned when the list of signers is requested for a block
	// that is not part of the local blockchain.
	errUnknownBlock = errors.New("unknown block")

	// errInvalidDifficulty is returned if the difficulty of a block is not 1
	errInvalidDifficulty = errors.New("invalid difficulty")
)

type Mock struct {
	fakeFail uint64
}

func NewMock() *Mock {
	return &Mock{}
}

func NewMockFail(failNum uint64) *Mock {
	return &Mock{
		fakeFail: failNum,
	}
}
func (m *Mock) Author(header *types.Header) (common.Address, error) {
	return header.Coinbase, nil
}

func (m *Mock) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	err := m.verifyHeader(chain, header, nil)
	if err != nil {
		log.Debug("VerifyHeader error", "hash", header.Hash(), "err", err.Error())
	} else {
		log.Debug("VerifyHeader NO error", "hash", header.Hash(), "number", header.Number.Int64())
	}
	return err
}

func (m *Mock) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{})
	results := make(chan error, len(headers))
	go func() {
		for i, header := range headers {
			err := m.verifyHeader(chain, header, headers[:i])

			select {
			case <-abort:
				return
			case results <- err:
			}
		}
	}()
	return abort, results
}

func (m *Mock) verifyHeader(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	if m.fakeFail == header.Number.Uint64() {
		return fmt.Errorf("fake Fail")
	}
	if header.Number == nil {
		return errUnknownBlock
	}
	number := header.Number.Uint64()

	// Ensure that the block's difficulty is meaningful (may not be correct at this point)
	if number > 0 {
		if header.Difficulty == nil || header.Difficulty.Cmp(big.NewInt(1)) != 0 {
			return errInvalidDifficulty
		}
	}
	// All basic checks passed, verify cascading fields
	return m.verifyCascadingFields(chain, header, parents)
}
func (m *Mock) verifyCascadingFields(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	// The genesis block is the always valid dead-end
	number := header.Number.Uint64()
	if number == 0 {
		return nil
	}
	// Ensure that the block's timestamp isn't too close to it's parent
	var parent *types.Header
	if len(parents) > 0 {
		parent = parents[len(parents)-1]
	} else {
		parent = chain.GetHeader(header.ParentHash, number-1)
	}
	if parent == nil || parent.Number.Uint64() != number-1 || parent.Hash() != header.ParentHash {
		return consensus.ErrUnknownAncestor
	}
	return nil
}

func (m *Mock) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	return nil
}

func (m *Mock) VerifyWitnesses(header *types.Header, db *state.StateDB, parent *types.Header) error {
	return nil
}

func (m *Mock) VerifyCommitMsg(block *types.Block) error {
	return nil
}

func (m *Mock) Prepare(chain consensus.ChainReader, header *types.Header) error {
	header.Difficulty = big.NewInt(1)
	return nil
}

func (m *Mock) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	receipts []*types.Receipt) (*types.Block, error) {
	state.AddBalance(header.Coinbase, big.NewInt(5e18))
	header.Root = state.IntermediateRoot(true)

	// Assemble and return the final block for sealing
	return types.NewBlock(header, txs, receipts), nil
}

func (m *Mock) Seal(chain consensus.ChainReader, block *types.Block, stop <-chan struct{}) (*types.Block, error) {
	return block, nil
}

func (m *Mock) CalcDifficulty(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {
	return common.Big1
}

func (m *Mock) HandleBftMsg(chain consensus.ChainReader, msg types.ConsensusMsg) {
}

// APIs returns the RPC APIs this consensus engine provides.
func (m *Mock) APIs(chain consensus.ChainReader) []rpc.API {
	return nil
}
