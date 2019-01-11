package types

import (
	"math/big"

	"fmt"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/crypto/sha3"
	"github.com/vntchain/go-vnt/rlp"
)

type BftMsgType uint8

const (
	BftPreprepareMessage BftMsgType = iota
	BftPrepareMessage
	BftCommitMessage
)

func (msg BftMsgType) String() string {
	switch msg {
	case BftPreprepareMessage:
		return "BftPreprepareMessage"
	case BftPrepareMessage:
		return "BftPrepareMessage"
	case BftCommitMessage:
		return "BftCommitMessage"
	default:
		return "Unknown bft message type"
	}
}

type BftMsg struct {
	BftType BftMsgType
	Msg     ConsensusMsg
}

type ConsensusMsg interface {
	Type() BftMsgType
	GetRound() uint32
	GetBlockNum() *big.Int
	Hash() common.Hash
}

type PreprepareMsg struct {
	Round uint32
	Block *Block
}

func (msg *PreprepareMsg) Type() BftMsgType {
	return BftPreprepareMessage
}
func (msg *PreprepareMsg) GetBlockNum() *big.Int {
	return msg.Block.Number()
}

func (msg *PreprepareMsg) GetRound() uint32 {
	return msg.Round
}

func (msg *PreprepareMsg) Hash() (hash common.Hash) {
	hasher := sha3.NewKeccak256()

	rlp.Encode(hasher, []interface{}{
		msg.Round,
		msg.Block.Hash(),
	})

	hasher.Sum(hash[:0])
	return
}

type PrepareMsg struct {
	Round       uint32
	PrepareAddr common.Address
	BlockNumber *big.Int
	BlockHash   common.Hash
	PrepareSig  []byte
}

func (msg *PrepareMsg) Type() BftMsgType {
	return BftPrepareMessage
}

func (msg *PrepareMsg) GetBlockNum() *big.Int {
	return msg.BlockNumber
}

func (msg *PrepareMsg) GetRound() uint32 {
	return msg.Round
}

func (msg *PrepareMsg) Hash() (hash common.Hash) {
	hasher := sha3.NewKeccak256()

	// 加入消息类型是为了区别Prepare消息和commit消息
	rlp.Encode(hasher, []interface{}{
		BftPrepareMessage,
		msg.Round,
		msg.PrepareAddr,
		msg.BlockNumber,
		msg.BlockHash,
	})

	hasher.Sum(hash[:0])
	return
}

type CommitMsg struct {
	Round       uint32
	Commiter    common.Address
	BlockNumber *big.Int
	BlockHash   common.Hash
	CommitSig   []byte
}

func (msg *CommitMsg) Type() BftMsgType {
	return BftCommitMessage
}

func (msg *CommitMsg) GetBlockNum() *big.Int {
	return msg.BlockNumber
}

func (msg *CommitMsg) GetRound() uint32 {
	return msg.Round
}

func (msg *CommitMsg) Hash() (hash common.Hash) {
	hasher := sha3.NewKeccak256()

	rlp.Encode(hasher, []interface{}{
		BftCommitMessage,
		msg.Round,
		msg.Commiter,
		msg.BlockNumber,
		msg.BlockHash,
	})

	hasher.Sum(hash[:0])
	return
}

// Size returns the approximate memory used by all internal contents.
func (msg *CommitMsg) Size() int {
	return len(msg.Commiter) + len(msg.CommitSig) + common.HashLength + msg.BlockNumber.BitLen()/8
}

func (msg *CommitMsg) Dump() {
	fmt.Println("----------------- Dump Commit Message -----------------")
	fmt.Printf("committer: %s\n", msg.Commiter.String())
	fmt.Printf("number: %d\nround: %d\n", msg.BlockNumber.Int64(), msg.Round)
	fmt.Printf("hash: %s\n", msg.BlockHash.String())
}

func CopyCmtMsg(msg *CommitMsg) *CommitMsg {
	cpy := *msg
	if cpy.BlockNumber = new(big.Int); msg.BlockNumber != nil {
		cpy.BlockNumber.Set(msg.BlockNumber)
	}

	if len(msg.CommitSig) > 0 {
		cpy.CommitSig = make([]byte, len(msg.CommitSig))
		copy(cpy.CommitSig, msg.CommitSig)
	}
	return &cpy
}
