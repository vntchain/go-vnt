package wavm

import (
	"math/big"
	"github.com/vntchain/go-vnt/core/state"
	"github.com/vntchain/go-vnt/core/vm/interface"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/accounts/abi"
	"github.com/vntchain/go-vnt/core/wavm/gas"
	"github.com/vntchain/go-vnt/params"
	"github.com/vntchain/go-vnt/core/wavm/contract"
	"github.com/vntchain/go-vnt/core/wavm/storage"
)

type ChainContext struct {
	// CanTransfer returns whether the account contains
	// sufficient ether to transfer the value
	CanTransfer func(inter.StateDB, common.Address, *big.Int) bool
	// Transfer transfers ether from one account to the other
	Transfer func(inter.StateDB, common.Address, common.Address, *big.Int)
	// GetHash returns the hash corresponding to n
	GetHash func(uint64) common.Hash
	// Message information
	Origin   common.Address // Provides information for ORIGIN
	GasPrice *big.Int       // Provides information for GASPRICE

	// Block information
	Coinbase       common.Address // Provides information for COINBASE
	GasLimit       uint64         // Provides information for GASLIMIT
	BlockNumber    *big.Int       // Provides information for NUMBER
	Time           *big.Int       // Provides information for TIME
	Difficulty     *big.Int       // Provides information for DIFFICULTY
	StateDB        *state.StateDB
	Contract       *contract.WASMContract
	Code           []byte  //Wasm contract code
	Abi            abi.ABI //Wasm contract abi
	Wavm           *WAVM
	IsCreated      bool
	StorageMapping map[uint64]storage.StorageMapping
	GasRule        gas.Gas
	GasCounter     gas.GasCounter
	GasTable       params.GasTable
}
