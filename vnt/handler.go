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
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/consensus"
	"github.com/vntchain/go-vnt/core"
	"github.com/vntchain/go-vnt/core/types"
	"github.com/vntchain/go-vnt/event"
	"github.com/vntchain/go-vnt/log"
	"github.com/vntchain/go-vnt/node"
	"github.com/vntchain/go-vnt/params"
	"github.com/vntchain/go-vnt/rlp"
	"github.com/vntchain/go-vnt/vnt/downloader"
	"github.com/vntchain/go-vnt/vnt/fetcher"
	"github.com/vntchain/go-vnt/vntdb"
	"github.com/vntchain/go-vnt/vntp2p"

	libp2p "github.com/libp2p/go-libp2p-peer"
)

const (
	softResponseLimit = 2 * 1024 * 1024 // Target maximum size of returned blocks, headers or node data.
	estHeaderRlpSize  = 500             // Approximate size of an RLP encoded block header

	// txChanSize is the size of channel listening to NewTxsEvent.
	// The number is referenced from the size of tx pool.
	txChanSize = 4096
)

// errIncompatibleConfig is returned if the requested protocols and configs are
// not compatible (low protocol version restrictions and high requirements).
var errIncompatibleConfig = errors.New("incompatible configuration")

func errResp(code errCode, format string, v ...interface{}) error {
	return fmt.Errorf("%v - %v", code, fmt.Sprintf(format, v...))
}

type ProtocolManager struct {
	networkId uint64

	fastSync  uint32 // Flag whether fast sync is enabled (gets disabled if we already have blocks)
	acceptTxs uint32 // Flag whether we're considered synchronised (enables transaction processing)

	txpool      txPool
	blockchain  *core.BlockChain
	chainconfig *params.ChainConfig
	maxPeers    int

	downloader *downloader.Downloader
	fetcher    *fetcher.Fetcher
	peers      *peerSet
	node       *node.Node

	SubProtocols []vntp2p.Protocol

	eventMux      *event.TypeMux
	txsCh         chan core.NewTxsEvent
	txsSub        event.Subscription
	minedBlockSub *event.TypeMuxSubscription
	bftMsgSub     *event.TypeMuxSubscription
	bftPeerSub    *event.TypeMuxSubscription

	// channels for fetcher, syncer, txsyncLoop
	newPeerCh   chan *peer
	txsyncCh    chan *txsync
	quitSync    chan struct{}
	noMorePeers chan struct{}

	// wait group is used for graceful shutdowns during downloading
	// and processing
	wg sync.WaitGroup

	urlsCh chan []string // 传递p2p urls of witnesses
}

