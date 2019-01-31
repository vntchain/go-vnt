package election

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"sort"
	"strings"
	"unicode"

	"github.com/pkg/errors"
	"github.com/vntchain/go-vnt/accounts/abi"
	"github.com/vntchain/go-vnt/common"
	inter "github.com/vntchain/go-vnt/core/vm/interface"
	"github.com/vntchain/go-vnt/log"
	"github.com/vntchain/go-vnt/vntp2p"
)

const (
	voteLimit = 30
	oneDay    = int64(24) * 3600
	oneWeek   = oneDay * 7
	year2019  = 1546272000
)

var (
	ErrCandiNameLenInvalid    = errors.New("the length of candidate's name should between [3, 20]")
	ErrCandiUrlLenInvalid     = errors.New("the length of candidate's website url should between [3, 60]")
	ErrCandiNameInvalid       = errors.New("candidate's name should consist of digits and lowercase letters")
	ErrCandiInfoDup           = errors.New("candidate's name, website url or node url is duplicated with a registered candidate")
	ErrCandiAlreadyRegistered = errors.New("candidate is already registered")
)

var (
	electionAddr = common.BytesToAddress([]byte{9})
	emptyAddress = common.Address{}
	eraTimeStamp = big.NewInt(year2019)

	// stake minimum time period
	unstakePeriod   = big.NewInt(oneDay)
	baseBounty      = big.NewInt(0).Mul(big.NewInt(1e+18), big.NewInt(1000))
	restTotalBounty = big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(1e9))
)

type Election struct{}

type electionContext struct {
	context inter.ChainContext
}

type Voter struct {
	Owner          common.Address   // 投票人的地址
	IsProxy        bool             // 是否是代理人
	ProxyVoteCount *big.Int         // 收到的代理的票数
	Proxy          common.Address   // 代理人
	LastVoteCount  *big.Int         // 票数
	TimeStamp      *big.Int         // 时间戳
	VoteCandidates []common.Address // 投了哪些人
}

// Candidate information of witness candidates.
// Tips: Modify CandidateList.Swap when adding element of Candidate.
type Candidate struct {
	Owner           common.Address // 候选人地址
	VoteCount       *big.Int       // 收到的票数
	Active          bool           // 当前是否是候选人
	Url             []byte         // 节点的URL
	TotalBounty     *big.Int       // 总奖励金额
	ExtractedBounty *big.Int       // 已提取奖励金额
	LastExtractTime *big.Int       // 上次提权时间
	Website         []byte         // 节点网站地址
	Name            []byte         // 节点名字
}

func (c *Candidate) String() string {
	return fmt.Sprintf("candidate, addr:%s, votes:%s, active:%v, url:%s, totalBounty: %v, extractedBounty: %v, lastExtractTime: %v, WebSite: %s, Name: %s\n",
		c.Owner.String(), c.VoteCount.String(), c.Active, string(c.Url), c.TotalBounty, c.ExtractedBounty, c.LastExtractTime, string(c.Website), string(c.Name))
}

func newVoter() Voter {
	return Voter{
		Owner:          emptyAddress,
		IsProxy:        false,
		ProxyVoteCount: big.NewInt(0),
		Proxy:          emptyAddress,
		LastVoteCount:  big.NewInt(0),
		TimeStamp:      big.NewInt(0),
		VoteCandidates: nil,
	}
}

func newCandidate() Candidate {
	return Candidate{
		Owner:     emptyAddress,
		VoteCount: big.NewInt(0),
		Active:    false,
	}
}

func (c *Candidate) votes() *big.Int {
	if c.Active {
		return c.VoteCount
	}

	one := big.NewInt(-1)
	return one.Mul(c.VoteCount, one)
}

// Equal two object is equal
func (c *Candidate) equal(d *Candidate) bool {
	return reflect.DeepEqual(c, d)
}

type CandidateList []Candidate

func (c CandidateList) Len() int {
	return len(c)
}

// Less for Sort interface, actually implement of c[i] more than c[j]
// Rule 1: 票数越多排名越靠前
// Rule 2: 票数相等，地址越小越靠前
//
// sort.Stable对于big.Int并不能真正的stable，所以排序还参考地址，并且排序不再使用stable
func (c CandidateList) Less(i, j int) bool {
	ret := c[i].votes().Cmp(c[j].votes())
	if ret != 0 {
		return ret > 0
	}

	return bytes.Compare(c[i].Owner.Bytes(), c[j].Owner.Bytes()) < 0
}

