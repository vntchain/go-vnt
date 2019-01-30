package gas

import (
	"errors"
	"math/big"

	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/common/math"
	"github.com/vntchain/go-vnt/core/state"
	"github.com/vntchain/go-vnt/core/vm"
	"github.com/vntchain/go-vnt/core/wavm/contract"
	"github.com/vntchain/go-vnt/params"
	ops "github.com/vntchain/vnt-wasm/wasm/operators"
)

const (
	WasmCostsRegular       = 1
	WasmCostsDiv           = 16
	WasmCostsMul           = 4
	WasmCostsMem           = 2
	WasmCostsStaticU256    = 64
	WasmCostsStaticHash    = 64
	WasmCostsStaticAddress = 40
	/// Memory stipend. Amount of free memory (in 64kb pages) each contract can use for stack.
	WasmCostsInitialMem = 4096
	/// Grow memory cost, per page (64kb)
	WasmCostsGrowMem        = 8192
	WasmCostsMemcpy         = 1
	WasmCostsMaxStackHeight = 64 * 1024
	WasmCostsOpcodesMul     = 3
	WasmCostsOpcodesDiv     = 8
)

const ErrorGasLimit = "Invocation resulted in gas limit violated"
const ErrorDisableFloatingPoint = "Wasm contract error: disabled floating point"

var errGasUintOverflow = errors.New("gas uint64 overflow")

type GasValue struct {
	Metering Metering
	Value    uint64
}

type Gas struct {
	Ops   map[byte]InstructionType
	Rules map[InstructionType]GasValue
}

type GasCounter struct {
	Contract *contract.WASMContract
	GasTable params.GasTable
}