// NewProtocolManager returns a new VNT sub protocol manager. The VNT sub protocol manages peers capable
// with the VNT network.
func NewProtocolManager(config *params.ChainConfig, mode downloader.SyncMode, networkId uint64, mux *event.TypeMux, txpool txPool, engine consensus.Engine, blockchain *core.BlockChain, chaindb vntdb.Database, node *node.Node) (*ProtocolManager, error) {
	// Create the protocol manager with the base fields
	manager := &ProtocolManager{
		networkId:   networkId,
		eventMux:    mux,
		txpool:      txpool,
		blockchain:  blockchain,
		chainconfig: config,
		peers:       newPeerSet(),
		newPeerCh:   make(chan *peer),
		noMorePeers: make(chan struct{}),
		txsyncCh:    make(chan *txsync),
		quitSync:    make(chan struct{}),
		node:        node,

		urlsCh: make(chan []string),
	}
	// Figure out whether to allow fast sync or not
	if mode == downloader.FastSync && blockchain.CurrentBlock().NumberU64() > 0 {
		log.Warn("Blockchain not empty, fast sync disabled")
		mode = downloader.FullSync
	}
	if mode == downloader.FastSync {
		manager.fastSync = uint32(1)
	}
	// Initiate a sub-protocol for every implemented version we can handle
	manager.SubProtocols = make([]vntp2p.Protocol, 0, len(ProtocolVersions))
	for i, version := range ProtocolVersions {
		// Skip protocol version if incompatible with the mode of operation
		if mode == downloader.FastSync && version < vnt63 {
			continue
		}
		// Compatible; initialise the sub-protocol
		version := version // Closure for the run
		manager.SubProtocols = append(manager.SubProtocols, vntp2p.Protocol{
			Name:    ProtocolName,
			Version: version,
			Length:  ProtocolLengths[i],
			Run: func(p *vntp2p.Peer, rw vntp2p.MsgReadWriter) error {
				peer := manager.newPeer(int(version), p, rw)
				select {
				case manager.newPeerCh <- peer:
					manager.wg.Add(1)
					defer manager.wg.Done()
					return manager.handle(peer)
				case <-manager.quitSync:
					return vntp2p.DiscQuitting
				}
			},
			NodeInfo: func() interface{} {
				return manager.NodeInfo()
			},
			PeerInfo: func(id libp2p.ID) interface{} {
				if p := manager.peers.Peer(id); p != nil {
					return p.Info()
				}
				return nil
			},
		})
	}
	if len(manager.SubProtocols) == 0 {
		return nil, errIncompatibleConfig
	}
	// Construct the different synchronisation mechanisms
	manager.downloader = downloader.New(mode, chaindb, manager.eventMux, blockchain, nil, manager.removePeer)

	validator := func(header *types.Header) error {
		return engine.VerifyHeader(blockchain, header, true)
	}
	heighter := func() uint64 {
		return blockchain.CurrentBlock().NumberU64()
	}
	inserter := func(blocks types.Blocks) (int, error) {
		// If fast sync is running, deny importing weird blocks
		if atomic.LoadUint32(&manager.fastSync) == 1 {
			log.Warn("Discarded bad propagated block", "number", blocks[0].Number(), "hash", blocks[0].Hash())
			return 0, nil
		}
		atomic.StoreUint32(&manager.acceptTxs, 1) // Mark initial sync done on any fetcher import
		return manager.blockchain.InsertChain(blocks)
	}
	manager.fetcher = fetcher.New(blockchain.GetBlockByHash, validator, manager.BroadcastBlock, heighter, inserter, manager.removePeer)

	return manager, nil
}

func (pm *ProtocolManager) removePeer(id libp2p.ID) {
	// Short circuit if the peer was already removed
	peer := pm.peers.Peer(id)
	if peer == nil {
		return
	}
	log.Debug("Removing VNT peer", "peer", id)

	// Unregister the peer from the downloader and VNT peer set
	pm.downloader.UnregisterPeer(id)
	if err := pm.peers.Unregister(id); err != nil {
		log.Error("Peer removal failed", "peer", id, "err", err)
	}
	// Hard disconnect at the networking layer
	if peer != nil {
		peer.Peer.Disconnect(vntp2p.DiscUselessPeer)
	}
}

// resetBftPeer update current bft peer connection. If node not has connection
// will them, will connecting to them. url format is:
// /ip4/192.168.102.2/tcp/5216/ipfs/1kHBzN17vVE75rwZA7vKAFfxUYS8XMh6QBYS6JWF13xHGX9
func (pm *ProtocolManager) resetBftPeer(urls []string) {
	pm.peers.lock.Lock()
	defer pm.peers.lock.Unlock()

	// Clean old records
	pm.peers.bftPeers = make(map[libp2p.ID]struct{})

	// Add new records, and addPeer if not connect
	selfID := pm.node.Server().NodeInfo().ID
	for _, url := range urls {
		node, err := vntp2p.ParseNode(url)
		if err != nil {
			log.Error("resetBftPeer invalid vnode:", "error", err)
			continue
		}
		if node.Id.ToString() == selfID {
			continue
		}

		pm.peers.bftPeers[node.Id] = struct{}{}
		if _, exists := pm.peers.peers[node.Id]; !exists {
			log.Debug("Reset bft peer, connecting to", "peer", url)
			go pm.node.Server().AddPeer(context.Background(), node)
		}
	}
}

