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

package utils_test

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/vntchain/go-vnt/core/wavm/utils"
)

func TestSplit(t *testing.T) {
	teststr := []byte("111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111")
	n, res := utils.Split(teststr)
	testres := []byte{}
	for i := 0; i < n; i++ {
		testres = append(testres, res[i]...)
	}

	if !bytes.Equal(teststr, new(big.Int).SetBytes(testres).Bytes()) {
		t.Errorf("want %s | get %s", teststr, new(big.Int).SetBytes(testres).Bytes())
	}
}

func TestGetU256(t *testing.T) {

	var mem []byte
	bigint := utils.GetU256(mem)
	if bigint.String() != "0" {
		t.Fatalf("want %s |get %s", "0", bigint.String())
	}

	mem = []byte(nil)
	bigint = utils.GetU256(mem)
	if bigint.String() != "0" {
		t.Fatalf("want %s |get %s", "0", bigint.String())
	}

	mem = []byte{}
	bigint = utils.GetU256(mem)
	if bigint.String() != "0" {
		t.Fatalf("want %s |get %s", "0", bigint.String())
	}

	mem = []byte("111111111")
	bigint = utils.GetU256(mem)
	if bigint.String() != "111111111" {
		t.Fatalf("want %s |get %s", "111111111", bigint.String())
	}

}
