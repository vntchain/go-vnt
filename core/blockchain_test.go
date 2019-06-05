// Copyright 2014 The go-ethereum Authors
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

package core

import (
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/consensus/mock"
	"github.com/vntchain/go-vnt/core/rawdb"
	"github.com/vntchain/go-vnt/core/state"
	"github.com/vntchain/go-vnt/core/types"
	"github.com/vntchain/go-vnt/core/vm"
	"github.com/vntchain/go-vnt/core/vm/election"
	"github.com/vntchain/go-vnt/crypto"
	"github.com/vntchain/go-vnt/params"
	"github.com/vntchain/go-vnt/vntdb"
)

// Test fork of length N starting from block i
func testFork(t *testing.T, blockchain *BlockChain, i, n int, full bool, comparator func(td1, td2 *big.Int)) {
	// Copy old chain up to #i into a new db
	db, blockchain2, err := newCanonical(mock.NewMock(), i, full)
	if err != nil {
		t.Fatal("could not make new canonical in testFork", err)
	}
	defer blockchain2.Stop()

	// Assert the chains have the same header/block at #i
	var hash1, hash2 common.Hash
	if full {
		hash1 = blockchain.GetBlockByNumber(uint64(i)).Hash()
		hash2 = blockchain2.GetBlockByNumber(uint64(i)).Hash()
	} else {
		hash1 = blockchain.GetHeaderByNumber(uint64(i)).Hash()
		hash2 = blockchain2.GetHeaderByNumber(uint64(i)).Hash()
	}
	if hash1 != hash2 {
		t.Errorf("chain content mismatch at %d: have hash %v, want hash %v", i, hash2, hash1)
	}
	// Extend the newly created chain
	var (
		blockChainB  []*types.Block
		headerChainB []*types.Header
	)
	if full {
		blockChainB = makeBlockChain(blockchain2.CurrentBlock(), n, mock.NewMock(), db, forkSeed)
		if _, err := blockchain2.InsertChain(blockChainB); err != nil {
			t.Fatalf("failed to insert forking chain: %v", err)
		}
	} else {
		headerChainB = makeHeaderChain(blockchain2.CurrentHeader(), n, mock.NewMock(), db, forkSeed)
		if _, err := blockchain2.InsertHeaderChain(headerChainB, 1); err != nil {
			t.Fatalf("failed to insert forking chain: %v", err)
		}
	}
	// Sanity check that the forked chain can be imported into the original
	var tdPre, tdPost *big.Int

	if full {
		tdPre = blockchain.GetTdByHash(blockchain.CurrentBlock().Hash())
		if err := testBlockChainImport(blockChainB, blockchain); err != nil {
			t.Fatalf("failed to import forked block chain: %v", err)
		}
		tdPost = blockchain.GetTdByHash(blockChainB[len(blockChainB)-1].Hash())
	} else {
		tdPre = blockchain.GetTdByHash(blockchain.CurrentHeader().Hash())
		if err := testHeaderChainImport(headerChainB, blockchain); err != nil {
			t.Fatalf("failed to import forked header chain: %v", err)
		}
		tdPost = blockchain.GetTdByHash(headerChainB[len(headerChainB)-1].Hash())
	}
	// Compare the total difficulties of the chains
	comparator(tdPre, tdPost)
}

func printChain(bc *BlockChain) {
	for i := bc.CurrentBlock().Number().Uint64(); i > 0; i-- {
		b := bc.GetBlockByNumber(uint64(i))
		fmt.Printf("\t%x %v\n", b.Hash(), b.Difficulty())
	}
}

// testBlockChainImport tries to process a chain of blocks, writing them into
// the database if successful.
func testBlockChainImport(chain types.Blocks, blockchain *BlockChain) error {
	for _, block := range chain {
		// Try and process the block
		err := blockchain.engine.VerifyHeader(blockchain, block.Header(), true)
		if err == nil {
			err = blockchain.validator.ValidateBody(block)
		}
		if err != nil {
			if err == ErrKnownBlock {
				continue
			}
			return err
		}
		statedb, err := state.New(blockchain.GetBlockByHash(block.ParentHash()).Root(), blockchain.stateCache)
		if err != nil {
			return err
		}
		receipts, _, usedGas, err := blockchain.Processor().Process(block, statedb, vm.Config{})
		if err != nil {
			blockchain.reportBlock(block, receipts, err)
			return err
		}
		err = blockchain.validator.ValidateState(block, blockchain.GetBlockByHash(block.ParentHash()), statedb, receipts, usedGas)
		if err != nil {
			blockchain.reportBlock(block, receipts, err)
			return err
		}
		blockchain.mu.Lock()
		rawdb.WriteTd(blockchain.db, block.Hash(), block.NumberU64(), new(big.Int).Add(block.Difficulty(), blockchain.GetTdByHash(block.ParentHash())))
		rawdb.WriteBlock(blockchain.db, block)
		statedb.Commit(false)
		blockchain.mu.Unlock()
	}
	return nil
}

// testHeaderChainImport tries to process a chain of header, writing them into
// the database if successful.
func testHeaderChainImport(chain []*types.Header, blockchain *BlockChain) error {
	for _, header := range chain {
		// Try and validate the header
		if err := blockchain.engine.VerifyHeader(blockchain, header, false); err != nil {
			return err
		}
		// Manually insert the header into the database, but don't reorganise (allows subsequent testing)
		blockchain.mu.Lock()
		rawdb.WriteTd(blockchain.db, header.Hash(), header.Number.Uint64(), new(big.Int).Add(header.Difficulty, blockchain.GetTdByHash(header.ParentHash)))
		rawdb.WriteHeader(blockchain.db, header)
		blockchain.mu.Unlock()
	}
	return nil
}

func insertChain(done chan bool, blockchain *BlockChain, chain types.Blocks, t *testing.T) {
	_, err := blockchain.InsertChain(chain)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	done <- true
}

func TestLastBlock(t *testing.T) {
	_, blockchain, err := newCanonical(mock.NewMock(), 0, true)
	if err != nil {
		t.Fatalf("failed to create pristine chain: %v", err)
	}
	defer blockchain.Stop()

	blocks := makeBlockChain(blockchain.CurrentBlock(), 1, mock.NewMock(), blockchain.db, 0)
	if _, err := blockchain.InsertChain(blocks); err != nil {
		t.Fatalf("Failed to insert block: %v", err)
	}
	if blocks[len(blocks)-1].Hash() != rawdb.ReadHeadBlockHash(blockchain.db) {
		t.Fatalf("Write/Get HeadBlockHash failed")
	}
}

// Tests that given a starting canonical chain of a given size, it can be extended
// with various length chains.
func TestExtendCanonicalHeaders(t *testing.T) { testExtendCanonical(t, false) }
func TestExtendCanonicalBlocks(t *testing.T)  { testExtendCanonical(t, true) }

func testExtendCanonical(t *testing.T, full bool) {
	length := 5

	// Make first chain starting from genesis
	_, processor, err := newCanonical(mock.NewMock(), length, full)
	if err != nil {
		t.Fatalf("failed to make new canonical chain: %v", err)
	}
	defer processor.Stop()

	// Define the difficulty comparator
	better := func(td1, td2 *big.Int) {
		if td2.Cmp(td1) <= 0 {
			t.Errorf("total difficulty mismatch: have %v, expected more than %v", td2, td1)
		}
	}
	// Start fork from current height
	testFork(t, processor, length, 1, full, better)
	testFork(t, processor, length, 2, full, better)
	testFork(t, processor, length, 5, full, better)
	testFork(t, processor, length, 10, full, better)
}

// Tests that given a starting canonical chain of a given size, creating shorter
// forks do not take canonical ownership.
func TestShorterForkHeaders(t *testing.T) { testShorterFork(t, false) }
func TestShorterForkBlocks(t *testing.T)  { testShorterFork(t, true) }

