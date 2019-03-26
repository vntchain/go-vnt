package wavm

// OpCode is an EVM opcode
type OpCode byte

func (op OpCode) IsPush() bool {
	return false
}

func (op OpCode) String() string {
	return ""
}

func (op OpCode) Byte() byte {
	return byte(op)
}
