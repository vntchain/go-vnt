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
	"reflect"
	"strconv"
	"testing"

	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core/state"
	inter "github.com/vntchain/go-vnt/core/vm/interface"
	"github.com/vntchain/go-vnt/vntdb"
)

var url = []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHGq5zZFRW5FBJ9YMbbvSiW4AzGg5CKMCtDeg6FNnjCbGS")

var InputCase = [][]byte{
	common.FromHex("c94ba774"),
	common.FromHex("68cc738800000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000a"),
	common.FromHex("95c96554"),
	common.FromHex("487a2abb"),
	common.FromHex("31f080ef"),
	common.FromHex("c9a63035"),
	common.FromHex("97107d6d000000000000000000000000a863d8efa01ece6fabfa7e8c85217a3c1af833a9"),
	common.FromHex("a694fc3a0000000000000000000000000000000000000000000000000000000000000064"),
	common.FromHex("73cf575a"),
	common.FromHex("65f7314e000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000e00000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000000000000000000000000000004d2f6970342f3132372e302e302e312f7463702f33303330332f697066732f316b484771357a5a4652573546424a39594d62627653695734417a476735434b4d437444656736464e6e6a436247530000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000317777772e746573746e65742e696e666f2e776562736974652e746573742e746573742e746573742e746573742e74657374000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000874657374696e666f000000000000000000000000000000000000000000000000"),
	common.FromHex("f67ab93e"),
}

type candiRegInfo struct {
	addr    common.Address
	name    []byte
	website []byte
	url     []byte
}

