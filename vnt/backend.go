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

// Package vnt implements the VNT protocol.
package vnt

import (
	"errors"
	"fmt"
	"math/big"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/vntchain/go-vnt/accounts"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/common/hexutil"
	"github.com/vntchain/go-vnt/consensus"
	"github.com/vntchain/go-vnt/consensus/dpos"
	"github.com/vntchain/go-vnt/core"
	"github.com/vntchain/go-vnt/core/bloombits"
	"github.com/vntchain/go-vnt/core/rawdb"
	"github.com/vntchain/go-vnt/core/types"
	"github.com/vntchain/go-vnt/core/vm"
	"github.com/vntchain/go-vnt/event"
	"github.com/vntchain/go-vnt/internal/vntapi"
	"github.com/vntchain/go-vnt/log"
	"github.com/vntchain/go-vnt/miner"
	"github.com/vntchain/go-vnt/node"
	"github.com/vntchain/go-vnt/params"
	"github.com/vntchain/go-vnt/rlp"
	"github.com/vntchain/go-vnt/rpc"
	"github.com/vntchain/go-vnt/vnt/downloader"
	"github.com/vntchain/go-vnt/vnt/filters"
	"github.com/vntchain/go-vnt/vnt/gasprice"
	"github.com/vntchain/go-vnt/vntdb"
	"github.com/vntchain/go-vnt/vntp2p"
)

type LesServer interface {
	Start(srvr *vntp2p.Server)
	Stop()
	Protocols() []vntp2p.Protocol
	SetBloomBitsIndexer(bbIndexer *core.ChainIndexer)
}

// VNT implements the VNT full node service.
type VNT struct {
	config      *Config
	chainConfig *params.ChainConfig

	// Channel for shutting down the service
	shutdownChan chan bool // Channel for shutting down the VNT

	// Handlers
	txPool          *core.TxPool
	blockchain      *core.BlockChain
	protocolManager *ProtocolManager
	lesServer       LesServer

	// DB interfaces
	chainDb vntdb.Database // Block chain database

	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager

	bloomRequests chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer  *core.ChainIndexer             // Bloom indexer operating during block imports

	APIBackend *VntAPIBackend

	miner    *miner.Miner
	gasPrice *big.Int
	coinbase common.Address

	networkId     uint64
	netRPCService *vntapi.PublicNetAPI

	lock sync.RWMutex // Protects the variadic fields (e.g. gas price and coinbase)
}

func (s *VNT) AddLesServer(ls LesServer) {
	s.lesServer = ls
	ls.SetBloomBitsIndexer(s.bloomIndexer)
}

// New creates a new VNT object (including the
// initialisation of the common VNT object)
func New(ctx *node.ServiceContext, config *Config, node *node.Node) (*VNT, error) {
	if config.SyncMode == downloader.LightSync {
		return nil, errors.New("can't run vnt.VNT in light sync mode, use les.LightVnt")
	}
	if !config.SyncMode.IsValid() {
		return nil, fmt.Errorf("invalid sync mode %d", config.SyncMode)
	}
	chainDb, err := CreateDB(ctx, config, "chaindata")
	if err != nil {
		return nil, err
	}
	chainConfig, genesisHash, genesisErr := core.SetupGenesisBlock(chainDb, config.Genesis)
	if _, ok := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !ok {
		return nil, genesisErr
	}
	log.Info("Initialised chain configuration", "config", chainConfig)

	vnt := &VNT{
		config:         config,
		chainDb:        chainDb,
		chainConfig:    chainConfig,
		eventMux:       ctx.EventMux,
		accountManager: ctx.AccountManager,
		engine:         CreateConsensusEngine(ctx, chainConfig, chainDb),
		shutdownChan:   make(chan bool),
		networkId:      config.NetworkId,
		gasPrice:       config.GasPrice,
		coinbase:       config.Coinbase,
		bloomRequests:  make(chan chan *bloombits.Retrieval),
		bloomIndexer:   NewBloomIndexer(chainDb, params.BloomBitsBlocks),
	}

	log.Info("Initialising VNT protocol", "versions", ProtocolVersions, "network", config.NetworkId)

	if !config.SkipBcVersionCheck {
		bcVersion := rawdb.ReadDatabaseVersion(chainDb)
		if bcVersion != core.BlockChainVersion && bcVersion != 0 {
			return nil, fmt.Errorf("Blockchain DB version mismatch (%d / %d). Run gvnt upgradedb.\n", bcVersion, core.BlockChainVersion)
		}
		rawdb.WriteDatabaseVersion(chainDb, core.BlockChainVersion)
	}
	var (
		vmConfig    = vm.Config{EnablePreimageRecording: config.EnablePreimageRecording}
		cacheConfig = &core.CacheConfig{Disabled: config.NoPruning, TrieNodeLimit: config.TrieCache, TrieTimeLimit: config.TrieTimeout}
	)
	vnt.blockchain, err = core.NewBlockChain(chainDb, cacheConfig, vnt.chainConfig, vnt.engine, vmConfig)
	if err != nil {
		return nil, err
	}
	// Rewind the chain in case of an incompatible config upgrade.
	if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		log.Warn("Rewinding chain to upgrade configuration", "err", compat)
		vnt.blockchain.SetHead(compat.RewindTo)
		rawdb.WriteChainConfig(chainDb, genesisHash, chainConfig)
	}
	vnt.bloomIndexer.Start(vnt.blockchain)

	if config.TxPool.Journal != "" {
		config.TxPool.Journal = ctx.ResolvePath(config.TxPool.Journal)
	}
	vnt.txPool = core.NewTxPool(config.TxPool, vnt.chainConfig, vnt.blockchain)

	if vnt.protocolManager, err = NewProtocolManager(vnt.chainConfig, config.SyncMode, config.NetworkId, vnt.eventMux, vnt.txPool, vnt.engine, vnt.blockchain, chainDb, node); err != nil {
		return nil, err
	}
	vnt.miner = miner.New(vnt, vnt.chainConfig, vnt.EventMux(), vnt.engine)
	vnt.miner.SetExtra(makeExtraData(config.ExtraData))

	vnt.APIBackend = &VntAPIBackend{vnt, nil}
	gpoParams := config.GPO
	if gpoParams.Default == nil {
		gpoParams.Default = config.GasPrice
	}
	vnt.APIBackend.gpo = gasprice.NewOracle(vnt.APIBackend, gpoParams)

	return vnt, nil
}

