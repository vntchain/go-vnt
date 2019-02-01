package wavm

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strconv"

	"github.com/vntchain/go-vnt/common/math"

	"github.com/vntchain/go-vnt/accounts/abi"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core/types"
	errormsg "github.com/vntchain/go-vnt/core/wavm/errors"
	"github.com/vntchain/go-vnt/core/wavm/storage"
	"github.com/vntchain/go-vnt/core/wavm/utils"
	"github.com/vntchain/go-vnt/crypto"
	"github.com/vntchain/go-vnt/log"
	"github.com/vntchain/go-vnt/params"
	"github.com/vntchain/vnt-wasm/exec"
	"github.com/vntchain/vnt-wasm/wasm"
)

var (
	errExceededArray = errors.New("array length exceeded")
)

var endianess = binary.LittleEndian

type EnvFunctions struct {
	ctx       *ChainContext
	funcTable map[string]wasm.Function
}

func (ef *EnvFunctions) InitFuncTable(context *ChainContext) {
	ef.ctx = context
	ef.funcTable = ef.getFuncTable()

	// process events
	if ef.ctx == nil {
		return
	}
	for _, event := range ef.ctx.Abi.Events {
		paramTypes := make([]wasm.ValueType, len(event.Inputs))
		for index, input := range event.Inputs {
			switch input.Type.String() {
			case "uint64", "int64":
				paramTypes[index] = wasm.ValueTypeI64
			case "uint32", "int32", "address", "string", "bool", "uint256":
				paramTypes[index] = wasm.ValueTypeI32
			default:
				panic("unsupported type " + input.Type.String())
			}
		}
		//ef.funcTable[event.Name] = reflect.ValueOf(ef.getEvent(len(event.Inputs), event.Name))
		ef.funcTable[event.Name] = wasm.Function{
			Host: reflect.ValueOf(ef.getEvent(event.Name)),
			Sig: &wasm.FunctionSig{
				ParamTypes:  paramTypes,
				ReturnTypes: []wasm.ValueType{},
			},
			Body: &wasm.FunctionBody{
				Code: []byte{},
			},
		}

	}

	// process contract calls
	for _, call := range ef.ctx.Abi.Calls {
		paramTypes := make([]wasm.ValueType, len(call.Inputs)+1)
		paramTypes[0] = wasm.ValueTypeI32
		for index, input := range call.Inputs {
			idx := index + 1
			switch input.Type.String() {
			case "uint64", "int64":
				paramTypes[idx] = wasm.ValueTypeI64
			case "uint32", "int32", "address", "string", "bool", "uint256":
				paramTypes[idx] = wasm.ValueTypeI32
			default:
				panic("unsupported type " + input.Type.String())
			}
		}
		returnTypes := make([]wasm.ValueType, len(call.Outputs))
		for index, output := range call.Outputs {
			switch output.Type.String() {
			case "uint64", "int64":
				returnTypes[index] = wasm.ValueTypeI64
			case "uint32", "int32", "address", "string", "bool", "uint256":
				returnTypes[index] = wasm.ValueTypeI32
			default:
				panic("unsupported type " + output.Type.String())
			}
		}
		ef.funcTable[call.Name] = wasm.Function{
			Host: reflect.ValueOf(ef.getContractCall(call.Name)),
			Sig: &wasm.FunctionSig{
				ParamTypes:  paramTypes,
				ReturnTypes: returnTypes,
			},
			Body: &wasm.FunctionBody{
				Code: []byte{},
			},
		}
	}
}

func (ef *EnvFunctions) GetFuncTable() map[string]wasm.Function {
	return ef.funcTable
}

//todo uint64 =>uint256
func (ef *EnvFunctions) GetBalanceFromAddress(proc *exec.WavmProcess, locIndex uint64) uint64 {
	ef.ctx.GasCounter.GasGetBalanceFromAddress()
	ctx := ef.ctx
	addr := common.BytesToAddress(proc.ReadAt(locIndex))
	balance := ctx.StateDB.GetBalance(addr)
	return ef.returnU256(proc, balance)
}

func (ef *EnvFunctions) GetBlockNumber(proc *exec.WavmProcess) uint64 {
	ef.ctx.GasCounter.GasGetBlockNumber()
	return ef.ctx.BlockNumber.Uint64()
}

func (ef *EnvFunctions) GetGas(proc *exec.WavmProcess) uint64 {
	//当前剩余gas
	ef.ctx.GasCounter.GasGetGas()
	return ef.ctx.Contract.Gas
}

func (ef *EnvFunctions) GetBlockHash(proc *exec.WavmProcess, blockNum uint64) uint64 {
	ctx := ef.ctx
	ctx.GasCounter.GasGetBlockHash()
	num := new(big.Int).SetUint64(blockNum)
	bigval := new(big.Int)
	n := bigval.Sub(ctx.BlockNumber, common.Big257)
	if num.Cmp(n) > 0 && num.Cmp(ctx.BlockNumber) < 0 {
		bhash := ctx.GetHash(num.Uint64())
		return ef.returnHash(proc, []byte(bhash.Hex()))
	} else {

		return ef.returnHash(proc, []byte(common.Hash{}.Hex()))
	}
}

func (ef *EnvFunctions) GetBlockProduser(proc *exec.WavmProcess) uint64 {
	ctx := ef.ctx
	ctx.GasCounter.GasGetBlockProduser()
	coinbase := ctx.Coinbase.Bytes()
	return ef.returnAddress(proc, coinbase)
}

func (ef *EnvFunctions) GetTimestamp(proc *exec.WavmProcess) uint64 {
	ef.ctx.GasCounter.GasGetTimestamp()
	return ef.ctx.Time.Uint64()
}