func testShorterFork(t *testing.T, full bool) {
	length := 10

	// Make first chain starting from genesis
	_, processor, err := newCanonical(mock.NewMock(), length, full)
	if err != nil {
		t.Fatalf("failed to make new canonical chain: %v", err)
	}
	defer processor.Stop()

	// Define the difficulty comparator
	worse := func(td1, td2 *big.Int) {
		if td2.Cmp(td1) >= 0 {
			t.Errorf("total difficulty mismatch: have %v, expected less than %v", td2, td1)
		}
	}
	// Sum of numbers must be less than `length` for this to be a shorter fork
	testFork(t, processor, 0, 3, full, worse)
	testFork(t, processor, 0, 7, full, worse)
	testFork(t, processor, 1, 1, full, worse)
	testFork(t, processor, 1, 7, full, worse)
	testFork(t, processor, 5, 3, full, worse)
	testFork(t, processor, 5, 4, full, worse)
}

// Tests that given a starting canonical chain of a given size, creating longer
// forks do take canonical ownership.
func TestLongerForkHeaders(t *testing.T) { testLongerFork(t, false) }
func TestLongerForkBlocks(t *testing.T)  { testLongerFork(t, true) }

func testLongerFork(t *testing.T, full bool) {
	length := 10

	// Make first chain starting from genesis
	_, processor, err := newCanonical(mock.NewMock(), length, full)
	if err != nil {
		t.Fatalf("failed to make new canonical chain: %v", err)
	}
	defer processor.Stop()

	// Define the difficulty comparator
	better := func(td1, td2 *big.Int) {
		if td2.Cmp(td1) <= 0 {
			t.Errorf("total difficulty mismatch: have %v, expected more than %v", td2, td1)
		}
	}
	// Sum of numbers must be greater than `length` for this to be a longer fork
	testFork(t, processor, 0, 11, full, better)
	testFork(t, processor, 0, 15, full, better)
	testFork(t, processor, 1, 10, full, better)
	testFork(t, processor, 1, 12, full, better)
	testFork(t, processor, 5, 6, full, better)
	testFork(t, processor, 5, 8, full, better)
}

// Tests that given a starting canonical chain of a given size, creating equal
// forks do take canonical ownership.
func TestEqualForkHeaders(t *testing.T) { testEqualFork(t, false) }
func TestEqualForkBlocks(t *testing.T)  { testEqualFork(t, true) }

func testEqualFork(t *testing.T, full bool) {
	length := 10

	// Make first chain starting from genesis
	_, processor, err := newCanonical(mock.NewMock(), length, full)
	if err != nil {
		t.Fatalf("failed to make new canonical chain: %v", err)
	}
	defer processor.Stop()

	// Define the difficulty comparator
	equal := func(td1, td2 *big.Int) {
		if td2.Cmp(td1) != 0 {
			t.Errorf("total difficulty mismatch: have %v, want %v", td2, td1)
		}
	}
	// Sum of numbers must be equal to `length` for this to be an equal fork
	testFork(t, processor, 0, 10, full, equal)
	testFork(t, processor, 1, 9, full, equal)
	testFork(t, processor, 2, 8, full, equal)
	testFork(t, processor, 5, 5, full, equal)
	testFork(t, processor, 6, 4, full, equal)
	testFork(t, processor, 9, 1, full, equal)
}

// Tests that chains missing links do not get accepted by the processor.
func TestBrokenHeaderChain(t *testing.T) { testBrokenChain(t, false) }
func TestBrokenBlockChain(t *testing.T)  { testBrokenChain(t, true) }

func testBrokenChain(t *testing.T, full bool) {
	// Make chain starting from genesis
	db, blockchain, err := newCanonical(mock.NewMock(), 10, full)
	if err != nil {
		t.Fatalf("failed to make new canonical chain: %v", err)
	}
	defer blockchain.Stop()

	// Create a forked chain, and try to insert with a missing link
	if full {
		chain := makeBlockChain(blockchain.CurrentBlock(), 5, mock.NewMock(), db, forkSeed)[1:]
		if err := testBlockChainImport(chain, blockchain); err == nil {
			t.Errorf("broken block chain not reported")
		}
	} else {
		chain := makeHeaderChain(blockchain.CurrentHeader(), 5, mock.NewMock(), db, forkSeed)[1:]
		if err := testHeaderChainImport(chain, blockchain); err == nil {
			t.Errorf("broken header chain not reported")
		}
	}
}

// 正常产块情况，在当前主链后产生新区块，当前新区块应当成为头区块
// Header chain测试
func TestReorgLongHeaders(t *testing.T) { testReorgLong(t, false) }

// Block chain测试
func TestReorgLongBlocks(t *testing.T) { testReorgLong(t, true) }

func testReorgLong(t *testing.T, full bool) {
	testReorg(t, []int64{0, 0, 0}, []int64{0, 0, 0, 0}, 4, full, true, false)
}

// 测试攻击链的场景，填充过去某个时间点缺失的区块，创造更长的链
// Test some block producers make up fake chain by filling skipped blocks.
// first : A ->  -> B -> C
// second: A -> B' -> C' -> D'
func TestReorgFakeChainFillMissTimeAndLongHeaders(t *testing.T) {
	testReorgFakeChainFillMissTimeAndLong(t, false)
}
func TestReorgFakeChainFillMissTimeAndLongBlocks(t *testing.T) {
	testReorgFakeChainFillMissTimeAndLong(t, true)
}
func testReorgFakeChainFillMissTimeAndLong(t *testing.T, full bool) {
	mainChain := []int64{0, 2, 0, 0}
	fakeChain := []int64{0, 0, 0, 0, 0} // LIB之前的区块插入失败
	testReorg(t, mainChain, fakeChain, 4, full, false, true)
}

// 临时分叉情况，网络突然中断，然后下一轮结束前恢复
// Case 1.1 two blocks in the same height for bad network
// first:  A -> B, 代表本地链
// second: A -> -> C, 代表远端链
// first is still main chain
func TestReorgSameHeightFirstIsEarlyHeader(t *testing.T) {
	testReorgSameHeightFirstIsEarly(t, false)
}
func TestReorgSameHeightFirstIsEarlyBlock(t *testing.T) {
	testReorgSameHeightFirstIsEarly(t, true)
}
func testReorgSameHeightFirstIsEarly(t *testing.T, full bool) {
	first := []int64{0, 0}
	second := []int64{0, 2}
	testReorg(t, first, second, 2, full, true, true)
}

// 临时分叉情况，网络突然中断，然后下一轮结束前恢复
// Case 1.2 two blocks in the same height for bad network
// 与Case 1.1相反的过程，但都应当选择`A -> B`
// first:  A -> -> C
// second: A -> B
func TestReorgSameHeightSecondIsEarlyHeader(t *testing.T) {
	testReorgSameHeightSecondIsEarly(t, false)
}
func TestReorgSameHeightSecondIsEarlyBlock(t *testing.T) {
	testReorgSameHeightSecondIsEarly(t, true)
}
func testReorgSameHeightSecondIsEarly(t *testing.T, full bool) {
	first := []int64{0, 2}
	second := []int64{0, 0}
	testReorg(t, first, second, 2, full, true, false)
}

// 临时分叉情况，网络在第N轮突然中断，然后在N+2轮结束恢复
// Case 2.1
// first:  A -> B
// second: A ->  -> C -> D
func TestReorgFirstIsEarlyButShortHeader(t *testing.T) {
	testReorgFirstIsEarlyButShort(t, false)
}
func TestReorgFirstIsEarlyButShortBlock(t *testing.T) {
	testReorgFirstIsEarlyButShort(t, true)
}
func testReorgFirstIsEarlyButShort(t *testing.T, full bool) {
	first := []int64{0, 0}
	second := []int64{0, 2, 0}
	testReorg(t, first, second, 3, full, true, false)
}

