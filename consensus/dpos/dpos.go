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

package dpos

import (
	"errors"
	"math/big"
	"sync"
	"time"

	"bytes"
	"encoding/binary"
	"fmt"

	lru "github.com/hashicorp/golang-lru"
	"github.com/vntchain/go-vnt/accounts"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/common/math"
	"github.com/vntchain/go-vnt/consensus"
	"github.com/vntchain/go-vnt/core"
	"github.com/vntchain/go-vnt/core/state"
	"github.com/vntchain/go-vnt/core/types"
	"github.com/vntchain/go-vnt/core/vm/election"
	"github.com/vntchain/go-vnt/crypto"
	"github.com/vntchain/go-vnt/crypto/sha3"
	"github.com/vntchain/go-vnt/log"
	"github.com/vntchain/go-vnt/params"
	"github.com/vntchain/go-vnt/rlp"
	"github.com/vntchain/go-vnt/rpc"
	"github.com/vntchain/go-vnt/vntdb"
)

const (
	inMemorySignatures = 4096 // Number of recent block signatures to keep in memory
	updateTimeLen      = 8    // Number of bytes the witnesses list update time take up
)

var (
	VortexBlockReward     *big.Int = big.NewInt(6e+18)
	VortexCandidatesBonus *big.Int = big.NewInt(6e+18)
	// 2 seconds one block, 3 years producing about 47304000 blocks
	stageTwoBlkNr = big.NewInt(47304000)
	// 2 seconds one block, 6 years producing about 94608000 blocks
	stageThreeBlkNr = big.NewInt(94608000)
)

// Various error messages to mark blocks invalid. These should be private to
// prevent engine specific errors from being referenced in the remainder of the
// codebase, inherently breaking if the engine is swapped out. Please put common
// error types into the consensus package.
var (
	// errUnknownBlock is returned when the list of signers is requested for a block
	// that is not part of the local blockchain.
	errUnknownBlock = errors.New("unknown block")

	// block has a beneficiary set to non-zeroes.
	errInvalidCoinBase = errors.New("coinbase in block non-zero")

	// errInvalidDifficulty is returned if the difficulty of a block is not 1
	errInvalidDifficulty = errors.New("invalid difficulty")

	// errInvalidTimestamp is returned if the timestamp of a block is lower than
	// the previous block's timestamp + the minimum block period.
	errInvalidTimestamp = errors.New("invalid timestamp")

	// witness should be same with the parent
	errWitnesses = errors.New("witnesses is different from parent")

	// witness should in turn
	errOutTurn = errors.New("witness is out turn")

	// errNoPreviousWitness is returned if no previous witness still in current witness list
	errNoPreviousWitness = errors.New("no previous witness still in current witness list")

	// errInvalidExtraLen is returned if extra length is invalid
	errInvalidExtraLen = errors.New("invalid Extra length")
)

type SignerFn func(accounts.Account, []byte) ([]byte, error)

// getHeaderFromParentsFn get header from previous headers
type getHeaderFromParentsFn func(hash common.Hash, num uint64) *types.Header

type Dpos struct {
	config         *params.DposConfig
	bft            *BftManager
	db             vntdb.Database // Database to store and retrieve dpos temp data, current not used
	signatures     *lru.ARCCache  // Signatures of recent blocks to speed up block producing
	signer         common.Address // VNT address of the signing key
	signFn         SignerFn       // Signer function to authorize hashes with
	lock           sync.RWMutex   // Protects the signer fields
	updateInterval *big.Int       // Duration of update witnesses list
	lastBounty     lastBountyInfo // 上次发放激励的信息

	sendBftPeerUpdateFn func(urls []string)
}

type lastBountyInfo struct {
	bountyHeight *big.Int // 上次发送激励的高度
	updateHeight *big.Int // 更新当前数据的高度
	sync.RWMutex          // 存在并发访问，加锁保护
}

// sigHash returns the hash which is used as input for the proof-of-authority
// signing. It is the hash of the entire header apart from the 65 byte signature
// contained at the end of the extra data.
//
// Note, the method requires the extra data to be at least 65 bytes, otherwise it
// panics. This is done to avoid accidentally using both forms (signature present
// or not), which could be abused to produce different hashes for the same header.
func sigHash(header *types.Header) (hash common.Hash, err error) {
	hasher := sha3.NewKeccak256()

	err = rlp.Encode(hasher, []interface{}{
		header.ParentHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra,
		header.Witnesses,
	})
	if err != nil {
		return common.Hash{}, err
	}

	hasher.Sum(hash[:0])
	return hash, nil
}

