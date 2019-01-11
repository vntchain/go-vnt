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

package miner

import (
	"bytes"
	"encoding/json"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/consensus"
	"github.com/vntchain/go-vnt/consensus/dpos"
	"github.com/vntchain/go-vnt/core"
	"github.com/vntchain/go-vnt/core/state"
	"github.com/vntchain/go-vnt/core/types"
	"github.com/vntchain/go-vnt/core/vm"
	"github.com/vntchain/go-vnt/event"
	"github.com/vntchain/go-vnt/log"
	"github.com/vntchain/go-vnt/params"
	"github.com/vntchain/go-vnt/vntdb"
)

const (
	resultQueueSize  = 10
	miningLogAtDepth = 5

	// txChanSize is the size of channel listening to NewTxsEvent.
	// The number is referenced from the size of tx pool.
	txChanSize = 4096
	// chainHeadChanSize is the size of channel listening to ChainHeadEvent.
	chainHeadChanSize = 10
	// chainSideChanSize is the size of channel listening to ChainSideEvent.
	chainSideChanSize = 10
)

// Agent can register themself with the worker
type Agent interface {
	Work() chan<- *Work
	SetReturnCh(chan<- *Result)
	Stop()
	Start()
}

// Work is the workers current environment and holds
// all of the current state information
type Work struct {
	config *params.ChainConfig
	signer types.Signer

	state   *state.StateDB // apply state changes here
	tcount  int            // tx count in cycle
	gasPool *core.GasPool  // available gas used to pack transactions

	Block *types.Block // the new block

	header   *types.Header
	txs      []*types.Transaction
	receipts []*types.Receipt

	createdAt time.Time
}

type Result struct {
	Work  *Work
	Block *types.Block
}

// worker is the main object which takes care of applying messages to the new state
type worker struct {
	config *params.ChainConfig
	engine consensus.Engine

	mu sync.Mutex

	// update loop
	mux          *event.TypeMux
	txsCh        chan core.NewTxsEvent
	txsSub       event.Subscription
	chainHeadCh  chan core.ChainHeadEvent
	chainHeadSub event.Subscription
	chainSideCh  chan core.ChainSideEvent
	chainSideSub event.Subscription
	recBftMsgSub *event.TypeMuxSubscription
	wg           sync.WaitGroup

	agents map[Agent]struct{}
	recv   chan *Result

	vnt     Backend
	chain   *core.BlockChain
	proc    core.Validator
	chainDb vntdb.Database

	coinbase common.Address
	extra    []byte

	currentMu sync.Mutex
	current   *Work

	snapshotMu    sync.RWMutex
	snapshotBlock *types.Block
	snapshotState *state.StateDB

	unconfirmed *unconfirmedBlocks // set of locally mined blocks pending canonicalness confirmations

	// atomic status counters
	mining int32
	atWork int32

	roundTimer      *time.Timer // Timer to trigger each round of producing block
	resetTimerEvent chan *big.Int
	minerStop       chan struct{}
}

func newWorker(config *params.ChainConfig, engine consensus.Engine, coinbase common.Address, vnt Backend, mux *event.TypeMux) *worker {
	worker := &worker{
		config:          config,
		engine:          engine,
		vnt:             vnt,
		mux:             mux,
		txsCh:           make(chan core.NewTxsEvent, txChanSize),
		chainHeadCh:     make(chan core.ChainHeadEvent, chainHeadChanSize),
		chainSideCh:     make(chan core.ChainSideEvent, chainSideChanSize),
		chainDb:         vnt.ChainDb(),
		recv:            make(chan *Result, resultQueueSize),
		chain:           vnt.BlockChain(),
		proc:            vnt.BlockChain().Validator(),
		coinbase:        coinbase,
		agents:          make(map[Agent]struct{}),
		unconfirmed:     newUnconfirmedBlocks(vnt.BlockChain(), miningLogAtDepth),
		roundTimer:      time.NewTimer(time.Second),
		resetTimerEvent: make(chan *big.Int, 1),
		minerStop:       make(chan struct{}, 1),
	}
	worker.stopRoundTimer()

	// Subscribe NewTxsEvent for tx pool
	worker.txsSub = vnt.TxPool().SubscribeNewTxsEvent(worker.txsCh)
	// Subscribe events for blockchain
	worker.chainHeadSub = vnt.BlockChain().SubscribeChainHeadEvent(worker.chainHeadCh)
	worker.chainSideSub = vnt.BlockChain().SubscribeChainSideEvent(worker.chainSideCh)
	worker.recBftMsgSub = worker.mux.Subscribe(core.RecBftMsgEvent{})

	go worker.recBftMsg()
	go worker.update()
	go worker.wait()

	worker.commitNewWork()

	return worker
}

