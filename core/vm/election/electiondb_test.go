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

package election

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"

	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core/state"
	"github.com/vntchain/go-vnt/vntdb"
)

var (
	voter = Voter{
		Owner: common.HexToAddress("9ee97d274eb4c215f23238fee1f103d9ea10a234"),
		VoteCandidates: []common.Address{
			common.HexToAddress("9ee97d274eb4c215f23238fee1f103d9ea10a234"),
			common.BytesToAddress([]byte{10}),
		},
		ProxyVoteCount: big.NewInt(102),
		LastVoteCount:  big.NewInt(5),
		IsProxy:        true,
		TimeStamp:      big.NewInt(1531454152),
	}
	candidate = Candidate{
		Owner:       common.HexToAddress("9ee97d274eb4c215f23238fee1f103d9ea10a234"),
		Binder:      binder,
		Beneficiary: beneficiary,
		Registered:  true,
		Bind:        true,
		VoteCount:   big.NewInt(0),
		Url:         []byte("/ip4/192.168.9.102/tcp/5210/ipfs/1kHaMUmZgTpjGEhxcGATr1UVWy6iKkygFuknWEtW7LiLrev"),
		Website:     []byte("www.testwebsite.net/test/witness/website"),
		Name:        []byte("testNet"),
	}
	stake = Stake{
		Owner:      common.HexToAddress("9ee97d274eb4c215f23238fee1f103d9ea10a234"),
		StakeCount: big.NewInt(230),
		TimeStamp:  big.NewInt(1531454152),
	}
	bounty = Reward{
		Rest: big.NewInt(1e18),
	}
)

func print(key common.Hash, value common.Hash) {
	fmt.Printf("%x:%x\n", key, value)
}

func sameVoter(voter *Voter, voter1 *Voter) (bool, error) {
	if !bytes.Equal(voter.Owner[:], voter1.Owner[:]) {
		return false, fmt.Errorf("Error, owner before %v and after %v is different", voter.Owner, voter1.Owner)
	} else if voter.TimeStamp.Cmp(voter1.TimeStamp) != 0 {
		return false, fmt.Errorf("Error, timestamp before %v and after %v is different", voter.TimeStamp, voter1.TimeStamp)
	} else if voter.Proxy != voter1.Proxy {
		return false, fmt.Errorf("Error, proxy before %v and after %v is different", voter.Proxy, voter1.Proxy)
	} else if voter.IsProxy != voter1.IsProxy {
		return false, fmt.Errorf("Error, isproxy before %v and after %v is different", voter.IsProxy, voter1.IsProxy)
	} else if voter.ProxyVoteCount.Cmp(voter1.ProxyVoteCount) != 0 {
		return false, fmt.Errorf("Error, ProxyVoteCount before %v and after %v is different", voter.ProxyVoteCount, voter1.ProxyVoteCount)
	} else if voter.LastVoteCount.Cmp(voter1.LastVoteCount) != 0 {
		return false, fmt.Errorf("Error, LastVoteCount before %v and after %v is different", voter.LastVoteCount, voter1.LastVoteCount)
	} else {
		if len(voter.VoteCandidates) != len(voter1.VoteCandidates) {
			return false, fmt.Errorf("Error, the length of VoteCandidates before %v and after %v is different", len(voter.VoteCandidates), len(voter1.VoteCandidates))
		} else {
			for i, candi := range voter.VoteCandidates {
				if !bytes.Equal(candi.Bytes(), voter1.VoteCandidates[i].Bytes()) {
					return false, fmt.Errorf("Error,  VoteCandidates[%d] before %v and after %v is different", i, candi, voter1.VoteCandidates[i])
				}
			}
		}
	}
	return true, nil
}

