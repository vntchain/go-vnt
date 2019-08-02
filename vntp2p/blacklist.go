package vntp2p

import (
	"github.com/hashicorp/golang-lru"
	"github.com/libp2p/go-libp2p-peer"
	"github.com/vntchain/go-vnt/log"
)

type BlackList struct {
	cache	*lru.Cache
	rwPid	chan peer.ID
}

func NewPeerBlackList() *BlackList {
	blacklist, err := lru.New(1024)
	if err != nil {
		panic(err)
	}

	return &BlackList{
		blacklist,
		make(chan peer.ID),
	}
}

var blacklist = NewPeerBlackList()

func (b *BlackList) write(pid peer.ID) {
	b.rwPid <- pid
}

func (b *BlackList) run() {
	for {
		select {
		case pid := <- b.rwPid:
			log.Info("Add to blacklist:", "pid", pid)
			if !b.cache.Contains(pid) {
				b.cache.Add(pid, true)
			}
		}
	}
}

func (b *BlackList) exists(pid peer.ID) bool {
	return b.cache.Contains(pid)
}