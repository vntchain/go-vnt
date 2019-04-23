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
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core/types"
	"github.com/vntchain/go-vnt/log"
)

// BFT step
const (
	newRound    = uint32(0)
	prePrePared = uint32(1)
	preparing   = uint32(2)
	prepared    = uint32(3)
	committing  = uint32(4)
	committed   = uint32(5)
	done        = uint32(6)
)

type BftManager struct {
	dp       *Dpos // DPoS object
	quorum   int   // 2f+1
	coinBase common.Address

	mp      *msgPool // message pool of all future round bft messages, and not verified
	roundMp *msgPool // message pool of current round, and been verified

	// BFT state
	h           *big.Int                    // local block chain height, protect by newRoundRWLock
	r           uint32                      // local BFT round, protect by newRoundRWLock
	step        uint32                      // local BFT round, protect by atomic operation
	witnessList map[common.Address]struct{} // current witness list, rely on producing

	newRoundRWLock sync.RWMutex // RW lock for switch to new round

	blockRound uint32 // round of sealing block, no need lock
	producing  uint32 // producing or not, atomic read and write

	// callbacks
	sendBftMsg  func(types.ConsensusMsg)
	verifyBlock func(*types.Block) (types.Receipts, []*types.Log, uint64, error)
	writeBlock  func(*types.Block) error
}

func newBftManager(dp *Dpos) *BftManager {
	n := dp.config.WitnessesNum
	q := n - (n-1)/3 // N-f
	return &BftManager{
		dp:          dp,
		quorum:      q,
		mp:          newMsgPool(q, "msg pool"),
		roundMp:     newMsgPool(q, "round msg pool"),
		h:           big.NewInt(0),
		r:           0,
		step:        newRound,
		witnessList: make(map[common.Address]struct{}, dp.config.WitnessesNum),
		producing:   0,
	}
}

// startPrePrepare will send pre-prepare msg and prepare msg
func (b *BftManager) startPrePrepare(block *types.Block) {
	log.Trace("Start PrePrepare")

	prePreMsg := b.makePrePrepareMsg(block, b.blockRound)

	// This node is a witness, which can seal block, no need check again
	b.sendBftMsg(prePreMsg)

	// startPrePrepare may be execute before newRound, so calling handleBftMsg
	// to check round
	go b.handleBftMsg(prePreMsg)
}

func (b *BftManager) handleBftMsg(msg types.ConsensusMsg) error {
	msgBlkNum, msgHash, msgType, msgRound := msg.GetBlockNum(), msg.Hash(), msg.Type(), msg.GetRound()
	if msgBlkNum == nil {
		return fmt.Errorf("bft msg's height is empty")
	}

	log.Trace("Handle bft msg", "type", msgType.String(),
		"number", msgBlkNum.Uint64(), "round", msgRound, "hash", msgHash)

	// Read lock for critical section
	b.newRoundRWLock.RLock()
	defer func() {
		b.newRoundRWLock.RUnlock()
		log.Trace("handle msg end", "hash", msgHash)
	}()

	// If not producing, no need deal with these bft msg right now,
	// save to msg pool, if you are witness and producing after sync,
	// you will fast deal with these msg.
	if atomic.LoadUint32(&b.producing) == 0 {
		_ = b.mp.addMsg(msg)
		log.Debug("HandleBftMsg: return for not producing")
		return nil
	}

	// 比较高度
	blkNumCmp := msgBlkNum.Cmp(b.h)
	if blkNumCmp > 0 {
		if err := b.mp.addMsg(msg); err != nil {
			log.Error("add prepare msg to msg pool error", "err", err)
			return err
		}
		if prePrepareMsg, ok := msg.(*types.PreprepareMsg); ok {
			go b.startSync(prePrepareMsg.Block)
		}
		return nil
	} else if blkNumCmp < 0 {
		// 比当前高度消息的msg，暂时先直接舍弃掉
		return fmt.Errorf("the height of msg is lower than current height, msg height :%d, current height : %d", msgBlkNum, b.h)
	}

	// 高度一致比较轮次
	if msgRound < b.r {
		return fmt.Errorf("the round of msg is lower than current round, msg round :%d, current round : %d", msgRound, b.r)
	} else if msgRound > b.r {
		if err := b.mp.addMsg(msg); err != nil {
			log.Error("add msg to msg pool error", "err", err)
			return err
		}
		return nil
	}

	switch msgType {
	case types.BftPreprepareMessage:
		return b.handlePrePrepareMsg(msg.(*types.PreprepareMsg))
	case types.BftPrepareMessage:
		return b.handlePrepareMsg(msg.(*types.PrepareMsg))
	case types.BftCommitMessage:
		return b.handleCommitMsg(msg.(*types.CommitMsg))
	default:
		log.Error("unknown bft message", "type", msgType.String())
		return fmt.Errorf("unknown bft message type: %s", msgType.String())
	}
}

