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
	"github.com/vntchain/go-vnt/core/vm/election"
	"github.com/vntchain/go-vnt/params"
)

func TestUpdateTime(t *testing.T) {
	ti := big.NewInt(10000)
	encoded := encodeUpdateTime(ti)
	var up updateTime
	copy(up[:], encoded[:])
	ret := up.bigInt()
	if ret.Cmp(ti) != 0 {
		t.Errorf("updateTime encode and decode error")
	}
}

func TestNeedUpdateWitnesses(t *testing.T) {
	cfg := &params.DposConfig{
		WitnessesNum: 3,
		Period:       2,
	}

	dp := New(cfg, nil)

	tests := []struct {
		cur  int64
		last int64
		ret  bool
	}{
		{1, 0, false},
		{10, 0, false},
		{17, 0, false},
		{18, 0, true},
		{19, 0, true},
		{20, 0, true},
		{21, 0, true},
		{-1, 0, false},
		{-10, 0, false},
		{-100, 0, false},
	}

	for i, ts := range tests {
		if dp.needUpdateWitnesses(big.NewInt(ts.cur), big.NewInt(ts.last)) != ts.ret {
			t.Errorf("test:%d, failed. cur = %d, last = %d, want=%v\n", i, ts.cur, ts.last, ts.ret)
		}
	}
}

func TestUpdatedWitnessCheckByTime(t *testing.T) {
	cfg := &params.DposConfig{
		WitnessesNum: 4,
		Period:       2,
	}

	dp := New(cfg, nil)

	updatedHeader := &types.Header{Time: big.NewInt(103023930)}
	updatedHeader.Extra = make([]byte, updateTimeLen)
	copy(updatedHeader.Extra, encodeUpdateTime(updatedHeader.Time))

	noneUpdatedHeader := &types.Header{Time: big.NewInt(103023930)}
	noneUpdatedHeader.Extra = make([]byte, updateTimeLen)
	copy(noneUpdatedHeader.Extra, encodeUpdateTime(big.NewInt(203930394)))

	tests := []struct {
		header *types.Header
		ret    bool
	}{
		{updatedHeader, true},
		{noneUpdatedHeader, false},
	}

	for i, ts := range tests {
		if dp.updatedWitnessCheckByTime(ts.header) != ts.ret {
			t.Errorf("test: %d failed", i)
		}
	}
}

func TestCurHeightBonus(t *testing.T) {
	blkNr1 := big.NewInt(100)
	blkNr2 := big.NewInt(57304000)
	blkNr3 := big.NewInt(104608000)

	tests := []struct {
		nr    *big.Int
		bonus *big.Int
	}{
		{blkNr1, big.NewInt(0).Set(VortexBlockReward)},
		{blkNr2, big.NewInt(0).Div(big.NewInt(0).Set(VortexBlockReward), big.NewInt(2))},
		{blkNr3, big.NewInt(0).Div(big.NewInt(0).Set(VortexBlockReward), big.NewInt(4))},
	}

	for i, ts := range tests {
		if ret := curHeightBonus(ts.nr, VortexBlockReward); ret.Cmp(ts.bonus) != 0 {
			t.Errorf("test: %d failed, want: %s, get: %s", i, ts.bonus.String(), ret.String())
		}
	}
}

func TestCalcVoteBounty(t *testing.T) {
	cfg := &params.DposConfig{
		WitnessesNum: 4,
		Period:       2,
	}

	dp := New(cfg, nil)

	// 情况1：人不够,但都是active
	// 情况2：人不够，有inactive
	// 情况3：没投票
	ca1 := election.Candidate{
		Owner:      common.BytesToAddress([]byte{1}),
		VoteCount:  big.NewInt(0),
		Registered: true,
		Bind:       true,
	}
	ca2 := election.Candidate{
		Owner:      common.BytesToAddress([]byte{2}),
		VoteCount:  big.NewInt(0),
		Registered: true,
		Bind:       true,
	}
	ca3 := election.Candidate{
		Owner:      common.BytesToAddress([]byte{3}),
		VoteCount:  big.NewInt(0),
		Registered: true,
		Bind:       true,
	}
	ca4 := election.Candidate{
		Owner:      common.BytesToAddress([]byte{4}),
		VoteCount:  big.NewInt(0),
		Registered: true,
		Bind:       true,
	}
	ca5 := election.Candidate{
		Owner:      common.BytesToAddress([]byte{5}),
		VoteCount:  big.NewInt(0),
		Registered: false,
	}
	ca6 := election.Candidate{
		Owner:      common.BytesToAddress([]byte{5}),
		VoteCount:  big.NewInt(0),
		Registered: true,
		Bind:       true,
	}

	candis1 := election.CandidateList{ca1, ca2, ca3}
	candis2 := election.CandidateList{ca1, ca2, ca3, ca5}
	candis3 := election.CandidateList{ca2, ca3, ca4, ca6}

	failTests := []election.CandidateList{candis1, candis2, candis3}
	allBonus := big.NewInt(100)
	rewardsExpEmpty := make(map[common.Address]*big.Int)
	for i, ts := range failTests {
		dp.calcVoteBounty(ts, allBonus, rewardsExpEmpty)
		if len(rewardsExpEmpty) != 0 {
			t.Errorf("test: %d, want: reward is empty, get: len(rewardsExpEmpty)=%v", i, len(rewardsExpEmpty))
		}
	}

	// 情况4：人够，按比例分配，其中1个inactive
	ca1.VoteCount = big.NewInt(10)
	ca2.VoteCount = big.NewInt(40)
	ca3.VoteCount = big.NewInt(20)
	ca4.VoteCount = big.NewInt(30)
	ca5.VoteCount = big.NewInt(15)
	candis4 := election.CandidateList{ca1, ca2, ca3, ca4, ca5}
	rewardsExpNoneEmpty := make(map[common.Address]*big.Int)
	rewardsExpNoneEmpty[ca1.Owner] = big.NewInt(3) // ca1有3个的已有激励
	dp.calcVoteBounty(candis4, allBonus, rewardsExpNoneEmpty)
	if len(rewardsExpNoneEmpty) == 0 {
		t.Errorf("want: rewardsExpNoneEmpty, get: len(rewardsExpNoneEmpty) == 0")
	}
	for i, ca := range candis4 {
		if !ca.Active() {
			continue
		}
		caBonus := rewardsExpNoneEmpty[ca.Owner]
		if caBonus == nil {
			t.Errorf("can %d, want bonus: %s, get: nil", i, ca.VoteCount.String())
			continue
		}
		if ca.Owner == ca1.Owner {
			// ca1的已有激励不能被覆盖
			if ca.VoteCount.Add(ca.VoteCount, big.NewInt(3)).Cmp(caBonus) != 0 {
				t.Errorf("can %d, want bonus: %s+3, get: %s", i, ca.VoteCount.String(), caBonus.String())
			}
		} else {
			if ca.VoteCount.Cmp(caBonus) != 0 {
				t.Errorf("can %d, want bonus: %s, get: %s", i, ca.VoteCount.String(), caBonus.String())
			}
		}

	}
}
