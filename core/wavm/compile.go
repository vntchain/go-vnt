// Copyright 2017 The go-interpreter Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package compile is used internally by wagon to convert standard structured
// WebAssembly bytecode into an unstructured form suitable for execution by
// it's VM.
// The conversion process consists of translating block instruction sequences
// and branch operators (br, br_if, br_table) to absolute jumps to PC values.
// For instance, an instruction sequence like:
//     loop
//       i32.const 1
//       get_local 0
//       i32.add
//       set_local 0
//       get_local 1
//       i32.const 1
//       i32.add
//       tee_local 1
//       get_local 2
//       i32.eq
//       br_if 0
//     end
// Is "compiled" to:
//     i32.const 1
//     i32.add
//     set_local 0
//     get_local 1
//     i32.const 1
//     i32.add
//     tee_local 1
//     get_local 2
//     i32.eq
//     jmpnz <addr> <preserve> <discard>
// Where jmpnz is a jump-if-not-zero operator that takes certain arguments
// plus the jump address as immediates.
// This is in contrast with original WebAssembly bytecode, where the target
// of branch operators are relative block depths instead.
package wavm

import (
	"bytes"
	"encoding/binary"

	"github.com/vntchain/go-vnt/core/wavm/gas"
	"github.com/vntchain/go-vnt/core/wavm/utils"
	"github.com/vntchain/go-vnt/log"
	"github.com/vntchain/vnt-wasm/disasm"
	"github.com/vntchain/vnt-wasm/vnt"
	"github.com/vntchain/vnt-wasm/wasm"
	ops "github.com/vntchain/vnt-wasm/wasm/operators"
)

// A small note on the usage of discard instructions:
// A control operator sequence isn't allowed to access nor modify (pop) operands
// that were pushed outside it. Therefore, each sequence has its own stack
// that may or may not push a value to the original stack, depending on the
// block's signature.
// Instead of creating a new stack every time we enter a control structure,
// we record the current stack height on encountering a control operator.
// After we leave the sequence, the stack height is restored using the discard
// operator. A block with a signature will push a value of that type on the parent
// stack (that is, the stack of the parent block where this block started). The
// OpDiscardPreserveTop operator allows us to preserve this value while
// discarding the remaining ones.

// Branches are rewritten as
//     <jmp> <addr>
// Where the address is an 8 byte address, initially set to zero. It is
// later "patched" by patchOffset.

var (
	// OpJmp unconditionally jumps to the provided address.
	OpJmp byte = 0x0c
	// OpJmpZ jumps to the given address if the value at the top of the stack is zero.
	OpJmpZ byte = 0x03
	// OpJmpNz jumps to the given address if the value at the top of the
	// stack is not zero. It also discards elements and optionally preserves
	// the topmost value on the stack
	OpJmpNz byte = 0x0d
	// OpDiscard discards a given number of elements from the execution stack.
	OpDiscard byte = 0x0b
	// OpDiscardPreserveTop discards a given number of elements from the
	// execution stack, while preserving the value on the top of the stack.
	OpDiscardPreserveTop byte = 0x05
)

// block stores the information relevant for a block created by a control operator
// sequence (if...else...end, loop...end, and block...end)
type block struct {
	// the byte offset to which the continuation of the label
	// created by the block operator is located
	// for 'loop', this is the offset of the loop operator itself
	// for 'if', 'else', 'block', this is the 'end' operator
	offset int64

	// Whether this block is created by an 'if' operator
	// in that case, the 'offset' field is set to the byte offset
	// of the else branch, once the else operator is reached.
	ifBlock bool
	// if ... else ... end is compiled to
	// jmpnz <else-addr> ... jmp <end-addr> ... <discard>
	// elseAddrOffset is the byte offset of the else-addr address
	// in the new/compiled byte buffer.
	elseAddrOffset int64

	// Whether this block is created by a 'loop' operator
	// in that case, the 'offset' field is set at the end of the block
	loopBlock bool

	patchOffsets []int64 // A list of offsets in the bytecode stream that need to be patched with the correct jump addresses

	discard      disasm.StackInfo   // Information about the stack created in this block, used while creating Discard instructions
	branchTables []*vnt.BranchTable // All branch tables that were defined in this block.
}