// handlePrePrepareMsg
// WARN: msg pool has only one pre-prepare msg position, only the first correct one can be
// saved to msg pool. No need to care about you will vote for two pre-prepare msg. You only
// vote for the pre-prepare msg in msg pool.
func (b *BftManager) handlePrePrepareMsg(msg *types.PreprepareMsg) error {
	// check bft step, should be at newRound
	stp := atomic.LoadUint32(&b.step)
	if stp != newRound {
		log.Debug("Pre-prepare msg bft step not match", "local round", stp)
		return nil
	}

	// Verify msg
	if err := b.verifyPrePrepareMsg(msg); err != nil {
		log.Debug("Pre-prepare msg is invalid", "error", err)
		return err
	}
	if _, _, _, err := b.verifyBlock(msg.Block); err != nil {
		log.Debug("Pre-prepare's block is invalid", "error", err)
		return nil
	} else {
		log.Debug("Pre-prepare block is valid")
	}

	// Add msg to round msg pool, instead of msg pool
	if err := b.roundMp.addMsg(msg); err != nil {
		log.Warn("Add pre-prepare msg failed", "error", err)
		return err
	}

	// Go to next step
	if ok := atomic.CompareAndSwapUint32(&b.step, newRound, prePrePared); ok {
		return b.startPrepare()
	}
	return nil
}

// startPrepare enter prepare step, whether
func (b *BftManager) startPrepare() error {
	log.Trace("Start Prepare")
	// check our state and make a prepare msg
	if atomic.LoadUint32(&b.step) != prePrePared {
		return nil
	}
	prePreMsg, err := b.roundMp.getPrePrepareMsg(b.h, b.r)
	if err != nil {
		log.Error("Get pre-prepare msg from msg pool in prepare", "error", err)
		return fmt.Errorf("get pre-prepare msg from msg pool, error: %s", err)
	}
	preMsg, err := b.makePrepareMsg(prePreMsg)
	if err != nil {
		return err
	}
	if err := b.roundMp.addMsg(preMsg); err != nil {
		return err
	}

	// The one of first changing step, send the msg
	if ok := atomic.CompareAndSwapUint32(&b.step, prePrePared, preparing); ok {
		b.sendMsg(preMsg)
	}

	// Maybe have enough prepare msg to enter commit, when pre-prepare msg comes after prepare msg
	return b.tryCommitStep()
}

// tryCommitStep check whether can enter commit step
func (b *BftManager) tryCommitStep() error {
	stp := atomic.LoadUint32(&b.step)
	if stp < preparing || stp > prepared {
		log.Debug("tryCommitStep step not match", "step", stp)
		return nil
	}

	var (
		prePreMsg   *types.PreprepareMsg
		prepareMsgs []*types.PrepareMsg
		err         error
	)
	if prepareMsgs, err = b.roundMp.getTwoThirdMajorityPrepareMsg(b.h, b.r); err != nil {
		return nil
	}
	if prePreMsg, err = b.roundMp.getPrePrepareMsg(b.h, b.r); err != nil {
		log.Error("Get pre-prepare msg from msg pool in try commit", "error", err)
		return fmt.Errorf("get pre-prepare msg from msg pool, error: %s", err)
	}
	if prePreMsg.Block.Hash() != prepareMsgs[0].BlockHash {
		log.Error("Majority prepare msg is not match with pre-prepare msg", "block in prepare",
			prepareMsgs[0].BlockHash, "block in pre-prepare", prePreMsg.Block.Hash())
		return fmt.Errorf("majority prepare msg is not match with pre-prepare msg")
	}
	atomic.CompareAndSwapUint32(&b.step, preparing, prepared)
	return b.startCommit(prePreMsg)
}

// startCommit build commit message and send it
func (b *BftManager) startCommit(prePreMsg *types.PreprepareMsg) error {
	log.Trace("Start commit")
	if atomic.LoadUint32(&b.step) != prepared {
		return nil
	}

	commitMsg, err := b.makeCommitMsg(prePreMsg)
	if err != nil {
		return err
	}
	if err := b.roundMp.addMsg(commitMsg); err != nil {
		return err
	}

	// The one of first changing step, send the msg
	if ok := atomic.CompareAndSwapUint32(&b.step, prepared, committing); ok {
		b.sendMsg(commitMsg)
	}

	return b.tryWriteBlockStep()
}

// sendMsg only witness send bft message.
// Caller make sure has the newRoundRWLock.
func (b *BftManager) sendMsg(msg types.ConsensusMsg) {
	log.Trace("bft sendMsg start")
	if _, ok := b.witnessList[b.coinBase]; ok {
		b.sendBftMsg(msg)
	}
	log.Trace("bft sendMsg success")
}

