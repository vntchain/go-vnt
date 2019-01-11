package vntp2p

import (
	"context"
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	util "github.com/ipfs/go-ipfs-util"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	peer "github.com/libp2p/go-libp2p-peer"
	routing "github.com/libp2p/go-libp2p-routing"
)

const (
	refreshInterval = 30 * time.Second
	searchTimeOut   = 1 * time.Minute
)

type DhtTable interface {
	Start(ctx context.Context) error
	Lookup(ctx context.Context, targetID NodeID) []*NodeID
	Update(ctx context.Context, id peer.ID) error
	RandomPeer() []peer.ID
}

type VNTDht struct {
	mutex sync.Mutex
	table *dht.IpfsDHT
	self  peer.ID
}

func NewDHTTable(dht *dht.IpfsDHT, id peer.ID) *VNTDht {
	return &VNTDht{
		table: dht,
		self:  id,
	}
}

func (vdht *VNTDht) Start(ctx context.Context) error {
	// init

	// loop
	go vdht.loop(ctx)
	return nil
}

func (vdht *VNTDht) Update(ctx context.Context, id peer.ID) error {
	vdht.table.Update(ctx, id)
	return nil
}

func randomID() peer.ID {
	id := make([]byte, 16)
	rand.Read(id)
	id = util.Hash(id)

	// var aid NodeID
	// copy(aid, id)
	return peer.ID(id)
}

func (vdht *VNTDht) loop(ctx context.Context) {
	var (
		refresh     = time.NewTicker(refreshInterval)
		refreshDone = make(chan struct{})
	)
	go vdht.doRefresh(ctx, refreshDone)
	// loop:
	for {
		// 开始搜寻

		select {
		case <-refresh.C:
			go vdht.doRefresh(ctx, refreshDone)
		}
		// 刷新K桶
	}
}

func (vdht *VNTDht) Lookup(ctx context.Context, targetID NodeID) []*NodeID {
	// vdht.table.GetClosestPeers(vdht.Context, )

	results := vdht.lookup(ctx, targetID.PeerID())

	nodeids := []*NodeID{}

	for _, result := range results {
		nodeid := PeerIDtoNodeID(result)
		nodeids = append(nodeids, &nodeid)
	}

	return nodeids
}

func (vdht *VNTDht) lookup(ctx context.Context, targetid peer.ID) []peer.ID {
	// 开始搜寻

	// fmt.Println("Begin Lookup")
	cctx, cancel := context.WithTimeout(ctx, searchTimeOut)
	defer cancel()

	var merr util.MultiErr

	runQuery := func(ctxs context.Context, id peer.ID) {
		p, err := vdht.table.FindPeer(ctxs, id)
		if err == routing.ErrNotFound {
		} else if err != nil {
			merr = append(merr, err)
		} else {
			err := fmt.Errorf("Bootstrap peer error: Actually FOUND peer. (%s, %s)", id, p)
			fmt.Println("Warning ", err)
			merr = append(merr, err)
		}
	}

	runQuery(cctx, targetid)

	return nil
}

func (vdht *VNTDht) doRefresh(ctx context.Context, done chan struct{}) {
	// defer close(done)

	// Load nodes from the database and insert
	// them. This should yield a few previously seen nodes that are
	// (hopefully) still alive.
	// tab.loadSeedNodes(true)

	// Run self lookup to discover new neighbor nodes.
	// fmt.Println("@@@CCCC find my self", vdht.self)
	vdht.lookup(ctx, vdht.self)

	// The Kademlia paper specifies that the bucket refresh should
	// perform a lookup in the least recently used bucket. We cannot
	// adhere to this because the findnode target is a 512bit value
	// (not hash-sized) and it is not easily possible to generate a
	// sha3 preimage that falls into a chosen bucket.
	// We perform a few lookups with a random target instead.
	for i := 0; i < 3; i++ {
		target := randomID()
		// fmt.Println("random id: ", target)
		vdht.lookup(ctx, target)
	}
}

func (vdht *VNTDht) RandomPeer() []peer.ID {
	return vdht.table.GetRandomPeers()
}
