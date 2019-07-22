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
	"crypto/rand"
	"sync"
	"time"

	util "github.com/ipfs/go-ipfs-util"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	peer "github.com/libp2p/go-libp2p-peer"
	routing "github.com/libp2p/go-libp2p-routing"
	"github.com/vntchain/go-vnt/log"
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
	GetDhtTable() *dht.IpfsDHT
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
	var bootStrapConfig = dht.DefaultBootstrapConfig
	bootStrapConfig.Period = time.Duration(refreshInterval)
	bootStrapConfig.Timeout = time.Duration(searchTimeOut)
	proc, err := vdht.table.BootstrapWithConfig(bootStrapConfig)
	if err != nil {
		log.Debug("Start refresh k-bucket error", "error", err)
		return err
	}

	// wait till ctx or dht.Context exits.
	// we have to do it this way to satisfy the Routing interface (contexts)
	go func() {
		defer proc.Close()
		select {
		case <-ctx.Done():
		case <-vdht.table.Context().Done():
		}
	}()

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

	runQuery := func(ctxs context.Context, id peer.ID) {
		p, err := vdht.table.FindPeer(ctxs, id)
		if err == routing.ErrNotFound {
		} else if err != nil {
			log.Debug("lookup peer occurs error", "error", err)
		} else {
			log.Debug("lookup peer find peer", "id", id.ToString(), "peer", p.ID.ToString())
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

func (vdht *VNTDht) GetDhtTable() *dht.IpfsDHT {
	return vdht.table
}
