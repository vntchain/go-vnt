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
	"github.com/pkg/errors"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core/types"
	"github.com/vntchain/go-vnt/log"
	"math/big"
	"sync"
)

const (
	bftMsgBufSize    = 30
	msgCleanInterval = 100
)

// msgPool store all bft consensus message of each height, and these message grouped by height.
type msgPool struct {
	name       string
	pool       map[uint64]*heightMsgPool
	quorum     int // 2f+1
	lock       sync.RWMutex
	msgHashSet map[common.Hash]uint64 //value为高度方便按高度进行删除
}

func newMsgPool(q int, n string) *msgPool {
	mp := &msgPool{
		name:       n,
		pool:       make(map[uint64]*heightMsgPool),
		quorum:     q,
		msgHashSet: make(map[common.Hash]uint64),
	}
	return mp
}

func (mp *msgPool) addMsg(msg types.ConsensusMsg) error {
	msgHash := msg.Hash()
	h := msg.GetBlockNum()
	if h == nil {
		return fmt.Errorf("addMsg msg's height is nil, msg: %s", msgHash.Hex())
	}
	r := msg.GetRound()

	mp.lock.Lock()
	defer mp.lock.Unlock()

	if _, exists := mp.msgHashSet[msgHash]; exists {
		return fmt.Errorf("addMsg msg already exists, msg: %s", msgHash.Hex())
	}

	rmp := mp.getOrNewRoundMsgPool(h, r)

	if err := rmp.addMsg(msg); err != nil {
		log.Warn("Msg pool add msg failed", "pool name", mp.name, "msg type", msg.Type().String(), "error", err)
		return err
	}

	mp.msgHashSet[msgHash] = msg.GetBlockNum().Uint64()
	return nil
}

func (mp *msgPool) getPrePrepareMsg(h *big.Int, r uint32) (*types.PreprepareMsg, error) {
	mp.lock.RLock()
	defer mp.lock.RUnlock()

	rmp, err := mp.getRoundMsgPool(h, r)
	if err != nil {
		return nil, err
	}

	if rmp.prePreMsg == nil {
		return nil, fmt.Errorf("round (%d,%d) has no pre-prepare msg", h.Uint64(), r)
	}
	return rmp.prePreMsg, nil
}

func (mp *msgPool) getAllMsgOf(h *big.Int, r uint32) []types.ConsensusMsg {
	msg := make([]types.ConsensusMsg, 0, bftMsgBufSize*2+1)

	mp.lock.RLock()
	defer mp.lock.RUnlock()

	rmp, _ := mp.getRoundMsgPool(h, r)
	if rmp == nil {
		return msg
	}

	if rmp.prePreMsg != nil {
		msg = append(msg, rmp.prePreMsg)
	}
	for _, m := range rmp.preMsgs {
		msg = append(msg, m)
	}
	for _, m := range rmp.commitMsgs {
		msg = append(msg, m)
	}
	return msg
}

// getTwoThirdMajorityPrepareMsg get the majority prepare message, and the count of these
// message must is bigger than 2f. otherwise, return nil, nil
func (mp *msgPool) getTwoThirdMajorityPrepareMsg(h *big.Int, r uint32) ([]*types.PrepareMsg, error) {
	mp.lock.RLock()
	defer mp.lock.RUnlock()

	rmp, _ := mp.getRoundMsgPool(h, r)
	if rmp == nil {
		return nil, nil
	}
	msgs := rmp.preMsgs

	// too less commit message
	if len(msgs) < mp.quorum {
		return nil, errors.New("too less prepare message")
	}

	// count
	cnt := make(map[common.Hash]int)
	var maxCntHash common.Hash
	maxCnt := 0
	for _, msg := range msgs {
		bh := msg.BlockHash
		if _, ok := cnt[bh]; !ok {
			cnt[bh] = 1
		} else {
			cnt[bh] += 1
		}

		if cnt[bh] > maxCnt {
			maxCnt = cnt[bh]
			maxCntHash = bh
		}
	}

	// not enough
	if maxCnt < mp.quorum {
		return nil, errors.New("majority prepare message is too less")
	}

	// get prepare massage
	matchedMsgs := make([]*types.PrepareMsg, 0, maxCnt)
	for _, msg := range msgs {
		if msg.BlockHash == maxCntHash {
			matchedMsgs = append(matchedMsgs, msg)
		}
	}
	return matchedMsgs, nil
}

// getTwoThirdMajorityCommitMsg get the majority commit message, and the count of these
// // message must is bigger than 2f. otherwise, return nil, nil
func (mp *msgPool) getTwoThirdMajorityCommitMsg(h *big.Int, r uint32) ([]*types.CommitMsg, error) {
	mp.lock.RLock()
	defer mp.lock.RUnlock()

	rmp, _ := mp.getRoundMsgPool(h, r)
	if rmp == nil {
		return nil, nil
	}
	msgs := rmp.commitMsgs

	// too less commit message
	if len(msgs) < mp.quorum {
		return nil, errors.New("too less commit message")
	}

	// count
	cnt := make(map[common.Hash]int)
	var maxCntHash common.Hash
	maxCnt := 0
	for _, msg := range msgs {
		bh := msg.BlockHash
		if _, ok := cnt[bh]; !ok {
			cnt[bh] = 1
		} else {
			cnt[bh] += 1
		}

		if cnt[bh] > maxCnt {
			maxCnt = cnt[bh]
			maxCntHash = bh
		}
	}

	// not enough
	if maxCnt < mp.quorum {
		return nil, errors.New("majority commit message is too less")
	}

	// get prepare massage
	matchedMsgs := make([]*types.CommitMsg, 0, maxCnt)
	for _, msg := range msgs {
		if msg.BlockHash == maxCntHash {
			matchedMsgs = append(matchedMsgs, msg)
		}
	}
	return matchedMsgs, nil
}

