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
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"

	peer "github.com/libp2p/go-libp2p-peer"
	ma "github.com/multiformats/go-multiaddr"
)

// 包内都用peerID,对外方法使用NodeID
const NodeIDBits = 512

type NodeID [NodeIDBits / 8]byte

// type Node struct {
// 	IP       net.IP // len 4 for IPv4 or 16 for IPv6
// 	UDP, TCP uint16 // port numbers
// 	ID       NodeID // the node's public key

// 	sha common.Hash

// 	// Time when the node was added to the table.
// 	addedAt time.Time
// }

type Node struct {
	Addr ma.Multiaddr
	Id   peer.ID
}

func (n *Node) String() string {
	return n.Addr.String() + "/ipfs/" + n.Id.ToString()
}

func NewNode(id peer.ID, ip net.IP, udpPort, tcpPort uint16) *Node {
	peerid := id
	target := ""
	if ipv4 := ip.To4(); ipv4 != nil {
		target += "/ip4/" + ip.String() + "/tcp/" + strconv.Itoa(int(tcpPort)) + "/ipfs/" + peerid.ToString()
	} else {
		target += "/ip6/" + ip.String() + "/tcp/" + strconv.Itoa(int(tcpPort)) + "/ipfs/" + peerid.ToString()
	}

	targetAddr, peerid, err := GetAddr(target)
	if err != nil {
		return nil
	}

	return &Node{Addr: targetAddr, Id: peerid}
}

func ParseNode(url string) (*Node, error) {
	addr, peerid, err := GetAddr(url)
	if err != nil {
		// log
		return nil, err
	}

	return &Node{Addr: addr, Id: peerid}, nil
}

func MustParseNode(rawurl string) *Node {
	n, err := ParseNode(rawurl)
	if err != nil {
		panic("invalid node URL: " + err.Error())
	}
	return n
}

// for toml unmarshal
func (n *Node) UnmarshalText(data []byte) error {
	var err error
	var nd *Node

	nd, err = ParseNode(string(data))
	n.Addr = nd.Addr
	n.Id = nd.Id

	return err
}

// for toml marshal
func (n *Node) MarshalText() ([]byte, error) {
	url := n.String()
	ret := make([]byte, len(url))
	copy(ret, url)

	return ret, nil
}

/* func PubkeyID(pub *ecdsa.PublicKey) NodeID {
	// var id NodeID
	// pbytes := elliptic.Marshal(pub.Curve, pub.X, pub.Y)
	// if len(pbytes)-1 != len(id) {
	// 	panic(fmt.Errorf("need %d bit pubkey, got %d bits", (len(id)+1)*8, len(pbytes)))
	// }
	// copy(id[:], pbytes[1:])
	// return id
	id, err := peer.IDFromPublicKey(pub)
	if err != nil {
		panic("wrong publick key")
	}
	return PeerIDtoNodeID(id)
} */

func PeerIDtoNodeID(n peer.ID) NodeID {
	var id NodeID
	copy(id[:], []byte(n))
	return id
}

func (n NodeID) PeerID() peer.ID {
	return peer.ID(n.Bytes())
}

func (n NodeID) Pubkey() (*ecdsa.PublicKey, error) {
	// 通过ID如何生成公钥

	return n.PeerID().ExtractPublicKey()
}

func (n NodeID) Bytes() []byte {
	return n[:]
}

// NodeID prints as a long hexadecimal number.
func (n NodeID) String() string {
	return string(n.Bytes())
}

// The Go syntax representation of a NodeID is a call to HexID.
func (n NodeID) GoString() string {
	return fmt.Sprintf("discover.HexID(\"%x\")", n.Bytes())
}

// TerminalString returns a shortened hex string for terminal logging.
func (n NodeID) TerminalString() string {
	return hex.EncodeToString(n.Bytes()[:8])
}

// MarshalText implements the encoding.TextMarshaler interface.
func (n NodeID) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(n.Bytes())), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (n *NodeID) UnmarshalText(text []byte) error {
	id, err := HexID(string(text))
	if err != nil {
		return err
	}
	*n = id
	return nil
}

func HexID(in string) (NodeID, error) {
	var id NodeID
	b, err := hex.DecodeString(strings.TrimPrefix(in, "0x"))
	if err != nil {
		return id, err
	} else if len(b) != len(id) {
		return id, fmt.Errorf("wrong length, want %d hex chars", len(id)*2)
	}
	copy(id[:], b)
	return id, nil
}