func (pm *ProtocolManager) Start(maxPeers int) {
	pm.maxPeers = maxPeers

	// broadcast transactions
	pm.txsCh = make(chan core.NewTxsEvent, txChanSize)
	pm.txsSub = pm.txpool.SubscribeNewTxsEvent(pm.txsCh)
	go pm.txBroadcastLoop()

	// broadcast mined blocks
	pm.minedBlockSub = pm.eventMux.Subscribe(core.NewMinedBlockEvent{})
	pm.bftMsgSub = pm.eventMux.Subscribe(core.SendBftMsgEvent{})
	pm.bftPeerSub = pm.eventMux.Subscribe(core.BftPeerChangeEvent{})
	go pm.minedBroadcastLoop()
	go pm.bftBroadcastLoop()

	go pm.resetBftPeerLoop()

	// start sync handlers
	go pm.syncer()
	go pm.txsyncLoop()
	go pm.bftPeerLoop()
}

func (pm *ProtocolManager) Stop() {
	log.Info("Stopping VNT protocol")

	pm.txsSub.Unsubscribe()        // quits txBroadcastLoop
	pm.minedBlockSub.Unsubscribe() // quits blockBroadcastLoop
	pm.bftMsgSub.Unsubscribe()
	pm.bftPeerSub.Unsubscribe()

	close(pm.urlsCh)

	// Quit the sync loop.
	// After this send has completed, no new peers will be accepted.
	pm.noMorePeers <- struct{}{}

	// Quit fetcher, txsyncLoop.
	close(pm.quitSync)

	// Disconnect existing sessions.
	// This also closes the gate for any new registrations on the peer set.
	// sessions which are already established but not added to pm.peers yet
	// will exit when they try to register.
	pm.peers.Close()

	// Wait for all peer handler goroutines and the loops to come down.
	pm.wg.Wait()

	log.Info("VNT protocol stopped")
}

func (pm *ProtocolManager) newPeer(pv int, p *vntp2p.Peer, rw vntp2p.MsgReadWriter) *peer {
	return newPeer(pv, p, newMeteredMsgWriter(rw))
}

// handle is the callback invoked to manage the life cycle of an vnt peer. When
// this function terminates, the peer is disconnected.
func (pm *ProtocolManager) handle(p *peer) error {
	// Ignore maxPeers if this is a trusted peer

	// if pm.peers.Len() >= pm.maxPeers && !p.Peer.Info().Network.Trusted {
	// 	return vntp2p.DiscTooManyPeers
	// }

	// p.Log().Debug("VNT peer connected", "name", p.Name())

	// Execute the VNT handshake
	var (
		genesis = pm.blockchain.Genesis()
		head    = pm.blockchain.CurrentHeader()
		hash    = head.Hash()
		number  = head.Number.Uint64()
		td      = pm.blockchain.GetTd(hash, number)
	)
	if err := p.Handshake(pm.networkId, td, hash, genesis.Hash()); err != nil {
		p.Log().Debug("VNT handshake failed", "err", err)
		return err
	}
	if rw, ok := p.rw.(*meteredMsgReadWriter); ok {
		rw.Init(p.version)
	}
	// Register the peer locally
	if err := pm.peers.Register(p); err != nil {
		p.Log().Error("VNT peer registration failed", "err", err)
		return err
	}
	defer pm.removePeer(p.id)

	// Register the peer in the downloader. If the downloader considers it banned, we disconnect
	if err := pm.downloader.RegisterPeer(p.id, p.version, p); err != nil {
		return err
	}
	// Propagate existing transactions. new transactions appearing
	// after this will be sent via broadcasts.
	pm.syncTransactions(p)

	// main loop. handle incoming messages.
	for {
		if err := pm.handleMsg(p); err != nil {
			p.Log().Debug("VNT message handling failed", "err", err)
			return err
		}
	}
}