func (ef *EnvFunctions) GetOrigin(proc *exec.WavmProcess) uint64 {
	ctx := ef.ctx
	ctx.GasCounter.GasGetOrigin()
	origin := ctx.Origin.Bytes()
	return ef.returnAddress(proc, origin)
}

func (ef *EnvFunctions) GetSender(proc *exec.WavmProcess) uint64 {
	ctx := ef.ctx
	ctx.GasCounter.GasGetSender()
	sender := ctx.Contract.CallerAddress.Bytes()
	return ef.returnAddress(proc, sender)
}

func (ef *EnvFunctions) GetGasLimit(proc *exec.WavmProcess) uint64 {
	ctx := ef.ctx
	ctx.GasCounter.GasGetGasLimit()
	return ctx.GasLimit
}

//todo 不能转成uint64 必须是uint256
func (ef *EnvFunctions) GetValue(proc *exec.WavmProcess) uint64 {
	ctx := ef.ctx
	ctx.GasCounter.GasGetValue()
	val := ctx.Contract.Value()
	return ef.returnU256(proc, val)
}

func (ef *EnvFunctions) SHA3(proc *exec.WavmProcess, dataIdx uint64) uint64 {
	data := proc.ReadAt(dataIdx)
	ef.ctx.GasCounter.GasSHA3(uint64(len(data)))
	hash := []byte(crypto.Keccak256Hash(data).Hex())
	return uint64(proc.SetBytes(hash))
}

func (ef *EnvFunctions) GetContractAddress(proc *exec.WavmProcess) uint64 {
	ctx := ef.ctx
	ctx.GasCounter.GasGetContractAddress()
	addr := ctx.Contract.Address().Bytes()
	return ef.returnAddress(proc, addr)
}

func (ef *EnvFunctions) Assert(proc *exec.WavmProcess, condition uint64, msgIdx uint64) {
	ef.ctx.GasCounter.GasAssert()
	msg := proc.ReadAt(msgIdx)
	if condition != 1 {
		err := fmt.Sprintf("%s: %s", errormsg.ErrExecutionAssert, string(msg))
		log.Error(err)
		panic(err)
	}
}

func (ef *EnvFunctions) SendFromContract(proc *exec.WavmProcess, addrIdx uint64, amountIdx uint64) {
	log.Debug("instructions", "func", "SendFromContract")
	ef.forbiddenMutable(proc)
	addr := common.BytesToAddress(proc.ReadAt(addrIdx))
	amount := utils.GetU256(proc.ReadAt(amountIdx))
	if ef.ctx.CanTransfer(ef.ctx.StateDB, ef.ctx.Contract.Address(), amount) {
		_, returnGas, err := ef.ctx.Wavm.Call(ef.ctx.Contract, addr, nil, params.CallStipend, amount)
		ef.ctx.GasCounter.Charge(returnGas)
		if err != nil {
			panic(errormsg.ErrExecutionReverted)
		}
	} else {
		panic(errormsg.ErrExecutionReverted)
	}
}

func (ef *EnvFunctions) TransferFromContract(proc *exec.WavmProcess, addrIdx uint64, amountIdx uint64) uint64 {
	log.Debug("instructions", "func", "TransferFromContract")
	ef.forbiddenMutable(proc)
	// ef.ctx.GasCounter.GasSendFromContract()
	addr := common.BytesToAddress(proc.ReadAt(addrIdx))
	amount := utils.GetU256(proc.ReadAt(amountIdx))
	if ef.ctx.CanTransfer(ef.ctx.StateDB, ef.ctx.Contract.Address(), amount) {
		_, returnGas, err := ef.ctx.Wavm.Call(ef.ctx.Contract, addr, nil, params.CallStipend, amount)
		ef.ctx.GasCounter.Charge(returnGas)
		if err != nil {
			return 0
		}
		return 1
	}
	return 0
}

func (ef *EnvFunctions) fromI64(proc *exec.WavmProcess, value uint64) uint64 {
	ef.ctx.GasCounter.GasFromI64()
	amount := int(value)
	str := strconv.Itoa(amount)
	return uint64(proc.SetBytes([]byte(str)))
}

func (ef *EnvFunctions) fromU64(proc *exec.WavmProcess, amount uint64) uint64 {
	ef.ctx.GasCounter.GasFromU64()
	str := strconv.FormatUint(amount, 10)
	return uint64(proc.SetBytes([]byte(str)))
}

func (ef *EnvFunctions) toI64(proc *exec.WavmProcess, strIdx uint64) int64 {
	ef.ctx.GasCounter.GasToI64()
	b := proc.ReadAt(strIdx)
	amount, err := strconv.Atoi(string(b))
	if err != nil {
		return 0
	} else {
		return int64(amount)
	}
}

func (ef *EnvFunctions) toU64(proc *exec.WavmProcess, strIdx uint64) uint64 {
	ef.ctx.GasCounter.GasToU64()
	b := proc.ReadAt(strIdx)
	amount, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		return 0
	} else {
		return amount
	}
}

func (ef *EnvFunctions) Concat(proc *exec.WavmProcess, str1Idx uint64, str2Idx uint64) uint64 {
	str1 := proc.ReadAt(str1Idx)
	str2 := proc.ReadAt(str2Idx)
	ef.ctx.GasCounter.GasConcat(uint64(len(str1) + len(str2)))
	res := append(str1, str2...)
	return uint64(proc.SetBytes(res))
}

func (ef *EnvFunctions) Equal(proc *exec.WavmProcess, str1Idx uint64, str2Idx uint64) uint64 {
	ef.ctx.GasCounter.GasEqual()
	str1 := proc.ReadAt(str1Idx)
	str2 := proc.ReadAt(str2Idx)
	res := bytes.Equal(str1, str2)
	if res {
		return 1
	} else {
		return 0
	}
}