// ecrecover extracts the VNT account address from a signed header.
func ecrecover(header *types.Header, sigcache *lru.ARCCache) (common.Address, error) {
	// If the signature's already cached, return that
	hash := header.Hash()
	if address, known := sigcache.Get(hash); known {
		return address.(common.Address), nil
	}

	signature := header.Signature

	// Recover the public key and the VNT address
	sh, err := sigHash(header)
	if err != nil {
		return common.Address{}, err
	}
	pubkey, err := crypto.Ecrecover(sh.Bytes(), signature)
	if err != nil {
		return common.Address{}, fmt.Errorf("ecrecover fialed: %s", err.Error())
	}
	var signer common.Address
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])

	sigcache.Add(hash, signer)
	return signer, nil
}

// New creates a Delegated proof-of-stake consensus engine with the initial
// signers set to the ones provided by the user.
func New(config *params.DposConfig, db vntdb.Database) *Dpos {
	signatures, _ := lru.NewARC(inMemorySignatures)

	d := &Dpos{
		config:         config,
		bft:            nil,
		db:             db,
		signatures:     signatures,
		updateInterval: nil,

		lastBounty: lastBountyInfo{
			bountyHeight: big.NewInt(0),
			updateHeight: big.NewInt(0),
		},
	}

	d.bft = newBftManager(d)
	d.setUpdateInterval()

	return d
}

func (d *Dpos) InitBft(sendBftMsg func(types.ConsensusMsg), SendPeerUpdate func(urls []string), verifyBlock func(*types.Block) (types.Receipts, []*types.Log, uint64, error), writeBlock func(*types.Block) error) {
	d.sendBftPeerUpdateFn = SendPeerUpdate

	// Init bft function
	d.bft.sendBftMsg = sendBftMsg
	d.bft.verifyBlock = verifyBlock
	d.bft.writeBlock = writeBlock

	// Init bft field
	d.bft.coinBase = d.coinBase()

	d.bft.producingStart()
}

// Author implements consensus.Engine, returning the VNT address recovered
// from the signature in the header's extra-data section.
func (d *Dpos) Author(header *types.Header) (common.Address, error) {
	return ecrecover(header, d.signatures)
}

// VerifyHeader checks whether a header conforms to the consensus rules.
func (d *Dpos) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {

	err := d.verifyHeader(chain, header, nil)
	if err != nil {
		log.Debug("VerifyHeader error", "hash", header.Hash().String(), "err", err.Error())
	} else {
		log.Debug("VerifyHeader NO error", "hash", header.Hash().String(), "number", header.Number.Int64())
	}
	return err
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers. The
// method returns a quit channel to abort the operations and a results channel to
// retrieve the async verifications (the order is that of the input slice).
func (d *Dpos) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{})
	results := make(chan error, len(headers))
	go func() {
		for i, header := range headers {
			err := d.verifyHeader(chain, header, headers[:i])

			select {
			case <-abort:
				return
			case results <- err:
			}
		}
	}()
	return abort, results
}

// verifyHeader checks whether a header conforms to the consensus rules.The
// caller may optionally pass in a batch of parents (ascending order) to avoid
// looking those up from the database. This is useful for concurrently verifying
// a batch of new headers.
func (d *Dpos) verifyHeader(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	if header.Number == nil {
		return errUnknownBlock
	}
	number := header.Number.Uint64()

	// todo verify header.time after implement bft without time trigger produce block
	// Don't waste time checking blocks from the future
	// if header.Time.Cmp(big.NewInt(time.Now().Unix())) > 0 {
	// 	return consensus.ErrFutureBlock
	// }

	// Ensure extra has correct length' value checked in verify witnesses
	if len(header.Extra) != updateTimeLen {
		return errInvalidExtraLen
	}

	// Ensure that the block's difficulty is meaningful (may not be correct at this point)
	if number > 0 {
		if header.Difficulty == nil || header.Difficulty.Cmp(big.NewInt(1)) != 0 {
			return errInvalidDifficulty
		}
	}
	// All basic checks passed, verify cascading fields
	return d.verifyCascadingFields(chain, header, parents)
}