func (c CandidateList) Swap(i, j int) {
	c[i].Owner, c[j].Owner = c[j].Owner, c[i].Owner
	c[i].VoteCount, c[j].VoteCount = c[j].VoteCount, c[i].VoteCount
	c[i].Active, c[j].Active = c[j].Active, c[i].Active
	c[i].Url, c[j].Url = c[j].Url, c[i].Url
	c[i].TotalBounty, c[j].TotalBounty = c[j].TotalBounty, c[i].TotalBounty
	c[i].ExtractedBounty, c[j].ExtractedBounty = c[j].ExtractedBounty, c[i].ExtractedBounty
	c[i].LastExtractTime, c[j].LastExtractTime = c[j].LastExtractTime, c[i].LastExtractTime
	c[i].Website, c[j].Website = c[j].Website, c[i].Website
	c[i].Name, c[j].Name = c[j].Name, c[i].Name
}

// Sort
func (c CandidateList) Sort() {
	sort.Sort(c)
}

func (c CandidateList) dump() {
	fmt.Println("dump candidats list")
	for i, ca := range c {
		fmt.Printf("can:%d, addr:%s, votes:%s, active:%v \n", i, ca.Owner.String(), ca.VoteCount.String(), ca.Active)
	}
}

type Stake struct {
	Owner      common.Address // 抵押人地址
	StakeCount *big.Int       // 抵押的数量
	TimeStamp  *big.Int       // 时间戳
}

type Bounty struct {
	RestTotalBounty *big.Int // 剩余总激励，初始值10亿VNT
}

func newElectionContext(ctx inter.ChainContext) electionContext {
	return electionContext{
		context: ctx,
	}
}

func (e *Election) RequiredGas(input []byte) uint64 {
	return 0
}