func (ef *EnvFunctions) getEvent(funcName string) interface{} {
	fnDef := func(proc *exec.WavmProcess, vars ...uint64) {
		ef.forbiddenMutable(proc)
		Abi := ef.ctx.Abi

		var event abi.Event
		var ok bool
		if event, ok = Abi.Events[funcName]; !ok {
			panic(fmt.Sprintf("event execution failed: there is no event '%s' in abi", funcName))
		}

		abiParamLen := len(event.Inputs)
		paramLen := len(vars)

		if abiParamLen != paramLen {
			panic(fmt.Sprintf("event execution failed: there is no event '%s' in abi", funcName))
		}

		topics := make([]common.Hash, 0)
		data := make([]byte, 0)

		topics = append(topics, event.Id())

		strStartIndex := make([]int, 0)
		strData := make([][]byte, 0)

		for i := 0; i < paramLen; i++ {
			input := event.Inputs[i]
			indexed := input.Indexed
			paramType := input.Type.String()
			param := vars[i]
			var value []byte
			switch paramType {
			case "address":
				value = proc.ReadAt(param)
			case "string":
				value = proc.ReadAt(param)
			case "uint64", "int64":
				// value = abi.U256(new(big.Int).SetUint64(param))
				value = make([]byte, 8)
				binary.BigEndian.PutUint64(value, uint64(param))
			case "uint32", "int32", "bool":
				value = make([]byte, 4)
				binary.BigEndian.PutUint32(value, uint32(param))
			case "uint256":
				// ef.readU256FromMemory(proc, param)
				mem := proc.ReadAt(param)
				value = abi.U256(utils.GetU256(mem))
			}

			if indexed {
				topic := common.BytesToHash(value)
				topics = append(topics, topic)
			} else {
				if paramType == "string" {
					strStartIndex = append(strStartIndex, len(data))
					data = append(data, make([]byte, 32)...)
					strData = append(strData, value)
				} else {
					data = append(data, common.LeftPadBytes(value, 32)...)
				}
			}
		}

		// append the string data at the end of the data, and
		// update the start position of string data
		if len(strStartIndex) > 0 {
			for i := range strStartIndex {
				value := strData[i]
				startPos := abi.U256(new(big.Int).SetUint64(uint64(len(data))))
				copy(data[strStartIndex[i]:], startPos)

				size := abi.U256(new(big.Int).SetUint64(uint64(len(value))))
				data = append(data, size...)
				data = append(data, common.RightPadBytes(value, (len(value)+31)/32*32)...)
			}
		}

		log.Debug("Will add event log: ", "topics", topics, "data", data)
		ef.ctx.StateDB.AddLog(&types.Log{
			Address:     ef.ctx.Contract.Address(),
			Topics:      topics,
			Data:        data,
			BlockNumber: ef.ctx.BlockNumber.Uint64(),
		})
		ef.ctx.GasCounter.GasLog(uint64(len(data)), uint64(len(topics)))
		log.Debug("Added event log.")
	}

	return fnDef
}

