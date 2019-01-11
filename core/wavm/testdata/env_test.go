package tests

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	T "time"

	"github.com/vntchain/go-vnt/accounts/abi"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core"
	"github.com/vntchain/go-vnt/core/state"
	"github.com/vntchain/go-vnt/core/vm"
	"github.com/vntchain/go-vnt/core/vm/interface"
	"github.com/vntchain/go-vnt/core/wavm"
	errorsmsg "github.com/vntchain/go-vnt/core/wavm/errors"
	"github.com/vntchain/go-vnt/log"
	"github.com/vntchain/go-vnt/params"
	"github.com/vntchain/go-vnt/vntdb"
)

var envJsonPath = filepath.Join("", "env.json")

type ENVTest struct {
	json          vmJSON
	statedb       *state.StateDB
	createCost    float64
	callCost      float64
	compileCost   float64
	nocompileCost float64
	createRunCost float64
	callRunCost   float64
}

func (t *ENVTest) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &t.json)
	if err != nil {
		return err
	}
	return nil
}

func parseInput(args []argument) []interface{} {
	var input []interface{}
	for _, v := range args {
		input = append(input, parseData(v.Data, v.DataType))
	}
	return input
}

func parseData(data string, dataType string) interface{} {
	var parse interface{}
	var err error
	switch dataType {
	case "uint32":
		var res uint64
		res, err = strconv.ParseUint(data, 10, 32)
		if err == nil {
			parse = uint32(res)
		} else {
			fmt.Printf("err %s\n", err.Error())
		}
	case "int32":
		var res int64
		res, err = strconv.ParseInt(data, 10, 32)
		if err == nil {
			parse = int32(res)
		}
	case "uint64":
		var res uint64
		res, err = strconv.ParseUint(data, 10, 64)
		if err == nil {
			parse = uint64(res)
		}
	case "int64":
		var res int64
		res, err = strconv.ParseInt(data, 10, 64)
		if err == nil {
			parse = int64(res)
		}
	case "uint256":
		bigint := new(big.Int)
		_, flag := bigint.SetString(data, 10)
		if flag == false {
			panic("Illegal uint256 input " + data)
		}
		parse = bigint
	case "string":
		parse = data
	case "address":
		if data[0:2] != "0x" {
			parse = common.BytesToAddress([]byte(data))
		} else {
			parse = common.HexToAddress(data)
		}

		fmt.Printf("data %s address %v\n", data, parse)
	case "bool":
		if data == "true" {
			parse = true
		} else {
			parse = false
		}
	default:
		err = errors.New(fmt.Sprintf("unsupport data type %s", dataType))
	}
	if err != nil {
		panic(err)
	}
	return parse
}