func (e *Election) Run(ctx inter.ChainContext, input []byte) ([]byte, error) {
	nonce := ctx.GetStateDb().GetNonce(electionAddr)
	if nonce == 0 {
		setRestBounty(ctx.GetStateDb(), Bounty{restTotalBounty})
	}
	ctx.GetStateDb().SetNonce(electionAddr, nonce+1)
	abiJSON := `[
{"inputs":[{"name":"nodeUrl","type":"bytes"},{"name":"website","type":"bytes"},{"name":"nodeName","type":"bytes"}],"name":"registerWitness","outputs":[],"type":"function"},
{"inputs":[],"name":"unregisterWitness","outputs":[],"type":"function"},
{"inputs":[{"name":"candidate","type":"address[]"}],"name":"voteWitnesses","outputs":[],"type":"function"},
{"inputs":[],"name":"cancelVote","outputs":[],"type":"function"},
{"inputs":[],"name":"startProxy","outputs":[],"type":"function"},
{"inputs":[],"name":"stopProxy","outputs":[],"type":"function"},
{"inputs":[],"name":"cancelProxy","outputs":[],"type":"function"},
{"inputs":[{"name":"proxy","type":"address"}],"name":"setProxy","outputs":[],"type":"function"},
{"inputs":[{"name":"stakeCount","type":"uint256"}],"name":"stake","outputs":[],"type":"function"},
{"inputs":[],"name":"unStake","outputs":[],"type":"function"},
{"inputs":[],"name":"extractOwnBounty","outputs":[],"type":"function"}
]`
	electionABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, err
	}

	c := newElectionContext(ctx)
	methodName := "None"

	// input的组成见abi.Pack函数
	methodId := input[:4]
	methodArgs := input[4:]
	switch {
	case bytes.Equal(methodId, electionABI.Methods["registerWitness"].Id()):
		methodName = "registerWitness"
		type NodeInfo struct {
			NodeUrl  []byte
			Website  []byte
			NodeName []byte
		}
		var nodeInfo NodeInfo
		if err = electionABI.UnpackInput(&nodeInfo, "registerWitness", methodArgs); err == nil {
			err = c.registerWitness(ctx.GetOrigin(), nodeInfo.NodeUrl, nodeInfo.Website, nodeInfo.NodeName)
		}

	case bytes.Equal(methodId, electionABI.Methods["unregisterWitness"].Id()):
		methodName = "unregisterWitness"
		err = c.unregisterWitness(ctx.GetOrigin())

	case bytes.Equal(methodId, electionABI.Methods["voteWitnesses"].Id()):
		methodName = "voteWitnesses"
		var candidates []common.Address
		if err = electionABI.UnpackInput(&candidates, "voteWitnesses", methodArgs); err == nil {
			err = c.voteWitnesses(ctx.GetOrigin(), candidates)
		}

	case bytes.Equal(methodId, electionABI.Methods["cancelVote"].Id()):
		methodName = "cancelVote"
		err = c.cancelVote(ctx.GetOrigin())

	case bytes.Equal(methodId, electionABI.Methods["startProxy"].Id()):
		methodName = "startProxy"
		err = c.startProxy(ctx.GetOrigin())

	case bytes.Equal(methodId, electionABI.Methods["stopProxy"].Id()):
		methodName = "stopProxy"
		err = c.stopProxy(ctx.GetOrigin())

	case bytes.Equal(methodId, electionABI.Methods["cancelProxy"].Id()):
		methodName = "cancelProxy"
		err = c.cancelProxy(ctx.GetOrigin())

	case bytes.Equal(methodId, electionABI.Methods["setProxy"].Id()):
		methodName = "setProxy"
		var proxy common.Address
		if err = electionABI.UnpackInput(&proxy, "setProxy", methodArgs); err == nil {
			err = c.setProxy(ctx.GetOrigin(), proxy)
		}
	case bytes.Equal(methodId, electionABI.Methods["stake"].Id()):
		methodName = "stake"
		var stakeCount *big.Int
		if err = electionABI.UnpackInput(&stakeCount, "stake", methodArgs); err == nil {
			err = c.stake(ctx.GetOrigin(), stakeCount)
		}
	case bytes.Equal(methodId, electionABI.Methods["unStake"].Id()):
		methodName = "unStake"
		err = c.unStake(ctx.GetOrigin())
	case bytes.Equal(methodId, electionABI.Methods["extractOwnBounty"].Id()):
		methodName = "extractOwnBounty"
		err = c.extractOwnBounty(ctx.GetOrigin())
	}
	if err != nil {
		log.Error("call election contract err:", "method", methodName, "err", err)
	} else if methodName == "None" {
		log.Error("call election contract err: method doesn't exist")
		err = fmt.Errorf("call election contract err: method doesn't exist")
	}
	return nil, err
}

func (ec electionContext) registerWitness(address common.Address, url []byte, website []byte, name []byte) error {
	// get candidate from db
	candidate := ec.getCandidate(address)

	// if candidate already exists
	if bytes.Equal(candidate.Owner.Bytes(), address.Bytes()) {
		// if candidate is already active, just ignore
		if candidate.Active {
			log.Warn("registerWitness witness already exists", "address", address.Hex())
			return ErrCandiAlreadyRegistered
		}
	} else {
		// if candidate is not found in db
		// make a new candidate
		candidate.Owner = address
		candidate.VoteCount = big.NewInt(0)
	}

	// Sanity check
	if err := ec.checkCandi(address, string(name), string(website), string(url)); err != nil {
		return err
	}

	// Mark candidate as active
	candidate.Active = true
	candidate.Url = url
	candidate.Website = website
	candidate.Name = name

	// save candidate info db
	err := ec.setCandidate(candidate)
	if err != nil {
		log.Error("registerWitness setCandidate err.", "address", address.Hex(), "err", err)
		return err
	}

	return nil
}

// checkCandi 候选人基本参数的校验
func (ec electionContext) checkCandi(addr common.Address, name string, website string, url string) error {
	// length check
	if len(name) < 3 || len(name) > 20 {
		return ErrCandiNameLenInvalid
	}
	if len(website) < 3 || len(website) > 60 {
		return ErrCandiUrlLenInvalid
	}

	digitalAndLower := func(s string) bool {
		for _, ru := range s {
			if !unicode.IsDigit(ru) && !unicode.IsLower(ru) {
				return false
			}
		}
		return true
	}
	if !digitalAndLower(name) {
		return ErrCandiNameInvalid
	}

	// p2p node url format check
	if _, err := vntp2p.ParseNode(url); err != nil {
		return fmt.Errorf("registerWitness node url is error: %s", err)
	}

	// duplication check
	wits := getAllCandidate(ec.context.GetStateDb())
	for _, w := range wits {
		if w.Owner != addr && (string(w.Name) == name || string(w.Website) == website || string(w.Url) == url) {
			return ErrCandiInfoDup
		}
	}
	return nil
}

