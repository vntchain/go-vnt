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
func (server *Server) HandleStream(s inet.Stream) {
	for {
		log.Info("yhx-test, stream data comming")
		peer := server.GetPeerByRemoteID(s)
		if peer == nil {
			log.Info("HandleStream", "localPeerID", s.Conn().LocalPeer(), "remotePeerID", s.Conn().RemotePeer(), "this remote peer is nil, don't handle it")
			return
		}
		msgHeaderByte := make([]byte, MessageHeaderLength)
		_, err := io.ReadFull(s, msgHeaderByte)
		if err != nil {
			//log.Error("handleStream", "read error", err)
			notifyError(peer.messenger, err)
			return
		}
		bodySize := binary.LittleEndian.Uint32(msgHeaderByte)

		msgBodyByte := make([]byte, bodySize)
		_, err = io.ReadFull(s, msgBodyByte)
		if err != nil {
			log.Error("handleStream", "read msgBody error", err)
			notifyError(peer.messenger, err)
			return
		}
		msgBody := &MsgBody{Payload: &rlp.EncReader{}}
		err = json.Unmarshal(msgBodyByte, msgBody)
		if err != nil {
			log.Error("handleSteam", "unmarshal msgBody error", err)
			notifyError(peer.messenger, err)
			return
		}
		msgBody.ReceivedAt = time.Now()
		//log.Info("yhx-test", "RECEIVED MESSAGE", msgBody)

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

		//handler, err := msgBody.handleForMsgType()
		//if err != nil {
		//	log.Error("handleStream", "handleForMsgType error", err)
		//	return
		//}
		//handler()
	}
}

func notifyError(messengers map[string]*VNTMessenger, err error) {
	for _, m := range messengers {
		m.err <- err
	}
}
