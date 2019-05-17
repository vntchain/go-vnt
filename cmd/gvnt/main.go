// Copyright 2014 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

// gvnt is the official command-line client for VNT.
package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
	godebug "runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/gosigar"
	"github.com/vntchain/go-vnt/accounts"
	"github.com/vntchain/go-vnt/accounts/keystore"
	"github.com/vntchain/go-vnt/cmd/utils"
	"github.com/vntchain/go-vnt/console"
	"github.com/vntchain/go-vnt/internal/debug"
	"github.com/vntchain/go-vnt/log"
	"github.com/vntchain/go-vnt/metrics"
	"github.com/vntchain/go-vnt/node"
	"github.com/vntchain/go-vnt/vnt"
	"github.com/vntchain/go-vnt/vntclient"
	cli "gopkg.in/urfave/cli.v1"
)

const (
	clientIdentifier = "gvnt" // Client identifier to advertise over the network
)

var (
	// Git SHA1 commit hash of the release (set via linker flags)
	gitCommit = ""
	// The app that holds all commands and flags.
	app = utils.NewApp(gitCommit, "the go-vnt command line interface")
	// flags that configure the node
	nodeFlags = []cli.Flag{
		utils.IdentityFlag,
		utils.UnlockedAccountFlag,
		utils.PasswordFileFlag,
		utils.FindNodeFlag,
		utils.VNTBootnodeFlag,
		utils.BootnodesFlag,
		utils.BootnodesV4Flag,
		utils.BootnodesV5Flag,
		utils.DataDirFlag,
		utils.KeyStoreDirFlag,
		// utils.EthashCacheDirFlag,
		// utils.EthashCachesInMemoryFlag,
		// utils.EthashCachesOnDiskFlag,
		// utils.EthashDatasetDirFlag,
		// utils.EthashDatasetsInMemoryFlag,
		// utils.EthashDatasetsOnDiskFlag,
		utils.TxPoolNoLocalsFlag,
		utils.TxPoolJournalFlag,
		utils.TxPoolRejournalFlag,
		utils.TxPoolPriceLimitFlag,
		utils.TxPoolPriceBumpFlag,
		utils.TxPoolAccountSlotsFlag,
		utils.TxPoolGlobalSlotsFlag,
		utils.TxPoolAccountQueueFlag,
		utils.TxPoolGlobalQueueFlag,
		utils.TxPoolLifetimeFlag,
		utils.SyncModeFlag,
		utils.GCModeFlag,
		utils.LightServFlag,
		utils.LightPeersFlag,
		utils.LightKDFFlag,
		utils.CacheFlag,
		utils.CacheDatabaseFlag,
		utils.CacheGCFlag,
		utils.TrieCacheGenFlag,
		utils.ListenPortFlag,
		utils.MaxPeersFlag,
		utils.MaxPendingPeersFlag,
		utils.CoinbaseFlag,
		utils.GasPriceFlag,
		utils.ProducingEnabledFlag,
		utils.TargetGasLimitFlag,
		utils.NATFlag,
		utils.NoDiscoverFlag,
		utils.DiscoveryV5Flag,
		utils.NetrestrictFlag,
		utils.NodeKeyFileFlag,
		utils.NodeKeyHexFlag,
		utils.VMEnableDebugFlag,
		utils.NetworkIdFlag,
		utils.RPCCORSDomainFlag,
		utils.RPCVirtualHostsFlag,
		utils.EthStatsURLFlag,
		utils.MetricsEnabledFlag,
		utils.NoCompactionFlag,
		utils.GpoBlocksFlag,
		utils.GpoPercentileFlag,
		utils.ExtraDataFlag,
		configFileFlag,
	}

	rpcFlags = []cli.Flag{
		utils.RPCEnabledFlag,
		utils.RPCListenAddrFlag,
		utils.RPCPortFlag,
		utils.RPCApiFlag,
		utils.WSEnabledFlag,
		utils.WSListenAddrFlag,
		utils.WSPortFlag,
		utils.WSApiFlag,
		utils.WSAllowedOriginsFlag,
		utils.IPCDisabledFlag,
		utils.IPCPathFlag,
	}

	whisperFlags = []cli.Flag{
		utils.WhisperEnabledFlag,
		utils.WhisperMaxMessageSizeFlag,
		utils.WhisperMinPOWFlag,
	}
)

