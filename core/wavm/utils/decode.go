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
