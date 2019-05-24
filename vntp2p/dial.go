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

	"errors"
	"time"

	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/vntchain/go-vnt/log"
)

var (
	errAlreadyDialing   = errors.New("already dialing")
	errAlreadyConnected = errors.New("already connected")
)

type taskstate struct {
	maxDynDials   int
	table         DhtTable
	bootnodes     []peer.ID
	lookupNode    []peer.ID
	lookupRunning bool
	static        map[peer.ID]*dialTask
	dialmap       map[peer.ID]dialFlag
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
	targets []peer.ID
}

type waitExpireTask struct {
	time.Duration
}

func (s *taskstate) newTasks(peers map[peer.ID]*Peer) []task {
	var newtasks []task

	addDial := func(flag dialFlag, n peer.ID) bool {
		if err := s.checkDial(n, peers); err != nil {
			log.Trace("Skip dial peer", "id", n, "err", err)
			return false
		}
		s.dialmap[n] = flag
		newtasks = append(newtasks, &dialTask{target: n, pid: PID})
		return true
	}

	needDial := s.maxDynDials
	// dial

	for _, flag := range s.dialmap {
		if flag&dynDialedDail != 0 {
			needDial--
		}
	}

	// newtasks = append(newtasks, &dialTask{})
	for id, task := range s.static {
		if err := s.checkDial(id, peers); err == nil {
			s.dialmap[id] = task.flag
			newtasks = append(newtasks, task)
		}
	}

	for _, bootnode := range s.bootnodes {
		// fmt.Println("bootnode: ", bootnode)
		// for k, _ := range s.dialmap {
		// 	fms.Println("dialmap: ", k)
		// }

		// for k, _ := range peers {
		// 	fmt.Println("peers: ", k)
		// }

		if addDial(staticDialedDail, bootnode) {
			needDial--
		}
	}

	randomDial := needDial / 2

	if randomDial > 0 {
		randompeerlist := s.table.RandomPeer()
		for i := 0; i < randomDial && i < len(randompeerlist); i++ {
			if addDial(dynDialedDail, randompeerlist[i]) {
				needDial--
			}
		}
	}

	// lookup
	// if still need to dial more peer, create dynamic dials from random
	// lookup results
	i := 0
	for ; i < len(s.lookupNode) && needDial > 0; i++ {
		if addDial(dynDialedDail, s.lookupNode[i]) {
			needDial--
		}
	}

	// update lookupNode, if need more, launch lookup task
	s.lookupNode = s.lookupNode[i:]
	if len(s.lookupNode) < needDial && !s.lookupRunning {
		s.lookupRunning = true
		newtasks = append(newtasks, &lookupTask{})
	}

	// waitExpireTask
	// newtasks = append(newtasks, &waitExpireTask{})

	if len(newtasks) == 0 {
		newtasks = append(newtasks, &waitExpireTask{1 * time.Second})
	}
	// fmt.Println("tasks: ", newtasks)
	return newtasks
}

func (s *taskstate) checkDial(n peer.ID, peers map[peer.ID]*Peer) error {
	_, dialing := s.dialmap[n]
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
		delete(s.dialmap, t.target)
	case *lookupTask:
		s.lookupRunning = false
		s.lookupNode = append(s.lookupNode, t.targets...)
	}
}

func (s *taskstate) addStatic(n *Node) {
	s.static[n.Id] = &dialTask{flag: staticDialedDail, target: n.Id, pid: PID}
}

func newTaskState(maxdial int, bootnodes []peer.ID, dht DhtTable) *taskstate {
	s := &taskstate{
		maxDynDials: maxdial,
		bootnodes:   make([]peer.ID, len(bootnodes)),
		dialmap:     make(map[peer.ID]dialFlag),
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

	log.Trace("Dial task", "target", t.target)
	_ = t.dial(ctx, server, t.target, t.pid)
}

func (t *dialTask) checkTarget() bool {
	return t.target != ""
}

func (t *dialTask) dial(ctx context.Context, server *Server, target peer.ID, pid string) (err error) {
	if err = server.SetupStream(ctx, target, pid); err != nil {
		log.Trace("Dial failed", "error", err)
	}
	return
}

func (t *lookupTask) Do(ctx context.Context, server *Server) {
	target := randomID()
	peers, err := server.table.GetDhtTable().GetClosestPeers(ctx, string(target))
	if err != nil {
		log.Trace("lookupTask failed", "error", err)
	}

	for p := range peers {
		t.targets = append(t.targets, p)
	}
}

func (t *waitExpireTask) Do(ctx context.Context, server *Server) {
	time.Sleep(t.Duration)
}