//todo 如果一个unmutable的方法跨合约调用了一个mutable的方法 则会报错
func (ef *EnvFunctions) getContractCall(funcName string) interface{} {
	Abi := ef.ctx.Abi

	var dc abi.Method
	var ok bool
	if dc, ok = Abi.Calls[funcName]; !ok {
		panic(fmt.Sprintf("call execution failed: Can not find call '%s' in abi", funcName))
	}

	fnDef := func(proc *exec.WavmProcess, vars ...uint64) interface{} {
		// ef.forbiddenMutable(proc)
		abiParamLen := len(dc.Inputs)
		paramLen := len(vars)
		if abiParamLen+1 != paramLen {
			panic(fmt.Sprintf("call execution failed: there is no such call '%s' in abi", funcName))
		}

		args := []interface{}{}
		for i := 1; i < paramLen; i++ {
			input := dc.Inputs[i-1]
			paramType := input.Type.String()
			param := vars[i]
			var value []byte

			switch paramType {
			case "address":
				value = proc.ReadAt(param)
				addr := common.BytesToAddress(value)
				args = append(args, addr)
			case "string":
				value = proc.ReadAt(param)
				args = append(args, string(value))
			case "uint64":
				args = append(args, uint64(param))
			case "int64":
				args = append(args, int64(param))
			case "uint32":
				args = append(args, uint32(param))
			case "int32":
				args = append(args, int32(param))
			case "uint256":
				mem := proc.ReadAt(param)
				bigint := utils.GetU256(mem)
				args = append(args, bigint)
			case "bool":
				arg := false
				if param == 1 {
					arg = true
				}
				args = append(args, arg)
			default:
				panic("unsupport type " + paramType)
			}
		}
		var res []byte
		var err error
		if len(args) == 0 {
			res, err = Abi.Pack(funcName)
			if err != nil {
				panic(err.Error())
			}
		} else {
			res, err = Abi.Pack(funcName, args...)
			if err != nil {
				panic(err.Error())
			}
		}
		toAddr := common.BytesToAddress(proc.ReadAt(uint64(endianess.Uint32(proc.GetData()[vars[0]:]))))
		amount := readU256FromMemory(proc, uint64(endianess.Uint32(proc.GetData()[vars[0]+4:])))
		gascost := endianess.Uint64(proc.GetData()[vars[0]+8:])
		// Get arguments from the memory.
		// ef.ctx.GetWavm().GetCallGasTemp
		//todo
		var gaslimit *big.Int
		gaslimit = new(big.Int).SetUint64(gascost)
		gas := ef.ctx.GasCounter.GasCall(toAddr, amount, gaslimit, ef.ctx.BlockNumber, ef.ctx.Wavm.GetChainConfig(), ef.ctx.StateDB)
		ef.ctx.Wavm.SetCallGasTemp(gas)
		//免费提供额外的gas
		if amount.Sign() != 0 {
			gas += params.CallStipend
		}
		ret, returnGas, err := ef.ctx.Wavm.Call(ef.ctx.Contract, toAddr, res, gas, amount)
		log.Debug("instructions", "func", "contractcall", "ret", ret, "gas", gas, "returnGas", returnGas, "err", err, "gasused", gas-returnGas)
		failError := errors.New("failed to get result in contract call.")
		if err != nil {
			e := fmt.Errorf("%s Reason : %s", failError, err)
			panic(e)
		} else {
			ef.ctx.Contract.Gas += returnGas
			if len(dc.Outputs) == 0 {
				return nil
			} else {
				t := dc.Outputs[0].Type
				switch t.String() {
				case "string":
					var unpackres string
					if err := Abi.Unpack(&unpackres, funcName, ret); err != nil {
						panic(failError)
					} else {
						return uint32(proc.SetBytes([]byte(unpackres)))
					}
				case "address":
					var unpackres common.Address
					if err := Abi.Unpack(&unpackres, funcName, ret); err != nil {
						panic(failError)
					} else {
						return uint32(proc.SetBytes(unpackres.Bytes()))
					}
				case "uint64":
					var unpackres uint64
					if err := Abi.Unpack(&unpackres, funcName, ret); err != nil {
						panic(failError)
					} else {
						return uint64(unpackres)
					}
				case "int64":
					var unpackres int64
					if err := Abi.Unpack(&unpackres, funcName, ret); err != nil {
						panic(failError)
					} else {
						return int64(unpackres)
					}
				case "uint32":
					var unpackres uint32
					if err := Abi.Unpack(&unpackres, funcName, ret); err != nil {
						panic(failError)
					} else {
						return uint32(unpackres)
					}
				case "int32":
					var unpackres int32
					if err := Abi.Unpack(&unpackres, funcName, ret); err != nil {
						panic(failError)
					} else {
						return int32(unpackres)
					}
				case "uint256":
					var unpackres *big.Int
					if err := Abi.Unpack(&unpackres, funcName, ret); err != nil {
						panic(failError)
					} else {
						return uint32(proc.SetBytes([]byte(unpackres.String())))
					}
				case "bool":
					var unpackres bool
					if err := Abi.Unpack(&unpackres, funcName, ret); err != nil {
						panic(failError)
					} else {
						if unpackres == true {
							return int32(1)
						} else {
							return int32(0)
						}

					}
				default:
					panic("unsupport type " + t.String())
				}
			}

		}
	}

	funcVoid := func(proc *exec.WavmProcess, vars ...uint64) {
		fnDef(proc, vars...)
	}

	funcUint64 := func(proc *exec.WavmProcess, vars ...uint64) uint64 {
		return fnDef(proc, vars...).(uint64)
	}

	funcUint32 := func(proc *exec.WavmProcess, vars ...uint64) uint32 {
		return fnDef(proc, vars...).(uint32)
	}

	funcInt64 := func(proc *exec.WavmProcess, vars ...uint64) int64 {
		return fnDef(proc, vars...).(int64)
	}

	funcInt32 := func(proc *exec.WavmProcess, vars ...uint64) int32 {
		return fnDef(proc, vars...).(int32)
	}

	if len(dc.Outputs) == 0 {
		return funcVoid
	} else {
		switch dc.Outputs[0].Type.String() {
		case "uint64":
			return funcUint64
		case "string", "address", "uint32", "uint256":
			return funcUint32
		case "int64":
			return funcInt64
		case "int32":
			return funcInt32
		default:
			return nil
		}
	}

	//return makeFunc(fnDef)
}

// End the line
func (ef *EnvFunctions) printLine(msg string) error {
	funcName := ef.ctx.Wavm.Wavm.GetFuncName()
	log.Info("Contract Debug >>>>", "func", funcName, "message", msg)
	return nil
}

func (ef *EnvFunctions) GetPrintRemark(proc *exec.WavmProcess, remarkIdx uint64) string {
	strValue := proc.ReadAt(remarkIdx)
	return string(strValue)
}

// Print an Address
func (ef *EnvFunctions) PrintAddress(proc *exec.WavmProcess, remarkIdx uint64, strIdx uint64) {
	addrValue := proc.ReadAt(strIdx)
	msg := fmt.Sprint(ef.GetPrintRemark(proc, remarkIdx), common.BytesToAddress(addrValue).String())
	ef.printLine(msg)
}

// Print a string
func (ef *EnvFunctions) PrintStr(proc *exec.WavmProcess, remarkIdx uint64, strIdx uint64) {
	strValue := proc.ReadAt(strIdx)
	msg := fmt.Sprint(ef.GetPrintRemark(proc, remarkIdx), string(strValue))
	ef.printLine(msg)
}

// Print a string
func (ef *EnvFunctions) PrintQStr(proc *exec.WavmProcess, remarkIdx uint64, strIdx uint64) {

	size := endianess.Uint32(proc.GetData()[strIdx : strIdx+4])
	offset := endianess.Uint32(proc.GetData()[strIdx+4 : strIdx+8])
	strValue := proc.GetData()[offset : offset+size]
	length := len(proc.GetData())
	if length > 128 {
		length = 128
	}
	log.Debug("memory", "data", proc.GetData()[0:length])
	log.Debug("PrintQStr", "remarkIdx", remarkIdx, "strIdx", strIdx, "offset", offset, "size", size, "data", strValue)
	// msg := fmt.Sprint(ef.GetPrintRemark(proc, remarkIdx), string(strValue))
	msg := fmt.Sprint(ef.GetPrintRemark(proc, remarkIdx), hex.EncodeToString(strValue))
	ef.printLine(msg)
}