// getOrNewRoundMsgPool if round msg pool not exist, it will create.
// WARN: caller should lock the msg pool
func (mp *msgPool) getOrNewRoundMsgPool(h *big.Int, r uint32) *roundMsgPool {
	uh := h.Uint64()
	if _, ok := mp.pool[uh]; !ok {
		mp.pool[uh] = newHeightMsgPool()
	}
	hmp := mp.pool[uh]
	if _, ok := hmp.pool[r]; !ok {
		hmp.pool[r] = newRoundMsgPool()
	}
	return hmp.pool[r]
}

// getRoundMsgPool just to get round message pool. If not exist, it will return error
// WARN: caller should lock the msg pool
func (mp *msgPool) getRoundMsgPool(h *big.Int, r uint32) (*roundMsgPool, error) {
	uh := h.Uint64()
	if _, ok := mp.pool[uh]; !ok {
		return nil, fmt.Errorf("hight manager is nil, h: %d", uh)
	}
	hmp := mp.pool[uh]
	if _, ok := hmp.pool[r]; !ok {
		return nil, fmt.Errorf("round manager is nil, (h,r): (%d, %d)", uh, r)
	}
	return hmp.pool[r], nil
}

func (mp *msgPool) cleanMsgOfHeight(h *big.Int) error {
	mp.lock.Lock()
	defer mp.lock.Unlock()

	delete(mp.pool, h.Uint64())
	for k, height := range mp.msgHashSet {
		if h.Uint64() == height {
			delete(mp.msgHashSet, k)
		}
	}

	return nil
}

func (mp *msgPool) cleanAllMessage() {
	mp.lock.Lock()
	defer mp.lock.Unlock()

	mp.pool = make(map[uint64]*heightMsgPool)
	mp.msgHashSet = make(map[common.Hash]uint64)

}

func (mp *msgPool) cleanOldMessage(h *big.Int) {
	uh := h.Uint64()

	if uh%msgCleanInterval == 0 {
		log.Debug("Message pool clean old message")
		mp.lock.Lock()
		defer mp.lock.Unlock()

		oldPool := mp.pool
		oldHashSet := mp.msgHashSet
		mp.pool = make(map[uint64]*heightMsgPool)
		mp.msgHashSet = make(map[common.Hash]uint64)
		for mh, hp := range oldPool {
			if mh > uh {
				mp.pool[mh] = hp
			}
		}
		for k, h := range oldHashSet {
			if h > uh {
				mp.msgHashSet[k] = h
			}
		}
		log.Debug("Message pool clean old message done", "num. of height cleaned", len(oldPool)-len(mp.pool))
	}
}

// heightMsgPool store all bft message of each height, and these message grouped by round index.
// WARN: heightMsgPool do not support lock, but MsgPool support lock
type heightMsgPool struct {
	pool map[uint32]*roundMsgPool
}

func newHeightMsgPool() *heightMsgPool {
	return &heightMsgPool{
		pool: make(map[uint32]*roundMsgPool),
	}
}

func (hmp *heightMsgPool) addMsg(msg types.ConsensusMsg) error {
	r := msg.GetRound()
	if _, ok := hmp.pool[r]; !ok {
		hmp.pool[r] = newRoundMsgPool()
	}
	return hmp.pool[r].addMsg(msg)
}

// roundMsgPool store all bft message of each round, and these message grouped by message type.
// WARN: heightMsgPool do not support lock, but MsgPool support lock
type roundMsgPool struct {
	prePreMsg  *types.PreprepareMsg
	preMsgs    []*types.PrepareMsg
	commitMsgs []*types.CommitMsg
}

func newRoundMsgPool() *roundMsgPool {
	return &roundMsgPool{
		prePreMsg:  nil,
		preMsgs:    make([]*types.PrepareMsg, 0, bftMsgBufSize),
		commitMsgs: make([]*types.CommitMsg, 0, bftMsgBufSize),
	}
}

func (rmp *roundMsgPool) addMsg(msg types.ConsensusMsg) error {
	switch msg.Type() {
	case types.BftPreprepareMessage:
		if rmp.prePreMsg == nil {
			// fmt.Println("not save prepre")
			rmp.prePreMsg = msg.(*types.PreprepareMsg)
		} else {
			return fmt.Errorf("already save a pre-prepare msg at round: (%d,%d), added: %s, adding: %s",
				msg.GetBlockNum().Uint64(), msg.GetRound(), rmp.prePreMsg.Hash().Hex(), msg.Hash().Hex())
		}

	case types.BftPrepareMessage:
		rmp.preMsgs = append(rmp.preMsgs, msg.(*types.PrepareMsg))

	case types.BftCommitMessage:
		rmp.commitMsgs = append(rmp.commitMsgs, msg.(*types.CommitMsg))

	default:
		return fmt.Errorf("unknow bft message type: %d, hash: %s", msg.Type(), msg.Hash().Hex())
	}
	return nil
}

func (rmp *roundMsgPool) clean() {
	rmp.prePreMsg = nil
	rmp.preMsgs = make([]*types.PrepareMsg, 0, bftMsgBufSize)
	rmp.commitMsgs = make([]*types.CommitMsg, 0, bftMsgBufSize)
}
