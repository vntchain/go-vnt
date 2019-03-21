// Copyright 2017 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"time"

	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core"
	"github.com/vntchain/go-vnt/log"
	"github.com/vntchain/go-vnt/params"
)

// makeGenesis creates a new genesis struct based on some user input.
func (w *wizard) makeGenesis() {
	// Construct a default genesis block
	genesis := &core.Genesis{
		Timestamp:  uint64(time.Now().Unix()),
		GasLimit:   4700000,
		Difficulty: big.NewInt(524288),
		Alloc:      make(core.GenesisAlloc),
		Config: &params.ChainConfig{
			HubbleBlock: big.NewInt(1),
		},
	}
	// Figure out which consensus engine to choose
	fmt.Println()

	genesis.Difficulty = big.NewInt(1)
	genesis.Config.Dpos = &params.DposConfig{
		Period:       3,
		WitnessesNum: 15,
	}
	fmt.Println()
	fmt.Println("How many seconds should blocks take? (default = 3)")
	genesis.Config.Dpos.Period = uint64(w.readDefaultInt(3))
	fmt.Println()
	fmt.Println("How many witnesses produce blocks (default = 15)?")
	genesis.Config.Dpos.WitnessesNum = w.readDefaultInt(15)

	// We also need the initial list of signers
RECONFIGWITNESS:
	fmt.Println()
	fmt.Println("Which accounts are allowed to seal? (The number of accounts should be witnesses number)")
	signers := make([]common.Address, 0)
	for {
		if address := w.readAddress(); address != nil {
			signers = append(signers, *address)
			continue
		}
		if len(signers) > 0 {
			break
		}
	}
	fmt.Println()
	fmt.Println("the node info of witnesses? (The number of accounts should be witnesses number)")
	nodeInfos := make([]string, 0)
	for {
		if nodeInfo := w.readNodeInfo(); nodeInfo != "" {
			nodeInfos = append(nodeInfos, nodeInfo)
			continue
		}
		if len(nodeInfos) > 0 {
			break
		}
	}
	if genesis.Config.Dpos.WitnessesNum != len(signers) {
		fmt.Printf("Error: the number of witnesses is not match previous config(%d). Please reconfig\n", genesis.Config.Dpos.WitnessesNum)
		goto RECONFIGWITNESS
	}
	genesis.Witnesses = make([]common.Address, len(signers))
	genesis.Config.Dpos.WitnessesUrl = make([]string, len(nodeInfos))
	for i, signer := range signers {
		genesis.Witnesses[i] = signer
	}
	for i, nodeInfo := range nodeInfos {
		genesis.Config.Dpos.WitnessesUrl[i] = nodeInfo
	}

	// Consensus all set, just ask for initial funds and go
	fmt.Println()
	fmt.Println("Which accounts should be pre-funded? (advisable at least one)")
	for {
		// Read the address of the account to fund
		if address := w.readAddress(); address != nil {
			genesis.Alloc[*address] = core.GenesisAccount{
				Balance: new(big.Int).Lsh(big.NewInt(1), 256-7), // 2^256 / 128 (allow many pre-funds without balance overflows)
			}
			continue
		}
		break
	}
	// Add a batch of precompile balances to avoid them getting deleted
	for i := int64(0); i < 256; i++ {
		genesis.Alloc[common.BigToAddress(big.NewInt(i))] = core.GenesisAccount{Balance: big.NewInt(1)}
	}
	// Query the user for some custom extras
	fmt.Println()
	fmt.Println("Specify your chain/network ID if you want an explicit one (default = random)")
	genesis.Config.ChainID = new(big.Int).SetUint64(uint64(w.readDefaultInt(rand.Intn(65536))))

	// All done, store the genesis and flush to disk
	log.Info("Configured new genesis block")

	w.conf.Genesis = genesis
	w.conf.flush()
}

// manageGenesis permits the modification of chain configuration parameters in
// a genesis config and the export of the entire genesis spec.
func (w *wizard) manageGenesis() {
	// Figure out whether to modify or export the genesis
	fmt.Println()
	fmt.Println(" 1. Modify existing fork rules")
	fmt.Println(" 2. Export genesis configuration")
	fmt.Println(" 3. Remove genesis configuration")

	choice := w.read()
	switch {
	case choice == "1":
		// Fork rule updating requested, iterate over each fork
		fmt.Println()
		fmt.Printf("Which block should Hubble come into effect? (default = %v)\n", w.conf.Genesis.Config.HubbleBlock)
		w.conf.Genesis.Config.HubbleBlock = w.readDefaultBigInt(w.conf.Genesis.Config.HubbleBlock)

		out, _ := json.MarshalIndent(w.conf.Genesis.Config, "", "  ")
		fmt.Printf("Chain configuration updated:\n\n%s\n", out)

	case choice == "2":
		// Save whatever genesis configuration we currently have
		fmt.Println()
		fmt.Printf("Which file to save the genesis into? (default = %s.json)\n", w.network)
		out, _ := json.MarshalIndent(w.conf.Genesis, "", "  ")
		if err := ioutil.WriteFile(w.readDefaultString(fmt.Sprintf("%s.json", w.network)), out, 0644); err != nil {
			log.Error("Failed to save genesis file", "err", err)
		}
		log.Info("Exported existing genesis block")

	case choice == "3":
		// Make sure we don't have any services running
		if len(w.conf.servers()) > 0 {
			log.Error("Genesis reset requires all services and servers torn down")
			return
		}
		log.Info("Genesis block destroyed")

		w.conf.Genesis = nil
		w.conf.flush()

	default:
		log.Error("That's not something I can do")
	}
}
