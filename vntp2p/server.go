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

package vntp2p

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"net"
	"sync"

	libp2p "github.com/libp2p/go-libp2p"
	p2phost "github.com/libp2p/go-libp2p-host"
	inet "github.com/libp2p/go-libp2p-net"
	protocol "github.com/libp2p/go-libp2p-protocol"
	"github.com/vntchain/go-vnt/event"
	"github.com/vntchain/go-vnt/log"

	// inet "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	// kb "github.com/libp2p/go-libp2p-kbucket"
)

const (
	maxActiveDialTasks = 16
	defaultDialRatio   = 3
)

var errServerStopped = errors.New("server stopped")

type dialFlag int

const (
	dynDialedDail dialFlag = 1 << iota
	staticDialedDail
)

type Config struct {
	PrivateKey      *ecdsa.PrivateKey `toml:"-"`
	MaxPeers        int
	MaxPendingPeers int `toml:",omitempty"`
	DialRatio       int `toml:",omitempty"`
	NoDiscovery     bool
	Name            string `toml:"-"`

	BootstrapNodes []*Node
	StaticNodes    []*Node
	TrustedNodes   []*Node

	NetRestrict  []*net.IPNet `toml:",omitempty"`
	NodeDatabase string       `toml:",omitempty"`
	Protocols    []Protocol   `toml:"-"`
	ListenAddr   string
	NAT          libp2p.Option `toml:"-"`
	// Dialer NodeDialer `toml:"-"`
	// NoDial bool `toml:",omitempty"`

	EnableMsgEvents bool
	Logger          log.Logger `toml:",omitempty"`
}

type Server struct {
	Config
	table   DhtTable
	host    p2phost.Host
	running bool

	peerFeed event.Feed
	loopWG   sync.WaitGroup
	cancel   context.CancelFunc

	lock sync.Mutex

	quit         chan struct{}
	addStatic    chan *Node
	removeStatic chan *Node

	peerOp chan peerOpFunc // channel to operation peers in run

	protocols map[string][]Protocol // vntp2p.PID to vnt sub protocols
}

type peerOpFunc func(map[peer.ID]*Peer)

// func NewServer() (*Server, error) {
// 	return newServer()
// }

// func newServer() (*Server, error) {
// 	return &Server{}, nil
// }

type peerDrop struct {
	*Peer
	err       error
	requested bool // true if signaled by the peer
}

// close server.quit to broadcast the server shutdown
func (server *Server) Close() {
	log.Info("P2P server is being closed...")
	close(server.quit)
}

func (server *Server) Start() error {
	log.Trace("P2P server starting...", "protocols", server.Protocols)
	if server.running {
		return errors.New("server already running")
	}

	server.lock.Lock()
	defer server.lock.Unlock()

	server.addStatic = make(chan *Node)
	server.removeStatic = make(chan *Node)
	server.quit = make(chan struct{})
	server.peerOp = make(chan peerOpFunc)

	// 协议映射初始化
	server.protocols = make(map[string][]Protocol)

	server.protocols[PID] = server.Protocols

	// Listen
	// run
	if server.ListenAddr == "" {
		return fmt.Errorf("P2P Server can't start for no listening")
	}

	listenPort := server.Config.ListenAddr[1:]
	ctx, cancel := context.WithCancel(context.Background())
	server.cancel = cancel

	d := server.NodeDatabase
	vdht, host, err := ConstructDHT(ctx, MakePort(listenPort), nil, d, server.Config.NetRestrict, server.Config.NAT)
	if err != nil {
		log.Error("ConstructDHT failed", "error", err)
		return err
	}

	// setStreamHandler can only handle request message
	// it can not hear response
	host.SetStreamHandler(PID, server.HandleStream)

	server.table = NewDHTTable(vdht, host.ID())
	server.host = host

	bootnodes := server.LoadConfig(ctx)

	maxdials := server.maxDialedConns()

	taskState := newTaskState(maxdials, bootnodes, server.table)

	server.loopWG.Add(1)
	go server.run(ctx, taskState)

	server.running = true
	return nil
}

func (server *Server) LoadConfig(ctx context.Context) []peer.ID {
	// 创建初始连接

	bootnodes := []peer.ID{}

	for _, bootnode := range server.Config.BootstrapNodes {
		server.host.Peerstore().AddAddrs(bootnode.Id, []ma.Multiaddr{bootnode.Addr}, peerstore.PermanentAddrTTL)
		_ = server.table.Update(ctx, bootnode.Id)

		bootnodes = append(bootnodes, bootnode.Id)
	}

	return bootnodes

}