type Mutable map[uint32]bool

type CodeBlock struct {
	stack map[int]*Stack
}

type Code struct {
	Body     disasm.Instr
	Children []*Code
}

func (c Code) Recursive() []disasm.Instr {
	if len(c.Children) == 0 {
		return []disasm.Instr{c.Body}
	} else {
		merge := []disasm.Instr{c.Body}
		for _, v := range c.Children {
			child := v.Body
			tmp := merge
			merge = append([]disasm.Instr{child}, tmp...)
		}
		return merge
	}
}

// func (c Code) String() string {
// 	if len(c.Children) == 0 {
// 		return fmt.Sprintf("( %s )", c.Body)
// 	} else {
// 		merge := fmt.Sprintf("( %s", c.Body)
// 		for _, v := range c.Children {
// 			child := fmt.Sprintf("%s", v.String())
// 			merge = fmt.Sprintf("%s\n%s", merge, child)
// 		}
// 		merge = fmt.Sprintf("%s\n)", merge)
// 		return merge
// 	}
// }

type Stack struct {
	slice []*Code
}

func (s *Stack) Push(b *Code) {
	s.slice = append(s.slice, b)
}

func (s *Stack) Pop() *Code {
	v := s.Top()
	s.slice = s.slice[:len(s.slice)-1]
	return v
}

func (s *Stack) Top() *Code {
	return s.slice[len(s.slice)-1]
}

func (s Stack) Len() int {
	return len(s.slice)
}

func (cb *CodeBlock) buildCode(blockDepth int, n int) *Code {
	code := &Code{}
	stack := cb.stack[blockDepth]
	for i := 0; i < n; i++ {
		pop := stack.Pop()
		tmp := code.Children
		code.Children = append([]*Code{pop}, tmp...)
		//w.code.Children = append(w.code.Children, pop)
	}

	return code
}

func CompileModule(module *wasm.Module, chainctx ChainContext, mutable Mutable) ([]vnt.Compiled, error) {
	Compiled := make([]vnt.Compiled, len(module.FunctionIndexSpace))
	for i, fn := range module.FunctionIndexSpace {
		// Skip native methods as they need not be
		// disassembled; simply add them at the end
		// of the `funcs` array as is, as specified
		// in the spec. See the "host functions"
		// section of:
		// https://webassembly.github.io/spec/core/exec/modules.html#allocation
		if fn.IsHost() {
			continue
		}

		var code []byte
		var table []*vnt.BranchTable
		var maxDepth int
		totalLocalVars := 0

		disassembly, err := disasm.Disassemble(fn, module)
		if err != nil {
			return nil, err
		}

		maxDepth = disassembly.MaxDepth

		totalLocalVars += len(fn.Sig.ParamTypes)
		for _, entry := range fn.Body.Locals {
			totalLocalVars += int(entry.Count)
		}
		disassembly.Code = gas.InjectCounter(disassembly.Code, module, chainctx.GasRule)
		code, table = Compile(disassembly.Code, module, mutable)
		Compiled[i] = vnt.Compiled{
			Code:           code,
			Table:          table,
			MaxDepth:       maxDepth,
			TotalLocalVars: totalLocalVars,
		}
	}
	return Compiled, nil
}

// func (cb *CodeBlock) addChild() {
// 	cb.code.Children = append([]code)
// }

