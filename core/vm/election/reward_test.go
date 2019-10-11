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
	"fmt"
	"math/big"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/vntchain/go-vnt/common"
)

type grantCase struct {
	name          string                      // case名称
	balance       *big.Int                    // 合约余额
	allLockBalance *big.Int                    // 锁仓和抵押总额
	cans          CandidateList               // 当前的候选人列表
	rewards       map[common.Address]*big.Int // 待发放余额
	errExpOfGrant error                       // 有error时意味着回滚，也要匹配具体error
	matchTotal    bool                        // 匹配总金额的变化，只有当全部分配或全部未分配时才使用
	// 匹配每个受益账号增加的金额与rewardBalance是否匹配，只有rewards中所有Active受益人能收到激励的测试场景才设置为true
	matchBeneficiary bool
}

func TestGrantBounty(t *testing.T) {
	be := beneficiary.String()
	be = be[:len(be)-1]
	// binder地址实际不使用，可以相同，Beneficiary地址不能相同
	ca1 := Candidate{Owner: addr1, Binder: binder, Beneficiary: common.HexToAddress(be + "1"), Registered: true, Bind: true}
	ca2 := Candidate{Owner: addr2, Binder: binder, Beneficiary: common.HexToAddress(be + "2"), Registered: true, Bind: true}
	ca3 := Candidate{Owner: addr3, Binder: binder, Beneficiary: common.HexToAddress(be + "3"), Registered: true, Bind: false}

	// 	用例1：有充足余额，有充足剩余激励，余额减少，收益人余额增加
	{
		caList := CandidateList{ca1, ca2}
		rewards := map[common.Address]*big.Int{
			ca1.Owner: vnt2wei(1),                                     // 1VNT
			ca2.Owner: big.NewInt(0).Div(vnt2wei(15), big.NewInt(10))} // 1.5VNT
		cas := grantCase{"case1", vnt2wei(100), vnt2wei(90), caList, rewards, nil, true, true}
		testGrantBounty(t, &cas)
	}

	// 	用例2：有充足余额，剩余激励够1个账号的，不足2个账号，有账号的受益人余额增加，有的未增加，剩余激励为0
	{
		caList := CandidateList{ca1, ca2}
		rewards := map[common.Address]*big.Int{
			ca1.Owner: vnt2wei(10), // 10VNT
			ca2.Owner: vnt2wei(20)} // 20VNT
		cas := grantCase{"case2", vnt2wei(100), vnt2wei(75), caList, rewards, nil, false, false}
		testGrantBounty(t, &cas)
	}

	// 	用例3：有充足余额，有充足剩余激励，候选人有1个未激活，被跳过，检查剩余激励正确
	{
		caList := CandidateList{ca1, ca2, ca3}
		rewards := map[common.Address]*big.Int{
			ca1.Owner: vnt2wei(10), // 10VNT
			ca2.Owner: vnt2wei(20), // 20VNT
			ca3.Owner: vnt2wei(20)} // 20VNT，非法情况，正常调用不会存在非Active的候选人分激励
		cas := grantCase{"case3", vnt2wei(100), vnt2wei(50), caList, rewards, nil, false, true}
		testGrantBounty(t, &cas)
	}
}