func sameCandidate(candidate *Candidate, candidate1 *Candidate) (bool, error) {
	if !bytes.Equal(candidate.Owner[:], candidate1.Owner[:]) {
		return false, fmt.Errorf("Error,owner before %v and after %v is different", candidate.Owner, candidate1.Owner)
	} else if candidate.Registered != candidate1.Registered {
		return false, fmt.Errorf("Error,registered before %v and after %v is different", candidate.Registered, candidate1.Registered)
	} else if candidate.Binder != candidate1.Binder {
		return false, fmt.Errorf("Error,binder before %v and after %v is different", candidate.Binder, candidate1.Binder)
	} else if candidate.Bind != candidate1.Bind {
		return false, fmt.Errorf("Error,bind before %v and after %v is different", candidate.Bind, candidate1.Bind)
	} else if candidate.Beneficiary != candidate1.Beneficiary {
		return false, fmt.Errorf("Error,beneficiary before %v and after %v is different", candidate.Beneficiary, candidate1.Beneficiary)
	} else if candidate.VoteCount.Cmp(candidate1.VoteCount) != 0 {
		return false, fmt.Errorf("Error,voteCount before %v and after %v is different", candidate.VoteCount, candidate1.VoteCount)
	} else if !bytes.Equal(candidate.Url, candidate1.Url) {
		return false, fmt.Errorf("Error, url before %x and after %x is different", candidate.Url, candidate1.Url)
	} else if !bytes.Equal(candidate.Website, candidate1.Website) {
		return false, fmt.Errorf("Error, Website before %v and after %v is different", candidate.Website, candidate1.Website)
	} else if !bytes.Equal(candidate.Name, candidate1.Name) {
		return false, fmt.Errorf("Error, Name before %v and after %v is different", candidate.Name, candidate1.Name)
	}
	return true, nil
}

func sameStake(stake *Stake, stake1 *Stake) (bool, error) {
	if !bytes.Equal(stake.Owner[:], stake1.Owner[:]) {
		return false, fmt.Errorf("Error, owner before %v and after %v is different", stake.Owner, stake1.Owner)
	} else if stake.StakeCount.Cmp(stake1.StakeCount) != 0 {
		return false, fmt.Errorf("Error, stakeCount before %v and after %v is different", stake.StakeCount, stake1.StakeCount)
	} else if stake.TimeStamp.Cmp(stake1.TimeStamp) != 0 {
		return false, fmt.Errorf("Error, timestamp before %v and after %v is different", stake.TimeStamp, stake1.TimeStamp)
	}
	return true, nil
}

func TestConvertToKV(t *testing.T) {
	err := convertToKV(VOTERPREFIX, voter, print)
	if err != nil {
		t.Error(err)
	}

	err = convertToKV(CANDIDATEPREFIX, candidate, print)
	if err != nil {
		t.Error(err)
	}

	err = convertToKV(STAKEPREFIX, stake, print)
	if err != nil {
		t.Error(err)
	}

	err = convertToKV(REWARDPREFIX, bounty, print)
	if err != nil {
		t.Error(err)
	}
}

