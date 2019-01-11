package vntp2p

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/vntchain/go-vnt/log"
	"github.com/vntchain/go-vnt/rlp"
	inet "github.com/libp2p/go-libp2p-net"
)

type MsgReadWriter interface {
	MsgReader
	MsgWriter
}

type MsgReader interface {
	ReadMsg() (Msg, error)
}

type MsgWriter interface {
	WriteMsg(Msg) error
}

// MessageHeaderLength define message header length
const MessageHeaderLength = 5

// MessageType define vnt p2p protocol message type
type MessageType uint64

const (
	// GoodMorning say good morning protocol
	GoodMorning MessageType = iota
	// GoodAfternoon say good afternoon protocol
	GoodAfternoon
	// GoodNight say good night protocol
	GoodNight
)

// Msg message struct
type Msg struct {
	Header MsgHeader
	Body   MsgBody
}

// MsgHeader store the size of MsgBody
type MsgHeader [MessageHeaderLength]byte

// MsgBody message body
type MsgBody struct {
	ProtocolID  string //Protocol name
	Type        MessageType
	ReceivedAt  time.Time
	PayloadSize uint32
	Payload     io.Reader
}

// GoodMorningMsg message for goodmorning protocol
type GoodMorningMsg struct {
	Greet     string
	Timestamp string
}

// HandleMessage implement VNTMessage interface
func (gmm *GoodMorningMsg) HandleMessage() error {
	fmt.Printf("Receive Message: greet = %s, at %s\n", gmm.Greet, gmm.Timestamp)
	return nil
}

type msgHandler func() error

/*func (mb *MsgBody) handleForMsgType() (msgHandler, error) {
	switch mb.Type {
	case GoodMorning:
		gmm := &GoodMorningMsg{}
		err := json.Unmarshal(mb.Payload, gmm)
		if err != nil {
			return nil, err
		}
		return gmm.HandleMessage, nil
	default:
		return nil, errors.New("invalid message type")
	}
}*/

// Send is used to send message payload with specific messge type
func Send(w MsgWriter, protocolID string, msgType MessageType, data interface{}) error {
	// 还是要使用rlp进行序列化，因为类型多变，rlp已经有完整的支持
	log.Info("yhx-test", "send message", data)
	size, r, err := rlp.EncodeToReader(data)
	if err != nil {
		log.Error("Send()", "rlp encode error", err)
		return err
	}

	msgBody := MsgBody{
		ProtocolID:  protocolID,
		Type:        msgType,
		PayloadSize: uint32(size),
		Payload:     r,
	}
	msgBodyByte, err := json.Marshal(msgBody)
	if err != nil {
		log.Error("Send()", "marshal msgBody error", err)
		return err
	}
	msgBodySize := len(msgBodyByte)
	msgHeaderByte := make([]byte, MessageHeaderLength)
	binary.LittleEndian.PutUint32(msgHeaderByte, uint32(msgBodySize))

	var msgHeader MsgHeader
	copy(msgHeader[:], msgHeaderByte)

	msg := Msg{
		Header: msgHeader,
		Body:   msgBody,
	}

	return w.WriteMsg(msg)
}

// SendItems can send many payload in one function call
func SendItems(w MsgWriter, protocolID string, msgType MessageType, elems ...interface{}) error {
	return Send(w, protocolID, msgType, elems)
}

// Decode using json unmarshal decode msg payload
func (msg Msg) Decode(val interface{}) error {
	s := rlp.NewStream(msg.Body.Payload, uint64(msg.Body.PayloadSize))
	err := s.Decode(val)
	if err != nil {
		log.Error("Decode()", "err", err, "message type", msg.Body.Type, "payload size", msg.Body.PayloadSize)
		return err
	}
	return nil
}

// GetBodySize get message body size in uint32
func (msg *Msg) GetBodySize() uint32 {
	header := msg.Header
	bodySize := binary.LittleEndian.Uint32(header[:])
	return bodySize
}

// VNTMessenger vnt chain message readwriter
type VNTMessenger struct {
	protocol Protocol
	in       chan Msg
	err      chan error
	w        inet.Stream
}

// WriteMsg implement MsgReadWriter interface
func (rw *VNTMessenger) WriteMsg(msg Msg) (err error) {
	//if uint64(msg.Body.Type) >= rw.Length {
	//	return newPeerError(errInvalidMsgCode, "not handled")
	//}
	// 暂时先不管主动关闭需要告知对方的情况，目前聚焦于发送消息这件基本工作
	msgHeaderByte := msg.Header[:]
	msgBodyByte, err := json.Marshal(msg.Body)
	if err != nil {
		log.Error("WriteMsg()", "marshal msgbody error", err)
		return err
	}
	m := append(msgHeaderByte, msgBodyByte...)
	//log.Info("yhx-test", "MESSAGE", string(m))

	_, err = rw.w.Write(m)
	if err != nil {
		log.Error("WriteMsg()", "write msg error", err)
		return err
	}
	return nil
}

// ReadMsg implement MsgReadWriter interface
func (rw *VNTMessenger) ReadMsg() (Msg, error) {
	select {
	case msg := <-rw.in:
		log.Info("yhx-test", "incoming message", msg)
		return msg, nil
	case err := <-rw.err:
		return Msg{}, err
	}
}

// ExpectMsg COMMENT: this function is just for _test.go files
// ExpectMsg reads a message from r and verifies that its
// code and encoded RLP content match the provided values.
// If content is nil, the payload is discarded and not verified.
func ExpectMsg(r MsgReader, code MessageType, content interface{}) error {
	msg, err := r.ReadMsg()
	if err != nil {
		return err
	}
	if msg.Body.Type != code {
		return fmt.Errorf("message code mismatch: got %d, expected %d", msg.Body.Type, code)
	}
	if content == nil {
		return nil
	}
	contentEnc, err := rlp.EncodeToBytes(content)
	if err != nil {
		panic("content encode error: " + err.Error())
	}
	if int(msg.Body.PayloadSize) != len(contentEnc) {
		return fmt.Errorf("message size mismatch: got %d, want %d", msg.Body.PayloadSize, len(contentEnc))
	}
	actualContent, err := ioutil.ReadAll(msg.Body.Payload)
	if err != nil {
		return err
	}
	if !bytes.Equal(actualContent, contentEnc) {
		return fmt.Errorf("message payload mismatch:\ngot:  %x\nwant: %x", actualContent, contentEnc)
	}
	return nil
}