func (server *Server) run(ctx context.Context, tasker taskworker) {
	defer server.loopWG.Done()
	_ = server.table.Start(ctx)
	var (
		runningTasks []task
		queuedTasks  []task
		taskdone     = make(chan task, maxActiveDialTasks)
		peers        = make(map[peer.ID]*Peer)
	)

	delTask := func(t task) {
		for i := range runningTasks {
			if runningTasks[i] == t {
				runningTasks = append(runningTasks[:i], runningTasks[i+1:]...)
				break
			}
		}
	}
	startTasks := func(ts []task) (rest []task) {
		i := 0
		for ; len(runningTasks) < maxActiveDialTasks && i < len(ts); i++ {
			t := ts[i]
			go func() { t.Do(ctx, server); taskdone <- t }()
			runningTasks = append(runningTasks, t)
		}
		return ts[i:]
	}
	scheduleTasks := func() {
		queuedTasks = append(queuedTasks[:0], startTasks(queuedTasks)...)
		if len(runningTasks) < maxActiveDialTasks {
			// fmt.Println("begin new task")
			nt := tasker.newTasks(peers)
			queuedTasks = append(queuedTasks, startTasks(nt)...)
		}
	}

	for {
		scheduleTasks()
		select {
		case t := <-taskdone:
			tasker.taskDone(t)
			delTask(t)

		case t := <-server.addStatic:
			log.Debug("Adding static", "peer id", t.Id)
			tasker.addStatic(t)
			log.Debug("Added static", "peer id", t.Id)
		case t := <-server.removeStatic:
			tasker.removeStatic(t)
			if p, ok := peers[t.Id]; ok {
				p.Disconnect(DiscRequested)
			}

		case op := <-server.peerOp:
			op(peers)
		}
	}
}

func (server *Server) Stop() {
	log.Info("Server is Stopping!")
	defer server.cancel()
}

func (server *Server) AddStaticPeer(ctx context.Context, node *Node) {
	server.host.Peerstore().AddAddrs(node.Id, []ma.Multiaddr{node.Addr}, peerstore.PermanentAddrTTL)
	_ = server.table.Update(ctx, node.Id)

	select {
	case <-server.quit:
	case server.addStatic <- node:
	}
}

func (server *Server) RemovePeer(node *Node) {
	select {
	case <-server.quit:
	case server.removeStatic <- node:
	}
}

func (server *Server) SubscribeEvents(ch chan *PeerEvent) event.Subscription {
	return server.peerFeed.Subscribe(ch)
}

func (server *Server) PeersInfo() []*PeerInfo {
	infos := make([]*PeerInfo, 0, server.PeerCount())

	for _, peer := range server.Peers() {
		if peer != nil {
			infos = append(infos, peer.Info())
		}
	}
	for i := 0; i < len(infos); i++ {
		for j := i + 1; j < len(infos); j++ {
			if infos[i].ID > infos[j].ID {
				infos[i], infos[j] = infos[j], infos[i]
			}
		}
	}
	return infos
}

func (server *Server) Self() *Node {
	server.lock.Lock()
	defer server.lock.Unlock()

	if !server.running {
		return &Node{}
	}

	// hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", server.host.ID().Pretty()))

	addr := server.host.Addrs()[0]
	// fullAddr := addr.Encapsulate(hostAddr)

	return &Node{Addr: addr, Id: server.host.ID()}
}

func (server *Server) getPeer(pid peer.ID) *Peer {
	retCh := make(chan *Peer)

	query := func(peers map[peer.ID]*Peer) {
		if p, ok := peers[pid]; ok {
			retCh <- p
		}
	}

	select {
	case <-server.quit:
		return nil
	case server.peerOp <- query:
		return <-retCh
	}
}

// addPeer add a peer and do call protocol stack to handshake.
// If it already added, just return the peer.
func (server *Server) addPeer(s inet.Stream) *Peer {
	pid := s.Conn().RemotePeer()
	retCh := make(chan *Peer)
	stream := &Stream{stream: s, Protocols: server.protocols[PID]}

	add := func(peers map[peer.ID]*Peer) {
		var (
			p  *Peer
			ok bool
		)

		defer func() {
			retCh <- p
		}()

		log.Debug("Adding peer", "peer id", pid)
		if p, ok = peers[pid]; ok {
			p.log.Debug("Already exist peer")
			return
		}

		p = newPeer(stream, server)
		if server.EnableMsgEvents {
			p.events = &server.peerFeed
		}
		go server.runPeer(p)
		peers[pid] = p
		p.log.Debug("Added peer", "peers", peers)
	}

	select {
	case <-server.quit:
		return nil
	case server.peerOp <- add:
		return <-retCh
	}
}

