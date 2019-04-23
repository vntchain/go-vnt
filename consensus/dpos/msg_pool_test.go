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
	"math/big"
	"testing"

	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core/types"
)

func TestMsgPool_GetOrNewRoundMsgPool(t *testing.T) {
	// n = 4, f = 1
	quo := 3
	mp := newMsgPool(quo, "test")
	h := big.NewInt(1)
	r := uint32(0)

	if _, ok := mp.pool[h.Uint64()]; ok {
		t.Errorf("mp should empty")
	}
	if rmp := mp.getOrNewRoundMsgPool(h, r); rmp == nil {
		t.Errorf("rmp should be valid")
	}

	// should have
	if _, ok := mp.pool[h.Uint64()]; !ok {
		t.Errorf("mp should have h manger")
	}

	if _, ok := mp.pool[h.Uint64()].pool[r]; !ok {
		t.Errorf("mp should have r manager")
	}

	// make Pre-prepareMsg as a tag
	mp.pool[h.Uint64()].pool[r].prePreMsg = &types.PreprepareMsg{Round: r, Block: nil}

	if rmp, err := mp.getRoundMsgPool(h, r); err != nil || rmp == nil {
		t.Errorf("get roundMsgPool failed, %s", err)
	} else {
		if rmp.prePreMsg == nil {
			t.Errorf("pre prepare msg should not be nil")
		} else {
			t.Logf("tag round: %d", rmp.prePreMsg.Round)
		}
	}
}

func TestMsgPool_AddMsgAndMajoritySuccess(t *testing.T) {
	// n = 4, f = 1
	quo := 3
	mp := newMsgPool(quo, "test")
	h := big.NewInt(1)
	r := uint32(0)

	header := &types.Header{
		Number:   big.NewInt(0).Set(h),
		Coinbase: common.BytesToAddress([]byte{1}),
	}
	b := types.NewBlock(header, nil, nil)

	// add 1 pre-prepare msg
	prePre := &types.PreprepareMsg{
		Round: r,
		Block: b,
	}
	if err := mp.addMsg(prePre); err != nil {
		t.Errorf("add pre prepare message error: %s", err)
	}

	if msg, err := mp.getPrePrepareMsg(h, r); msg == nil || err != nil {
		t.Errorf("should success, but got: %s", err)
	} else {
		t.Logf("msg: %d", msg.Round)
	}

	header1 := &types.Header{
		Number:   big.NewInt(0).Set(h),
		Coinbase: common.BytesToAddress([]byte{2}),
	}
	b1 := types.NewBlock(header1, nil, nil)

	prePre1 := &types.PreprepareMsg{
		Round: r,
		Block: b1,
	}

	// can not add pre-prepare msg again
	if err := mp.addMsg(prePre1); err == nil {
		t.Errorf("Adding pre-premessage the second time should failed")
	}

	// add 2 prepare message for same block hash
	preMsg := &types.PrepareMsg{
		BlockNumber: h,
		Round:       0,
		BlockHash:   common.HexToHash("503290d0c4dd2d72202521e4701e89daecf048d400b2fbb8cbad1f15a4ec2e8d"),
		PrepareAddr: common.BytesToAddress([]byte{1}),
	}
	if err := mp.addMsg(preMsg); err != nil {
		t.Errorf("should success, but get: %s", err)
	}
	// 相同消息,添加会失败
	if err := mp.addMsg(preMsg); err == nil {
		t.Errorf("Adding same message the second time should failed")
	}

	preMsg1 := &types.PrepareMsg{
		BlockNumber: h,
		Round:       0,
		BlockHash:   common.HexToHash("503290d0c4dd2d72202521e4701e89daecf048d400b2fbb8cbad1f15a4ec2e8d"),
		PrepareAddr: common.BytesToAddress([]byte{2}),
	}
	if err := mp.addMsg(preMsg1); err != nil {
		t.Errorf("should success, but get: %s", err)
	}

	// can not get majority prepare message
	if _, err := mp.getTwoThirdMajorityPrepareMsg(h, r); err == nil {
		t.Errorf("should failed, but success")
	}

	preMsg2 := &types.PrepareMsg{
		BlockNumber: h,
		Round:       0,
		BlockHash:   common.HexToHash("503290d0c4dd2d72202521e4701e89daecf048d400b2fbb8cbad1f15a4ec2e8d"),
		PrepareAddr: common.BytesToAddress([]byte{3}),
	}
	// add 1 prepare message for the same block hash
	if err := mp.addMsg(preMsg2); err != nil {
		t.Errorf("should success, but get: %s", err)
	}

	// success get majority prepare message
	if msgs, err := mp.getTwoThirdMajorityPrepareMsg(h, r); err != nil {
		t.Errorf("should success, but get: %s", err)
	} else {
		for i, m := range msgs {
			t.Logf("prepare msg: %d, for block: %s", i, m.BlockHash.String())
		}
	}

	// clean messages of previous height
	_ = mp.cleanMsgOfHeight(h)

	// below should failed
	if rmp, err := mp.getRoundMsgPool(h, r); rmp != nil || err == nil {
		t.Errorf("should failed, but success")
	}
	if pp, err := mp.getPrePrepareMsg(h, r); pp != nil || err == nil {
		t.Errorf("should failed, but success")
	}
}