func NewGas(disableFloatingPoint bool) Gas {
	rules := Gas{
		Ops: map[byte]InstructionType{
			ops.Unreachable:  InstructionTypeUnreachable,
			ops.Nop:          InstructionTypeNop,
			ops.Block:        InstructionTypeControlFlow,
			ops.Loop:         InstructionTypeControlFlow,
			ops.If:           InstructionTypeControlFlow,
			ops.Else:         InstructionTypeControlFlow,
			ops.End:          InstructionTypeControlFlow,
			ops.Br:           InstructionTypeControlFlow,
			ops.BrIf:         InstructionTypeControlFlow,
			ops.BrTable:      InstructionTypeControlFlow,
			ops.Return:       InstructionTypeControlFlow,
			ops.Call:         InstructionTypeControlFlow,
			ops.CallIndirect: InstructionTypeControlFlow,
			ops.Drop:         InstructionTypeControlFlow,
			ops.Select:       InstructionTypeControlFlow,

			ops.GetLocal:  InstructionTypeLocal,
			ops.SetLocal:  InstructionTypeLocal,
			ops.TeeLocal:  InstructionTypeLocal,
			ops.GetGlobal: InstructionTypeLocal,
			ops.SetGlobal: InstructionTypeLocal,

			ops.I32Load:    InstructionTypeLoad,
			ops.I64Load:    InstructionTypeLoad,
			ops.F32Load:    InstructionTypeLoad,
			ops.F64Load:    InstructionTypeLoad,
			ops.I32Load8s:  InstructionTypeLoad,
			ops.I32Load8u:  InstructionTypeLoad,
			ops.I32Load16s: InstructionTypeLoad,
			ops.I32Load16u: InstructionTypeLoad,
			ops.I64Load8s:  InstructionTypeLoad,
			ops.I64Load8u:  InstructionTypeLoad,
			ops.I64Load16s: InstructionTypeLoad,
			ops.I64Load16u: InstructionTypeLoad,
			ops.I64Load32s: InstructionTypeLoad,
			ops.I64Load32u: InstructionTypeLoad,

			ops.I32Store:   InstructionTypeStore,
			ops.I64Store:   InstructionTypeStore,
			ops.F32Store:   InstructionTypeStore,
			ops.F64Store:   InstructionTypeStore,
			ops.I32Store8:  InstructionTypeStore,
			ops.I32Store16: InstructionTypeStore,
			ops.I64Store8:  InstructionTypeStore,
			ops.I64Store16: InstructionTypeStore,
			ops.I64Store32: InstructionTypeStore,

			ops.CurrentMemory: InstructionTypeCurrentMemory,
			ops.GrowMemory:    InstructionTypeGrowMemory,

			ops.I32Const: InstructionTypeConst,
			ops.I64Const: InstructionTypeConst,

			ops.F32Const: InstructionTypeFloatConst,
			ops.F64Const: InstructionTypeFloatConst,

			ops.I32Eqz: InstructionTypeIntegerComparsion,
			ops.I32Eq:  InstructionTypeIntegerComparsion,
			ops.I32Ne:  InstructionTypeIntegerComparsion,
			ops.I32LtS: InstructionTypeIntegerComparsion,
			ops.I32LtU: InstructionTypeIntegerComparsion,
			ops.I32GtS: InstructionTypeIntegerComparsion,
			ops.I32GtU: InstructionTypeIntegerComparsion,
			ops.I32LeS: InstructionTypeIntegerComparsion,
			ops.I32LeU: InstructionTypeIntegerComparsion,
			ops.I32GeS: InstructionTypeIntegerComparsion,
			ops.I32GeU: InstructionTypeIntegerComparsion,

			ops.I64Eqz: InstructionTypeIntegerComparsion,
			ops.I64Eq:  InstructionTypeIntegerComparsion,
			ops.I64Ne:  InstructionTypeIntegerComparsion,
			ops.I64LtS: InstructionTypeIntegerComparsion,
			ops.I64LtU: InstructionTypeIntegerComparsion,
			ops.I64GtS: InstructionTypeIntegerComparsion,
			ops.I64GtU: InstructionTypeIntegerComparsion,
			ops.I64LeS: InstructionTypeIntegerComparsion,
			ops.I64LeU: InstructionTypeIntegerComparsion,
			ops.I64GeS: InstructionTypeIntegerComparsion,
			ops.I64GeU: InstructionTypeIntegerComparsion,

			ops.F32Eq: InstructionTypeFloatComparsion,
			ops.F32Ne: InstructionTypeFloatComparsion,
			ops.F32Lt: InstructionTypeFloatComparsion,
			ops.F32Gt: InstructionTypeFloatComparsion,
			ops.F32Le: InstructionTypeFloatComparsion,
			ops.F32Ge: InstructionTypeFloatComparsion,

			ops.F64Eq: InstructionTypeFloatComparsion,
			ops.F64Ne: InstructionTypeFloatComparsion,
			ops.F64Lt: InstructionTypeFloatComparsion,
			ops.F64Gt: InstructionTypeFloatComparsion,
			ops.F64Le: InstructionTypeFloatComparsion,
			ops.F64Ge: InstructionTypeFloatComparsion,

			ops.I32Clz:    InstructionTypeBit,
			ops.I32Ctz:    InstructionTypeBit,
			ops.I32Popcnt: InstructionTypeBit,
			ops.I32Add:    InstructionTypeAdd,
			ops.I32Sub:    InstructionTypeAdd,
			ops.I32Mul:    InstructionTypeMul,
			ops.I32DivS:   InstructionTypeDiv,
			ops.I32DivU:   InstructionTypeDiv,
			ops.I32RemS:   InstructionTypeDiv,
			ops.I32RemU:   InstructionTypeDiv,
			ops.I32And:    InstructionTypeBit,
			ops.I32Or:     InstructionTypeBit,
			ops.I32Xor:    InstructionTypeBit,
			ops.I32Shl:    InstructionTypeBit,
			ops.I32ShrS:   InstructionTypeBit,
			ops.I32ShrU:   InstructionTypeBit,
			ops.I32Rotl:   InstructionTypeBit,
			ops.I32Rotr:   InstructionTypeBit,

			ops.I64Clz:    InstructionTypeBit,
			ops.I64Ctz:    InstructionTypeBit,
			ops.I64Popcnt: InstructionTypeBit,
			ops.I64Add:    InstructionTypeAdd,
			ops.I64Sub:    InstructionTypeAdd,
			ops.I64Mul:    InstructionTypeMul,
			ops.I64DivS:   InstructionTypeDiv,
			ops.I64DivU:   InstructionTypeDiv,
			ops.I64RemS:   InstructionTypeDiv,
			ops.I64RemU:   InstructionTypeDiv,
			ops.I64And:    InstructionTypeBit,
			ops.I64Or:     InstructionTypeBit,
			ops.I64Xor:    InstructionTypeBit,
			ops.I64Shl:    InstructionTypeBit,
			ops.I64ShrS:   InstructionTypeBit,
			ops.I64ShrU:   InstructionTypeBit,
			ops.I64Rotl:   InstructionTypeBit,
			ops.I64Rotr:   InstructionTypeBit,

			ops.F32Abs:      InstructionTypeFloat,
			ops.F32Neg:      InstructionTypeFloat,
			ops.F32Ceil:     InstructionTypeFloat,
			ops.F32Floor:    InstructionTypeFloat,
			ops.F32Trunc:    InstructionTypeFloat,
			ops.F32Nearest:  InstructionTypeFloat,
			ops.F32Sqrt:     InstructionTypeFloat,
			ops.F32Add:      InstructionTypeFloat,
			ops.F32Sub:      InstructionTypeFloat,
			ops.F32Mul:      InstructionTypeFloat,
			ops.F32Div:      InstructionTypeFloat,
			ops.F32Min:      InstructionTypeFloat,
			ops.F32Max:      InstructionTypeFloat,
			ops.F32Copysign: InstructionTypeFloat,
			ops.F64Abs:      InstructionTypeFloat,
			ops.F64Neg:      InstructionTypeFloat,
			ops.F64Ceil:     InstructionTypeFloat,
			ops.F64Floor:    InstructionTypeFloat,
			ops.F64Trunc:    InstructionTypeFloat,
			ops.F64Nearest:  InstructionTypeFloat,
			ops.F64Sqrt:     InstructionTypeFloat,
			ops.F64Add:      InstructionTypeFloat,
			ops.F64Sub:      InstructionTypeFloat,
			ops.F64Mul:      InstructionTypeFloat,
			ops.F64Div:      InstructionTypeFloat,
			ops.F64Min:      InstructionTypeFloat,
			ops.F64Max:      InstructionTypeFloat,
			ops.F64Copysign: InstructionTypeFloat,

			ops.I32WrapI64:    InstructionTypeConversion,
			ops.I64ExtendSI32: InstructionTypeConversion,
			ops.I64ExtendUI32: InstructionTypeConversion,

			ops.I32TruncSF32:   InstructionTypeFloatConversion,
			ops.I32TruncUF32:   InstructionTypeFloatConversion,
			ops.I32TruncSF64:   InstructionTypeFloatConversion,
			ops.I32TruncUF64:   InstructionTypeFloatConversion,
			ops.I64TruncSF32:   InstructionTypeFloatConversion,
			ops.I64TruncUF32:   InstructionTypeFloatConversion,
			ops.I64TruncSF64:   InstructionTypeFloatConversion,
			ops.I64TruncUF64:   InstructionTypeFloatConversion,
			ops.F32ConvertSI32: InstructionTypeFloatConversion,
			ops.F32ConvertUI32: InstructionTypeFloatConversion,
			ops.F32ConvertSI64: InstructionTypeFloatConversion,
			ops.F32ConvertUI64: InstructionTypeFloatConversion,
			ops.F32DemoteF64:   InstructionTypeFloatConversion,
			ops.F64ConvertSI32: InstructionTypeFloatConversion,
			ops.F64ConvertUI32: InstructionTypeFloatConversion,
			ops.F64ConvertSI64: InstructionTypeFloatConversion,
			ops.F64ConvertUI64: InstructionTypeFloatConversion,
			ops.F64PromoteF32:  InstructionTypeFloatConversion,

			ops.I32ReinterpretF32: InstructionTypeReinterpretation,
			ops.I64ReinterpretF64: InstructionTypeReinterpretation,
			ops.F32ReinterpretI32: InstructionTypeReinterpretation,
			ops.F64ReinterpretI64: InstructionTypeReinterpretation,
		},
		Rules: map[InstructionType]GasValue{
			InstructionTypeLoad:  GasValue{Metering: MeteringFixed, Value: WasmCostsMem},
			InstructionTypeStore: GasValue{Metering: MeteringFixed, Value: WasmCostsMem},
			InstructionTypeDiv:   GasValue{Metering: MeteringFixed, Value: WasmCostsDiv},
			InstructionTypeMul:   GasValue{Metering: MeteringFixed, Value: WasmCostsMul},
		},
	}
	if disableFloatingPoint {
		rules.Rules[InstructionTypeFloat] = GasValue{
			Metering: MeteringForbidden, Value: 0,
		}
		rules.Rules[InstructionTypeFloatComparsion] = GasValue{
			Metering: MeteringForbidden, Value: 0,
		}
		rules.Rules[InstructionTypeFloatConst] = GasValue{
			Metering: MeteringForbidden, Value: 0,
		}
		rules.Rules[InstructionTypeFloatConversion] = GasValue{
			Metering: MeteringForbidden, Value: 0,
		}
	}
	return rules
}

