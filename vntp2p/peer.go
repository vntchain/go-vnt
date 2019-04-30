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
	// "crypto/ecdsa"
	// "crypto/elliptic"
	"bufio"
	"crypto/ecdsa"
	"encoding/json"
	"sync"

	"net"

	inet "github.com/libp2p/go-libp2p-net"
	libp2p "github.com/libp2p/go-libp2p-peer"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
	"github.com/vntchain/go-vnt/event"
	"github.com/vntchain/go-vnt/log"
)

type PeerEventType string

const (
	// PeerEventTypeAdd is the type of event emitted when a peer is added
	// to a p2p.Server
	PeerEventTypeAdd PeerEventType = "add"

	// PeerEventTypeDrop is the type of event emitted when a peer is
	// dropped from a p2p.Server
	PeerEventTypeDrop PeerEventType = "drop"

	// PeerEventTypeMsgSend is the type of event emitted when a
	// message is successfully sent to a peer
	PeerEventTypeMsgSend PeerEventType = "msgsend"

	// PeerEventTypeMsgRecv is the type of event emitted when a
	// message is received from a peer
	PeerEventTypeMsgRecv PeerEventType = "msgrecv"
)

// PeerEvent is an event emitted when peers are either added or dropped from
// a p2p.Server or when a message is sent or received on a peer connection
type PeerEvent struct {
	Type     PeerEventType `json:"type"`
	Peer     libp2p.ID     `json:"peer"`
	Error    string        `json:"error,omitempty"`
	Protocol string        `json:"protocol,omitempty"`
	MsgCode  *uint64       `json:"msg_code,omitempty"`
	MsgSize  *uint32       `json:"msg_size,omitempty"`
}

type PeerInfo struct {
	ID      string   `json:"id"`   // Unique node identifier (also the encryption key)
	Name    string   `json:"name"` // Name of the node, including client type, version, OS, custom data
	Caps    []string `json:"caps"` // Sum-protocols advertised by this particular peer
	Network struct {
		LocalAddress  string `json:"localAddress"`  // Local endpoint of the TCP data connection
		RemoteAddress string `json:"remoteAddress"` // Remote endpoint of the TCP data connection
		Inbound       bool   `json:"inbound"`
		Trusted       bool   `json:"trusted"`
		Static        bool   `json:"static"`
	} `json:"network"`
	Protocols map[string]interface{} `json:"protocols"` // Sub-protocol specific metadata fields
}

type Peer struct {
	rw        inet.Stream
	log       log.Logger
	events    *event.Feed
	err       chan error
	closed    bool
	messenger map[string]*VNTMessenger // protocolName - vntMessenger
	server	  *Server
	wg        sync.WaitGroup
	// need to add wg
}

func newPeer(conn *Stream, server *Server) *Peer {
	m := make(map[string]*VNTMessenger)
	for i := range conn.Protocols {
		proto := conn.Protocols[i]
		vntMessenger := &VNTMessenger{
			protocol: proto,
			in:       make(chan Msg),
			err:      make(chan error),
			w:        conn.Conn,
		}
		m[proto.Name] = vntMessenger
	}

	p := &Peer{
		rw:        conn.Conn,
		log:       log.New(),
		err:       make(chan error),
		closed:    false,
		messenger: m,
		server:	   server,
	}
	for _, msger := range p.messenger {
		msger.peerPointer = p
	}

	return p
}

// LocalID return local PeerID for upper application
func (p *Peer) LocalID() libp2p.ID {
	return p.rw.Conn().LocalPeer()
}

// RemoteID return remote PeerID for upper application
func (p *Peer) RemoteID() libp2p.ID {
	return p.rw.Conn().RemotePeer()
}

func (p *Peer) Log() log.Logger {
	return p.log
}

func (p *Peer) LocalAddr() net.Addr {
	lma := p.rw.Conn().LocalMultiaddr()
	return parseMultiaddr(lma)
}

func (p *Peer) RemoteAddr() net.Addr {
	rma := p.rw.Conn().RemoteMultiaddr()
	return parseMultiaddr(rma)
}

func parseMultiaddr(maddr ma.Multiaddr) net.Addr {
	network, host, err := manet.DialArgs(maddr)
	if err != nil {
		log.Error("parseMultiaddr()", "dialArgs error", err)
		return nil
	}

	switch network {
	case "tcp", "tcp4", "tcp6":
		na, err := net.ResolveTCPAddr(network, host)
		if err != nil {
			log.Error("parseMultiaddr()", "resolveTCPAddr error", err)
			return nil
		}
		return na
	case "udp", "udp4", "udp6":
		na, err := net.ResolveUDPAddr(network, host)
		if err != nil {
			log.Error("parseMultiaddr()", "ResolveUDPAddr error", err)
			return nil
		}
		return na
	case "ip", "ip4", "ip6":
		na, err := net.ResolveIPAddr(network, host)
		if err != nil {
			log.Error("parseMultiaddr()", "ResolveIPAddr error", err)
			return nil
		}
		return na
	}
	log.Error("parseMultiaddr()", "network not supported", network)
	return nil
}

func (p *Peer) Disconnect(reason DiscReason) {
	// test for it
	p.rw.Conn().Close()
	// p.rw.Close()
}

func (p *Peer) Info() *PeerInfo {
	info := &PeerInfo{
		ID: p.RemoteID().String(),
	}
	info.Network.LocalAddress = p.rw.Conn().LocalMultiaddr().String()
	info.Network.RemoteAddress = p.rw.Conn().RemoteMultiaddr().String()

	// 此处暂时不处理状态
	// info.Network.Static = p.rw.Conn().RemotePeer()
	// info.Network.Trusted =
	// info.Network.Inbound =

	return info
}

func (p *Peer) run() (remoteRequested bool, err error) {
	for _, msger := range p.messenger {
		proto := msger.protocol
		m := msger
		go func() {
			p.wg.Add(1)
			err := proto.Run(p, m)
			log.Info("p2p-test", "run protocol error log", err)
			if !p.closed {
				p.err <- err
			}
			p.wg.Done()
		}()
	}

	err = <-p.err
	remoteRequested = true

	p.closed = true
	p.rw.Close()
	//log.Info("p2p-test remote peer request close, but we need to wait for other protocol", "peerid", p.RemoteID())
	//p.wg.Wait()
	log.Info("p2p-test wait complete!", "peerid", p.RemoteID())

	return remoteRequested, err
}

type Stream struct {
	Conn      inet.Stream
	Protocols []Protocol
}

//临时测试使用
func WritePackage(rw *bufio.ReadWriter, code uint64, data []byte) ([]byte, error) {
	//msg := &Msg{Code: code, Data: data, Size: uint32(len(data))}
	msg := &Msg{}
	// logg.Printf("Msg: %v", msg)
	return json.Marshal(msg)
}

func PubkeyID(pub *ecdsa.PublicKey) libp2p.ID {
	// var id NodeID
	// pbytes := elliptic.Marshal(pub.Curve, pub.X, pub.Y)
	// if len(pbytes)-1 != len(id) {
	// 	panic(fmt.Errorf("need %d bit pubkey, got %d bits", (len(id)+1)*8, len(pbytes)))
	// }
	// copy(id[:], pbytes[1:])
	// return id
	id, err := libp2p.IDFromPublicKey(pub)
	if err != nil {
		panic("wrong publick key")
	}
	return id
}