// 临时分叉情况，网络在第N轮突然中断，然后在N+2轮结束恢复
// Case 2.2
// first:  A ->  -> C -> D
// second: A -> B
func TestReorgSecondIsEarlyButShortHeader(t *testing.T) {
	testReorgSecondIsEarlyButShort(t, false)
}
func TestReorgSecondIsEarlyButShortBlock(t *testing.T) {
	testReorgSecondIsEarlyButShort(t, true)
}
func testReorgSecondIsEarlyButShort(t *testing.T, full bool) {
	first := []int64{0, 2, 0}
	second := []int64{0, 0}
	testReorg(t, first, second, 3, full, false, true)
}

// 临时分叉情况，络在第N轮突然中断，然后在N+1轮结束恢复，然后在N+2轮产生了新区块，
// 在N+1轮产生区块的少数节点，如果没再case 1.1切回到first，在first生成新区块后切回first
// Case 3.1
// first:  A -> B ->  -> D
// second: A ->  -> C
func TestReorgFirstIsEarlyAndLongHeader(t *testing.T) {
	testReorgFirstIsEarlyAndLong(t, false)
}
func TestReorgFirstIsEarlyAndLongBlock(t *testing.T) {
	testReorgFirstIsEarlyAndLong(t, true)
}
func testReorgFirstIsEarlyAndLong(t *testing.T, full bool) {
	first := []int64{0, 0, 2}
	second := []int64{0, 2}
	testReorg(t, first, second, 3, full, false, true)
}

// 临时分叉情况，络在第N轮突然中断，然后在N+1轮结束恢复，然后在N+2轮产生了新区块，
// 在N+1轮产生区块的少数节点，如果没再case 1.1切回到first，在first生成新区块后切回first
// Case 3.2
// first:  A ->  -> C
// second: A -> B ->  -> D
func TestReorgSecondIsEarlyAndLongHeader(t *testing.T) {
	testReorgSecondIsEarlyAndLong(t, false)
}
func TestReorgSecondIsEarlyAndLongBlock(t *testing.T) {
	testReorgSecondIsEarlyAndLong(t, true)
}
func testReorgSecondIsEarlyAndLong(t *testing.T, full bool) {
	first := []int64{0, 2}
	second := []int64{0, 0, 2}
	testReorg(t, first, second, 3, full, true, false)
}

// first: 第一条链，是当前主链的时间戳偏移，在当前生成的区块上进一步做时间戳偏移
// second: 第二条链，分叉链、攻击链的时间戳偏移
// td: reorg后成为主链的区块的总高度
// full: 是否是header chain测试
// firstIsMain: reorg后，第一条链是主链
func testReorg(t *testing.T, first, second []int64, td int64, full, secondInsertSuccess, firstIsMain bool) {
	// Create a pristine chain and database
	db, blockchain, err := newCanonical(mock.NewMock(), 0, full)
	if err != nil {
		t.Fatalf("failed to create pristine chain: %v", err)
	}
	defer blockchain.Stop()

	// 生成区块
	firstBlocks, _ := GenerateChain(params.TestChainConfig, blockchain.CurrentBlock(), mock.NewMock(), db, len(first), func(i int, b *BlockGen) {
		b.OffsetTime(first[i])
	})
	secondBlocks, _ := GenerateChain(params.TestChainConfig, blockchain.CurrentBlock(), mock.NewMock(), db, len(second), func(i int, b *BlockGen) {
		b.OffsetTime(second[i])
	})

	printBlocks(t, "first  blocks:", firstBlocks)
	printBlocks(t, "second blocks:", secondBlocks)

	// 调用测试函数
	testFun := testReorgBlock
	if !full {
		testFun = testReorgHeader
	}
	testFun(t, blockchain, firstBlocks, secondBlocks, td, secondInsertSuccess, firstIsMain)
}

// testReorgBlock block chain的reorg测试
func testReorgBlock(t *testing.T, blockchain *BlockChain, firstBlocks, secondBlocks types.Blocks, td int64, secondInsertSuccess, firstIsMain bool) {
	// 插入的区块链，first代表当前主链，必须成功
	if _, err := blockchain.InsertChain(firstBlocks); err != nil {
		t.Fatalf("failed to insert first chain: %v", err)
	}
	// second是攻击链或临时分叉，按要求执行，可能成功，可能失败
	{
		_, err := blockchain.InsertChain(secondBlocks)
		if secondInsertSuccess {
			if err != nil {
				t.Logf("failed to insert second chain: %v", err)
			}
		} else {
			if err == nil {
				t.Logf("insert second chain successed, want failed")
			}
		}
	}

	// 检查链是否是连接正确的
	block := blockchain.CurrentBlock()
	for prev := blockchain.GetBlockByNumber(blockchain.CurrentBlock().NumberU64() - 1); prev.NumberU64() != 0; block, prev = prev, blockchain.GetBlockByNumber(prev.NumberU64()-1) {
		// t.Logf("block hash: %s\n", block.Hash().String())
		if block.ParentHash() != prev.Hash() {
			t.Errorf("parent block hash mismatch: have %x, want %x", block.ParentHash(), prev.Hash())
		}
	}

	// 链难度检查
	wantTd := new(big.Int).Add(blockchain.genesisBlock.Difficulty(), big.NewInt(td))
	if have := blockchain.GetTdByHash(blockchain.CurrentBlock().Hash()); have.Cmp(wantTd) != 0 {
		t.Errorf("total difficulty mismatch: have %v, want %v", have, wantTd)
	}

	// 检查主链的区块是否匹配
	targetChain := firstBlocks
	if !firstIsMain {
		targetChain = secondBlocks
	}
	blocks := make([]*types.Block, 0)
	blk := blockchain.CurrentBlock()
	for prev := blockchain.GetBlockByNumber(blockchain.CurrentBlock().NumberU64() - 1); blk.NumberU64() != 0; blk, prev = prev, blockchain.GetBlockByNumber(prev.NumberU64()-1) {
		blocks = append(blocks, blk)
	}
	// 反序，高度从低到高
	for i, j := 0, len(blocks)-1; i < j; i, j = i+1, j-1 {
		blocks[i], blocks[j] = blocks[j], blocks[i]
	}
	printBlocks(t, "main chain:", blocks)

	// 打印难度日志
	genesis := blockchain.genesisBlock
	td0 := blockchain.GetTdByHash(genesis.Hash())
	t.Logf("genesis block %s, total diff: %v\n", genesis.Hash().String(), td0.String())
	for i, b := range blocks {
		td := blockchain.GetTdByHash(b.Hash())
		t.Logf("block [%d] %s, total diff: %v\n", i, b.Hash().String(), td.String())
	}

	if len(blocks) != len(targetChain) {
		t.Fatalf("total length mismatch: have %v, want %v", len(blocks), len(targetChain))
	}
	for i, b := range targetChain {
		bc := blocks[i]
		if b.Hash() != bc.Hash() {
			t.Errorf("block [%d] mismatch: have %v, want: %v", i, bc.Hash().String(), b.Hash().String())
		}
	}
}