func (gas Gas) GasCost(op byte) uint64 {
	metering := gas.Rules[gas.Ops[op]]
	switch metering.Metering {
	case MeteringForbidden:
		panic(ErrorDisableFloatingPoint)
	case MeteringFixed:
		return metering.Value
	default:
		return WasmCostsRegular
	}
}

func constGasFunc(gas uint64) uint64 {
	return gas
}

func NewGasCounter(contract *contract.WASMContract, gasTable params.GasTable) GasCounter {
	return GasCounter{
		Contract: contract,
		GasTable: gasTable,
	}

}

func (gas GasCounter) Charge(amount uint64) {
	if !gas.ChargeGas(amount) {
		panic(ErrorGasLimit)
	}
}

func (gas GasCounter) ChargeGas(amount uint64) bool {
	return gas.Contract.UseGas(amount)
}

func (gas GasCounter) AdjustedCharge(amount uint64) {
	// gas.Charge(amount * WasmCostsOpcodesMul / WasmCostsOpcodesDiv)
	gas.Charge(amount)
}

// memoryGasCosts calculates the quadratic gas for memory expansion. It does so
// only for the memory region that is expanded, not the total memory.
// func memoryGasCost(mem *wasm.WavmMemory, newMemSize uint64) (uint64, error) {

// 	if newMemSize == 0 {
// 		return 0, nil
// 	}
// 	// The maximum that will fit in a uint64 is max_word_count - 1
// 	// anything above that will result in an overflow.
// 	// Additionally, a newMemSize which results in a
// 	// newMemSizeWords larger than 0x7ffffffff will cause the square operation
// 	// to overflow.
// 	// The constant 0xffffffffe0 is the highest number that can be used without
// 	// overflowing the gas calculation
// 	if newMemSize > 0xffffffffe0 {
// 		return 0, errGasUintOverflow
// 	}