func TestMsgPool_MajorityFail(t *testing.T) {
	// n = 4, f = 1
	quo := 3
	mp := newMsgPool(quo, "test")
	h := big.NewInt(1)
	r := uint32(0)

	// add 2 prepare message for same block hash
	preMsg := &types.PrepareMsg{
		BlockNumber: h,
		Round:       0,
		BlockHash:   common.HexToHash("503290d0c4dd2d72202521e4701e89daecf048d400b2fbb8cbad1f15a4ec2e8d"),
	}
	if err := mp.addMsg(preMsg); err != nil {
		t.Errorf("should success, but get: %s", err)
	}
	if err := mp.addMsg(preMsg); err == nil {
		t.Errorf("Adding same message the second time should failed")
	}

	// can not get majority prepare message
	if _, err := mp.getTwoThirdMajorityPrepareMsg(h, r); err == nil {
		t.Errorf("should failed, but success")
	}

	preMsg1 := &types.PrepareMsg{
		BlockNumber: h,
		Round:       0,
		BlockHash:   common.HexToHash("503290d0c4dd2d72202521e4701e89daecf048d400b2fbb8cbad1f15a4ec2e8d"),
		PrepareAddr: common.BytesToAddress([]byte{2}),
	}
	if err := mp.addMsg(preMsg1); err != nil {
		t.Errorf("should success, but get: %s", err)
	}

	// get 1 prepare message for the same block hash
	preMsg2 := &types.PrepareMsg{
		BlockNumber: h,
		Round:       0,
		BlockHash:   common.HexToHash("233290d0c4dd2d72202521e4701e89daecf048d400b2fbb8cbad1f15a4ec2e8d"),
		PrepareAddr: common.BytesToAddress([]byte{3}),
	}
	if err := mp.addMsg(preMsg2); err != nil {
		t.Errorf("should success, but get: %s", err)
	}

	// fail get majority prepare message
	if pms, err := mp.getTwoThirdMajorityPrepareMsg(h, r); pms != nil || err == nil {
		t.Errorf("should fail, but success. err: %v", err)
		for i, m := range pms {
			t.Logf("prepare msg: %d, block: %s\n", i, m.BlockHash.Hex())
		}
	}
}

func TestMsgPool_cleanOldMessage(t *testing.T) {
	// n = 4, f = 1
	quo := 3
	mp := newMsgPool(quo, "test")
	h1 := big.NewInt(100)
	h2 := big.NewInt(101)
	r := uint32(0)

	hd1 := &types.Header{
		Number: h1,
	}
	hd2 := &types.Header{
		Number: h2,
	}
	b1 := types.NewBlockWithHeader(hd1)
	b2 := types.NewBlockWithHeader(hd2)
	// add msg of height 100 and 101 to msg pool
	preMsg := &types.PreprepareMsg{
		Block: b1,
		Round: 0,
	}
	if err := mp.addMsg(preMsg); err != nil {
		t.Errorf("should success, but get: %s", err)
	}
	preMsg2 := &types.PreprepareMsg{
		Block: b2,
		Round: 0,
	}
	if err := mp.addMsg(preMsg2); err != nil {
		t.Errorf("should success, but get: %s", err)
	}

	// 	get preMsg && preMsg2 should success
	if msg, err := mp.getPrePrepareMsg(h1, r); err != nil || msg == nil {
		t.Errorf("Should get preMsg, but failed, err: %s", err)
	} else if msg.Hash() != preMsg.Hash() {
		t.Errorf("want header: %x, got: %x", preMsg.Hash(), msg.Hash())
	}
	if msg, err := mp.getPrePrepareMsg(h2, r); err != nil || msg == nil {
		t.Errorf("Should get preMsg2, but failed, err: %s", err)
	}

	// clean message
	mp.cleanOldMessage(h1)

	// get preMsg should failed
	if msg, err := mp.getPrePrepareMsg(h1, r); err == nil || msg != nil {
		t.Errorf("Still get preMsg, cleanOldMessage failed, err: %s, number: %s", err, msg.GetBlockNum().String())
	}
	// 	get preMsg2 should success
	if msg, err := mp.getPrePrepareMsg(h2, r); err != nil || msg == nil {
		t.Errorf("Should still get preMsg2, but failed, err: %s", err)
	}
}