// 把上面单元测试打印的结果添加到下面，如果pass说明解析对了
func TestConvertToStruct(t *testing.T) {
	kvMap := make(map[common.Hash]common.Hash)
	// voter
	kvMap[common.HexToHash("000000009ee97d274eb4c215f23238fee1f103d9ea10a2340000000000000000")] = common.HexToHash("0000000000000000000000949ee97d274eb4c215f23238fee1f103d9ea10a234") // Owner
	kvMap[common.HexToHash("000000009ee97d274eb4c215f23238fee1f103d9ea10a2340000000000000001")] = common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001") // IsProxy
	kvMap[common.HexToHash("000000009ee97d274eb4c215f23238fee1f103d9ea10a2340000000000000002")] = common.HexToHash("0000000000000000000000000000000000000000000000000000000000000066") // ProxyVoteCount
	kvMap[common.HexToHash("000000009ee97d274eb4c215f23238fee1f103d9ea10a2340000000000000003")] = common.HexToHash("0000000000000000000000940000000000000000000000000000000000000000") // Proxy
	kvMap[common.HexToHash("000000009ee97d274eb4c215f23238fee1f103d9ea10a2340000000000000004")] = common.HexToHash("0000000000000000000000000000000000000000000000000000000000000005") // LastStakeCount
	kvMap[common.HexToHash("000000009ee97d274eb4c215f23238fee1f103d9ea10a2340000000000000005")] = common.HexToHash("0000000000000000000000000000000000000000000000000000000000000005") // LastVoteCount
	kvMap[common.HexToHash("000000009ee97d274eb4c215f23238fee1f103d9ea10a2340000000000000006")] = common.HexToHash("000000000000000000000000000000000000000000000000000000845b4822c8") // TimeStamp
	kvMap[common.HexToHash("000000009ee97d274eb4c215f23238fee1f103d9ea10a2340000000100000007")] = common.HexToHash("0000000000000000000000949ee97d274eb4c215f23238fee1f103d9ea10a234") // VoteCandidates
	kvMap[common.HexToHash("000000009ee97d274eb4c215f23238fee1f103d9ea10a2340000000200000007")] = common.HexToHash("000000000000000000000094000000000000000000000000000000000000000a") // VoteCandidates
	kvMap[common.HexToHash("000000009ee97d274eb4c215f23238fee1f103d9ea10a2340000000000000007")] = common.HexToHash("0000000000000000000000000000000000000000000000000000000000000002") // VoteCandidates
	// candidate
	kvMap[common.HexToHash("010000009ee97d274eb4c215f23238fee1f103d9ea10a2340000000000000000")] = common.HexToHash("0000000000000000000000949ee97d274eb4c215f23238fee1f103d9ea10a234") // owner
	kvMap[common.HexToHash("010000009ee97d274eb4c215f23238fee1f103d9ea10a2340000000000000001")] = common.HexToHash("0000000000000000000000940000000000000000000000923839919383938289") // Binder
	kvMap[common.HexToHash("010000009ee97d274eb4c215f23238fee1f103d9ea10a2340000000000000002")] = common.HexToHash("0000000000000000000000940000000000000000000000923839919383938281") // Beneficiary
	kvMap[common.HexToHash("010000009ee97d274eb4c215f23238fee1f103d9ea10a2340000000000000003")] = common.HexToHash("0000000000000000000000000000000000000000000000000000000000000080") // vote count
	kvMap[common.HexToHash("010000009ee97d274eb4c215f23238fee1f103d9ea10a2340000000000000004")] = common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001") // register
	kvMap[common.HexToHash("010000009ee97d274eb4c215f23238fee1f103d9ea10a2340000000000000005")] = common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001") // bind
	kvMap[common.HexToHash("010000009ee97d274eb4c215f23238fee1f103d9ea10a2340000000000000006")] = common.HexToHash("0000000000000000000000000000b8502f6970342f3139322e3136382e392e31") // url
	kvMap[common.HexToHash("010000009ee97d274eb4c215f23238fee1f103d9ea10a2340000000100000006")] = common.HexToHash("30322f7463702f353231302f697066732f316b48614d556d5a6754706a474568") // url
	kvMap[common.HexToHash("010000009ee97d274eb4c215f23238fee1f103d9ea10a2340000000200000006")] = common.HexToHash("786347415472315556577936694b6b796746756b6e57457457374c694c726576") // url
	kvMap[common.HexToHash("010000009ee97d274eb4c215f23238fee1f103d9ea10a2340000000100000007")] = common.HexToHash("776562736974652e6e65742f746573742f7769746e6573732f77656273697465") // Website
	kvMap[common.HexToHash("010000009ee97d274eb4c215f23238fee1f103d9ea10a2340000000000000007")] = common.HexToHash("0000000000000000000000000000000000000000000000a87777772e74657374") // Website
	kvMap[common.HexToHash("010000009ee97d274eb4c215f23238fee1f103d9ea10a2340000000000000008")] = common.HexToHash("00000000000000000000000000000000000000000000000087746573744e6574") // name
	// stake
	kvMap[common.HexToHash("020000009ee97d274eb4c215f23238fee1f103d9ea10a2340000000000000000")] = common.HexToHash("0000000000000000000000949ee97d274eb4c215f23238fee1f103d9ea10a234") // Owner
	kvMap[common.HexToHash("020000009ee97d274eb4c215f23238fee1f103d9ea10a2340000000000000001")] = common.HexToHash("00000000000000000000000000000000000000000000000000000000000081e6") // StakeCount
	kvMap[common.HexToHash("020000009ee97d274eb4c215f23238fee1f103d9ea10a2340000000000000002")] = common.HexToHash("00000000000000000000000000000000000000000000000000000000000081e6") // Vnt
	kvMap[common.HexToHash("020000009ee97d274eb4c215f23238fee1f103d9ea10a2340000000000000003")] = common.HexToHash("000000000000000000000000000000000000000000000000000000845b4822c8") // TimeStamp
	// bounty
	kvMap[common.HexToHash("0300000000000000000000000000000000000000000000090000000000000000")] = common.HexToHash("0000000000000000000000000000000000000000000000880de0b6b3a7640000") // RestTotalBounty

	getFn := func(hash common.Hash) common.Hash {
		return kvMap[hash]
	}

	voter1 := Voter{}
	err := convertToStruct(VOTERPREFIX, common.HexToAddress("9ee97d274eb4c215f23238fee1f103d9ea10a234"), &voter1, getFn)
	if err != nil {
		t.Errorf(err.Error())
	}
	if same, err := sameVoter(&voter, &voter1); !same {
		t.Errorf(err.Error())
	}

	var candidate1 Candidate
	err = convertToStruct(CANDIDATEPREFIX, common.HexToAddress("9ee97d274eb4c215f23238fee1f103d9ea10a234"), &candidate1, getFn)
	if err != nil {
		t.Error(err)
	}
	if same, err := sameCandidate(&candidate, &candidate1); !same {
		t.Errorf(err.Error())
	}

	var stake1 Stake
	err = convertToStruct(STAKEPREFIX, common.HexToAddress("9ee97d274eb4c215f23238fee1f103d9ea10a234"), &stake1, getFn)
	if err != nil {
		t.Error(err)
	}
	if same, err := sameStake(&stake, &stake1); !same {
		t.Errorf(err.Error())
	}

	var bounty1 Reward
	err = convertToStruct(REWARDPREFIX, contractAddr, &bounty1, getFn)
	if err != nil {
		t.Error(err)
	}
	if bounty1.Rest.Cmp(bounty.Rest) != 0 {
		t.Errorf("Error: the reset total Reward before is %v after is %v", bounty.Rest, bounty1.Rest)
	}
}

