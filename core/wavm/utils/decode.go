package utils

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/vntchain/go-vnt/core/wavm/contract"
	"github.com/vntchain/go-vnt/rlp"
)

func DecodeContractCode(input []byte) (contract.WasmCode, []byte, error) {
	buf := bytes.NewReader(input)
	magic := make([]byte, 4)
	_, err := buf.Read(magic)
	if err != nil {
		return contract.WasmCode{}, nil, err
	}
	magicBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(magicBytes, MAGIC)
	// magicNum := binary.LittleEndian.Uint32(magic)
	if !bytes.Equal(magic, magicBytes) {
		return contract.WasmCode{}, nil, errors.New("Magic number mismatch")
	}

	cps := []byte{}
	err = rlp.Decode(buf, &cps)
	if err != nil {
		return contract.WasmCode{}, nil, err
	}
	decom, err := DeCompress(cps)
	if err != nil {
		return contract.WasmCode{}, nil, err
	}

	dec := contract.WasmCode{}
	err = rlp.Decode(bytes.NewReader(decom), &dec)
	if err != nil {
		return contract.WasmCode{}, nil, err
	}
	return dec, input[int(buf.Size())-buf.Len():], nil
}
