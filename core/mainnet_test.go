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

package core

import (
	"math/big"
	"testing"

	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/consensus/mock"
	"github.com/vntchain/go-vnt/core/types"
	"github.com/vntchain/go-vnt/core/vm"
	"github.com/vntchain/go-vnt/core/vm/election"
	"github.com/vntchain/go-vnt/crypto"
	"github.com/vntchain/go-vnt/params"
	"github.com/vntchain/go-vnt/vntdb"
)

// TestBanTransaction transaction is not allowed before main net started.
func TestBanTransaction(t *testing.T) {
	var (
		key, _   = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr     = crypto.PubkeyToAddress(key.PublicKey)
		db       = vntdb.NewMemDatabase()
		gspec    = &Genesis{Config: params.TestChainConfig, Alloc: GenesisAlloc{addr: {Balance: big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18))}}}
		genesis  = gspec.MustCommit(db)
		signer   = types.NewHubbleSigner(gspec.Config.ChainID)
		receiver = common.BytesToAddress([]byte{111})
		amount   = big.NewInt(1000)
	)

	// 恢复到未激活状态的标记
	election.ResetActive()

	blockchain, _ := NewBlockChain(db, nil, gspec.Config, mock.NewMock(), vm.Config{})
	defer blockchain.Stop()

	chain, _ := GenerateChain(params.TestChainConfig, genesis, mock.NewMock(), db, 3, func(i int, gen *BlockGen) {
		switch i {
		case 2:
			tx, err := types.SignTx(types.NewTransaction(gen.TxNonce(addr), receiver, amount, 1000000, big.NewInt(18000000000), nil), signer, key)
			if err != nil {
				t.Fatalf("failed to create tx: %v", err)
			}
			gen.AddTx(tx)
		}
	})
	if _, err := blockchain.InsertChain(chain); err != nil {
		t.Fatalf("failed to insert chain: %v", err)
	}

	// 	接收人账号余额应当为0
	stateDb, err := blockchain.State()
	if err != nil {
		t.Fatalf("get state db error: %v", err)
	}
	if bal := stateDb.GetBalance(receiver); bal.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("receiver balance got: %v want: %v", bal.String(), 0)
	}
}

// TestAllowTransaction transaction is allowed after main net started.
func TestAllowTransaction(t *testing.T) {
	var (
		key, _   = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr     = crypto.PubkeyToAddress(key.PublicKey)
		db       = vntdb.NewMemDatabase()
		gspec    = &Genesis{Config: params.TestChainConfig, Alloc: GenesisAlloc{addr: {Balance: big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18))}}}
		genesis  = gspec.MustCommit(db)
		signer   = types.NewHubbleSigner(gspec.Config.ChainID)
		receiver = common.BytesToAddress([]byte{111})
		amount   = big.NewInt(1000)
	)

	// 恢复到未激活状态的标记
	election.ResetActive()

	blockchain, _ := NewBlockChain(db, nil, gspec.Config, mock.NewMock(), vm.Config{})
	defer blockchain.Stop()

	signTx := func(tx *types.Transaction) (*types.Transaction, error) {
		return types.SignTx(tx, signer, key)
	}
	chain, _ := GenerateChain(params.TestChainConfig, genesis, mock.NewMock(), db, 3, func(i int, gen *BlockGen) {
		switch i {
		case 0:
			if err := StartFakeMainNet(gen, addr, signTx); err != nil {
				t.Fatalf(err.Error())
			}
		case 2:
			if !election.MainNetActive(gen.statedb) {
				t.Fatalf("main is inactive")
			}

			tx, err := types.SignTx(types.NewTransaction(gen.TxNonce(addr), receiver, amount, 1000000, big.NewInt(18000000000), nil), signer, key)
			if err != nil {
				t.Fatalf("failed to create tx: %v", err)
			}
			gen.AddTx(tx)
		}
	})
	if _, err := blockchain.InsertChain(chain); err != nil {
		t.Fatalf("failed to insert chain: %v", err)
	}

	// 应当能查询到交易，并且成功，接收人应当有等于转账金额的余额
	stateDb, err := blockchain.State()
	if err != nil {
		t.Fatalf("get state db error: %v", err)
	}
	if bal := stateDb.GetBalance(receiver); bal.Cmp(amount) != 0 {
		t.Errorf("receiver balance got: %v want: %v", bal.String(), 0)
	}
}