var (
	addr1 = common.HexToAddress("41b0db166cfdf1c4ba3ce657171482a9aa55cc93")
	addr2 = common.HexToAddress("08b467a881ec34b668254aa956e0c46f9c3b2b83")
	addr3 = common.HexToAddress("0c0292587ccdc76b8f449002a017bc9479ff0a88")
	addr4 = common.HexToAddress("0a0292587ccdc76b8f449002a017bc9479ff0a88")
	addr5 = common.HexToAddress("0b0292587ccdc76b8f449002a017bc9479ff0a88")
	addr6 = common.HexToAddress("0d0292587ccdc76b8f449002a017bc9479ff0a88")
	addr7 = common.HexToAddress("0e0292587ccdc76b8f449002a017bc9479ff0a88")
	addr8 = common.HexToAddress("0e0292587ccdc76b8f449002a017bc9479ff0a81")
	addr9 = common.HexToAddress("0e0292587ccdc76b8f449002a017bc9479ff0a82")

	candiInfo1 = candiRegInfo{addr1, []byte("node1"), []byte("www.node1.com"), []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHNAAfnqXNsxMwJf6QjJFRmVK7iB32U9owwK9KfeLFxEA7")}
	candiInfo2 = candiRegInfo{addr2, []byte("node2"), []byte("www.node2.com"), []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHcch6yuBCgC5nPPSK3Yp7Es4c4eenxAeK167pYwUvNjRo")}
	candiInfo3 = candiRegInfo{addr3, []byte("node3"), []byte("www.node3.com"), []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHJFKr2bzUnMr1NbeyYbYJa3RXT18cEu7cNDrHWjg8XYKB")}
	candiInfo4 = candiRegInfo{addr4, []byte("node4"), []byte("www.node4.com"), []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHfop9dnUHHmtBXVkLB5UauAmACtrsEX5H5t6oCRpdL198")}
	candiInfo5 = candiRegInfo{addr5, []byte("node5"), []byte("www.node5.com"), []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHGG8L1DTrVG3Cad479Q32oGmFAiEjLFwxzNyXH3ehGo73")}
	candiInfo6 = candiRegInfo{addr6, []byte("node6"), []byte("www.node6.com"), []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHmExX4yutwBZLbRsYHq59KfgiM1LUJFW2JSPeSCcBf7rH")}
	candiInfo7 = candiRegInfo{addr7, []byte("nodd7"), []byte("www.node7.com"), []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHSMhv82q5thJkdeJzxCVW8tdXwaDThBZWsH2Q9KUGGFUq")}
	candiInfo8 = candiRegInfo{addr8, []byte("nodd8"), []byte("www.node8.com"), []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHYwV52itn31V5fRXzERMygHFDx6PrSFS8puEr3N4Ujv69")}
	candiInfo9 = candiRegInfo{addr9, []byte("nodd9"), []byte("www.node9.com"), []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHiEQS9qnK1YJN25eNyTjDRUukTzComAMWYAowTAYCu1K4")}

	candidates = []common.Address{addr1, addr2, addr3, addr4, addr5, addr6, addr7, addr8, addr9}
	candiInfos = []candiRegInfo{candiInfo1, candiInfo2, candiInfo3, candiInfo4, candiInfo5, candiInfo6, candiInfo7, candiInfo8, candiInfo9}
)

type testContext struct {
	Origin  common.Address
	Time    *big.Int
	StateDB inter.StateDB
}

func (tc *testContext) GetOrigin() common.Address {
	return tc.Origin
}

func (tc *testContext) GetStateDb() inter.StateDB {
	return tc.StateDB
}

func (tc *testContext) GetTime() *big.Int {
	return tc.Time
}

func (tc *testContext) SetTime(t *big.Int) {
	tc.Time = t
}

func newcontext() inter.ChainContext {
	db := vntdb.NewMemDatabase()
	stateDB, _ := state.New(common.Hash{}, state.NewDatabase(db))
	c := testContext{
		Origin:  common.BytesToAddress([]byte{111}),
		Time:    big.NewInt(1531328510),
		StateDB: stateDB,
	}
	return &c
}

func getAllVoter(t *testing.T, db inter.StateDB) []*Voter {
	var result []*Voter
	voters := make(map[common.Hash]common.Hash)
	addrs := make(map[common.Address]struct{})

	// 根据voter前缀，在db中获取所有的voter地址
	db.ForEachStorage(contractAddr, func(key common.Hash, value common.Hash) bool {
		if key[0] == VOTERPREFIX {
			voters[key] = value

			var addr common.Address
			copy(addr[:], key[PREFIXLENGTH:PREFIXLENGTH+common.AddressLength])
			addrs[addr] = struct{}{}
		}
		return true
	})

	getFn := func(key common.Hash) common.Hash {
		return voters[key]
	}

	// 根据已获得的voter地址获取所有的voter结构体
	for addr := range addrs {
		var voter Voter
		if err := convertToStruct(VOTERPREFIX, addr, &voter, getFn); err != nil {
			t.Errorf("addr: %s, error: %s", addr.String(), err)
		}
		result = append(result, &voter)

	}
	return result
}

// 计算所有投出的票数，跟候选人所拥有的票数是否想相等
func checkValid(t *testing.T, c electionContext) (bool, error) {
	// 保存原先context的时间
	currentTime := c.context.GetTime()

	proxyVote := make(map[common.Address]*big.Int)
	voteCount := make(map[common.Address]*big.Int)
	voters := getAllVoter(t, c.context.GetStateDb())
	// 循环一遍voter，做一遍初步的检查
	for _, voter := range voters {
		// 如果proxy不为空
		if !bytes.Equal(voter.Proxy.Bytes(), emptyAddress.Bytes()) {
			// 则isProxy为假，且voteCandidates为空
			if voter.IsProxy || len(voter.VoteCandidates) != 0 {
				return false, fmt.Errorf("voter owner: %x, proxy: %x, isProxy: %t, voteCandidates: %v\n", voter.Owner, voter.Proxy, voter.IsProxy, voter.VoteCandidates)
			}
			// 统计投给代理的票
			if count, ok := proxyVote[voter.Proxy]; ok {
				count.Add(count, voter.LastVoteCount)
				count.Add(count, voter.ProxyVoteCount)
			}
		}

		if voter.LastVoteCount != nil && voter.LastVoteCount.Sign() > 0 {
			if voter.TimeStamp == nil || voter.TimeStamp.Cmp(big.NewInt(0)) == 0 {
				return false, fmt.Errorf("lastVoteCount is not zero. timeStamp must not be nil")

			}
			// 比对lastVoteCount，是否是当时的抵押数兑换所得
			stake := c.getStake(voter.Owner)
			if stake.TimeStamp == nil || stake.TimeStamp.Cmp(big.NewInt(0)) == 0 {
				return false, fmt.Errorf("lastVoteCount is not zero. stake.timeStamp must not be nil")

			}
			// 抵押数大于0，且抵押时间小于投票时间，说明上次投票后抵押数没有变
			if stake.StakeCount.Sign() > 0 && stake.TimeStamp.Cmp(voter.TimeStamp) <= 0 {
				if ctx, ok := c.context.(*testContext); ok {
					ctx.SetTime(voter.TimeStamp)
				}
				// 计算抵押数可以兑换的票数，与上次投票所得是否一致
				calculateCount := c.calculateVoteCount(stake.StakeCount)
				if voter.LastVoteCount.Cmp(calculateCount) != 0 {
					if ctx, ok := c.context.(*testContext); ok {
						ctx.SetTime(currentTime)
					}
					return false, fmt.Errorf("time: %v, lastVoteCount : %d, stakeCount: %d,calculateCount : %d ",
						voter.TimeStamp, voter.LastVoteCount, stake.StakeCount, calculateCount)
				}
			}
			// 统计总投票
			for _, candi := range voter.VoteCandidates {
				if voteCount[candi] == nil {
					voteCount[candi] = big.NewInt(0)
				}
				voteCount[candi].Add(voteCount[candi], voter.LastVoteCount)
				voteCount[candi].Add(voteCount[candi], voter.ProxyVoteCount)
			}
		}

	}
	if ctx, ok := c.context.(*testContext); ok {
		ctx.SetTime(currentTime)
	}

	// 检查代理投票数
	for _, voter := range voters {
		if _, ok := proxyVote[voter.Owner]; !ok {
			continue
		}
		if voter.ProxyVoteCount.Cmp(proxyVote[voter.Owner]) != 0 {
			return false, fmt.Errorf("proxyVoteCount is wrong, proxyVoteCount in db: %d, expect proxyVote: %d\n", voter.ProxyVoteCount, proxyVote[voter.Owner])

		}
	}

	candidates := getAllCandidate(c.context.GetStateDb())
	for _, candidate := range candidates {
		if voteCount[candidate.Owner] == nil {
			voteCount[candidate.Owner] = big.NewInt(0)
		}
		if candidate.VoteCount == nil || candidate.VoteCount.Cmp(voteCount[candidate.Owner]) != 0 {
			return false, fmt.Errorf("voteCount is wrong. candidate address: %x,voteCount in db: %d, expect voteCount : %d", candidate.Owner, candidate.VoteCount, voteCount[candidate.Owner])
		}
	}

	return true, nil
}

func TestInput(t *testing.T) {
	var e Election
	context := newcontext()

	for _, input := range InputCase {
		_, err := e.Run(context, input)
		if err != nil && err == fmt.Errorf("call election contract err: method doesn't exist") {
			t.Error(err)
		}
	}
}

func TestCandidate_votes(t *testing.T) {
	var addr1 common.Address
	c1 := &Candidate{
		Owner:     addr1,
		VoteCount: big.NewInt(10),
		Active:    true,
	}

	if c1.votes().Cmp(big.NewInt(10)) != 0 {
		t.Errorf("votes() error. want = %v, got = %s", 10, c1.votes().String())
	}

	c1.Active = false
	if c1.votes().Cmp(big.NewInt(-10)) != 0 {
		t.Errorf("votes() error. want = %v, got = %s", -10, c1.votes().String())
	}
}

func TestCandidate_equal(t *testing.T) {
	addr1 := common.HexToAddress("0x122369f04f32269598789998de33e3d56e2c507a")
	addr2 := common.HexToAddress("0x42a875ac43f2b4e6d17f54d288071f5952bf8911")
	c1 := Candidate{Owner: addr1, VoteCount: big.NewInt(10), Active: true}
	c2 := Candidate{Owner: addr2, VoteCount: big.NewInt(20), Active: false}

	if c1.equal(&c2) {
		t.Errorf("two Candidate should not equal")
	}

	c1.Owner = addr2
	c1.Active = false
	c1.VoteCount = big.NewInt(20)

	if c1.equal(&c2) == false {
		t.Errorf("two Candidate should equal")
	}
}

func TestCandidateList_Less(t *testing.T) {
	addr1 := common.HexToAddress("0x522369f04f32269598789998de33e3d56e2c507a")
	addr2 := common.HexToAddress("0x42a875ac43f2b4e6d17f54d288071f5952bf8911")
	addr3 := common.HexToAddress("0x18a875ac43f2b4e6d17f54d288071f5952bf8911")
	c1 := Candidate{Owner: addr1, VoteCount: big.NewInt(10), Active: true}
	c2 := Candidate{Owner: addr2, VoteCount: big.NewInt(20), Active: false}
	c3 := Candidate{Owner: addr3, VoteCount: big.NewInt(10), Active: true}

	cl := CandidateList{c1, c2, c3}

	if cl.Less(0, 1) == false {
		t.Errorf("c1 should greater than c2, c1= %s, c2=%s", c1.String(), c2.String())
	}
	if cl.Less(0, 2) == true {
		t.Errorf("c1 should less than c3, c1= %s, c3=%s", c1.String(), c3.String())
	}
}

func TestCandidateList_Swap(t *testing.T) {
	c1 := Candidate{common.HexToAddress("0x1"), big.NewInt(100),
		false, []byte("/p2p/1"), big.NewInt(10000), big.NewInt(200),
		big.NewInt(1548664636), []byte("node1.com"), []byte("node1")}
	c2 := Candidate{common.HexToAddress("0x2"), big.NewInt(20),
		true, []byte("/p2p/2"), big.NewInt(10000), big.NewInt(200),
		big.NewInt(1548664636), []byte("node2.com"), []byte("node2")}

	candidates := CandidateList{c1, c2}
	swaped := CandidateList{c2, c1}
	candidates.Swap(0, 1)
	for i, tt := range candidates {
		if tt.equal(&swaped[i]) == false {
			t.Errorf("index: %d, expect: %s, got: %s", i, swaped[i].String(), tt.String())
		}
	}
}

func TestCandidateList_Sort(t *testing.T) {
	// c1票数为负，c2与c5票数相等，c3票数最多
	c1 := Candidate{common.HexToAddress("0x1"), big.NewInt(100),
		false, []byte("/p2p/1"), big.NewInt(10000), big.NewInt(200),
		big.NewInt(1548664636), []byte("node1.com"), []byte("node1")}
	c2 := Candidate{common.HexToAddress("0x2"), big.NewInt(20),
		true, []byte("/p2p/2"), big.NewInt(10000), big.NewInt(200),
		big.NewInt(1548664636), []byte("node2.com"), []byte("node2")}
	c3 := Candidate{common.HexToAddress("0x3"), big.NewInt(90),
		true, []byte("/p2p/3"), big.NewInt(10000), big.NewInt(200),
		big.NewInt(1548664636), []byte("node3.com"), []byte("node3")}
	c4 := Candidate{common.HexToAddress("0x4"), big.NewInt(40),
		true, []byte("/p2p/4"), big.NewInt(10000), big.NewInt(200),
		big.NewInt(1548664636), []byte("node4.com"), []byte("node4")}
	c5 := Candidate{common.HexToAddress("0x5"), big.NewInt(20),
		true, []byte("/p2p/5"), big.NewInt(10000), big.NewInt(200),
		big.NewInt(1548664636), []byte("node5.com"), []byte("node5")}
	candidates := CandidateList{c1, c2, c3, c4, c5}
	sorted := CandidateList{c3, c4, c2, c5, c1}

	candidates.Sort()
	for i, tt := range candidates {
		if tt.equal(&sorted[i]) == false {
			t.Errorf("index: %d, expect: %s, got: %s", i, sorted[i].String(), tt.String())
		}
	}
}

// Test voteWitnesses
// 投的候选人过多返回错误
func TestVoteTooManyCandidates(t *testing.T) {
	context := newcontext()
	c := newElectionContext(context)
	addr := common.BytesToAddress([]byte{111})

	var candidates []common.Address
	for i := 1; i <= VoteLimit+1; i++ {
		candidate := common.BytesToAddress([]byte{byte(i)})
		candidates = append(candidates, candidate)
		website := "www.testnet.info" + strconv.Itoa(i)
		name := "testinfo" + strconv.Itoa(i)
		p2pUrl := []byte(string(url)[:13] + strconv.Itoa(i) + string(url)[14:])
		// t.Logf("url: %s", string(p2pUrl))
		if err := c.registerWitness(candidate, p2pUrl, []byte(website), []byte(name)); err != nil {
			t.Errorf("register failed, addr: %s, error: %s", candidate.String(), err)
		}
	}
	err := c.voteWitnesses(addr, candidates)
	if err.Error() != fmt.Sprintf("you voted too many candidates: the limit is %d, you voted %d", VoteLimit, len(candidates)) {
		t.Error(err)
	}
}

// Test voteWitnesses
// 当前无抵押返回错误
func TestVoteWithoutStake(t *testing.T) {
	context := newcontext()
	c := newElectionContext(context)

	addr := common.BytesToAddress([]byte{111})

	// 不抵押投票，返回未抵押的错误
	err := c.voteWitnesses(addr, candidates)
	if err.Error() != "you must stake before vote" {
		t.Error(err)
	}
}

// Test voteWitnesses
// 距数据库中上次操作时间不足24小时，返回错误
func TestVoteWithoutEnoughTimeGap(t *testing.T) {
	context := newcontext()
	c := newElectionContext(context)
	addr := common.BytesToAddress([]byte{111})

	// 数据库中塞一条voter数据
	voter := Voter{
		Owner:     addr,
		TimeStamp: big.NewInt(1531328500),
	}
	if err := c.setVoter(voter); err != nil {
		t.Errorf("addr: %s, error: %s", voter.Owner, err)
	}

	// 抵押
	c.context.GetStateDb().AddBalance(addr, big.NewInt(1e18))
	if err := c.stake(addr, big.NewInt(1)); err != nil {
		t.Errorf("stake failed, addr: %s, error: %s", addr.String(), err)
	}
	// 投票
	err := c.voteWitnesses(addr, candidates)
	if err.Error() != fmt.Sprintf("it's less than 24h after your last vote or setProxy, lastTime: %v, now: %v", voter.TimeStamp, c.context.GetTime()) {
		t.Error(err)
	}
}

// Test voteWitnesses
// 原先没有投票也没有代理的情况
func TestVoteCandidatesFistTime(t *testing.T) {
	context := newcontext()
	c := newElectionContext(context)

	addr := common.BytesToAddress([]byte{111})
	c.context.GetStateDb().AddBalance(addr, big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)))
	if err := c.stake(addr, big.NewInt(10)); err != nil {
		t.Errorf("stake failed, addr: %s, error: %s", addr.String(), err)
	}

	// 候选人注册
	for i := 0; i < len(candidates); i++ {
		website := "www.testnet.info" + strconv.Itoa(i)
		name := "testinfo" + strconv.Itoa(i)
		p2pUrl := []byte(string(url)[:13] + strconv.Itoa(i) + string(url)[14:])
		if err := c.registerWitness(candidates[i], p2pUrl, []byte(website), []byte(name)); err != nil {
			t.Errorf("register failed, addr: %s, error: %s", candidates[i].String(), err)
		}
	}

	// 投票
	err := c.voteWitnesses(addr, candidates)
	if err != nil {
		t.Error(err)
	}

	if _, err := checkValid(t, c); err != nil {
		t.Error(err)
	}
}

// Test cancelVote
// 数据库中无记录
func TestCancelVoteNoRecord(t *testing.T) {
	context := newcontext()
	c := newElectionContext(context)

	addr := common.BytesToAddress([]byte{111})
	err := c.cancelVote(addr)
	if err.Error() != fmt.Sprintf("the voter %x doesn't exist", addr) {
		t.Error(err)
	}
}

// Test cancelVote
// 数据库中无记录
func TestCancelVoteWithProxy(t *testing.T) {
	context := newcontext()
	c := newElectionContext(context)

	addr := common.BytesToAddress([]byte{111})
	// 数据库中塞条记录
	voter := Voter{
		Owner:     addr,
		Proxy:     common.BytesToAddress([]byte{10}),
		TimeStamp: big.NewInt(1531328510),
	}
	if err := c.setVoter(voter); err != nil {
		t.Errorf("addr: %s, error: %s", voter.Owner, err)
	}

	err := c.cancelVote(addr)
	if err.Error() != fmt.Sprintf("must cancel proxy first, proxy: %x", voter.Proxy) {
		t.Error(err)
	}
}

// Test cancelVote
func TestCancelVote(t *testing.T) {
	context := newcontext()
	c := newElectionContext(context)
	addr := common.BytesToAddress([]byte{111})
	addr1 := common.BytesToAddress([]byte{10})

	// 抵押1
	c.context.GetStateDb().AddBalance(addr, big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)))
	if err := c.stake(addr, big.NewInt(10)); err != nil {
		t.Errorf("stake failed, addr: %s, error: %s", addr.String(), err)
	}

	// 抵押2
	c.context.GetStateDb().AddBalance(addr1, big.NewInt(0).Mul(big.NewInt(100), big.NewInt(1e18)))
	if err := c.stake(addr1, big.NewInt(100)); err != nil {
		t.Errorf("stake failed, addr: %s, error: %s", addr.String(), err)
	}

	// 设置候选人
	for i := 0; i < len(candidates); i++ {
		website := "www.testnet.info" + strconv.Itoa(i)
		name := "testinfo" + strconv.Itoa(i)
		p2pUrl := []byte(string(url)[:13] + strconv.Itoa(i) + string(url)[14:])
		if err := c.registerWitness(candidates[i], p2pUrl, []byte(website), []byte(name)); err != nil {
			t.Errorf("register failed, addr: %s, error: %s", candidates[i].String(), err)
		}
	}

	// 投票1
	if err := c.voteWitnesses(addr, candidates); err != nil {
		t.Errorf("vote addr: %s, error: %s", addr.String(), err)
	}
	voteCount := c.calculateVoteCount(big.NewInt(10))
	if _, err := checkValid(t, c); err != nil {
		t.Error(err)
	}

	// 投票2
	if err := c.voteWitnesses(addr1, candidates); err != nil {
		t.Errorf("vote addr: %s, error: %s", addr1.String(), err)
	}
	voteCount1 := c.calculateVoteCount(big.NewInt(100))
	totalVoteCount := new(big.Int).Set(voteCount)
	totalVoteCount.Add(totalVoteCount, voteCount1)

	// 验证2
	if _, err := checkValid(t, c); err != nil {
		t.Error(err)
	}

	// 取消投票
	if err := c.cancelVote(addr); err != nil {
		t.Errorf("cancelVote, addr: %s, error: %s", addr.String(), err)
	}

	if _, err := checkValid(t, c); err != nil {
		t.Error(err)
	}
}

// Test setProxy
// 设置代理人为自身返回错误
func TestProxySelf(t *testing.T) {
	context := newcontext()
	c := newElectionContext(context)
	addr := common.BytesToAddress([]byte{111})

	err := c.setProxy(addr, addr)
	if err.Error() != "cannot proxy to self" {
		t.Error(err)
	}
}

// Test setProxy
// 设置代理人为自身
func TestProxySelfIsProxy(t *testing.T) {
	context := newcontext()
	c := newElectionContext(context)
	addr := common.BytesToAddress([]byte{111})
	proxy := common.BytesToAddress([]byte{10})
	// addr 成为代理
	if err := c.startProxy(addr); err != nil {
		t.Errorf("start proxy, addr: %s, error: %s", addr.String(), err)
	}
	if err := c.setProxy(addr, proxy); err.Error() != "account registered as a proxy is not allowed to use a proxy" {
		t.Error(err)
	}
}

// Test setProxy
// 当前无抵押返回错误
func TestProxyWithoutStake(t *testing.T) {
	context := newcontext()
	c := newElectionContext(context)

	addr := common.BytesToAddress([]byte{111})
	proxy := common.BytesToAddress([]byte{10})
	// 不抵押投票，返回未抵押的错误
	err := c.setProxy(addr, proxy)
	if err.Error() != "you must stake before vote" {
		t.Error(err)
	}
}

// Test setProxy
// 距数据库中上次操作时间不足24小时，返回错误
func TestProxyWithoutEnoughTimeGap(t *testing.T) {
	context := newcontext()
	c := newElectionContext(context)
	addr := common.BytesToAddress([]byte{111})
	proxy := common.BytesToAddress([]byte{10})

	// 数据库中塞一条voter数据
	voter := Voter{
		Owner:     addr,
		TimeStamp: big.NewInt(1531328500),
	}
	if err := c.setVoter(voter); err != nil {
		t.Errorf("addr: %s, error: %s", voter.Owner, err)
	}

	// 抵押
	c.context.GetStateDb().AddBalance(addr, big.NewInt(1e18))
	if err := c.stake(addr, big.NewInt(1)); err != nil {
		t.Errorf("stake failed, addr: %s, error: %s", addr.String(), err)
	}
	// 设置代理
	err := c.setProxy(addr, proxy)
	if err.Error() != fmt.Sprintf("it's less than 24h after your last vote or setProxy, lastTime: %v, now: %v", voter.TimeStamp, c.context.GetTime()) {
		t.Error(err)
	}
}

// Test setProxy
// 设置的代理人不是代理
func TestProxyIsNotProxy(t *testing.T) {
	context := newcontext()
	c := newElectionContext(context)
	addr := common.BytesToAddress([]byte{111})
	proxy := common.BytesToAddress([]byte{10})

	// 抵押
	c.context.GetStateDb().AddBalance(addr, big.NewInt(1e18))
	if err := c.stake(addr, big.NewInt(1)); err != nil {
		t.Errorf("stake failed, addr: %s, error: %s", addr.String(), err)
	}

	// 设置代理
	err := c.setProxy(addr, proxy)
	if err.Error() != fmt.Sprintf("%x is not a proxy", proxy) {
		t.Error(err)
	}
}

func TestSetProxy(t *testing.T) {
	context := newcontext()
	c := newElectionContext(context)

	err := setProxy(t, c)
	if err != nil {
		t.Error(err)
	}
}

// Test cancelProxy
// 数据库中无记录
func TestCancelProxyNoRecord(t *testing.T) {
	context := newcontext()
	c := newElectionContext(context)
	addr := common.BytesToAddress([]byte{111})

	err := c.cancelProxy(addr)
	if err.Error() != "not set proxy" {
		t.Error(err)
	}
}

// Test cancelProxy
// 数据库中无记录
func TestCancelProxyNoProxy(t *testing.T) {
	context := newcontext()
	c := newElectionContext(context)
	addr := common.BytesToAddress([]byte{111})
	// 数据库中塞条记录
	voter := Voter{
		Owner:     addr,
		TimeStamp: big.NewInt(1531328500),
	}
	if err := c.setVoter(voter); err != nil {
		t.Errorf("addr: %s, error: %s", voter.Owner, err)
	}

	err := c.cancelProxy(addr)
	if err.Error() != "not set proxy" {
		t.Error(err)
	}
}

// Test cancelProxy
func TestCancelProxy(t *testing.T) {
	// 111->10
	context := newcontext()
	c := newElectionContext(context)

	addr := common.BytesToAddress([]byte{111})
	err := setProxy(t, c)
	if err != nil {
		t.Error(err)
	}

	// 取消代理
	err = c.cancelProxy(addr)
	if err != nil {
		t.Error(err)
	}

	voteCount := c.calculateVoteCount(big.NewInt(100))
	for i := 0; i < len(candidates); i++ {
		candi := c.getCandidate(candidates[i])
		if candi.VoteCount.Cmp(voteCount) != 0 {
			t.Errorf("The vote count %v is Wrong!", candi.VoteCount)
		}
	}

	if _, err := checkValid(t, c); err != nil {
		t.Error(err)
	}
}

// Test startProxy
// 已经是代理了
func TestStartProxyWithIsProxy(t *testing.T) {
	context := newcontext()
	c := newElectionContext(context)
	err := setProxy(t, c)
	if err != nil {
		t.Error(err)
	}

	err = c.startProxy(common.BytesToAddress([]byte{10}))
	if err.Error() != "startProxy proxy is already started" {
		t.Error(err)
	}
}

// Test startProxy
// 设置了代理的不可以成为代理
func TestStartProxyWithSetProxy(t *testing.T) {
	context := newcontext()
	c := newElectionContext(context)
	err := setProxy(t, c)
	if err != nil {
		t.Error(err)
	}

	err = c.startProxy(common.BytesToAddress([]byte{111}))
	if err.Error() != "account that uses a proxy is not allowed to become a proxy" {
		t.Error(err)
	}
}

// Test stopProxy
// 数据库中无记录
func TestStopProxyNoRecord(t *testing.T) {
	context := newcontext()
	c := newElectionContext(context)
	addr := common.BytesToAddress([]byte{111})
	err := c.stopProxy(addr)
	if err.Error() != "stopProxy proxy does not exist." {
		t.Error(err)
	}
}

// Test stopProxy
// 不是代理
func TestStopProxyNotProxy(t *testing.T) {
	context := newcontext()
	c := newElectionContext(context)

	voter := newVoter()
	voter.Owner = common.BytesToAddress([]byte{111})
	if err := c.setVoter(voter); err != nil {
		t.Errorf("addr: %s, error: %s", voter.Owner, err)
	}

	err := c.stopProxy(common.BytesToAddress([]byte{111}))
	if err.Error() != "stopProxy address is not proxy" {
		t.Error(err)
	}
}

func TestStartAndStopProxy(t *testing.T) {
	// addr: common.BytesToAddress([]byte{111}) 有10票
	// proxy: common.BytesToAddress([]byte{10}) 有100票
	// addr1: common.BytesToAddress([]byte{50}) 有20票
	// 一开始proxy是代理
	// addr和addr1都设置proxy为其代理
	// proxy停止代理，这个时候proxy身上有自身的票，addr和addr1的票
	// addr取消代理， 这个时候proxy身上有自身的票，addr1的票
	// proxy重新开始代理,接受新的人给它的代理

	context := newcontext()
	c := newElectionContext(context)
	addr := common.BytesToAddress([]byte{111})
	addr1 := common.BytesToAddress([]byte{50})
	proxy := common.BytesToAddress([]byte{10})
	c.context.GetStateDb().AddBalance(addr1, big.NewInt(0).Mul(big.NewInt(20), big.NewInt(1e18)))
	if err := c.stake(addr1, big.NewInt(20)); err != nil {
		t.Errorf("stake failed, addr: %s, error: %s", addr.String(), err)
	}

	// addr 设置 proxy为代理
	err := setProxy(t, c)
	if err != nil {
		t.Error(err)
	}

	// addr1 设置 proxy为代理
	err = c.setProxy(addr1, proxy)
	if err != nil {
		t.Error(err)
	}
	if _, err := checkValid(t, c); err != nil {
		t.Error(err)
	}

	// 停止代理，原先代理的票都还有效
	if err := c.stopProxy(proxy); err != nil {
		t.Errorf("stop proxy, addr: %s, error: %s", proxy.String(), err)
	}
	if ctx, ok := c.context.(*testContext); ok {
		ctx.SetTime(big.NewInt(1531795552))
	}
	if err := c.voteWitnesses(proxy, candidates); err != nil {
		t.Errorf("vote, addr: %s, error: %s", proxy.String(), err)
	}
	if _, err := checkValid(t, c); err != nil {
		t.Error(err)
	}

	// 取消addr交给proxy的代理
	if err := c.cancelProxy(addr); err != nil {
		t.Errorf("cancel proxy, addr: %s, error: %s", addr.String(), err)
	}
	if _, err := checkValid(t, c); err != nil {
		t.Error(err)
	}

	// 重新开始代理
	if err := c.startProxy(proxy); err != nil {
		t.Errorf("start proxy, addr: %s, error: %s", proxy.String(), err)
	}
	if err := c.unStake(addr1); err != nil {
		t.Errorf("unstake, addr: %s, error: %s", addr1.String(), err)
	}
	c.context.GetStateDb().AddBalance(addr1, big.NewInt(0).Mul(big.NewInt(20), big.NewInt(1e18)))
	if err := c.stake(addr1, big.NewInt(30)); err != nil {
		t.Errorf("stake, addr: %s, error: %s", addr1.String(), err)
	}

	// 设置代理
	if err := c.setProxy(addr1, proxy); err != nil {
		t.Errorf("set proxy, addr: %s, error: %s", addr1.String(), err)
	}
	if _, err := checkValid(t, c); err != nil {
		t.Error(err)
	}
}

func TestStopAndSetProxy(t *testing.T) {
	// addr: common.BytesToAddress([]byte{111}) 有10票
	// proxy: common.BytesToAddress([]byte{10}) 有100票
	// addr1: common.BytesToAddress([]byte{50}) 有20票
	// 一开始proxy和addr1都是代理
	// addr开始设置proxy为其代理
	// proxy停止代理，并设置addr1为其代理，  这个时候addr1身上有自身的票，proxy的票和addr的票
	// addr1取消代理， 这个时候addr1身上有自身的票，proxy的票
	context := newcontext()
	c := newElectionContext(context)
	addr := common.BytesToAddress([]byte{111})
	addr1 := common.BytesToAddress([]byte{50})
	proxy := common.BytesToAddress([]byte{10})

	// addr 设置 proxy为代理
	err := setProxy(t, c)
	if err != nil {
		t.Error(err)
	}

	voteCount := c.calculateVoteCount(big.NewInt(10))
	voteCount1 := c.calculateVoteCount(big.NewInt(100))
	proxyVoteCount := big.NewInt(0)
	proxyVoteCount.Add(proxyVoteCount, voteCount)

	stake := Stake{
		Owner:      addr1,
		StakeCount: big.NewInt(20),
	}
	if err := c.setStake(stake); err != nil {
		t.Errorf("stake, addr: %s, error: %s", stake.Owner.String(), err)
	}

	// addr1 投票
	// addr1 开始代理
	if err := c.startProxy(addr1); err != nil {
		t.Errorf("start proxy, addr: %s, error: %s", addr1.String(), err)
	}
	err = c.voteWitnesses(addr1, candidates)
	voteCount2 := c.calculateVoteCount(big.NewInt(20))
	totalVoteCount := big.NewInt(0)
	totalVoteCount.Add(voteCount, voteCount1)
	totalVoteCount.Add(totalVoteCount, voteCount2)

	if err != nil {
		t.Error(err)
	}
	for i := 0; i < len(candidates); i++ {
		candi := c.getCandidate(candidates[i])
		if candi.VoteCount.Cmp(totalVoteCount) != 0 {
			t.Errorf("The vote count %v is Wrong! Expected: %d", candi.VoteCount, totalVoteCount)
		}
	}

	// proxy 停止代理
	err = c.stopProxy(proxy)
	if err != nil {
		t.Error(err)
	}

	// proxy 设置代理
	if ctx, ok := c.context.(*testContext); ok {
		ctx.SetTime(big.NewInt(1531795552))
	}
	err = c.setProxy(proxy, addr1)
	totalVoteCount.Sub(totalVoteCount, voteCount1)
	voteCount1 = c.calculateVoteCount(big.NewInt(100))
	totalVoteCount.Add(totalVoteCount, voteCount1)

	if err != nil {
		t.Error(err)
	}

	for i := 0; i < len(candidates); i++ {
		candi := c.getCandidate(candidates[i])
		if candi.VoteCount.Cmp(totalVoteCount) != 0 {
			t.Errorf("The vote count %v is Wrong!", candi.VoteCount)
		}
	}

	// addr 取消 proxy代理
	if err := c.cancelProxy(addr); err != nil {
		t.Errorf("cancel proxy, addr: %s, error: %s", addr.String(), err)
	}
	totalVoteCount.Sub(totalVoteCount, voteCount)

	for i := 0; i < len(candidates); i++ {
		candi := c.getCandidate(candidates[i])
		if candi.VoteCount.Cmp(totalVoteCount) != 0 {
			t.Errorf("The vote count %v is Wrong!", candi.VoteCount)
		}
	}
}

func setProxy(t *testing.T, c electionContext) error {
	addr := common.BytesToAddress([]byte{111})
	proxy := common.BytesToAddress([]byte{10})
	// 账户addr抵押
	c.context.GetStateDb().AddBalance(addr, big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)))
	if err := c.stake(addr, big.NewInt(10)); err != nil {
		t.Errorf("stake, addr: %s, error: %s", addr.String(), err)
	}

	// 账户proxy抵押
	c.context.GetStateDb().AddBalance(proxy, big.NewInt(0).Mul(big.NewInt(100), big.NewInt(1e18)))
	if err := c.stake(proxy, big.NewInt(100)); err != nil {
		t.Errorf("stake, addr: %s, error: %s", proxy.String(), err)
	}

	// 账户proxy注册成为代理
	err := c.startProxy(proxy)
	if err != nil {
		return err
	}

	// 账户addr设置proxy
	err = c.setProxy(addr, proxy)
	if err != nil {
		return err
	}

	// 设置候选人，借用candiInfos的website、name、url信息
	for i := 0; i < len(candidates); i++ {
		if err := c.registerWitness(candidates[i], candiInfos[i].url, candiInfos[i].website, candiInfos[i].name); err != nil {
			return fmt.Errorf("registerWitness addr: %s, error: %s", candidates[i].String(), err)
		}
	}

	// 代理人投票
	err = c.voteWitnesses(proxy, candidates)
	if err != nil {
		return err
	}

	if _, err := checkValid(t, c); err != nil {
		return err
	}
	return nil
}

func TestRegisterWitness(t *testing.T) {
	context := newcontext()
	ec := newElectionContext(context)

	addr1 := common.HexToAddress("41b0db166cfdf1c4ba3ce657171482a9aa55cc93")
	addr2 := common.HexToAddress("08b467a881ec34b668254aa956e0c46f9c3b2b83")
	addr3 := common.HexToAddress("0c0292587ccdc76b8f449002a017bc9479ff0a88")
	addr4 := common.HexToAddress("0a0292587ccdc76b8f449002a017bc9479ff0a88")
	addr5 := common.HexToAddress("0b0292587ccdc76b8f449002a017bc9479ff0a88")
	addr6 := common.HexToAddress("0d0292587ccdc76b8f449002a017bc9479ff0a88")
	addr7 := common.HexToAddress("0e0292587ccdc76b8f449002a017bc9479ff0a88")

	t.Logf("addr1: %v", addr1.Hex())
	t.Logf("addr2: %v", addr2.Hex())
	t.Logf("addr3: %v", addr3.Hex())
	t.Logf("addr4: %v", addr4.Hex())

	complicatedErr := fmt.Errorf("complicated error, not compare error content")
	// 注册见证人的测试用例，err为nil代表需要注册成功
	ts := []struct {
		addr     common.Address
		url      []byte
		website  []byte
		name     []byte
		err      error
		matchErr bool   // 是否对error内容进行匹配，p2p类别的错误不易的，可设置为false
		desc     string // 本case的描述
	}{
		{addr1, []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHNAAfnqXNsxMwJf6QjJFRmVK7iB32U9owwK9KfeLFxEA7"), []byte("www.testnet1.site"), []byte("node1"), nil, true, "node1 success"},
		{addr1, []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHNAAfnqXNsxMwJf6QjJFRmVK7iB32U9owwK9KfeLFxEA7"), []byte("www.testnet1.site"), []byte("node1"), ErrCandiAlreadyRegistered, true, "node1 dup-register"},
		{addr2, []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHcch6yuBCgC5nPPSK3Yp7Es4c4eenxAeK167pYwUvNjRo"), []byte("www.testnet2.site"), []byte("node2"), nil, true, "node2 success"},
		{addr3, []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHJFKr2bzUnMr1NbeyYbYJa3RXT18cEu7cNDrHWjg8XYKB"), []byte("www.testnet3.site"), []byte("node3"), nil, true, "node3 success"},
		{addr4, []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHfop9dnUHHmtBXVkLB5UauAmACtrsEX5H5t6oCRpdL198"), []byte("www.testnet4.site"), []byte("s"), ErrCandiNameLenInvalid, true, "node4 too short name"},
		{addr4, []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHfop9dnUHHmtBXVkLB5UauAmACtrsEX5H5t6oCRpdL198"), []byte("www.testnet4.site"), []byte("tooloooooooooooooname"), ErrCandiNameLenInvalid, true, "node4 too long name"},
		{addr4, []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHfop9dnUHHmtBXVkLB5UauAmACtrsEX5H5t6oCRpdL198"), []byte("ww"), []byte("right name"), ErrCandiUrlLenInvalid, true, "node4 too short website"},
		{addr4, []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHfop9dnUHHmtBXVkLB5UauAmACtrsEX5H5t6oCRpdL198"), []byte("www.looooooooooooooooooooooooooooooooooooongwebsite.com/looog"), []byte("right name"), ErrCandiUrlLenInvalid, true, "node4 too long website"},
		{addr4, []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHfop9dnUHHmtBXVkLB5UauAmACtrsEX5H5t6oCRpdL198"), []byte("www.testnet4.site"), []byte("ABCEFacd"), ErrCandiNameInvalid, true, "node4 name should lowercase"},
		{addr4, []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHfop9dnUHHmtBXVkLB5UauAmACtrsEX5H5t6oCRpdL198"), []byte("www.testnet4.site"), []byte("acd xyz"), ErrCandiNameInvalid, true, "node4 name should only contain lowercase letter and digits"},
		{addr4, []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHfop9dnUHHmtBXVkLB5UauAmACtrsEX5H5t6oCRpdL198"), []byte("www.testnet4.site"), []byte("acd.xyz"), ErrCandiNameInvalid, true, "node4 name should only contain lowercase letter and digits"},
		{addr4, []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHfop9dnUHHmtBXVkLB5UauAmACtrsEX5H5t6oCRpdL198"), []byte("www.testnet4.site"), []byte("node3"), ErrCandiInfoDup, true, "node4 dup name"},
		{addr4, []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHfop9dnUHHmtBXVkLB5UauAmACtrsEX5H5t6oCRpdL198"), []byte("www.testnet3.site"), []byte("node4"), ErrCandiInfoDup, true, "node4 dup website"},
		{addr4, []byte("/ip9/127.0.0.1/tcp/30303/ipfs/1kHfop9dnUHHmtBXVkLB5UauAmACtrsEX5H5t6oCRpdL198"), []byte("www.testnet4.site"), []byte("node4"), fmt.Errorf("registerWitness node url is error: no protocol with name ip9"), true, "node4 node url invalid: ip9"},
		{addr4, []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHfop9dnUHHmtskLfs5fdafemACtfsEX5H5t6oCRpdL1"), []byte("www.testnet4.site"), []byte("node4"), complicatedErr, false, "node4 node url invalid: id is less"},
		{addr4, []byte("/ip4/127.0.0.1/txp/30303/ipfs/1kHfop9dnUHHmtskLfs5fdafemACtfsEX5H5t6oCRpdL198"), []byte("www.testnet4.site"), []byte("node4"), fmt.Errorf("registerWitness node url is error: no protocol with name txp"), true, "node4 node url invalid: not tcp"},
		{addr4, []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1111111111111111111111111111111111111111RpdL198"), []byte("www.testnet4.site"), []byte("node4"), complicatedErr, false, "node4 node url invalid: id is invalid"},
		{addr4, []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHfop9dnUHHmtBXVkLB5UauAmACtrsEX5H5t6oCRpdL198"), []byte("www.testnet4.site"), []byte("node4"), nil, true, "node4 success"},
		{addr5, []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHGG8L1DTrVG3Cad479Q32oGmFAiEjLFwxzNyXH3ehGo73"), []byte("www.testnet5.site"), []byte("20charactornaaaaaame"), nil, true, "node5 success"},
		{addr6, []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHmExX4yutwBZLbRsYHq59KfgiM1LUJFW2JSPeSCcBf7rH"), []byte("www"), []byte("node6"), nil, true, "node6 success"},
		{addr7, []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHSMhv82q5thJkdeJzxCVW8tdXwaDThBZWsH2Q9KUGGFUq"), []byte("www.just60charactor.com/loooooooooooooooooooooooooooooooooog"), []byte("node7"), nil, true, "node7 success"},
	}

	for i, c := range ts {
		err := ec.registerWitness(c.addr, c.url, c.website, c.name)
		if !reflect.DeepEqual(err, c.err) {
			if c.matchErr {
				t.Errorf("TestRegisterWitness case %d, case discrition: %s, want err :%v, got:%v", i, c.desc, c.err, err)
			}
		}
	}

	candis := getAllCandidate(context.GetStateDb())
	for _, candi := range candis {
		t.Logf("333 addr: %v, voteCount: %v, active: %v", candi.Owner.Hex(), candi.VoteCount, candi.Active)
	}

	err := ec.unregisterWitness(addr1)
	if err != nil {
		t.Errorf("TestRegisterWitness unregisterWitness err:%v", err)
	}

	candis = getAllCandidate(context.GetStateDb())
	for _, candi := range candis {
		t.Logf("444 addr: %v, voteCount: %v, active: %v", candi.Owner.Hex(), candi.VoteCount, candi.Active)
	}

	err = ec.unregisterWitness(addr3)
	if err != nil {
		t.Errorf("TestRegisterWitness unregisterWitness err:%v", err)
	}

	candis = getAllCandidate(context.GetStateDb())
	for _, candi := range candis {
		t.Logf("555 addr: %v, voteCount: %v, active: %v", candi.Owner.Hex(), candi.VoteCount, candi.Active)
	}

	err = ec.unregisterWitness(addr2)
	if err != nil {
		t.Errorf("TestRegisterWitness unregisterWitness err:%v", err)
	}

	candis = getAllCandidate(context.GetStateDb())
	for _, candi := range candis {
		t.Logf("666 addr: %v, voteCount: %v, active: %v", candi.Owner.Hex(), candi.VoteCount, candi.Active)
	}
}

func TestRegisterProxy(t *testing.T) {
	context := newcontext()
	ec := newElectionContext(context)

	addr1 := common.HexToAddress("41b0db166cfdf1c4ba3ce657171482a9aa55cc93")
	addr2 := common.HexToAddress("08b467a881ec34b668254aa956e0c46f9c3b2b83")
	addr3 := common.HexToAddress("0c0292587ccdc76b8f449002a017bc9479ff0a88")

	t.Logf("addr1: %v", addr1.Hex())
	t.Logf("addr2: %v", addr2.Hex())
	t.Logf("addr3: %v", addr3.Hex())

	err := ec.startProxy(addr1)
	if err != nil {
		t.Errorf("TestRegisterProxy startProxy err: %v", err)
	}

	proxys := getAllProxy(context.GetStateDb())
	for _, proxy := range proxys {
		t.Logf("111 proxy %v", proxy.Owner.Hex())
	}

	err = ec.stopProxy(addr1)
	if err != nil {
		t.Errorf("TestRegisterProxy startProxy err: %v", err)
	}

	proxys = getAllProxy(context.GetStateDb())
	for _, proxy := range proxys {
		t.Logf("222 proxy %v", proxy.Owner.Hex())
	}

	err = ec.startProxy(addr2)
	if err != nil {
		t.Errorf("TestRegisterProxy startProxy err: %v", err)
	}

	err = ec.startProxy(addr3)
	if err != nil {
		t.Errorf("TestRegisterProxy startProxy err: %v", err)
	}

	proxys = getAllProxy(context.GetStateDb())
	for _, proxy := range proxys {
		t.Logf("333 proxy %v", proxy.Owner.Hex())
	}
}

func TestStake(t *testing.T) {
	context := newcontext()
	ec := newElectionContext(context)

	addr1 := common.HexToAddress("41b0db166cfdf1c4ba3ce657171482a9aa55cc93")

	context.GetStateDb().AddBalance(addr1, big.NewInt(0).Mul(big.NewInt(10000000000), big.NewInt(10000000000)))

	err := ec.stake(addr1, big.NewInt(20))
	if err != nil {
		t.Errorf("TestStake stake err:%v ", err)
	}

	t.Logf("111 addr1 balance: %v", context.GetStateDb().GetBalance(addr1))

	stake := ec.getStake(addr1)
	t.Logf("111 addr1 stake: %v", stake.StakeCount)

	err = ec.unStake(addr1)
	if err.Error() != "cannot unstake in 24 hours" {
		t.Errorf("TestStake unStake err:%v ", err)
	}

	t.Logf("222 addr1 balance after unStake: %v", context.GetStateDb().GetBalance(addr1))

	stake = ec.getStake(addr1)
	t.Logf("222 addr1 stake after unStake: %v", stake.StakeCount)

	if ctx, ok := context.(*testContext); ok {
		ctx.SetTime(big.NewInt(0).Add(context.GetTime(), big.NewInt(3600*24+1)))
	}
	err = ec.unStake(addr1)
	if err != nil {
		t.Errorf("TestStake unStake err:%v ", err)
	}

	t.Logf("333 addr1 balance after unStake: %v", context.GetStateDb().GetBalance(addr1))

	stake = ec.getStake(addr1)
	t.Logf("333 addr1 stake after unStake: %v", stake.StakeCount)

	err = ec.stake(addr1, big.NewInt(-20))
	if err.Error() != "stake stakeCount less than 0" {
		t.Errorf("TestStake stake err:%v ", err)
	}

	t.Logf("444 addr1 balance: %v", context.GetStateDb().GetBalance(addr1))

	stake = ec.getStake(addr1)
	t.Logf("444 addr1 stake: %v", stake.StakeCount)
}

func TestExtractBounty(t *testing.T) {
	context := newcontext()
	ec := newElectionContext(context)
	if err := ec.setCandidate(candidate); err != nil {
		t.Errorf("candiates: %s, error: %s", candidate.Owner, err)
	}
	if err := ec.extractOwnBounty(candidate.Owner); err != nil {
		t.Error(err)
	}
	candidate1 := ec.getCandidate(candidate.Owner)
	if candidate1.TotalBounty.Cmp(candidate1.ExtractedBounty) != 0 {
		t.Errorf("extracted bounty %v not equal to totalBouty %v", candidate1.ExtractedBounty, candidate1.TotalBounty)
	}
}

func TestGrantBounty(t *testing.T) {
	context := newcontext()

	if err := setRestBounty(context.GetStateDb(), bounty); err != nil {
		t.Error(err)
	}
	// enough to pay
	if rest, err := GrantBounty(context.GetStateDb(), big.NewInt(1e17)); err != nil {
		t.Error(err)
	} else if rest.Cmp(big.NewInt(9e17)) != 0 {
		t.Error("the rest of bounty error")
	}

	// not enough to pay
	if rest, err := GrantBounty(context.GetStateDb(), big.NewInt(1e18)); err == nil {
		t.Error("the rest of bounty should be not enough to pay")
	} else if rest.Cmp(big.NewInt(9e17)) != 0 {
		t.Error("the rest of bounty error")
	}

	// just to pay
	if rest, err := GrantBounty(context.GetStateDb(), big.NewInt(9e17)); err != nil {
		t.Log(err)
	} else if rest.Cmp(big.NewInt(0)) != 0 {
		t.Error("the rest of bounty error")
	}
}

func TestCalculateVote(t *testing.T) {
	context := newcontext()
	ec := newElectionContext(context)
	stakeCount := big.NewInt(10000000)

	tests := []struct {
		curTime *big.Int
		stake   *big.Int
		votes   *big.Int
	}{
		{eraTimeStamp, stakeCount, big.NewInt(10000000)},           // 半衰期开始时
		{big.NewInt(1562256000), stakeCount, big.NewInt(14142135)}, // 半衰期半个周期：26周
		{big.NewInt(1578067200), stakeCount, big.NewInt(20000000)}, // 半衰期一个周期：52周
		{big.NewInt(1609430400), stakeCount, big.NewInt(40000000)}, // 半衰期二个周期：104周
	}

	for i, ts := range tests {
		if ctx, ok := context.(*testContext); ok {
			ctx.SetTime(ts.curTime)
		}
		voteCount := ec.calculateVoteCount(stakeCount)
		if voteCount.Cmp(ts.votes) != 0 {
			t.Errorf("case %d error, time: %s, stake: %s, want votes: %s, got votes: %s", i, ts.curTime.String(), ts.stake.String(), ts.votes.String(), voteCount)
		}
	}
}

var operates []string            // 有哪些操作
var alreadySet map[byte]struct{} // alreadySet用来标记节点是否已经扩展过

func dfsState(t *testing.T, c electionContext, address common.Address, checkFn func(byte, string, byte) error) error {
	// 进入时保存当前数据库状态，退出时恢复
	snap := c.context.GetStateDb().Snapshot()
	defer c.context.GetStateDb().RevertToSnapshot(snap)
	currentState := checkState(c, address)
	for _, op := range operates {
		snap1 := c.context.GetStateDb().Snapshot()
		err := operateOnState(c, op, address)
		if err == nil {
			if _, err = checkValid(t, c); err != nil {
				return err
			}
			nextState := checkState(c, address)
			// 比较到达的状态是否与预期状态一致
			if err = checkFn(currentState, op, nextState); err != nil {
				return err
			}
			// 如果是新的状态加入到队列中
			if _, ok := alreadySet[nextState]; !ok {
				alreadySet[nextState] = struct{}{}
				err = dfsState(t, c, address, checkFn)
				if err != nil {
					return err
				}
			}
		}
		c.context.GetStateDb().RevertToSnapshot(snap1)
	}
	return nil
}

func TestVoteAndProxyState1(t *testing.T) {
	// 定义两种角色addr做为普通用户，操作只有投票，取消投票，设置代理，取消设置代理四种
	// proxy做为代理用户，操作有投票，取消投票，设置代理，取消设置代理四种，还有startProxy和stopProxy
	// addrStateMap存预先设置好addr相关的状态转化图，proxyStateMap存proxy的状态转化图
	addr := common.BytesToAddress([]byte{111})
	operates = []string{
		"voteWitnesses",
		"cancelVote",
		"setProxy",
		"cancelProxy",
	}

	// addr的状态转化图
	addrStateMap := make(map[byte]map[string]byte)

	// 0无投票无代理，4无投票有代理，8有投票无代理
	addrStateMap[0] = make(map[string]byte)
	addrStateMap[8] = make(map[string]byte)
	addrStateMap[2] = make(map[string]byte)

	addrStateMap[0]["voteWitnesses"] = 8
	addrStateMap[0]["setProxy"] = 2
	addrStateMap[0]["cancelVote"] = 0
	addrStateMap[8]["cancelVote"] = 0
	addrStateMap[8]["setProxy"] = 2
	addrStateMap[8]["voteWitnesses"] = 8
	addrStateMap[2]["voteWitnesses"] = 8
	addrStateMap[2]["cancelProxy"] = 0
	addrStateMap[2]["setProxy"] = 2

	// 抵押一类的初始操作
	context := newcontext()
	c := newElectionContext(context)
	initForStateTest(t, c)

	alreadySet = make(map[byte]struct{})
	alreadySet[0] = struct{}{}
	checkFn := func(current byte, op string, next byte) error {
		if addrStateMap[current][op] != next {
			return fmt.Errorf("state error ,current state %d, op %s, nextState %d, expected %d", current, op, next, addrStateMap[current][op])
		}
		return nil
	}

	if err := dfsState(t, c, addr, checkFn); err != nil {
		t.Error(err)
	}
}

func TestVoteAndProxyState2(t *testing.T) {
	// 定义两种角色addr做为普通用户，操作只有投票，取消投票，设置代理，取消设置代理四种
	// proxy做为代理用户，操作有投票，取消投票，设置代理，取消设置代理四种，还有startProxy和stopProxy
	// addrStateMap存预先设置好addr相关的状态转化图，proxyStateMap存proxy的状态转化图
	proxy := common.BytesToAddress([]byte{10})
	proxyStateMap := make(map[byte]map[string]byte)

	// 代理的一些状态转变不是自身发起的
	proxyOperates := []string{
		"voteWitnesses",
		"cancelVote",
		"setProxy",
		"cancelProxy",
		"startProxy",
		"stopProxy",
		"addrSetProxy",
		"addrCancelProxy",
	}
	operates = proxyOperates
	// 4个标志位，一共16种，其中有6中状态是不合法的
	for i := 0; i < 14; i++ {
		if i == 6 || i == 7 || i == 10 {
			continue
		}
		proxyStateMap[byte(i)] = make(map[string]byte)
	}
	proxyStateMap[0]["voteWitnesses"] = 8
	proxyStateMap[0]["cancelVote"] = 0
	proxyStateMap[0]["startProxy"] = 4
	proxyStateMap[0]["setProxy"] = 2
	proxyStateMap[0]["stopProxy"] = 0

	proxyStateMap[1]["voteWitnesses"] = 9
	proxyStateMap[1]["cancelVote"] = 1
	proxyStateMap[1]["startProxy"] = 5
	proxyStateMap[1]["setProxy"] = 3
	proxyStateMap[1]["stopProxy"] = 1
	proxyStateMap[1]["addrCancelProxy"] = 0

	proxyStateMap[2]["voteWitnesses"] = 8
	proxyStateMap[2]["setProxy"] = 2
	proxyStateMap[2]["stopProxy"] = 2
	proxyStateMap[2]["cancelProxy"] = 0

	proxyStateMap[3]["voteWitnesses"] = 9
	proxyStateMap[3]["setProxy"] = 3
	proxyStateMap[3]["stopProxy"] = 3
	proxyStateMap[3]["cancelProxy"] = 1
	proxyStateMap[3]["addrCancelProxy"] = 2

	proxyStateMap[4]["voteWitnesses"] = 12
	proxyStateMap[4]["cancelVote"] = 4
	proxyStateMap[4]["startProxy"] = 4
	proxyStateMap[4]["stopProxy"] = 0
	proxyStateMap[4]["addrSetProxy"] = 5

	proxyStateMap[5]["voteWitnesses"] = 13
	proxyStateMap[5]["cancelVote"] = 5
	proxyStateMap[5]["startProxy"] = 5
	proxyStateMap[5]["stopProxy"] = 1
	proxyStateMap[5]["addrSetProxy"] = 5
	proxyStateMap[5]["addrCancelProxy"] = 4

	proxyStateMap[8]["cancelVote"] = 0
	proxyStateMap[8]["voteWitnesses"] = 8
	proxyStateMap[8]["setProxy"] = 2
	proxyStateMap[8]["startProxy"] = 12
	proxyStateMap[8]["stopProxy"] = 8

	proxyStateMap[9]["cancelVote"] = 1
	proxyStateMap[9]["setProxy"] = 3
	proxyStateMap[9]["voteWitnesses"] = 9
	proxyStateMap[9]["stopProxy"] = 9
	proxyStateMap[9]["startProxy"] = 13
	proxyStateMap[9]["addrCancelProxy"] = 8

	proxyStateMap[12]["cancelVote"] = 4
	proxyStateMap[12]["voteWitnesses"] = 12
	proxyStateMap[12]["startProxy"] = 12
	proxyStateMap[12]["stopProxy"] = 8
	proxyStateMap[12]["addrSetProxy"] = 13

	proxyStateMap[13]["cancelVote"] = 5
	proxyStateMap[13]["voteWitnesses"] = 13
	proxyStateMap[13]["startProxy"] = 13
	proxyStateMap[13]["stopProxy"] = 9
	proxyStateMap[13]["addrSetProxy"] = 13
	proxyStateMap[13]["addrCancelProxy"] = 12

	checkFn := func(current byte, op string, next byte) error {
		if proxyStateMap[current][op] != next {
			return fmt.Errorf("state error ,current state %d, op %s, nextState %d, expected %d", current, op, next, proxyStateMap[current][op])
		}
		return nil
	}

	// 抵押一类的初始操作
	context := newcontext()
	c := newElectionContext(context)
	initForStateTest(t, c)

	alreadySet = make(map[byte]struct{})

	alreadySet[4] = struct{}{}
	if err := dfsState(t, c, proxy, checkFn); err != nil {
		t.Error(err)
	}
}

func initForStateTest(t *testing.T, c electionContext) {
	addr := common.BytesToAddress([]byte{111})
	proxy := common.BytesToAddress([]byte{10})
	proxy1 := common.BytesToAddress([]byte{50})
	if ctx, ok := c.context.(*testContext); ok {
		ctx.SetTime(new(big.Int).Set(eraTimeStamp))
	}

	c.context.GetStateDb().AddBalance(addr, big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)))
	if err := c.stake(addr, big.NewInt(10)); err != nil {
		t.Errorf("stake error, addr: %s, error: %s", addr.String(), err)
	}
	c.context.GetStateDb().AddBalance(proxy, big.NewInt(0).Mul(big.NewInt(100), big.NewInt(1e18)))
	if err := c.stake(proxy, big.NewInt(100)); err != nil {
		t.Errorf("stake error, addr: %s, error: %s", proxy.String(), err)
	}
	c.context.GetStateDb().AddBalance(proxy1, big.NewInt(0).Mul(big.NewInt(1000), big.NewInt(1e18)))
	if err := c.stake(proxy1, big.NewInt(1000)); err != nil {
		t.Errorf("stake error, addr: %s, error: %s", proxy1.String(), err)
	}
	if err := c.startProxy(proxy); err != nil {
		t.Errorf("start proxy, addr: %s, error: %s", proxy.String(), err)
	}
	if err := c.startProxy(proxy1); err != nil {
		t.Errorf("start proxy, addr: %s, error: %s", proxy1.String(), err)
	}
	if err := c.voteWitnesses(proxy1, candidates); err != nil {
		t.Errorf("vote, addr: %s, error: %s", proxy1.String(), err)
	}

	// 下面只是为了把addr塞进数据库
	if err := c.startProxy(addr); err != nil {
		t.Errorf("start proxy, addr: %s, error: %s", addr.String(), err)
	}
	if err := c.stopProxy(addr); err != nil {
		t.Errorf("start proxy, addr: %s, error: %s", addr.String(), err)
	}

	for i, candi := range candidates {
		website := "www.testnet.info" + strconv.Itoa(i)
		name := "testinfo" + strconv.Itoa(i)
		p2pUrl := []byte(string(url)[:13] + strconv.Itoa(i) + string(url)[14:])
		if err := c.registerWitness(candi, p2pUrl, []byte(website), []byte(name)); err != nil {
			t.Errorf("register failed, addr: %s, error: %s", candi.String(), err)
		}
	}
}

func operateOnState(c electionContext, op string, address common.Address) error {
	addr := address
	proxy := common.BytesToAddress([]byte{10})
	switch op {
	case "setProxy":
		if bytes.Equal(address.Bytes(), common.BytesToAddress([]byte{10}).Bytes()) {
			proxy = common.BytesToAddress([]byte{50})
		}
	case "addrSetProxy":
		fallthrough
	case "addrCancelProxy":
		addr = common.BytesToAddress([]byte{111})
	}
	return operate(c, op, addr, proxy, candidates)
}

func checkState(c electionContext, address common.Address) byte {
	voter := c.getVoter(address)
	if !bytes.Equal(voter.Owner.Bytes(), address.Bytes()) {
		return 0
	}
	var result byte
	// 有投票,判断voteCandidates,或者无代理却有投票数(上次投了个空票)
	if len(voter.VoteCandidates) > 0 {
		result |= 1 << 3
	} else if bytes.Equal(voter.Proxy.Bytes(), emptyAddress.Bytes()) && voter.LastVoteCount.Sign() > 0 {
		result |= 1 << 3
	}
	// 是代理
	if voter.IsProxy {
		result |= 1 << 2
	}
	// 有代理
	if !bytes.Equal(voter.Proxy.Bytes(), emptyAddress.Bytes()) {
		result |= 1 << 1
	}
	// 有代理投票
	if voter.ProxyVoteCount.Sign() > 0 {
		result |= 1
	}
	return result
}

func operate(c electionContext, op string, address common.Address, proxy common.Address, candidates []common.Address) error {
	var err error
	switch op {
	case "voteWitnesses":
		if ctx, ok := c.context.(*testContext); ok {
			t := c.context.GetTime()
			ctx.SetTime(t.Add(t, big.NewInt(24*3600+1)))
		}
		err = c.voteWitnesses(address, candidates)
	case "cancelVote":
		err = c.cancelVote(address)
	case "setProxy":
		if ctx, ok := c.context.(*testContext); ok {
			t := c.context.GetTime()
			ctx.SetTime(t.Add(t, big.NewInt(24*3600+1)))
		}
		err = c.setProxy(address, proxy)
	case "cancelProxy":
		err = c.cancelProxy(address)
	case "startProxy":
		err = c.startProxy(address)
	case "stopProxy":
		err = c.stopProxy(address)
	case "registerWitness":
		website := "www.testnet.info"
		name := "testinfo"
		err = c.registerWitness(address, url, []byte(website), []byte(name))
	case "unregisterWitness":
		err = c.unregisterWitness(address)
	default:
		err = fmt.Errorf("method not found")
	}
	return err
}