func (self *worker) setEtherbase(addr common.Address) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.coinbase = addr
}

func (self *worker) setExtra(extra []byte) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.extra = extra
}

func (self *worker) pending() (*types.Block, *state.StateDB) {
	if atomic.LoadInt32(&self.mining) == 0 {
		// return a snapshot to avoid contention on currentMu mutex
		self.snapshotMu.RLock()
		defer self.snapshotMu.RUnlock()
		return self.snapshotBlock, self.snapshotState.Copy()
	}

	self.currentMu.Lock()
	defer self.currentMu.Unlock()
	return self.current.Block, self.current.state.Copy()
}

func (self *worker) pendingBlock() *types.Block {
	if atomic.LoadInt32(&self.mining) == 0 {
		// return a snapshot to avoid contention on currentMu mutex
		self.snapshotMu.RLock()
		defer self.snapshotMu.RUnlock()
		return self.snapshotBlock
	}

	self.currentMu.Lock()
	defer self.currentMu.Unlock()
	return self.current.Block
}

func (self *worker) start() {
	self.mu.Lock()
	defer self.mu.Unlock()

	atomic.StoreInt32(&self.mining, 1)

	// Init bft
	if dp, ok := self.engine.(*dpos.Dpos); ok {
		dp.InitBft(self.SendBftMsg, self.SendBftPeerChangeMsg, self.chain.VerifyBlock, self.writeBlock)
		// 刚启动节点的bft节点设置
		currentRoot := self.chain.CurrentHeader().Root
		witnessesUrl := self.chain.Config().Dpos.WitnessesUrl
		if db, err := self.chain.StateAt(currentRoot); err != nil {
			log.Error("get current db error", "err", err)
		} else {
			_, urls := dp.GetWitnessesFromStateDB(db)
			if len(urls) > 0 {
				witnessesUrl = urls
			}
		}
		self.SendBftPeerChangeMsg(witnessesUrl)
	}

	// spin up agents
	for agent := range self.agents {
		agent.Start()
	}
}

func (self *worker) stop() {
	self.wg.Wait()

	self.minerStop <- struct{}{}

	self.mu.Lock()
	defer self.mu.Unlock()
	if atomic.LoadInt32(&self.mining) == 1 {
		for agent := range self.agents {
			agent.Stop()
		}
	}
	atomic.StoreInt32(&self.mining, 0)
	atomic.StoreInt32(&self.atWork, 0)

	if dp, ok := self.engine.(*dpos.Dpos); ok {
		dp.MiningStop()
	}
}

func (self *worker) register(agent Agent) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.agents[agent] = struct{}{}
	agent.SetReturnCh(self.recv)
}

func (self *worker) unregister(agent Agent) {
	self.mu.Lock()
	defer self.mu.Unlock()
	delete(self.agents, agent)
	agent.Stop()
}

func (self *worker) update() {
	defer self.txsSub.Unsubscribe()
	defer self.chainHeadSub.Unsubscribe()
	defer self.chainSideSub.Unsubscribe()
	defer self.recBftMsgSub.Unsubscribe()

	for {
		// A real event arrived, process interesting content
		select {
		// Handle ChainHeadEvent
		case headEvent := <-self.chainHeadCh:
			log.Debug("Worker: new block write finished", "block hash", headEvent.Block.Hash())
			if self.config.Dpos != nil {
				if dp, ok := self.engine.(*dpos.Dpos); ok {
					dp.CleanOldMsg(headEvent.Block.Number())
				}
			}

			// Handle ChainSideEvent
		case ev := <-self.chainSideCh:
			log.Info("Block fail to on chain", "block hash", ev.Block.Hash())

			// Handle NewTxsEvent
		case ev := <-self.txsCh:
			// Apply transactions to the pending state if we're not mining.
			//
			// Note all transactions received may not be continuous with transactions
			// already included in the current mining block. These transactions will
			// be automatically eliminated.
			if self.config.Dpos == nil && atomic.LoadInt32(&self.mining) == 0 {
				self.currentMu.Lock()
				txs := make(map[common.Address]types.Transactions)
				for _, tx := range ev.Txs {
					acc, _ := types.Sender(self.current.signer, tx)
					txs[acc] = append(txs[acc], tx)
				}
				txset := types.NewTransactionsByPriceAndNonce(self.current.signer, txs)
				self.current.commitTransactions(self.mux, txset, self.chain, self.coinbase)
				self.updateSnapshot()
				self.currentMu.Unlock()
			}

			// Only Dpos using
		case <-self.roundTimer.C:
			if self.config.Dpos != nil {
				self.commitNewWork()
			}

		case nextRoundTime := <-self.resetTimerEvent:
			self.resetRoundTimer(nextRoundTime)

		case <-self.minerStop:
			self.stopRoundTimer()

			// System stopped
		case <-self.txsSub.Err():
			self.stopRoundTimer()
			return
		case <-self.chainHeadSub.Err():
			self.stopRoundTimer()
			return
		case <-self.chainSideSub.Err():
			self.stopRoundTimer()
			return
		}
	}
}