func makeExtraData(extra []byte) []byte {
	if len(extra) == 0 {
		// create default extradata
		extra, _ = rlp.EncodeToBytes([]interface{}{
			uint(params.VersionMajor<<16 | params.VersionMinor<<8 | params.VersionPatch),
			"gvnt",
			runtime.Version(),
			runtime.GOOS,
		})
	}
	if uint64(len(extra)) > params.MaximumExtraDataSize {
		log.Warn("Miner extra data exceed limit", "extra", hexutil.Bytes(extra), "limit", params.MaximumExtraDataSize)
		extra = nil
	}
	return extra
}

// CreateDB creates the chain database.
func CreateDB(ctx *node.ServiceContext, config *Config, name string) (vntdb.Database, error) {
	log.Debug("backend", "func", "CreateDB", "ctx", ctx, "config", config, "name", name)
	db, err := ctx.OpenDatabase(name, config.DatabaseCache, config.DatabaseHandles)
	if err != nil {
		return nil, err
	}
	if db, ok := db.(*vntdb.LDBDatabase); ok {
		db.Meter("vnt/db/chaindata/")
	}
	return db, nil
}

// CreateConsensusEngine creates the required type of consensus engine instance for an VNT service
func CreateConsensusEngine(ctx *node.ServiceContext, chainConfig *params.ChainConfig, db vntdb.Database) consensus.Engine {
	// Otherwise assume DPoS
	cfg := chainConfig.Dpos
	// TODO vnt below test net config data should be in a single const variable
	// now, it's used for tests
	if chainConfig.Dpos == nil {
		cfg = &params.DposConfig{
			WitnessesNum: 4,
			Period:       2,
			WitnessesUrl: nil,
		}
	}
	return dpos.New(cfg, db)
}

// APIs return the collection of RPC services the ethereum package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *VNT) APIs() []rpc.API {
	apis := vntapi.GetAPIs(s.APIBackend)

	// Append any APIs exposed explicitly by the consensus engine
	apis = append(apis, s.engine.APIs(s.BlockChain())...)

	// Append all the local APIs and return
	return append(apis, []rpc.API{
		{
			Namespace: "core",
			Version:   "1.0",
			Service:   NewPublicVntAPI(s),
			Public:    true,
		}, {
			Namespace: "core",
			Version:   "1.0",
			Service:   NewPublicMinerAPI(s),
			Public:    true,
		}, {
			Namespace: "core",
			Version:   "1.0",
			Service:   downloader.NewPublicDownloaderAPI(s.protocolManager.downloader, s.eventMux),
			Public:    true,
		}, {
			Namespace: "bp",
			Version:   "1.0",
			Service:   NewPrivateMinerAPI(s),
			Public:    false,
		}, {
			Namespace: "core",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.APIBackend, false),
			Public:    true,
		}, {
			Namespace: "admin",
			Version:   "1.0",
			Service:   NewPrivateAdminAPI(s),
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPublicDebugAPI(s),
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPrivateDebugAPI(s.chainConfig, s),
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		},
	}...)
}

