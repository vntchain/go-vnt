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
    "config": {
        "chainId": 1,
        "dpos": {
            "period": 2,
            "witnessesnum": 19,
            "witnessesUrl": [
                "/ip4/47.106.71.114/tcp/3001/ipfs/1kHh6iu6GiXidWZCm3B7kw4HChM4CncLiccTbQtJSUrDpnR",
                "/ip4/47.108.69.101/tcp/3001/ipfs/1kHeyfXwiuXLbNFLrCsu54gBPkf3e2J8hvbXBsD5NTfwEA1",
                "/ip4/47.108.67.119/tcp/3001/ipfs/1kHhg1CeC5h8TT7UYtnrk5f6d27p89nqRSRvfX6uNBSYsVG",
                "/ip4/39.100.143.156/tcp/3001/ipfs/1kHd9mkdbw2smReu9G4dGxZ38JNyni6ajNSq9crTwsFiazg",
                "/ip4/118.190.59.122/tcp/3001/ipfs/1kHYCNwh1SVfTWGrgeSzYLH65NMrzVbjMCFKp2KKhCqfd42",
                "/ip4/118.190.59.100/tcp/3001/ipfs/1kHivFPKBXSwtLTkjTuzuMPpbHbDuVh6rQwGGzUXPjh1sSw",
                "/ip4/47.56.69.191/tcp/3001/ipfs/1kHHMELQGozJeaGckomHSMSymwgYWj2cRR2uSgn9y5eB7rV",
                "/ip4/39.97.171.233/tcp/3001/ipfs/1kHLWCTi4qqfZw13f393K79Qjmo7yNEVTPtozLpjvhTvwCs",
                "/ip4/47.103.107.188/tcp/3001/ipfs/1kHbmc5hvBcQRoWh5MhpLM3ryKiQdukRJZEF4CrVeKWYeHc",
                "/ip4/47.103.57.160/tcp/3001/ipfs/1kHG2ZxeGmVxrWXm18Y8eQsmNdwqofb1ExGResSM4P86RF8",
                "/ip4/47.254.235.57/tcp/3001/ipfs/1kHDWP8wPvZ9UTDthgfbJ4uygwsziYCoSVHypUPRqJfoGJb",
                "/ip4/120.77.236.120/tcp/3001/ipfs/1kHmBZUaPtEmEZPhL1wUFprEH27vSjoF75duK7Wv9dbYri2",
                "/ip4/47.111.131.2/tcp/3001/ipfs/1kHC41ck2NwkyNxtEuEsKcYpv5iyGC9j4ekgE3B7BcUNW3D",
                "/ip4/47.88.217.237/tcp/3001/ipfs/1kHevyN16xUnQR5yD8DWa4VtiQpJ9kWYnboTEPtRhB5MDRs",
                "/ip4/47.91.19.11/tcp/3001/ipfs/1kHQaPVKPkoSaoUiJbC2GTRXC5eEDCSRTJY4nw9irrFjoN3",
                "/ip4/47.254.20.76/tcp/3001/ipfs/1kHfn7yfdJx4x2f8fh7ZzxFuA5f5KBGyfUWgBUHXMQWKFaL",
                "/ip4/47.93.191.135/tcp/3001/ipfs/1kHCnrsiTwr9y7q8zBCxE6DdHNPTfLfsHYdCzdjGYvdKpYY",
                "/ip4/101.132.191.42/tcp/3001/ipfs/1kHdWEpRxfqYzc9K5SS617NwNHSQUYBcQJMBVp6QwqPzAgJ",
                "/ip4/39.104.62.26/tcp/3001/ipfs/1kHTiT8vJ73EQWpJC57dpsjJ4Erz1VoS61zpfPtaYuJ6iZt"
            ]
        }
    },
    "timestamp": "0x5d18dc80",
    "extraData": "0x546f206578706c6f726520737472616e6765206e657720776f726c6473efbc8c746f207365656b206f7574206e6577206c69666520616e64206e657720636976696c697a6174696f6e73efbc8c746f20626f6c646c7920676f207768657265206e6f206f6e652068617320676f6e65206265666f72652e",
    "gasLimit": "0x47b760",
    "difficulty": "0x1",
    "coinbase": "0x0000000000000000000000000000000000000000",
    "alloc": {
        "08257303e7f2ed7529cb81d2521d70a652e38238": {
            "balance": "0x48ec9b2d230d5a7c00000"
        },
        "09fd5d32fdff6e28e93cd9d9dfbaf5821d15e8d0": {
            "balance": "0x19d971e4fe8401e74000000"
        },
        "0acef1c9182e92ad80618e177f6c3bcac85a7a0a": {
            "balance": "0xf8277896582678ac000000"
        },
        "0e9a4daadb35f37a081f8d48b5d1fc793dc9ea7c": {
            "balance": "0xae09eedac3a019ec00000"
        },
        "0eefec6d3bf82d84dbd0666f5e1db6298a4b1888": {
            "balance": "0x68b01f58f4fdcc3c00000"
        },
        "11556936ebe002f4948c5af40c046565df0fbd74": {
            "balance": "0x33b2e3c9fd0803ce8000000"
        },
        "14f19ea14598e482905f0d1985d4bb37dd413c40": {
            "balance": "0xbe0d928eb933fdf000000"
        },
        "17d72a1ae6aab8d2b59886d3252bea33d1058cb8": {
            "balance": "0x1cda2096bc8fbca800000"
        },
        "28e77afbcdbc628de497a85baa477660e5aa10f4": {
            "balance": "0x1abc06b5f2d50a6800000"
        },
        "2f6fce36a2dda0374c3a292e55d319e5cfc809e5": {
            "balance": "0xf8277896582678ac000000"
        },
        "425fba17477cb54b8d4047bb5930e00d2d5100f5": {
            "balance": "0x4646fad426e3fbac00000"
        },
        "4279408c1ed0ce78fc19a23e22af98208d14374b": {
            "balance": "0x1c747bbc96bcbb3c00000"
        },
        "43a3ede687f505ed106b41626d5b70e1bf550f7d": {
            "balance": "0x1abc06b5f2d50a6800000"
        },
        "4417a444c00b007803196ade124a300470f5fd95": {
            "balance": "0x19d971e4fe8401e74000000"
        },
        "4721b8e0594ca55d21d5de430dc8076820ecb896": {
            "balance": "0x2421998b7201816400000"
        },
        "4bc6af201fad669cd315d18d503a3bb7f555b673": {
            "balance": "0x1c0ed6e270e9b9d000000"
        },
        "4eebc801e52fd63e90ec1660bbe492f6de27e167": {
            "balance": "0x9b60aacdd1e2d71800000"
        },
        "4f6b04717b0dbfd5397266f81fef868ca4c68efd": {
            "balance": "0xf8277896582678ac000000"
        },
        "59879b671e2e034f387f625e7e123828def0d798": {
            "balance": "0xb377112ac88e82b000000"
        },
        "5d4fad8937f72cec3ca46ce1c14349018058213e": {
            "balance": "0xc6205537ba4bc58400000"
        },
        "66a101458305ae2cec192a94751055f025648869": {
            "balance": "0x89a23c233886b411680000"
        },
        "77af92b3bd4acf994d7518c79ae2f303e9b4c6eb": {
            "balance": "0x1adde853ff70b58c00000"
        },
        "7831db47c51c01666a435021888f2dcd04dcf02f": {
            "balance": "0x62120e7a7e965f3400000"
        },
        "8d6ae3764e5d0fdbfcabdd443dcdf300b4574c03": {
            "balance": "0xa56fa5b99019a5c8000000"
        },
        "8e667c898c2f2d7ebf1265652a35638347067754": {
            "balance": "0x95f29ecfdc8d99ddc980000"
        },
        "8eb227bb885ef74b52e6c8a3b6662de641a93cdf": {
            "balance": "0x19d971e4fe8401e74000000"
        },
        "90e1814cbef1aaa66acdab66d94109b12ec5afad": {
            "balance": "0xb5facfe5b81c365c00000"
        },
        "92dc2b4b56ce19a394ca6ff7af42e70296b6813b": {
            "balance": "0xa9027164e484b29400000"
        },
        "93355b1d40b20b92231e1bf2fbcb813ed1810547": {
            "balance": "0x19d971e4fe8401e74000000"
        },
        "983c54b64024bb92f42a065ca55e2d6b408a17be": {
            "balance": "0x33b2e3c9fd0803ce8000000"
        },
        "9b3e39a565f31aa669574668b8302ff500e9b73c": {
            "balance": "0xc0b332e7b55d5cc000000"
        },
        "9d69b5d31f33172435ec20ac847dfb3ba6bf94a4": {
            "balance": "0xa56fa5b99019a5c8000000"
        },
        "a14c79d97e2c58de4499a552a2ec72823ed93db0": {
            "balance": "0x1c6f307be4c4687e6000000"
        },
        "b49eb8f256dce146d6ed2682469e3a97f469575d": {
            "balance": "0x24cb01a1b10bd91800000"
        },
        "c8db01d838946a1be94cc18ff454ed32194cb1b6": {
            "balance": "0x96592d57f2c76fc000000"
        },
        "c9024f109ace4b50c800f857c57538a34eb2728c": {
            "balance": "0x5b521bfdfb93470800000"
        },
        "cc6991f75fe5bd8b73407d4b1d8a5c7c4d2e1e3a": {
            "balance": "0x2bf1a8054a46d0092000000"
        },
        "cdfbb44442e595260e2148e8f0caadfe149daa4c": {
            "balance": "0x1dc74be914d16aa400000"
        },
        "cf4b31e13d8f31d693ce609d62ea11d0ae663a40": {
            "balance": "0xc5fe7399adb01a6000000"
        },
        "d699607661f8a38b1eaab41274cf3cbacf267592": {
            "balance": "0xba9ca88171649c4800000"
        },
        "d9b7c50f6509a5aa747535c5248a58ff504fd505": {
            "balance": "0x22f0aafd00887d2000000"
        },
        "dc3c71f9f694235eeac29254902d51b1cb36638a": {
            "balance": "0x1abc06b5f2d50a6800000"
        },
        "e02c15f9a2a893b9d6ef4fdae67cef85e7e74af7": {
            "balance": "0xbe2f742cc5cfa91400000"
        },
        "e7e5dcf58d99c5fbb12084974cd8c6d9c89c9551": {
            "balance": "0x9920af4efb8c79b400000"
        }
    },
    "witnesses": [
        "0x91837ff26639700c9688cf8f3fe92bd8b2ec806d",
        "0x3c60a032ba3c6177e50188748e55e5894fb241e4",
        "0xaa2b5f39fb2a4aee56db3ee19567f699d30df1a1",
        "0x61a6e04c737483d72c20de6e71dd8cbb6f6c747d",
        "0x186bae02dc3444d2bb3d39504fefdc9754860481",
        "0xf4c8fd44490493000b8776fd1597752bd9ede431",
        "0x4e94885ed5cfe31a00c7496176f59fdc5e5c7a71",
        "0x4b47c3262a9d2c309b692c9220898ff728054c00",
        "0x31ba9c8cf34d7cc0957a95744b245322af427786",
        "0x4dcfcd45b253119c0d3db6b9ba9e154167dd6a58",
        "0xe6c745142283dbbe4b4a03e969525e25031939fa",
        "0xc61a92dd1713f9ba2214f0ce92e3d408ba4d426d",
        "0xc221a4d0b30dee366bc7899dd29e0f7ac9a7e45a",
        "0xddfd32c4d33915685b926ba5eaab3860db1690cd",
        "0xd338d81c4723982c815a294de3b38608dad9962c",
        "0x6cd54fc6da0f044c43d4550d87ae10b9e1cea351",
        "0xd328d8864649ed050b3d8e9d77f94c75299fd243",
        "0x386dd85ad17b6bd60d2d142473b54bf9d5439842",
        "0x4b8a6cff7b9e008caa936aadd33d9be048623d53"
    ],
    "number": "0x0",
    "gasUsed": "0x0",
    "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000"
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
		result:  "0x546f206578706c6f726520737472616e6765206e657720776f726c6473efbc8c746f207365656b206f7574206e6577206c69666520616e64206e657720636976696c697a6174696f6e73efbc8c746f20626f6c646c7920676f207768657265206e6f206f6e652068617320676f6e65206265666f72652e",
	},
	// Genesis file with specific chain configurations
	{
		genesis: genesisJson,
		query:   "core.getBlock(0).gasLimit",
		result:  "4700000",
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