// verifyCascadingFields verifies all the header fields that are not standalone,
// rather depend on a batch of previous headers. The caller may optionally pass
// in a batch of parents (ascending order) to avoid looking those up from the
// database. This is useful for concurrently verifying a batch of new headers.
func (d *Dpos) verifyCascadingFields(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
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

	headerTime, parentTime := header.Time.Uint64(), parent.Time.Uint64()
	if headerTime <= parentTime || (headerTime-parentTime)%d.config.Period != 0 {
		log.Warn("Timestamp is invalid", "headerTime", headerTime, "parentTime", parentTime)
		return errInvalidTimestamp
	}

	// Verify that the gas limit is <= 2^63-1
	cap := uint64(0x7fffffffffffffff)
	if header.GasLimit > cap {
		return fmt.Errorf("invalid gasLimit: have %v, max %v", header.GasLimit, cap)
	}
	// Verify that the gasUsed is <= gasLimit
	if header.GasUsed > header.GasLimit {
		return fmt.Errorf("invalid gasUsed: have %d, gasLimit %d", header.GasUsed, header.GasLimit)
	}
	// Verify that the gas limit remains within allowed bounds
	diff := int64(parent.GasLimit) - int64(header.GasLimit)
	if diff < 0 {
		diff *= -1
	}
	limit := parent.GasLimit / params.GasLimitBoundDivisor
	if uint64(diff) >= limit || header.GasLimit < params.MinGasLimit {
		return fmt.Errorf("invalid gas limit: have %d, want %d += %d", header.GasLimit, parent.GasLimit, limit)
	}

	if len(header.Witnesses) != d.config.WitnessesNum {
		return errWitnesses
	}

	// All basic checks passed, verify the seal and return
	return d.verifySeal(chain, header, parents)
}

// VerifySeal implements consensus.Engine, checking whether the signature contained
// in the header satisfies the consensus protocol requirements.
func (d *Dpos) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	return d.verifySeal(chain, header, nil)
}

// verifySeal checks whether the signature contained in the header satisfies the
// consensus protocol requirements. The method accepts an optional list of parent
// headers that aren't yet part of the local blockchain to generate the snapshots
// from.
func (d *Dpos) verifySeal(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	// Verifying the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return errUnknownBlock
	}
	// Resolve the authorization key and check against signers
	signer, err := ecrecover(header, d.signatures)
	if err != nil {
		return err
	}

	if signer != header.Coinbase {
		return errInvalidCoinBase
	}

	// 确认轮次对不对，是不是该这个节点出块
	if !d.inTurn(header, signer, chain, parents) {
		return errOutTurn
	}
	return nil
}

// VerifyWitnesses Verify witness list and update time(header.Extra) for DPoS
func (d *Dpos) VerifyWitnesses(header *types.Header, db *state.StateDB, parent *types.Header) error {
	updated, localWitnesses := d.getWitnesses(header, db, parent)
	if len(localWitnesses) != len(header.Witnesses) {
		return fmt.Errorf("witnesses length not match")
	}

	// Check header.Extra
	if needSetUpdateTime(updated, header.Number.Uint64()) {
		if !d.updatedWitnessCheckByTime(header) {
			return fmt.Errorf("header.Extra is mismatch with header.Time when update")
		}
	} else {
		if !bytes.Equal(header.Extra, parent.Extra) {
			return fmt.Errorf("header.Extra is mismatch with parent.Time when NOT update")
		}
	}

	// Check the witnesses list
	for i := 0; i < len(localWitnesses); i++ {
		if localWitnesses[i] != header.Witnesses[i] {
			return fmt.Errorf("witnesses is not match")
		}
	}
	return nil
}