func (self *worker) wait() {
	for {
		for result := range self.recv {
			atomic.AddInt32(&self.atWork, -1)

			if result == nil {
				continue
			}
			block := result.Block
			work := result.Work

			// Update the block hash in all logs since it is now available and not when the
			// receipt/log of individual transactions were created.
			for _, r := range work.receipts {
				for _, l := range r.Logs {
					l.BlockHash = block.Hash()
				}
			}
			for _, log := range work.state.Logs() {
				log.BlockHash = block.Hash()
			}

			stat, err := self.chain.WriteBlockWithState(block, work.receipts, work.state)
			if err != nil {
				log.Error("Failed writing block to chain", "err", err)
				continue
			}

			// Broadcast the block and announce chain insertion event
			self.mux.Post(core.NewMinedBlockEvent{Block: block})

			var (
				events []interface{}
				logs   = work.state.Logs()
			)
			events = append(events, core.ChainEvent{Block: block, Hash: block.Hash(), Logs: logs})
			if stat == core.CanonStatTy {
				events = append(events, core.ChainHeadEvent{Block: block})
			}
			self.chain.PostChainEvents(events, logs)

			// Insert the block into the set of pending ones to wait for confirmations
			self.unconfirmed.Insert(block.NumberU64(), block.Hash())
		}
	}
}

func (self *worker) recBftMsg() {
	for obj := range self.recBftMsgSub.Chan() {
		switch ev := obj.Data.(type) {
		case core.RecBftMsgEvent:
			self.engine.HandleBftMsg(self.chain, ev.BftMsg.Msg)
		default:
			log.Warn("Receive bft msg, but type unknown")
		}
	}
}

// push sends a new work task to currently live miner agents.
func (self *worker) push(work *Work) {
	if atomic.LoadInt32(&self.mining) != 1 {
		return
	}
	for agent := range self.agents {
		atomic.AddInt32(&self.atWork, 1)
		if ch := agent.Work(); ch != nil {
			ch <- work
		}
	}
}

func (self *worker) SendBftMsg(msg types.ConsensusMsg) {
	self.mux.Post(core.SendBftMsgEvent{
		BftMsg: types.BftMsg{
			BftType: msg.Type(),
			Msg:     msg,
		}})
}

func (self *worker) SendBftPeerChangeMsg(urls []string) {
	self.mux.Post(core.BftPeerChangeEvent{
		Urls: urls,
	})
}

// writeBlock write block to block chain, and post NewMinedBlockEvent
func (self *worker) writeBlock(block *types.Block) error {
	if err := self.chain.WriteBlock(block); err != nil {
		log.Error("Failed writing block to chain", "err", err)
		return err
	}

	// Broadcast the block and announce chain insertion event
	self.mux.Post(core.NewMinedBlockEvent{Block: block})
	return nil
}

// makeCurrent creates a new environment for the current cycle.
func (self *worker) makeCurrent(parent *types.Block, header *types.Header) error {
	log.Trace("Enter make current")

	state, err := self.chain.StateAt(parent.Root())
	if err != nil {
		return err
	}
	work := &Work{
		config:    self.config,
		signer:    types.NewEIP155Signer(self.config.ChainID),
		state:     state,
		header:    header,
		createdAt: time.Now(),
	}

	// Keep track of transactions which return errors so they can be removed
	work.tcount = 0
	self.current = work

	return nil
}