func (ec electionContext) unregisterWitness(address common.Address) error {
	// get candidate from db
	candidate := ec.getCandidate(address)

	// if candidate is not found in db
	if !bytes.Equal(candidate.Owner.Bytes(), address.Bytes()) {
		log.Warn("unregisterWitness unregister unknown witness.", "address", address.Hex())
		return fmt.Errorf("unregisterWitness unregister unknown witness.")
	}

	// if candidate is already inactive, just ignore
	if !candidate.Active {
		log.Warn("unregisterWitness witness already inactive.", "address", address.Hex())
		return fmt.Errorf("unregisterWitness witness already inactive.")
	}

	// set candidate active false
	candidate.Active = false

	// save candidate info db
	err := ec.setCandidate(candidate)
	if err != nil {
		log.Error("unregisterWitness setCandidate err.", "address", address.Hex(), "err", err)
		return err
	}

	return nil
}

func (ec electionContext) voteWitnesses(address common.Address, candidates []common.Address) error {
	// 入参校验，如果投的候选人过多，返回错误
	if len(candidates) > voteLimit {
		return fmt.Errorf("you voted too many candidates: the limit is %d, you voted %d", voteLimit, len(candidates))
	}

	voter := ec.getVoter(address)
	var voteCount *big.Int
	var err error

	if voteCount, err = ec.prepareForVote(&voter, address); err != nil {
		return err
	}
	// 计算当前stake可以兑换得到的票数
	voter.LastVoteCount = new(big.Int).Set(voteCount)

	if voter.ProxyVoteCount != nil && voter.ProxyVoteCount.Sign() > 0 {
		voteCount.Add(voteCount, voter.ProxyVoteCount)
	}

	// 逐个检查是否投给了非候选者,并给相应的候选者加上票数
	candiSet := make(map[common.Address]struct{})
	voter.VoteCandidates = nil
	for _, candidate := range candidates {
		if _, ok := candiSet[candidate]; ok {
			continue
		}
		candiSet[candidate] = struct{}{}

		// 如果是候选人则增加相应的选票
		candi := ec.getCandidate(candidate)
		if bytes.Equal(candi.Owner.Bytes(), candidate.Bytes()) && candi.Active {
			voter.VoteCandidates = append(voter.VoteCandidates, candidate)
			candi.VoteCount.Add(candi.VoteCount, voteCount)
			err = ec.setCandidate(candi)
			if err != nil {
				return fmt.Errorf("setCandidate error: %s", err)
			}
		}
	}

	return ec.setVoter(voter)
}

func (ec electionContext) cancelVote(address common.Address) error {
	voter := ec.getVoter(address)
	if !bytes.Equal(voter.Owner.Bytes(), address.Bytes()) {
		return fmt.Errorf("the voter %x doesn't exist", address)
	}
	// 设置了代理，则返回错误，让其取消代理
	if !bytes.Equal(voter.Proxy.Bytes(), emptyAddress.Bytes()) {
		return fmt.Errorf("must cancel proxy first, proxy: %x", voter.Proxy)
	}
	// 投的候选人为空，不需要取消投票，返回
	if len(voter.VoteCandidates) == 0 {
		log.Warn("voteCandidates is nil, need not cancel", "address", address.Hex())
		return nil
	}
	// 减去原候选人得到的投票
	err := ec.subVoteFromCandidates(&voter)
	if err != nil {
		return fmt.Errorf("subVoteFromCandidates error: %s", err)
	}

	// 将上次投票信息置空
	voter.LastVoteCount = big.NewInt(0)
	voter.VoteCandidates = nil

	return ec.setVoter(voter)
}