// testReorgBlock header chain的reorg测试
func testReorgHeader(t *testing.T, blockchain *BlockChain, firstBlocks, secondBlocks types.Blocks, td int64, secondInsertSuccess, firstIsMain bool) {
	firstHeaders := make([]*types.Header, len(firstBlocks))
	for i, block := range firstBlocks {
		firstHeaders[i] = block.Header()
	}
	secondHeaders := make([]*types.Header, len(secondBlocks))
	for i, block := range secondBlocks {
		secondHeaders[i] = block.Header()
	}

	// first是主链的header chain，必须成功
	if _, err := blockchain.InsertHeaderChain(firstHeaders, 1); err != nil {
		t.Fatalf("failed to insert first header chain: %v", err)
	}
	// second是攻击链或临时分叉的header chain，按要求觉得成功与否
	{
		_, err := blockchain.InsertHeaderChain(secondHeaders, 1)
		if secondInsertSuccess {
			if err != nil {
				t.Fatalf("failed to insert second header chain: %v", err)
			}
		} else {
			if err == nil {
				t.Fatalf("insert second header chain successed, want failed")
			}
		}
	}

	// Check that the chain is valid number and link wise
	header := blockchain.CurrentHeader()
	for prev := blockchain.GetHeaderByNumber(blockchain.CurrentHeader().Number.Uint64() - 1); prev.Number.Uint64() != 0; header, prev = prev, blockchain.GetHeaderByNumber(prev.Number.Uint64()-1) {
		if header.ParentHash != prev.Hash() {
			t.Errorf("parent header hash mismatch: have %x, want %x", header.ParentHash, prev.Hash())
		}
	}

	// Make sure the chain total difficulty is the correct one
	wantTd := new(big.Int).Add(blockchain.genesisBlock.Difficulty(), big.NewInt(td))
	if have := blockchain.GetTdByHash(blockchain.CurrentHeader().Hash()); have.Cmp(wantTd) != 0 {
		t.Errorf("total difficulty mismatch: have %v, want %v", have, wantTd)
	}

	// 检查主链是否匹配
	targetHeader := firstHeaders
	if !firstIsMain {
		targetHeader = secondHeaders
	}
	chainHeader := make([]*types.Header, 0)
	header = blockchain.CurrentHeader()
	for prev := blockchain.GetHeaderByNumber(blockchain.CurrentHeader().Number.Uint64() - 1); header.Number.Uint64() != 0; header, prev = prev, blockchain.GetHeaderByNumber(prev.Number.Uint64()-1) {
		chainHeader = append(chainHeader, header)
	}
	for i, j := 0, len(chainHeader)-1; i < j; i, j = i+1, j-1 {
		chainHeader[i], chainHeader[j] = chainHeader[j], chainHeader[i]
	}
	printHeaders(t, "header chain", chainHeader)

	if len(targetHeader) != len(chainHeader) {
		t.Fatalf("header chain lengeth mismatch, have: %d, want: %d", len(chainHeader), len(targetHeader))
	}
	for i := 0; i < len(targetHeader); i++ {
		if targetHeader[i].Hash() != chainHeader[i].Hash() {
			t.Errorf("header chain mismatch %d, have: %s, want: %s\n", i, chainHeader[i].Hash(), targetHeader[i].Hash())
		}
	}
}

func printBlocks(t *testing.T, tag string, blocks []*types.Block) {
	t.Logf("%s\n", tag)
	for i, b := range blocks {
		t.Logf("b [%d], h: %s, time: %s, diff: %s, hash: %s\n", i, b.Number().String(), b.Time().String(), b.Difficulty().String(), b.Hash().String())
	}
}

func printHeaders(t *testing.T, tag string, headers []*types.Header) {
	t.Logf("%s\n", tag)
	for i, h := range headers {
		t.Logf("header [%d], h: %s, time: %s, diff: %s, hash: %s\n", i, h.Number.String(), h.Time.String(), h.Difficulty.String(), h.Hash().String())
	}
}

// Tests that the insertion functions detect banned hashes.
func TestBadHeaderHashes(t *testing.T) { testBadHashes(t, false) }
func TestBadBlockHashes(t *testing.T)  { testBadHashes(t, true) }

func testBadHashes(t *testing.T, full bool) {
	// Create a pristine chain and database
	db, blockchain, err := newCanonical(mock.NewMock(), 0, full)
	if err != nil {
		t.Fatalf("failed to create pristine chain: %v", err)
	}
	defer blockchain.Stop()

	// Create a chain, ban a hash and try to import
	if full {
		blocks := makeBlockChain(blockchain.CurrentBlock(), 3, mock.NewMock(), db, 10)

		BadHashes[blocks[2].Header().Hash()] = true
		defer func() { delete(BadHashes, blocks[2].Header().Hash()) }()

		_, err = blockchain.InsertChain(blocks)
	} else {
		headers := makeHeaderChain(blockchain.CurrentHeader(), 3, mock.NewMock(), db, 10)

		BadHashes[headers[2].Hash()] = true
		defer func() { delete(BadHashes, headers[2].Hash()) }()

		_, err = blockchain.InsertHeaderChain(headers, 1)
	}
	if err != ErrBlacklistedHash {
		t.Errorf("error mismatch: have: %v, want: %v", err, ErrBlacklistedHash)
	}
}

// Tests that bad hashes are detected on boot, and the chain rolled back to a
// good state prior to the bad hash.
func TestReorgBadHeaderHashes(t *testing.T) { testReorgBadHashes(t, false) }
func TestReorgBadBlockHashes(t *testing.T)  { testReorgBadHashes(t, true) }

func testReorgBadHashes(t *testing.T, full bool) {
	// Create a pristine chain and database
	db, blockchain, err := newCanonical(mock.NewMock(), 0, full)
	if err != nil {
		t.Fatalf("failed to create pristine chain: %v", err)
	}
	// Create a chain, import and ban afterwards
	headers := makeHeaderChain(blockchain.CurrentHeader(), 4, mock.NewMock(), db, 10)
	blocks := makeBlockChain(blockchain.CurrentBlock(), 4, mock.NewMock(), db, 10)

	if full {
		if _, err = blockchain.InsertChain(blocks); err != nil {
			t.Errorf("failed to import blocks: %v", err)
		}
		if blockchain.CurrentBlock().Hash() != blocks[3].Hash() {
			t.Errorf("last block hash mismatch: have: %x, want %x", blockchain.CurrentBlock().Hash(), blocks[3].Header().Hash())
		}
		BadHashes[blocks[3].Header().Hash()] = true
		defer func() { delete(BadHashes, blocks[3].Header().Hash()) }()
	} else {
		if _, err = blockchain.InsertHeaderChain(headers, 1); err != nil {
			t.Errorf("failed to import headers: %v", err)
		}
		if blockchain.CurrentHeader().Hash() != headers[3].Hash() {
			t.Errorf("last header hash mismatch: have: %x, want %x", blockchain.CurrentHeader().Hash(), headers[3].Hash())
		}
		BadHashes[headers[3].Hash()] = true
		defer func() { delete(BadHashes, headers[3].Hash()) }()
	}
	blockchain.Stop()

	// Create a new BlockChain and check that it rolled back the state.
	ncm, err := NewBlockChain(blockchain.db, nil, blockchain.chainConfig, mock.NewMock(), vm.Config{})
	if err != nil {
		t.Fatalf("failed to create new chain manager: %v", err)
	}
	if full {
		if ncm.CurrentBlock().Hash() != blocks[2].Header().Hash() {
			t.Errorf("last block hash mismatch: have: %x, want %x", ncm.CurrentBlock().Hash(), blocks[2].Header().Hash())
		}
		if blocks[2].Header().GasLimit != ncm.GasLimit() {
			t.Errorf("last  block gasLimit mismatch: have: %d, want %d", ncm.GasLimit(), blocks[2].Header().GasLimit)
		}
	} else {
		if ncm.CurrentHeader().Hash() != headers[2].Hash() {
			t.Errorf("last header hash mismatch: have: %x, want %x", ncm.CurrentHeader().Hash(), headers[2].Hash())
		}
	}
	ncm.Stop()
}

// Tests chain insertions in the face of one entity containing an invalid nonce.
func TestHeadersInsertNonceError(t *testing.T) { testInsertNonceError(t, false) }
func TestBlocksInsertNonceError(t *testing.T)  { testInsertNonceError(t, true) }

