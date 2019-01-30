package utils

import (
	"bytes"
	"errors"

	"github.com/vntchain/go-vnt/core/wavm/contract"
	"github.com/vntchain/go-vnt/rlp"
)

func DecodeContractCode(input []byte) (contract.WasmCode, []byte, error) {
	magic, _ := ReadMagic(input)
	if magic != MAGIC {
		return contract.WasmCode{}, nil, errors.New("Magic number mismatch")
	}
	input = input[4:]
	buf := bytes.NewReader(input)
	cps := []byte{}
	err := rlp.Decode(buf, &cps)
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
