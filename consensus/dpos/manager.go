package dpos

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/log"
	"math/big"
)

// Manager is used to judge weather a witness is block producer or not.
// - Check weather a witness is block producer or not
// - Check weather a witness is in witness list
type Manager struct {
	blockPeriod uint64           // Period of producing a block
	Witnesses   []common.Address // Witness list, all current witness
}

// NewManager creating a dpos witness manager
func NewManager(period uint64, witnesses []common.Address) *Manager {
	wm := &Manager{
		blockPeriod: period,
		Witnesses:   witnesses,
	}
	return wm
}

// inTurn check whether witness is the block producer at witTime
func (m *Manager) inTurn(witness, pWitness common.Address, witTime, pWitTime *big.Int) bool {
	// sanity check
	if witTime.Cmp(pWitTime) <= 0 {
		log.Debug("Time is early")
		return false
	}
	witIndex := m.indexOf(witness)
	pIndex := m.indexOf(pWitness)
	if witIndex == -1 {
		log.Warn("You are not in the witnesses list", "addr", witness)
		return false
	}
	if pIndex == -1 {
		log.Warn("Previous witness is not in the witnesses list", "preWitness", pWitness)
		return false
	}

	// calc offset with timestamp
	dur := new(big.Int).Sub(witTime, pWitTime)
	period := new(big.Int).SetUint64(m.blockPeriod)
	left := big.NewInt(0)
	nPeriod, left := new(big.Int).DivMod(dur, period, left)
	if left.Cmp(big.NewInt(0)) != 0 {
		nPeriod.Add(nPeriod, big.NewInt(1)) // witTime in next period
	}

	// make sure offset in an safety range:[0, len(m.Witnesses))
	offset := big.NewInt(0)
	nWitness := big.NewInt(int64(len(m.Witnesses)))
	_, offset = new(big.Int).DivMod(nPeriod, nWitness, offset)

	iOffset := int(offset.Int64())

	// get block producer's index by offset
	targetIndex, err := m.moveIndex(pIndex, iOffset)
	if err != nil {
		log.Warn(err.Error())
		return false
	}

	// match curIndex and witIndex
	if targetIndex == witIndex {
		return true
	}

	log.Debug("witness index is not match", "targetIndex", targetIndex, "witIndex", witIndex)
	return false
}

// indexOf get the index of witness in witness list
func (m *Manager) indexOf(witness common.Address) int {
	for i := 0; i < len(m.Witnesses); i++ {
		if addressEqual(m.Witnesses[i], witness) {
			return i
		}
	}
	return -1
}

// has check whether witness is in witness list
func (m *Manager) has(witness common.Address) bool {
	return -1 != m.indexOf(witness)
}

// moveIndex move the index with offset
func (m *Manager) moveIndex(cur int, off int) (int, error) {
	if cur < 0 || cur > len(m.Witnesses) {
		log.Warn("Current index is out of range", "index", cur, "range", len(m.Witnesses)-1)
		return -1, errors.New("invalid current index")
	}

	cur += off
	cur = cur % len(m.Witnesses)
	return cur, nil
}

// dump witness list
func (m *Manager) dump() {
	fmt.Println("Witness list:")
	for _, wit := range m.Witnesses {
		fmt.Printf("%x \n", wit)
	}
}

// addressEqual check wether two address is queal
func addressEqual(addr1, addr2 common.Address) bool {
	return bytes.Equal(addr1[:], addr2[:])
}
