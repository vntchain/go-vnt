package wavm

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"reflect"

	"github.com/vntchain/go-vnt/accounts/abi"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/vnt-wasm/wasm"
)

var (
	errBadBool = errors.New("abi: improperly encoded boolean value")
)

func GetAbi(abibyte []byte) (abi.ABI, error) {
	buf := bytes.NewBuffer(abibyte)
	res, err := abi.JSON(buf)
	return res, err
}

// reads the integer based on its kind
func readInteger(kind reflect.Kind, b []byte) interface{} {
	switch kind {
	case reflect.Uint8:
		return uint64(b[len(b)-1])
	case reflect.Uint16:
		return uint64(binary.BigEndian.Uint16(b[len(b)-2:]))
	case reflect.Uint32:
		return uint64(binary.BigEndian.Uint32(b[len(b)-4:]))
	case reflect.Uint64:
		return uint64(binary.BigEndian.Uint64(b[len(b)-8:]))
	case reflect.Int8:
		return uint64(b[len(b)-1])
	case reflect.Int16:
		return uint64(binary.BigEndian.Uint16(b[len(b)-2:]))
	case reflect.Int32:
		return uint64(binary.BigEndian.Uint32(b[len(b)-4:]))
	case reflect.Int64:
		return uint64(binary.BigEndian.Uint64(b[len(b)-8:]))
	default:
		return new(big.Int).SetBytes(b)
	}
}

// reads a bool
func readBool(word []byte) (uint64, error) {
	for _, b := range word[:31] {
		if b != 0 {
			return uint64(0), errBadBool
		}
	}
	switch word[31] {
	case 0:
		return uint64(0), nil
	case 1:
		return uint64(1), nil
	default:
		return uint64(0), errBadBool
	}
}

// packNum packs the given number (using the reflect value) and will cast it to appropriate number representation
func packNum(value reflect.Value) ([]byte, error) {
	switch kind := value.Kind(); kind {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return abi.U256(new(big.Int).SetUint64(value.Uint())), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return abi.U256(big.NewInt(value.Int())), nil
	case reflect.Ptr:
		return abi.U256(value.Interface().(*big.Int)), nil
	default:
		return nil, fmt.Errorf("abi: fatal error")
	}
}

func packBytesSlice(bytes []byte, l int) ([]byte, error) {
	len, err := packNum(reflect.ValueOf(l))
	return append(len, common.RightPadBytes(bytes, (l+31)/32*32)...), err
}

func lengthPrefixPointsTo(index int, output []byte) (start int, length int, err error) {
	bigOffsetEnd := big.NewInt(0).SetBytes(output[index : index+32])
	bigOffsetEnd.Add(bigOffsetEnd, common.Big32)
	outputLength := big.NewInt(int64(len(output)))
	if bigOffsetEnd.Cmp(outputLength) > 0 {
		return 0, 0, fmt.Errorf("abi: cannot marshal in to go slice: offset %v would go over slice boundary (len=%v)", bigOffsetEnd, outputLength)
	}

	if bigOffsetEnd.BitLen() > 63 {
		return 0, 0, fmt.Errorf("abi offset larger than int64: %v", bigOffsetEnd)
	}

	offsetEnd := int(bigOffsetEnd.Uint64())
	lengthBig := big.NewInt(0).SetBytes(output[offsetEnd-32 : offsetEnd])

	totalSize := big.NewInt(0)
	totalSize.Add(totalSize, bigOffsetEnd)
	totalSize.Add(totalSize, lengthBig)
	if totalSize.BitLen() > 63 {
		return 0, 0, fmt.Errorf("abi length larger than int64: %v", totalSize)
	}

	if totalSize.Cmp(outputLength) > 0 {
		return 0, 0, fmt.Errorf("abi: cannot marshal in to go type: length insufficient %v require %v", outputLength, totalSize)
	}
	start = int(bigOffsetEnd.Uint64())
	length = int(lengthBig.Uint64())
	return
}

func MutableFunction(abi abi.ABI, module *wasm.Module) Mutable {
	mutable := Mutable{}
	for k, v := range module.Export.Entries {
		if m, ok := abi.Methods[k]; ok {
			mutable[v.Index] = !m.Const
		}
	}
	return mutable
}
