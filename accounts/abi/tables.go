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