func (ec electionContext) startProxy(address common.Address) error {
	// get voter from db
	voter := ec.getVoter(address)

	// proxy already in db
	if bytes.Equal(voter.Owner.Bytes(), address.Bytes()) {

		// already registered as proxy
		if voter.IsProxy {
			log.Info("startProxy proxy is already started", "address", address.Hex())
			return fmt.Errorf("startProxy proxy is already started")
		}
		// 已经设置了代理，则不可以成为代理
		if !bytes.Equal(voter.Proxy.Bytes(), emptyAddress.Bytes()) {
			return fmt.Errorf("account that uses a proxy is not allowed to become a proxy")
		}

		// not registered as proxy yet
		voter.IsProxy = true
		// voter.ProxyVoteCount = big.NewInt(0)
	} else {
		// proxy not in db
		voter.Owner = address
		voter.IsProxy = true
	}

	// save voter into db
	err := ec.setVoter(voter)
	if err != nil {
		log.Error("startProxy setVoter err.", "address", address.Hex(), "err", err)
		return err
	}

	return nil
}

func (ec electionContext) stopProxy(address common.Address) error {
	// get voter from db
	voter := ec.getVoter(address)

	// proxy not in db
	if !bytes.Equal(voter.Owner.Bytes(), address.Bytes()) {
		log.Warn("stopProxy proxy does not exist.", "address", address.Hex())
		return fmt.Errorf("stopProxy proxy does not exist.")
	}

	// voter is not a proxy, just ignore
	if !voter.IsProxy {
		log.Warn("stopProxy address is not proxy", "address", address.Hex())
		return fmt.Errorf("stopProxy address is not proxy")
	}

	voter.IsProxy = false
	// voter.ProxyVoteCount = big.NewInt(0)

	// save voter into db
	err := ec.setVoter(voter)
	if err != nil {
		log.Error("stopProxy setVoter err.", "address", address.Hex(), "err", err)
		return err
	}

	return nil
}

func (ec electionContext) setProxy(address common.Address, proxy common.Address) error {
	// 不可以将自身设置为自己的代理
	if bytes.Equal(address.Bytes(), proxy.Bytes()) {
		return fmt.Errorf("cannot proxy to self")
	}

	voter := ec.getVoter(address)
	// 如果自己也是个代理，返回错误
	if voter.IsProxy {
		return fmt.Errorf("account registered as a proxy is not allowed to use a proxy")
	}

	var voteCount *big.Int
	var err error
	// 撤销上次的投票或者设置代理
	if voteCount, err = ec.prepareForVote(&voter, address); err != nil {
		return err
	}
	voter.LastVoteCount = new(big.Int).Set(voteCount)

	if voter.ProxyVoteCount != nil && voter.ProxyVoteCount.Sign() > 0 {
		voteCount.Add(voteCount, voter.ProxyVoteCount)
	}

	var proxyVoter Voter
	proxyVoter = ec.getVoter(proxy)
	if !proxyVoter.IsProxy {
		return fmt.Errorf("%x is not a proxy", proxy)
	}

	// 增加代理人投的票
	proxyVoter.ProxyVoteCount.Add(proxyVoter.ProxyVoteCount, voteCount)
	err = ec.setVoter(proxyVoter)
	if err != nil {
		return fmt.Errorf("setVoter error: %s", err)
	}

	// 找到了最终代理
	if bytes.Equal(proxyVoter.Proxy.Bytes(), emptyAddress.Bytes()) {
		// 把票数加到该代理人投的候选者身上
		if len(proxyVoter.VoteCandidates) > 0 {
			addOp := func(count *big.Int) {
				count.Add(count, voteCount)
			}
			ec.opCandidates(&proxyVoter, addOp)
		}
	}

	voter.VoteCandidates = nil
	voter.Proxy = proxy
	return ec.setVoter(voter)
}