func testInsertNonceError(t *testing.T, full bool) {
	for i := 1; i < 25 && !t.Failed(); i++ {
		// Create a pristine chain and database
		db, blockchain, err := newCanonical(mock.NewMock(), 0, full)
		if err != nil {
			t.Fatalf("failed to create pristine chain: %v", err)
		}
		defer blockchain.Stop()

		// Create and insert a chain with a failing nonce
		var (
			failAt  int
			failRes int
			failNum uint64
		)
		if full {
			blocks := makeBlockChain(blockchain.CurrentBlock(), i, mock.NewMock(), db, 0)
			fmt.Println(blocks)
			failAt = rand.Int() % len(blocks)
			failNum = blocks[failAt].NumberU64()

			blockchain.engine = mock.NewMockFail(failNum)
			failRes, err = blockchain.InsertChain(blocks)
		} else {
			headers := makeHeaderChain(blockchain.CurrentHeader(), i, mock.NewMock(), db, 0)

			failAt = rand.Int() % len(headers)
			failNum = headers[failAt].Number.Uint64()

			blockchain.engine = mock.NewMockFail(failNum)
			blockchain.hc.engine = blockchain.engine
			failRes, err = blockchain.InsertHeaderChain(headers, 1)
		}
		// Check that the returned error indicates the failure.
		if failRes != failAt {
			t.Errorf("test %d: failure index mismatch: have %d, want %d", i, failRes, failAt)
		}
		// Check that all no blocks after the failing block have been inserted.
		for j := 0; j < i-failAt; j++ {
			if full {
				if block := blockchain.GetBlockByNumber(failNum + uint64(j)); block != nil {
					t.Errorf("test %d: invalid block in chain: %v", i, block)
				}
			} else {
				if header := blockchain.GetHeaderByNumber(failNum + uint64(j)); header != nil {
					t.Errorf("test %d: invalid header in chain: %v", i, header)
				}
			}
		}
	}
}

// Tests that fast importing a block chain produces the same chain data as the
// classical full block processing.
func TestFastVsFullChains(t *testing.T) {
	// Configure and generate a sample block chain
	var (
		gendb   = vntdb.NewMemDatabase()
		key, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address = crypto.PubkeyToAddress(key.PublicKey)
		funds   = big.NewInt(1000000000)
		gspec   = &Genesis{
			Config: params.TestChainConfig,
			Alloc:  GenesisAlloc{address: {Balance: funds}},
		}
		genesis = gspec.MustCommit(gendb)
		signer  = types.NewHubbleSigner(gspec.Config.ChainID)
	)
	blocks, receipts := GenerateChain(gspec.Config, genesis, mock.NewMock(), gendb, 1024, func(i int, block *BlockGen) {
		block.SetCoinbase(common.Address{0x00})

		// If the block number is multiple of 3, send a few bonus transactions to the miner
		if i%3 == 2 {
			for j := 0; j < i%4+1; j++ {
				tx, err := types.SignTx(types.NewTransaction(block.TxNonce(address), common.Address{0x00}, big.NewInt(1000), params.TxGas, nil, nil), signer, key)
				if err != nil {
					panic(err)
				}
				block.AddTx(tx)
			}
		}
	})
	// Import the chain as an archive node for the comparison baseline
	archiveDb := vntdb.NewMemDatabase()
	gspec.MustCommit(archiveDb)
	archive, _ := NewBlockChain(archiveDb, nil, gspec.Config, mock.NewMock(), vm.Config{})
	defer archive.Stop()

	if n, err := archive.InsertChain(blocks); err != nil {
		t.Fatalf("failed to process block %d: %v", n, err)
	}
	// Fast import the chain as a non-archive node to test
	fastDb := vntdb.NewMemDatabase()
	gspec.MustCommit(fastDb)
	fast, _ := NewBlockChain(fastDb, nil, gspec.Config, mock.NewMock(), vm.Config{})
	defer fast.Stop()

	headers := make([]*types.Header, len(blocks))
	for i, block := range blocks {
		headers[i] = block.Header()
	}
	if n, err := fast.InsertHeaderChain(headers, 1); err != nil {
		t.Fatalf("failed to insert header %d: %v", n, err)
	}
	if n, err := fast.InsertReceiptChain(blocks, receipts); err != nil {
		t.Fatalf("failed to insert receipt %d: %v", n, err)
	}
	// Iterate over all chain data components, and cross reference
	for i := 0; i < len(blocks); i++ {
		num, hash := blocks[i].NumberU64(), blocks[i].Hash()

		if ftd, atd := fast.GetTdByHash(hash), archive.GetTdByHash(hash); ftd.Cmp(atd) != 0 {
			t.Errorf("block #%d [%x]: td mismatch: have %v, want %v", num, hash, ftd, atd)
		}
		if fheader, aheader := fast.GetHeaderByHash(hash), archive.GetHeaderByHash(hash); fheader.Hash() != aheader.Hash() {
			t.Errorf("block #%d [%x]: header mismatch: have %v, want %v", num, hash, fheader, aheader)
		}
		if fblock, ablock := fast.GetBlockByHash(hash), archive.GetBlockByHash(hash); fblock.Hash() != ablock.Hash() {
			t.Errorf("block #%d [%x]: block mismatch: have %v, want %v", num, hash, fblock, ablock)
		} else if types.DeriveSha(fblock.Transactions()) != types.DeriveSha(ablock.Transactions()) {
			t.Errorf("block #%d [%x]: transactions mismatch: have %v, want %v", num, hash, fblock.Transactions(), ablock.Transactions())
		}
		if freceipts, areceipts := rawdb.ReadReceipts(fastDb, hash, *rawdb.ReadHeaderNumber(fastDb, hash)), rawdb.ReadReceipts(archiveDb, hash, *rawdb.ReadHeaderNumber(archiveDb, hash)); types.DeriveSha(freceipts) != types.DeriveSha(areceipts) {
			t.Errorf("block #%d [%x]: receipts mismatch: have %v, want %v", num, hash, freceipts, areceipts)
		}
	}
	// Check that the canonical chains are the same between the databases
	for i := 0; i < len(blocks)+1; i++ {
		if fhash, ahash := rawdb.ReadCanonicalHash(fastDb, uint64(i)), rawdb.ReadCanonicalHash(archiveDb, uint64(i)); fhash != ahash {
			t.Errorf("block #%d: canonical hash mismatch: have %v, want %v", i, fhash, ahash)
		}
	}
}

// Tests that various import methods move the chain head pointers to the correct
// positions.
func TestLightVsFastVsFullChainHeads(t *testing.T) {
	// Configure and generate a sample block chain
	var (
		gendb   = vntdb.NewMemDatabase()
		key, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address = crypto.PubkeyToAddress(key.PublicKey)
		funds   = big.NewInt(1000000000)
		gspec   = &Genesis{Config: params.TestChainConfig, Alloc: GenesisAlloc{address: {Balance: funds}}}
		genesis = gspec.MustCommit(gendb)
	)
	height := uint64(1024)
	blocks, receipts := GenerateChain(gspec.Config, genesis, mock.NewMock(), gendb, int(height), nil)

	// Configure a subchain to roll back
	remove := []common.Hash{}
	for _, block := range blocks[height/2:] {
		remove = append(remove, block.Hash())
	}
	// Create a small assertion method to check the three heads
	assert := func(t *testing.T, kind string, chain *BlockChain, header uint64, fast uint64, block uint64) {
		if num := chain.CurrentBlock().NumberU64(); num != block {
			t.Errorf("%s head block mismatch: have #%v, want #%v", kind, num, block)
		}
		if num := chain.CurrentFastBlock().NumberU64(); num != fast {
			t.Errorf("%s head fast-block mismatch: have #%v, want #%v", kind, num, fast)
		}
		if num := chain.CurrentHeader().Number.Uint64(); num != header {
			t.Errorf("%s head header mismatch: have #%v, want #%v", kind, num, header)
		}
	}
	// Import the chain as an archive node and ensure all pointers are updated
	archiveDb := vntdb.NewMemDatabase()
	gspec.MustCommit(archiveDb)

	archive, _ := NewBlockChain(archiveDb, nil, gspec.Config, mock.NewMock(), vm.Config{})
	if n, err := archive.InsertChain(blocks); err != nil {
		t.Fatalf("failed to process block %d: %v", n, err)
	}
	defer archive.Stop()

	assert(t, "archive", archive, height, height, height)
	archive.Rollback(remove)
	assert(t, "archive", archive, height/2, height/2, height/2)

	// Import the chain as a non-archive node and ensure all pointers are updated
	fastDb := vntdb.NewMemDatabase()
	gspec.MustCommit(fastDb)
	fast, _ := NewBlockChain(fastDb, nil, gspec.Config, mock.NewMock(), vm.Config{})
	defer fast.Stop()

	headers := make([]*types.Header, len(blocks))
	for i, block := range blocks {
		headers[i] = block.Header()
	}
	if n, err := fast.InsertHeaderChain(headers, 1); err != nil {
		t.Fatalf("failed to insert header %d: %v", n, err)
	}
	if n, err := fast.InsertReceiptChain(blocks, receipts); err != nil {
		t.Fatalf("failed to insert receipt %d: %v", n, err)
	}
	assert(t, "fast", fast, height, height, 0)
	fast.Rollback(remove)
	assert(t, "fast", fast, height/2, height/2, 0)

	// Import the chain as a light node and ensure all pointers are updated
	lightDb := vntdb.NewMemDatabase()
	gspec.MustCommit(lightDb)

	light, _ := NewBlockChain(lightDb, nil, gspec.Config, mock.NewMock(), vm.Config{})
	if n, err := light.InsertHeaderChain(headers, 1); err != nil {
		t.Fatalf("failed to insert header %d: %v", n, err)
	}
	defer light.Stop()

	assert(t, "light", light, height, 0, 0)
	light.Rollback(remove)
	assert(t, "light", light, height/2, 0, 0)
}

