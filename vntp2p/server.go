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
	// "time"
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
	inboundDail
	trustedDail
)

type Config struct {
	PrivateKey *ecdsa.PrivateKey `toml:"-"`

	MaxPeers int

	MaxPendingPeers int `toml:",omitempty"`

	DialRatio int `toml:",omitempty"`

	NoDiscovery bool

	Name string `toml:"-"`

	BootstrapNodes []*Node

	StaticNodes []*Node

	TrustedNodes []*Node

	NetRestrict []*net.IPNet `toml:",omitempty"`

	NodeDatabase string `toml:",omitempty"`

	Protocols []Protocol `toml:"-"`

	ListenAddr string

	NAT libp2p.Option `toml:",omitempty"`

	// Dialer NodeDialer `toml:"-"`

	// NoDial bool `toml:",omitempty"`

	EnableMsgEvents bool

	Logger log.Logger `toml:",omitempty"`
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
	addstatic    chan *Node
	removestatic chan *Node

	addpeer chan *Stream
	delpeer chan peerDrop

	peerOp     chan peerOpFunc
	peerOpDone chan struct{}

	protomap map[string][]Protocol
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

func (server *Server) Start() error {
	log.Info("yhx-test", "server.Protocols", server.Protocols)
	if server.running {
		return errors.New("server already running")
	}

	server.lock.Lock()
	defer server.lock.Unlock()

	log.Info("server is Starting!")
	// init

	server.addpeer = make(chan *Stream)
	server.delpeer = make(chan peerDrop)
	server.addstatic = make(chan *Node)
	server.removestatic = make(chan *Node)
	server.quit = make(chan struct{})
	server.peerOp = make(chan peerOpFunc)
	server.peerOpDone = make(chan struct{})

	// 协议映射初始化
	server.protomap = make(map[string][]Protocol)

	server.protomap[PID] = server.Protocols

	// Listen
	// run
	if server.ListenAddr == "" {
		return fmt.Errorf("P2P Server can't start for no listening")
	}

	listenPort := server.Config.ListenAddr[1:]
	log.Info("startVNTNode()", "listenPort", listenPort)
	ctx, cancel := context.WithCancel(context.Background())
	server.cancel = cancel

	d := server.NodeDatabase
	vdht, host, err := ConstructDHT(ctx, MakePort(listenPort), nil, d, server.Config.NetRestrict, server.Config.NAT)
	if err != nil {
		log.Error("startVNTNode()", "constructDHT error", err)
		return err
	}

	// setStreamHandler can only handle request message
	// it can not hear response
	host.SetStreamHandler(PID, server.HandleStream)

	log.Info("startVNTNode()", "own nodeID", host.ID())
	server.table = NewDHTTable(vdht, host.ID())
	server.host = host

	bootnodes := server.LoadConfig(ctx)

	maxdails := server.maxDialedConns()

	taskState := newTaskState(maxdails, bootnodes, server.table)

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
		server.table.Update(ctx, bootnode.Id)

		bootnodes = append(bootnodes, bootnode.Id)
	}

	return bootnodes

}

func (server *Server) run(ctx context.Context, tasker taskworker) {
	defer server.loopWG.Done()
	server.table.Start(ctx)
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

		case t := <-server.addpeer:
			remoteID := t.Conn.Conn().RemotePeer()
			if _, ok := peers[remoteID]; ok { // this peer already exists
				break
			}
			p := newPeer(t)

			if server.EnableMsgEvents {
				p.events = &server.peerFeed
			}
			go server.runPeer(p)
			peers[p.RemoteID()] = p
			log.Info("yhx-test", "peers", peers)

		case t := <-server.addstatic:
			tasker.addStatic(t)
		case t := <-server.removestatic:
			tasker.removeStatic(t)
			if p, ok := peers[t.Id]; ok {
				p.Disconnect(DiscRequested)
			}

		case op := <-server.peerOp:
			// This channel is used by Peers and PeerCount.
			op(peers)
			server.peerOpDone <- struct{}{}

		case pd := <-server.delpeer:
			// A peer disconnected.
			//log.Debug("Removing p2p peer", "peers", len(peers)-1, "req", "err", pd.err)
			// fmt.Println("Del peer", pd.RemoteID())
			delete(peers, pd.RemoteID())
		}
	}
}

