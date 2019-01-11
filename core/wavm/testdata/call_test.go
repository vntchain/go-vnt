package tests

import (
	"encoding/json"
	"fmt"
	"math/big"
	"path/filepath"
	"testing"

	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core/vm"
	"github.com/vntchain/go-vnt/core/wavm"
	wasmContract "github.com/vntchain/go-vnt/core/wavm/contract"
	"github.com/vntchain/go-vnt/crypto"
)

var (
	basepath         = "./"
	callCode         = filepath.Join(basepath, "call/call.wasm")
	callAbi          = filepath.Join(basepath, "call/call.abi")
	callContractAddr = common.HexToAddress("0xdddddddddddddddddddddddddddddddddddddddf")
)

func newCallVM(iscreated bool) *wavm.Wavm {
	wavm := newWavm(callCode, callAbi, iscreated)
	code := getCode(callCode)
	contract := wasmContract.NewWASMContract(vm.AccountRef(caller), vm.AccountRef(callContractAddr), value, gas)
	contract.SetCallCode(&callContractAddr, crypto.Keccak256Hash(code), code)
	wavm.ChainContext.Contract = contract
	return wavm
}

func TestCall(t *testing.T) {
	erc20addr := createERC20()
	fmt.Println(erc20addr)
	calladdr := createCall()
	fmt.Println(calladdr)
	runCall(t, calladdr, erc20addr)
	//for i := 0; i < 100; i++ {
	//	runCall(t, calladdr, erc20addr)
	//}
	//
	//runCall(t, calladdr, common.HexToAddress("0x0"))
}

func createERC20() common.Address {
	initialSupply := new(big.Int)
	initialSupply.SetUint64(1000000000)
	tokenName := "bitcoin"
	tokenSymbol := "BTC"
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
	code := wavm.WasmCode{}
	code.Abi = getCode(erc20Abi)
	code.Code = getCode(erc20Code)

	input := packInput(getABI(erc20Abi), "", initialSupply, tokenName, tokenSymbol)
	res, _ := json.Marshal(code)
	res = append(res, input...)
	// hexres := hex.EncodeToString(res)
	res, ctrAddr, leftgas, err := wvm.Create(vm.AccountRef(caller), res, gas, new(big.Int).SetUint64(0))
	if err != nil {
		fmt.Errorf("ctrAddr %s leftGas %d err %s", ctrAddr.Hex(), leftgas, err)
	} else {
		fmt.Printf("ctrAddr %s leftGas %d err %s", ctrAddr.Hex(), leftgas, err)
	}
	return ctrAddr
}

func createCall() common.Address {
	ctx := vm.Context{
		GetHash: getHash,
		// Message information
		Origin:   origin,
		GasPrice: gasPrice,

		// Block information
		Coinbase:    coinbase,
		GasLimit:    gasLimit,
		BlockNumber: new(big.Int).Add(blockNumber, new(big.Int).SetInt64(1)),
		Time:        time,
		Difficulty:  difficulty,
		CanTransfer: CanTransfer,
		Transfer:    Transfer,
	}
	wvm := newWAVM(ctx)
	code := wavm.WasmCode{}
	code.Abi = getCode(callAbi)
	code.Code = getCode(callCode)
	input := packInput(getABI(callAbi), "")
	res, _ := json.Marshal(code)
	res = append(res, input...)
	res, ctrAddr, _, _ := wvm.Create(vm.AccountRef(caller), res, gas, new(big.Int).SetUint64(0))
	return ctrAddr

}

//calladdr 0x553E6c30Af61e7A3576f31311EA8a620F80D047e erc20addr 0x8Af6A7AF30d840ba137e8F3F34d54CfB8BEbA6E2
func runCall(t *testing.T, calladdr common.Address, erc20addr common.Address) {
	fmt.Printf("runCall calladdr %s erc20addr %s", calladdr.Hex(), erc20addr.Hex())
	ctx := vm.Context{
		GetHash: getHash,
		// Message information
		Origin:   origin,
		GasPrice: gasPrice,

		// Block information
		Coinbase:    coinbase,
		GasLimit:    gasLimit,
		BlockNumber: new(big.Int).Add(blockNumber, new(big.Int).SetInt64(2)),
		Time:        time,
		Difficulty:  difficulty,
		CanTransfer: CanTransfer,
		Transfer:    Transfer,
	}
	wavm := newWAVM(ctx)
	callgas := uint64(1000000)
	amount := new(big.Int).SetUint64(0)
	input := packInput(getABI(callAbi), "Test_GetTokenName", erc20addr, amount, callgas)
	res, leftgas, err := wavm.Call(vm.AccountRef(caller), calladdr, input, gas, new(big.Int).SetUint64(0))
	var str string
	unpackOutput(getABI(callAbi), &str, "Test_GetTokenName", res)
	fmt.Printf("res str %s\n", str)
	fmt.Printf("res %s gas %d leftGas %d err %s\n", res, gas, leftgas, err)
	if str != "bitcoin" {
		t.Errorf("unexpected value : want %s, got %s", "bitcoin", str)
	} else {
		t.Logf("get expected name %s", str)
	}
}