// Tests that chain reorganisations handle transaction removals and reinsertions.
func TestChainTxReorgs(t *testing.T) {
	var (
		key1, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		key2, _ = crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
		key3, _ = crypto.HexToECDSA("49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
		addr1   = crypto.PubkeyToAddress(key1.PublicKey)
		addr2   = crypto.PubkeyToAddress(key2.PublicKey)
		addr3   = crypto.PubkeyToAddress(key3.PublicKey)
		db      = vntdb.NewMemDatabase()
		gspec   = &Genesis{
			Config:   params.TestChainConfig,
			GasLimit: 3141592,
			Alloc: GenesisAlloc{
				addr1: {Balance: big.NewInt(1000000)},
				addr2: {Balance: big.NewInt(1000000)},
				addr3: {Balance: big.NewInt(1000000)},
			},
		}
		genesis = gspec.MustCommit(db)
		signer  = types.NewHubbleSigner(gspec.Config.ChainID)
	)

	// Create three transactions shared between the chains:
	//  - shared: transaction included at both chain
	//  - removed: transaction included at the forked chain, will be removed
	//  - added : transaction included at the main chain
	shared, _ := types.SignTx(types.NewTransaction(0, addr1, big.NewInt(1000), params.TxGas, nil, nil), signer, key1)
	removed, _ := types.SignTx(types.NewTransaction(1, addr1, big.NewInt(1000), params.TxGas, nil, nil), signer, key1)
	added, _ := types.SignTx(types.NewTransaction(0, addr2, big.NewInt(1000), params.TxGas, nil, nil), signer, key2)

	chain, _ := GenerateChain(gspec.Config, genesis, mock.NewMock(), db, 3, func(i int, gen *BlockGen) {
		switch i {
		case 2:
			gen.AddTx(shared)
			gen.AddTx(removed)

			gen.OffsetTime(2) // Later block timestamp to make it a forked chain
		}
	})
	// Import the chain. This runs all block validation rules.
	blockchain, _ := NewBlockChain(db, nil, gspec.Config, mock.NewMock(), vm.Config{})
	if i, err := blockchain.InsertChain(chain); err != nil {
		t.Fatalf("failed to insert original fork chain[%d]: %v", i, err)
	}
	defer blockchain.Stop()

	// generate main chain, overwrite the old chain
	chain, _ = GenerateChain(gspec.Config, genesis, mock.NewMock(), db, 3, func(i int, gen *BlockGen) {
		switch i {
		case 2:
			gen.AddTx(shared)
			gen.AddTx(added)
		}
	})
	if _, err := blockchain.InsertChain(chain); err != nil {
		t.Fatalf("failed to insert main chain: %v", err)
	}

	// removed tx
	for i, tx := range (types.Transactions{removed}) {
		if txn, _, _, _ := rawdb.ReadTransaction(db, tx.Hash()); txn != nil {
			t.Errorf("drop %d: tx %v found while shouldn't have been", i, txn)
		}
		if rcpt, _, _, _ := rawdb.ReadReceipt(db, tx.Hash()); rcpt != nil {
			t.Errorf("drop %d: receipt %v found while shouldn't have been", i, rcpt)
		}
	}
	// added tx
	for i, tx := range (types.Transactions{added}) {
		if txn, _, _, _ := rawdb.ReadTransaction(db, tx.Hash()); txn == nil {
			t.Errorf("add %d: expected tx to be found", i)
		}
		if rcpt, _, _, _ := rawdb.ReadReceipt(db, tx.Hash()); rcpt == nil {
			t.Errorf("add %d: expected receipt to be found", i)
		}
	}
	// shared tx
	for i, tx := range (types.Transactions{shared}) {
		if txn, _, _, _ := rawdb.ReadTransaction(db, tx.Hash()); txn == nil {
			t.Errorf("share %d: expected tx to be found", i)
		}
		if rcpt, _, _, _ := rawdb.ReadReceipt(db, tx.Hash()); rcpt == nil {
			t.Errorf("share %d: expected receipt to be found", i)
		}
	}
}

func TestLogReorgs(t *testing.T) {
	var (
		key1, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr1   = crypto.PubkeyToAddress(key1.PublicKey)
		db      = vntdb.NewMemDatabase()
		// this code generates a log
		code    = common.Hex2Bytes("0161736db9052f0100789c9c56db6f145518ff9d397ba6ed0eddeeb260b90418168818ebf6226db0722b21624d480c589110b29ddd39dd8eccced499d94221ddaec5101e7df4c542bcbc181f8cf1414a4cfc177cf1091e144d4c88461a0c2185c47c33b317290fc47998effb7ee7fbfdbe739b39e75157f966e71730fc0a03c0328949a556c324ab611298640b500cc1a5339b3eed59813c6d05d36fbb9613480f20b8e7a434cc35a83a669ac70d1f8c82ecb8630596615b97e4bb866719455bfa50a84578d2f5cae03cc11545490836c318139c415137f20536f6f1b584b68086b9daa9751c116a45565c6f4e8156284c4b63a650347cc959b250308dc02848c7e44ada9425dbf0a4f946d5290596eb401527c3521dc9fd5c81b61d477114e97abd5e8fbc0f232fbd4856db80f41572c6ead7ea75a43f0a41ad9b61acdea9a951a7f1806b795316abe582e54cb9bf3120017a12348f58875d881e37b60300146acb003801804fd06b32e40800db01fc486b70577c42e61c536f84d606a07e13ba3e65777c154a75d600dc2478a0eb0ac1df511e4b528de3007ea2a67bb719a06e634d776bd3d57690fb27b9eba83fc3007a88f357d429844a53007e2774a5a5b4d2525a6929fdd354ea059022ce83b631862def346a3c6caf4195ef13badaaab1daaa11badd27e279644f626a8aa847003c21148cb27a426eec6f6df9e9b0874ae86748691371388ba5c235530074314049247603d8006023856a38d177883004a453f1aa578c122d3c709575c790673865e9d3fa6601ec00f0325a4f1f803c80210087da70aabd1e00ed815d71fc666ccbb11d89ed586c3900eae4fe361d3c6e76c428163d39cb326c4f2a2b782addbb35956113bd80926719b65e1dcc025c804c220f42786a547b5d3bbc19107916052f6e3ebc59d700552002c6b340c76ef0d4214dd380ce17309e8d5a80ae1e507332d26f088c674951dbd726b06e1b1bcc52ce7bda190de8a670a2370a526d6cd160f7c4e505b1d3fbda824c934b53876f138d0fd2b61c9953e883fc8ce68bb1d554371916cd1eeb9ff0a5e7f75f90966f5897a60da7dc7fcc2d552bd209fcfea21b04b67c25907ed05f729dc0334a818fe7a614ab966db61131eb04b655cc4f53e5f0275498f164c9adcc58b6cc97c25d0711fe18f81d30d191dc29c4b4caaf2f2c0aa6f2bbd01506c620141a24bf4b3fc8e499eb57042b371a42ee3dc6444ff2a0d06e893dea8a183c20f62e890e95dfaf2d0a7188f4960557f9df4c57781b6d8531a1250f0a5dedd429fd5e4d1762204ad754fe90e94257df123d6a9fd06f88017551f42e096d8a52ffa84db52b3f7a4a793554de29b6a83f53f6af355d24f97596175c5de28b5446e5bf80c04f19092cfd673875a5399c2deab2d075b12b2f869649e9764d179c7fc9f22a7f3c7fa021c3f9e76c39ec40428945360a96dc19aa86107d2f9b044be684280b762b2af7b5928c378d1f7825db70cafaacf47ccb75f4fdf981fc80be37f0aace79fdd57d83af8d0cbcf4bff741ff3396fe79b570527e50b53c8992eb9856789e158ebaae8d8a5f4669daf0e0079ee594b1e6d89b181a1e39654cc913551b55cb098686479ad8316bb6e99faa169bfe9869e2bc9c9b9b7582f787e766a23d8b395cc425e072c2312a92059d58731f604f5f0594f816c09f7501484467bfc8140a170cbf522819b65d2805aee7ab6b4eef8ee8f0fefe87b39773543e379a0b815c5faee43a7e6038416e74cab07dd997b39c996ae0e746cf9eebcbb9d5a015047333440c095e952ae5e6fb9a7a5eac6738ae335771abfe5ac1666ec5a7cc582f9af91ce599f2a23463de7caba29c954e909b3f57ff170000ffff")
		gspec   = &Genesis{Config: params.TestChainConfig, Alloc: GenesisAlloc{addr1: {Balance: big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18))}}}
		genesis = gspec.MustCommit(db)
		signer  = types.NewHubbleSigner(gspec.Config.ChainID)
	)

	blockchain, _ := NewBlockChain(db, nil, gspec.Config, mock.NewMock(), vm.Config{})
	defer blockchain.Stop()

	signTx := func(tx *types.Transaction) (*types.Transaction, error) {
		return types.SignTx(tx, signer, key1)
	}
	rmLogsCh := make(chan RemovedLogsEvent)
	blockchain.SubscribeRemovedLogsEvent(rmLogsCh)
	chain, _ := GenerateChain(params.TestChainConfig, genesis, mock.NewMock(), db, 3, func(i int, gen *BlockGen) {
		switch i {
		case 0:
			if err := StartFakeMainNet(gen, addr1, signTx); err != nil {
				t.Fatalf(err.Error())
			}
		case 2:
			if !election.MainNetActive(gen.statedb) {
				t.Fatalf("main is inactive")
			}
			gen.OffsetTime(2)

			tx, err := types.SignTx(types.NewContractCreation(gen.TxNonce(addr1), new(big.Int), 1000000, new(big.Int), code), signer, key1)
			if err != nil {
				t.Fatalf("failed to create tx: %v", err)
			}
			gen.AddTx(tx)
		}
	})
	if _, err := blockchain.InsertChain(chain); err != nil {
		t.Fatalf("failed to insert chain: %v", err)
	}

	// fork chain should to be canonical chain
	chain, _ = GenerateChain(params.TestChainConfig, genesis, mock.NewMock(), db, 3, func(i int, gen *BlockGen) {
		if i == 0 {
			if err := StartFakeMainNet(gen, addr1, signTx); err != nil {
				t.Fatalf(err.Error())
			}
		}
	})
	if _, err := blockchain.InsertChain(chain); err != nil {
		t.Fatalf("failed to insert forked chain: %v", err)
	}

	timeout := time.NewTimer(1 * time.Second)
	select {
	case ev := <-rmLogsCh:
		if len(ev.Logs) == 0 {
			t.Error("expected logs")
		}
	case <-timeout.C:
		t.Fatal("Timeout. There is no RemovedLogsEvent has been sent.")
	}
}

