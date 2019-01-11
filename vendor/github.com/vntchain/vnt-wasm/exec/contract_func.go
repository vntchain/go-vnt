// Copyright 2017 The go-interpreter Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exec

import (
	"fmt"
	"reflect"

	"github.com/vntchain/vnt-wasm/vnt"
)

type contractFunction struct {
	val     reflect.Value
	typ     reflect.Type
	sig     int
	memory  *vnt.WavmMemory
	mutable *bool
}

func (fn contractFunction) call(vm *VM, index int64) {
	// numIn := fn.typ.NumIn()
	args := make([]reflect.Value, fn.sig+1)
	proc := NewWavmProcess(vm, fn.memory, fn.mutable)

	// Pass proc as an argument. Check that the function indeed
	// expects a *Process argument.
	if reflect.ValueOf(proc).Kind() != fn.typ.In(0).Kind() {
		panic(fmt.Sprintf("exec: the first argument of a host function was %s, expected %s", fn.typ.In(0).Kind(), reflect.ValueOf(vm).Kind()))
	}
	args[0] = reflect.ValueOf(proc)
	for i := fn.sig; i >= 1; i-- {
		raw := vm.popUint64()
		args[i] = reflect.ValueOf(raw)
	}
	rtrns := fn.val.Call(args)
	if rtrns == nil {
		return
	}
	for i, out := range rtrns {
		kind := out.Kind()
		switch kind {
		case reflect.Float64, reflect.Float32:
			vm.pushFloat64(out.Float())
		case reflect.Uint32, reflect.Uint64:
			vm.pushUint64(out.Uint())
		case reflect.Int32, reflect.Int64:
			vm.pushInt64(out.Int())
		default:
			panic(fmt.Sprintf("exec: return value %d invalid kind=%v", i, kind))
		}
	}
}
