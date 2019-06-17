package tests

import (
	"io/ioutil"
	"math/big"
	"os"

	"github.com/vntchain/go-vnt/accounts/abi"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/common/hexutil"
	"github.com/vntchain/go-vnt/core"
	"github.com/vntchain/go-vnt/core/state"
	"github.com/vntchain/go-vnt/core/wavm"
	"github.com/vntchain/go-vnt/crypto"
	"github.com/vntchain/go-vnt/crypto/sha3"
	"github.com/vntchain/go-vnt/log"
	"github.com/vntchain/go-vnt/rlp"
	"github.com/vntchain/go-vnt/vntdb"
)

var (
	basepath = "./"
)

var (
	activeKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f292")
	activeAddr   = crypto.PubkeyToAddress(activeKey.PublicKey)
)

var logger *log.Logger

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
}

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
	Function string        `json:"function"`
	Input    []argument    `json:"input"`
	RawInput hexutil.Bytes `json:"rawinput"`
	Wanted   argument      `json:"wanted"`
	Error    string        `json:"error"`
	Event    []argument    `json:"event"`
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
	activeAccount := core.GenesisAccount{
		Code:    []byte{},
		Storage: map[common.Hash]common.Hash{},
		Balance: big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18)),
		Nonce:   0,
	}
	accounts[activeAddr] = activeAccount
	for addr, a := range accounts {
		statedb.SetCode(addr, a.Code)
		statedb.SetNonce(addr, a.Nonce)
		statedb.SetBalance(addr, a.Balance)
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