func (ec electionContext) cancelProxy(address common.Address) error {
	voter := ec.getVoter(address)
	if !bytes.Equal(voter.Owner.Bytes(), address.Bytes()) || bytes.Equal(voter.Proxy.Bytes(), emptyAddress.Bytes()) {
		return fmt.Errorf("not set proxy")
	}
	proxy := voter.Proxy
	voteCount := new(big.Int).Set(voter.LastVoteCount)
	if voter.ProxyVoteCount != nil && voter.ProxyVoteCount.Sign() > 0 {
		voteCount.Add(voteCount, voter.ProxyVoteCount)
	}

	var proxyVoter Voter
	for {
		proxyVoter = ec.getVoter(proxy)
		// 减少其代理的票
		proxyVoter.ProxyVoteCount.Sub(proxyVoter.ProxyVoteCount, voteCount)
		err := ec.setVoter(proxyVoter)
		if err != nil {
			return fmt.Errorf("setVoter error: %s", err)
		}

		// 找到了最终代理
		if bytes.Equal(proxyVoter.Proxy.Bytes(), emptyAddress.Bytes()) {
			if len(proxyVoter.VoteCandidates) > 0 {
				subOp := func(count *big.Int) {
					count.Sub(count, voteCount)
				}
				ec.opCandidates(&proxyVoter, subOp)
			}
			break
		}

		proxy = proxyVoter.Proxy
	}

	voter.Proxy = emptyAddress
	voter.LastVoteCount = big.NewInt(0)
	return ec.setVoter(voter)
}

func (ec electionContext) stake(address common.Address, stakeCount *big.Int) error {
	if stakeCount.Sign() <= 0 {
		log.Error("stake stakeCount <= 0", "address", address.Hex(), "stakeCount", stakeCount)
		return fmt.Errorf("stake stakeCount less than 0")
	}

	// get address balance
	balance := ec.context.GetStateDb().GetBalance(address)

	// get the balance that need staking
	balanceNeedStake := big.NewInt(0).Mul(stakeCount, big.NewInt(1e+18))

	// if balance is not enough, just ignore
	if balance.Cmp(balanceNeedStake) < 0 {
		log.Error("stake not enough balance.", "address", address.Hex(), "balance", balance)
		return fmt.Errorf("stake not enough balance.")
	}

	// sub balance of staker
	ec.context.GetStateDb().SubBalance(address, balanceNeedStake)

	// get stake from db
	stake := ec.getStake(address)

	// if stake already exists, just add stakeCount to origin stake
	if bytes.Equal(stake.Owner.Bytes(), address.Bytes()) {
		// add stake of staker
		stake.StakeCount = big.NewInt(0).Add(stake.StakeCount, stakeCount)
	} else {
		// else set StakeCount as @StakeCount
		stake.Owner = address
		stake.StakeCount = stakeCount
	}

	// update last stake time
	stake.TimeStamp = ec.context.GetTime()

	// put stake into db
	err := ec.setStake(stake)
	if err != nil {
		log.Error("stake setStake err.", "address", address.Hex(), "err", err)
		return err
	}

	return nil
}

func (ec electionContext) unStake(address common.Address) error {
	// get stake from db
	stake := ec.getStake(address)

	// if stake is not found in db, just ignore
	if !bytes.Equal(stake.Owner.Bytes(), address.Bytes()) {
		log.Error("unStake stake is not found in db.", "address", address.Hex())
		return fmt.Errorf("unStake stake is not found in db.")
	}

	stakeCount := stake.StakeCount

	// no stake, no need to unstake, just ignore
	if stakeCount.Cmp(big.NewInt(0)) == 0 {
		log.Error("unStake 0 stakeCount.", "address", address.Hex())
		return fmt.Errorf("unStake 0 stakeCount.")
	}

	// get the time point that can unstake
	canUnstakeTime := big.NewInt(0).Add(stake.TimeStamp, unstakePeriod)

	// if time is less than minimum stake period, cannot untake, just ignore
	if ec.context.GetTime().Cmp(canUnstakeTime) < 0 {
		log.Error("cannot unstake in 24 hours", "address", address.Hex())
		return fmt.Errorf("cannot unstake in 24 hours")
	}

	// sub stakeCount of staker
	stake.StakeCount = big.NewInt(0)

	// save stake into db
	err := ec.setStake(stake)
	if err != nil {
		log.Error("unStake setStake err.", "address", address.Hex(), "err", err)
		return err
	}

	// add balance of staker
	ec.context.GetStateDb().AddBalance(address, big.NewInt(0).Mul(stakeCount, big.NewInt(1e+18)))

	return nil
}

