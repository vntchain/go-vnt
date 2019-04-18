package exec

import (
	"fmt"
	"math"
	"strings"

	"github.com/vntchain/vnt-wasm/disasm"
	"github.com/vntchain/vnt-wasm/exec/internal/compile"
	"github.com/vntchain/vnt-wasm/vnt"
	"github.com/vntchain/vnt-wasm/wasm"
)

type Interpreter struct {
	*VM
	Memory           *vnt.WavmMemory
	heapPointerIndex int64
	Mutable          *bool
}

func NewInterpreter(module *wasm.Module, compiled []vnt.Compiled, initMem func(m *vnt.WavmMemory, module *wasm.Module) error, captureOp func(pc uint64, op byte) error, captureEnvFunctionStart func(pc uint64, name string) error, captureEnvFunctionEnd func(pc uint64, name string) error, debug bool) (*Interpreter, error) {
	var inter Interpreter
	var vm VM
	vm.captureOp = captureOp
	vm.captureEnvFunctionStart = captureEnvFunctionStart
	vm.captureEnvFunctionEnd = captureEnvFunctionEnd
	vm.debug = debug
	inter.Memory = vnt.NewWavmMemory()
	inter.heapPointerIndex = -1
	mut := false
	inter.Mutable = &mut
	if module.Memory != nil && len(module.Memory.Entries) != 0 {
		if len(module.Memory.Entries) > 1 {
			return nil, ErrMultipleLinearMemories
		}
		vm.memory = make([]byte, uint(module.Memory.Entries[0].Limits.Initial)*wasmPageSize)
		copy(vm.memory, module.LinearMemoryIndexSpace[0])
	} else {
		vm.memory = make([]byte, 1*wasmPageSize)
	}

	inter.Memory.Memory = vm.memory
	// init linear memory with module data section
	err := initMem(inter.Memory, module)
	if err != nil {
		return nil, err
	}

	vm.funcs = make([]function, len(module.FunctionIndexSpace))
	vm.globals = make([]uint64, len(module.GlobalIndexSpace))
	vm.newFuncTable()
	vm.module = module

	nNatives := 0
	for i, fn := range module.FunctionIndexSpace {
		// Skip native methods as they need not be
		// disassembled; simply add them at the end
		// of the `funcs` array as is, as specified
		// in the spec. See the "host functions"
		// section of:
		// https://webassembly.github.io/spec/core/exec/modules.html#allocation
		if fn.IsHost() {
			vm.funcs[i] = contractFunction{
				typ:     fn.Host.Type(),
				val:     fn.Host,
				sig:     len(fn.Sig.ParamTypes),
				memory:  inter.Memory,
				mutable: inter.Mutable,
			}
			nNatives++
			continue
		}

		if len(compiled) != 0 {
			//now := time.Now()
			code := compiled[i].Code
			maxDepth := compiled[i].MaxDepth
			table := compiled[i].Table
			totalLocalVars := compiled[i].TotalLocalVars
			//duration := time.Since(now)
			//vm.NoCompileTimeCost += duration.Seconds()

			vm.funcs[i] = compiledFunction{
				code:           code,
				branchTables:   table,
				maxDepth:       maxDepth,
				totalLocalVars: totalLocalVars,
				args:           len(fn.Sig.ParamTypes),
				returns:        len(fn.Sig.ReturnTypes) != 0,
			}
			continue
		}

		disassembly, err := disasm.Disassemble(fn, module)
		if err != nil {
			return nil, err
		}

		totalLocalVars := 0
		totalLocalVars += len(fn.Sig.ParamTypes)
		for _, entry := range fn.Body.Locals {
			totalLocalVars += int(entry.Count)
		}
		code, table := compile.Compile(disassembly.Code)
		vm.funcs[i] = compiledFunction{
			code:           code,
			branchTables:   table,
			maxDepth:       disassembly.MaxDepth,
			totalLocalVars: totalLocalVars,
			args:           len(fn.Sig.ParamTypes),
			returns:        len(fn.Sig.ReturnTypes) != 0,
		}
	}

	for i, global := range module.GlobalIndexSpace {
		val, err := module.ExecInitExpr(global.Init)
		if err != nil {
			return nil, err
		}
		switch v := val.(type) {
		case int32:
			vm.globals[i] = uint64(v)
		case int64:
			vm.globals[i] = uint64(v)
		case float32:
			vm.globals[i] = uint64(math.Float32bits(v))
		case float64:
			vm.globals[i] = uint64(math.Float64bits(v))
		}
	}

	for k, v := range module.Export.Entries {
		if strings.Contains(k, "heap_pointer") || strings.Contains(k, "__heap_base") {
			if v.Kind == wasm.ExternalGlobal {
				inter.heapPointerIndex = int64(v.Index)
				inter.Memory.Pos = int(vm.globals[inter.heapPointerIndex])
			}
		}
	}

	if module.Start != nil {
		_, err := vm.ExecCode(int64(module.Start.Index))
		if err != nil {
			return nil, err
		}
	}

	inter.VM = &vm
	return &inter, nil
}

