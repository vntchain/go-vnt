package dpos

import (
	"bytes"
	"crypto/ecdsa"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/crypto"
	"github.com/vntchain/go-vnt/params"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestManager(t *testing.T) {
	tests := []struct {
		initWitness []string
		result      []string
	}{
		// Just Create a Manager, and test sort
		{
			initWitness: []string{"A", "B", "C", "D", "E"},
			result:      []string{"A", "B", "C", "D", "E"},
		},
	}

	for i, tt := range tests {
		ap := newTesterAccountPool()

		ws := ap.stringToAddressSorted(tt.initWitness)
		assert.Equal(t, len(tt.initWitness), len(ws))

		conf := &params.DposConfig{
			Period:       2,
			WitnessesNum: 3,
		}

		m := NewManager(conf.Period, ws)
		witness := m.Witnesses

		// results address
		results := ap.stringToAddressSorted(tt.result)

		if len(witness) != len(results) {
			t.Errorf("test:%d, witness lens != result", i)
		}

		for j := 0; j < len(results); j++ {
			if !bytes.Equal(witness[j][:], results[j][:]) {
				t.Errorf("test:%d, witness mismatch result[%d]: %x, %x", i, j, witness[j], results[j])
			}
		}
	}
}

// testerAccountPool maintains current active address
type testerAccountPool struct {
	accounts map[string]*ecdsa.PrivateKey
}

func newTesterAccountPool() *testerAccountPool {
	return &testerAccountPool{
		accounts: make(map[string]*ecdsa.PrivateKey),
	}
}

func (ap *testerAccountPool) address(account string) common.Address {
	// Ensure we have a persistent key for the account
	if ap.accounts[account] == nil {
		ap.accounts[account], _ = crypto.GenerateKey()
	}
	// Resolve and return the VNT address
	return crypto.PubkeyToAddress(ap.accounts[account].PublicKey)
}

func (ap *testerAccountPool) stringToAddress(accounts []string) []common.Address {
	var address []common.Address
	for _, acc := range accounts {
		address = append(address, ap.address(acc))
	}
	return address
}

func (ap *testerAccountPool) stringToAddressSorted(accounts []string) []common.Address {
	address := ap.stringToAddress(accounts)
	// sort
	for i := 0; i < len(address); i++ {
		for j := i + 1; j < len(address); j++ {
			if bytes.Compare(address[i][:], address[j][:]) > 0 {
				address[i], address[j] = address[j], address[i]
			}
		}
	}

	return address
}