func TestSetToDB(t *testing.T) {
	db := vntdb.NewMemDatabase()
	stateDB, _ := state.New(common.Hash{}, state.NewDatabase(db))

	ctx := testContext{StateDB: stateDB}
	c := newElectionContext(&ctx)

	err := c.setVoter(voter)
	if err != nil {
		t.Error(err)
	}

	err = c.setCandidate(candidate)
	if err != nil {
		t.Error(err)
	}

	err = c.setStake(stake)
	if err != nil {
		t.Error(err)
	}

	err = setReward(stateDB, bounty)
	if err != nil {
		t.Error(err)
	}

}

func TestGetFromDB(t *testing.T) {
	db := vntdb.NewMemDatabase()
	stateDB, _ := state.New(common.Hash{}, state.NewDatabase(db))

	ctx := testContext{StateDB: stateDB}
	c := newElectionContext(&ctx)

	err1 := c.setVoter(voter)
	err2 := c.setCandidate(candidate)
	err3 := c.setStake(stake)
	err4 := setReward(stateDB, bounty)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		t.Fatal("SetToDB err", err1, err2, err3, err4)
	}

	voter1 := c.getVoter(voter.Owner)

	if same, err := sameVoter(&voter, &voter1); !same {
		t.Errorf(err.Error())
	}

	candidate1 := c.getCandidate(candidate.Owner)
	if same, err := sameCandidate(&candidate, &candidate1); !same {
		t.Errorf(err.Error())
	}

	stake1 := c.getStake(stake.Owner)
	if same, err := sameStake(&stake, &stake1); !same {
		t.Errorf(err.Error())
	}

	bounty1 := getReward(stateDB)
	if bounty.Rest.Cmp(bounty1.Rest) != 0 {
		t.Errorf("Error: the reset total Reward before is %v after is %v", bounty.Rest, bounty1.Rest)
	}
}

func TestGetAllCandidate(t *testing.T) {
	db := vntdb.NewMemDatabase()
	stateDB, _ := state.New(common.Hash{}, state.NewDatabase(db))

	ctx := testContext{StateDB: stateDB}
	c := newElectionContext(&ctx)

	for i := 0; i < 255; i++ {
		candidate1 := candidate
		candidate1.Owner[0] = byte(i)
		if err := c.setCandidate(candidate1); err != nil {
			t.Errorf("candiates: %s, error: %s", candidate1.Owner, err)
		}
	}

	candidates := getAllCandidate(stateDB)
	if len(candidates) != 255 {
		t.Error("the number of candidates is wrong!")
	} else {
		for _, candidate := range candidates {
			if !candidate.Active() || candidate.VoteCount.Cmp(big.NewInt(0)) != 0 {
				t.Fatalf("Error: %s", candidate.String())
			}
		}
	}
}