func (server *Server) rmPeer(pid peer.ID, err error) {
	rm := func(peers map[peer.ID]*Peer) {
		if p, ok := peers[pid]; ok {
			p.log.Info("Remove p2p peer", "error", err)
			delete(peers, pid)
		}
	}

	select {
	case <-server.quit:
		return
	case server.peerOp <- rm:
		return
	}
}

func (server *Server) Peers() []*Peer {
	retCh := make(chan []*Peer)

	get := func(peers map[peer.ID]*Peer) {
		var ps []*Peer
		for _, p := range peers {
			ps = append(ps, p)
		}
		retCh <- ps
	}

	select {
	case <-server.quit:
		return nil
	case server.peerOp <- get:
		return <-retCh
	}
}

func (server *Server) PeerCount() int {
	retCh := make(chan int)
	count := func(ps map[peer.ID]*Peer) {
		retCh <- len(ps)
	}
	select {
	case <-server.quit:
		return 0
	case server.peerOp <- count:
		return <-retCh
	}
}

func (server *Server) maxDialedConns() int {
	r := server.DialRatio
	if r == 0 {
		r = defaultDialRatio
	}
	return server.MaxPeers / r
}

// SetupStream 主动发起连接
func (server *Server) SetupStream(ctx context.Context, target peer.ID, pid string) error {
	s, err := server.host.NewStream(ctx, target, protocol.ID(pid))
	if err != nil {
		// fmt.Println("SetupStream NewStream Error: ", err)
		return err
	}

	// handle response message
	go server.HandleStream(s)
	/* rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	vntMessenger := &VNTMsger{
		protocol: Protocol{},
		in:       make(chan Msg),
		w:        rw,
	}

	their, err := doProtocolHandshake(vntMessenger, server.ourHandshake)
	if err != nil {
		log.Error("SetupStream()", "failed protocolHandshake", err)
		return err
	} */

	// err = server.dispatch(&Stream{stream: s, Protocols: server.protomap[pid]}, server.addpeer)
	// if err != nil {
	// 	fmt.Println("SetupStream dispatch Error: ", err)
	// 	return err
	// }
	server.addPeer(s)
	return nil
}

func (server *Server) runPeer(p *Peer) {
	// broadcast peer add
	server.peerFeed.Send(&PeerEvent{
		Type: PeerEventTypeAdd,
		Peer: p.RemoteID(),
	})

	// run the protocol
	_, err := p.run()

	// broadcast peer drop
	server.peerFeed.Send(&PeerEvent{
		Type:  PeerEventTypeDrop,
		Peer:  p.RemoteID(),
		Error: err.Error(),
	})

	server.rmPeer(p.RemoteID(), err)
}

func (server *Server) dispatch(s *Stream, stage chan<- *Stream) error {
	select {
	case <-server.quit:
		return errServerStopped
	case stage <- s:
	}
	return nil
}

type NodeInfo struct {
	ID      string `json:"id"`    // Unique node identifier (also the encryption key)
	Name    string `json:"name"`  // Name of the node, including client type, version, OS, custom data
	VNTNode string `json:"vnode"` // Vnode URL for adding this peer from remote peers
	IP      string `json:"ip"`    // IP address of the node
	Ports   struct {
		Discovery int `json:"discovery"` // UDP listening port for discovery protocol
		Listener  int `json:"listener"`  // TCP listening port for RLPx
	} `json:"ports"`
	ListenAddr string                 `json:"listenAddr"`
	Protocols  map[string]interface{} `json:"protocols"`
}

func (server *Server) NodeInfo() *NodeInfo {
	node := server.Self()

	info := &NodeInfo{
		ID:         node.Id.ToString(),
		VNTNode:    node.String(),
		Name:       server.Name,
		IP:         GetIPfromAddr(node.Addr),
		ListenAddr: server.ListenAddr,
		Protocols:  make(map[string]interface{}),
	}

	// for _, proto := range server.Protocols {
	// 	if _, ok := info.Protocols[proto.Name]; !ok {
	// 		nodeInfo := interface{}("unknown")
	// 		if query := proto.NodeInfo; query != nil {
	// 			nodeInfo = proto.NodeInfo()
	// 		}
	// 		info.Protocols[proto.Name] = nodeInfo
	// 	}
	// }
	return info
}

type taskworker interface {
	newTasks(map[peer.ID]*Peer) []task
	addStatic(n *Node)
	removeStatic(n *Node)
	taskDone(t task)
}
