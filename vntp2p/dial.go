package vntp2p

import (
	"context"

	"errors"
	"fmt"
	"github.com/vntchain/go-vnt/log"
	"github.com/libp2p/go-libp2p-peer"
	"time"
)

var (
	errSelf             = errors.New("is self")
	errAlreadyDialing   = errors.New("already dialing")
	errAlreadyConnected = errors.New("already connected")
	errRecentlyDialed   = errors.New("recently dialed")
	errNotWhitelisted   = errors.New("not contained in netrestrict whitelist")
)

type taskstate struct {
	maxDynDials int
	table       DhtTable
	bootnodes   []peer.ID
	static      map[peer.ID]*dialTask
	dailmap     map[peer.ID]dialFlag
}

type task interface {
	Do(ctx context.Context, server *Server)
}

type dialTask struct {
	flag   dialFlag
	target peer.ID
	pid    string
}

type lookupTask struct {
}

type waitExpireTask struct {
	time.Duration
}

func (t *taskstate) newTasks(peers map[peer.ID]*Peer) []task {
	var newtasks []task

	addDial := func(flag dialFlag, n peer.ID) bool {
		if err := t.checkDial(n, peers); err != nil {
			// fmt.Println("dail skip")
			return false
		}
		t.dailmap[n] = flag
		// fmt.Println("begin to Add: ", n)
		newtasks = append(newtasks, &dialTask{target: n, pid: PID})
		return true
	}

	needdail := t.maxDynDials
	// dail

	for _, flag := range t.dailmap {
		if flag&dynDialedDail != 0 {
			needdail--
		}
	}

	// newtasks = append(newtasks, &dailTask{})
	for id, task := range t.static {
		err := t.checkDial(id, peers)
		switch err {
		case errNotWhitelisted, errSelf:
			log.Warn("Removing static dial candidate", "id", id, "err", err)
			delete(t.static, id)
		case nil:
			t.dailmap[id] = task.flag
			newtasks = append(newtasks, task)
		}
	}

	for _, bootnode := range t.bootnodes {
		// fmt.Println("bootnode: ", bootnode)
		// for k, _ := range t.dailmap {
		// 	fmt.Println("dailmap: ", k)
		// }

		// for k, _ := range peers {
		// 	fmt.Println("peers: ", k)
		// }

		if addDial(staticDialedDail, bootnode) {
			needdail--
		}
	}

	randomDail := needdail / 2

	if randomDail > 0 {
		randompeerlist := t.table.RandomPeer()
		for i := 0; i < randomDail && i < len(randompeerlist); i++ {
			if addDial(dynDialedDail, randompeerlist[i]) {
				needdail--
			}
		}

	}
	// lookup
	// newtasks = append(newtasks, &lookupTask{})

	// waitExpireTask
	// newtasks = append(newtasks, &waitExpireTask{})

	if len(newtasks) == 0 {
		newtasks = append(newtasks, &waitExpireTask{1 * time.Second})
	}
	// fmt.Println("tasks: ", newtasks)
	return newtasks
}

func (s *taskstate) checkDial(n peer.ID, peers map[peer.ID]*Peer) error {
	_, dialing := s.dailmap[n]
	switch {
	case dialing:
		return errAlreadyDialing
	case peers[n] != nil:
		return errAlreadyConnected
	}
	return nil
}

func (s *taskstate) removeStatic(n *Node) {
	delete(s.static, n.Id)
}

func (s *taskstate) taskDone(t task) {
	switch t := t.(type) {
	case *dialTask:
		// s.hist.add(t.dest.ID, now.Add(dialHistoryExpiration))
		// log.Debug("taskDone", "dialTask", t.target)
		// fmt.Println("taskDone dialTask", t.target)
		delete(s.dailmap, t.target)
	case *lookupTask:
		log.Debug("taskDone", "lookupTask")
	}
}

func (s *taskstate) addStatic(n *Node) {
	s.static[n.Id] = &dialTask{flag: staticDialedDail, target: n.Id, pid: PID}
}

func newTaskState(maxdail int, bootnodes []peer.ID, dht DhtTable) *taskstate {
	s := &taskstate{
		maxDynDials: maxdail,
		bootnodes:   make([]peer.ID, len(bootnodes)),
		dailmap:     make(map[peer.ID]dialFlag),
		static:      make(map[peer.ID]*dialTask),
		table:       dht,
	}

	copy(s.bootnodes, bootnodes)

	log.Debug("Task state", "bootnodes", s.bootnodes)

	return s
}

func (t *dialTask) Do(ctx context.Context, server *Server) {
	// 检验目的地有效
	if !t.checkTarget() {
		// 如果无效通过lookup获得地址
		return
	}

	// 直接连接
	// fmt.Println("it's time to dial")

	t.dial(ctx, server, t.target, t.pid)
}

func (t *dialTask) checkTarget() bool {
	if t.target == "" {
		return false
	}
	return true
}

func (t *dialTask) dial(ctx context.Context, server *Server, target peer.ID, pid string) error {
	return server.SetupStream(ctx, target, pid)
}

func (t *lookupTask) Do(ctx context.Context, server *Server) {
	fmt.Println("begin lookup")
	time.Sleep(1 * time.Second)
	fmt.Println("end lookup")
}

func (t *waitExpireTask) Do(ctx context.Context, server *Server) {
	time.Sleep(t.Duration)
}
