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

package wavm

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"regexp"

	"github.com/vntchain/go-vnt/accounts/abi"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/common/math"
	mat "github.com/vntchain/go-vnt/common/math"
	"github.com/vntchain/go-vnt/core/vm"
	"github.com/vntchain/go-vnt/core/wavm/gas"
	"github.com/vntchain/go-vnt/core/wavm/utils"
	"github.com/vntchain/go-vnt/log"
	"github.com/vntchain/vnt-wasm/exec"
	"github.com/vntchain/vnt-wasm/validate"
	"github.com/vntchain/vnt-wasm/vnt"
	"github.com/vntchain/vnt-wasm/wasm"
)

const maximum_linear_memory = 33 * 1024 * 1024 //bytes
const maximum_mutable_globals = 1024           //bytes
const maximum_table_elements = 1024            //elements
const maximum_linear_memory_init = 64 * 1024   //bytes
const maximum_func_local_bytes = 8192          //bytes
const wasm_page_size = 64 * 1024

const kPageSize = 64 * 1024
const AddressLength = 20

const FallBackFunctionName = "Fallback"
const FallBackPayableFunctionName = "$Fallback"

type InvalidFunctionNameError string

func (e InvalidFunctionNameError) Error() string {
	return fmt.Sprintf("Exec wasm error: Invalid function name \"%s\"", string(e))
}

type InvalidPayableFunctionError string

func (e InvalidPayableFunctionError) Error() string {
	return fmt.Sprintf("Exec wasm error: Invalid payable function: %s", string(e))
}

type IllegalInputError string

func (e IllegalInputError) Error() string {
	return fmt.Sprintf("Exec wasm error: Illegal input")
}

type UnknownABITypeError string

func (e UnknownABITypeError) Error() string {
	return fmt.Sprintf("Exec wasm error: Unknown abi type \"%s\"", e)
}

type UnknownTypeError string

func (e UnknownTypeError) Error() string {
	return fmt.Sprintf("Exec wasm error: Unknown type")
}

type NoFunctionError string

func (e NoFunctionError) Error() string {
	return fmt.Sprintf("Exec wasm error: Can't find function %s in abi", e)
}

type MismatchMutableFunctionError struct {
	parent  int
	current int
}

func (e MismatchMutableFunctionError) Error() string {
	parentStr := "unmutable"
	if e.parent == 1 {
		parentStr = "mutable"
	}
	currentStr := "unmutable"
	if e.current == 1 {
		currentStr = "mutable"
	}
	return fmt.Sprintf("Mismatch mutable type, parent function type : %s, current function type : %s", parentStr, currentStr)
}

type Wavm struct {
	VM              *exec.Interpreter
	Module          *wasm.Module
	ChainContext    ChainContext
	GasRules        gas.Gas
	WavmConfig      Config
	IsCreated       bool
	currentFuncName string
	MutableList     Mutable
	tempGasLeft     uint64
}

// type InstanceContext struct {
// 	memory *MemoryInstance
// }

func NewWavm(chainctx ChainContext, wavmConfig Config, iscreated bool) *Wavm {
	return &Wavm{
		ChainContext: chainctx,
		WavmConfig:   wavmConfig,
		IsCreated:    iscreated,
	}
}

func (wavm *Wavm) ResolveImports(name string) (*wasm.Module, error) {
	envModule := EnvModule{}
	envModule.InitModule(&wavm.ChainContext)
	return envModule.GetModule(), nil
}

func (wavm *Wavm) captureOp(pc uint64, op byte) error {
	if wavm.WavmConfig.Debug {
		wavm.Tracer().CaptureState(wavm.ChainContext.Wavm, pc, OpCode{Op: op}, wavm.ChainContext.Contract.Gas, 0, nil, nil, wavm.ChainContext.Contract, wavm.ChainContext.Wavm.depth, nil)
	}
	return nil
}

func (wavm *Wavm) captureEnvFunctionStart(pc uint64, funcName string) error {
	wavm.tempGasLeft = wavm.ChainContext.Contract.Gas
	return nil
}

