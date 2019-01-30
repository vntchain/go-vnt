package utils

import (
	"encoding/binary"
	"math"
	"math/big"
	"strings"

	"github.com/vntchain/go-vnt/common"
	m "github.com/vntchain/go-vnt/common/math"
	"github.com/vntchain/go-vnt/crypto"
	"github.com/vntchain/vnt-wasm/wasm"
)

var endianess = binary.LittleEndian

func BytesToF32(bytes []byte) float32 {
	bits := endianess.Uint32(bytes)

	return math.Float32frombits(bits)
}

func BytesToF64(bytes []byte) float64 {
	bits := endianess.Uint64(bytes)
	return math.Float64frombits(bits)
}

func BytesToI64(bytes []byte) uint64 {
	bits := endianess.Uint64(bytes)
	return bits
}

func I32ToBytes(i32 uint32) []byte {
	bytes := make([]byte, 4)
	endianess.PutUint32(bytes, i32)
	return bytes
}

func I64ToBytes(i64 uint64) []byte {
	bytes := make([]byte, 8)
	endianess.PutUint64(bytes, i64)
	return bytes
}
func F32ToBytes(f32 float32) []byte {
	bits := math.Float32bits(f32)
	bytes := make([]byte, 4)
	endianess.PutUint32(bytes, bits)
	return bytes
}

func F64ToBytes(f64 float64) []byte {
	bits := math.Float64bits(f64)
	bytes := make([]byte, 8)
	endianess.PutUint64(bytes, bits)
	return bytes
}

func Split(data []byte) (int, [][]byte) {
	l := len(data)
	n := math.Ceil(float64(l) / 32.0)
	h := make([]byte, int(n)*32)
	copy(h[len(h)-len(data):], data)
	res := [][]byte{}
	for i := 0; i < int(n); i++ {
		res = append(res, h[i*32:(i+1)*32])
	}
	return int(n), res
}

func ArrayLengthKey(symbol common.Hash) common.Hash {
	return common.BigToHash(new(big.Int).Add(symbol.Big(), common.BytesToHash([]byte("length")).Big()))
}

func MapLocation(sym, key []byte) common.Hash {
	return crypto.Keccak256Hash(sym, key)
}
func ArrayLocation(sym string, index int64, length int64) {}

func GetU256(mem []byte) *big.Int {
	bigint := new(big.Int)
	var toStr string
	if len(mem) == 0 {
		toStr = "0"
	} else {
		toStr = string(mem)
	}
	_, success := bigint.SetString(toStr, 10)
	if success == false {
		panic("Illegal uint256 input " + toStr)
	}
	return m.U256(bigint)
}

func GetIndex(m *wasm.Module) (writeIndex int, readIndex int, gasIndex int) {
	writeIndex = -1
	readIndex = -1
	gasIndex = -1
	for i, v := range m.Import.Entries {
		if v.FieldName == "WriteWithPointer" {
			writeIndex = i
		} else if v.FieldName == "ReadWithPointer" {
			readIndex = i
		} else if v.FieldName == "AddGas" {
			gasIndex = i
		}
		if writeIndex != -1 && readIndex != -1 && gasIndex != -1 {
			return
		}
	}
	return
}

func IsAddress(data []byte) (bool, []byte) {
	magic := "address1537182776"
	if strings.HasPrefix(string(data), magic) {
		return true, []byte(string(data)[len(magic):])
	}
	return false, nil
}

func IsU256(data []byte) (bool, []byte) {
	magic := "u2561537182776"
	if strings.HasPrefix(string(data), magic) {
		return true, []byte(string(data)[len(magic):])
	}
	return false, nil
}