// Print a uint64
func (ef *EnvFunctions) PrintUint64T(proc *exec.WavmProcess, remarkIdx uint64, intValue uint64) {
	msg := fmt.Sprint(ef.GetPrintRemark(proc, remarkIdx), intValue)
	ef.printLine(msg)
}

// Print a uint32
func (ef *EnvFunctions) PrintUint32T(proc *exec.WavmProcess, remarkIdx uint64, intValue uint64) {
	msg := fmt.Sprint(ef.GetPrintRemark(proc, remarkIdx), uint32(intValue))
	ef.printLine(msg)
}

// Print a int64
func (ef *EnvFunctions) PrintInt64T(proc *exec.WavmProcess, remarkIdx uint64, intValue uint64) {
	msg := fmt.Sprint(ef.GetPrintRemark(proc, remarkIdx), int64(intValue))
	ef.printLine(msg)
}

// Print a int32
func (ef *EnvFunctions) PrintInt32T(proc *exec.WavmProcess, remarkIdx uint64, intValue uint64) {
	msg := fmt.Sprint(ef.GetPrintRemark(proc, remarkIdx), int32(intValue))
	ef.printLine(msg)
}

func (ef *EnvFunctions) PrintUint256T(proc *exec.WavmProcess, remarkIdx uint64, idx uint64) {
	u256 := readU256FromMemory(proc, idx)
	msg := fmt.Sprint(ef.GetPrintRemark(proc, remarkIdx), u256.String())
	ef.printLine(msg)
}

func (ef *EnvFunctions) AddressFrom(proc *exec.WavmProcess, idx uint64) uint64 {
	ctx := ef.ctx
	ctx.GasCounter.GasQuickStep()
	addrStr := string(proc.ReadAt(idx))
	address := common.HexToAddress(addrStr).Bytes()
	return ef.returnAddress(proc, address)
}

func (ef *EnvFunctions) AddressToString(proc *exec.WavmProcess, idx uint64) uint64 {
	ctx := ef.ctx
	ctx.GasCounter.GasQuickStep()
	addrBytes := proc.ReadAt(idx)
	address := common.BytesToAddress(addrBytes)
	return uint64(proc.SetBytes([]byte(address.Hex())))
}

func (ef *EnvFunctions) U256From(proc *exec.WavmProcess, idx uint64) uint64 {
	ctx := ef.ctx
	ctx.GasCounter.GasQuickStep()
	u256Str := string(proc.ReadAt(idx))
	bigint, success := new(big.Int).SetString(u256Str, 10)
	if success != true {
		panic(fmt.Sprintf("Can't Convert strin %s to uint256", u256Str))
	}
	return ef.returnU256(proc, bigint)
}

func (ef *EnvFunctions) U256ToString(proc *exec.WavmProcess, idx uint64) uint64 {
	ctx := ef.ctx
	ctx.GasCounter.GasQuickStep()
	u256Bytes := proc.ReadAt(idx)
	return uint64(proc.SetBytes(u256Bytes))
}

// // Open for unit testing
// func (ef *EnvFunctions) TestWritePerType(proc *exec.WavmProcess, typ abi.Type, validx uint64, loc common.Hash) {
// 	ef.writePerType(proc, typ, validx, loc)
// }

// func (ef *EnvFunctions) TestReadPerType(proc *exec.WavmProcess, typ abi.Type, val []byte, loc common.Hash) uint64 {
// 	return ef.readPerType(proc, typ, val, loc)
// }

func (ef *EnvFunctions) AddKeyInfo(proc *exec.WavmProcess, valAddr, valType, keyAddr, keyType, isArrayIndex uint64) {
	isArr := false
	if isArrayIndex > 0 {
		isArr = true
	}
	key := storage.StorageKey{
		KeyAddress:   keyAddr,
		KeyType:      int32(keyType),
		IsArrayIndex: isArr,
	}

	val := storage.StorageValue{
		ValueAddress: valAddr,
		ValueType:    int32(valType),
	}
	storageMap := ef.ctx.StorageMapping
	keySym := fmt.Sprintf("%d%d%t", key.KeyAddress, key.KeyType, key.IsArrayIndex)
	if _, ok := storageMap[valAddr]; ok {
		temp := storageMap[valAddr]
		if _, ok := temp.StorageKeyMap[keySym]; !ok {
			temp.StorageKey = append(temp.StorageKey, key)
			temp.StorageKeyMap[keySym] = true
			storageMap[valAddr] = temp
		}
	} else {
		temp := storage.StorageMapping{
			StorageValue:  val,
			StorageKey:    []storage.StorageKey{key},
			StorageKeyMap: map[string]bool{keySym: true},
		}
		storageMap[valAddr] = temp
	}
}

func callStateDb(ef *EnvFunctions, proc *exec.WavmProcess, valAddr uint64, stateDbOp func(val storage.StorageMapping, keyHash common.Hash)) {
	keyHash := common.BytesToHash(nil)
	storageMap := ef.ctx.StorageMapping
	if val, ok := storageMap[valAddr]; ok {
		for _, v := range val.StorageKey {
			log.Debug("env_funcs", "func", "callStateDb", "key", v, "keytype", v.KeyType)
			var lengthKeyHash common.Hash
			if v.IsArrayIndex {
				lengthKeyHash = keyHash
				log.Debug("callStateDb", "Is Array Index", "true")
			} else {
				log.Debug("callStateDb", "Is Array Index", "false")
			}
			keyMem := getMemory(proc, v.KeyAddress, v.KeyType, v.IsArrayIndex, getArrayLength(ef, lengthKeyHash))
			if (keyHash == common.Hash{}) {
				keyHash = utils.MapLocation(keyMem, nil)
			} else {
				keyHash = utils.MapLocation(keyHash.Bytes(), keyMem)
			}
			log.Debug("env_funcs", "func", "callStateDb", "keyhash", keyHash.String())
		}
		stateDbOp(val, keyHash)
	}
}