func (wavm *Wavm) captureEnvFunctionEnd(pc uint64, funcName string) error {
	if wavm.WavmConfig.Debug {
		gas := wavm.tempGasLeft - wavm.ChainContext.Contract.Gas
		wavm.Tracer().CaptureState(wavm.ChainContext.Wavm, pc, OpCode{FuncName: funcName}, wavm.ChainContext.Contract.Gas, gas, nil, nil, wavm.ChainContext.Contract, wavm.ChainContext.Wavm.depth, nil)
	}
	return nil
}

func (wavm *Wavm) captrueFault(pc uint64, err error) error {
	if wavm.WavmConfig.Debug {
		wavm.Tracer().CaptureState(wavm.ChainContext.Wavm, pc, OpCode{FuncName: "error"}, wavm.ChainContext.Contract.Gas, 0, nil, nil, wavm.ChainContext.Contract, wavm.ChainContext.Wavm.depth, err)
	}
	return nil
}

func (wavm *Wavm) Tracer() vm.Tracer {
	return wavm.ChainContext.Wavm.wavmConfig.Tracer
}

func instantiateMemory(m *vnt.WavmMemory, module *wasm.Module) error {
	if module.Data != nil {
		var index int
		for _, v := range module.Data.Entries {
			expr, _ := module.ExecInitExpr(v.Offset)
			offset, ok := expr.(int32)
			if !ok {
				return wasm.InvalidValueTypeInitExprError{Wanted: reflect.Int32, Got: reflect.TypeOf(offset).Kind()}
			}
			index = int(offset) + len(v.Data)
			if bytes.Contains(v.Data, []byte{byte(0)}) {
				split := bytes.Split(v.Data, []byte{byte(0)})
				var tmpoffset = int(offset)
				for _, tmp := range split {
					tmplen := len(tmp)
					b, res := utils.IsAddress(tmp)
					if b == true {
						tmp = common.HexToAddress(string(res)).Bytes()
					}
					b, res = utils.IsU256(tmp)
					if b == true {

						bigint := utils.GetU256(res)
						tmp = []byte(bigint.String())
					}
					m.Set(uint64(tmpoffset), uint64(len(tmp)), tmp)
					tmpoffset += tmplen + 1
				}
			} else {
				m.Set(uint64(offset), uint64(len(v.Data)), v.Data)
			}
		}
		m.Pos = index
	} else {
		m.Pos = 0
	}
	return nil
}

func (wavm *Wavm) InstantiateModule(code []byte, memory []uint8) error {
	wasm.SetDebugMode(false)
	buf := bytes.NewReader(code)
	m, err := wasm.ReadModule(buf, wavm.ResolveImports)
	if err != nil {
		log.Error("could not read module", "err", err)
		return err
	}
	if wavm.IsCreated == true {
		err = validate.VerifyModule(m)
		if err != nil {
			log.Error("could not verify module", "err", err)
			return err
		}
	}
	if m.Export == nil {
		log.Error("module has no export section", "export", "nil")
		return errors.New("module has no export section")
	}
	wavm.Module = m
	// m.PrintDetails()
	return nil
}