func (ec electionContext) extractOwnBounty(addr common.Address) error {
	//24小时内提取1次
	//总激励-已提取激励：是本次可提取的VNT数量，每次至少1000VNT才可提取
	candidate := ec.getCandidate(addr)
	if !bytes.Equal(candidate.Owner.Bytes(), addr.Bytes()) {
		log.Warn("extractOwnBounty unknown witness.", "address", addr.Hex())
		return fmt.Errorf("extractOwnBounty unknown witness.")
	}
	now := ec.context.GetTime()
	if now.Cmp(candidate.LastExtractTime) < 0 || now.Cmp(new(big.Int).Add(candidate.LastExtractTime, big.NewInt(oneDay))) < 0 {
		return fmt.Errorf("it's less than 24h after your last extract bounty,lastExtractTime: %v , now: %v", candidate.LastExtractTime, now)
	}

	restBounty := new(big.Int).Sub(candidate.TotalBounty, candidate.ExtractedBounty)

	if restBounty.Cmp(baseBounty) < 0 {
		log.Warn("the rest of bounty is not enough 1000 vnt", restBounty)
		return fmt.Errorf("the rest of bounty %v wei is not enough 1000 vnt", restBounty)
	}

	candidate.ExtractedBounty.Add(candidate.ExtractedBounty, restBounty)
	candidate.LastExtractTime = now
	if err := ec.setCandidate(candidate); err != nil {
		return fmt.Errorf("set Candidate error %s", err)
	}
	ec.context.GetStateDb().AddBalance(addr, restBounty)
	return nil
}

func (ec electionContext) prepareForVote(voter *Voter, address common.Address) (*big.Int, error) {
	now := ec.context.GetTime()
	stake := ec.getStake(address)
	// 查看当前是否有抵押，无抵押返回无权投票的错误
	if !bytes.Equal(stake.Owner.Bytes(), address.Bytes()) || stake.StakeCount == nil || stake.StakeCount.Sign() <= 0 {
		return nil, fmt.Errorf("you must stake before vote")
	}
	voteCount := ec.calculateVoteCount(stake.StakeCount)
	// 第一次投票
	if !bytes.Equal(voter.Owner.Bytes(), address.Bytes()) {
		voter.Owner = address
		voter.TimeStamp = now
	} else {
		// 如果距离上次投票时间不足24小时，拒绝投票
		if now.Cmp(voter.TimeStamp) < 0 || now.Cmp(new(big.Int).Add(voter.TimeStamp, big.NewInt(oneDay))) < 0 {
			return nil, fmt.Errorf("it's less than 24h after your last vote or setProxy, lastTime: %v, now: %v", voter.TimeStamp, ec.context.GetTime())
		} else {
			voter.TimeStamp = now
		}
		// 如果当前设置了代理，要先取消代理,或者取消之前的投票
		if !bytes.Equal(voter.Proxy.Bytes(), emptyAddress.Bytes()) {
			voter.Proxy = emptyAddress
			return voteCount, ec.cancelProxy(voter.Owner)
		} else {
			// 代理的票数和自身的票数
			return voteCount, ec.subVoteFromCandidates(voter)
		}
	}
	return voteCount, nil
}

func (ec electionContext) subVoteFromCandidates(voter *Voter) error {
	lastVoteCount := new(big.Int).Set(voter.LastVoteCount)
	if voter.ProxyVoteCount != nil && voter.ProxyVoteCount.Sign() > 0 {
		lastVoteCount.Add(lastVoteCount, voter.ProxyVoteCount)
	}
	subOp := func(count *big.Int) {
		count.Sub(count, lastVoteCount)
	}
	return ec.opCandidates(voter, subOp)
}

func (ec electionContext) opCandidates(voter *Voter, opFn func(*big.Int)) error {
	for _, candidate := range voter.VoteCandidates {
		candi := ec.getCandidate(candidate)
		if !bytes.Equal(candi.Owner.Bytes(), candidate.Bytes()) {
			return fmt.Errorf("The candidate %x doesn't exist.", candidate)
		}

		if candi.VoteCount == nil {
			candi.VoteCount = big.NewInt(0)
		}
		// 操作候选人的相应投票
		opFn(candi.VoteCount)
		if candi.VoteCount.Sign() < 0 {
			return fmt.Errorf("the voteCount %v of candidate %x is negative", candi.VoteCount, candi.Owner)
		}
		err := ec.setCandidate(candi)
		if err != nil {
			return fmt.Errorf("setCandidate error: %s", err)
		}
	}
	return nil
}