func getArrayLength(ef *EnvFunctions, lengthKeyHash common.Hash) uint64 {
	length := ef.ctx.StateDB.GetState(ef.ctx.Contract.Address(), lengthKeyHash).Bytes()
	log.Debug("env_funcs", "func", "getArrayLength", "hash", lengthKeyHash.String(), "len", length)
	return endianess.Uint64(length[len(length)-8:])
}

func inBounds(memoryData []byte, end uint64) {
	if len(memoryData) < int(end) {
		panic("error:out of memory bound")
	}
}

func getMemory(proc *exec.WavmProcess, addr uint64, addrType int32, isArrayIndex bool, length uint64) []byte {
	log.Debug("func", "getMemory", "")
	var mem []byte
	memoryData := proc.GetData()
	switch addrType {
	case abi.TY_INT32:
		inBounds(memoryData, addr+4)
		mem = memoryData[addr : addr+4]
		log.Debug("getMemory", "int32", int32(endianess.Uint32(mem)))
	case abi.TY_INT64:
		inBounds(memoryData, addr+8)
		mem = memoryData[addr : addr+8]
		log.Debug("getMemory", "int64", int64(endianess.Uint64(mem)))
	case abi.TY_UINT32:
		inBounds(memoryData, addr+4)
		mem = memoryData[addr : addr+4]
		log.Debug("getMemory", "uint32", endianess.Uint32(mem))
	case abi.TY_UINT64:
		inBounds(memoryData, addr+8)
		mem = memoryData[addr : addr+8]
		if isArrayIndex {
			index := endianess.Uint64(mem)
			log.Debug("getMemory", "array index", index, "length", length)
			if index+1 > length {
				panic(errExceededArray)
			}
		}
		log.Debug("getMemory", "uint64", endianess.Uint64(mem))
	case abi.TY_UINT256:
		inBounds(memoryData, addr+4)
		ptr := endianess.Uint32(memoryData[addr : addr+4])
		mem = []byte(readU256FromMemory(proc, uint64(ptr)).String())
		// mem = readU256FromMemory(proc, uint64(ptr)).Bytes()
		log.Debug("getMemory", "uint256", string(mem))
	case abi.TY_STRING:
		inBounds(memoryData, addr+4)
		ptr := endianess.Uint32(memoryData[addr : addr+4])
		mem = proc.ReadAt(uint64(ptr))
		log.Debug("getMemory", "string", string(mem))
	case abi.TY_ADDRESS:
		inBounds(memoryData, addr+4)
		ptr := endianess.Uint32(memoryData[addr : addr+4])
		mem = proc.ReadAt(uint64(ptr))
		log.Debug("getMemory", "address", common.BytesToAddress(mem).Hex())
	case abi.TY_BOOL:
		inBounds(memoryData, addr+4)
		mem = memoryData[addr : addr+4]
		log.Debug("getMemory", "bool", endianess.Uint32(mem))
	case abi.TY_POINTER:
		mem = make([]byte, 8)
		binary.BigEndian.PutUint64(mem, addr)
		log.Debug("getMemory", "pointer", mem)
	}
	return mem
}
func (ef *EnvFunctions) WriteWithPointer(proc *exec.WavmProcess, offsetAddr, baseAddr uint64) {
	log.Debug("instruction", "func", ">>>>>>>WriteWithPointer<<<<<<<")
	valAddr := offsetAddr + baseAddr
	storageMap := ef.ctx.StorageMapping
	if _, ok := storageMap[valAddr]; ok {
		ef.forbiddenMutable(proc)
	}
	op := func(val storage.StorageMapping, keyHash common.Hash) {

		valMem := getMemory(proc, val.StorageValue.ValueAddress, val.StorageValue.ValueType, false, 0)
		statedb := ef.ctx.StateDB
		contractAddr := ef.ctx.Contract.Address()
		if val.StorageValue.ValueType == abi.TY_STRING {
			n, s := utils.Split(valMem)
			statedb.SetState(contractAddr, keyHash, common.BigToHash(new(big.Int).SetInt64(int64(n))))
			ef.ctx.GasCounter.GasStore(statedb, contractAddr, keyHash, common.BigToHash(new(big.Int).SetInt64(int64(n))))
			for i := 1; i <= n; i++ {
				loc0 := new(big.Int).Add(keyHash.Big(), new(big.Int).SetInt64(int64(i)))
				statedb.SetState(contractAddr, common.BigToHash(loc0), common.BytesToHash(s[i-1]))
				ef.ctx.GasCounter.GasStore(statedb, contractAddr, common.BigToHash(loc0), common.BytesToHash(s[i-1]))
			}
		} else if val.StorageValue.ValueType == abi.TY_UINT256 {
			bigint := utils.GetU256(valMem)
			statedb.SetState(contractAddr, keyHash, common.BigToHash(bigint))
			ef.ctx.GasCounter.GasStore(statedb, contractAddr, keyHash, common.BytesToHash(valMem))
		} else {
			statedb.SetState(contractAddr, keyHash, common.BytesToHash(valMem))
			ef.ctx.GasCounter.GasStore(statedb, contractAddr, keyHash, common.BytesToHash(valMem))
		}
	}
	callStateDb(ef, proc, valAddr, op)
}