func packInput(abiobj abi.ABI, name string, args ...interface{}) []byte {
	if name == "testfunctionnoexist" {
		return []byte("")
	}
	abires := abiobj
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

func unpackOutput(abiobj abi.ABI, v interface{}, name string, output []byte) interface{} {
	abires := abiobj
	err := abires.Unpack(v, name, output)
	if err != nil {
		panic(err)
	}
	return v
}

func (t *ENVTest) newWAVM(statedb *state.StateDB, vmconfig vm.Config) vm.VM {
	canTransfer := func(db inter.StateDB, address common.Address, amount *big.Int) bool {
		return core.CanTransfer(db, address, amount)
	}
	transfer := func(db inter.StateDB, sender, recipient common.Address, amount *big.Int) {
		core.Transfer(db, sender, recipient, amount)
	}
	context := vm.Context{
		CanTransfer: canTransfer,
		Transfer:    transfer,
		GetHash:     vmTestBlockHash,
		Origin:      t.json.Exec.Origin,
		Coinbase:    t.json.Env.Coinbase,
		BlockNumber: new(big.Int).SetUint64(t.json.Env.Number),
		Time:        new(big.Int).SetUint64(t.json.Env.Timestamp),
		GasLimit:    t.json.Env.GasLimit,
		Difficulty:  t.json.Env.Difficulty,
		GasPrice:    t.json.Exec.GasPrice,
	}
	return wavm.NewWAVM(context, statedb, params.MainnetChainConfig, vmconfig)
}

func (t *ENVTest) Run(vmconfig vm.Config, data []byte, iscreate bool, test *testing.T) ([]byte, error) {

	if t.statedb == nil {
		db := vntdb.NewMemDatabase()
		statedb := MakePreState(db, t.json.Pre)
		t.statedb = statedb
	}
	// now := T.Now()
	ret, _, err := t.exec(t.statedb, vmconfig, data, iscreate)
	if err != nil {
		return nil, err
	}
	// duration := T.Since(now)
	// test.Logf("time duration %f", duration.Seconds())
	// t.timeUsed += duration.Seconds()

	// if t.json.GasRemaining == nil {
	// 	if err == nil {
	// 		return fmt.Errorf("gas unspecified (indicating an error), but VM returned no error")
	// 	}
	// 	if gasRemaining > 0 {
	// 		return fmt.Errorf("gas unspecified (indicating an error), but VM returned gas remaining > 0")
	// 	}
	// 	return nil
	// }
	// Test declares gas, expecting outputs to match.
	// if !bytes.Equal(ret, t.json.Out) {
	// 	return fmt.Errorf("return data mismatch: got %x, want %x", ret, t.json.Out)
	// }
	// if gasRemaining != uint64(*t.json.GasRemaining) {
	// 	return fmt.Errorf("remaining gas %v, want %v", gasRemaining, *t.json.GasRemaining)
	// }
	// for addr, account := range t.json.Post {
	// 	for k, wantV := range account.Storage {
	// 		if haveV := statedb.GetState(addr, k); haveV != wantV {
	// 			return fmt.Errorf("wrong storage value at %x:\n  got  %x\n  want %x", k, haveV, wantV)
	// 		}
	// 	}
	// }
	// if root := statedb.IntermediateRoot(false); root != t.json.PostStateRoot {
	// 	return fmt.Errorf("post state root mismatch, got %x, want %x", root, t.json.PostStateRoot)
	// }
	// if logs := rlpHash(statedb.Logs()); logs != common.Hash(t.json.Logs) {
	// 	return fmt.Errorf("post state logs hash mismatch: got %x, want %x", logs, t.json.Logs)
	// }
	return ret, nil
}

func (t *ENVTest) exec(statedb *state.StateDB, vmconfig vm.Config, data []byte, isCreated bool) ([]byte, uint64, error) {
	wavmobj := t.newWAVM(statedb, vmconfig)
	e := t.json.Exec
	if isCreated {
		now := T.Now()
		res, addr, gas, err := wavmobj.Create(vm.AccountRef(e.Caller), data, e.GasLimit, e.Value)
		duration := T.Since(now)
		t.createCost += duration.Seconds()
		t.json.Exec.Address = addr
		// t.compileCost += wavmobj.(*wavm.WAVM).Wavm.VM.CompileTimeCost
		// t.createRunCost += wavmobj.(*wavm.WAVM).CreateTimeCost
		fmt.Printf("create gas cost %d\n", gas)
		return res, gas, err
	} else {
		now := T.Now()
		res, gas, err := wavmobj.Call(vm.AccountRef(e.Caller), e.Address, data, e.GasLimit, e.Value)
		duration := T.Since(now)
		t.callCost += duration.Seconds()
		// t.nocompileCost += wavmobj.(*wavm.WAVM).Wavm.VM.NoCompileTimeCost
		// t.callRunCost += wavmobj.(*wavm.WAVM).CallTimeCost
		return res, gas, err
	}

}

func run(t *testing.T, jspath string) {
	jsonfile, err := ioutil.ReadFile(jspath)
	if err != nil {
		t.Fatalf(err.Error())
	}
	vmconfig := vm.Config{}
	envtest := new(ENVTest)
	envtest.callCost = 0
	envtest.createCost = 0
	err = envtest.UnmarshalJSON(jsonfile)
	if err != nil {
		t.Fatalf(err.Error())
	}
	for i := 0; i < 1; i++ {
		for _, v := range envtest.json.TestCase {

			//init
			if v.InitCase.NeedInit == true {
				code := wavm.WasmCode{}
				code.Code = readFile(filepath.Join(v.Code))
				code.Abi = readFile(filepath.Join(v.Abi))
				parseinput := parseInput(v.InitCase.Input)
				input := packInput(getABI(filepath.Join(v.Abi)), "", parseinput...)
				c, err := json.Marshal(code)
				if err != nil {
					panic(err)
				}
				c = append(c, input...)
				// pre := envtest.json.Pre[envtest.json.Exec.Address]
				// pre.Code = c //[]byte(hexutil.Encode(c))
				// envtest.json.Pre[envtest.json.Exec.Address] = pre
				_, err = envtest.Run(vmconfig, c, true, t)
				if err != nil {
					t.Fatalf(err.Error())
				}
			} else {

			}

			for _, testcase := range v.Tests {
				input := parseInput(testcase.Input)
				abiobj := getABI(filepath.Join(v.Abi))
				pack := packInput(abiobj, testcase.Function, input...)
				ret, err := envtest.Run(vmconfig, pack, false, t)
				if err != nil {
					if strings.HasPrefix(err.Error(), errorsmsg.ErrExecutionAssert.Error()) {
						t.Logf("%s", err.Error())
					} else if err.Error() == errorsmsg.ErrExecutionReverted.Error() {
						t.Logf("%s", errorsmsg.ErrExecutionReverted)
					} else {
						t.Fatal(err)
					}
				}
				verify(t, ret, testcase.Wanted, abiobj, testcase.Function)
				if testcase.Event != nil {
					fmt.Printf("logs %s\n", rlpHash(envtest.statedb.Logs()).Hex())
					res := envtest.statedb.Logs()[0].Data
					fmt.Printf("data %v\n", res)
					type testevent struct {
						Str  string
						Addr common.Address
						U64  uint64
						U32  uint32
						I64  int64
						I32  int32
						U256 *big.Int
						B    bool
					}
					var test1 testevent
					err := abiobj.Unpack(&test1, "TESTEVENT", res)
					if err != nil {
						panic(err)
					}
					fmt.Printf("test1 %+v\n", test1)
					// verifyEvent(t, ret, testcase.Event, abiobj, testcase.Function)
				}
			}

		}
	}

	t.Logf("create cost %f", envtest.createCost/10000.0)
	t.Logf("call cost %f", envtest.callCost/10000.0)
	t.Logf("compile cost %f", envtest.compileCost/10000.0)
	t.Logf("nocompile cost %f", envtest.nocompileCost/10000.0)
	t.Logf("create run cost %f", envtest.compileCost/10000.0)
	t.Logf("call run cost %f", envtest.compileCost/10000.0)

}

func verify(t *testing.T, ret []byte, wanted argument, abiobj abi.ABI, funcName string) {
	t.Logf("funcName %s\n", funcName)
	data := wanted.Data
	dataType := wanted.DataType
	if wanted.DataType == "" {
		return
	}
	parse := parseData(data, dataType)
	switch dataType {
	case "uint32":
		want := parse.(uint32)
		var got uint32
		unpackOutput(abiobj, &got, funcName, ret)
		if got != want {
			t.Fatalf("wavm result mismatch, got %d, want %d", got, want)
		} else {
			t.Logf("wavm result match, got %d, want %d", got, want)
		}
	case "int32":
		want := parse.(int32)
		var got int32
		unpackOutput(abiobj, &got, funcName, ret)
		if got != want {
			t.Fatalf("wavm result mismatch, got %d, want %d", got, want)
		} else {
			t.Logf("wavm result match, got %d, want %d", got, want)
		}
	case "uint64":
		want := parse.(uint64)
		var got uint64
		log.Debug("111", "funcName", funcName, "ret", ret)
		unpackOutput(abiobj, &got, funcName, ret)
		if got != want {
			t.Fatalf("wavm result mismatch, got %d, want %d", got, want)
		} else {
			t.Logf("wavm result match, got %d, want %d", got, want)
		}
	case "int64":
		want := parse.(int64)
		var got int64
		unpackOutput(abiobj, &got, funcName, ret)
		if got != want {
			t.Fatalf("wavm result mismatch, got %d, want %d", got, want)
		} else {
			t.Logf("wavm result match, got %d, want %d", got, want)
		}
	case "uint256":
		want := parse.(*big.Int)
		var got *big.Int
		unpackOutput(abiobj, &got, funcName, ret)
		if got.Cmp(want) != 0 {
			t.Fatalf("wavm result mismatch, got %d, want %d", got, want)
		} else {
			t.Logf("wavm result match, got %d, want %d", got, want)
		}
	case "string":
		want := parse.(string)
		var got string
		unpackOutput(abiobj, &got, funcName, ret)
		if got != want {
			t.Fatalf("wavm result mismatch, got %s, want %s", got, want)
		} else {
			t.Logf("wavm result match, got %s, want %s", got, want)
		}
	case "address":
		want := parse.(common.Address)
		var got common.Address
		unpackOutput(abiobj, &got, funcName, ret)
		if got != want {
			t.Fatalf("wavm result mismatch, got %s, want %s", got.Hex(), want.Hex())
		} else {
			t.Logf("wavm result match, got %s, want %s", got.Hex(), want.Hex())
		}
	case "bool":
		want := parse.(bool)
		var got bool
		unpackOutput(abiobj, &got, funcName, ret)
		if got != want {
			t.Fatalf("wavm result mismatch, got %t, want %t", got, want)
		} else {
			t.Logf("wavm result match, got %t, want %t", got, want)
		}
	}

}

func verifyEvent(t *testing.T, ret []byte, wanted argument, abiobj abi.ABI, funcName string) {
	t.Logf("funcName %s\n", funcName)
	data := wanted.Data
	dataType := wanted.DataType
	if wanted.DataType == "" {
		return
	}
	parse := parseData(data, dataType)
	switch dataType {
	case "uint32":
		want := parse.(uint32)
		var got uint32
		unpackOutput(abiobj, &got, funcName, ret)
		if got != want {
			t.Fatalf("wavm result mismatch, got %d, want %d", got, want)
		} else {
			t.Logf("wavm result match, got %d, want %d", got, want)
		}
	case "int32":
		want := parse.(int32)
		var got int32
		unpackOutput(abiobj, &got, funcName, ret)
		if got != want {
			t.Fatalf("wavm result mismatch, got %d, want %d", got, want)
		} else {
			t.Logf("wavm result match, got %d, want %d", got, want)
		}
	case "uint64":
		want := parse.(uint64)
		var got uint64
		log.Debug("111", "funcName", funcName, "ret", ret)
		unpackOutput(abiobj, &got, funcName, ret)
		if got != want {
			t.Fatalf("wavm result mismatch, got %d, want %d", got, want)
		} else {
			t.Logf("wavm result match, got %d, want %d", got, want)
		}
	case "int64":
		want := parse.(int64)
		var got int64
		unpackOutput(abiobj, &got, funcName, ret)
		if got != want {
			t.Fatalf("wavm result mismatch, got %d, want %d", got, want)
		} else {
			t.Logf("wavm result match, got %d, want %d", got, want)
		}
	case "uint256":
		want := parse.(*big.Int)
		var got *big.Int
		unpackOutput(abiobj, &got, funcName, ret)
		if got.Cmp(want) != 0 {
			t.Fatalf("wavm result mismatch, got %d, want %d", got, want)
		} else {
			t.Logf("wavm result match, got %d, want %d", got, want)
		}
	case "string":
		want := parse.(string)
		var got string
		unpackOutput(abiobj, &got, funcName, ret)
		if got != want {
			t.Fatalf("wavm result mismatch, got %s, want %s", got, want)
		} else {
			t.Logf("wavm result match, got %s, want %s", got, want)
		}
	case "address":
		want := parse.(common.Address)
		var got common.Address
		unpackOutput(abiobj, &got, funcName, ret)
		if got != want {
			t.Fatalf("wavm result mismatch, got %s, want %s", got.Hex(), want.Hex())
		} else {
			t.Logf("wavm result match, got %s, want %s", got.Hex(), want.Hex())
		}
	case "bool":
		want := parse.(bool)
		var got bool
		unpackOutput(abiobj, &got, funcName, ret)
		if got != want {
			t.Fatalf("wavm result mismatch, got %t, want %t", got, want)
		} else {
			t.Logf("wavm result match, got %t, want %t", got, want)
		}
	}

}

func TestEnv(t *testing.T) {
	run(t, envJsonPath)
}