// 	newMemSizeWords := toWordSize(newMemSize)
// 	newMemSize = newMemSizeWords * 32

// 	if newMemSize > uint64(mem.Len()) {
// 		square := newMemSizeWords * newMemSizeWords
// 		linCoef := newMemSizeWords * params.MemoryGas
// 		quadCoef := square / params.QuadCoeffDiv
// 		newTotalFee := linCoef + quadCoef

// 		fee := newTotalFee - mem.lastGasCost
// 		mem.lastGasCost = newTotalFee

// 		return fee, nil
// 	}
// 	return 0, nil
// }

// func (vm *VM) GasGetBalanceFromAddress() error {
// 	return vm.adjustedCharge(constGasFunc(20))
// }

// func (vm *VM) GasStorageWrite() error {
// 	idx := big.NewInt(int64(vm.ctx.locals[0])).Uint64()
// 	loc := common.BytesToHash(vm.Memory().GetPtr(idx))
// 	validx := big.NewInt(int64(vm.ctx.locals[1])).Uint64()
// 	val := common.BytesToHash(vm.Memory().GetPtr(validx))
// 	formerVal := vm.StateDB.GetState(vm.Contract.Address(), loc)
// 	var err error
// 	if common.EmptyHash(formerVal) && !common.EmptyHash(val) {
// 		log.Debug("---------", "set", "-------------")
// 		err = vm.adjustedCharge(params.SstoreSetGas)
// 	} else {
// 		log.Debug("---------", "reset", "-------------")
// 		err = vm.adjustedCharge(params.SstoreResetGas)
// 	}
// 	return err
// }

// func (vm *VM) GasStorageRead() error {
// 	return vm.adjustedCharge(params.SloadGas)
// }