func (ef *EnvFunctions) ReadWithPointer(proc *exec.WavmProcess, offsetAddr, baseAddr uint64) {
	log.Debug("instruction", "func", ">>>>>>ReadWithPointer<<<<<<<<<")
	valAddr := offsetAddr + baseAddr
	op := func(val storage.StorageMapping, keyHash common.Hash) {
		stateVal := []byte{}
		statedb := ef.ctx.StateDB
		contractAddr := ef.ctx.Contract.Address()
		if val.StorageValue.ValueType == abi.TY_STRING {
			n := statedb.GetState(contractAddr, keyHash).Big().Int64()
			ef.ctx.GasCounter.GasLoad()
			for i := 1; i <= int(n); i++ {
				loc0 := new(big.Int).Add(keyHash.Big(), new(big.Int).SetInt64(int64(i)))
				val0 := statedb.GetState(contractAddr, common.BigToHash(loc0)).Big().Bytes()
				stateVal = append(stateVal, val0...)
				ef.ctx.GasCounter.GasLoad()
			}

		} else if val.StorageValue.ValueType == abi.TY_UINT256 {
			stateVal = []byte(statedb.GetState(contractAddr, keyHash).Big().String())
			ef.ctx.GasCounter.GasLoad()
		} else {
			stateVal = statedb.GetState(contractAddr, keyHash).Bytes()
			ef.ctx.GasCounter.GasLoad()
		}
		memoryData := proc.GetData()
		switch val.StorageValue.ValueType {
		case abi.TY_STRING:
			offset := proc.SetBytes(stateVal)
			endianess.PutUint32(memoryData[valAddr:], uint32(offset))
		case abi.TY_ADDRESS:
			offset := proc.SetBytes(stateVal[12:32])
			endianess.PutUint32(memoryData[valAddr:], uint32(offset))
		case abi.TY_UINT256:
			offset := proc.SetBytes(stateVal)
			endianess.PutUint32(memoryData[valAddr:], uint32(offset))
		case abi.TY_INT32, abi.TY_UINT32, abi.TY_BOOL:
			res := endianess.Uint32(stateVal[len(stateVal)-4:])
			endianess.PutUint32(memoryData[valAddr:], res)
		case abi.TY_INT64, abi.TY_UINT64:
			res := endianess.Uint64(stateVal[len(stateVal)-8:])
			endianess.PutUint64(memoryData[valAddr:], res)
		}

	}
	callStateDb(ef, proc, valAddr, op)
}

func (ef *EnvFunctions) InitializeVariables(proc *exec.WavmProcess) {
	// 普通类型初始化，忽略mapping和array
	log.Debug("EnvFunctions", "call", "InitializeVariables")
	//need to ignore array type because array init need array length
	storageMap := ef.ctx.StorageMapping
	for k, v := range storageMap {
		containArray := false
		for _, storageKey := range v.StorageKey {
			if storageKey.IsArrayIndex == true {
				containArray = true
				break
			}
		}
		if containArray == false {
			ef.WriteWithPointer(proc, k, 0)
		}
	}
}

func readU256FromMemory(proc *exec.WavmProcess, offset uint64) *big.Int {
	mem := proc.ReadAt(offset)
	return utils.GetU256(mem)
}

func (ef *EnvFunctions) U256FromU64(proc *exec.WavmProcess, x uint64) uint64 {
	bigint := new(big.Int)
	bigint.SetUint64(x)
	ef.ctx.GasCounter.GasFastestStep()
	return ef.returnU256(proc, bigint)
}
func (ef *EnvFunctions) U256FromI64(proc *exec.WavmProcess, x uint64) uint64 {
	bigint := new(big.Int)
	bigint.SetInt64(int64(x))
	ef.ctx.GasCounter.GasFastestStep()
	return ef.returnU256(proc, bigint)
}

func (ef *EnvFunctions) U256Add(proc *exec.WavmProcess, x, y uint64) uint64 {
	bigx := readU256FromMemory(proc, x)
	bigy := readU256FromMemory(proc, y)
	res := math.U256(bigy.Add(bigx, bigy))
	ef.ctx.GasCounter.GasFastestStep()
	return ef.returnU256(proc, res)
}

func (ef *EnvFunctions) U256Sub(proc *exec.WavmProcess, x, y uint64) uint64 {
	bigx := readU256FromMemory(proc, x)
	bigy := readU256FromMemory(proc, y)
	res := math.U256(bigy.Sub(bigx, bigy))
	ef.ctx.GasCounter.GasFastestStep()
	return ef.returnU256(proc, res)
}

func (ef *EnvFunctions) U256Mul(proc *exec.WavmProcess, x, y uint64) uint64 {
	bigx := readU256FromMemory(proc, x)
	bigy := readU256FromMemory(proc, y)
	res := math.U256(bigy.Mul(bigx, bigy))
	ef.ctx.GasCounter.GasFastestStep()
	return ef.returnU256(proc, res)
}

func (ef *EnvFunctions) U256Div(proc *exec.WavmProcess, x, y uint64) uint64 {
	bigx := readU256FromMemory(proc, x)
	bigy := readU256FromMemory(proc, y)
	res := new(big.Int)
	if bigy.Sign() != 0 {
		res = math.U256(bigy.Div(bigx, bigy))
	} else {
		res = bigy.SetUint64(0)
	}
	ef.ctx.GasCounter.GasFastestStep()
	return ef.returnU256(proc, res)
}

func (ef *EnvFunctions) U256Mod(proc *exec.WavmProcess, x, y uint64) uint64 {
	bigx := readU256FromMemory(proc, x)
	bigy := readU256FromMemory(proc, y)
	res := new(big.Int)
	if bigy.Sign() != 0 {
		res = math.U256(bigy.Mod(bigx, bigy))
	} else {
		res = bigy.SetUint64(0)
	}
	ef.ctx.GasCounter.GasFastestStep()
	return ef.returnU256(proc, res)
}

func (ef *EnvFunctions) U256Pow(proc *exec.WavmProcess, base, exponent uint64) uint64 {
	b := readU256FromMemory(proc, base)
	e := readU256FromMemory(proc, exponent)
	res := math.Exp(b, e)
	ef.ctx.GasCounter.GasPow(e)
	return ef.returnU256(proc, res)
}

