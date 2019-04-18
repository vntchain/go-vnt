package wavm

import ops "github.com/vntchain/vnt-wasm/wasm/operators"

// OpCode is an WAVM opcode

var (
	OpReturn  byte = 0x0f
	OpBrTable byte = 0x0e
)

type OpCode struct {
	Op       byte
	FuncName string
}

func (op OpCode) IsPush() bool {
	return false
}

func (op OpCode) String() string {
	if op.FuncName == "" {
		switch op.Op {
		case OpReturn:
			return "Return"
		case OpJmp:
			return "Jmp"
		case OpJmpZ:
			return "JmpZ"
		case OpJmpNz:
			return "JmpNz"
		case OpBrTable:
			return "BrTable"
		case OpDiscard:
			return "Discard"
		case OpDiscardPreserveTop:
			return "DiscardPreserveTop"
		default:
			o, err := ops.New(op.Op)
			if err != nil {
				return ""
			}
			return o.Name
		}
	} else {
		return op.FuncName
	}
}

func (op OpCode) Byte() byte {
	return op.Op
}
