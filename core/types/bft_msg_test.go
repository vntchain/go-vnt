package types

import (
	"bytes"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/rlp"
	"math/big"
	"reflect"
	"testing"
)

func check(t *testing.T, f string, got, want interface{}) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s mismatch: got %v, want %v", f, got, want)
	}
}

func TestPreprepareMsgEncoding(t *testing.T) {
	header := &Header{
		Difficulty: big.NewInt(131072),
		GasLimit:   uint64(3141592),
		GasUsed:    uint64(21000),
		Coinbase:   common.HexToAddress("8888f1f195afa192cfee860698584c030f4c9db1"),
		Root:       common.HexToHash("ef1552a40b7165c3cd773806b9e0c165b75356e0314bf0706f279c729f51e017"),
		Time:       big.NewInt(1426516743),
		CmtMsges:   make([]*CommitMsg, 0),
	}

	blk := NewBlockWithHeader(header)
	originMsg := &PreprepareMsg{
		0,
		blk,
	}

	var (
		msgEnc []byte
		err    error
	)
	if msgEnc, err = rlp.EncodeToBytes(originMsg); err != nil {
		t.Fatal("encode error:", err)
	}
	// msgEnc := common.FromHex("f902660af90262f901fba00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000948888f1f195afa192cfee860698584c030f4c9db1a0ef1552a40b7165c3cd773806b9e0c165b75356e0314bf0706f279c729f51e017a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000b90100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008302000080832fefd8825208845506eb0780a00000000000000000000000000000000000000000000000000000000000000000880000000000000000c080f861f85f800a82c35094095e7baea6a6c7c4c2dfeb977efac326af552d870a801ba09bea4c4daac7c7c52e093e6a4c35dbbcf8856f1af7b059ba20253e70848d094fa08a8fae537ce25ed8cb5af9adac3f141af69bd515bd2ba031522df09b97dd72b1c0")
	var msg PreprepareMsg

	if err := rlp.DecodeBytes(msgEnc, &msg); err != nil {
		t.Fatal("decode error: ", err)
	}

	check(t, "Difficulty", msg.Block.Difficulty(), header.Difficulty)
	check(t, "GasLimit", msg.Block.GasLimit(), header.GasLimit)
	check(t, "GasUsed", msg.Block.GasUsed(), header.GasUsed)
	check(t, "Coinbase", msg.Block.Coinbase(), header.Coinbase)
	check(t, "Root", msg.Block.Root(), header.Root)
	check(t, "Hash", msg.Block.Hash(), header.Hash())
	check(t, "Time", msg.Block.Time(), header.Time)

	ourMsgEnc, err := rlp.EncodeToBytes(&msg)
	if err != nil {
		t.Fatal("encode error: ", err)
	}
	if !bytes.Equal(ourMsgEnc, msgEnc) {
		t.Errorf("encoded blockProposalMsg mismatch:\ngot:  %x\nwant: %x", ourMsgEnc, msgEnc)
	}
}

func TestPrepareMsgEncoding(t *testing.T) {
	msgEnc := common.FromHex("f87b0a948888f1f195afa192cfee860698584c030f4c9db10aa0503290d0c4dd2d72202521e4701e89daecf048d400b2fbb8cbad1f15a4ec2e8db8419bea4c4daac7c7c52e093e6a4c35dbbcf8856f1af7b059ba20253e70848d094f8a8fae537ce25ed8cb5af9adac3f141af69bd515bd2ba031522df09b97dd72b100")
	var msg PrepareMsg
	if err := rlp.DecodeBytes(msgEnc, &msg); err != nil {
		t.Fatal("decode error: ", err)
	}
	check(t, "PrepareAddr", msg.PrepareAddr, common.HexToAddress("8888f1f195afa192cfee860698584c030f4c9db1"))
	check(t, "BlockNumber", msg.BlockNumber, big.NewInt(10))
	check(t, "BlockHash", msg.BlockHash, common.HexToHash("503290d0c4dd2d72202521e4701e89daecf048d400b2fbb8cbad1f15a4ec2e8d"))
	check(t, "PrepareSig", msg.PrepareSig, common.Hex2Bytes("9bea4c4daac7c7c52e093e6a4c35dbbcf8856f1af7b059ba20253e70848d094f8a8fae537ce25ed8cb5af9adac3f141af69bd515bd2ba031522df09b97dd72b100"))

	ourMsgEnc, err := rlp.EncodeToBytes(&msg)
	if err != nil {
		t.Fatal("encode error: ", err)
	}
	if !bytes.Equal(ourMsgEnc, msgEnc) {
		t.Errorf("encoded blockEndorseMsg mismatch:\ngot:  %x\nwant: %x", ourMsgEnc, msgEnc)
	}

}

func TestCommitMsgEncoding(t *testing.T) {
	msgEnc := common.FromHex("f87b0a948888f1f195afa192cfee860698584c030f4c9db10aa0503290d0c4dd2d72202521e4701e89daecf048d400b2fbb8cbad1f15a4ec2e8db8419bea4c4daac7c7c52e093e6a4c35dbbcf8856f1af7b059ba20253e70848d094f8a8fae537ce25ed8cb5af9adac3f141af69bd515bd2ba031522df09b97dd72b100")
	var msg CommitMsg
	if err := rlp.DecodeBytes(msgEnc, &msg); err != nil {
		t.Fatal("decode error: ", err)
	}
	check(t, "Commiter", msg.Commiter, common.HexToAddress("8888f1f195afa192cfee860698584c030f4c9db1"))
	check(t, "BlockNumber", msg.BlockNumber, big.NewInt(10))
	check(t, "BlockHash", msg.BlockHash, common.HexToHash("503290d0c4dd2d72202521e4701e89daecf048d400b2fbb8cbad1f15a4ec2e8d"))
	check(t, "CommitSig", msg.CommitSig, common.Hex2Bytes("9bea4c4daac7c7c52e093e6a4c35dbbcf8856f1af7b059ba20253e70848d094f8a8fae537ce25ed8cb5af9adac3f141af69bd515bd2ba031522df09b97dd72b100"))

	ourMsgEnc, err := rlp.EncodeToBytes(&msg)
	if err != nil {
		t.Fatal("encode error", err)
	}
	if !bytes.Equal(ourMsgEnc, msgEnc) {
		t.Errorf("encoded blockCommitMsg mismatch:\ngot:  %x\nwant: %x", ourMsgEnc, msgEnc)
	}
}
