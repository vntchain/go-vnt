package tests

import (
	"path/filepath"
	"testing"
)

var (
	erc20Code = filepath.Join(basepath, "erc20/TokenERC20.compress")
	// erc20Code = filepath.Join(basepath, "qlang/PERC20.wasm")
	erc20Abi = filepath.Join(basepath, "erc20/abi.json")
)

// func newErc20VM(iscreated bool) *wavm.Wavm {
// 	return newWavm(erc20Code, erc20Abi, iscreated)
// }
// func TestERC20(t *testing.T) {
// 	initialSupply := new(big.Int)
// 	initialSupply.SetString("1000000000000", 10)
// 	tokenName := "bitcoin"
// 	tokenSymbol := "BTC"
// 	amount1 := new(big.Int)
// 	amount1.SetString("10000000", 10)
// 	//to1 := common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
// 	to1 := common.HexToAddress("0x02")
// 	vm := newErc20VM(true)
// 	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
// 	testCreate := func(initialSupply *big.Int, tokenName, tokenSymbol string) {
// 		funcName := ""
// 		res, err := vm.Apply(pack(vm, funcName, initialSupply, tokenName, tokenSymbol), nil, mutable)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		t.Logf("create contract %s", res)
// 	}
// 	testTransfer := func(to common.Address, value *big.Int) {
// 		vm.ChainContext.IsCreated = false
// 		funcName := "transfer"
// 		res, err := vm.Apply(pack(vm, funcName, to, value), nil, mutable)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		var flag bool
// 		unPack(vm, &flag, funcName, res)
// 		if flag != true {
// 			t.Errorf("unexpected value : want %t, got %t", true, flag)
// 		} else {
// 			t.Logf("%s success get %t", funcName, flag)
// 		}
// 		type TransferEvent struct {
// 			From  common.Address
// 			To    common.Address
// 			Value *big.Int
// 		}
// 		var ev TransferEvent
// 		logs := vm.ChainContext.StateDB.GetLogs(common.HexToHash("0x1111"))
// 		res, _ = logs[0].MarshalJSON()
// 		t.Logf("success get event json %s", res)
// 		unPack(vm, &ev, "Transfer", logs[0].Data)
// 		t.Logf("success get event %v", ev)
// 	}

// 	testGetTokenName := func() {
// 		vm.ChainContext.IsCreated = false
// 		funcName := "GetTokenName"
// 		res, err := vm.Apply(pack(vm, funcName), nil, mutable)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		var name string
// 		unPack(vm, &name, funcName, res)
// 		if name != tokenName {
// 			t.Errorf("unexpected value : want %s, got %s", tokenName, name)
// 		} else {
// 			t.Logf("%s success get %s", funcName, name)
// 		}
// 	}

// 	testGetTotalSupply := func() {
// 		vm.ChainContext.IsCreated = false
// 		funcName := "GetTotalSupply"
// 		res, err := vm.Apply(pack(vm, funcName), nil, mutable)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		var supply *big.Int
// 		unPack(vm, &supply, funcName, res)
// 		a := new(big.Int)
// 		if supply.Cmp(a.Mul(initialSupply, new(big.Int).SetUint64(100000000))) != 0 {
// 			t.Errorf("unexpected value : want %d, got %d", a, supply)
// 		} else {
// 			t.Logf("%s success get %d", funcName, supply)
// 		}
// 	}

// 	testGetSymbol := func() {
// 		vm.ChainContext.IsCreated = false
// 		funcName := "GetSymbol"
// 		res, err := vm.Apply(pack(vm, funcName), nil, mutable)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		var sym string
// 		unPack(vm, &sym, funcName, res)
// 		if sym != tokenSymbol {
// 			t.Errorf("unexpected value : want %s, got %s", tokenSymbol, sym)
// 		} else {
// 			t.Logf("%s success get %s", funcName, sym)
// 		}
// 	}

// 	testGetDecimals := func() {
// 		vm.ChainContext.IsCreated = false
// 		funcName := "GetDecimals"
// 		res, err := vm.Apply(pack(vm, funcName), nil, mutable)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		var deci *big.Int
// 		unPack(vm, &deci, funcName, res)
// 		if deci.Uint64() != 8 {
// 			t.Errorf("unexpected value : want %d, got %d", 8, deci)
// 		} else {
// 			t.Logf("%s success get %d", funcName, deci)
// 		}
// 	}

// 	testGetAmount := func(addr common.Address) *big.Int {
// 		vm.ChainContext.IsCreated = false
// 		funcName := "GetAmount"
// 		res, err := vm.Apply(pack(vm, funcName, addr), nil, mutable)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		var amount *big.Int
// 		unPack(vm, &amount, funcName, res)
// 		return amount
// 		// if amount != 4 {
// 		// 	t.Errorf("unexpected value : want %d, got %d", 4, amount)
// 		// } else {
// 		// 	t.Logf("%s success get %d", funcName, amount)
// 		// }
// 	}

// 	testCreate(initialSupply, tokenName, tokenSymbol)
// 	amount := testGetAmount(caller)
// 	a := new(big.Int)
// 	if amount.Cmp(a.Mul(initialSupply, new(big.Int).SetUint64(100000000))) != 0 {
// 		t.Errorf("unexpected value : want %d, got %d", a, amount)
// 	} else {
// 		t.Logf("%s success get %d", "GetAmount", amount)
// 	}
// 	now := T.Now()
// 	testTransfer(to1, amount1)
// 	after := T.Since(now)
// 	t.Logf("time %s", after.String())

// 	testGetTokenName()
// 	testGetTotalSupply()
// 	testGetSymbol()
// 	testGetDecimals()
// 	amount = testGetAmount(to1)
// 	if amount.Cmp(amount1) != 0 {
// 		t.Errorf("unexpected value : want %d, got %d", amount1, amount)
// 	} else {
// 		t.Logf("%s success get %d", "GetAmount", amount)
// 	}
// 	amount = testGetAmount(caller)
// 	b := new(big.Int)
// 	if new(big.Int).Add(amount, amount1).Cmp(b.Mul(initialSupply, new(big.Int).SetUint64(100000000))) != 0 {
// 		t.Errorf("unexpected value : want %d, got %d", b.Sub(b, amount1), amount)
// 	} else {
// 		t.Logf("%s success get %d", "GetAmount", amount)
// 	}
// }

var ercJsonPath = filepath.Join("", "erc20.json")

func TestERC(t *testing.T) {
	run(t, ercJsonPath)
}

func FromUInt64(n uint64) (out []byte) {
	more := true
	for more {
		b := byte(n & 0x7F)
		n >>= 7
		if n == 0 {
			more = false
		} else {
			b = b | 0x80
		}
		out = append(out, b)
	}
	return
}
