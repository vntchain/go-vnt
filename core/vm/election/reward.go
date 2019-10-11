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

	"github.com/vntchain/go-vnt/common"
	inter "github.com/vntchain/go-vnt/core/vm/interface"
	"github.com/vntchain/go-vnt/log"
)

type AllLock struct {
	Amount *big.Int // 锁仓和抵押总额
}

func (ec electionContext) depositReward(address common.Address, value *big.Int) error {
	if value.Cmp(common.Big0) <= 0 {
		log.Error("depositReward less than 0 VNT", "address", address.Hex(), "VNT", value.String())
		return fmt.Errorf("deposit reward less than 0 VNT")
	}
	return nil
}

// GrantReward 发放激励给该候选节点的受益人，返回错误。
// 发放激励的接口不区分是产块激励还是投票激励，超级节点必须是Active，否则无收益。
// 激励金额不足发放时为正常情况不返回error，返回nil。
// 返回错误时，数据状态恢复到原始情况，即所有激励都不发放。
func GrantReward(stateDB inter.StateDB, rewards map[common.Address]*big.Int) (err error) {
	// 无激励即可返回
	rest := QueryRestReward(stateDB)
	if rest.Cmp(common.Big0) <= 0 {
		return nil
	}

	// 退出时，如果存在错误，恢复原始状态
	snap := stateDB.Snapshot()
	defer func() {
		if err != nil {
			stateDB.RevertToSnapshot(snap)
		}
	}()

	for addr, amount := range rewards {
		// 激励不能超过剩余金额
		if rest.Cmp(amount) < 0 {
			amount = rest
		}
		can := GetCandidate(stateDB, addr)
		// 再检查：跳过不存在或未激活的候选人
		if can == nil || !can.Active() {
			log.Error("Not find candidate or inactive when granting reward", "addr", addr.String())
			continue
		}
		// 发送错误退出
		if err = transfer(stateDB, contractAddr, can.Beneficiary, amount); err != nil {
			return err
		}
		rest = rest.Sub(rest, amount)
		// 发放到无剩余激励
		if rest.Cmp(common.Big0) <= 0 {
			break
		}
	}

	// 激励正常发放完毕
	return nil
}

// QueryRestReward returns the value of left reward for candidates.
func QueryRestReward(stateDB inter.StateDB) *big.Int {
	totalLock, err := getLock(stateDB)
	if err != nil {
		log.Error("QueryRestReward failed", "err", err)
		return common.Big0;
	}
	totalBalance := stateDB.GetBalance(contractAddr)
	if rest := big.NewInt(0).Sub(totalBalance, totalLock.Amount); rest.Cmp(common.Big0) > 0 {
		return rest;
	} else {
		return common.Big0;
	}
}