func (self *worker) commitNewWork() {
	log.Trace("commitNewWork start")

	self.mu.Lock()
	defer self.mu.Unlock()
	self.currentMu.Lock()
	defer self.currentMu.Unlock()

	tstart := time.Now()
	tstamp := tstart.Unix()
	parent := self.chain.CurrentBlock()

	// Do not work too try before parent block
	wait := time.Unix(parent.Time().Int64(), 0).Sub(tstart)
	if wait > 0 {
		log.Info("CommitNewWork start before parent's block time, wait", "wait", wait)
		time.Sleep(wait)
	}

	num := parent.Number()
	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     num.Add(num, common.Big1),
		GasLimit:   core.CalcGasLimit(parent),
		Extra:      self.extra,
		Time:       big.NewInt(tstamp),
	}
	// Only set the coinbase if we are mining (avoid spurious block rewards)
	if atomic.LoadInt32(&self.mining) == 1 {
		header.Coinbase = self.coinbase
	}

	preErr := self.engine.Prepare(self.chain, header)
	if atomic.LoadInt32(&self.mining) == 1 {
		self.resetTimerEvent <- header.Time
	}
	if time.Unix(header.Time.Int64(), 0).Sub(time.Now()) <= 0 {
		log.Warn("Prepare use too much time, missing out your turn")
		return
	}

	// If we are care about TheDAO hard-fork check whether to override the extra-data or not
	if daoBlock := self.config.DAOForkBlock; daoBlock != nil {
		// Check whether the block is among the fork extra-override range
		limit := new(big.Int).Add(daoBlock, params.DAOForkExtraRange)
		if header.Number.Cmp(daoBlock) >= 0 && header.Number.Cmp(limit) < 0 {
			// Depending whether we support or oppose the fork, override differently
			if self.config.DAOForkSupport {
				header.Extra = common.CopyBytes(params.DAOForkBlockExtra)
			} else if bytes.Equal(header.Extra, params.DAOForkBlockExtra) {
				header.Extra = []byte{} // If miner opposes, don't let it use the reserved extra-data
			}
		}
	}
	// Could potentially happen if starting to mine in an odd state.
	err := self.makeCurrent(parent, header)
	if err != nil {
		log.Error("Failed to create mining context", "err", err)
		return
	}
	// Create the current work task and check any fork transitions needed
	work := self.current
	pending, err := self.vnt.TxPool().Pending()
	if err != nil {
		log.Error("Failed to fetch pending transactions", "err", err)
		return
	}
	txs := types.NewTransactionsByPriceAndNonce(self.current.signer, pending)
	work.commitTransactions(self.mux, txs, self.chain, self.coinbase)

	// Create the new block to seal with the consensus engine
	if work.Block, err = self.engine.Finalize(self.chain, header, work.state, work.txs, work.receipts); err != nil {
		log.Error("Failed to finalize block for sealing", "err", err)
		return
	}
	blockheaderjson, _ := json.Marshal(work.Block.Header())
	blocktxjson, _ := json.Marshal(work.Block.Transactions())
	log.Debug("worker", "func", "commitNewWork", "block header", string(blockheaderjson), "block tx", string(blocktxjson))
	// We only care about logging if we're actually mining.
	if atomic.LoadInt32(&self.mining) == 1 {
		log.Info("Commit new mining work", "number", work.Block.Number(), "txs", work.tcount, "elapsed", common.PrettyDuration(time.Since(tstart)))
		self.unconfirmed.Shift(work.Block.NumberU64() - 1)
	}

	// This is time consuming. The max time has been used is 4.3ms with 220txs in a block.
	self.updateSnapshot()

	if preErr != nil {
		log.Debug("Failed to prepare header for mining", "preErr", preErr)
		return
	}

	// updateSnapshot() is time consuming. If push() before updateSnapshot(), Gvnt may be
	// stop for concurrent map iteration and map write. After push(), a block generated
	// will be write to statedb, and may be updateSnapshot() still read statedb. Then
	// an error occurs.
	self.push(work)
}

// Reset the clock for the next period
func (self *worker) resetRoundTimer(nextRoundTime *big.Int) {
	dur := time.Unix(nextRoundTime.Int64(), 0).Sub(time.Now())
	// Always make sure timer stoped and cleaned before Reset()
	self.stopRoundTimer()
	self.roundTimer.Reset(dur)
	log.Debug("Reset round timer", "header.time", nextRoundTime, "time", time.Now().Unix(), "dur", dur)
}