func TestGetFirstXCandidates_1(t *testing.T) {
	db := vntdb.NewMemDatabase()
	stateDB, _ := state.New(common.Hash{}, state.NewDatabase(db))

	ctx := testContext{StateDB: stateDB}
	c := newElectionContext(&ctx)

	type addrPair struct {
		addrPre byte
		votes   int64
	}

	witNum := 4
	tests := []addrPair{
		{byte(1), 200},
		{byte(2), 100},
		{byte(3), 50},
		{byte(4), 10},
		{byte(5), 100},
		{byte(6), 5},
	}
	rets := []addrPair{
		{byte(1), 200},
		{byte(2), 100},
		{byte(5), 100},
		{byte(3), 50},
		// {byte(5), 10},
		// {byte(6), 5},
	}

	// set to db
	for i := 0; i < len(tests); i++ {
		candidate1 := candidate
		candidate1.Owner[0] = byte(tests[i].addrPre)
		candidate1.VoteCount = big.NewInt(tests[i].votes)
		if err := c.setCandidate(candidate1); err != nil {
			t.Errorf("candiates: %s, error: %s", candidate1.Owner, err)
		}
	}

	witsAddr, _ := GetFirstNCandidates(stateDB, witNum)
	if len(witsAddr) != len(rets) {
		t.Errorf("lenght not match, want:%d, got:%d", witNum, len(witsAddr))
	}
	baseAddr := candidate.Owner
	for i := 0; i < witNum; i++ {
		can := baseAddr
		can[0] = byte(rets[i].addrPre)
		ret := bytes.Compare(can.Bytes(), witsAddr[i].Bytes())
		if ret != 0 {
			t.Errorf("candidates nots match at index:%d, ret:%d, want: %x, got:%x", i, ret, can, witsAddr[i])
		}
	}

	// candidates := getAllCandidate(stateDB)
	// for _, candi := range candidates {
	// 	fmt.Printf("candidate owner: %x, active: %v, voteCount : %v\n", candi.Owner, candi.Active, candi.VoteCount)
	// }

}

func TestGetFirstXCandidates_2(t *testing.T) {
	db := vntdb.NewMemDatabase()
	stateDB, _ := state.New(common.Hash{}, state.NewDatabase(db))

	ctx := testContext{StateDB: stateDB}
	c := newElectionContext(&ctx)

	type addrPair struct {
		addrPre byte
		votes   int64
	}

	witNum := 4
	tests := []addrPair{
		{byte(1), 200},
		{byte(2), 100},
		{byte(3), 50},
		{byte(4), 100},
		{byte(5), 10},
		{byte(6), 5},
	}
	rets := []addrPair{
		{byte(1), 200},
		{byte(2), 100},
		{byte(4), 100},
		{byte(3), 50},
		// {byte(5), 10},
		// {byte(6), 5},
	}

	// set to db
	for i := 0; i < len(tests); i++ {
		candidate1 := candidate
		candidate1.Owner[0] = byte(tests[i].addrPre)
		candidate1.VoteCount = big.NewInt(tests[i].votes)
		if err := c.setCandidate(candidate1); err != nil {
			t.Errorf("candiates: %s, error: %s", candidate1.Owner, err)
		}
	}

	witsAddr, _ := GetFirstNCandidates(stateDB, witNum)
	if len(witsAddr) != len(rets) {
		t.Errorf("lenght not match, want:%d, got:%d", witNum, len(witsAddr))
	}
	baseAddr := candidate.Owner
	for i := 0; i < witNum; i++ {
		can := baseAddr
		can[0] = byte(rets[i].addrPre)
		ret := bytes.Compare(can.Bytes(), witsAddr[i].Bytes())
		if ret != 0 {
			t.Errorf("candidates nots match at index:%d, ret:%d, want: %x, got:%x", i, ret, can, witsAddr[i])
		}
	}

	candidates := getAllCandidate(stateDB)
	for _, candi := range candidates {
		t.Logf("candidate owner: %x, active: %v, voteCount : %v\n", candi.Owner, candi.Active(), candi.VoteCount)
	}
}