func testGrantBounty(t *testing.T, cas *grantCase) {
	ec := newTestElectionCtx()
	db := ec.context.GetStateDb()

	// 设置余额
	db.AddBalance(contractAddr, cas.balance)

	// 设置alllock amount
	err := setLock(db, AllLock{cas.allLockBalance})
	assert.Equal(t, err, nil, fmt.Sprintf("%v, set alllock amount error: %v", cas.name, err))

	restReward := QueryRestReward(db)
	expRestReward := common.Big0
	if cas.balance.Cmp(cas.allLockBalance) > 0 {
		expRestReward = big.NewInt(0).Sub(cas.balance, cas.allLockBalance)
	}
	assert.Equal(t, restReward, expRestReward, fmt.Sprintf("%v, rest bounty amount error: %v", cas.name, err))


	// 设置候选人
	for i, can := range cas.cans {
		err = ec.setCandidate(can)
		assert.Equal(t, err, nil, fmt.Sprintf("%v, [%d] set candidate error: %v", cas.name, i, err))
	}

	// 执行分激励
	err = GrantReward(db, cas.rewards)
	assert.Equal(t, err, cas.errExpOfGrant, fmt.Sprintf("%v, grant bounty error mismatch", cas.name))

	// 校验回滚
	if cas.errExpOfGrant != nil {
		// 校验余额，应当不变
		assert.Equal(t, db.GetBalance(contractAddr), cas.balance, ",", cas.name, ", reverted, balance of contract should not change")
		// 校验剩余激励，应当不变
		acLockAmount, _ := getLock(db)
		assert.Equal(t, acLockAmount.Amount, cas.allLockBalance, ",", cas.name, ", reverted, left reward should not change")
		// 	各候选人账号的收益账号应当为0
		for addr, _ := range cas.rewards {
			can := ec.getCandidate(addr)
			assert.Equal(t, db.GetBalance(can.Beneficiary), common.Big0, ",", cas.name, ", reverted, balance of beneficiary should be 0")
		}
	}

	// 校验余额和受益人余额增加额是否匹配，剩余激励是否匹配
	totalReward := big.NewInt(0)
	for _, re := range cas.rewards {
		totalReward = totalReward.Add(totalReward, re)
	}

	reminReward := QueryRestReward(db)
	reducedBalance := big.NewInt(0).Sub(cas.balance, db.GetBalance(contractAddr))
	reducedReward := big.NewInt(0).Sub(restReward, reminReward)
	assert.Equal(t, reducedBalance, reducedReward, ",", cas.name, "reduced balance should always equal to reduces reward")
	if cas.matchTotal {
		assert.Equal(t, reducedBalance, totalReward, ",", cas.name, "reduced contract balance should equal total reward")
		assert.Equal(t, reducedReward, totalReward, ",", cas.name, "reduced contract reward should equal total reward")
	}

	// 	校验每个受益账户增加的余额
	if cas.matchBeneficiary {
		// 统计受益账户应收到的总激励
		benes := make(map[common.Address]*big.Int)
		for addr, amount := range cas.rewards {
			can := ec.getCandidate(addr)
			assert.Equal(t, can.Owner, addr, "%v, candidate address not match", cas.name)
			if can.Active() {
				if _, ok := benes[can.Beneficiary]; ok {
					benes[can.Beneficiary] = big.NewInt(0).Add(benes[can.Beneficiary], amount)
				} else {
					benes[can.Beneficiary] = big.NewInt(0).Set(amount)
				}
			}
		}

		// 	匹配实际收到金额
		for addr, reward := range benes {
			// 所有受益账户初始都没有余额
			addedBal := db.GetBalance(addr)
			assert.Equal(t, addedBal, reward, fmt.Sprintf("%v, beneficiary reward mismtach: %v", cas.name, addr.String()))
		}
	}
}

func TestDepositReward(t *testing.T) {
	ec := newTestElectionCtx()
	sender := common.HexToAddress("0x123456")

	// 初始值应当为0
	got, _:= getLock(ec.context.GetStateDb())
	assert.Equal(t, got.Amount, common.Big0, "rest reward should be 0 at first")

	// 存负值
	amount := vnt2wei(-1000)
	err := ec.depositReward(sender, amount)
	assert.Equal(t, err, fmt.Errorf("deposit reward less than 0 VNT"), fmt.Sprintf("deposit reward error: %v", err))

	// 得1000VNT
	amount = vnt2wei(1000)
	err = ec.depositReward(sender, amount)
	assert.Equal(t, err, nil, fmt.Sprintf("deposit reward error: %v", err))
	got, _ = getLock(ec.context.GetStateDb())
	assert.Equal(t, got.Amount, common.Big0, "deposit reward not equal")
}
