package tests

import (
	"fmt"
	"io/ioutil"
	"math/big"

	"github.com/vntchain/go-vnt/accounts/abi"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core"
	"github.com/vntchain/go-vnt/core/state"
	"github.com/vntchain/go-vnt/core/wavm"
	"github.com/vntchain/go-vnt/crypto"
	"github.com/vntchain/go-vnt/crypto/sha3"
	"github.com/vntchain/go-vnt/rlp"
	"github.com/vntchain/go-vnt/vntdb"
)

type vmJSON struct {
	Env      stEnv             `json:"env"`
	Exec     vmExec            `json:"exec"`
	Pre      core.GenesisAlloc `json:"pre"`
	TestCase []testCase        `json:"testcase"`
}

type stEnv struct {
	Coinbase   common.Address `json:"currentCoinbase"`
	Difficulty *big.Int       `json:"currentDifficulty"`
	GasLimit   uint64         `json:"currentGasLimit"`
	Number     uint64         `json:"currentNumber"`
	Timestamp  uint64         `json:"currentTimestamp"`
}

type vmExec struct {
	Address  common.Address `json:"address"`
	Value    *big.Int       `json:"value"`
	GasLimit uint64         `json:"gas"`
	Caller   common.Address `json:"caller"`
	Origin   common.Address `json:"origin"`
	GasPrice *big.Int       `json:"gasPrice"`
}

type testCase struct {
	Code     string   `json:"code"`
	Abi      string   `json:"abi"`
	InitCase initcase `json:"initcase"`
	Tests    []tests  `json:"tests"`
}

type initcase struct {
	NeedInit bool       `json:"needinit"`
	Input    []argument `json:"input"`
}

type tests struct {
	Function string     `json:"function"`
	Input    []argument `json:"input"`
	Wanted   argument   `json:"wanted"`
	Event    []argument `json:"event"`
}

type argument struct {
	Data     string `json:"data"`
	DataType string `json:"type"`
}

func vmTestBlockHash(n uint64) common.Hash {
	return common.BytesToHash(crypto.Keccak256([]byte(big.NewInt(int64(n)).String())))
}

func MakePreState(db vntdb.Database, accounts core.GenesisAlloc) *state.StateDB {
	sdb := state.NewDatabase(db)
	statedb, _ := state.New(common.Hash{}, sdb)
	for addr, a := range accounts {
		statedb.SetCode(addr, a.Code)
		statedb.SetNonce(addr, a.Nonce)
		statedb.SetBalance(addr, a.Balance)
		fmt.Printf("addr %s a.balance %d\n", addr.Hex(), a.Balance)
		for k, v := range a.Storage {
			statedb.SetState(addr, k, v)
		}
	}
	// Commit and re-open to start with a clean state.
	root, _ := statedb.Commit(false)
	statedb, _ = state.New(root, sdb)
	return statedb
}

func readFile(filepath string) []byte {
	code, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	return code
}

func getABI(filepath string) abi.ABI {
	abi, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	abiobj, err := wavm.GetAbi(abi)
	if err != nil {
		panic(err)
	}
	return abiobj
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}