func (wavm *Wavm) Apply(input []byte, compiled []vnt.Compiled, mutable Mutable) (res []byte, err error) {
	// Catch all the panic and transform it into an error
	defer func() {
		if r := recover(); r != nil {
			log.Error("Got error during wasm execution.", "err", r)
			res = nil
			err = fmt.Errorf("%s", r)
			if wavm.WavmConfig.Debug == true {
				if wavm.VM == nil {
					wavm.captrueFault(uint64(0), err)
				} else {
					wavm.captrueFault(uint64(wavm.VM.Pc()), err)
				}
			}
		}
	}()
	wavm.MutableList = mutable

	//initialize the gas cost for initial memory when create contract before create Interpreter
	//todo memory grow内存消耗
	if wavm.ChainContext.IsCreated == true {
		memSize := uint64(1)
		if len(wavm.Module.Memory.Entries) != 0 {
			memSize = uint64(wavm.Module.Memory.Entries[0].Limits.Initial)
		}
		wavm.ChainContext.GasCounter.GasInitialMemory(memSize)
	}

	var vm *exec.Interpreter
	vm, err = exec.NewInterpreter(wavm.Module, compiled, instantiateMemory, wavm.captureOp, wavm.captureEnvFunctionStart, wavm.captureEnvFunctionEnd, wavm.WavmConfig.Debug)
	if err != nil {
		log.Error("Could not create VM: ", "error", err)
		return nil, fmt.Errorf("Could not create VM: %s", err)
	}

	wavm.VM = vm
	// gas := wavm.ChainContext.Contract.Gas
	// adjustedGas := uint64(gas * exec.WasmCostsOpcodesDiv / exec.WasmCostsOpcodesMul)
	// if adjustedGas > math.MaxUint64 {
	// 	return nil, fmt.Errorf("Wasm interpreter cannot run contracts with gas (wasm adjusted) >= 2^64")
	// }
	//
	// vm.Contract.Gas = adjustedGas

	res, err = wavm.ExecCodeWithFuncName(input)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (wavm *Wavm) GetFallBackFunction() (int64, string) {
	index := int64(-1)
	for name, e := range wavm.VM.Module().Export.Entries {
		if name == FallBackFunctionName {
			index = int64(e.Index)
			return index, FallBackFunctionName
		}
		if name == FallBackPayableFunctionName {
			index = int64(e.Index)
			return index, FallBackPayableFunctionName
		}
	}
	return index, ""
}

func (wavm *Wavm) ExecCodeWithFuncName(input []byte) ([]byte, error) {
	wavm.ChainContext.Wavm.depth++
	defer func() { wavm.ChainContext.Wavm.depth-- }()
	index := int64(0)
	matched := false
	funcName := ""
	VM := wavm.VM
	module := VM.Module()
	Abi := wavm.ChainContext.Abi
	if wavm.ChainContext.IsCreated == true {
		val := Abi.Constructor
		for name, e := range module.Export.Entries {
			if name == val.Name {
				index = int64(e.Index)
				funcName = val.Name
				matched = true
			}
		}
	} else {
		//TODO: do optimization on function searching
		if len(input) < 4 {
			matched = false
		} else {
			sig := input[:4]
			input = input[4:]
			for name, e := range module.Export.Entries {
				if val, ok := Abi.Methods[name]; ok {
					res := val.Id()
					if bytes.Equal(sig, res) {
						matched = true
						funcName = name
						index = int64(e.Index)
						break
					}
				}
			}
		}
	}

	if matched == false {
		//查找是否有fallback方法
		index, funcName = wavm.GetFallBackFunction()
		if index == -1 {
			return nil, InvalidFunctionNameError(funcName)
		}
	}

	if wavm.payable(funcName) != true {
		if wavm.ChainContext.Contract.Value().Cmp(new(big.Int).SetUint64(0)) > 0 {
			return nil, InvalidPayableFunctionError(funcName)
		}
	}
	wavm.currentFuncName = funcName
	var method abi.Method
	if wavm.ChainContext.IsCreated == true {
		method = Abi.Constructor
	} else {
		method = Abi.Methods[funcName]
	}
	var args []uint64

	// if funcName == InitFuntionName {
	// 	input = vm.ChainContext.Input
	// }

	for i, v := range method.Inputs {
		if len(input) < 32*(i+1) {
			return nil, IllegalInputError("")
		}
		arg := input[(32 * i):(32 * (i + 1))]
		switch v.Type.T {
		case abi.StringTy: // variable arrays are written at the end of the return bytes
			output := input[:]
			begin, end, err := lengthPrefixPointsTo(i*32, output)
			if err != nil {
				return nil, err
			}
			value := output[begin : begin+end]
			offset := VM.Memory.SetBytes(value)
			VM.AddHeapPointer(uint64(len(value)))
			args = append(args, uint64(offset))
		case abi.IntTy, abi.UintTy:
			a := readInteger(v.Type.Kind, arg)
			val := reflect.ValueOf(a)
			if val.Kind() == reflect.Ptr { //uint256
				u256 := math.U256(a.(*big.Int))
				value := []byte(u256.String())
				// args = append(args, a.(uint64))
				offset := VM.Memory.SetBytes(value)
				VM.AddHeapPointer(uint64(len(value)))
				args = append(args, uint64(offset))
			} else {
				args = append(args, a.(uint64))
			}
		case abi.BoolTy:
			res, err := readBool(arg)
			if err != nil {
				return nil, err
			}
			args = append(args, res)
		case abi.AddressTy:
			addr := common.BytesToAddress(arg)
			idx := VM.Memory.SetBytes(addr.Bytes())
			VM.AddHeapPointer(uint64(len(addr.Bytes())))
			args = append(args, uint64(idx))
		default:
			return nil, UnknownABITypeError(v.Type.String())
		}
	}
	if wavm.ChainContext.IsCreated == true {
		*VM.Mutable = true
	} else if funcName == FallBackFunctionName || funcName == FallBackPayableFunctionName {
		*VM.Mutable = true
	} else {
		if v, ok := wavm.MutableList[uint32(index)]; ok {
			*VM.Mutable = v
		} else {
			*VM.Mutable = false
		}
	}
	if wavm.ChainContext.Wavm.mutable == -1 {
		if *VM.Mutable == true {
			wavm.ChainContext.Wavm.mutable = 1
		} else {
			wavm.ChainContext.Wavm.mutable = 0
		}
	} else {
		if wavm.ChainContext.Wavm.mutable == 0 && *VM.Mutable == true {
			return nil, MismatchMutableFunctionError{0, 1}
		}
	}

	res, err := VM.ExecContractCode(index, args...)
	if err != nil {
		return nil, err
	}

	// vm.GetGasCost()
	funcType := module.GetFunction(int(index)).Sig
	if len(funcType.ReturnTypes) == 0 {
		return nil, nil
	}

	if val, ok := Abi.Methods[funcName]; ok {
		outputs := val.Outputs
		if len(outputs) != 0 {
			output := outputs[0].Type.T
			switch output {
			case abi.StringTy:
				v := VM.Memory.GetPtr(res)
				l, err := packNum(reflect.ValueOf(32))
				if err != nil {
					return nil, err
				}
				s, err := packBytesSlice(v, len(v))
				if err != nil {
					return nil, err
				}
				return append(l, s...), nil
			case abi.UintTy, abi.IntTy:
				if outputs[0].Type.Kind == reflect.Ptr {
					mem := VM.Memory.GetPtr(res)
					bigint := utils.GetU256(mem)
					return abi.U256(bigint), nil
				} else if output == abi.UintTy {
					return abi.U256(new(big.Int).SetUint64(res)), nil
				} else {
					if outputs[0].Type.Size == 32 {
						return abi.U256(big.NewInt(int64(int32(res)))), nil
					} else {
						return abi.U256(big.NewInt(int64(res))), nil
					}
				}
			case abi.BoolTy:
				if res != 0 {
					return mat.PaddedBigBytes(common.Big1, 32), nil
				}
				return mat.PaddedBigBytes(common.Big0, 32), nil
			case abi.AddressTy:
				v := VM.Memory.GetPtr(res)
				return common.LeftPadBytes(v, 32), nil
			default:
				//todo 所有类型处理
				return nil, UnknownTypeError("")
			}
		} else { //无返回类型
			return utils.I32ToBytes(0), nil
		}
	}
	return nil, NoFunctionError(funcName)
}

func (wavm *Wavm) GetFuncName() string {
	return wavm.currentFuncName
}

func (wavm *Wavm) SetFuncName(name string) {
	wavm.currentFuncName = name
}

func (wavm *Wavm) payable(funcName string) bool {
	reg := regexp.MustCompile(`^\$`)
	res := reg.FindAllString(funcName, -1)
	if len(res) != 0 {
		return true
	}
	return false
}
