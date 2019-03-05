// Copyright 2019 The go-vnt Authors
// This file is part of go-vnt.
//
// go-vnt is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-vnt is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-vnt. If not, see <http://www.gnu.org/licenses/>.

package vnt

// Target is the "target" of a br_table instruction.
// Unlike other control instructions, br_table does jumps and discarding all
// by itself.
type Target struct {
	Addr        int64 // The absolute address of the target
	Discard     int64 // The number of elements to discard
	PreserveTop bool  // Whether the top of the stack is to be preserved
	Return      bool  // Whether to return in order to take this branch/target
}

// BranchTable is the structure pointed to by a rewritten br_table instruction.
// A rewritten br_table instruction is of the format:
//     br_table <table_index>
// where <table_index> is the index to an array of
// BranchTable objects stored by the VM.
type BranchTable struct {
	Targets       []Target // A list of targets, br_table pops an int value, and jumps to Targets[val]
	DefaultTarget Target   // If val > len(Targets), the VM will jump here
	PatchedAddrs  []int64  // A list of already patched addresses
	BlocksLen     int      // The length of the blocks map in Compile when this table was initialized
}

func (table *BranchTable) PatchTable(block int, addr int64) {
	if block < 0 {
		panic("Invalid block value")
	}

	for i, target := range table.Targets {
		if !table.isAddr(target.Addr) && target.Addr == int64(block) {
			table.Targets[i].Addr = addr
		}
	}

	if table.DefaultTarget.Addr == int64(block) {
		table.DefaultTarget.Addr = addr
	}
	table.PatchedAddrs = append(table.PatchedAddrs, addr)
}

// Whether the given value is an instruction (or the block depth)
func (table *BranchTable) isAddr(addr int64) bool {
	for _, t := range table.PatchedAddrs {
		if t == addr {
			return true
		}
	}
	return false
}

type Compiled struct {
	Code           []byte
	Table          []*BranchTable
	MaxDepth       int
	TotalLocalVars int
}