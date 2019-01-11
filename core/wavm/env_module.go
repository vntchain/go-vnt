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