// Prepare implements consensus.Engine, preparing all the consensus fields of the
// header for running the transactions on top.
func (d *Dpos) Prepare(chain consensus.ChainReader, header *types.Header) error {
	var (
		updated bool
		err     error
	)

	number := header.Number.Uint64()
	if number == 0 {
		return errUnknownBlock
	}

	d.lock.RLock()
	header.Coinbase = d.signer
	d.lock.RUnlock()

	// Set the correct difficulty
	header.Difficulty = big.NewInt(1)

	// Try to sleep if can not find parent header, try to stop commitNewWork start again immediately.
	// WARN: there must be some db write or read error
	parent := chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	if parent == nil {
		log.Error("Miss parent in dos.Prepare()", "hash", header.ParentHash.String(), "number", header.Number.Uint64()-1)
		time.Sleep(time.Minute * 2)
		return consensus.ErrUnknownAncestor
	}

	// Put next time in header
	produceTime, nPeriod, err := d.nextProduceTime(parent.Time)
	if err != nil {
		return err
	}
	header.Time = produceTime

	// Update witness list if needed，and set Extra with update value
	updated, header.Witnesses, err = d.getWitnessesForProduce(header, chain, parent)
	if err != nil {
		return err
	}

	// Start a new round of bft
	r := uint32(nPeriod.Uint64()) - 1
	d.bft.blockRound = r
	go d.bft.newRound(header.Number, r, header.Witnesses)

	// Make sure self is the current block producer before produce
	witness := header.Coinbase
	if !d.inTurn(header, witness, chain, nil) {
		log.Debug("Prepare failed", "err", errOutTurn)
		return fmt.Errorf("node is out of turn")
	}

	// Fill Extra with the update time
	// If this updated the witnesses list in this block, extra = this header time
	// else, extra = last update time(get from parent's block)
	header.Extra = make([]byte, updateTimeLen)
	if needSetUpdateTime(updated, number) {
		copy(header.Extra, encodeUpdateTime(header.Time))
	} else {
		copy(header.Extra, parent.Extra)
	}

	return nil
}

// needSetUpdateTime block 1 is special, should use it's header time instead of
// genesis's extra, because genesis's extra is nil
func needSetUpdateTime(update bool, number uint64) bool {
	return update || number == 1
}

// Finalize implements consensus.Engine,  grants reward and returns the final block.
func (d *Dpos) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, receipts []*types.Receipt) (*types.Block, error) {
	// Granting bounty, if any left
	if err := d.grantingReward(chain, header, state); err != nil {
		return nil, err
	}

	// Commit db
	header.Root = state.IntermediateRoot(true)

	// Assemble and return the final block for sealing
	return types.NewBlock(header, txs, receipts), nil
}

// grantingReward granting producing reward to the current block producer for producing this block,
// and granting vote reward to all the active witness candidates. the vote reward, which each witness
// earned, in direct proportion to it's vote percentage.
// WARN: There is no reward if no VNT bounty left.
func (d *Dpos) grantingReward(chain consensus.ChainReader, header *types.Header, state *state.StateDB) error {
	if restBounty := election.QueryRestVNTBounty(state); restBounty.Cmp(common.Big0) > 0 {
		var err error
		// Reward BP for producing this block
		reward := curHeightBonus(header.Number, VortexBlockReward)
		if restBounty.Cmp(reward) < 0 {
			reward = restBounty
		}
		if restBounty, err = election.GrantBounty(state, reward); err == nil {
			state.AddBalance(header.Coinbase, reward)
		}

		// Reward all witness candidates, when update witness list, if has any bounty
		if d.updatedWitnessCheckByTime(header) && restBounty.Cmp(common.Big0) > 0 {
			candis, allBonus, err := d.voteBonusPreWork(chain, header, state)
			if err != nil {
				return err
			}

			// the amount of bounty granted must not greater than the left bounty
			actualBonus := math.BigMin(allBonus, restBounty)
			log.Debug("Vote bounty", "each bounty(wei)", actualBonus.String())
			if bonus := d.calcVoteBounty(candis, actualBonus); bonus != nil {
				return election.AddCandidatesBounty(state, bonus, actualBonus)
			}
		}
	}
	return nil
}

// Authorize injects a private key into the consensus engine to mint new blocks
// with.
func (d *Dpos) Authorize(signer common.Address, signFn SignerFn) {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.signer = signer
	d.signFn = signFn
}

