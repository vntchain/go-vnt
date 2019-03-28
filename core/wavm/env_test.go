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
	"testing"

	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"

	"github.com/stretchr/testify/assert"
	"github.com/vntchain/go-vnt/accounts/abi"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core/state"
	"github.com/vntchain/go-vnt/core/vm"
	"github.com/vntchain/go-vnt/core/wavm/contract"
	g "github.com/vntchain/go-vnt/core/wavm/gas"
	"github.com/vntchain/go-vnt/log"
	"github.com/vntchain/go-vnt/params"
	"github.com/vntchain/go-vnt/trie"
	"github.com/vntchain/go-vnt/vntdb"
	"github.com/vntchain/vnt-wasm/exec"
	"github.com/vntchain/vnt-wasm/validate"
	"github.com/vntchain/vnt-wasm/wasm"
)

var debugCodePath = filepath.Join("testdata/debug", "program.wasm")
var debugAbiPath = filepath.Join("testdata/debug", "abi.json")

var eventCodePath = filepath.Join("testdata/event", "program.wasm")
var eventAbiPath = filepath.Join("testdata/event", "abi.json")

var logMsg = ""
var logCtx []interface{} = nil
var logHandler = log.FuncHandler(func(r *log.Record) error {
	logMsg = r.Msg
	logCtx = r.Ctx
	return nil
})

var tHash common.Hash
var bHash common.Hash

func clearLog() {
	logMsg = ""
	logCtx = nil
}

func init() {
	//log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
	log.Root().SetHandler(logHandler)
}

type fakeDB struct {
}

func (db *fakeDB) OpenTrie(root common.Hash) (state.Trie, error) {
	return nil, nil
}

func (db *fakeDB) OpenStorageTrie(addrHash, root common.Hash) (state.Trie, error) {
	return nil, nil
}

// CopyTrie returns an independent copy of the given trie.
func (db *fakeDB) CopyTrie(state.Trie) state.Trie {
	return nil
}

// ContractCode retrieves a particular contract's code.
func (db *fakeDB) ContractCode(addrHash, codeHash common.Hash) ([]byte, error) {
	return nil, nil
}

// ContractCodeSize retrieves a particular contracts code's size.
func (db *fakeDB) ContractCodeSize(addrHash, codeHash common.Hash) (int, error) {
	return 0, nil
}

// TrieDB retrieves the low level trie database used for data storage.
func (db *fakeDB) TrieDB() *trie.Database {
	return nil
}

func readAbi(abiPath string) abi.ABI {
	abiData, err := ioutil.ReadFile(abiPath)
	if err != nil {
		log.Crit("could not read abi: ", "error", err)
	}

	abi, err := GetAbi(abiData)
	if err != nil {
		log.Crit("could not read abi: ", "error", err)
	}

	return abi
}

