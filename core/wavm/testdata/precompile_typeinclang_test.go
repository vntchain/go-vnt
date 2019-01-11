package tests

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/vntchain/go-vnt/core/wavm"
)

var (
	precodepath2 = filepath.Join(basepath, "precompile/wasm/main5.wasm")
	preabipath2  = filepath.Join(basepath, "precompile/wasm/main5.json")
)

func TestMappinp2(t *testing.T) {
	vm := newWavm(precodepath2, preabipath2, false)
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	for i := 0; i < 1; i++ {
		key := fmt.Sprintf("testkey%d", i)
		val := fmt.Sprintf("testvalue%d", i)
		res, err := vm.Apply(pack(vm, "writemapping", key, val), nil, mutable)
		if err != nil {
			t.Error(err)
		}
		t.Logf("writemapping %s", res)
	}
	for i := 0; i < 1; i++ {
		func() {
			key := fmt.Sprintf("testkey%d", i)
			val := fmt.Sprintf("testvalue%d", i)
			res, err := vm.Apply(pack(vm, "readmapping", key), nil, mutable)
			if err != nil {
				t.Error(err)
			}
			t.Logf("readmapping %s", res)
			var str string
			t.Logf("readmapping %s", res)
			unPack(vm, &str, "readmapping", res)
			if str != val {
				t.Errorf("unexpected value : want %s, got %s", val, str)
			}
			t.Logf("readmapping %s", str)
		}()
	}
}

//[236 17 32 201 140 212 97 170 145 54 19 134 226 187 12 116 30 136 78 127 98 137 99 35 0 16 30 67 139 170 142 48]
//[126 69 3 134 144 202 185 230 248 48 189 225 1 83 79 235 19 150 224 252 157 126 151 222 58 200 246 29 15 214 34 132]
func TestMapping_int32_2(t *testing.T) {
	vm := newWavm(precodepath2, preabipath2, false)
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	for i := 0; i < 100; i++ {
		key := i
		val := i - 100
		res, err := vm.Apply(pack(vm, "writemapping_int32", int32(key), int32(val)), nil, mutable)
		if err != nil {
			t.Error(err)
		}
		t.Logf("writemapping_int32 %s", res)
	}
	for i := 0; i < 100; i++ {
		func() {
			key := i
			val := i - 100
			res, err := vm.Apply(pack(vm, "readmapping_int32", int32(key)), nil, mutable)
			if err != nil {
				t.Error(err)
			}
			t.Logf("readmapping_int32 %s", res)
			var result int32
			t.Logf("readmapping_int32 %s", res)
			unPack(vm, &result, "readmapping_int32", res)
			if result != int32(val) {
				t.Errorf("unexpected value : want %d, got %d", val, result)
			}
			t.Logf("readmapping_int32 %d", result)
		}()
	}
}

func TestMapping_int64_2(t *testing.T) {
	vm := newWavm(precodepath2, preabipath2, false)
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	for i := 0; i < 1; i++ {
		key := i
		val := i - 9223372036854775807
		res, err := vm.Apply(pack(vm, "writemapping_int64", int64(key), int64(val)), nil, mutable)
		if err != nil {
			t.Error(err)
		}
		t.Logf("writemapping_int64 %s", res)
	}
	for i := 0; i < 1; i++ {
		func() {
			key := i
			val := i - 9223372036854775807
			res, err := vm.Apply(pack(vm, "readmapping_int64", int64(key)), nil, mutable)
			if err != nil {
				t.Error(err)
			}
			t.Logf("readmapping_int64 %s", res)
			var result int64
			t.Logf("readmapping_int64 %s", res)
			unPack(vm, &result, "readmapping_int64", res)
			if result != int64(val) {
				t.Errorf("unexpected value : want %d, got %d", val, result)
			}
			t.Logf("readmapping_int64 %d", result)
		}()
	}
}

func TestArray2(t *testing.T) {
	vm := newWavm(precodepath2, preabipath2, false)
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	length := uint64(1)
	for i := 0; i < 1; i++ {
		key := i
		val := fmt.Sprintf("array_%d", i+100)

		res, err := vm.Apply(pack(vm, "writearray", uint64(key), val, uint64(length)), nil, mutable)
		if err != nil {
			t.Error(err)
		}
		t.Logf("writearray %s", res)
	}
	for i := 0; i < 1; i++ {
		func() {
			key := i
			val := fmt.Sprintf("array_%d", i+100)
			res, err := vm.Apply(pack(vm, "readarray", uint64(key)), nil, mutable)
			if err != nil {
				t.Error(err)
			}
			var result string
			t.Logf("readarray %s", res)
			unPack(vm, &result, "readarray", res)
			if result != val {
				t.Errorf("unexpected value : want %s, got %s", val, result)
			}
			t.Logf("readarray %s", result)
		}()
	}

	func() {
		res, err := vm.Apply(pack(vm, "readarraylength", uint64(1)), nil, mutable)
		if err != nil {
			t.Error(err)
		}
		var result uint64
		t.Logf("readarraylength %s", res)
		unPack(vm, &result, "readarraylength", res)
		if result != length {
			t.Errorf("unexpected value : want %d, got %d", length, result)
		}
		t.Logf("readarraylength %d", result)
	}()

	func() {
		res, err := vm.Apply(pack(vm, "TestPop"), nil, mutable)
		if err != nil {
			t.Error(err)
		}
		var result string
		t.Logf("TestPop %s", res)
		unPack(vm, &result, "TestPop", res)
		if result != "100" {
			t.Errorf("unexpected value : want %s, got %s", "100", result)
		}
		t.Logf("TestPop %s", result)
	}()

	// func() {
	// 	res, err := vm.Apply(pack(vm, "readarraylength", uint64(1)))
	// 	if err != nil {
	// 		t.Error(err)
	// 	}
	// 	var result uint64
	// 	t.Logf("readarraylength %s", res)
	// 	unPack(vm, &result, "readarraylength", res)
	// 	if result != length-1 {
	// 		t.Errorf("unexpected value : want %d, got %d", length-1, result)
	// 	}
	// 	t.Logf("readarraylength %d", result)
	// }()
	// pushval := "pushval"
	// func() {
	// 	res, err := vm.Apply(pack(vm, "TestPush", pushval))
	// 	if err != nil {
	// 		t.Error(err)
	// 	}
	// 	t.Logf("TestPop %s", res)
	// }()

	// func() {
	// 	res, err := vm.Apply(pack(vm, "readarraylength", uint64(1)))
	// 	if err != nil {
	// 		t.Error(err)
	// 	}
	// 	var result uint64
	// 	t.Logf("readarraylength %s", res)
	// 	unPack(vm, &result, "readarraylength", res)
	// 	if result != length {
	// 		t.Errorf("unexpected value : want %d, got %d", length, result)
	// 	}
	// 	t.Logf("readarraylength %d", result)
	// }()

	// func() {
	// 	res, err := vm.Apply(pack(vm, "readarray", uint64(0)))
	// 	if err != nil {
	// 		t.Error(err)
	// 	}
	// 	var result string
	// 	t.Logf("readarray %s", res)
	// 	unPack(vm, &result, "readarray", res)
	// 	if result != pushval {
	// 		t.Errorf("unexpected value : want %s, got %s", pushval, result)
	// 	}
	// 	t.Logf("readarray %s", result)
	// }()
}