// Tests if the canonical block can be fetched from the database during chain insertion.
func TestCanonicalBlockRetrieval(t *testing.T) {
	_, blockchain, err := newCanonical(mock.NewMock(), 0, true)
	if err != nil {
		t.Fatalf("failed to create pristine chain: %v", err)
	}
	defer blockchain.Stop()

	chain, _ := GenerateChain(blockchain.chainConfig, blockchain.genesisBlock, mock.NewMock(), blockchain.db, 10, func(i int, gen *BlockGen) {})

	var pend sync.WaitGroup
	pend.Add(len(chain))

	for i := range chain {
		go func(block *types.Block) {
			defer pend.Done()

			// try to retrieve a block by its canonical hash and see if the block data can be retrieved.
			for {
				ch := rawdb.ReadCanonicalHash(blockchain.db, block.NumberU64())
				if ch == (common.Hash{}) {
					continue // busy wait for canonical hash to be written
				}
				if ch != block.Hash() {
					t.Fatalf("unknown canonical hash, want %s, got %s", block.Hash().Hex(), ch.Hex())
				}
				fb := rawdb.ReadBlock(blockchain.db, ch, block.NumberU64())
				if fb == nil {
					t.Fatalf("unable to retrieve block %d for canonical hash: %s", block.NumberU64(), ch.Hex())
				}
				if fb.Hash() != block.Hash() {
					t.Fatalf("invalid block hash for block %d, want %s, got %s", block.NumberU64(), block.Hash().Hex(), fb.Hash().Hex())
				}
				return
			}
		}(chain[i])

		if _, err := blockchain.InsertChain(types.Blocks{chain[i]}); err != nil {
			t.Fatalf("failed to insert block %d: %v", i, err)
		}
	}
	pend.Wait()
}

func TestEIP155Transition(t *testing.T) {
	// Configure and generate a sample block chain
	var (
		db         = vntdb.NewMemDatabase()
		key, _     = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address    = crypto.PubkeyToAddress(key.PublicKey)
		funds      = big.NewInt(1000000000)
		deleteAddr = common.Address{1}
		gspec      = &Genesis{
			Config: &params.ChainConfig{ChainID: big.NewInt(1), HubbleBlock: new(big.Int)},
			Alloc:  GenesisAlloc{address: {Balance: funds}, deleteAddr: {Balance: new(big.Int)}},
		}
		genesis = gspec.MustCommit(db)
	)

	blockchain, _ := NewBlockChain(db, nil, gspec.Config, mock.NewMock(), vm.Config{})
	defer blockchain.Stop()

	blocks, _ := GenerateChain(gspec.Config, genesis, mock.NewMock(), db, 4, func(i int, block *BlockGen) {
		var (
			tx      *types.Transaction
			err     error
			basicTx = func(signer types.Signer) (*types.Transaction, error) {
				return types.SignTx(types.NewTransaction(block.TxNonce(address), common.Address{}, new(big.Int), 21000, new(big.Int), nil), signer, key)
			}
		)
		switch i {
		case 0:
			tx, err = basicTx(types.NewHubbleSigner(gspec.Config.ChainID))
			if err != nil {
				t.Fatal(err)
			}
			block.AddTx(tx)
		case 2:
			tx, err = basicTx(types.NewHubbleSigner(gspec.Config.ChainID))
			if err != nil {
				t.Fatal(err)
			}
			block.AddTx(tx)

			tx, err = basicTx(types.NewHubbleSigner(gspec.Config.ChainID))
			if err != nil {
				t.Fatal(err)
			}
			block.AddTx(tx)
		case 3:
			tx, err = basicTx(types.NewHubbleSigner(gspec.Config.ChainID))
			if err != nil {
				t.Fatal(err)
			}
			block.AddTx(tx)

			tx, err = basicTx(types.NewHubbleSigner(gspec.Config.ChainID))
			if err != nil {
				t.Fatal(err)
			}
			block.AddTx(tx)
		}
	})

	if _, err := blockchain.InsertChain(blocks); err != nil {
		t.Fatal(err)
	}
	block := blockchain.GetBlockByNumber(1)
	if !block.Transactions()[0].Protected() {
		t.Error("Expected block[0].txs[0] to be replay protected")
	}

	block = blockchain.GetBlockByNumber(3)
	if !block.Transactions()[0].Protected() {
		t.Error("Expected block[3].txs[0] to be replay protected")
	}
	if !block.Transactions()[1].Protected() {
		t.Error("Expected block[3].txs[1] to be replay protected")
	}
	if _, err := blockchain.InsertChain(blocks[4:]); err != nil {
		t.Fatal(err)
	}

	// generate an invalid chain id transaction
	config := &params.ChainConfig{ChainID: big.NewInt(2), HubbleBlock: new(big.Int)}
	blocks, _ = GenerateChain(config, blocks[len(blocks)-1], mock.NewMock(), db, 4, func(i int, block *BlockGen) {
		var (
			tx      *types.Transaction
			err     error
			basicTx = func(signer types.Signer) (*types.Transaction, error) {
				return types.SignTx(types.NewTransaction(block.TxNonce(address), common.Address{}, new(big.Int), 21000, new(big.Int), nil), signer, key)
			}
		)
		switch i {
		case 0:
			tx, err = basicTx(types.NewHubbleSigner(big.NewInt(2)))
			if err != nil {
				t.Fatal(err)
			}
			block.AddTx(tx)
		}
	})
	_, err := blockchain.InsertChain(blocks)
	if err != types.ErrInvalidChainId {
		t.Error("expected error:", types.ErrInvalidChainId)
	}
}

