// Copyright 2019 The go-vnt Authors
// This file is part of the go-vnt library.
//
// The go-vnt library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-vnt library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-vnt library. If not, see <http://www.gnu.org/licenses/>.

package abi

import (
	"encoding/json"
	"fmt"
)

type Table struct {
	Name        string
	Type        Type
	StorageType uint64
	Tables      Tables
}

type Tables []Table

type Extarg struct {
	Name        string
	Type        string
	StorageType uint64
	Tables      []Extarg
}

// UnmarshalJSON implements json.Unmarshaler interface
func (table *Table) UnmarshalJSON(data []byte) error {
	var ext Extarg
	err := json.Unmarshal(data, &ext)
	if err != nil {
		return fmt.Errorf("argument json err: %v", err)
	}

	fmt.Printf("========ext %+v\n", ext)
	err = ext.recursive(table)
	if err != nil {
		return err
	}

	return nil
}

func (ext *Extarg) recursive(table *Table) error {
	fmt.Printf("recursive %v\n", ext)
	if ext.Type == "" {
		table.Type = Type{}
	} else {
		var err error
		table.Type, err = NewType(ext.Type)
		if err != nil {
			return err
		}
	}
	table.Name = ext.Name
	table.StorageType = ext.StorageType

	// table.Tables = ext.Name
	if len(ext.Tables) != 0 {
		for _, v := range ext.Tables {
			t := &Table{}
			err := v.recursive(t)
			if err != nil {
				return err
			}
			table.Tables = append(table.Tables, *t)
		}
	}
	return nil
}

func (tbl Table) Traversal(sym string, key *Key) {
	if len(tbl.Tables) == 0 {
		key.Types[sym] = tbl.Type
		key.Keys[sym] = tbl
		fmt.Printf("######sym %s ######\n", sym)
	} else {
		for _, v := range tbl.Tables {
			s := fmt.Sprintf("%s.%s", sym, v.Name)
			v.Traversal(s, key)
		}
	}
}