func importer(name string) (*wasm.Module, error) {
	f, err := os.Open(name + ".wasm")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m, err := wasm.ReadModule(f, nil)
	if err != nil {
		return nil, err
	}
	err = validate.VerifyModule(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func getVM(codeFile string, abiPath string) (*exec.Interpreter, EnvFunctions) {
	var fileHandler, err = os.Open(codeFile)
	defer fileHandler.Close()
	abi := readAbi(abiPath)
	addr := common.BytesToAddress([]byte("0xd2be7e0d40c1a73ec1709f00b11cb5e24c784077"))

	chainconfig := &params.ChainConfig{HubbleBlock: big.NewInt(0)}
	gasRule := g.NewGas(false)
	gasTable := chainconfig.GasTable(new(big.Int).SetInt64(10000))
	contract := contract.NewWASMContract(vm.AccountRef(addr),
		vm.AccountRef(addr), big.NewInt(100), 200000)
	gasCounter := g.NewGasCounter(contract, gasTable)
	cc := ChainContext{
		BlockNumber: big.NewInt(1),
		Contract:    contract,
		Abi:         abi,
		GasRule:     gasRule,
		GasCounter:  gasCounter,
		GasLimit:    10000000,
		Wavm: &WAVM{
			wavmConfig: Config{Debug: true, Tracer: NewWasmLogger(nil)},
			Wavm:       &Wavm{},
		},
	}

	cc.StateDB = prepareState()
	cc.StateDB.GetOrNewStateObject(cc.Contract.Address())

	envModule := EnvModule{}
	envModule.InitModule(&cc)

	resolveHost := func(name string) (*wasm.Module, error) {
		return envModule.GetModule(), nil
	}

	m, err := wasm.ReadModule(fileHandler, resolveHost)

	fmt.Println(err)
	if err != nil {
		log.Crit("could not read module: ", "error", err)
	}

	//compiled, err := CompileModule(m, cc)
	//compiled := make([]vnt.Compiled, 0)

	vm, err := exec.NewInterpreter(m, nil, instantiateMemory, cc.Wavm.Wavm.captureOp, cc.Wavm.Wavm.captureEnvFunction, false)
	if err != nil {
		log.Crit("failed to create vm: ", "error", err)
	}
	vm.ResetContext()

	return vm, envModule.GetEnvFunctions()
}

func handlePanic(t *testing.T, msg string) {
	if r := recover(); r != nil {
		assert.Equal(t, msg, r)
	} else {
		assert.Equal(t, 1, 2)
	}
}

func TestVM_PrintAddress(t *testing.T) {
	log.Root().SetHandler(logHandler)
	defer clearLog()

	vm, ef := getVM(debugCodePath, debugAbiPath)

	ef.ctx.Wavm.Wavm.SetFuncName("init")

	var mutable = true

	proc := exec.NewWavmProcess(vm.VM, vm.Memory, &mutable)

	fmt.Println("testing...")

	remarkIdx := uint64(vm.Memory.SetBytes([]byte("The sender is: ")))
	strIdx := uint64(vm.Memory.SetBytes(common.HexToAddress("0x880d84da2bE4D02830b03FF4CF0840924Be6B0A6").Bytes()))

	ef.PrintAddress(proc, remarkIdx, strIdx)

	expectedMsg := "Contract Debug >>>>"

	assert.Equal(t, expectedMsg, logMsg)
	assert.Equal(t, 4, len(logCtx))
	assert.Equal(t, "init", logCtx[1])
	assert.Equal(t, "The sender is: 0x880d84da2bE4D02830b03FF4CF0840924Be6B0A6", logCtx[3])
}

func TestVM_PrintStr(t *testing.T) {
	log.Root().SetHandler(logHandler)
	defer clearLog()

	vm, ef := getVM(debugCodePath, debugAbiPath)

	ef.ctx.Wavm.Wavm.SetFuncName("init")
	var mutable = true

	proc := exec.NewWavmProcess(vm.VM, vm.Memory, &mutable)

	fmt.Println("testing...")

	remarkIdx := uint64(vm.Memory.SetBytes([]byte("The contract name is: ")))
	strIdx := uint64(vm.Memory.SetBytes([]byte("VNT Token")))

	ef.PrintStr(proc, remarkIdx, strIdx)

	expectedMsg := "Contract Debug >>>>"

	assert.Equal(t, expectedMsg, logMsg)
	assert.Equal(t, 4, len(logCtx))
	assert.Equal(t, "init", logCtx[1])
	assert.Equal(t, "The contract name is: VNT Token", logCtx[3])
}

func TestVM_PrintInt64T(t *testing.T) {
	log.Root().SetHandler(logHandler)
	defer clearLog()
	vm, ef := getVM(debugCodePath, debugAbiPath)
	var mutable = true

	proc := exec.NewWavmProcess(vm.VM, vm.Memory, &mutable)

	fmt.Println("testing...")

	remarkIdx := uint64(vm.Memory.SetBytes([]byte("The value is: ")))
	intValue := uint64(0x0f)

	expectedMsg := "Contract Debug >>>>"

	ef.PrintUint64T(proc, remarkIdx, intValue)

	assert.Equal(t, expectedMsg, logMsg)
	assert.Equal(t, 4, len(logCtx))
	assert.Equal(t, "The value is: 15", logCtx[3])
}

func TestVM_PrintInt32T(t *testing.T) {
	log.Root().SetHandler(logHandler)
	defer clearLog()
	vm, ef := getVM(debugCodePath, debugAbiPath)

	var mutable = true

	proc := exec.NewWavmProcess(vm.VM, vm.Memory, &mutable)

	fmt.Println("testing...")

	remarkIdx := uint64(vm.Memory.SetBytes([]byte("The value is: ")))
	intValue := uint64(0x01)

	expectedMsg := "Contract Debug >>>>"

	ef.PrintUint64T(proc, remarkIdx, intValue)

	assert.Equal(t, expectedMsg, logMsg)
	assert.Equal(t, 4, len(logCtx))
	assert.Equal(t, "The value is: 1", logCtx[3])
}

func TestVM_PrintUint64T(t *testing.T) {
	log.Root().SetHandler(logHandler)
	defer clearLog()
	vm, ef := getVM(debugCodePath, debugAbiPath)

	var mutable = true

	proc := exec.NewWavmProcess(vm.VM, vm.Memory, &mutable)

	fmt.Println("testing...")

	remarkIdx := uint64(vm.Memory.SetBytes([]byte("The value is: ")))
	intValue := uint64(0x01)

	expectedMsg := "Contract Debug >>>>"

	ef.PrintUint64T(proc, remarkIdx, intValue)

	assert.Equal(t, expectedMsg, logMsg)
	assert.Equal(t, 4, len(logCtx))
	assert.Equal(t, "The value is: 1", logCtx[3])
}

func TestVM_PrintUint32T(t *testing.T) {
	log.Root().SetHandler(logHandler)
	defer clearLog()
	vm, ef := getVM(debugCodePath, debugAbiPath)

	var mutable = true

	proc := exec.NewWavmProcess(vm.VM, vm.Memory, &mutable)

	fmt.Println("testing...")

	remarkIdx := uint64(vm.Memory.SetBytes([]byte("The value is: ")))
	intValue := uint64(0x01)

	expectedMsg := "Contract Debug >>>>"

	ef.PrintUint64T(proc, remarkIdx, intValue)

	assert.Equal(t, expectedMsg, logMsg)
	assert.Equal(t, 4, len(logCtx))
	assert.Equal(t, "The value is: 1", logCtx[3])
}

func TestVM_getPrintRemark(t *testing.T) {
	defer clearLog()
	vm, ef := getVM(debugCodePath, debugAbiPath)

	var mutable = true

	proc := exec.NewWavmProcess(vm.VM, vm.Memory, &mutable)

	fmt.Println("testing...")

	remarkIdx := uint64(vm.Memory.SetBytes([]byte("The value is: ")))

	remark := ef.getPrintRemark(proc, remarkIdx)

	assert.Equal(t, "The value is: ", remark)
}

//
//func TestVM_Event_One_Topic_No_Data(t *testing.T) {
//	defer clearLog()
//
//	vm, _ := getVM(eventCodePath, eventAbiPath)
//	//proc := exec.NewWavmProcess(vm.VM, vm.Memory)

//fmt.Println("testing...")
//
//	// make a fake db
//	db := &fakeDB{}
//
//	// make a fake statedb object
//	value := make([]byte, 1)
//	value[0] = 0x01
//	hash := common.BytesToHash(value)
//	stateDB, _ := state.New(hash, db)
//	tHash := common.BytesToHash(value)
//	bHash := common.BytesToHash(value)
//	stateDB.Prepare(tHash, bHash, 1)
//	vm.StateDB = stateDB
//
//	// construct the memorydata
//	evtName := "INIT"
//	addr := "0xd2be7e0d40c1a73ec1709f00b11cb5e24c784077"
//	address := common.HexToAddress(addr)
//
//	// construct the locals
//	locals := vm.GetContext().GetLocals()
//	locals = append(locals, uint64(vm.Memory.SetBytes([]byte(evtName))))
//	locals = append(locals, uint64(vm.Memory.SetBytes(address.Bytes())))
//	vm.GetContext().SetLocals(locals)
//
//	// Execute the target method
//	vm.Event()
//
//	// check the result
//	addedLog := stateDB.GetLogs(tHash)
//	assert.Equal(t, 1, len(addedLog))
//	assert.Equal(t, tHash, addedLog[0].TxHash)
//	assert.Equal(t, bHash, addedLog[0].BlockHash)
//	assert.Equal(t, uint(1), addedLog[0].TxIndex)
//	assert.Equal(t, uint(0), addedLog[0].Index)
//	assert.Equal(t, 0, len(addedLog[0].Data))
//	assert.Equal(t, 2, len(addedLog[0].Topics))
//	assert.Equal(t, vm.Abi.Events[evtName].Id(), addedLog[0].Topics[0])
//	assert.Equal(t, common.BytesToHash(address.Bytes()), addedLog[0].Topics[1])
//}
//
//func TestVM_Event_Two_Topic_One_Data(t *testing.T) {
//	defer clearLog()
//	vm, ef := getVM(eventCodePath, eventAbiPath)
//
//	// construct the memorydata
//	evtName := "TESTEVENT"
//	addrFrom := "0xd2be7e0d40c1a73ec1709f00b11cb5e24c784077"
//	addrTo := "0xd2be7e0d40c1a73ec1709f00b11cb5e24c784078"
//	addressFrom := common.HexToAddress(addrFrom)
//	addressTo := common.HexToAddress(addrTo)
//
//	var mutable = true
//
//	proc := exec.NewWavmProcess(vm.VM, vm.Memory, &mutable)
//
//	fmt.Println("testing...")
//
//	// construct the locals
//	var1 := uint64(proc.SetBytes(addressFrom.Bytes()))
//	var2 := uint64(proc.SetBytes(addressTo.Bytes()))
//	var3 := uint64(1000)
//	var4 := uint64(proc.SetBytes([]byte("this is a balance transfer")))
//
//	fn := ef.GetFuncTable()[evtName].Host
//	args := make([]reflect.Value, 5)
//	args[0] = reflect.ValueOf(proc)
//	//args[1] = reflect.ValueOf(var1)
//	//args[2] = reflect.ValueOf(var2)
//	//args[3] = reflect.ValueOf(var3)
//	//args[4] = reflect.ValueOf(var4)
//
//	//arg1 := reflect.New(fn.Type().In(0)).Elem()
//	//arg2 := reflect.New(fn.Type().In(1)).Elem()
//	//arg3 := reflect.New(fn.Type().In(2)).Elem()
//	//arg4 := reflect.New(fn.Type().In(3)).Elem()
//	//
//	//arg1.SetUint(var1)
//	//arg2.SetUint(var2)
//	//arg3.SetUint(var3)
//	//arg4.SetUint(var4)
//
//	args[1] = reflect.ValueOf(var1)
//	args[2] = reflect.ValueOf(var2)
//	args[3] = reflect.ValueOf(var3)
//	args[4] = reflect.ValueOf(var4)
//
//	valueByte := make([]byte, 8)
//	binary.BigEndian.PutUint64(valueByte, uint64(1000))
//
//	// Execute the target method
//	fn.Call(args)
//
//	// check the result
//	addedLog := ef.ctx.StateDB.GetLogs(tHash)
//	assert.Equal(t, 1, len(addedLog))
//	assert.Equal(t, tHash, addedLog[0].TxHash)
//	assert.Equal(t, bHash, addedLog[0].BlockHash)
//	assert.Equal(t, uint(1), addedLog[0].TxIndex)
//	assert.Equal(t, uint(0), addedLog[0].Index)
//	assert.Equal(t, 128, len(addedLog[0].Data))
//	assert.Equal(t, valueByte, addedLog[0].Data[24:32])
//	assert.Equal(t, 3, len(addedLog[0].Topics))
//	assert.Equal(t, ef.ctx.Abi.Events[evtName].Id(), addedLog[0].Topics[0])
//	assert.Equal(t, common.BytesToHash(addressFrom.Bytes()), addedLog[0].Topics[1])
//	assert.Equal(t, common.BytesToHash(addressTo.Bytes()), addedLog[0].Topics[2])
//}

func prepareState() *state.StateDB {
	// make a fake db
	//db := &fakeDB{}
	//
	//// make a fake statedb object
	//value := make([]byte, 1)
	//value[0] = 0x01
	//hash := common.BytesToHash(value)
	//stateDB, _ := state.New(hash, db)
	//tHash := common.BytesToHash(value)
	//bHash := common.BytesToHash(value)
	//stateDB.Prepare(tHash, bHash, 1)

	db := vntdb.NewMemDatabase()
	value := make([]byte, 1)
	value[0] = 0x01
	state, _ := state.New(common.Hash{}, state.NewDatabase(db))
	tHash = common.BytesToHash(value)
	bHash = common.BytesToHash(value)
	state.Prepare(tHash, bHash, 1)

	return state
}

func TestVM_Address(t *testing.T) {
	//log.Root().SetHandler(logHandler)
	//defer clearLog()
	vm, ef := getVM(eventCodePath, eventAbiPath)
	var mutable = true

	proc := exec.NewWavmProcess(vm.VM, vm.Memory, &mutable)

	fmt.Println("testing...")

	address := string("0x0523029b179009a28a7fae478cd0c2e5ba2adc38")

	offset := vm.Memory.SetBytes([]byte(address))

	tOffset := ef.AddressFrom(proc, uint64(offset))

	strBytes := vm.Memory.GetPtr(tOffset)
	assert.Equal(t, common.HexToAddress(address).String(), common.BytesToAddress(strBytes).String())
}

func TestVM_AddressNoPrefix(t *testing.T) {
	//log.Root().SetHandler(logHandler)
	//defer clearLog()
	vm, ef := getVM(eventCodePath, eventAbiPath)

	mutable := true

	proc := exec.NewWavmProcess(vm.VM, vm.Memory, &mutable)

	fmt.Println("testing...")

	address := string("0523029b179009a28a7fae478cd0c2e5ba2adc38")

	offset := vm.Memory.SetBytes([]byte(address))

	tOffset := ef.AddressFrom(proc, uint64(offset))

	strBytes := vm.Memory.GetPtr(tOffset)

	assert.Equal(t, common.HexToAddress(address).String(), common.BytesToAddress(strBytes).String())
}

func TestVM_AddressTooShort(t *testing.T) {
	//log.Root().SetHandler(logHandler)
	//defer clearLog()
	//defer handlePanic(t, "wrong format of address string literal '0x0523029b179009a28a7fae478cd0c2e5ba2adc' with length '40'")
	vm, ef := getVM(eventCodePath, eventAbiPath)

	mutable := true

	proc := exec.NewWavmProcess(vm.VM, vm.Memory, &mutable)

	fmt.Println("testing...")

	address := string("0x0523029b179009a28a7fae478cd0c2e5ba2adc")

	offset := vm.Memory.SetBytes([]byte(address))

	dest := string("0x00523029b179009a28a7fae478cd0c2e5ba2adc")

	tOffset := ef.AddressFrom(proc, uint64(offset))

	strBytes := vm.Memory.GetPtr(tOffset)

	assert.Equal(t, common.HexToAddress(dest).String(), common.BytesToAddress(strBytes).String())

}

func TestVM_AddressTooLong(t *testing.T) {
	//log.Root().SetHandler(logHandler)
	//defer clearLog()
	//defer handlePanic(t, "wrong format of address string literal '0x0523029b179009a28a7fae478cd0c2e5ba2adc' with length '40'")
	vm, ef := getVM(eventCodePath, eventAbiPath)

	mutable := true

	proc := exec.NewWavmProcess(vm.VM, vm.Memory, &mutable)

	fmt.Println("testing...")

	address := string("0x0523029b179009a28a7fae478cd0c2e5ba2adc00")

	offset := vm.Memory.SetBytes([]byte(address))

	dest := string("0x523029b179009a28a7fae478cd0c2e5ba2adc00")

	tOffset := ef.AddressFrom(proc, uint64(offset))

	strBytes := vm.Memory.GetPtr(tOffset)

	assert.Equal(t, common.HexToAddress(dest).String(), common.BytesToAddress(strBytes).String())

}
