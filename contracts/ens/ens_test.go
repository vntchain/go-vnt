// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package ens

import (
	"context"
	"math/big"
	"testing"

	"github.com/vntchain/go-vnt/accounts/abi/bind"
	"github.com/vntchain/go-vnt/accounts/abi/bind/backends"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/contracts/ens/contract"
	"github.com/vntchain/go-vnt/core"
	"github.com/vntchain/go-vnt/core/types"
	"github.com/vntchain/go-vnt/core/vm/election"
	"github.com/vntchain/go-vnt/crypto"
	"github.com/vntchain/go-vnt/params"
)

var (
	chainID      = params.TestChainConfig.ChainID
	key, _       = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	name         = "my name on ENS"
	hash         = crypto.Keccak256Hash([]byte("my content")).Hex()
	addr         = crypto.PubkeyToAddress(key.PublicKey)
	testAddr     = common.HexToAddress("0x1234123412341234123412341234123412341234")
	activeKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f292")
	activeAddr   = crypto.PubkeyToAddress(activeKey.PublicKey)
)

func mainnetActive(backend bind.ContractBackend) ([]*types.Transaction, error) {
	nonce, err := backend.PendingNonceAt(context.Background(), activeAddr)
	if err != nil {
		return nil, err
	}
	txs, err := election.GenFakeStartedTxs(nonce, []common.Address{activeAddr})
	return txs, err
}

func TestENS(t *testing.T) {
	contractBackend := backends.NewSimulatedBackend(core.GenesisAlloc{addr: {Balance: big.NewInt(1000000000)}, activeAddr: {Balance: big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18))}})
	activeTxs, err := mainnetActive(contractBackend)
	if err != nil {
		t.Fatalf("can't create active tx: %v", err)
	}
	for _, v := range activeTxs {
		signTx, err := types.SignTx(v, types.NewHubbleSigner(chainID), activeKey)
		if err != nil {
			t.Fatalf("sign tx error: %v", err)
		}
		err = contractBackend.SendTransaction(context.Background(), signTx)
		if err != nil {
			t.Fatalf("can't send active tx: %v", err)
		}

		contractBackend.Commit()
	}

	transactOpts := bind.NewKeyedTransactor(key, chainID)

	ensAddr, ens, err := DeployENS(transactOpts, contractBackend)
	if err != nil {
		t.Fatalf("can't deploy root registry: %v", err)
	}
	contractBackend.Commit()

	// Set ourself as the owner of the name.
	if _, err := ens.Register(name); err != nil {
		t.Fatalf("can't register: %v", err)
	}
	contractBackend.Commit()

	// Deploy a resolver and make it responsible for the name.
	resolverAddr, _, _, err := contract.DeployPublicResolver(transactOpts, contractBackend, ensAddr)
	if err != nil {
		t.Fatalf("can't deploy resolver: %v", err)
	}
	if _, err := ens.SetResolver(EnsNode(name), resolverAddr); err != nil {
		t.Fatalf("can't set resolver: %v", err)
	}
	contractBackend.Commit()

	// Set the content hash for the name.
	if _, err = ens.SetContentHash(name, hash); err != nil {
		t.Fatalf("can't set content hash: %v", err)
	}
	contractBackend.Commit()

	// Try to resolve the name.
	vhost, err := ens.Resolve(name)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if vhost != hash {
		t.Fatalf("resolve error, expected %v, got %v", hash, vhost)
	}

	// set the address for the name
	if _, err = ens.SetAddr(name, testAddr); err != nil {
		t.Fatalf("can't set address: %v", err)
	}
	contractBackend.Commit()

	// Try to resolve the name to an address
	recoveredAddr, err := ens.Addr(name)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if vhost != hash {
		t.Fatalf("resolve error, expected %v, got %v", testAddr.Hex(), recoveredAddr.Hex())
	}
}
