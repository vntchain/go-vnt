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
	"github.com/vntchain/go-vnt/core/types"
	inter "github.com/vntchain/go-vnt/core/vm/interface"
)

// MainNetActive returns whether the main net is started.
func MainNetActive(stateDB inter.StateDB) bool {
	if !mainActive {
		mv := getMainNetVotes(stateDB)
		if mv.Active {
			mainActive = true
		}
	}

	return mainActive
}

// GetMainNetVotes return a pointer of main net vote information.
func GetMainNetVotes(stateDB inter.StateDB) *MainNetVotes {
	mv := getMainNetVotes(stateDB)
	return &mv
}

// modifyMainNetVotes modify the votes of main net and judge whether
// the main net match start condition.
func modifyMainNetVotes(stateDB inter.StateDB, num *big.Int, add bool) error {
	mv := getMainNetVotes(stateDB)
	if add {
		mv.VoteStake = big.NewInt(0).Add(mv.VoteStake, num)
	} else {
		mv.VoteStake = big.NewInt(0).Sub(mv.VoteStake, num)
	}

	// 判断是否激活，并且只执行1次
	if !mv.Active && mv.VoteStake.Cmp(big.NewInt(5e8)) >= 0 {
		mv.Active = true
	}

	return setMainNetVotes(stateDB, mv)
}

// GenFakeStartedTxs generate 3 fake transaction to start the main net.
func GenFakeStartedTxs(nextNonce uint64, witness []common.Address) ([]*types.Transaction, error) {
	electAbi, err := GetElectionABI()
	if err != nil {
		return nil, err
	}
	// 抵押交易
	txData, err := PackInput(electAbi, "stake", big.NewInt(6e8))
	if err != nil {
		return nil, fmt.Errorf("stake error: %v", err)
	}
	stakeTx := types.NewTransaction(nextNonce, common.HexToAddress(ContractAddr), common.Big0, 30000, big.NewInt(18000000000), txData)
	// 注册交易
	txData, err = PackInput(electAbi, "registerWitness", []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHGq5zZFRW5FBJ9YMbbvSiW4xzGg5CKMCtDeg6FNnjCbGS"),
		[]byte("www.vnt.com"), []byte("mocktestnode"))
	if err != nil {
		return nil, fmt.Errorf("reg error: %v", err)
	}
	regTx := types.NewTransaction(nextNonce+1, common.HexToAddress(ContractAddr), common.Big0, 30000, big.NewInt(18000000000), txData)
	// 投票交易
	txData, err = PackInput(electAbi, "voteWitnesses", witness)
	if err != nil {
		return nil, fmt.Errorf("vote error: %v", err)
	}
	voteTx := types.NewTransaction(nextNonce+2, common.HexToAddress(ContractAddr), common.Big0, 30000, big.NewInt(18000000000), txData)

	return []*types.Transaction{stakeTx, regTx, voteTx}, nil
}

// ResetActive is reset the main net state in memory.
// Only used for tests.
func ResetActive() {
	mainActive = false
}