func (b *BftManager) tryWriteBlockStep() error {
	if atomic.LoadUint32(&b.step) != committing {
		return nil
	}
	// commit消息满足数量要求后写区块
	if commitMsgs, err := b.roundMp.getTwoThirdMajorityCommitMsg(b.h, b.r); err != nil {
		return nil
	} else if ok := atomic.CompareAndSwapUint32(&b.step, committing, committed); ok {
		if prePrepareMsg, err := b.roundMp.getPrePrepareMsg(b.h, b.r); err != nil {
			log.Error("Get pre-prepare msg from msg pool in try write block", "err", err)
			return fmt.Errorf("get pre-prepare msg from msg pool, error: %s", err)
		} else {
			if err := b.writeBlockWithSig(prePrepareMsg, commitMsgs); err != nil {
				log.Error("Write block to chain", "err", err)
				return fmt.Errorf("write block to chain error: %s", err)
			}
			atomic.CompareAndSwapUint32(&b.step, committed, done)
		}
	}
	return nil
}

func (b *BftManager) handlePrepareMsg(msg *types.PrepareMsg) error {
	stp := atomic.LoadUint32(&b.step)
	// 当前阶段大于preparing，则prepare消息已经没用了，直接舍弃
	if stp > preparing {
		log.Debug("prepare msg bft step not match", "local step", stp)
		return nil
	}

	if err := b.verifyPrepareMsg(msg); err != nil {
		log.Error("failed to verify prepare msg", "err", err)
		return err
	}
	if err := b.roundMp.addMsg(msg); err != nil {
		log.Error("failed to add prepare msg", "height", b.h, "round", b.r, "err", err)
		return err
	}

	return b.tryCommitStep()
}

func (b *BftManager) handleCommitMsg(msg *types.CommitMsg) error {
	stp := atomic.LoadUint32(&b.step)
	// 当前阶段大于committing，则commit消息已经没用了，直接舍弃
	if stp > committing {
		log.Debug("commit msg bft step not match", "local step", stp)
		return nil
	}
	if err := b.verifyCommitMsg(msg); err != nil {
		log.Error("failed to verify commit msg", "err", err)
		return err
	}
	if err := b.roundMp.addMsg(msg); err != nil {
		log.Error("failed to add commit msg", "height", b.h, "round", b.r, "err", err)
		return err
	}

	return b.tryWriteBlockStep()
}

// writeBlock to block chain
func (b *BftManager) writeBlockWithSig(msg *types.PreprepareMsg, cmtMsg []*types.CommitMsg) error {
	block := msg.Block
	// Match pre-prepare msg and commit msg
	if block.Hash() != cmtMsg[0].BlockHash {
		return fmt.Errorf("writeBlockWithSig error, commit msg for block: %s, not for block: %s", cmtMsg[0].BlockHash.String(), block.Hash().String())
	}

	block.FillBftMsg(cmtMsg)
	log.Trace("writeBlockWithSig", "h", b.h.String(), "r", b.r, "hash", block.Hash().Hex())
	return b.writeBlock(block)
}

// newRound has lock, it maybe time consuming at sometime, call it by routine
func (b *BftManager) newRound(h *big.Int, r uint32, witList []common.Address) {
	log.Trace("New round switch start")
	b.newRoundRWLock.Lock()

	// Update witness list
	if b.h.Cmp(h) != 0 {
		b.witnessList = make(map[common.Address]struct{})
		for _, wit := range witList {
			b.witnessList[wit] = struct{}{}
		}
	}

	// Reset state and round msg pool
	b.h = h
	b.r = r
	b.step = newRound

	// Reset round msg pool
	b.roundMp.cleanAllMessage()

	// Switch to new round finished
	b.newRoundRWLock.Unlock()

	log.Trace("New round switch finish", "h", b.h.String(), "r", b.r, "time", time.Now().Unix())

	// New round switch finished, must return right now
	go b.importCurRoundMsg()
}

// importCurRoundMsg import consensus messages, but can not directly import to round msg pool
func (b *BftManager) importCurRoundMsg() {
	b.newRoundRWLock.RLock()
	msg := b.mp.getAllMsgOf(b.h, b.r)
	b.newRoundRWLock.RUnlock()
	for _, m := range msg {
		log.Trace("Import Msg", "type", m.Type(), "hash", m.Hash())
		go b.handleBftMsg(m)
	}
}

func (b *BftManager) startSync(block *types.Block) {
	// 	TODO Not emergency, sync will be triggered by block hash msg
	log.Debug("Bft manager startPrePrepare sync")
}

func (b *BftManager) producingStop() {
	atomic.StoreUint32(&b.producing, 0)
	log.Debug("BFT stop producing")
}

func (b *BftManager) producingStart() {
	atomic.StoreUint32(&b.producing, 1)
	log.Debug("BFT start producing")
}

func (b *BftManager) validWitness(wit common.Address) bool {
	_, ok := b.witnessList[wit]
	return ok
}

// cleanOldMsg clean msg pool and keep future message. cleaning only
// on height % 100 == 0.
func (b *BftManager) cleanOldMsg(h *big.Int) {
	b.mp.cleanOldMessage(h)
}