func (ec electionContext) calculateVoteCount(stakeCount *big.Int) *big.Int {
	deltaTime := big.NewInt(0)
	deltaTime.Sub(ec.context.GetTime(), eraTimeStamp)
	deltaTime.Div(deltaTime, big.NewInt(oneWeek))

	weight := float64(deltaTime.Uint64()) / 52

	votes := float64(stakeCount.Uint64()) * math.Exp2(weight)
	return big.NewInt(int64(votes))
}

// GetFirstNCandidates get candidates with most votes as witness from specific stateDB
func GetFirstNCandidates(stateDB inter.StateDB, witnessesNum int) ([]common.Address, []string) {
	var witnesses []common.Address
	var urls []string
	candidates := getAllCandidate(stateDB)
	if candidates == nil {
		log.Warn("There is no witness candidates. If you want to be a witness, please register now.")
		return nil, nil
	}
	if len(candidates) < witnessesNum {
		log.Warn("Witness candidates is too less. If you want to be a witness, please register now.", "num of candidates", len(candidates), "want", witnessesNum)
		return nil, nil
	}

	candidates.Sort()
	witnessSet := make(map[common.Address]struct{})
	for i := 0; i < len(candidates) && len(witnesses) < witnessesNum; i++ {
		if candidates[i].VoteCount.Cmp(big.NewInt(0)) >= 0 && candidates[i].Active {
			witnesses = append(witnesses, candidates[i].Owner)
			witnessSet[candidates[i].Owner] = struct{}{}
			urls = append(urls, string(candidates[i].Url))
		}
	}
	if len(witnessSet) != witnessesNum {
		log.Warn("Valid witness candidates is too less. If you want to be a witness, please register now.", "num of valid candidates", len(witnessSet), "want", witnessesNum)
		return nil, nil
	}

	return witnesses, urls
}

// GetAllCandidates return the list of all candidate. Candidates will be
// sort by votes and address, if sorted is true.
func GetAllCandidates(stateDB inter.StateDB, sorted bool) CandidateList {
	candidates := getAllCandidate(stateDB)
	if sorted {
		candidates.Sort()
	}
	return candidates
}

// GetVoter returns a voter's information
func GetVoter(stateDB inter.StateDB, addr common.Address) *Voter {
	getFromDB := func(key common.Hash) common.Hash {
		return stateDB.GetState(electionAddr, key)
	}

	v := getVoterFrom(addr, getFromDB)
	return &v
}

// GetStake returns a user's information
func GetStake(stateDB inter.StateDB, addr common.Address) *Stake {
	getFromDB := func(key common.Hash) common.Hash {
		return stateDB.GetState(electionAddr, key)
	}

	s := getStakeFrom(addr, getFromDB)
	return &s
}

func AddCandidatesBounty(stateDB inter.StateDB, bonus map[common.Address]*big.Int) error {
	for addr, bu := range bonus {
		if err := addCandidateBounty(stateDB, addr, bu); err != nil {
			return err
		}
	}
	return nil
}

// GrantBounty grants VNT bounty. Returns an error, if RestTotalBounty is less
// than grantAmount.
func GrantBounty(stateDB inter.StateDB, grantAmount *big.Int) (*big.Int, error) {
	bounty := getRestBounty(stateDB)
	if bounty.RestTotalBounty.Cmp(grantAmount) < 0 {
		return bounty.RestTotalBounty, fmt.Errorf("rest bounty %v is not enough to pay %v", bounty.RestTotalBounty, grantAmount)
	}
	newRestBounty := new(big.Int).Sub(bounty.RestTotalBounty, grantAmount)
	err := setRestBounty(stateDB, Bounty{newRestBounty})
	return newRestBounty, err
}

// QueryRestVNTBounty returns the value of RestTotalBounty.
func QueryRestVNTBounty(stateDB inter.StateDB) *big.Int {
	if !stateDB.Exist(electionAddr) {
		stateDB.SetNonce(electionAddr, 1)
		setRestBounty(stateDB, Bounty{restTotalBounty})
		return restTotalBounty
	}
	bounty := getRestBounty(stateDB)
	return bounty.RestTotalBounty
}
