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
	"encoding/binary"
	"encoding/json"
	"io"
	"time"

	inet "github.com/libp2p/go-libp2p-net"
	libp2p "github.com/libp2p/go-libp2p-peer"
	"github.com/vntchain/go-vnt/log"
	"github.com/vntchain/go-vnt/rlp"
)

// Protocol 以太坊自带代码，别的地方要用到
// 目前依然沿用eth的子协议结构，减少上层的改动
type Protocol struct {
	Name     string
	Version  uint
	Length   uint64
	Run      func(peer *Peer, rw MsgReadWriter) error
	NodeInfo func() interface{}
	PeerInfo func(id libp2p.ID) interface{}
}

// HandleStream handle all message which is from anywhere
// 主、被动连接都走的流程
func (server *Server) HandleStream(s inet.Stream) {
	// 发生错误时才会退出
	defer func() {
		log.Debug("HandleStream reset stream before exit")
		s.Reset()
	}()

	// peer信息只获取1次即可
	log.Debug("p2p-test, stream data comming")
	peer := server.GetPeerByRemoteID(s)
	if peer == nil {
		log.Debug("HandleStream", "localPeerID", s.Conn().LocalPeer(), "remotePeerID", s.Conn().RemotePeer(), "this remote peer is nil, don't handle it")
		return
	}

	// stream未关闭则连接正常可持续读取消息
	for {
		// 读取消息
		msgHeaderByte := make([]byte, MessageHeaderLength)
		_, err := io.ReadFull(s, msgHeaderByte)
		if err != nil {
			log.Error("handleStream", "read header error", err, "peer", peer.RemoteID().ToString())
			notifyError(peer.messenger, err)
			return
		}
		bodySize := binary.LittleEndian.Uint32(msgHeaderByte)

		msgBodyByte := make([]byte, bodySize)
		_, err = io.ReadFull(s, msgBodyByte)
		if err != nil {
			log.Error("handleStream", "read msgBody error", err, "peer", peer.RemoteID().ToString())
			notifyError(peer.messenger, err)
			return
		}
		msgBody := &MsgBody{Payload: &rlp.EncReader{}}
		err = json.Unmarshal(msgBodyByte, msgBody)
		if err != nil {
			log.Error("handleSteam", "unmarshal msgBody error", err, "peer", peer.RemoteID().ToString())
			notifyError(peer.messenger, err)
			return
		}
		msgBody.ReceivedAt = time.Now()

		// 传递给messenger
		var msgHeader MsgHeader
		copy(msgHeader[:], msgHeaderByte)

		msg := Msg{
			Header: msgHeader,
			Body:   *msgBody,
		}
		if messenger, ok := peer.messenger[msgBody.ProtocolID]; ok { // this node support protocolID
			messenger.in <- msg
		} else {
			log.Warn("handleStream", "receive Unknown Message", msg)
		}
	}
}

func notifyError(messengers map[string]*VNTMessenger, err error) {
	log.Trace("notifyError enter")
	defer log.Trace("notifyError exit")
	for _, m := range messengers {
		m.err <- err
	}
}
