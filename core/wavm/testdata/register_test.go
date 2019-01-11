package tests

import (
	"bytes"
	"math"
	"path/filepath"
	"testing"

	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core/wavm"
)

var (
	codepath = filepath.Join(basepath, "precompile/wasm/main7.wasm")
	abipath  = filepath.Join(basepath, "precompile/wasm/main7.json")
)

func newVM(iscreated bool) *wavm.Wavm {
	return newWavm(codepath, abipath, iscreated)
}

func TestRegister(t *testing.T) {
	vm := newVM(false)
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	aa := "value1111111111sdjsaldjasldasdlksdjlaskdjaslkdjsalkdjaslsajdsakdkasl"
	res, err := vm.Apply(pack(vm, "testWrite"), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	t.Logf("contractname %s", res)

	func() {
		vm.ChainContext.IsCreated = false

		res, err := vm.Apply(pack(vm, "testRead"), nil, mutable)
		if err != nil {
			t.Error(err)
		}
		t.Logf("testRead %s", res)
		var str string
		unPack(vm, &str, "testRead", res)
		if str != aa {
			t.Errorf("unexpected value : want %s, got %s", aa, str)
		}
		t.Logf("testRead %s", str)
	}()

}

func TestNewArray(t *testing.T) {
	vm := newVM(false)
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	res, err := vm.Apply(pack(vm, "testArrayWrite", uint64(3)), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	t.Logf("testArrayWrite %s", res)

	testread := func(index uint64, wanted common.Address) {
		vm.ChainContext.IsCreated = false
		res, err := vm.Apply(pack(vm, "testArrayRead", index), nil, mutable)
		if err != nil {
			t.Error(err)
		}
		t.Logf("testArrayRead res %v", res)
		var str common.Address
		aa := wanted
		unPack(vm, &str, "testArrayRead", res)
		if !bytes.Equal(str.Bytes(), aa.Bytes()) {
			t.Errorf("unexpected value : want %s, got %s", aa.Hex(), str.Hex())
		}
		t.Logf("testArrayRead %s", str.Hex())
	}
	testread(0, common.HexToAddress("0x000"))
	testread(1, common.HexToAddress("0x111"))
	testread(2, common.HexToAddress("0x222"))
}

func TestNewStruct(t *testing.T) {
	vm := newVM(false)
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	a1 := int32(math.MaxInt32)
	a2 := int32(math.MinInt32)
	b1 := int64(math.MaxInt64)
	b2 := int64(math.MinInt64)
	c1 := uint32(math.MaxUint32)
	c2 := uint32(0)
	d1 := uint64(math.MaxUint64)
	d2 := uint64(0)
	e1 := "test string"
	e2 := "long string long string long string long string long string long string long string long string long string long string long string long string"
	f1 := common.BytesToAddress([]byte("0xaddress"))
	f2 := common.BytesToAddress([]byte("0xaddress"))
	g1 := false
	g2 := true
	test := func(a int32, b int64, c uint32, d uint64, e string, f common.Address, g bool) {
		func() {
			funcName := "testWriteStruct"
			res, err := vm.Apply(pack(vm, funcName, a, b, c, d, e, f, g), nil, mutable)
			if err != nil {
				t.Error(err)
			}
			t.Logf("%s %s", funcName, res)
		}()

		func() {
			vm.ChainContext.IsCreated = false
			funcName := "testReadStructA"
			res, err := vm.Apply(pack(vm, funcName), nil, mutable)
			if err != nil {
				t.Error(err)
			}
			t.Logf("%s %s", funcName, res)
			var i32 int32
			unPack(vm, &i32, funcName, res)
			if i32 != a {
				t.Errorf("unexpected value : want %d, got %d", a, i32)
			} else {
				t.Logf("%s success get %d", funcName, i32)
			}

		}()

		func() {
			vm.ChainContext.IsCreated = false
			funcName := "testReadStructB"
			res, err := vm.Apply(pack(vm, funcName), nil, mutable)
			if err != nil {
				t.Error(err)
			}
			t.Logf("%s %s", funcName, res)
			var i64 int64
			unPack(vm, &i64, funcName, res)
			if i64 != b {
				t.Errorf("unexpected value : want %d, got %d", b, i64)
			} else {
				t.Logf("%s success get %d", funcName, i64)
			}

		}()

		func() {
			vm.ChainContext.IsCreated = false
			funcName := "testReadStructC"
			res, err := vm.Apply(pack(vm, funcName), nil, mutable)
			if err != nil {
				t.Error(err)
			}
			t.Logf("%s %s", funcName, res)
			var u32 uint32
			unPack(vm, &u32, funcName, res)
			if u32 != c {
				t.Errorf("unexpected value : want %d, got %d", c, u32)
			} else {
				t.Logf("%s success get %d", funcName, u32)
			}

		}()

		func() {
			vm.ChainContext.IsCreated = false
			funcName := "testReadStructD"
			res, err := vm.Apply(pack(vm, funcName), nil, mutable)
			if err != nil {
				t.Error(err)
			}
			t.Logf("%s %s", funcName, res)
			var u64 uint64
			unPack(vm, &u64, funcName, res)
			if u64 != d {
				t.Errorf("unexpected value : want %d, got %d", d, u64)
			} else {
				t.Logf("%s success get %d", funcName, u64)
			}

		}()

		func() {
			vm.ChainContext.IsCreated = false
			funcName := "testReadStructE"
			res, err := vm.Apply(pack(vm, funcName), nil, mutable)
			if err != nil {
				t.Error(err)
			}
			t.Logf("%s %s", funcName, res)
			var str string
			unPack(vm, &str, funcName, res)
			if str != e {
				t.Errorf("unexpected value : want %s, got %s", e, str)
			} else {
				t.Logf("%s success get %s", funcName, str)
			}

		}()

		func() {
			vm.ChainContext.IsCreated = false
			funcName := "testReadStructF"
			res, err := vm.Apply(pack(vm, funcName), nil, mutable)
			if err != nil {
				t.Error(err)
			}
			t.Logf("%s %s", funcName, res)
			var addr common.Address
			unPack(vm, &addr, funcName, res)
			if !bytes.Equal(addr.Bytes(), f.Bytes()) {
				t.Errorf("unexpected value : want %s, got %s", f, addr)
			} else {
				t.Logf("%s success get %s", funcName, addr)
			}

		}()

		func() {
			vm.ChainContext.IsCreated = false
			funcName := "testReadStructG"
			res, err := vm.Apply(pack(vm, funcName), nil, mutable)
			if err != nil {
				t.Error(err)
			}
			t.Logf("%s %s", funcName, res)
			var flag bool
			unPack(vm, &flag, funcName, res)
			if flag != g {
				t.Errorf("unexpected value : want %t, got %t", g, flag)
			} else {
				t.Logf("%s success get %t", funcName, flag)
			}

		}()
	}
	test(a1, b1, c1, d1, e1, f1, g1)
	test(a2, b2, c2, d2, e2, f2, g2)
}

func TestNewMapWithStruct(t *testing.T) {
	vm := newVM(false)
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	key1 := "key1"
	key2 := "key2"
	a1 := int32(math.MaxInt32)
	a2 := int32(math.MinInt32)
	b1 := int64(math.MaxInt64)
	b2 := int64(math.MinInt64)
	c1 := uint32(math.MaxUint32)
	c2 := uint32(0)
	d1 := uint64(math.MaxUint64)
	d2 := uint64(0)
	e1 := "test string"
	e2 := "long string long string long string long string long string long string long string long string long string long string long string long string"
	f1 := common.BytesToAddress([]byte("0xaddress"))
	f2 := common.BytesToAddress([]byte("0xaddress"))
	g1 := false
	g2 := true
	test := func(key string, a int32, b int64, c uint32, d uint64, e string, f common.Address, g bool) {
		func() {
			funcName := "testWriteMapStruct"
			res, err := vm.Apply(pack(vm, funcName, key, a, b, c, d, e, f, g), nil, mutable)
			if err != nil {
				t.Error(err)
			}
			t.Logf("%s %s", funcName, res)
		}()

		func() {
			vm.ChainContext.IsCreated = false
			funcName := "testReadMapStructA"
			res, err := vm.Apply(pack(vm, funcName, key), nil, mutable)
			if err != nil {
				t.Error(err)
			}
			t.Logf("%s %s", res, funcName)
			var i32 int32
			unPack(vm, &i32, funcName, res)
			if i32 != a {
				t.Errorf("unexpected value : want %d, got %d", a, i32)
			} else {
				t.Logf("%s success get %d", funcName, i32)
			}

		}()

		func() {
			vm.ChainContext.IsCreated = false
			funcName := "testReadMapStructB"
			res, err := vm.Apply(pack(vm, funcName, key), nil, mutable)
			if err != nil {
				t.Error(err)
			}
			t.Logf("%s%s", funcName, res)
			var i64 int64
			unPack(vm, &i64, funcName, res)
			if i64 != b {
				t.Errorf("unexpected value : want %d, got %d", b, i64)
			} else {
				t.Logf("%s success get %d", funcName, i64)
			}

		}()

		func() {
			vm.ChainContext.IsCreated = false
			funcName := "testReadMapStructC"
			res, err := vm.Apply(pack(vm, funcName, key), nil, mutable)
			if err != nil {
				t.Error(err)
			}
			t.Logf("%s %s", funcName, res)
			var u32 uint32
			unPack(vm, &u32, funcName, res)
			if u32 != c {
				t.Errorf("unexpected value : want %d, got %d", c, u32)
			} else {
				t.Logf("%s success get %d", funcName, u32)
			}

		}()

		func() {
			vm.ChainContext.IsCreated = false
			funcName := "testReadMapStructD"
			res, err := vm.Apply(pack(vm, funcName, key), nil, mutable)
			if err != nil {
				t.Error(err)
			}
			t.Logf("%s %s", funcName, res)
			var u64 uint64
			unPack(vm, &u64, funcName, res)
			if u64 != d {
				t.Errorf("unexpected value : want %d, got %d", d, u64)
			} else {
				t.Logf("%s success get %d", funcName, u64)
			}

		}()

		func() {
			vm.ChainContext.IsCreated = false
			funcName := "testReadMapStructE"
			res, err := vm.Apply(pack(vm, funcName, key), nil, mutable)
			if err != nil {
				t.Error(err)
			}
			t.Logf("%s %s", funcName, res)
			var str string
			unPack(vm, &str, funcName, res)
			if str != e {
				t.Errorf("unexpected value : want %s, got %s", e, str)
			} else {
				t.Logf("%s success get %s", funcName, str)
			}

		}()

		func() {
			vm.ChainContext.IsCreated = false
			funcName := "testReadMapStructF"
			res, err := vm.Apply(pack(vm, funcName, key), nil, mutable)
			if err != nil {
				t.Error(err)
			}
			t.Logf("%s %s", funcName, res)
			var addr common.Address
			unPack(vm, &addr, funcName, res)
			if !bytes.Equal(addr.Bytes(), f.Bytes()) {
				t.Errorf("unexpected value : want %s, got %s", f, addr)
			} else {
				t.Logf("%s success get %s", funcName, addr)
			}

		}()

		func() {
			vm.ChainContext.IsCreated = false
			funcName := "testReadMapStructG"
			res, err := vm.Apply(pack(vm, funcName, key), nil, mutable)
			if err != nil {
				t.Error(err)
			}
			t.Logf("%s %s", funcName, res)
			var flag bool
			unPack(vm, &flag, funcName, res)
			if flag != g {
				t.Errorf("unexpected value : want %t, got %t", g, flag)
			} else {
				t.Logf("%s success get %t", funcName, flag)
			}

		}()
	}
	test(key1, a1, b1, c1, d1, e1, f1, g1)
	test(key2, a2, b2, c2, d2, e2, f2, g2)
}