func TestEIP161AccountRemoval(t *testing.T) {
	// Configure and generate a sample block chain
	var (
		db      = vntdb.NewMemDatabase()
		key, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address = crypto.PubkeyToAddress(key.PublicKey)
		funds   = big.NewInt(1000000000)
		theAddr = common.Address{1}
		gspec   = &Genesis{
			Config: &params.ChainConfig{
				ChainID:     big.NewInt(1),
				HubbleBlock: new(big.Int),
			},
			Alloc: GenesisAlloc{address: {Balance: funds}},
		}
		genesis = gspec.MustCommit(db)
	)
	blockchain, _ := NewBlockChain(db, nil, gspec.Config, mock.NewMock(), vm.Config{})
	defer blockchain.Stop()

	blocks, _ := GenerateChain(gspec.Config, genesis, mock.NewMock(), db, 3, func(i int, block *BlockGen) {
		var (
			tx     *types.Transaction
			err    error
			signer = types.NewHubbleSigner(gspec.Config.ChainID)
		)
		switch i {
		case 0:
			tx, err = types.SignTx(types.NewTransaction(block.TxNonce(address), theAddr, new(big.Int), 21000, new(big.Int), nil), signer, key)
		case 1:
			tx, err = types.SignTx(types.NewTransaction(block.TxNonce(address), theAddr, new(big.Int), 21000, new(big.Int), nil), signer, key)
		case 2:
			tx, err = types.SignTx(types.NewTransaction(block.TxNonce(address), theAddr, new(big.Int), 21000, new(big.Int), nil), signer, key)
		}
		if err != nil {
			t.Fatal(err)
		}
		block.AddTx(tx)
	})

	// account must exist pre eip 161
	if _, err := blockchain.InsertChain(types.Blocks{blocks[0]}); err != nil {
		t.Fatal(err)
	}
	if st, _ := blockchain.State(); st.Exist(theAddr) {
		t.Error("account should not exist")
	}

	// account needs to be deleted post eip 161
	if _, err := blockchain.InsertChain(types.Blocks{blocks[1]}); err != nil {
		t.Fatal(err)
	}
	if st, _ := blockchain.State(); st.Exist(theAddr) {
		t.Error("account should not exist")
	}

	// account musn't be created post eip 161
	if _, err := blockchain.InsertChain(types.Blocks{blocks[2]}); err != nil {
		t.Fatal(err)
	}
	if st, _ := blockchain.State(); st.Exist(theAddr) {
		t.Error("account should not exist")
	}
}

// Benchmarks large blocks with value transfers to non-existing accounts
func benchmarkLargeNumberOfValueToNonexisting(b *testing.B, numTxs, numBlocks int, recipientFn func(uint64) common.Address, dataFn func(uint64) []byte) {
	var (
		signer          = types.HubbleSigner{}
		testBankKey, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		testBankAddress = crypto.PubkeyToAddress(testBankKey.PublicKey)
		bankFunds       = big.NewInt(100000000000000000)
		gspec           = Genesis{
			Config: params.TestChainConfig,
			Alloc: GenesisAlloc{
				testBankAddress: {Balance: bankFunds},
				common.HexToAddress("0xc0de"): {
					Code:    []byte{0x60, 0x01, 0x50},
					Balance: big.NewInt(0),
				}, // push 1, pop
			},
			GasLimit: 100e6, // 100 M
		}
	)
	// Generate the original common chain segment and the two competing forks
	engine := mock.NewMock()
	db := vntdb.NewMemDatabase()
	genesis := gspec.MustCommit(db)

	blockGenerator := func(i int, block *BlockGen) {
		block.SetCoinbase(common.Address{1})
		for txi := 0; txi < numTxs; txi++ {
			uniq := uint64(i*numTxs + txi)
			recipient := recipientFn(uniq)
			//recipient := common.BigToAddress(big.NewInt(0).SetUint64(1337 + uniq))
			tx, err := types.SignTx(types.NewTransaction(uniq, recipient, big.NewInt(1), params.TxGas, big.NewInt(1), nil), signer, testBankKey)
			if err != nil {
				b.Error(err)
			}
			block.AddTx(tx)
		}
	}

	shared, _ := GenerateChain(params.TestChainConfig, genesis, engine, db, numBlocks, blockGenerator)
	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Import the shared chain and the original canonical one
		diskdb := vntdb.NewMemDatabase()
		gspec.MustCommit(diskdb)

		chain, err := NewBlockChain(diskdb, nil, params.TestChainConfig, engine, vm.Config{})
		if err != nil {
			b.Fatalf("failed to create tester chain: %v", err)
		}
		b.StartTimer()
		if _, err := chain.InsertChain(shared); err != nil {
			b.Fatalf("failed to insert shared chain: %v", err)
		}
		b.StopTimer()
		if got := chain.CurrentBlock().Transactions().Len(); got != numTxs*numBlocks {
			b.Fatalf("Transactions were not included, expected %d, got %d", (numTxs * numBlocks), got)

		}
	}
}
func BenchmarkBlockChain_1x1000ValueTransferToNonexisting(b *testing.B) {
	var (
		numTxs    = 1000
		numBlocks = 1
	)

	recipientFn := func(nonce uint64) common.Address {
		return common.BigToAddress(big.NewInt(0).SetUint64(1337 + nonce))
	}
	dataFn := func(nonce uint64) []byte {
		return nil
	}

	benchmarkLargeNumberOfValueToNonexisting(b, numTxs, numBlocks, recipientFn, dataFn)
}
func BenchmarkBlockChain_1x1000ValueTransferToExisting(b *testing.B) {
	var (
		numTxs    = 1000
		numBlocks = 1
	)
	b.StopTimer()
	b.ResetTimer()

	recipientFn := func(nonce uint64) common.Address {
		return common.BigToAddress(big.NewInt(0).SetUint64(1337))
	}
	dataFn := func(nonce uint64) []byte {
		return nil
	}

	benchmarkLargeNumberOfValueToNonexisting(b, numTxs, numBlocks, recipientFn, dataFn)
}
func BenchmarkBlockChain_1x1000Executions(b *testing.B) {
	var (
		numTxs    = 1000
		numBlocks = 1
	)
	b.StopTimer()
	b.ResetTimer()

	recipientFn := func(nonce uint64) common.Address {
		return common.BigToAddress(big.NewInt(0).SetUint64(0xc0de))
	}
	dataFn := func(nonce uint64) []byte {
		return nil
	}

	benchmarkLargeNumberOfValueToNonexisting(b, numTxs, numBlocks, recipientFn, dataFn)
}