func (server *Server) Stop() {
	log.Info("Server is Stopping!")
	defer server.cancel()
	return
}

func (server *Server) AddPeer(ctx context.Context, node *Node) {

	server.host.Peerstore().AddAddrs(node.Id, []ma.Multiaddr{node.Addr}, peerstore.PermanentAddrTTL)
	server.table.Update(ctx, node.Id)

	select {
	case server.addstatic <- node:
	case <-server.quit:
	}
}

func (server *Server) RemovePeer(node *Node) {
	select {
	case server.removestatic <- node:
	case <-server.quit:
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

// GetPeerByRemoteID get specific peer by remoteID
// if it doesn't exist, new it
// this function guarantee get the wanted peer
func (server *Server) GetPeerByRemoteID(s inet.Stream) *Peer {
	var p *Peer

	// always try to new this peer
	err := server.dispatch(&Stream{Conn: s, Protocols: server.protomap[PID]}, server.addpeer)
	if err != nil {
		log.Error("GetPeerByRemoteID()", "new peer error", err)
		return nil
	}

	select {
	case server.peerOp <- func(peers map[peer.ID]*Peer) {
		remoteID := s.Conn().RemotePeer()
		if val, ok := peers[remoteID]; ok {
			p = val
		}
	}:
		<-server.peerOpDone
	case <-server.quit:
	}

	pid := s.Conn().RemotePeer()
	log.Info("yhx-test", "GetPeerByRemoteID peerid", pid, "peer got", p)

	return p
}

func (server *Server) Peers() []*Peer {
	var ps []*Peer
	select {
	case server.peerOp <- func(peers map[peer.ID]*Peer) {
		for _, p := range peers {
			ps = append(ps, p)
		}
	}:
		<-server.peerOpDone
	case <-server.quit:
	}
	return ps
}

func (server *Server) PeerCount() int {
	var count int
	select {
	case server.peerOp <- func(ps map[peer.ID]*Peer) { count = len(ps) }:
		<-server.peerOpDone
	case <-server.quit:
	}
	return count
}

func (server *Server) maxDialedConns() int {
	r := server.DialRatio
	if r == 0 {
		r = defaultDialRatio
	}
	return server.MaxPeers / r
}

func (server *Server) SetupStream(ctx context.Context, target peer.ID, pid string) error {
	// log.Info("yhx-test", "SetupStream target", target, "pid", pid)
	s, err := server.host.NewStream(ctx, target, protocol.ID(pid))
	if err != nil {
		// fmt.Println("SetupStream NewStream Error: ", err)
		return err
	}

	// handle response message
	go server.HandleStream(s)
	/* rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	vntMessenger := &VNTMessenger{
		protocol: Protocol{},
		in:       make(chan Msg),
		w:        rw,
	}

	their, err := doProtocolHandshake(vntMessenger, server.ourHandshake)
	if err != nil {
		log.Error("SetupStream()", "failed protocolHandshake", err)
		return err
	} */

	err = server.dispatch(&Stream{Conn: s, Protocols: server.protomap[pid]}, server.addpeer)
	if err != nil {
		fmt.Println("SetupStream dispatch Error: ", err)
		return err
	}
	return nil
}

func (server *Server) runPeer(p *Peer) {
	// broadcast peer add
	server.peerFeed.Send(&PeerEvent{
		Type: PeerEventTypeAdd,
		Peer: p.RemoteID(),
	})

	// run the protocol
	remoteRequested, err := p.run()

	// broadcast peer drop
	server.peerFeed.Send(&PeerEvent{
		Type:  PeerEventTypeDrop,
		Peer:  p.RemoteID(),
		Error: err.Error(),
	})

	// Note: run waits for existing peers to be sent on srv.delpeer
	// before returning, so this send should not select on srv.quit.
	server.delpeer <- peerDrop{p, err, remoteRequested}
}

func (server *Server) dispatch(s *Stream, stage chan<- *Stream) error {
	select {
	case stage <- s:
	case <-server.quit:
		return errServerStopped
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
		VNTNode:    node.Addr.String() + "/ipfs/" + node.Id.ToString(),
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
