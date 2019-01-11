// Copyright 2015 The go-ethereum Authors
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

package vnt

import (
	"context"
	"math/big"

	"github.com/vntchain/go-vnt/accounts"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/common/math"
	"github.com/vntchain/go-vnt/core"
	"github.com/vntchain/go-vnt/core/bloombits"
	"github.com/vntchain/go-vnt/core/rawdb"
	"github.com/vntchain/go-vnt/core/state"
	"github.com/vntchain/go-vnt/core/types"
	"github.com/vntchain/go-vnt/core/vm"
	"github.com/vntchain/go-vnt/event"
	"github.com/vntchain/go-vnt/params"
	"github.com/vntchain/go-vnt/rpc"
	"github.com/vntchain/go-vnt/vnt/downloader"
	"github.com/vntchain/go-vnt/vnt/gasprice"
	"github.com/vntchain/go-vnt/vntdb"
)

// VntAPIBackend implements vntapi.Backend for full nodes
type VntAPIBackend struct {
	vnt *VNT
	gpo *gasprice.Oracle
}

func (b *VntAPIBackend) ChainConfig() *params.ChainConfig {
	return b.vnt.chainConfig
}

func (b *VntAPIBackend) CurrentBlock() *types.Block {
	return b.vnt.blockchain.CurrentBlock()
}

func (b *VntAPIBackend) SetHead(number uint64) {
	b.vnt.protocolManager.downloader.Cancel()
	b.vnt.blockchain.SetHead(number)
}

func (b *VntAPIBackend) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.vnt.miner.PendingBlock()
		return block.Header(), nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.vnt.blockchain.CurrentBlock().Header(), nil
	}
	return b.vnt.blockchain.GetHeaderByNumber(uint64(blockNr)), nil
}

func (b *VntAPIBackend) BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Block, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.vnt.miner.PendingBlock()
		return block, nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.vnt.blockchain.CurrentBlock(), nil
	}
	return b.vnt.blockchain.GetBlockByNumber(uint64(blockNr)), nil
}

func (b *VntAPIBackend) StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *types.Header, error) {
	// Pending state is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block, state := b.vnt.miner.Pending()
		return state, block.Header(), nil
	}
	// Otherwise resolve the block number and return its state
	header, err := b.HeaderByNumber(ctx, blockNr)
	if header == nil || err != nil {
		return nil, nil, err
	}
	stateDb, err := b.vnt.BlockChain().StateAt(header.Root)
	return stateDb, header, err
}

func (b *VntAPIBackend) GetBlock(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return b.vnt.blockchain.GetBlockByHash(hash), nil
}

func (b *VntAPIBackend) GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts, error) {
	if number := rawdb.ReadHeaderNumber(b.vnt.chainDb, hash); number != nil {
		return rawdb.ReadReceipts(b.vnt.chainDb, hash, *number), nil
	}
	return nil, nil
}

func (b *VntAPIBackend) GetLogs(ctx context.Context, hash common.Hash) ([][]*types.Log, error) {
	number := rawdb.ReadHeaderNumber(b.vnt.chainDb, hash)
	if number == nil {
		return nil, nil
	}
	receipts := rawdb.ReadReceipts(b.vnt.chainDb, hash, *number)
	if receipts == nil {
		return nil, nil
	}
	logs := make([][]*types.Log, len(receipts))
	for i, receipt := range receipts {
		logs[i] = receipt.Logs
	}
	return logs, nil
}

func (b *VntAPIBackend) GetTd(blockHash common.Hash) *big.Int {
	return b.vnt.blockchain.GetTdByHash(blockHash)
}

func (b *VntAPIBackend) GetVM(ctx context.Context, msg core.Message, state *state.StateDB, header *types.Header, vmCfg vm.Config) (vm.VM, func() error, error) {
	state.SetBalance(msg.From(), math.MaxBig256)
	vmError := func() error { return nil }

	context := core.NewVMContext(msg, header, b.vnt.BlockChain(), nil)
	vmGet := core.GetVM(msg, context, state, b.vnt.chainConfig, vmCfg)
	return vmGet, vmError, nil
}

func (b *VntAPIBackend) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	return b.vnt.BlockChain().SubscribeRemovedLogsEvent(ch)
}

func (b *VntAPIBackend) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	return b.vnt.BlockChain().SubscribeChainEvent(ch)
}

func (b *VntAPIBackend) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return b.vnt.BlockChain().SubscribeChainHeadEvent(ch)
}

func (b *VntAPIBackend) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
	return b.vnt.BlockChain().SubscribeChainSideEvent(ch)
}

func (b *VntAPIBackend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return b.vnt.BlockChain().SubscribeLogsEvent(ch)
}

func (b *VntAPIBackend) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	return b.vnt.txPool.AddLocal(signedTx)
}

func (b *VntAPIBackend) GetPoolTransactions() (types.Transactions, error) {
	pending, err := b.vnt.txPool.Pending()
	if err != nil {
		return nil, err
	}
	var txs types.Transactions
	for _, batch := range pending {
		txs = append(txs, batch...)
	}
	return txs, nil
}

func (b *VntAPIBackend) GetPoolTransaction(hash common.Hash) *types.Transaction {
	return b.vnt.txPool.Get(hash)
}

func (b *VntAPIBackend) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	return b.vnt.txPool.State().GetNonce(addr), nil
}

func (b *VntAPIBackend) Stats() (pending int, queued int) {
	return b.vnt.txPool.Stats()
}

func (b *VntAPIBackend) TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	return b.vnt.TxPool().Content()
}

func (b *VntAPIBackend) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) event.Subscription {
	return b.vnt.TxPool().SubscribeNewTxsEvent(ch)
}

func (b *VntAPIBackend) Downloader() *downloader.Downloader {
	return b.vnt.Downloader()
}

func (b *VntAPIBackend) ProtocolVersion() int {
	return b.vnt.EthVersion()
}

func (b *VntAPIBackend) SuggestPrice(ctx context.Context) (*big.Int, error) {
	return b.gpo.SuggestPrice(ctx)
}

func (b *VntAPIBackend) ChainDb() vntdb.Database {
	return b.vnt.ChainDb()
}

func (b *VntAPIBackend) EventMux() *event.TypeMux {
	return b.vnt.EventMux()
}

func (b *VntAPIBackend) AccountManager() *accounts.Manager {
	return b.vnt.AccountManager()
}

func (b *VntAPIBackend) BloomStatus() (uint64, uint64) {
	sections, _, _ := b.vnt.bloomIndexer.Sections()
	return params.BloomBitsBlocks, sections
}

func (b *VntAPIBackend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	for i := 0; i < bloomFilterThreads; i++ {
		go session.Multiplex(bloomRetrievalBatch, bloomRetrievalWait, b.vnt.bloomRequests)
	}
}
