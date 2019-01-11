package tests

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/vntchain/go-vnt/accounts/abi"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core/state"
	"github.com/vntchain/go-vnt/core/vm"
	"github.com/vntchain/go-vnt/core/vm/interface"
	"github.com/vntchain/go-vnt/core/wavm"
	wasmContract "github.com/vntchain/go-vnt/core/wavm/contract"
	g "github.com/vntchain/go-vnt/core/wavm/gas"
	"github.com/vntchain/go-vnt/core/wavm/storage"
	"github.com/vntchain/go-vnt/crypto"
	"github.com/vntchain/go-vnt/log"
	"github.com/vntchain/go-vnt/params"
	"github.com/vntchain/go-vnt/vntdb"
)

var logger *log.Logger

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
}

// func TestInvokeName(t *testing.T) {
// 	res := wavm.InvokeName([]byte("11111"))
// 	fmt.Printf("%d", res)
// }

// func TestStringToName(t *testing.T) {
// 	res := wavm.StringToName([]byte("issue"))
// 	log.Debug("TestStringToName", "uint", res)
// }

// func TestWasm(t *testing.T) {
// 	code, err := ioutil.ReadFile(filepath.Join("./getnumber/main.wasm")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	log.Debug("wavm_test", "code", code)
// 	wm := wavm.newERC20Wavm()
// 	err = wm.InstantiateModule(code, []uint8{})
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	ctx := &wavm.ApplyContext{
// 		Receiver: 111,
// 		Code:     111,
// 		Action:   111,
// 	}
// 	res, err := wm.Apply(ctx)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	log.Debug("AAA", "res", res)
// }

// func TestSlotTypeToString(t *testing.T) {
// 	value := []wasm.ValueType{wasm.ValueType(0), wasm.ValueType(1), wasm.ValueType(2)}
// 	res := strings.Replace(strings.Trim(fmt.Sprint(value), "[]"), " ", ",", -1)
// 	log.Debug(res)
// }

func getHash(uint64) common.Hash {
	return common.BytesToHash([]byte("0x11111111111111111111111111111111"))
}

var (
	erc20codepath  = filepath.Join(basepath, "erc20/erc20_2_0.wasm")
	erc20abipath   = filepath.Join(basepath, "erc20/abi.json")
	wrcodepath     = filepath.Join(basepath, "readwrite/wr.wasm")
	wrabipath      = filepath.Join(basepath, "readwrite/abi.json")
	callcodepath   = filepath.Join(basepath, "call/call.wasm")
	callabipath    = filepath.Join(basepath, "call/call.abi")
	structcodepath = filepath.Join(basepath, "struct/struct.wasm")
	structabipath  = filepath.Join(basepath, "struct/struct.abi")
)

func getCode(filepath string) []byte {
	code, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	return code
}

func getERC20Code() []byte {
	return getCode(erc20codepath)
}

func getERC20ABI() abi.ABI {
	return getABI(erc20abipath)
}

func getWRCode() []byte {
	return getCode(wrcodepath)
}

func getWRABI() abi.ABI {
	return getABI(wrabipath)
}

func getCALLCode() []byte {
	return getCode(callcodepath)
}

func getCALLABI() abi.ABI {
	return getABI(callabipath)
}

func getStructCode() []byte {
	return getCode(structcodepath)
}

func getStructABI() abi.ABI {
	return getABI(structabipath)
}

var (
	caller          = common.HexToAddress("0xcccccccccccccccccccccccccccccccccccccccc")
	contractAddr    = common.HexToAddress("0xdddddddddddddddddddddddddddddddddddddddd")
	contractBalance = big.NewInt(99999999)
	value           = big.NewInt(1000000)
	gas             = uint64(5000000)
	origin          = common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	gasPrice        = big.NewInt(1000000)
	coinbase        = common.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	gasLimit        = uint64(5000000)
	blockNumber     = big.NewInt(10000)
	time            = big.NewInt(12345678)
	difficulty      = big.NewInt(23456789)
	stateDB         *state.StateDB
)

func newContract(code []byte) *wasmContract.WASMContract {

	contract := wasmContract.NewWASMContract(vm.AccountRef(caller), vm.AccountRef(contractAddr), value, gas)
	contract.SetCallCode(&contractAddr, crypto.Keccak256Hash(code), code)
	return contract
}

func newStatedb() *state.StateDB {
	if stateDB == nil {
		db := vntdb.NewMemDatabase()
		sdb := state.NewDatabase(db)
		stateDB, _ = state.New(common.Hash{}, sdb)
		stateDB.SetBalance(common.HexToAddress("0x553E6c30Af61e7A3576f31311EA8a620F80D047e"), big.NewInt(1150000000))
		root, _ := stateDB.Commit(false)
		stateDB, _ = state.New(root, sdb)
	}
	return stateDB
}
func newWAVM(ctx vm.Context) *wavm.WAVM {
	chainconfig := &params.ChainConfig{HomesteadBlock: big.NewInt(1150000)}
	statedb := newStatedb()
	wavm := wavm.NewWAVM(ctx, statedb, chainconfig, vm.Config{})

	return wavm
}

func newWavm(code, abi string, iscreate bool) *wavm.Wavm {
	statedb := newStatedb()
	tHash := common.HexToHash("0x1111")
	bHash := common.HexToHash("0x2222")
	statedb.Prepare(tHash, bHash, 1)
	statedb.SetBalance(contractAddr, contractBalance)

	contract := newContract(getCode(code))
	config := vm.Config{}
	chainconfig := &params.ChainConfig{HomesteadBlock: big.NewInt(1150000)}
	gasRule := g.NewGas(config.DisableFloatingPoint)
	gasTable := chainconfig.GasTable(blockNumber)
	gasCounter := g.NewGasCounter(contract, gasTable)
	ctx := vm.Context{
		GetHash: getHash,
		// Message information
		Origin:   origin,
		GasPrice: gasPrice,

		// Block information
		Coinbase:    coinbase,
		GasLimit:    gasLimit,
		BlockNumber: blockNumber,
		Time:        time,
		Difficulty:  difficulty,
		CanTransfer: CanTransfer,
		Transfer:    Transfer,
	}
	wvm := newWAVM(ctx)

	chainctx := wavm.ChainContext{
		GetHash: getHash,
		// Message information
		Origin:   origin,
		GasPrice: gasPrice,

		// Block information
		Coinbase:       coinbase,
		GasLimit:       gasLimit,
		BlockNumber:    blockNumber,
		Time:           time,
		Difficulty:     difficulty,
		Contract:       contract,
		StateDB:        statedb,
		Code:           getCode(code),
		Abi:            getABI(abi),
		IsCreated:      iscreate,
		StorageMapping: make(map[uint64]storage.StorageMapping),
		GasRule:        gasRule,
		GasCounter:     gasCounter,
		GasTable:       gasTable,
		Wavm:           wvm,
		CanTransfer:    CanTransfer,
		Transfer:       Transfer,
	}

	vm := wavm.NewWavm(chainctx, config, iscreate)
	err := vm.InstantiateModule(chainctx.Code, nil)
	if err != nil {
		panic(err)
	}
	return vm
}

func newERC20Wavm() *wavm.Wavm {
	return newWavm(erc20codepath, erc20abipath, false)
}

func CanTransfer(db inter.StateDB, addr common.Address, amount *big.Int) bool {
	fmt.Printf("db.GetBalance(addr) %d addr %s\n", db.GetBalance(addr), addr.Hex())
	return db.GetBalance(addr).Cmp(amount) >= 0
}

// Transfer subtracts amount from sender and adds amount to recipient using the given Db
func Transfer(db inter.StateDB, sender, recipient common.Address, amount *big.Int) {
	db.SubBalance(sender, amount)
	db.AddBalance(recipient, amount)
}

func newCallWavm() *wavm.Wavm {
	return newWavm(callcodepath, callabipath, false)
}

func newStructWavm() *wavm.Wavm {
	return newWavm(structcodepath, structabipath, false)
}

func newWRWavm() *wavm.Wavm {
	return newWavm(wrcodepath, wrabipath, false)
}

func pack(vm *wavm.Wavm, name string, args ...interface{}) []byte {
	abires := vm.ChainContext.Abi
	//fmt.Printf("%+v", abires)
	var res []byte
	var err error
	if len(args) == 0 {
		res, err = abires.Pack(name)
	} else {
		res, err = abires.Pack(name, args...)
	}
	if err != nil {
		panic(err)
	}
	return res
}

func unPack(vm *wavm.Wavm, v interface{}, name string, output []byte) interface{} {
	abires := vm.ChainContext.Abi
	err := abires.Unpack(v, name, output)
	if err != nil {
		panic(err)
	}
	return v
}

func TestGetBlockHash(t *testing.T) {
	vm := newERC20Wavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	res, err := vm.Apply(pack(vm, "TestGetBlockHash", uint64(10000)), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	t.Logf("TestGetBlockHash %s", res)
	res, err = vm.Apply(pack(vm, "TestGetBlockHash", uint64(9999)), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	t.Logf("TestGetBlockHash %s", res)
	res, err = vm.Apply(pack(vm, "TestGetBlockHash", uint64(10000-257)), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	t.Logf("TestGetBlockHash %s", res)
}

func TestGetBlockNumber(t *testing.T) {
	vm := newERC20Wavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	res, err := vm.Apply(pack(vm, "TestGetBlockNumber"), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var Int uint64
	unPack(vm, &Int, "TestGetBlockNumber", res)
	if Int != blockNumber.Uint64() {
		t.Errorf("unexpected value : want %x, got %x", blockNumber.Uint64(), Int)
	}
	t.Logf("TestGetBlockNumber %d", Int)
}

func TestGetTimestamp(t *testing.T) {
	vm := newERC20Wavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	res, err := vm.Apply(pack(vm, "TestGetTimestamp"), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var Int uint64
	unPack(vm, &Int, "TestGetTimestamp", res)
	if Int != time.Uint64() {
		t.Errorf("unexpected value : want %x, got %x", time.Uint64(), Int)
	}
	t.Logf("TestGetTimestamp %d", Int)
}

func TestGetCoinBase(t *testing.T) {
	vm := newERC20Wavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	res, err := vm.Apply(pack(vm, "TestGetCoinBase"), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var addr common.Address
	unPack(vm, &addr, "TestGetCoinBase", res)
	if !bytes.Equal(addr.Bytes(), coinbase.Bytes()) {
		t.Errorf("unexpected value : want %s, got %s", coinbase, addr)
	}
	t.Logf("TestGetCoinBase %s", addr)
}

func TestGetGas(t *testing.T) {
	vm := newERC20Wavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	res, err := vm.Apply(pack(vm, "TestGetGas"), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var Int uint64
	unPack(vm, &Int, "TestGetGas", res)
	if Int != 999979 {
		t.Errorf("unexpected value : want %d, got %d", 999979, Int)
	}
	t.Logf("TestGetGas %d", Int)
}

func TestGetGasLimit(t *testing.T) {
	vm := newERC20Wavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	res, err := vm.Apply(pack(vm, "TestGetGasLimit"), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var Int uint64
	unPack(vm, &Int, "TestGetGasLimit", res)
	if Int != gasLimit {
		t.Errorf("unexpected value : want %d, got %d", gasLimit, Int)
	}
	t.Logf("TestGetGasLimit %d", Int)
}

func TestGetContractAddress(t *testing.T) {
	vm := newERC20Wavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	res, err := vm.Apply(pack(vm, "TestGetContractAddress"), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var addr common.Address
	unPack(vm, &addr, "TestGetContractAddress", res)
	if !bytes.Equal(addr.Bytes(), contractAddr.Bytes()) {
		t.Errorf("unexpected value : want %s, got %s", contractAddr, addr)
	}
	t.Logf("TestGetContractAddress %s", addr)
}

func TestGetSender(t *testing.T) {
	vm := newERC20Wavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	res, err := vm.Apply(pack(vm, "TestGetSender"), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var addr common.Address
	unPack(vm, &addr, "TestGetSender", res)
	if !bytes.Equal(addr.Bytes(), caller.Bytes()) {
		t.Errorf("unexpected value : want %s, got %s", caller, addr)
	}
	t.Logf("TestGetSender %s", addr)
	fmt.Printf("TestGetSender %s", addr.Hex())
}

func TestGetOrigin(t *testing.T) {
	vm := newERC20Wavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	res, err := vm.Apply(pack(vm, "TestGetOrigin"), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var addr common.Address
	unPack(vm, &addr, "TestGetOrigin", res)
	if !bytes.Equal(addr.Bytes(), origin.Bytes()) {
		t.Errorf("unexpected value : want %s, got %s", origin, addr)
	}
	t.Logf("TestGetOrigin %s", addr.Hex())
}

func TestGetValue(t *testing.T) {
	vm := newERC20Wavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	res, err := vm.Apply(pack(vm, "TestGetValue"), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var Int uint64
	unPack(vm, &Int, "TestGetValue", res)
	if Int != value.Uint64() {
		t.Errorf("unexpected value : want %d, got %d", value.Uint64(), Int)
	}
	t.Logf("TestGetValue %d", Int)
}

func TestGetDifficulty(t *testing.T) {
	vm := newERC20Wavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	res, err := vm.Apply(pack(vm, "TestGetDifficulty"), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var Int uint64
	unPack(vm, &Int, "TestGetDifficulty", res)
	if Int != difficulty.Uint64() {
		t.Errorf("unexpected value : want %d, got %d", difficulty.Uint64(), Int)
	}
	t.Logf("TestGetDifficulty %d", Int)
}

func TestGetContractValue(t *testing.T) {
	vm := newERC20Wavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	res, err := vm.Apply(pack(vm, "TestGetContractValue"), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var Int uint64
	unPack(vm, &Int, "TestGetContractValue", res)
	if Int != contractBalance.Uint64() {
		t.Errorf("unexpected value : want %d, got %d", contractBalance.Uint64(), Int)
	}
	t.Logf("TestGetContractValue %d", Int)
}

func TestSHA3(t *testing.T) {
	vm := newERC20Wavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	res, err := vm.Apply(pack(vm, "TestSHA3", "TEST"), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var str string
	unPack(vm, &str, "TestSHA3", res)
	if str != "TEST" {
		t.Errorf("unexpected value : want %s, got %s", "TEST", str)
	}
	t.Logf("TestSHA3 %s", str)
}

func TestStorageReadUint64(t *testing.T) {
	vm := newERC20Wavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	res, err := vm.Apply(pack(vm, "TestStorageReadUint64", "key", uint64(10000)), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var Int uint64
	unPack(vm, &Int, "TestStorageReadUint64", res)
	if Int != 10000 {
		t.Errorf("unexpected value : want %d, got %d", 10000, Int)
	}
	t.Logf("TestStorageReadUint64 %d", Int)
}

func TestGenerateKey(t *testing.T) {
	vm := newERC20Wavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	res, err := vm.Apply(pack(vm, "TestGenerateKey", "key1", "key2", "key3"), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var str string
	unPack(vm, &str, "TestGenerateKey", res)
	if str != "key1_key2_key3" {
		t.Errorf("unexpected value : want %s, got %s", "key1_key2_key3", str)
	}
	t.Logf("TestGenerateKey %s", str)
}

// [120 64 116 251 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 64 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 128 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 3 107 101 121 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 5 118 97 108 117 101 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
// 36mDEBUG[0m[07-11|16:37:14] packBytesSlice                           [36mlen[0m="[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 3]" [36ml[0m=32 [36mbyte[0m="[107 101 121]"
// [36mDEBUG[0m[07-11|16:37:14] abi                                      [36mtype[0m=string [36mvalue[0m=string [36mT[0m=3
// [36mDEBUG[0m[07-11|16:37:14] abi                                      [36mpack[0m=3
// [36mDEBUG[0m[07-11|16:37:14] packBytesSlice                           [36mlen[0m="[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 5]" [36ml[0m=32 [36mbyte[0m="[118 97 108 117 101]"
// 3 107 101 121 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
// 5 118 97 108 117 101 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
// 120 64 116 251
// 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 64
// 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 128
// 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 3
// 107 101 121 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
// 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 44
// 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 49 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
func TestStorageReadString(t *testing.T) {
	vm := newERC20Wavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	val := ""
	for i := 0; i < 256; i++ {
		val = val + "1"
	}
	fmt.Printf("val___ %s", val)
	res, err := vm.Apply(pack(vm, "TestStorageReadString", "key", val), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var str string
	unPack(vm, &str, "TestStorageReadString", res)
	if str != val {
		t.Errorf("unexpected value : want %s, got %s", val, str)
	}
	t.Logf("TestStorageReadString %s", str)
}

// func TestInvokeName(t *testing.T) {
// 	res := wavm.InvokeName([]byte("11111"))
// 	fmt.Printf("%d", res)
// }

//[166 116 129 167 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 32 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 4 84 69 83 84 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
//[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 32 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 4 84 69 83 84 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
//[84 69 83 84]
//0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 4
//84 69 83 84 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0

//[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 64 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 128 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 3 107 101 121 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 5 118 97 108 117 101 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]

//0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 64
//0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 128
//0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 3
func TestStringWithInput(t *testing.T) {
	vm := newERC20Wavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	//input := []byte("TEST1111111111111111111111111TEST1111111111111111111111111TEST1111111111111111111111111TEST1111111111111111111111111TEST1111111111111111111111111")
	input := []byte("test")
	log.Debug("TestStringWithInput", "input", input)
	res, err := vm.Apply(pack(vm, "TestStringWithInput", string(input)), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var str string
	unPack(vm, &str, "TestStringWithInput", res)
	if str != string(input) {
		t.Errorf("unexpected value : want %s, got %s", string(input), str)
	}
	t.Logf("TestStringWithInput %s", str)
}

func TestString1(t *testing.T) {
	vm := newERC20Wavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	res, err := vm.Apply(pack(vm, "TestString1"), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var str string
	unPack(vm, &str, "TestString1", res)
	if str != "TESTSTRING" {
		t.Errorf("unexpected value : want %s, got %s", "TESTSTRING", str)
	}
	t.Logf("TestString1 %s", str)
}

func TestString2(t *testing.T) {
	vm := newERC20Wavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	res, err := vm.Apply(pack(vm, "TestString2"), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var str string
	unPack(vm, &str, "TestString2", res)
	if str != "TESTSTRING" {
		t.Errorf("unexpected value : want %s, got %s", "TESTSTRING", str)
	}
	t.Logf("TestString2 %s", str)
}

func TestUnpackString(t *testing.T) {
	abistr := []byte(`[ {"name":"init","constant":false,"inputs":[],"outputs":[{"name":"bbb","type":"string"},{"name":"aaa","type":"uint64"}],"payable":false,"stateMutability":"nonpayable","type":"function"}]`)
	buffer := bytes.NewBuffer(abistr)
	abires, err := abi.JSON(buffer)
	if err != nil {
		panic(err)
	}
	var v struct {
		str    string
		inrres uint64
	}
	err = abires.Unpack(&v, "init", []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 84, 69, 83, 84, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	if err != nil {
		panic(err)
	}
	fmt.Print(v)
}

func TestIntMax(t *testing.T) {
	const INT_MAX = int(^uint(0) >> 1)
	const INT32_MAX = int32(^uint32(0) >> 1)
	const UINT_MAX = ^uint(0)
	const UINT32_MAX = ^uint32(0)
	const INT_MIN = ^INT_MAX
	//conts UINT_MIN = 0
	fmt.Println(INT_MAX)
	fmt.Println(INT32_MAX)
	fmt.Println(UINT_MAX)
	fmt.Println(UINT32_MAX)
	fmt.Println(INT_MIN)
}

// 9223372036854775807
// 2147483647
// 18446744073709551615
// 4294967295
// -9223372036854775808
func TestToI32(t *testing.T) {
	vm := newERC20Wavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	i32str := "-2147483647"
	res, err := vm.Apply(pack(vm, "TestToI32", i32str), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var i32 int32
	unPack(vm, &i32, "TestToI32", res)
	if i32 != -2147483647 {
		t.Errorf("unexpected value : want %s, got %d", "-2147483647", i32)
	}
	t.Logf("TestToI32 %d", i32)
}

func TestToI64(t *testing.T) {
	vm := newERC20Wavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	i64str := "-9223372036854775808"
	res, err := vm.Apply(pack(vm, "TestToI64", i64str), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var i64 int64
	unPack(vm, &i64, "TestToI64", res)
	if i64 != -9223372036854775808 {
		t.Errorf("unexpected value : want %s, got %d", "-9223372036854775808", i64)
	}
	t.Logf("TestToI64 %d", i64)
}

func TestToU32(t *testing.T) {
	vm := newERC20Wavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	u32str := "4294967295"
	res, err := vm.Apply(pack(vm, "TestToU32", u32str), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var u32 uint32
	unPack(vm, &u32, "TestToU32", res)
	if u32 != 4294967295 {
		t.Errorf("unexpected value : want %s, got %d", "4294967295", u32)
	}
	t.Logf("TestToU32 %d", u32)
}

func TestToU64(t *testing.T) {
	vm := newERC20Wavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	u64str := "18446744073709551615"
	res, err := vm.Apply(pack(vm, "TestToU64", u64str), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var u64 uint64
	unPack(vm, &u64, "TestToU64", res)
	if u64 != 18446744073709551615 {
		t.Errorf("unexpected value : want %s, got %d", "18446744073709551615", u64)
	}
	t.Logf("TestToU64 %d", u64)
}

func TestFromI32(t *testing.T) {
	vm := newERC20Wavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	i32 := int32(-2147483647)
	res, err := vm.Apply(pack(vm, "TestfromI32", i32), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var i32str string
	unPack(vm, &i32str, "TestfromI32", res)
	if i32str != "-2147483647" {
		t.Errorf("unexpected value : want %s, got %s", "4294967295", i32str)
	}
	t.Logf("TestfromI32 %s", i32str)
}

func TestFromI64(t *testing.T) {
	vm := newERC20Wavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	i64 := int64(-9223372036854775808)
	res, err := vm.Apply(pack(vm, "TestfromI64", i64), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var i64str string
	unPack(vm, &i64str, "TestfromI64", res)

	if i64str != "-9223372036854775808" {
		t.Errorf("unexpected value : want %s, got %s", "-9223372036854775808", i64str)
	}
	t.Logf("TestfromI64 %s", i64str)
}

func TestFromU32(t *testing.T) {
	vm := newERC20Wavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	u32 := uint32(4294967295)
	res, err := vm.Apply(pack(vm, "TestfromU32", u32), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var u32str string
	unPack(vm, &u32str, "TestfromU32", res)
	if u32str != "4294967295" {
		t.Errorf("unexpected value : want %s, got %s", "4294967295", u32str)
	}
	t.Logf("TestfromU32 %s", u32str)
}

func TestFromU64(t *testing.T) {
	vm := newERC20Wavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	u64 := uint64(18446744073709551615)
	res, err := vm.Apply(pack(vm, "TestfromU64", u64), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var u64str string
	unPack(vm, &u64str, "TestfromU64", res)
	if u64str != "18446744073709551615" {
		t.Errorf("unexpected value : want %s, got %s", "18446744073709551615", u64str)
	}
	t.Logf("TestfromU64 %s", u64str)
}

func TestStruct(t *testing.T) {
	vm := newStructWavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	_, err := vm.Apply(pack(vm, "main"), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	// var u64str string
	// unPack(vm, &u64str, "TestfromU64", res)
	// if u64str != "18446744073709551615" {
	// 	t.Errorf("unexpected value : want %s, got %s", "18446744073709551615", u64str)
	// }
	// t.Logf("TestfromU64 %s", u64str)
}

//todo 参数为负数的时候有问题
func TestRwInt64(t *testing.T) {
	vm := newWRWavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	res, err := vm.Apply(pack(vm, "testrwint64", int64(1000000)), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var i64 int64
	unPack(vm, &i64, "testrwint64", res)
	if i64 != 9223372036854775807 {
		t.Errorf("unexpected value : want %s, got %d", "9223372036854775807", i64)
	}
	t.Logf("testrwint64 %d", i64)
}

func TestRwString(t *testing.T) {
	vm := newWRWavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	param := "111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111"
	res, err := vm.Apply(pack(vm, "testrwstring", param), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var str string
	unPack(vm, &str, "testrwstring", res)
	if str != param {
		t.Errorf("unexpected value : want %s, got %s", param, str)
	}
	t.Logf("testrwstring %s", str)
}

func TestRwAddress(t *testing.T) {
	vm := newWRWavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	param := common.HexToAddress("0x7EF5A6135f1FD6a02593eEdC869c6D41D934aef8")
	t.Logf("address %s", param.String())
	res, err := vm.Apply(pack(vm, "testrwaddress", param), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var addr common.Address
	unPack(vm, &addr, "testrwaddress", res)
	if addr != param {
		t.Errorf("unexpected value : want %s, got %s", param, addr)
	}
	t.Logf("testrwaddress %s", addr.Hex())
	fmt.Printf("testrwaddress %s", addr.Hex())
}

func haTestA(t *testing.T) {
	vm := newWRWavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	res, err := vm.Apply(pack(vm, "testa", int64(-1)), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var i64 int64
	unPack(vm, &i64, "testa", res)
	if i64 != 9223372036854775807 {
		t.Errorf("unexpected value : want %s, got %d", "9223372036854775807", i64)
	}
	t.Logf("TestfromU64 %d", i64)
}

func TestTestmap(t *testing.T) {
	vm := newWRWavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	res, err := vm.Apply(pack(vm, "test1111"), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	var u64 uint64
	unPack(vm, &u64, "test1111", res)
	if u64 != 0 {
		t.Errorf("unexpected value : want %s, got %d", "0", u64)
	}
	t.Logf("TestfromU64 %d", u64)
}

func TestReada(t *testing.T) {
	vm := newStructWavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	_, err := vm.Apply(pack(vm, "main"), nil, mutable)
	if err != nil {
		t.Error(err)
	}
	// var u64str string
	// unPack(vm, &u64str, "TestfromU64", res)
	// if u64str != "18446744073709551615" {
	// 	t.Errorf("unexpected value : want %s, got %s", "18446744073709551615", u64str)
	// }
	// t.Logf("TestfromU64 %s", u64str)
}

func TestLittle(t *testing.T) {
	var endianess = binary.LittleEndian
	var b = make([]byte, 10)
	fmt.Printf("%v", b)
	a := -11111
	endianess.PutUint64(b, uint64(a))
	fmt.Printf("b  %v", b)
	u64 := endianess.Uint64(b)
	fmt.Println(int64(u64))
	//补码
	// var x uint64 = 18446744073709551615
	// var y int64 = int64(x)
	// fmt.Println(y)
}

func TestWriteAndReadInt64(t *testing.T) {
	vm := newWRWavm()
	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
	func() {
		res, err := vm.Apply(pack(vm, "testwriteint64", int64(111)), nil, mutable)
		if err != nil {
			t.Error(err)
		}
		var i64 int64
		unPack(vm, &i64, "testwriteint64", res)
		if i64 != 111 {
			t.Errorf("unexpected value : want %s, got %d", "111", i64)
		}
		t.Logf("testwriteint64 %d", i64)
	}()
	func() {
		res, err := vm.Apply(pack(vm, "testreadint64", int64(111)), nil, mutable)
		if err != nil {
			t.Error(err)
		}
		var i64 int64
		unPack(vm, &i64, "testreadint64", res)
		if i64 != 111 {
			t.Errorf("unexpected value : want %s, got %d", "111", i64)
		}
		t.Logf("testwriteint64 %d", i64)
		fmt.Printf("testwriteint64 %d", i64)
	}()
}

func TestBigInt(t *testing.T) {
	loc0 := new(big.Int).SetBytes([]byte("address##0"))
	loc1 := new(big.Int).SetBytes([]byte("address##1"))
	fmt.Printf("%s", common.BigToHash(loc0))
	fmt.Printf("%s", common.BigToHash(loc1))
}