func TestCleanOldMsg(t *testing.T) {
	// n = 4, f = 1
	quo := 3
	mp := newMsgPool(quo, "test")

	// less
	preMsg1 := &types.PrepareMsg{
		BlockNumber: big.NewInt(msgCleanInterval - 10),
		Round:       0,
		BlockHash:   common.HexToHash("503290d0c4dd2d72202521e4701e89daecf048d400b2fbb8cbad1f15a4ec2e8d"),
		PrepareAddr: common.BytesToAddress([]byte{2}),
	}
	if err := mp.addMsg(preMsg1); err != nil {
		t.Errorf("should success, but get: %s", err)
	}

	// equal
	preMsg2 := &types.PrepareMsg{
		BlockNumber: big.NewInt(msgCleanInterval),
		Round:       0,
		BlockHash:   common.HexToHash("503290d0c4dd2d72202521e4701e89daecf048d400b2fbb8cbad1f15a4ec2e8d"),
		PrepareAddr: common.BytesToAddress([]byte{2}),
	}
	if err := mp.addMsg(preMsg2); err != nil {
		t.Errorf("should success, but get: %s", err)
	}

	// bigger
	preMsg3 := &types.PrepareMsg{
		BlockNumber: big.NewInt(msgCleanInterval + 10),
		Round:       0,
		BlockHash:   common.HexToHash("503290d0c4dd2d72202521e4701e89daecf048d400b2fbb8cbad1f15a4ec2e8d"),
		PrepareAddr: common.BytesToAddress([]byte{2}),
	}
	if err := mp.addMsg(preMsg3); err != nil {
		t.Errorf("should success, but get: %s", err)
	}

	// clean
	mp.cleanOldMessage(big.NewInt(msgCleanInterval))

	// check result
	if _, ok := mp.msgHashSet[preMsg1.Hash()]; ok {
		t.Errorf("premsg1 should not exist")
	}
	if rmp, _ := mp.getRoundMsgPool(preMsg1.BlockNumber, 0); rmp != nil {
		t.Errorf("rmp of (%d,%d) should be nil", preMsg1.BlockNumber.Uint64(), preMsg1.Round)
	}

	if _, ok := mp.msgHashSet[preMsg2.Hash()]; ok {
		t.Errorf("premsg2 should not exist")
	}
	if rmp, _ := mp.getRoundMsgPool(preMsg2.BlockNumber, 0); rmp != nil {
		t.Errorf("rmp of (%d,%d) should be nil", preMsg2.BlockNumber.Uint64(), preMsg2.Round)
	}

	if _, ok := mp.msgHashSet[preMsg3.Hash()]; !ok {
		t.Errorf("premsg3 should exist, but not, before clean")
	}
	if rmp, _ := mp.getRoundMsgPool(preMsg3.BlockNumber, 0); rmp == nil {
		t.Errorf("rmp of (%d,%d) should be valid, but nil", preMsg3.BlockNumber.Uint64(), preMsg3.Round)
	} else {
		find := false
		if len(rmp.preMsgs) == 0 {
			t.Errorf("round msg pool is empty")
		}
		for _, msg := range rmp.preMsgs {
			if msg.Hash() == preMsg3.Hash() {
				find = true
				break
			}
		}
		if find == false {
			t.Errorf("premsg3 should exist, but not in list")
		}
	}
}