func (gas GasCounter) GasQuickStep() {
	gas.Charge(constGasFunc(vm.GasQuickStep))
}

func (gas GasCounter) GasFastestStep() {
	gas.Charge(constGasFunc(vm.GasFastestStep))
}

func (gas GasCounter) GasGetBlockNumber() {
	gas.Charge(constGasFunc(vm.GasQuickStep))
}

func (gas GasCounter) GasGetBalanceFromAddress() {
	gas.Charge(constGasFunc(vm.GasQuickStep))
}

func (gas GasCounter) GasMemoryCost(size uint64) {
	gas.Charge(constGasFunc(vm.GasQuickStep * size))
}

func (gas GasCounter) GasGetGas() {
	gas.Charge(constGasFunc(vm.GasQuickStep))
}

func (gas GasCounter) GasGetBlockHash() {
	gas.Charge(constGasFunc(vm.GasExtStep))
}

func (gas GasCounter) GasGetBlockProduser() {
	gas.Charge(constGasFunc(vm.GasQuickStep))
}

func (gas GasCounter) GasGetTimestamp() {
	gas.Charge(constGasFunc(vm.GasQuickStep))
}

func (gas GasCounter) GasGetOrigin() {
	gas.Charge(constGasFunc(vm.GasQuickStep))
}

func (gas GasCounter) GasGetSender() {
	gas.Charge(constGasFunc(vm.GasQuickStep))
}

func (gas GasCounter) GasGetGasLimit() {
	gas.Charge(constGasFunc(vm.GasQuickStep))
}

// func (vm *VM) GasGenerateKey() error {
// 	return vm.adjustedCharge(constGasFunc(20))
// }

// func (vm *VM) GasGetStorageCount() error {
// 	return nil
// }

func (gas GasCounter) GasGetValue() {
	gas.Charge(constGasFunc(vm.GasQuickStep))
}

//todo 验证sha3 gas消耗，和evm保持一致
func (gas GasCounter) GasSHA3(size uint64) {
	gas.Charge(params.Sha3Gas)
	gas.Charge(params.Sha3WordGas * size)
}

func (gas GasCounter) GasGetContractAddress() {
	gas.Charge(constGasFunc(vm.GasQuickStep))
}

func (gas GasCounter) GasAssert() {
	gas.Charge(constGasFunc(vm.GasQuickStep))
}

func (gas GasCounter) GasRevert() {
	gas.Charge(constGasFunc(vm.GasQuickStep))
}

func (gas GasCounter) GasSendFromContract() {
	gas.Charge(constGasFunc(params.CallStipend))
}

// func (gas GasCounter) GasGetContractValue() {
// 	gas.Charge(constGasFunc(params.CallValueTransferGas))
// }

func (gas GasCounter) GasFromI64() {
	gas.Charge(constGasFunc(vm.GasQuickStep))
}

func (gas GasCounter) GasFromU64() {
	gas.Charge(constGasFunc(vm.GasQuickStep))
}

func (gas GasCounter) GasToI64() {
	gas.Charge(constGasFunc(vm.GasQuickStep))
}

func (gas GasCounter) GasToU64() {
	gas.Charge(constGasFunc(vm.GasQuickStep))
}

func (gas GasCounter) GasConcat(size uint64) {
	gas.Charge(constGasFunc(WasmCostsMem * size))
}

func (gas GasCounter) GasEqual() {
	gas.Charge(constGasFunc(vm.GasQuickStep))
}

func (gas GasCounter) GasLog(size uint64, topics uint64) {
	requestedSize, overflow := vm.BigUint64(new(big.Int).SetUint64(size))
	if overflow {
		panic(errGasUintOverflow)
	}
	costgas := uint64(0)
	if costgas, overflow = math.SafeAdd(costgas, params.LogGas); overflow {
		panic(errGasUintOverflow)
	}
	if costgas, overflow = math.SafeAdd(costgas, topics*params.LogTopicGas); overflow {
		panic(errGasUintOverflow)
	}

	var memorySizeGas uint64
	if memorySizeGas, overflow = math.SafeMul(requestedSize, params.LogDataGas); overflow {
		panic(errGasUintOverflow)
	}
	if costgas, overflow = math.SafeAdd(costgas, memorySizeGas); overflow {
		panic(errGasUintOverflow)
	}
	gas.Charge(costgas)
}