// handleMsg is invoked whenever an inbound message is received from a remote
// peer. The remote connection is torn down upon returning any error.
func (pm *ProtocolManager) handleMsg(p *peer) error {
	// Read the next message from the remote peer, and ensure it's fully consumed
	msg, err := p.rw.ReadMsg()
	if err != nil {
		return err
	}
	size := msg.GetBodySize()
	if size > ProtocolMaxMsgSize {
		return errResp(ErrMsgTooLarge, "%v > %v", size, ProtocolMaxMsgSize)
	}

	//按理说，新版的协议处理方式，不会有残留数据得不到处理
	//defer msg.Discard()

	// Handle the message depending on its contents
	switch {
	case msg.Body.Type == StatusMsg:
		// Status messages should never arrive after the handshake
		return errResp(ErrExtraStatusMsg, "uncontrolled status message")

	// Block header query, collect the requested headers and reply
	case msg.Body.Type == GetBlockHeadersMsg:
		// Decode the complex header query
		var query getBlockHeadersData
		if err := msg.Decode(&query); err != nil {
			return errResp(ErrDecode, "%v: %v", msg, err)
		}
		hashMode := query.Origin.Hash != (common.Hash{})
		first := true
		maxNonCanonical := uint64(100)

		// Gather headers until the fetch or network limits is reached
		var (
			bytes   common.StorageSize
			headers []*types.Header
			unknown bool
		)
		for !unknown && len(headers) < int(query.Amount) && bytes < softResponseLimit && len(headers) < downloader.MaxHeaderFetch {
			// Retrieve the next header satisfying the query
			var origin *types.Header
			if hashMode {
				if first {
					first = false
					origin = pm.blockchain.GetHeaderByHash(query.Origin.Hash)
					if origin != nil {
						query.Origin.Number = origin.Number.Uint64()
					}
				} else {
					origin = pm.blockchain.GetHeader(query.Origin.Hash, query.Origin.Number)
				}
			} else {
				origin = pm.blockchain.GetHeaderByNumber(query.Origin.Number)
			}
			if origin == nil {
				break
			}
			headers = append(headers, origin)
			bytes += estHeaderRlpSize

			// Advance to the next header of the query
			switch {
			case hashMode && query.Reverse:
				// Hash based traversal towards the genesis block
				ancestor := query.Skip + 1
				if ancestor == 0 {
					unknown = true
				} else {
					query.Origin.Hash, query.Origin.Number = pm.blockchain.GetAncestor(query.Origin.Hash, query.Origin.Number, ancestor, &maxNonCanonical)
					unknown = (query.Origin.Hash == common.Hash{})
				}
			case hashMode && !query.Reverse:
				// Hash based traversal towards the leaf block
				var (
					current = origin.Number.Uint64()
					next    = current + query.Skip + 1
				)
				if next <= current {
					infos, _ := json.MarshalIndent(p.Peer.Info(), "", "  ")
					p.Log().Warn("GetBlockHeaders skip overflow attack", "current", current, "skip", query.Skip, "next", next, "attacker", infos)
					unknown = true
				} else {
					if header := pm.blockchain.GetHeaderByNumber(next); header != nil {
						nextHash := header.Hash()
						expOldHash, _ := pm.blockchain.GetAncestor(nextHash, next, query.Skip+1, &maxNonCanonical)
						if expOldHash == query.Origin.Hash {
							query.Origin.Hash, query.Origin.Number = nextHash, next
						} else {
							unknown = true
						}
					} else {
						unknown = true
					}
				}
			case query.Reverse:
				// Number based traversal towards the genesis block
				if query.Origin.Number >= query.Skip+1 {
					query.Origin.Number -= query.Skip + 1
				} else {
					unknown = true
				}

			case !query.Reverse:
				// Number based traversal towards the leaf block
				query.Origin.Number += query.Skip + 1
			}
		}
		return p.SendBlockHeaders(headers)

	case msg.Body.Type == BlockHeadersMsg:
		// A batch of headers arrived to one of our previous requests
		var headers []*types.Header
		if err := msg.Decode(&headers); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Filter out any explicitly requested headers, deliver the rest to the downloader
		filter := len(headers) == 1
		if filter {
			// Irrelevant of the fork checks, send the header to the fetcher just in case
			headers = pm.fetcher.FilterHeaders(p.id, headers, time.Now())
		}
		if len(headers) > 0 || !filter {
			err := pm.downloader.DeliverHeaders(p.id, headers)
			if err != nil {
				log.Debug("Failed to deliver headers", "err", err)
			}
		}

	case msg.Body.Type == GetBlockBodiesMsg:
		// Decode the retrieval message
		msgStream := rlp.NewStream(msg.Body.Payload, uint64(msg.Body.PayloadSize))
		if _, err := msgStream.List(); err != nil {
			return err
		}
		// Gather blocks until the fetch or network limits is reached
		var (
			hash   common.Hash
			bytes  int
			bodies []rlp.RawValue
		)
		for bytes < softResponseLimit && len(bodies) < downloader.MaxBlockFetch {
			// Retrieve the hash of the next block
			if err := msgStream.Decode(&hash); err == rlp.EOL {
				break
			} else if err != nil {
				return errResp(ErrDecode, "msg %v: %v", msg, err)
			}
			// Retrieve the requested block body, stopping if enough was found
			if data := pm.blockchain.GetBodyRLP(hash); len(data) != 0 {
				bodies = append(bodies, data)
				bytes += len(data)
			}
		}
		return p.SendBlockBodiesRLP(bodies)

	case msg.Body.Type == BlockBodiesMsg:
		// A batch of block bodies arrived to one of our previous requests
		var request blockBodiesData
		if err := msg.Decode(&request); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Deliver them all to the downloader for queuing
		transactions := make([][]*types.Transaction, len(request))

		for i, body := range request {
			transactions[i] = body.Transactions
		}
		// Filter out any explicitly requested bodies, deliver the rest to the downloader
		filter := len(transactions) > 0
		if filter {
			transactions = pm.fetcher.FilterBodies(p.id, transactions, time.Now())
		}
		err := pm.downloader.DeliverBodies(p.id, transactions)
		if err != nil {
			log.Debug("Failed to deliver bodies", "err", err)
		}

	case p.version >= vnt63 && msg.Body.Type == GetNodeDataMsg:
		// Decode the retrieval message
		msgStream := rlp.NewStream(msg.Body.Payload, uint64(msg.Body.PayloadSize))
		if _, err := msgStream.List(); err != nil {
			return err
		}
		// Gather state data until the fetch or network limits is reached
		var (
			hash  common.Hash
			bytes int
			data  [][]byte
		)
		for bytes < softResponseLimit && len(data) < downloader.MaxStateFetch {
			// Retrieve the hash of the next state entry
			if err := msgStream.Decode(&hash); err == rlp.EOL {
				break
			} else if err != nil {
				return errResp(ErrDecode, "msg %v: %v", msg, err)
			}
			// Retrieve the requested state entry, stopping if enough was found
			if entry, err := pm.blockchain.TrieNode(hash); err == nil {
				data = append(data, entry)
				bytes += len(entry)
			}
		}
		return p.SendNodeData(data)

	case p.version >= vnt63 && msg.Body.Type == NodeDataMsg:
		// A batch of node state data arrived to one of our previous requests
		var data [][]byte
		if err := msg.Decode(&data); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Deliver all to the downloader
		if err := pm.downloader.DeliverNodeData(p.id, data); err != nil {
			log.Debug("Failed to deliver node state data", "err", err)
		}

	case p.version >= vnt63 && msg.Body.Type == GetReceiptsMsg:
		// Decode the retrieval message
		msgStream := rlp.NewStream(msg.Body.Payload, uint64(msg.Body.PayloadSize))
		if _, err := msgStream.List(); err != nil {
			return err
		}
		// Gather state data until the fetch or network limits is reached
		var (
			hash     common.Hash
			bytes    int
			receipts []rlp.RawValue
		)
		for bytes < softResponseLimit && len(receipts) < downloader.MaxReceiptFetch {
			// Retrieve the hash of the next block
			if err := msgStream.Decode(&hash); err == rlp.EOL {
				break
			} else if err != nil {
				return errResp(ErrDecode, "msg %v: %v", msg, err)
			}
			// Retrieve the requested block's receipts, skipping if unknown to us
			results := pm.blockchain.GetReceiptsByHash(hash)
			if results == nil {
				if header := pm.blockchain.GetHeaderByHash(hash); header == nil || header.ReceiptHash != types.EmptyRootHash {
					continue
				}
			}
			// If known, encode and queue for response packet
			if encoded, err := rlp.EncodeToBytes(results); err != nil {
				log.Error("Failed to encode receipt", "err", err)
			} else {
				receipts = append(receipts, encoded)
				bytes += len(encoded)
			}
		}
		return p.SendReceiptsRLP(receipts)

	case p.version >= vnt63 && msg.Body.Type == ReceiptsMsg:
		// A batch of receipts arrived to one of our previous requests
		var receipts [][]*types.Receipt
		if err := msg.Decode(&receipts); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Deliver all to the downloader
		if err := pm.downloader.DeliverReceipts(p.id, receipts); err != nil {
			log.Debug("Failed to deliver receipts", "err", err)
		}

	case msg.Body.Type == NewBlockHashesMsg:
		var announces newBlockHashesData
		if err := msg.Decode(&announces); err != nil {
			return errResp(ErrDecode, "%v: %v", msg, err)
		}
		// Mark the hashes as present at the remote node
		for _, block := range announces {
			log.Debug("Receive announce", "high", block.Number, "hash", block.Hash, "parent", block.ParentHash, "parent td", block.ParentTD, "from", p.id)
			p.MarkBlock(block.Hash)
		}
		// Schedule all the unknown hashes for retrieval
		unknown := make(newBlockHashesData, 0, len(announces))
		for _, block := range announces {
			if !pm.blockchain.HasBlock(block.Hash, block.Number) {
				unknown = append(unknown, block)
			}
		}

		maxTd := big.NewInt(0)
		var maxHash common.Hash
		for _, block := range unknown {
			pm.fetcher.Notify(p.id, block.Hash, block.Number, time.Now(), p.RequestOneHeader, p.RequestBodies)
			if block.ParentTD.Cmp(maxTd) > 0 {
				maxTd = block.ParentTD
				maxHash = block.ParentHash
			}
		}
		// Update with peer when fall behind
		pm.updatePeerHeadAndSync(p, maxHash, maxTd)

	case msg.Body.Type == NewBlockMsg:
		// This message is forbid. The peer is malicious and will be removed.
		log.Info("Receive NewBlockMsg from", "peer", p.id)
		pm.removePeer(p.id)

	case msg.Body.Type == TxMsg:
		// Transactions arrived, make sure we have a valid and fresh chain to handle them
		if atomic.LoadUint32(&pm.acceptTxs) == 0 {
			break
		}
		// Transactions can be processed, parse all of them and deliver to the pool
		var txs []*types.Transaction
		if err := msg.Decode(&txs); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		for i, tx := range txs {
			// Validate and mark the remote transaction
			if tx == nil {
				return errResp(ErrDecode, "transaction %d is nil", i)
			}
			p.MarkTransaction(tx.Hash())
		}
		pm.txpool.AddRemotes(txs)
	case msg.Body.Type == BftPreprepareMsg:
		bftMsg := types.PreprepareMsg{}
		if err := msg.Decode(&bftMsg); err != nil {
			log.Error("Decode bftMsg Error", "err", err)
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		pm.postRecBftEvent(&bftMsg)
	case msg.Body.Type == BftPrepareMsg:
		bftMsg := types.PrepareMsg{}
		if err := msg.Decode(&bftMsg); err != nil {
			log.Error("Decode bftMsg Error", "err", err)
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		pm.postRecBftEvent(&bftMsg)
	case msg.Body.Type == BftCommitMsg:
		bftMsg := types.CommitMsg{}
		if err := msg.Decode(&bftMsg); err != nil {
			log.Error("Decode bftMsg Error", "err", err)
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		pm.postRecBftEvent(&bftMsg)
	default:
		return errResp(ErrInvalidMsgCode, "%v", msg.Body.Type)
	}
	return nil
}

// updatePeerHeadAndSync will update peer's head and start a sync if local fall behind of peer.
func (pm *ProtocolManager) updatePeerHeadAndSync(p *peer, parentHash common.Hash, parentTd *big.Int) {
	// Update the peers total difficulty if better than the previous
	if _, td := p.Head(); parentTd.Cmp(td) > 0 {
		p.SetHead(parentHash, parentTd)

		// Schedule a sync if above ours. Note, this will not fire a sync for a gap of
		// a singe block (as the true TD is below the propagated block), however this
		// scenario should easily be covered by the fetcher.
		currentBlock := pm.blockchain.CurrentBlock()
		currentTd := pm.blockchain.GetTd(currentBlock.Hash(), currentBlock.NumberU64())
		if parentTd.Cmp(currentTd) > 0 {
			log.Debug("updatePeerHeadAndSync: local behind peer", "local td", currentTd.Int64(), "parent td", parentTd.Int64())
			go pm.synchronise(p)
		}
	}
}

func (pm *ProtocolManager) postRecBftEvent(msg types.ConsensusMsg) {
	log.Debug("Post RecBftEvent", "type", msg.Type(),
		"h", msg.GetBlockNum(), "r", msg.GetRound(), "hash", msg.Hash())
	pm.eventMux.Post(core.RecBftMsgEvent{
		BftMsg: types.BftMsg{
			BftType: msg.Type(),
			Msg:     msg,
		},
	})
}

// BroadcastBlock will either propagate a block to a subset of it's peers, or
// will only announce it's availability (depending what's requested).
func (pm *ProtocolManager) BroadcastBlock(block *types.Block, propagate bool) {
	hash := block.Hash()
	peers := pm.peers.PeersWithoutBlock(hash)

	// Otherwise if the block is indeed in out own chain, announce it
	parentTD := pm.blockchain.GetTd(block.ParentHash(), block.NumberU64()-1)
	if pm.blockchain.HasBlock(hash, block.NumberU64()) {
		for _, peer := range peers {
			log.Debug("Broadcast announce", "high", block.NumberU64(), "hash", block.Hash(), "to peer", peer)
			peer.AsyncSendNewBlockHash(block, parentTD)
		}
		log.Trace("Announced block", "hash", hash, "recipients", len(peers), "duration", common.PrettyDuration(time.Since(block.ReceivedAt)))
	}
}

// BroadcastTxs will propagate a batch of transactions to all peers which are not known to
// already have the given transaction.
func (pm *ProtocolManager) BroadcastTxs(txs types.Transactions) {
	var txset = make(map[*peer]types.Transactions)

	// Broadcast transactions to a batch of peers not knowing about it
	for _, tx := range txs {
		peers := pm.peers.PeersWithoutTx(tx.Hash())
		for _, peer := range peers {
			txset[peer] = append(txset[peer], tx)
		}
		log.Trace("Broadcast transaction", "hash", tx.Hash(), "recipients", len(peers))
	}
	// FIXME include this again: peers = peers[:int(math.Sqrt(float64(len(peers))))]
	for peer, txs := range txset {
		peer.AsyncSendTransactions(txs)
	}
}

func (pm *ProtocolManager) BroadcastBftMsg(bftMsg types.BftMsg) {
	peers := pm.peers.PeersForBft()
	log.Trace("BroadcastBftMsg", "type", bftMsg.BftType, "hash", bftMsg.Msg.Hash(), "number of bft peer", len(peers))

	for _, p := range peers {
		// using goroutine for each peer for peer may connection
		go func(p *peer) {
			log.Trace("BroadcastBftMsg", "to peer", p.id.ToString())
			err := p.SendBftMsg(bftMsg)
			if err != nil {
				log.Error("BroadcastBftMsg error", "to peer", p.id.ToString(), "error", err)
			} else {
				log.Trace("BroadcastBftMsg success", "to peer", p.id.ToString())
			}
		}(p)
	}

	log.Trace("BroadcastBftMsg exit")
}

// Mined broadcast loop
func (pm *ProtocolManager) minedBroadcastLoop() {
	// automatically stops if unsubscribe
	for obj := range pm.minedBlockSub.Chan() {
		switch ev := obj.Data.(type) {
		case core.NewMinedBlockEvent:
			log.Debug("PM receive NewMinedBlockEvent")
			pm.BroadcastBlock(ev.Block, false) // Only announce all the peer
		}
	}
}

func (pm *ProtocolManager) txBroadcastLoop() {
	for {
		select {
		case event := <-pm.txsCh:
			pm.BroadcastTxs(event.Txs)

		// Err() channel will be closed when unsubscribing.
		case <-pm.txsSub.Err():
			return
		}
	}
}

func (pm *ProtocolManager) bftBroadcastLoop() {
	for obj := range pm.bftMsgSub.Chan() {
		switch ev := obj.Data.(type) {
		case core.SendBftMsgEvent:
			pm.BroadcastBftMsg(ev.BftMsg) // First propagate block to peers
		}
	}
}

func (pm *ProtocolManager) bftPeerLoop() {
	for obj := range pm.bftPeerSub.Chan() {
		switch ev := obj.Data.(type) {
		case core.BftPeerChangeEvent:
			log.Trace("Receive BftPeerChangeEvent")
			pm.urlsCh <- ev.Urls
			// pm.resetBftPeer(ev.Urls) // First propagate block to peers
		}
	}
}

// NodeInfo represents a short summary of the VNT sub-protocol metadata
// known about the host peer.
type NodeInfo struct {
	Network    uint64              `json:"network"`    // VNT network ID (1=Frontier)
	Difficulty *big.Int            `json:"difficulty"` // Total difficulty of the host's blockchain
	Genesis    common.Hash         `json:"genesis"`    // SHA3 hash of the host's genesis block
	Config     *params.ChainConfig `json:"config"`     // Chain configuration for the fork rules
	Head       common.Hash         `json:"head"`       // SHA3 hash of the host's best owned block
}

// NodeInfo retrieves some protocol metadata about the running host node.
func (pm *ProtocolManager) NodeInfo() *NodeInfo {
	currentBlock := pm.blockchain.CurrentBlock()
	return &NodeInfo{
		Network:    pm.networkId,
		Difficulty: pm.blockchain.GetTd(currentBlock.Hash(), currentBlock.NumberU64()),
		Genesis:    pm.blockchain.Genesis().Hash(),
		Config:     pm.blockchain.Config(),
		Head:       currentBlock.Hash(),
	}
}

func (pm *ProtocolManager) resetBftPeerLoop() {
	log.Debug("resetBftPeerLoop start")
	defer log.Debug("resetBftPeerLoop exit")

	for urls := range pm.urlsCh {
		pm.resetBftPeer(urls)
	}
}
