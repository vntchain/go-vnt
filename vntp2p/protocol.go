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
	// peer信息只获取1次即可
	log.Debug("Stream data coming...")
	peer := server.GetPeerByRemoteID(s)
	if peer == nil {
		log.Debug("HandleStream", "localPeerID", s.Conn().LocalPeer(), "remotePeerID", s.Conn().RemotePeer(), "this remote peer is nil, don't handle it")
		_ = s.Reset()
		return
	}

	// 发生错误时才会退出
	defer func() {
		log.Debug("HandleStream reset stream before exit")
		peer.Reset()
	}()

	// stream未关闭则连接正常可持续读取消息
	for {
		// 读取消息
		msgHeaderByte := make([]byte, MessageHeaderLength)
		_, err := io.ReadFull(s, msgHeaderByte)
		if err != nil {
			log.Error("HandleStream", "read msg header error", err, "peer", peer.RemoteID().ToString())
			notifyError(peer.msgers, err)
			return
		}
		bodySize := binary.LittleEndian.Uint32(msgHeaderByte)

		msgBodyByte := make([]byte, bodySize)
		_, err = io.ReadFull(s, msgBodyByte)
		if err != nil {
			log.Error("HandleStream", "read msg Body error", err, "peer", peer.RemoteID().ToString())
			notifyError(peer.msgers, err)
			return
		}
		msgBody := &MsgBody{Payload: &rlp.EncReader{}}
		err = json.Unmarshal(msgBodyByte, msgBody)
		if err != nil {
			log.Error("HandleStream", "unmarshal msg Body error", err, "peer", peer.RemoteID().ToString())
			notifyError(peer.msgers, err)
			return
		}
		msgBody.ReceivedAt = time.Now()

		// 传递给msger
		var msgHeader MsgHeader
		copy(msgHeader[:], msgHeaderByte)

		msg := Msg{
			Header: msgHeader,
			Body:   *msgBody,
		}
		if msger, ok := peer.msgers[msgBody.ProtocolID]; ok { // this node support protocolID
			// 非阻塞向上层协议传递消息，如果2s还未被读取，认为上层协议有故障
			select {
			case msger.in <- msg:
				log.Trace("HandleStream send message to messager success")
			case <-time.NewTimer(time.Second * 2).C:
				log.Trace("HandleStream send message to messager timeout")
			}
		} else {
			log.Warn("HandleStream", "receive unknown message", msg)
		}
	}
}

func notifyError(msgers map[string]*VNTMsger, err error) {
	for _, m := range msgers {
		m.err <- err
	}
}