// 存在不active的节点
func TestGetFirstXCandidates_3(t *testing.T) {
	db := vntdb.NewMemDatabase()
	stateDB, _ := state.New(common.Hash{}, state.NewDatabase(db))

	ctx := testContext{StateDB: stateDB}
	c := newElectionContext(&ctx)

	type addrPair struct {
		addrPre byte
		votes   int64
		bind    bool // register始终是true，bind是true则active为true
	}

	witNum := 4
	tests := []addrPair{
		{byte(1), 200, true},
		{byte(2), 100, false},
		{byte(3), 50, true},
		{byte(4), 100, true},
		{byte(5), 10, true},
		{byte(6), 5, true},
	}
	rets := []addrPair{
		{byte(1), 200, true},
		{byte(4), 100, true},
		{byte(3), 50, true},
		{byte(5), 10, true},
	}

	// set to db
	for i := 0; i < len(tests); i++ {
		candidate1 := candidate
		candidate1.Owner[0] = byte(tests[i].addrPre)
		candidate1.VoteCount = big.NewInt(tests[i].votes)
		candidate1.Registered = true
		candidate1.Bind = tests[i].bind
		if err := c.setCandidate(candidate1); err != nil {
			t.Errorf("candiates: %s, error: %s", candidate1.Owner, err)
		}
	}

	witsAddr, _ := GetFirstNCandidates(stateDB, witNum)
	if len(witsAddr) != len(rets) {
		t.Errorf("lenght not match, want:%d, got:%d", witNum, len(witsAddr))
	}
	baseAddr := candidate.Owner
	for i := 0; i < witNum; i++ {
		can := baseAddr
		can[0] = byte(rets[i].addrPre)
		ret := bytes.Compare(can.Bytes(), witsAddr[i].Bytes())
		if ret != 0 {
			t.Errorf("candidates nots match at index:%d, ret:%d, want: %x, got:%x", i, ret, can, witsAddr[i])
		}
	}

	candidates := getAllCandidate(stateDB)
	for _, candi := range candidates {
		t.Logf("candidate owner: %x, active: %v, voteCount : %v\n", candi.Owner, candi.Active(), candi.VoteCount)
	}
}

// 使用registerWitness注册见证人，每个人应当0票，但按地址排序
func TestGetFirstXCandidates_4(t *testing.T) {
	db := vntdb.NewMemDatabase()
	stateDB, _ := state.New(common.Hash{}, state.NewDatabase(db))

	ctx := testContext{StateDB: stateDB}
	c := newElectionContext(&ctx)

	type addrPair struct {
		addrPre byte
		votes   int64
	}

	// 忽略票数
	witNum := 4
	tests := []addrPair{
		{byte(1), 0},
		{byte(2), 0},
		{byte(3), 0},
		{byte(4), 0},
		{byte(5), 0},
		{byte(6), 0},
	}
	rets := []addrPair{
		{byte(1), 0},
		{byte(2), 0},
		{byte(3), 0},
		{byte(4), 0},
	}

	// 设置到数据库
	baseAddr := candidate.Owner
	for i := 0; i < len(tests); i++ {
		candidate1 := candidate
		candidate1.Owner[0] = byte(tests[i].addrPre)
		candidate1.VoteCount = big.NewInt(tests[i].votes)
		candidate1.Registered = true
		candidate1.Bind = true
		if err := c.setCandidate(candidate1); err != nil {
			t.Errorf("candiates: %s, error: %s", candidate1.Owner, err)
		}
	}

	witsAddr, _ := GetFirstNCandidates(stateDB, witNum)
	if len(witsAddr) != len(rets) {
		t.Errorf("lenght not match, want:%d, got:%d", witNum, len(witsAddr))
		t.FailNow()
	}
	for i := 0; i < witNum; i++ {
		can := baseAddr
		can[0] = byte(rets[i].addrPre)
		ret := bytes.Compare(can.Bytes(), witsAddr[i].Bytes())
		if ret != 0 {
			t.Errorf("candidates nots match at index:%d, ret:%d, want: %x, got:%x", i, ret, can, witsAddr[i])
		}
	}
}