// Compile rewrites WebAssembly bytecode from its disassembly.
// TODO(vibhavp): Add options for optimizing code. Operators like i32.reinterpret/f32
// are no-ops, and can be safely removed.
func Compile(disassembly []disasm.Instr, module *wasm.Module, mutable Mutable) ([]byte, []*vnt.BranchTable) {
	buffer := new(bytes.Buffer)
	branchTables := []*vnt.BranchTable{}

	curBlockDepth := -1
	blocks := make(map[int]*block) // maps nesting depths (labels) to blocks

	blocks[-1] = &block{}

	writeIndex, readIndex, _ := utils.GetIndex(module)
	codeBlock := &CodeBlock{stack: map[int]*Stack{}}

	newInstr := []disasm.Instr{}

	for _, instr := range disassembly {
		// fmt.Printf("compile instr %+v blockinfo %+v\n", instr, instr.Block)
		var readInstr []disasm.Instr
		var writeInstr []disasm.Instr
		if instr.Unreachable {
			continue
		}
		if codeBlock.stack[curBlockDepth] == nil {
			codeBlock.stack[curBlockDepth] = &Stack{}
		}
		switch instr.Op.Code {
		case ops.I32Const, ops.I64Const, ops.F32Const, ops.F64Const:
			codeBlock.stack[curBlockDepth].Push(&Code{Body: instr})
		case ops.I32Add, ops.I32Sub, ops.I32Mul, ops.I32DivS, ops.I32DivU, ops.I32RemS, ops.I32RemU, ops.I32And, ops.I32Or, ops.I32Xor, ops.I32Shl, ops.I32ShrS, ops.I32ShrU, ops.I32Rotl, ops.I32Rotr,
			ops.I32Eq, ops.I32Ne, ops.I32LtS, ops.I32LtU, ops.I32LeS, ops.I32LeU, ops.I32GtS, ops.I32GtU, ops.I32GeS, ops.I32GeU,
			ops.I64Add, ops.I64Sub, ops.I64Mul, ops.I64DivS, ops.I64DivU, ops.I64RemS, ops.I64RemU, ops.I64And, ops.I64Or, ops.I64Xor, ops.I64Shl, ops.I64ShrS, ops.I64ShrU, ops.I64Rotl, ops.I64Rotr,
			ops.I64Eq, ops.I64Ne, ops.I64LtS, ops.I64LtU, ops.I64LeS, ops.I64LeU, ops.I64GtS, ops.I64GtU, ops.I64GeS, ops.I64GeU,
			ops.F32Add, ops.F32Sub, ops.F32Mul, ops.F32Div, ops.F32Min, ops.F32Max, ops.F32Copysign,
			ops.F32Eq, ops.F32Ne, ops.F32Lt, ops.F32Le, ops.F32Gt, ops.F32Ge,
			ops.F64Add, ops.F64Sub, ops.F64Mul, ops.F64Div, ops.F64Min, ops.F64Max, ops.F64Copysign,
			ops.F64Eq, ops.F64Ne, ops.F64Lt, ops.F64Le, ops.F64Gt, ops.F64Ge:
			code := codeBlock.buildCode(curBlockDepth, 2)
			code.Body = instr
			codeBlock.stack[curBlockDepth].Push(code)
		case ops.I32Clz, ops.I32Ctz, ops.I32Popcnt, ops.I32Eqz,
			ops.I64Clz, ops.I64Ctz, ops.I64Popcnt, ops.I64Eqz,
			ops.F32Sqrt, ops.F32Ceil, ops.F32Floor, ops.F32Trunc, ops.F32Nearest, ops.F32Abs, ops.F32Neg,
			ops.F64Sqrt, ops.F64Ceil, ops.F64Floor, ops.F64Trunc, ops.F64Nearest, ops.F64Abs, ops.F64Neg,
			ops.I32WrapI64, ops.I64ExtendUI32, ops.I64ExtendSI32,
			ops.I32TruncUF32, ops.I32TruncUF64, ops.I64TruncUF32, ops.I64TruncUF64,
			ops.I32TruncSF32, ops.I32TruncSF64, ops.I64TruncSF32, ops.I64TruncSF64,
			ops.F32DemoteF64, ops.F64PromoteF32,
			ops.F32ConvertUI32, ops.F32ConvertUI64, ops.F64ConvertUI32, ops.F64ConvertUI64,
			ops.F32ConvertSI32, ops.F32ConvertSI64, ops.F64ConvertSI32, ops.F64ConvertSI64,
			ops.I32ReinterpretF32, ops.I64ReinterpretF64,
			ops.F32ReinterpretI32, ops.F64ReinterpretI64:
			code := codeBlock.buildCode(curBlockDepth, 1)
			code.Body = instr
			codeBlock.stack[curBlockDepth].Push(code)
		case ops.Drop:
			code := codeBlock.buildCode(curBlockDepth, 1)
			code.Body = instr
		case ops.GetLocal, ops.GetGlobal:
			codeBlock.stack[curBlockDepth].Push(&Code{Body: instr})
		case ops.SetLocal, ops.SetGlobal:
			code := codeBlock.buildCode(curBlockDepth, 1)
			code.Body = instr
		case ops.TeeLocal:
			code := codeBlock.buildCode(curBlockDepth, 1)
			code.Body = instr
			codeBlock.stack[curBlockDepth].Push(code)
		case ops.I32Load, ops.I64Load, ops.F32Load, ops.F64Load, ops.I32Load8s, ops.I32Load8u, ops.I32Load16s, ops.I32Load16u, ops.I64Load8s, ops.I64Load8u, ops.I64Load16s, ops.I64Load16u, ops.I64Load32s, ops.I64Load32u:
			// memory_immediate has two fields, the alignment and the offset.
			// The former is simply an optimization hint and can be safely
			// discarded.
			instr.Immediates = []interface{}{instr.Immediates[1].(uint32)}

			arg := codeBlock.stack[curBlockDepth].slice[codeBlock.stack[curBlockDepth].Len()-1]
			if arg.Body.Op.Code == ops.I32Const {
				constBaseInstr := arg.Body
				constOffsetOp, _ := ops.New(ops.I32Const)
				constInstr := disasm.Instr{Op: constOffsetOp, Immediates: []interface{}{int32(instr.Immediates[0].(uint32))}}
				callOp, _ := ops.New(ops.Call)
				callInstr := disasm.Instr{Op: callOp, Immediates: []interface{}{uint32(readIndex)}}
				readInstr = []disasm.Instr{constBaseInstr, constInstr, callInstr}
			}
			code := codeBlock.buildCode(curBlockDepth, 1)
			code.Body = instr
			codeBlock.stack[curBlockDepth].Push(code)
		case ops.I32Store, ops.I64Store, ops.F32Store, ops.F64Store, ops.I32Store8, ops.I32Store16, ops.I64Store8, ops.I64Store16, ops.I64Store32:
			// memory_immediate has two fields, the alignment and the offset.
			// The former is simply an optimization hint and can be safely
			// discarded.
			instr.Immediates = []interface{}{instr.Immediates[1].(uint32)}

			arg := codeBlock.stack[curBlockDepth].slice[codeBlock.stack[curBlockDepth].Len()-2]
			if arg.Body.Op.Code == ops.I32Const {
				constBaseInstr := arg.Body
				constOffsetOp, _ := ops.New(ops.I32Const)
				constoffsetInstr := disasm.Instr{Op: constOffsetOp, Immediates: []interface{}{int32(instr.Immediates[0].(uint32))}}
				callOp, _ := ops.New(ops.Call)
				callInstr := disasm.Instr{Op: callOp, Immediates: []interface{}{uint32(writeIndex)}}
				writeInstr = []disasm.Instr{constBaseInstr, constoffsetInstr, callInstr}
			}
			code := codeBlock.buildCode(curBlockDepth, 2)
			code.Body = instr
		case ops.Call, ops.CallIndirect:
			index := instr.Immediates[0].(uint32)
			sig := module.GetFunction(int(index)).Sig
			if instr.Op.Code == ops.CallIndirect {
				sig = &module.Types.Entries[int(index)]
			}
			parms := len(sig.ParamTypes)
			returns := len(sig.ReturnTypes)
			code := codeBlock.buildCode(curBlockDepth, parms)
			code.Body = instr
			//codeBlock.stack.Push(codeBlock.code)
			if returns != 0 {
				codeBlock.stack[curBlockDepth].Push(code)
			}
		case ops.If:
			curBlockDepth++
			buffer.WriteByte(OpJmpZ)
			blocks[curBlockDepth] = &block{
				ifBlock:        true,
				elseAddrOffset: int64(buffer.Len()),
			}
			// the address to jump to if the condition for `if` is false
			// (i.e when the value on the top of the stack is 0)
			binary.Write(buffer, binary.LittleEndian, int64(0))

			op, err := ops.New(OpJmpZ)
			if err != nil {
				panic(err)
			}
			ins := disasm.Instr{
				Op:         op,
				Immediates: [](interface{}){},
			}
			ins.Immediates = append(ins.Immediates, int64(0))
			newInstr = append(newInstr, ins)

			sig := instr.Immediates[0].(wasm.BlockType)
			code := codeBlock.buildCode(curBlockDepth-1, 1)
			if sig != wasm.BlockTypeEmpty {
				code.Body = instr
				codeBlock.stack[curBlockDepth-1].Push(code)
			}
			// else {
			// 	if curBlockDepth == 0 {
			// 		code := &Code{Body: instr}
			// 		codeBlock.stack[curBlockDepth].Push(code)
			// 	} else {
			// 		parentCode := codeBlock.stack[curBlockDepth-1].Top()
			// 		code := &Code{Body: instr}
			// 		parentCode.Children = append(parentCode.Children, code)
			// 	}
			// }
			continue
		case ops.Loop:
			// there is no condition for entering a loop block
			curBlockDepth++
			blocks[curBlockDepth] = &block{
				offset:    int64(buffer.Len()),
				ifBlock:   false,
				loopBlock: true,
				discard:   *instr.NewStack,
			}

			sig := instr.Immediates[0].(wasm.BlockType)
			if sig != wasm.BlockTypeEmpty {
				code := &Code{Body: instr}
				codeBlock.stack[curBlockDepth-1].Push(code)
			}
			// else {
			// 	if curBlockDepth == 0 {
			// 		code := &Code{Body: instr}
			// 		codeBlock.stack[curBlockDepth].Push(code)
			// 	} else {
			// 		parentCode := codeBlock.stack[curBlockDepth-1].Top()
			// 		code := &Code{Body: instr}
			// 		parentCode.Children = append(parentCode.Children, code)
			// 	}
			// }

			continue
		case ops.Block:
			curBlockDepth++
			blocks[curBlockDepth] = &block{
				ifBlock: false,
				discard: *instr.NewStack,
			}

			sig := instr.Immediates[0].(wasm.BlockType)
			if sig != wasm.BlockTypeEmpty {
				code := &Code{Body: instr}
				codeBlock.stack[curBlockDepth-1].Push(code)
			}
			// else {
			// 	if curBlockDepth == 0 {
			// 		code := &Code{Body: instr}
			// 		codeBlock.stack[curBlockDepth].Push(code)
			// 	} else {
			// 		parentCode := codeBlock.stack[curBlockDepth-1].Top()
			// 		code := &Code{Body: instr}
			// 		parentCode.Children = append(parentCode.Children, code)
			// 	}
			// }
			continue
		case ops.Else:
			ifInstr := disassembly[instr.Block.ElseIfIndex] // the corresponding `if` instruction for this else
			if ifInstr.NewStack != nil && ifInstr.NewStack.StackTopDiff != 0 {
				// add code for jumping out of a taken if branch
				if ifInstr.NewStack.PreserveTop {
					buffer.WriteByte(OpDiscardPreserveTop)
				} else {
					buffer.WriteByte(OpDiscard)
				}
				binary.Write(buffer, binary.LittleEndian, ifInstr.NewStack.StackTopDiff)
			}
			buffer.WriteByte(OpJmp)
			ifBlockEndOffset := int64(buffer.Len())
			binary.Write(buffer, binary.LittleEndian, int64(0))

			curOffset := int64(buffer.Len())
			ifBlock := blocks[curBlockDepth]
			code := buffer.Bytes()

			buffer = patchOffset(code, ifBlock.elseAddrOffset, curOffset)
			// this is no longer an if block
			ifBlock.ifBlock = false
			ifBlock.patchOffsets = append(ifBlock.patchOffsets, ifBlockEndOffset)

			op, err := ops.New(OpJmp)
			if err != nil {
				panic(err)
			}
			ins := disasm.Instr{
				Op:         op,
				Immediates: [](interface{}){},
			}
			ins.Immediates = append(ins.Immediates, int64(0))
			newInstr = append(newInstr, ins)

			continue
		case ops.End:
			depth := curBlockDepth
			block := blocks[depth]

			if instr.NewStack.StackTopDiff != 0 {
				// when exiting a block, discard elements to
				// restore stack height.
				var op ops.Op
				var err error
				var ins disasm.Instr
				if instr.NewStack.PreserveTop {
					// this is true when the block has a
					// signature, and therefore pushes
					// a value on to the stack
					buffer.WriteByte(OpDiscardPreserveTop)

					op, err = ops.New(OpDiscardPreserveTop)
					if err != nil {
						panic(err)
					}
					ins = disasm.Instr{
						Op:         op,
						Immediates: [](interface{}){},
					}

				} else {
					buffer.WriteByte(OpDiscard)

					op, err = ops.New(OpDiscard)
					if err != nil {
						panic(err)
					}
					ins = disasm.Instr{
						Op:         op,
						Immediates: [](interface{}){},
					}
				}
				binary.Write(buffer, binary.LittleEndian, instr.NewStack.StackTopDiff)

				ins.Immediates = append(ins.Immediates, instr.NewStack.StackTopDiff)
				newInstr = append(newInstr, ins)
			}

			if !block.loopBlock { // is a normal block
				block.offset = int64(buffer.Len())
				if block.ifBlock {
					code := buffer.Bytes()
					buffer = patchOffset(code, block.elseAddrOffset, int64(block.offset))
				}
			}

			for _, offset := range block.patchOffsets {
				code := buffer.Bytes()
				buffer = patchOffset(code, offset, block.offset)
			}

			for _, table := range block.branchTables {
				table.PatchTable(table.BlocksLen-depth-1, int64(block.offset))
			}

			delete(blocks, curBlockDepth)
			curBlockDepth--
			continue
		case ops.Br:
			if instr.NewStack != nil && instr.NewStack.StackTopDiff != 0 {
				var op ops.Op
				var err error
				var ins disasm.Instr
				if instr.NewStack.PreserveTop {
					buffer.WriteByte(OpDiscardPreserveTop)

					op, err = ops.New(OpDiscardPreserveTop)
					if err != nil {
						panic(err)
					}
					ins = disasm.Instr{
						Op:         op,
						Immediates: [](interface{}){},
					}

				} else {
					buffer.WriteByte(OpDiscard)

					op, err = ops.New(OpDiscard)
					if err != nil {
						panic(err)
					}
					ins = disasm.Instr{
						Op:         op,
						Immediates: [](interface{}){},
					}
				}
				binary.Write(buffer, binary.LittleEndian, instr.NewStack.StackTopDiff)

				ins.Immediates = append(ins.Immediates, instr.NewStack.StackTopDiff)
				newInstr = append(newInstr, ins)
			}
			buffer.WriteByte(OpJmp)
			label := int(instr.Immediates[0].(uint32))
			block := blocks[curBlockDepth-int(label)]
			block.patchOffsets = append(block.patchOffsets, int64(buffer.Len()))
			// write the jump address
			binary.Write(buffer, binary.LittleEndian, int64(0))

			op, err := ops.New(OpJmp)
			if err != nil {
				panic(err)
			}
			ins := disasm.Instr{
				Op:         op,
				Immediates: [](interface{}){},
			}
			ins.Immediates = append(ins.Immediates, int64(0))
			newInstr = append(newInstr, ins)
			continue
		case ops.BrIf:
			buffer.WriteByte(OpJmpNz)
			label := int(instr.Immediates[0].(uint32))
			block := blocks[curBlockDepth-int(label)]
			block.patchOffsets = append(block.patchOffsets, int64(buffer.Len()))
			// write the jump address
			binary.Write(buffer, binary.LittleEndian, int64(0))

			op, err := ops.New(OpJmpNz)
			if err != nil {
				panic(err)
			}
			ins := disasm.Instr{
				Op:         op,
				Immediates: [](interface{}){},
			}

			var stackTopDiff int64
			// write whether we need to preserve the top
			if instr.NewStack == nil || !instr.NewStack.PreserveTop || instr.NewStack.StackTopDiff == 0 {
				buffer.WriteByte(byte(0))
				ins.Immediates = append(ins.Immediates, false)
			} else {
				stackTopDiff = instr.NewStack.StackTopDiff
				buffer.WriteByte(byte(1))
				ins.Immediates = append(ins.Immediates, true)
			}
			// write the number of elements on the stack we need to discard
			binary.Write(buffer, binary.LittleEndian, stackTopDiff)

			ins.Immediates = append(ins.Immediates, int64(0))
			ins.Immediates = append(ins.Immediates, stackTopDiff)
			newInstr = append(newInstr, ins)
			continue
		case ops.BrTable:
			branchTable := &vnt.BranchTable{
				// we subtract one for the implicit block created by
				// the function body
				BlocksLen: len(blocks) - 1,
			}
			targetCount := instr.Immediates[0].(uint32)
			branchTable.Targets = make([]vnt.Target, targetCount)
			for i := range branchTable.Targets {
				// The first immediates is the number of targets, so we ignore that
				label := int64(instr.Immediates[i+1].(uint32))
				branchTable.Targets[i].Addr = label
				branch := instr.Branches[i]

				branchTable.Targets[i].Return = branch.IsReturn
				branchTable.Targets[i].Discard = branch.StackTopDiff
				branchTable.Targets[i].PreserveTop = branch.PreserveTop
			}
			defaultLabel := int64(instr.Immediates[len(instr.Immediates)-1].(uint32))
			branchTable.DefaultTarget.Addr = defaultLabel
			defaultBranch := instr.Branches[targetCount]
			branchTable.DefaultTarget.Return = defaultBranch.IsReturn
			branchTable.DefaultTarget.Discard = defaultBranch.StackTopDiff
			branchTable.DefaultTarget.PreserveTop = defaultBranch.PreserveTop
			branchTables = append(branchTables, branchTable)
			for _, block := range blocks {
				block.branchTables = append(block.branchTables, branchTable)
			}

			buffer.WriteByte(ops.BrTable)
			binary.Write(buffer, binary.LittleEndian, int64(len(branchTables)-1))

			op, err := ops.New(ops.BrTable)
			if err != nil {
				panic(err)
			}
			ins := disasm.Instr{
				Op:         op,
				Immediates: [](interface{}){},
			}
			ins.Immediates = append(ins.Immediates, int64(len(branchTables)-1))
			newInstr = append(newInstr, ins)
		}
		if len(readInstr) != 0 {
			if readIndex != -1 {
				for _, instr := range readInstr {
					buffer.WriteByte(instr.Op.Code)
					for _, imm := range instr.Immediates {
						err := binary.Write(buffer, binary.LittleEndian, imm)
						if err != nil {
							panic(err)
						}
					}
				}
				newInstr = append(newInstr, readInstr...)
			} else {
				log.Warn("Compile warning", "Msg", "Can't find ReadWithPointer env function!!")
			}
		}
		buffer.WriteByte(instr.Op.Code)
		for _, imm := range instr.Immediates {
			err := binary.Write(buffer, binary.LittleEndian, imm)
			if err != nil {
				panic(err)
			}
		}
		newInstr = append(newInstr, instr)
		if len(writeInstr) != 0 {
			if writeIndex != -1 {
				for _, instr := range writeInstr {
					buffer.WriteByte(instr.Op.Code)
					for _, imm := range instr.Immediates {
						err := binary.Write(buffer, binary.LittleEndian, imm)
						if err != nil {
							panic(err)
						}
					}
				}
				newInstr = append(newInstr, writeInstr...)
			} else {
				log.Warn("Compile warning", "Msg", "Can't find WriteWithPointer env function!!")
			}
		}
	}

	// writing nop as the last instructions allows us to branch out of the
	// function (ie, return)
	addr := buffer.Len()
	buffer.WriteByte(ops.Nop)

	// patch all references to the "root" block of the function body
	for _, offset := range blocks[-1].patchOffsets {
		code := buffer.Bytes()
		buffer = patchOffset(code, offset, int64(addr))
	}

	for _, table := range branchTables {
		table.PatchedAddrs = nil
	}
	return buffer.Bytes(), branchTables
}

// replace the address starting at start with addr
func patchOffset(code []byte, start int64, addr int64) *bytes.Buffer {
	var shift uint
	for i := int64(0); i < 8; i++ {
		code[start+i] = byte(addr >> shift)
		shift += 8
	}

	buf := new(bytes.Buffer)
	buf.Write(code)
	return buf
}