func (ef *EnvFunctions) U256Cmp(proc *exec.WavmProcess, x, y uint64) uint64 {
	bigx := readU256FromMemory(proc, x)
	bigy := readU256FromMemory(proc, y)
	res := bigx.Cmp(bigy)
	// ef.ctx.GasCounter.GasPow(e)
	ef.ctx.GasCounter.GasFastestStep()
	return uint64(res)
}

func (ef *EnvFunctions) Pow(proc *exec.WavmProcess, base, exponent uint64) uint64 {
	b := new(big.Int)
	b.SetUint64(base)
	e := new(big.Int)
	e.SetUint64(exponent)
	res := math.Exp(b, e)
	ef.ctx.GasCounter.GasPow(e)
	return res.Uint64()
}

func (ef *EnvFunctions) AddGas(proc *exec.WavmProcess, cost uint64) {
	ef.ctx.GasCounter.AdjustedCharge(cost)
}

//todo 考虑revert的完整实现 contractcall里需要用到revert
func (ef *EnvFunctions) Revert(proc *exec.WavmProcess, msgIdx uint64) {
	ctx := ef.ctx
	ctx.GasCounter.GasRevert()
	msg := proc.ReadAt(msgIdx)
	ctx.GasCounter.GasMemoryCost(uint64(len(msg)))
	log.Info("Contract Revert >>>>", "message", string(msg))
	panic(errormsg.ErrExecutionReverted)
}

func (ef *EnvFunctions) returnPointer(proc *exec.WavmProcess, input []byte) uint64 {
	ctx := ef.ctx
	ctx.GasCounter.GasReturnPointer(uint64(len(input)))
	return uint64(proc.SetBytes(input))
}

func (ef *EnvFunctions) returnAddress(proc *exec.WavmProcess, input []byte) uint64 {
	ctx := ef.ctx
	ctx.GasCounter.GasReturnAddress()
	return uint64(proc.SetBytes(input))
}

func (ef *EnvFunctions) returnU256(proc *exec.WavmProcess, bigint *big.Int) uint64 {
	ctx := ef.ctx
	ctx.GasCounter.GasReturnU256()
	return uint64(proc.SetBytes([]byte(bigint.String())))
}

func (ef *EnvFunctions) returnHash(proc *exec.WavmProcess, hash []byte) uint64 {
	ctx := ef.ctx
	ctx.GasCounter.GasReturnHash()
	return uint64(proc.SetBytes(hash))
}

func (ef *EnvFunctions) Sender(proc *exec.WavmProcess, ptr uint64) {
	sender := ef.ctx.Contract.Address().Bytes()
	log.Debug("EnvFunctions", "func", "Sender", "ptr", ptr, "data", sender)
	proc.WriteAt(sender, int64(ptr))
}

func (ef *EnvFunctions) Load(proc *exec.WavmProcess, keyptr uint64, dataptr uint64) uint64 {
	log.Debug("EnvFunctions", "func", "Load")
	keyData := ef.getQString(proc, keyptr)
	log.Debug("EnvFunctions", "func", "Load", "key data", keyData, "data ptr", dataptr)
	keyHash := common.BytesToHash(keyData)
	statedb := ef.ctx.StateDB
	contractAddr := ef.ctx.Contract.Address()
	n := statedb.GetState(contractAddr, keyHash).Big().Int64()
	stateVal := []byte{}
	for i := 1; i <= int(n); i++ {
		loc0 := new(big.Int).Add(keyHash.Big(), new(big.Int).SetInt64(int64(i)))
		val0 := statedb.GetState(contractAddr, common.BigToHash(loc0)).Big().Bytes()
		stateVal = append(stateVal, val0...)
	}
	log.Debug("EnvFunctions", "func", "Load", "value data", stateVal, "size", len(stateVal))
	proc.WriteAt(stateVal, int64(dataptr))
	return uint64(len(stateVal))
}

func (ef *EnvFunctions) Store(proc *exec.WavmProcess, keyptr uint64, dataptr uint64) {
	log.Debug("EnvFunctions", "func", "Store")
	keyData := ef.getQString(proc, keyptr)
	keyHash := common.BytesToHash(keyData)
	valueData := ef.getQString(proc, dataptr)
	log.Debug("EnvFunctions", "func", "Store", "key ptr", keyptr, "key data", keyData, "value ptr", dataptr, "value data", valueData)
	statedb := ef.ctx.StateDB
	contractAddr := ef.ctx.Contract.Address()
	n, s := utils.Split(valueData)
	statedb.SetState(contractAddr, keyHash, common.BigToHash(new(big.Int).SetInt64(int64(n))))
	for i := 1; i <= n; i++ {
		loc0 := new(big.Int).Add(keyHash.Big(), new(big.Int).SetInt64(int64(i)))
		statedb.SetState(contractAddr, common.BigToHash(loc0), common.BytesToHash(s[i-1]))
	}
}

func (ef *EnvFunctions) getQString(proc *exec.WavmProcess, strPtr uint64) []byte {
	size := endianess.Uint32(proc.GetData()[strPtr : strPtr+4])
	offset := endianess.Uint32(proc.GetData()[strPtr+4 : strPtr+8])
	strData := proc.GetData()[offset : offset+size]
	log.Debug("EnvFunctions", "func", "getQString", "str_ptr", strPtr, "size", size, "offset", offset, "string data", strData)
	return strData
}

func (ef *EnvFunctions) forbiddenMutable(proc *exec.WavmProcess) {
	if proc.Mutable() == false {
		log.Debug("ForbiddenMutable", "msg", "this function is not a mutable function")
		err := errors.New("Mutable Forbidden: This function is not a mutable function")
		panic(err)
	}
}
