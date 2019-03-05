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

package wavm

import (
	"github.com/vntchain/vnt-wasm/wasm"
)

type EnvModule struct {
	functions EnvFunctions
	module    wasm.Module
}

func (m *EnvModule) InitModule(ctx *ChainContext) {
	m.functions = EnvFunctions{}
	m.functions.InitFuncTable(ctx)

	funcTable := m.functions.GetFuncTable()

	m.module = wasm.Module{
		FunctionIndexSpace: make([]wasm.Function, 0),
		Export: &wasm.SectionExports{
			Entries: make(map[string]wasm.ExportEntry),
		},
	}

	index := uint32(0)
	for name, function := range funcTable {
		m.module.Export.Entries[name] = wasm.ExportEntry{
			FieldStr: name,
			Kind:     wasm.ExternalFunction,
			Index:    index,
		}
		m.module.FunctionIndexSpace = append(m.module.FunctionIndexSpace, function)
		index++
	}
}

func (m *EnvModule) GetModule() *wasm.Module {
	return &m.module
}

func (m *EnvModule) GetEnvFunctions() EnvFunctions {
	return m.functions
}
