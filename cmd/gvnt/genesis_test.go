// Copyright 2016 The go-ethereum Authors
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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var genesisJson = `{
	    "alloc":{
	    },
	    "coinbase":"0x0000000000000000000000000000000000000000",
	    "difficulty":"0x1",
	    "extraData":"0x",
	    "gasLimit":"0x2fefd8",
	    "parentHash":"0x0000000000000000000000000000000000000000000000000000000000000000",
	    "timestamp":"0x00",
	    "config":{
	        "chainId":1333,
	        "dpos":{
	            "period":2,
	            "witnessesnum":4,
	            "witnessesUrl":[
	                "/ip4/127.0.0.1/tcp/5210/ipfs/1kHcch6yuBCgC5nPPSK3Yp7Es4c4eenxAeK167pYwUvNjRo",
	                "/ip4/127.0.0.1/tcp/5211/ipfs/1kHJFKr2bzUnMr1NbeyYbYJa3RXT18cEu7cNDrHWjg8XYKB",
	                "/ip4/127.0.0.1/tcp/5212/ipfs/1kHfop9dnUHHmtBXVkLB5UauAmACtrsEX5H5t6oCRpdL198",
	                "/ip4/127.0.0.1/tcp/5213/ipfs/1kHHWuQNUVV2wgE8SqzQjWhiFQcfpkP5tRVTdJXAPWVj4nR"
	            ]
	        }
	    },
	    "witnesses":[
	        "0x122369f04f32269598789998de33e3d56e2c507a",
	        "0x42a875ac43f2b4e6d17f54d288071f5952bf8911",
	        "0x3dcf0b3787c31b2bdf62d5bc9128a79c2bb18829",
	        "0xbf66d398226f200467cd27b14e85b25a8c232384"
	    ]
	}`

var customGenesisTests = []struct {
	genesis string
	query   string
	result  string
}{
	// Plain genesis file without anything extra
	// in real environment the genesis is worked, and there are blocks produced
	{
		genesis: genesisJson,
		query:   "core.getBlock(0).difficulty",
		result:  "1",
	},
	// Genesis file with an empty chain configuration (ensure missing fields work)
	{
		genesis: genesisJson,
		query:   "core.getBlock(0).extraData",
		result:  "0x",
	},
	// Genesis file with specific chain configurations
	{
		genesis: genesisJson,
		query:   "core.getBlock(0).gasLimit",
		result:  "3141592",
	},
}

// Tests that initializing Gvnt with a custom genesis block and chain definitions
// work properly.
func TestCustomGenesis(t *testing.T) {
	for i, tt := range customGenesisTests {
		// Create a temporary data directory to use and inspect later
		datadir := tmpdir(t)
		defer os.RemoveAll(datadir)

		// Initialize the data directory with the custom genesis block
		json := filepath.Join(datadir, "genesis.json")
		if err := ioutil.WriteFile(json, []byte(tt.genesis), 0600); err != nil {
			t.Fatalf("test %d: failed to write genesis file: %v", i, err)
		}
		runGvnt(t, "--datadir", datadir, "init", json).WaitExit()

		// Query the custom genesis block
		gvnt := runGvnt(t,
			"--datadir", datadir, "--maxpeers", "0", "--port", "0",
			"--nodiscover", "--nat", "none", "--ipcdisable",
			"--exec", tt.query, "console")
		gvnt.ExpectRegexp(tt.result)
		gvnt.ExpectExit()
	}
}
