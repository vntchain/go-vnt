package tests

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/vntchain/go-vnt/common/math"
	"github.com/vntchain/go-vnt/crypto"
)

func TestSig(t *testing.T) {
	res := crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))
	fmt.Printf("res %s\n", res.Hex())
	str, _ := hex.DecodeString("a9059cbb0000000000000000000000003dcf0b3787c31b2bdf62d5bc9128a79c2bb1882900000000000000000000000000000000000000000000000000000000000003e8")
	fmt.Printf("str %+v\n", str)
}

func TestSafeAdd(t *testing.T) {
	b := new(big.Int)
	b.SetUint64(2)
	e := new(big.Int)
	e.SetUint64(255)

	res := math.Exp(b, e)
	fmt.Printf("res %d\n", res)
	subRes := math.U256(new(big.Int).Sub(res, new(big.Int).SetUint64(0)))
	addRes := math.U256(new(big.Int).Add(subRes, res))
	fmt.Printf("addRes %d\n", addRes)
}