func (s *VNT) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *VNT) Coinbase() (eb common.Address, err error) {
	s.lock.RLock()
	coinbase := s.coinbase
	s.lock.RUnlock()

	if coinbase != (common.Address{}) {
		return coinbase, nil
	}
	if wallets := s.AccountManager().Wallets(); len(wallets) > 0 {
		if accounts := wallets[0].Accounts(); len(accounts) > 0 {
			coinbase := accounts[0].Address

			s.lock.Lock()
			s.coinbase = coinbase
			s.lock.Unlock()

			log.Info("Coinbase automatically configured", "address", coinbase)
			return coinbase, nil
		}
	}
	return common.Address{}, fmt.Errorf("coinbase must be explicitly specified")
}

// SetCoinbase sets the block producing reward address.
func (s *VNT) SetCoinbase(coinbase common.Address) {
	s.lock.Lock()
	s.coinbase = coinbase
	s.lock.Unlock()

	s.miner.SetCoinbase(coinbase)
}

func (s *VNT) StartProducing(local bool) error {
	eb, err := s.Coinbase()
	if err != nil {
		log.Error("Cannot start block producing without coinbase", "err", err)
		return fmt.Errorf("coinbase missing: %v", err)
	}

	if dpos, ok := s.engine.(*dpos.Dpos); ok {
		wallet, err := s.accountManager.Find(accounts.Account{Address: eb})
		if wallet == nil || err != nil {
			log.Error("Coinbase account unavailable locally", "err", err)
			return fmt.Errorf("signer missing: %v", err)
		}
		dpos.Authorize(eb, wallet.SignHash)
	}
	if local {
		// If local (CPU) block producing is started, we can disable the transaction rejection
		// mechanism introduced to speed sync times. CPU block producing on mainnet is ludicrous
		// so none will ever hit this path, whereas marking sync done on CPU block producing
		// will ensure that private networks work in single miner mode too.
		atomic.StoreUint32(&s.protocolManager.acceptTxs, 1)
	}
	go s.miner.Start(eb)
	return nil
}

func (s *VNT) StopProducing()      { s.miner.Stop() }
func (s *VNT) IsProducing() bool   { return s.miner.Producing() }
func (s *VNT) Miner() *miner.Miner { return s.miner }

func (s *VNT) AccountManager() *accounts.Manager  { return s.accountManager }
func (s *VNT) BlockChain() *core.BlockChain       { return s.blockchain }
func (s *VNT) TxPool() *core.TxPool               { return s.txPool }
func (s *VNT) EventMux() *event.TypeMux           { return s.eventMux }
func (s *VNT) Engine() consensus.Engine           { return s.engine }
func (s *VNT) ChainDb() vntdb.Database            { return s.chainDb }
func (s *VNT) IsListening() bool                  { return true } // Always listening
func (s *VNT) EthVersion() int                    { return int(s.protocolManager.SubProtocols[0].Version) }
func (s *VNT) NetVersion() uint64                 { return s.networkId }
func (s *VNT) Downloader() *downloader.Downloader { return s.protocolManager.downloader }

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (s *VNT) Protocols() []vntp2p.Protocol {
	if s.lesServer == nil {
		return s.protocolManager.SubProtocols
	}
	return append(s.protocolManager.SubProtocols, s.lesServer.Protocols()...)
}

// Start implements node.Service, starting all internal goroutines needed by the
// VNT protocol implementation.
func (s *VNT) Start(srvr *vntp2p.Server) error {
	// Start the bloom bits servicing goroutines
	s.startBloomHandlers()

	// Start the RPC service
	s.netRPCService = vntapi.NewPublicNetAPI(srvr, s.NetVersion())

	// Figure out a max peers count based on the server limits
	maxPeers := srvr.MaxPeers
	if s.config.LightServ > 0 {
		if s.config.LightPeers >= srvr.MaxPeers {
			return fmt.Errorf("invalid peer config: light peer count (%d) >= total peer count (%d)", s.config.LightPeers, srvr.MaxPeers)
		}
		maxPeers -= s.config.LightPeers
	}
	// Start the networking layer and the light server if requested
	s.protocolManager.Start(maxPeers)
	if s.lesServer != nil {
		s.lesServer.Start(srvr)
	}
	return nil
}

// Stop implements node.Service, terminating all internal goroutines used by the
// VNT protocol.
func (s *VNT) Stop() error {
	s.bloomIndexer.Close()
	s.blockchain.Stop()
	s.protocolManager.Stop()
	if s.lesServer != nil {
		s.lesServer.Stop()
	}
	s.txPool.Stop()
	s.miner.Stop()
	s.eventMux.Stop()

	s.chainDb.Close()
	close(s.shutdownChan)

	return nil
}