func (inter *Interpreter) Pc() int64 {
	return inter.ctx.pc
}

// ExecContractCode calls the function with the given index and arguments.
// fnIndex should be a valid index into the function index space of
// the VM's module.
func (vm *VM) ExecContractCode(fnIndex int64, args ...uint64) (ret uint64, err error) {
	// If used as a library, client code should set vm.RecoverPanic to true
	// in order to have an error returned.
	if vm.RecoverPanic {
		defer func() {
			if r := recover(); r != nil {
				switch e := r.(type) {
				case error:
					err = e
				default:
					err = fmt.Errorf("exec: %v", e)
				}
			}
		}()
	}
	if int(fnIndex) > len(vm.funcs) {
		return 0, InvalidFunctionIndexError(fnIndex)
	}
	if len(vm.module.GetFunction(int(fnIndex)).Sig.ParamTypes) != len(args) {
		return 0, ErrInvalidArgumentCount
	}
	compiled, ok := vm.funcs[fnIndex].(compiledFunction)
	if !ok {
		panic(fmt.Sprintf("exec: function at index %d is not a compiled function", fnIndex))
	}
	if len(vm.ctx.stack) < compiled.maxDepth {
		vm.ctx.stack = make([]uint64, 0, compiled.maxDepth)
	}
	vm.ctx.locals = make([]uint64, compiled.totalLocalVars)
	vm.ctx.pc = 0
	vm.ctx.code = compiled.code
	vm.ctx.curFunc = fnIndex

	for i, arg := range args {
		vm.ctx.locals[i] = arg
	}

	res := vm.execCode(compiled)

	return res, nil
}

func (vm *VM) Module() *wasm.Module {
	return vm.module
}

func (inter *Interpreter) AddHeapPointer(size uint64) {
	if inter.heapPointerIndex == -1 {
		//log.Debug("AddHeapPointer", "MSG", "Can't find head_pointer")
	} else {
		res := inter.globals[inter.heapPointerIndex]
		//log.Debug("AddHeapPointer", "MSG", "current heap", res, "size", size)
		res = res + size
		inter.globals[inter.heapPointerIndex] = res
	}
}

func (vm *VM) ResetContext() {
	vm.ctx = context{
		stack:   make([]uint64, 0),
		locals:  make([]uint64, 0),
		code:    make([]byte, 0),
		pc:      0,
		curFunc: 0,
	}
}

// Process is a proxy passed to host functions in order to access
// things such as memory and control.
type WavmProcess struct {
	vm      *VM
	memory  *vnt.WavmMemory
	mutable *bool
}

// NewProcess creates a VM interface object for host functions
func NewWavmProcess(vm *VM, memory *vnt.WavmMemory, mutable *bool) *WavmProcess {
	return &WavmProcess{
		vm,
		memory,
		mutable,
	}
}

func (proc *WavmProcess) ReadAt(off uint64) []byte {
	return proc.memory.GetPtr(off)
}

func (proc *WavmProcess) WriteAt(p []byte, off int64) (int, error) {
	proc.memory.Set(uint64(off), uint64(len(p)), p)
	return 0, nil
}

func (proc *WavmProcess) SetBytes(value []byte) (offset int) {
	return proc.memory.SetBytes(value)
}

func (proc *WavmProcess) GetData() []byte {
	return proc.memory.Data()
}

// Terminate stops the execution of the current module.
func (proc *WavmProcess) Terminate() {
	proc.vm.abort = true
}

func (proc *WavmProcess) Mutable() bool {
	return *proc.mutable
}
