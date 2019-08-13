package vntp2p

import (
	"github.com/libp2p/go-libp2p-peer"
	"github.com/vntchain/go-vnt/log"
	"github.com/bluele/gcache"
	"time"
)

type BlackList struct {
	cache	gcache.Cache
	rwPid	chan peer.ID
}

func NewPeerBlackList() *BlackList {
	blacklist := gcache.New(1024).LRU().Expiration(30 * time.Second).Build()

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
			if !b.cache.Has(pid) {
				b.cache.Set(pid, true)
			}
		}
	}
}

func (b *BlackList) exists(pid peer.ID) bool {
	return b.cache.Has(pid)
}