func init() {
	// Initialize the CLI app and start Gvnt
	app.Action = gvnt
	app.HideVersion = true // we have a command to print the version
	app.Copyright = "Copyright 2018-2019 The go-vnt Authors"
	app.Commands = []cli.Command{
		// See chaincmd.go:
		initCommand,
		importCommand,
		exportCommand,
		importPreimagesCommand,
		exportPreimagesCommand,
		copydbCommand,
		removedbCommand,
		dumpCommand,
		// See monitorcmd.go:
		monitorCommand,
		// See accountcmd.go:
		accountCommand,
		// See consolecmd.go:
		consoleCommand,
		attachCommand,
		javascriptCommand,
		// See misccmd.go:
		versionCommand,
		bugCommand,
		licenseCommand,
		// See config.go
		dumpConfigCommand,
	}
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Flags = append(app.Flags, nodeFlags...)
	app.Flags = append(app.Flags, rpcFlags...)
	app.Flags = append(app.Flags, consoleFlags...)
	app.Flags = append(app.Flags, debug.Flags...)
	app.Flags = append(app.Flags, whisperFlags...)

	app.Before = func(ctx *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		if err := debug.Setup(ctx); err != nil {
			return err
		}
		// Cap the cache allowance and tune the garbage colelctor
		var mem gosigar.Mem
		if err := mem.Get(); err == nil {
			allowance := int(mem.Total / 1024 / 1024 / 3)
			if cache := ctx.GlobalInt(utils.CacheFlag.Name); cache > allowance {
				log.Warn("Sanitizing cache to Go's GC limits", "provided", cache, "updated", allowance)
				ctx.GlobalSet(utils.CacheFlag.Name, strconv.Itoa(allowance))
			}
		}
		// Ensure Go's GC ignores the database cache for trigger percentage
		cache := ctx.GlobalInt(utils.CacheFlag.Name)
		gogc := math.Max(20, math.Min(100, 100/(float64(cache)/1024)))

		log.Debug("Sanitizing Go's GC trigger", "percent", int(gogc))
		godebug.SetGCPercent(int(gogc))

		// Start system runtime metrics collection
		go metrics.CollectProcessMetrics(3 * time.Second)

		utils.SetupNetwork(ctx)
		return nil
	}

	app.After = func(ctx *cli.Context) error {
		debug.Exit()
		console.Stdin.Close() // Resets terminal mode.
		return nil
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// gvnt is the main entry point into the system if no special subcommand is ran.
// It creates a default node based on the command line arguments and runs it in
// blocking mode, waiting for it to be shut down.
func gvnt(ctx *cli.Context) error {
	node := makeFullNode(ctx)
	// go startVNTNode(node)
	startNode(ctx, node)
	node.Wait()
	return nil
}

/* func startVNTNode(stack *node.Node) {
	findnode := stack.Config().P2P.FindNode
	bootnode := stack.Config().P2P.VNTBootnode
	listenPort := stack.Config().P2P.ListenAddr[1:]
	log.Info("startVNTNode()", "listenPort", listenPort)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	vdht, host, err := vntp2p.ConstructDHT(ctx, vntp2p.MakePort(listenPort), nil)
	if err != nil {
		log.Error("startVNTNode()", "constructDHT error", err)
		return
	}
	log.Info("startVNTNode()", "own nodeID", host.ID())

	h := &vntprotocol.HostWrapper{
		Host: host,
	}

	host.SetStreamHandler(vntp2p.PID, h.HandleStream)
	bootnodeAddr, bootnodeID, err := vntp2p.GetAddr(bootnode)
	if err != nil {
		log.Error("startVNTNode()", "getAddr error", err)
		return
	}
	host.Peerstore().AddAddrs(bootnodeID, []ma.Multiaddr{bootnodeAddr}, peerstore.PermanentAddrTTL)
	vdht.Update(ctx, bootnodeID)

	// test vntdb
	//err = vdht.PutValue(ctx, "/v/hello2", []byte("helloworld2"))
	//if err != nil {
	//	log.Error("startVNTNode()", "putValue error", err)
	//}

	contentKey := "/v/hello2"
	log.Info("begin getValue now!")
	best, err := vdht.GetValue(ctx, contentKey)
	if err != nil {
		log.Error("startVNTNode()", "getValue error", err)
	} else {
		log.Info("startVNTNode(), value gotten", "key", contentKey, "value", string(best))
	}

	content := host.Peerstore().Peers()
	for i := range content {
		log.Info("startVNTNode()", "---> index", i, "PeerInfo", host.Peerstore().PeerInfo(content[i]))
	}

	//sayHelloToBootnode(ctx, host, bootnodeID)

	if findnode != "" {
		findNode(ctx, findnode, vdht)
	}

	select {}
} */

/* func findNode(ctx context.Context, findnode string, vdht *dht.IpfsDHT) {
	_, findnodeID, err := vntp2p.GetAddr(findnode)
	if err != nil {
		log.Error("findNode()", "getAddr error", err)
		return
	}
	log.Info("findNode()", "findnodeID", findnodeID)
	targetPeerInfo, err := vdht.FindPeer(ctx, findnodeID)
	if err != nil {
		log.Error("findNode()", "findPeer error", err)
		return
	}
	log.Info("findNode()", "find peerid SUCCESS with info", targetPeerInfo)
} */

/* func sayHelloToBootnode(ctx context.Context, host p2phost.Host, nodeID peer.ID) {
	s, err := host.NewStream(ctx, nodeID, vntp2p.PID)
	if err != nil {
		log.Error("sayHelloToBootnode()", "NewStream Error", err)
		return
	}

	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	// 封包，并且发送
	msg, err := vntp2p.WritePackage(rw, 1, []byte("yoyoyo\n"))
	if err != nil {
		log.Error("sayHelloToBootnode()", "WritePackage Error", err)
		return
	}

	log.Info("sayHelloToBootnode()", "msg", string(msg))
	_, err = rw.Write(msg)
	if err != nil {
		log.Error("sayHelloToBootnode()", "Write Error", err)
		return
	}
	rw.WriteByte('\n')
	rw.Flush()
} */

// startNode boots up the system node and all registered protocols, after which
// it unlocks any requested accounts, and starts the RPC/IPC interfaces and the
// miner.
func startNode(ctx *cli.Context, stack *node.Node) {
	debug.Memsize.Add("node", stack)

	// Start up the node itself
	utils.StartNode(stack)

	// Unlock any account specifically requested
	ks := stack.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
	passwords := utils.MakePasswordList(ctx)
	unlocks := strings.Split(ctx.GlobalString(utils.UnlockedAccountFlag.Name), ",")
	for i, account := range unlocks {
		if trimmed := strings.TrimSpace(account); trimmed != "" {
			unlockAccount(ctx, ks, trimmed, i, passwords)
		}
	}
	// Register wallet event handlers to open and auto-derive wallets
	events := make(chan accounts.WalletEvent, 16)
	stack.AccountManager().Subscribe(events)

	go func() {
		// Create a chain state reader for self-derivation
		rpcClient, err := stack.Attach()
		if err != nil {
			utils.Fatalf("Failed to attach to self: %v", err)
		}
		stateReader := vntclient.NewClient(rpcClient)

		// Open any wallets already attached
		for _, wallet := range stack.AccountManager().Wallets() {
			if err := wallet.Open(""); err != nil {
				log.Warn("Failed to open wallet", "url", wallet.URL(), "err", err)
			}
		}
		// Listen for wallet event till termination
		for event := range events {
			switch event.Kind {
			case accounts.WalletArrived:
				if err := event.Wallet.Open(""); err != nil {
					log.Warn("New wallet appeared, failed to open", "url", event.Wallet.URL(), "err", err)
				}
			case accounts.WalletOpened:
				status, _ := event.Wallet.Status()
				log.Info("New wallet appeared", "url", event.Wallet.URL(), "status", status)

				if event.Wallet.URL().Scheme == "ledger" {
					event.Wallet.SelfDerive(accounts.DefaultLedgerBaseDerivationPath, stateReader)
				} else {
					event.Wallet.SelfDerive(accounts.DefaultBaseDerivationPath, stateReader)
				}

			case accounts.WalletDropped:
				log.Info("Old wallet dropped", "url", event.Wallet.URL())
				event.Wallet.Close()
			}
		}
	}()
	// Start auxiliary services if enabled
	if ctx.GlobalBool(utils.ProducingEnabledFlag.Name) {
		// Producing only makes sense if a full VNT node is running
		if ctx.GlobalString(utils.SyncModeFlag.Name) == "light" {
			utils.Fatalf("Light clients do not support block producing")
		}
		var vnt *vnt.VNT
		if err := stack.Service(&vnt); err != nil {
			utils.Fatalf("VNT service not running: %v", err)
		}

		// Set the gas price to the limits from the CLI and start producing
		vnt.TxPool().SetGasPrice(utils.GlobalBig(ctx, utils.GasPriceFlag.Name))
		if err := vnt.StartProducing(true); err != nil {
			utils.Fatalf("Failed to start block producing: %v", err)
		}
	}
}
