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
	return fmt.Sprintf("Invalid function name: %s", string(e))
}

type InvalidPayableFunctionError string

func (e InvalidPayableFunctionError) Error() string {
	return fmt.Sprintf("Invalid payable function: %s", string(e))
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
	return fmt.Sprintf("Mismatch mutable type , parent function type : %s , current function type : %s", parentStr, currentStr)
}

type Wavm struct {
	VM              *exec.Interpreter
	Module          *wasm.Module
	ChainContext    ChainContext
	GasRules        gas.Gas
	VmConfig        vm.Config
	IsCreated       bool
	currentFuncName string
	MutableList     Mutable
}

// type InstanceContext struct {
// 	memory *MemoryInstance
// }

type ActionName string

const (
	ActionNameInit  = "init"
	ActionNameApply = "deploy"
	ActionNameQuery = "query"
)

func NewWavm(chainctx ChainContext, vmconfig vm.Config, iscreated bool) *Wavm {
	return &Wavm{
		ChainContext: chainctx,
		VmConfig:     vmconfig,
		IsCreated:    iscreated,
	}
}

func (wavm *Wavm) ResolveImports(name string) (*wasm.Module, error) {
	envModule := EnvModule{}
	envModule.InitModule(&wavm.ChainContext)
	return envModule.GetModule(), nil
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

	//create需要验证，call不需要
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
	log.Debug("Wavm", "func", "apply")
	defer func() {
		if r := recover(); r != nil {
			log.Error("Got error during wasm execution.", "err", r)
			res = nil
			err = fmt.Errorf("%s", r)
		}
	}()
	wavm.MutableList = mutable
	vm, err := exec.NewInterpreter(wavm.Module, compiled, instantiateMemory)
	if err != nil {
		log.Error("could not create VM: ", "error", err)
		return nil, err
	}

	//initialize the gas cost for initial memory when create contract
	//todo memory grow内存消耗
	if wavm.ChainContext.IsCreated == true {
		memSize := uint64(1)
		if len(wavm.Module.Memory.Entries) != 0 {
			memSize = uint64(wavm.Module.Memory.Entries[0].Limits.Initial)
		}
		wavm.ChainContext.GasCounter.GasInitialMemory(memSize)
	}

	wavm.VM = vm
	// gas := wavm.ChainContext.Contract.Gas
	// adjustedGas := uint64(gas * exec.WasmCostsOpcodesDiv / exec.WasmCostsOpcodesMul)
	// if adjustedGas > math.MaxUint64 {
	// 	return nil, fmt.Errorf("Wasm interpreter cannot run contracts with gas (wasm adjusted) >= 2^64")
	// }
	//
	// vm.Contract.Gas = adjustedGas
	log.Debug("GAS", "NORMAL", wavm.ChainContext.Contract.Gas)

	res, err = wavm.ExecCodeWithFuncName(input)
	if err != nil {
		log.Error("wavm", "call", err)
		return nil, err
	}

	log.Debug("GAS", "GASLEFT", wavm.ChainContext.Contract.Gas)
	log.Debug("======wavm======", "====res====", res)
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
	log.Debug("VM", "func", ">>>ExecCodeWithFuncName", "input", input)
	log.Debug("VM", "func", ">>>ExecCodeWithFuncName", "GasLimit", wavm.ChainContext.Contract.GasLimit, "Gas", wavm.ChainContext.Contract.Gas)
	wavm.ChainContext.Wavm.depth++
	defer func() { wavm.ChainContext.Wavm.depth-- }()
	index := int64(0)
	//foo(string,string)
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
			// //查找是否有fallback方法
			// index, funcName = wavm.GetFallBackFunction()
			// if index == -1 {
			// 	return nil, fmt.Errorf("%s", "Illegal input")
			// }
			// // funcName = FallBackFunctionName
		} else {
			sig := input[:4]
			input = input[4:]
			for name, e := range module.Export.Entries {
				log.Debug("vm", "func", "ExecCodeWithFuncName", "sig", sig, "name", name)
				if val, ok := Abi.Methods[name]; ok {
					res := val.Id()
					log.Debug("vm", "func", "ExecCodeWithFuncName", "sig", sig, "res", res)
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
		// funcName = FallBackFunctionName
	}

	if wavm.payable(funcName) != true {
		if wavm.ChainContext.Contract.Value().Cmp(new(big.Int).SetUint64(0)) > 0 {
			return nil, InvalidPayableFunctionError(funcName)
		}
	}

	log.Debug("vm", "func", "ExecCodeWithFuncName", "funcName", funcName, "funcIndex", index)
	wavm.currentFuncName = funcName
	var method abi.Method
	if wavm.ChainContext.IsCreated == true {
		method = Abi.Constructor
	} else {
		method = Abi.Methods[funcName]
	}
	log.Debug("vm", "funcName", funcName)
	var args []uint64

	// if funcName == InitFuntionName {
	// 	input = vm.ChainContext.Input
	// }

	log.Debug("vm", "func", "Inputs", "len", len(method.Inputs), "input", input)
	for i, v := range method.Inputs {
		if len(input) < 32*(i+1) {
			return nil, fmt.Errorf("%s", "Illegal input")
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
			// log.Debug("vm", "func", "ExecCodeWithFuncName", "address", addr.Hex())
			// log.Debug("vm", "func", "ExecCodeWithFuncName", "address", addr.Bytes())
			idx := VM.Memory.SetBytes(addr.Bytes())
			VM.AddHeapPointer(uint64(len(addr.Bytes())))
			args = append(args, uint64(idx))
		default:
			return nil, fmt.Errorf("abi: unknown type %v", v.Type.T)
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
		log.Debug("Methods", "funcname", funcName, "value", val)
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
				} else {
					return abi.U256(new(big.Int).SetUint64(res)), nil
				}

			case abi.BoolTy:
				if res == 1 {
					return mat.PaddedBigBytes(common.Big1, 32), nil
				}
				return mat.PaddedBigBytes(common.Big0, 32), nil
			case abi.AddressTy:
				v := VM.Memory.GetPtr(res)
				return common.LeftPadBytes(v, 32), nil
			default:
				//todo 所有类型处理
				return nil, errors.New("unknown type")
			}
		} else { //无返回类型
			return utils.I32ToBytes(0), nil
		}
	}
	return nil, fmt.Errorf("can't find entrypoint %s in abi", funcName)
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
