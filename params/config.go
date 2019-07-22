// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package params

import (
	"fmt"
	"math/big"

	"github.com/vntchain/go-vnt/common"
)

// Genesis hashes to enforce below configs on.
var (
	MainnetGenesisHash = common.HexToHash("0xf2ba87b7a6b3c3ff2fc2c0e2b2985b9cfc7ca24daaa39ac04677855b1583e5ad")
)

var (
	// MainnetChainConfig is the chain parameters to run a node on the main network.
	MainnetChainConfig = &ChainConfig{
		ChainID: big.NewInt(1),
		Dpos: &DposConfig{
			Period:       2,
			WitnessesNum: 19,
			WitnessesUrl: []string{
				"/ip4/47.106.71.114/tcp/3001/ipfs/1kHfbk8u12U1HSyaqAe6f622wVMESHcvFZ8VcbGKwsrtT6H",
				"/ip4/47.108.69.101/tcp/3001/ipfs/1kHb7UpvD2zbEgCCPzboTtJLENQ1YhEZ1H2A7QNTb18sHd4",
				"/ip4/47.108.67.119/tcp/3001/ipfs/1kHjhzoNTxEpFvRky1oVvCUAPBUuGzBvBcXmkroDU9NhACg",
				"/ip4/39.100.143.156/tcp/3001/ipfs/1kHG1essWxbjSUjwKMrJcva6Y7XfPoBN86pUemd4wnd2X2A",
				"/ip4/118.190.59.122/tcp/3001/ipfs/1kHBXTtDX4JqStm4qBNCcCZar9isCyu74BPgX6b2odm7zw7",

				"/ip4/118.190.59.100/tcp/3001/ipfs/1kHWbCuCgjLnERsaVr7FCSAWFiA7Sx1CdzstMQKKkEWdYSr",
				"/ip4/47.56.69.191/tcp/3001/ipfs/1kHcupmVj3eLe6QdgrXRn5qetpQQc4XhYc6LXNhQUkRCDK5",
				"/ip4/39.97.171.233/tcp/3001/ipfs/1kHDnV8YSUrHNLf6NGWY7QHxBjMLnKPxKFKvMdPLt1Gg7a6",
				"/ip4/47.103.107.188/tcp/3001/ipfs/1kHTKpK29Kw2EvJ3C817yn6wNDgpXw6oWAkUz8AzAozrVHD",
				"/ip4/47.103.57.160/tcp/3001/ipfs/1kHD3RiEFGZE2SQQgGhX1yzipcR8sE6cSaTj3Xi1yvF1EBL",

				"/ip4/47.254.235.57/tcp/3001/ipfs/1kHMwGapAtV92rXgkxRx4dgjUZnXVtE6wUa8Bk3yzFYLz4h",
				"/ip4/120.77.236.120/tcp/3001/ipfs/1kHHW1DrdwETgtZmDdERZEDgKiAfP2SDfxz2oKuzjuJK2B9",
				"/ip4/47.111.131.2/tcp/3001/ipfs/1kHCcTMZ8EjRm23nffDsYEEARohWwF4ks6zdsaBs3JbzKnA",
				"/ip4/47.88.217.237/tcp/3001/ipfs/1kHVYJ6tckDB5gChvDFo46esBWzz6aWsMaEXfxJPNwruHWZ",
				"/ip4/47.91.19.11/tcp/3001/ipfs/1kHBa9E1onVKmruWvefJHSNFGokkNc3ESZebnq2oJFRDFDG",

				"/ip4/47.254.20.76/tcp/3001/ipfs/1kHKUGfEQ4nzpWG4SkQAUPicrKHCVo32WTgNEQbM75rNEbx",
				"/ip4/47.93.191.135/tcp/3001/ipfs/1kHVzyb8mczCYNB3suPCXac7HMPEc4XyvibHE2hQsh8ehEk",
				"/ip4/101.132.191.42/tcp/3001/ipfs/1kHCoMGo8ANpjFNMyYFZLTHruBJ5DrAXYVBvfQP2CZWKf6h",
				"/ip4/39.104.62.26/tcp/3001/ipfs/1kHHvSatMYVHDdDFNmMyta8tNY4c4VvtifkBAxfjNG7wSpU",
			},
		},
	}

	MainnetChainWitnesses = []common.Address{
		common.HexToAddress("0x91837ff26639700c9688cf8f3fe92bd8b2ec806d"),
		common.HexToAddress("0x3c60a032ba3c6177e50188748e55e5894fb241e4"),
		common.HexToAddress("0xaa2b5f39fb2a4aee56db3ee19567f699d30df1a1"),
		common.HexToAddress("0x61a6e04c737483d72c20de6e71dd8cbb6f6c747d"),
		common.HexToAddress("0x186bae02dc3444d2bb3d39504fefdc9754860481"),

		common.HexToAddress("0xf4c8fd44490493000b8776fd1597752bd9ede431"),
		common.HexToAddress("0x4e94885ed5cfe31a00c7496176f59fdc5e5c7a71"),
		common.HexToAddress("0x4b47c3262a9d2c309b692c9220898ff728054c00"),
		common.HexToAddress("0x31ba9c8cf34d7cc0957a95744b245322af427786"),
		common.HexToAddress("0x4dcfcd45b253119c0d3db6b9ba9e154167dd6a58"),

		common.HexToAddress("0xe6c745142283dbbe4b4a03e969525e25031939fa"),
		common.HexToAddress("0xc61a92dd1713f9ba2214f0ce92e3d408ba4d426d"),
		common.HexToAddress("0xc221a4d0b30dee366bc7899dd29e0f7ac9a7e45a"),
		common.HexToAddress("0xddfd32c4d33915685b926ba5eaab3860db1690cd"),
		common.HexToAddress("0xd338d81c4723982c815a294de3b38608dad9962c"),

		common.HexToAddress("0x6cd54fc6da0f044c43d4550d87ae10b9e1cea351"),
		common.HexToAddress("0xd328d8864649ed050b3d8e9d77f94c75299fd243"),
		common.HexToAddress("0x386dd85ad17b6bd60d2d142473b54bf9d5439842"),
		common.HexToAddress("0x4b8a6cff7b9e008caa936aadd33d9be048623d53"),
	}

	// TestnetChainConfig contains the chain parameters to run a node on the Ropsten test network.
	TestnetChainConfig = &ChainConfig{
		ChainID: big.NewInt(3),
		Dpos: &DposConfig{
			Period:       2,
			WitnessesNum: 19,
		},
	}

	// AllCliqueProtocolChanges contains every protocol change (EIPs) introduced
	// and accepted by the hubble core developers into the Clique consensus.
	//
	// This configuration is intentionally not using keyed fields to force anyone
	// adding flags to the config to also have to set these fields.
	AllCliqueProtocolChanges = &ChainConfig{
		big.NewInt(1337),
		big.NewInt(0),
		&DposConfig{
			Period:       2,
			WitnessesNum: 4,
		},
	}

	TestChainConfig = &ChainConfig{
		big.NewInt(1),
		big.NewInt(0),
		&DposConfig{
			Period:       2,
			WitnessesNum: 4,
		}}
	TestRules = TestChainConfig.Rules(new(big.Int))
)