func (gas GasCounter) GasCall(address common.Address, value, gasLimit, blockNumber *big.Int, chainConfig *params.ChainConfig, statedb *state.StateDB) uint64 {
	var (
		callgas        = gas.GasTable.Calls
		transfersValue = value.Sign() != 0
		eip158         = chainConfig.IsEIP158(blockNumber)
	)
	if eip158 {
		if transfersValue && statedb.Empty(address) {
			callgas += params.CallNewAccountGas
		}
	} else if !statedb.Exist(address) {
		callgas += params.CallNewAccountGas
	}
	if transfersValue {
		callgas += params.CallValueTransferGas
	}
	tempgas, err := gas.callGas(gas.GasTable, gas.Contract.Gas, callgas, gasLimit)
	if err != nil {
		panic(err.Error())
	}

	var overflow bool
	if callgas, overflow = math.SafeAdd(callgas, tempgas); overflow {
		panic(errGasUintOverflow)
	}
	gas.Charge(callgas)
	return tempgas
}

func (gas GasCounter) callGas(gasTable params.GasTable, availableGas, base uint64, callCost *big.Int) (uint64, error) {
	if gasTable.CreateBySuicide > 0 {
		availableGas = availableGas - base
		gas := availableGas - availableGas/64
		// If the bit length exceeds 64 bit we know that the newly calculated "gas" for EIP150
		// is smaller than the requested amount. Therefor we return the new gas instead
		// of returning an error.
		if callCost.BitLen() > 64 || gas < callCost.Uint64() {
			return gas, nil
		}
	}
	if callCost.BitLen() > 64 {
		return 0, errGasUintOverflow
	}
	return callCost.Uint64(), nil
}

func (gas GasCounter) GasStore(stateDb *state.StateDB, contractAddr common.Address, loc common.Hash, value common.Hash) {
	var (
		y, x = value, loc
		val  = stateDb.GetState(contractAddr, x)
	)
	// This checks for 3 scenario's and calculates gas accordingly
	// 1. From a zero-value address to a non-zero value         (NEW VALUE)
	// 2. From a non-zero value address to a zero-value address (DELETE)
	// 3. From a non-zero to a non-zero                         (CHANGE)
	if (val == common.Hash{} && y != common.Hash{}) {
		// 0 => non 0
		gas.Charge(constGasFunc(params.SstoreSetGas))
		// return params.SstoreSetGas, nil
	} else if (val != common.Hash{} && y == common.Hash{}) {
		stateDb.AddRefund(params.SstoreRefundGas)
		gas.Charge(constGasFunc(params.SstoreClearGas))
		// return params.SstoreClearGas, nil
	} else {
		// non 0 => non 0 (or 0 => 0)
		gas.Charge(constGasFunc(params.SstoreResetGas))
		// return params.SstoreResetGas, nil
	}
}

func (gas GasCounter) GasLoad() {
	gas.Charge(constGasFunc(gas.GasTable.SLoad))
}

func (gas GasCounter) GasPow(exponent *big.Int) {
	expByteLen := uint64((exponent.BitLen() + 7) / 8)
	var (
		costgas  = expByteLen * gas.GasTable.ExpByte // no overflow check required. Max is 256 * ExpByte gas
		overflow bool
	)
	if costgas, overflow = math.SafeAdd(costgas, vm.GasQuickStep); overflow {
		panic(errGasUintOverflow)
	}
	gas.Charge(constGasFunc(costgas))
}

func (gas GasCounter) GasReturnAddress() {
	gas.AdjustedCharge(constGasFunc(WasmCostsStaticAddress))
}

func (gas GasCounter) GasReturnU256() {
	gas.AdjustedCharge(constGasFunc(WasmCostsStaticU256))
}

func (gas GasCounter) GasReturnHash() {
	gas.AdjustedCharge(constGasFunc(WasmCostsStaticHash))
}

func (gas GasCounter) GasReturnPointer(size uint64) {
	gas.AdjustedCharge(constGasFunc(size))
}

func (gas GasCounter) GasInitialMemory(initial uint64) {
	gas.AdjustedCharge(constGasFunc(initial * WasmCostsInitialMem))
}