func (self *worker) stopRoundTimer() {
	// The stop command may be happen when the round timer is timeout but not deal with
	// the timeout event, so cleaning the channel of roundTimer is needed.
	if false == self.roundTimer.Stop() && len(self.roundTimer.C) > 0 {
		<-self.roundTimer.C
		log.Warn("worker.roundTimer.C still has a expired event, now has been cleaned")
	}
}

func (self *worker) updateSnapshot() {
	self.snapshotMu.Lock()
	defer self.snapshotMu.Unlock()

	self.snapshotBlock = types.NewBlock(
		self.current.header,
		self.current.txs,
		self.current.receipts,
	)
	self.snapshotState = self.current.state.Copy()
}

func (env *Work) commitTransactions(mux *event.TypeMux, txs *types.TransactionsByPriceAndNonce, bc *core.BlockChain, coinbase common.Address) {
	if env.gasPool == nil {
		env.gasPool = new(core.GasPool).AddGas(env.header.GasLimit)
	}

	var coalescedLogs []*types.Log

	for {
		// If we don't have enough gas for any further transactions then we're done
		if env.gasPool.Gas() < params.TxGas {
			log.Trace("Not enough gas for further transactions", "have", env.gasPool, "want", params.TxGas)
			break
		}
		// Retrieve the next transaction and abort if all done
		tx := txs.Peek()
		if tx == nil {
			break
		}
		// Error may be ignored here. The error has already been checked
		// during transaction acceptance is the transaction pool.
		//
		// We use the eip155 signer regardless of the current hf.
		from, _ := types.Sender(env.signer, tx)
		// Check whether the tx is replay protected. If we're not in the EIP155 hf
		// phase, start ignoring the sender until we do.
		if tx.Protected() && !env.config.IsEIP155(env.header.Number) {
			log.Trace("Ignoring reply protected transaction", "hash", tx.Hash(), "eip155", env.config.EIP155Block)

			txs.Pop()
			continue
		}
		// Start executing the transaction
		env.state.Prepare(tx.Hash(), common.Hash{}, env.tcount)

		err, logs := env.commitTransaction(tx, bc, coinbase, env.gasPool)
		switch err {
		case core.ErrGasLimitReached:
			// Pop the current out-of-gas transaction without shifting in the next from the account
			log.Trace("Gas limit exceeded for current block", "sender", from)
			txs.Pop()

		case core.ErrNonceTooLow:
			// New head notification data race between the transaction pool and miner, shift
			log.Trace("Skipping transaction with low nonce", "sender", from, "nonce", tx.Nonce())
			txs.Shift()

		case core.ErrNonceTooHigh:
			// Reorg notification data race between the transaction pool and miner, skip account =
			log.Trace("Skipping account with hight nonce", "sender", from, "nonce", tx.Nonce())
			txs.Pop()

		case nil:
			// Everything ok, collect the logs and shift in the next transaction from the same account
			coalescedLogs = append(coalescedLogs, logs...)
			env.tcount++
			txs.Shift()

		default:
			// Strange error, discard the transaction and get the next in line (note, the
			// nonce-too-high clause will prevent us from executing in vain).
			log.Debug("Transaction failed, account skipped", "hash", tx.Hash(), "err", err)
			txs.Shift()
		}
	}

	if len(coalescedLogs) > 0 || env.tcount > 0 {
		// make a copy, the state caches the logs and these logs get "upgraded" from pending to mined
		// logs by filling in the block hash when the block was mined by the local miner. This can
		// cause a race condition if a log was "upgraded" before the PendingLogsEvent is processed.
		cpy := make([]*types.Log, len(coalescedLogs))
		for i, l := range coalescedLogs {
			cpy[i] = new(types.Log)
			*cpy[i] = *l
		}
		go func(logs []*types.Log, tcount int) {
			if len(logs) > 0 {
				mux.Post(core.PendingLogsEvent{Logs: logs})
			}
			if tcount > 0 {
				mux.Post(core.PendingStateEvent{})
			}
		}(cpy, env.tcount)
	}
}

func (env *Work) commitTransaction(tx *types.Transaction, bc *core.BlockChain, coinbase common.Address, gp *core.GasPool) (error, []*types.Log) {
	snap := env.state.Snapshot()

	receipt, _, err := core.ApplyTransaction(env.config, bc, &coinbase, gp, env.state, env.header, tx, &env.header.GasUsed, vm.Config{})
	if err != nil {
		env.state.RevertToSnapshot(snap)
		return err, nil
	}
	env.txs = append(env.txs, tx)
	env.receipts = append(env.receipts, receipt)

	return nil, receipt.Logs
}