// ChainConfig is the core config which determines the blockchain settings.
//
// ChainConfig is stored in the database on a per block basis. This means
// that any network, identified by its genesis block, can have its own
// set of configuration options.
type ChainConfig struct {
	ChainID *big.Int `json:"chainId"` // chainId identifies the current chain and is used for replay protection

	HubbleBlock *big.Int `json:"HubbleBlock,omitempty"` // Hubble switch block (nil = no fork, 0 = already hubble)

	// Various consensus engines
	Dpos *DposConfig `json:"dpos,omitempty"`
}

type DposConfig struct {
	Period       uint64   `json:"period"`       // Number of seconds between blocks to enforce
	WitnessesNum int      `json:"witnessesnum"` // Number of witnesses
	WitnessesUrl []string `json:"witnessesUrl"`
}

// String implements the stringer interface, returning the consensus engine details.
func (c *DposConfig) String() string {
	return "dpos"
}

// String implements the fmt.Stringer interface.
func (c *ChainConfig) String() string {
	var engine interface{}
	switch {
	case c.Dpos != nil:
		engine = c.Dpos
	default:
		engine = "unknown"
	}

	return fmt.Sprintf("{ChainID: %v Hubble: %v Engine: %v}",
		c.ChainID,
		c.HubbleBlock,
		engine,
	)
}