// Seal implements consensus.Engine, attempting to create a sealed block using
// the local signing credentials.
func (d *Dpos) Seal(chain consensus.ChainReader, block *types.Block, stop <-chan struct{}) (*types.Block, error) {
	header := block.Header()

	// Sealing the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return nil, errUnknownBlock
	}

	// Don't hold the witness fields for the entire sealing procedure
	d.lock.RLock()
	witness, signFn := d.signer, d.signFn
	d.lock.RUnlock()

	// Sign all the things without Signature
	sh, err := sigHash(header)
	if err != nil {
		return nil, err
	}
	sighash, err := signFn(accounts.Account{Address: witness}, sh.Bytes())
	if err != nil {
		return nil, err
	}
	header.Signature = make([]byte, len(sighash))
	copy(header.Signature[:], sighash)

	log.Info("Seal block", "at time", time.Now().Unix())
	d.bft.startPrePrepare(block.WithSeal(header))
	// DPoS no need return block
	return nil, nil
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
// that a new block should have based on the previous blocks in the chain and the
// current signer.
func (d *Dpos) CalcDifficulty(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {
	return common.Big1
}

// APIs implements consensus.Engine, returning the user facing RPC API to allow
// controlling the signer voting.
func (d *Dpos) APIs(chain consensus.ChainReader) []rpc.API {
	return []rpc.API{{
		Namespace: "dpos",
		Version:   "1.0",
		Service:   &API{chain: chain, dpos: d},
		Public:    false,
	}}
}

// nextProduceTime calculate next block produce time with previous block time
// 		cur_time = time()
// 		dur = cur_time - parent_time
// 		nPeriod = dur / interval
// 		// always up bound
// 		nPeriod++
// 		return parent_time + diff_index * interval
func (d *Dpos) nextProduceTime(preBlockTime *big.Int) (produceTime *big.Int, nPeriod *big.Int, err error) {
	now := time.Now().Unix()
	dur := new(big.Int).Sub(new(big.Int).SetInt64(now), preBlockTime)
	period := new(big.Int).SetUint64(d.config.Period)
	// the unit is second, even no left of DivMod, but current time is in new period
	nPeriod = new(big.Int).Div(dur, period)
	nPeriod.Add(nPeriod, common.Big1)

	nextTime := new(big.Int).Mul(nPeriod, period)
	nextTime.Add(nextTime, preBlockTime)

	return nextTime, nPeriod, nil
}

func (d *Dpos) inTurn(header *types.Header, witness common.Address, chain consensus.ChainReader, parents []*types.Header) bool {
	var (
		preWitness common.Address
		preTime    *big.Int
		manager    *Manager
		err        error
	)

	getHeaderFromParents := func(hash common.Hash, num uint64) *types.Header {
		if len(parents) == 0 {
			return nil
		}

		for i := len(parents) - 1; i >= 0; i-- {
			if parents[i].Hash() == hash && parents[i].Number.Uint64() == num {
				return parents[i]
			}
		}
		return nil
	}

	number := header.Number.Uint64()
	// using current block's witness list create manager
	if manager, err = d.manager(header); err != nil {
		log.Warn("Not find manager", "err", err.Error(), "number", number)
		return false
	}

	// first block always belongs to first witness
	if number == 1 {
		if len(manager.Witnesses) > 0 && manager.Witnesses[0] == witness {
			return true
		} else {
			return false
		}
	}

	preWitness, preTime, err = d.previousWitness(manager, chain, header.ParentHash, number-1, getHeaderFromParents)
	if err == errNoPreviousWitness {
		return manager.Witnesses[0] == witness
	} else if err != nil {
		log.Warn("Not find preWitness", "err", err.Error())
		return false
	}

	return manager.inTurn(witness, preWitness, header.Time, preTime)
}

// previousWitness find the previous witness who still in current witness list
func (d *Dpos) previousWitness(manager *Manager, chain consensus.ChainReader, hash common.Hash, number uint64, getHeaderFromParents getHeaderFromParentsFn) (witness common.Address, produceTime *big.Int, err error) {
	var header *types.Header

	find := false
	for !find {
		if number == 0 {
			return witness, produceTime, errNoPreviousWitness
		}

		header = getHeaderFromParents(hash, number)
		if header == nil {
			header = chain.GetHeader(hash, number)
		}
		if header == nil {
			return witness, produceTime, fmt.Errorf("can not find block header hash: %x, at hight: %d", hash, number)
		}

		// check is witness still valid
		witness = header.Coinbase
		find = manager.has(witness)
		if find {
			break
		}

		hash, number = header.ParentHash, number-1
	}

	// produce time from parent
	produceTime = header.Time

	return witness, produceTime, nil
}

// manager create a witness list manager using header
func (d *Dpos) manager(header *types.Header) (*Manager, error) {
	if len(header.Witnesses) == 0 {
		return nil, fmt.Errorf("header.Witnesses is empty")
	}

	// using header's witness list create manager
	return NewManager(d.config.Period, header.Witnesses), nil
}

// getWitnessesForProduce Get the first N candidates as witnesses from chain
func (d *Dpos) getWitnessesForProduce(header *types.Header, chain consensus.ChainReader, parent *types.Header) (bool, []common.Address, error) {
	var (
		bc *core.BlockChain
		ok bool
	)

	// get state db from parent's root
	if bc, ok = chain.(*core.BlockChain); !ok {
		return false, nil, fmt.Errorf("getWitnessesForProduce, get block chain instance error")
	}
	db, err := bc.StateAt(parent.Root)
	if db == nil {
		return false, nil, err
	}

	updated, witnesses := d.getWitnesses(header, db, parent)
	return updated, witnesses, nil
}

// getWitnesses 根据当前情况，判断从指定的state db读取或者使用前一个区块的
func (d *Dpos) getWitnesses(header *types.Header, db *state.StateDB, parent *types.Header) (bool, []common.Address) {
	var (
		witnesses      []common.Address
		lastUpdateTime *big.Int
		urls           []string
	)

	// Get last update witnesses list time from parent block
	if parent.Number.Int64() == 0 {
		lastUpdateTime = parent.Time
	} else {
		var upTime updateTime
		copy(upTime[:], parent.Extra[:updateTimeLen])
		lastUpdateTime = upTime.bigInt()
	}

	need := d.needUpdateWitnesses(header.Time, lastUpdateTime)
	if need {
		log.Debug("Get new witness from db", "height", header.Number.String())
		witnesses, urls = d.GetWitnessesFromStateDB(db)
	}

	// Using parent's witnesses, when update failed or No need update
	updated := need
	if len(witnesses) == 0 {
		witnesses = parent.Witnesses
		updated = false
	}
	if updated && d.sendBftPeerUpdateFn != nil {
		d.sendBftPeerUpdateFn(urls)
	}
	return updated, witnesses
}

// GetWitnessesFromStateDB Get the first N candidates as witnesses from stateDB
// It's can be used for get produce block and verify witnesses
func (d *Dpos) GetWitnessesFromStateDB(stateDB *state.StateDB) ([]common.Address, []string) {
	if stateDB == nil {
		log.Error("GetWitnessesFromStateDB, stateDB is nil")
	}

	return election.GetFirstNCandidates(stateDB, d.config.WitnessesNum)
}

// needUpdateWitnesses weather current time needs update witnesses list
func (d *Dpos) needUpdateWitnesses(t *big.Int, lastUpdateTime *big.Int) bool {
	log.Debug("needUpdateWitnesses", "last", lastUpdateTime.String(), "current", t.String())
	dur := new(big.Int).Sub(t, lastUpdateTime)
	return dur.Cmp(d.updateInterval) >= 0
}

// setUpdateInterval only called when start up
func (d *Dpos) setUpdateInterval() {
	d.updateInterval = new(big.Int).SetUint64(3 * uint64(d.config.WitnessesNum) * d.config.Period)
}

// coinBase get the address of this miner
func (d *Dpos) coinBase() (cb common.Address) {
	d.lock.RLock()
	cb = d.signer
	d.lock.RUnlock()
	return
}

// HandleBftMsg handle the bft message received from peer.
func (d *Dpos) HandleBftMsg(chain consensus.ChainReader, msg types.ConsensusMsg) {
	go d.bft.handleBftMsg(msg)
}

func (d *Dpos) CleanOldMsg(h *big.Int) {
	d.bft.cleanOldMsg(h)
}

func (d *Dpos) VerifyCommitMsg(block *types.Block) error {
	return d.bft.VerifyCmtMsgOf(block)
}

func (d *Dpos) ProducingStop() {
	d.bft.producingStop()
}

type updateTime [updateTimeLen]byte

func encodeUpdateTime(uTime *big.Int) []byte {
	var upTime updateTime
	binary.BigEndian.PutUint64(upTime[:], uTime.Uint64())
	return upTime[:]
}

func (upTime *updateTime) bigInt() *big.Int {
	uTime := binary.BigEndian.Uint64(upTime[:])
	return big.NewInt(0).SetUint64(uTime)
}

// calcVoteBounty returns a map, which contains the vote bonus of each candidates
// in this period. If active candidates less than WitnessesNum it will return a
// nil map.
func (d *Dpos) calcVoteBounty(candis election.CandidateList, allBonus *big.Int) map[common.Address]*big.Int {
	totalVotes := big.NewInt(0)
	activeCnt := 0
	for _, can := range candis {
		if !can.Active {
			continue
		}
		totalVotes.Add(totalVotes, can.VoteCount)
		activeCnt++
	}
	// Too less candidates before main net start, but it's normal
	if activeCnt < d.config.WitnessesNum || totalVotes.Cmp(common.Big0) == 0 {
		return nil
	}

	// Calc each candidates' bonus
	bonus := make(map[common.Address]*big.Int, activeCnt)
	for _, can := range candis {
		if !can.Active {
			continue
		}

		tmp := big.NewInt(0).Mul(allBonus, can.VoteCount)
		tmp.Div(tmp, totalVotes)
		bonus[can.Owner] = tmp
	}
	return bonus
}

// voteBonusPreWork
// 1) calculate vote bonus
// 2) get witness candidates of last update witness
func (d *Dpos) voteBonusPreWork(chain consensus.ChainReader, header *types.Header,
	curStateDB *state.StateDB) (election.CandidateList, *big.Int, error) {
	var (
		bc *core.BlockChain
		ok bool
	)
	// Block 1, no need bonus, it's no error
	if header.Number.Cmp(common.Big1) <= 0 {
		return make(election.CandidateList, 0), big.NewInt(0), nil
	}

	// Get state db from previous block
	if bc, ok = chain.(*core.BlockChain); !ok {
		return nil, nil, fmt.Errorf("voteBonusPreWork, get block chain instance error")
	}

	// Calc all vote bonus
	// the last block number of calculate vote reward is the last block number of updating witness list
	lastCalcBountyBlkNr := d.lastBountyBlkNr(header, bc)
	log.Debug("Bounus", "lastCalcBountyBlkNr", lastCalcBountyBlkNr.String())
	allBonus := big.NewInt(0).Sub(header.Number, lastCalcBountyBlkNr)
	if allBonus.Sign() <= 0 {
		return make(election.CandidateList, 0), big.NewInt(0), nil
	}
	allBonus.Mul(allBonus, curHeightBonus(header.Number, VortexCandidatesBonus))

	// Get all witnesses candidates
	lastCandis := election.GetAllCandidates(curStateDB, false)

	return lastCandis, allBonus, nil
}

// lastBountyBlkNr returns the block number of last update witness list. Returns
// the current block number if not find. This prevents excessive incentives.
// 见证人列表长时间未更新时，可能查找耗时，所以设置缓存数据，
// 同过度的下次查找将不再耗时。
// 不同高度的话，必然已经进行了下一轮见证人列表更新，以链上数据
// 为准，所以进行查找，然后记录数据。
func (d *Dpos) lastBountyBlkNr(header *types.Header, bc *core.BlockChain) (bh *big.Int) {
	// 数据老旧时，尝试更新数据
	d.lastBounty.RLock()
	outdated := d.lastBounty.updateHeight.Cmp(header.Number) < 0
	d.lastBounty.RUnlock()
	if outdated {
		h := bc.CurrentHeader()
		for !d.updatedWitnessCheckByTime(h) && h != nil {
			h = bc.GetHeaderByHash(h.ParentHash)
		}
		if h != nil {
			bh = new(big.Int).Set(h.Number)
		} else {
			bh = new(big.Int).Set(header.Number)
		}

		d.lastBounty.Lock()
		d.lastBounty.bountyHeight.Set(bh)
		d.lastBounty.updateHeight.Set(header.Number)
		d.lastBounty.Unlock()
		return
	}

	// 本高度已查找过，使用已查得数据
	d.lastBounty.RLock()
	bh = big.NewInt(0).Set(d.lastBounty.bountyHeight)
	d.lastBounty.RUnlock()

	return bh
}

// updatedWitnessCheckByTime using time to check whether this block updated
// witness list or not.
func (d *Dpos) updatedWitnessCheckByTime(header *types.Header) bool {
	if len(header.Extra) < updateTimeLen {
		return false
	}

	var upTime updateTime
	copy(upTime[:], header.Extra[:updateTimeLen])
	uTime := upTime.bigInt()
	return uTime.Cmp(header.Time) == 0
}

// curHeightBonus return the VNT bonus at blkNr block number.
func curHeightBonus(blkNr *big.Int, initBonus *big.Int) *big.Int {
	var denominator *big.Int
	if blkNr.Cmp(stageTwoBlkNr) < 0 {
		return initBonus
	} else if blkNr.Cmp(stageThreeBlkNr) < 0 {
		denominator = common.Big2
	} else {
		denominator = common.Big4
	}

	return big.NewInt(0).Div(initBonus, denominator)
}
