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