// IsHubble returns whether num is either equal to the hubble block or greater.
func (c *ChainConfig) IsHubble(num *big.Int) bool {
	return isForked(c.HubbleBlock, num)
}

// GasTable returns the gas table corresponding to the current phase .
//
// The returned GasTable's fields shouldn't, under any circumstances, be changed.
func (c *ChainConfig) GasTable(num *big.Int) GasTable {
	if num == nil {
		return GasTableHubble
	}
	switch {
	default:
		return GasTableHubble
	}
}

// CheckCompatible checks whether scheduled fork transitions have been imported
// with a mismatching chain configuration.
func (c *ChainConfig) CheckCompatible(newcfg *ChainConfig, height uint64) *ConfigCompatError {
	bhead := new(big.Int).SetUint64(height)

	// Iterate checkCompatible to find the lowest conflict.
	var lasterr *ConfigCompatError
	for {
		err := c.checkCompatible(newcfg, bhead)
		if err == nil || (lasterr != nil && err.RewindTo == lasterr.RewindTo) {
			break
		}
		lasterr = err
		bhead.SetUint64(err.RewindTo)
	}
	return lasterr
}

func (c *ChainConfig) checkCompatible(newcfg *ChainConfig, head *big.Int) *ConfigCompatError {
	if isForkIncompatible(c.HubbleBlock, newcfg.HubbleBlock, head) {
		return newCompatError("Hubble fork block", c.HubbleBlock, newcfg.HubbleBlock)
	}
	return nil
}

// isForkIncompatible returns true if a fork scheduled at s1 cannot be rescheduled to
// block s2 because head is already past the fork.
func isForkIncompatible(s1, s2, head *big.Int) bool {
	return (isForked(s1, head) || isForked(s2, head)) && !configNumEqual(s1, s2)
}

// isForked returns whether a fork scheduled at block s is active at the given head block.
func isForked(s, head *big.Int) bool {
	if s == nil || head == nil {
		return false
	}
	return s.Cmp(head) <= 0
}

func configNumEqual(x, y *big.Int) bool {
	if x == nil {
		return y == nil
	}
	if y == nil {
		return x == nil
	}
	return x.Cmp(y) == 0
}

// ConfigCompatError is raised if the locally-stored blockchain is initialised with a
// ChainConfig that would alter the past.
type ConfigCompatError struct {
	What string
	// block numbers of the stored and new configurations
	StoredConfig, NewConfig *big.Int
	// the block number to which the local chain must be rewound to correct the error
	RewindTo uint64
}

func newCompatError(what string, storedblock, newblock *big.Int) *ConfigCompatError {
	var rew *big.Int
	switch {
	case storedblock == nil:
		rew = newblock
	case newblock == nil || storedblock.Cmp(newblock) < 0:
		rew = storedblock
	default:
		rew = newblock
	}
	err := &ConfigCompatError{what, storedblock, newblock, 0}
	if rew != nil && rew.Sign() > 0 {
		err.RewindTo = rew.Uint64() - 1
	}
	return err
}

func (err *ConfigCompatError) Error() string {
	return fmt.Sprintf("mismatching %s in database (have %d, want %d, rewindto %d)", err.What, err.StoredConfig, err.NewConfig, err.RewindTo)
}

// Rules wraps ChainConfig and is merely syntatic sugar or can be used for functions
// that do not have or require information about the block.
//
// Rules is a one time interface meaning that it shouldn't be used in between transition
// phases.
type Rules struct {
	ChainID  *big.Int
	IsHubble bool
}

// Rules ensures c's ChainID is not nil.
func (c *ChainConfig) Rules(num *big.Int) Rules {
	chainID := c.ChainID
	if chainID == nil {
		chainID = new(big.Int)
	}
	return Rules{ChainID: new(big.Int).Set(chainID), IsHubble: c.IsHubble(num)}
}
