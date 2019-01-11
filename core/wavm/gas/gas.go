package gas

import (
	"sort"

	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core/wavm/internal/stack"
	"github.com/vntchain/go-vnt/core/wavm/utils"
	"github.com/vntchain/vnt-wasm/disasm"
	"github.com/vntchain/vnt-wasm/wasm"
	ops "github.com/vntchain/vnt-wasm/wasm/operators"
)

type BlockEntry struct {
	/// Index of the first instruction (aka `Opcode`) in the block.
	starPos int
	/// Sum of costs of all instructions until end of the block.
	cost  uint64
	index []int
}

type BlockEntrys []*BlockEntry

func (block BlockEntrys) Len() int           { return len(block) }
func (block BlockEntrys) Less(i, j int) bool { return block[i].starPos-block[j].starPos < 0 }
func (block BlockEntrys) Swap(i, j int)      { block[i], block[j] = block[j], block[i] }

type Counter struct {
	/// All blocks in the order of theirs start position.
	blocks BlockEntrys

	// Stack of blocks. Each element is an index to a `self.blocks` vector.
	stack *stack.Stack
}

func NewCounter() *Counter {
	return &Counter{
		blocks: BlockEntrys{},
		stack:  &stack.Stack{},
	}
}

func (counter *Counter) Begin(cursor int) {
	blockIdx := len(counter.blocks)
	counter.blocks = append(counter.blocks, &BlockEntry{
		starPos: cursor,
		cost:    1,
		index:   []int{},
	})
	counter.stack.Push(uint64(blockIdx))
}

func (counter *Counter) Finalize() {
	counter.stack.Pop()
}

func (counter *Counter) Increment(value uint64, index int) {
	// var top uint64
	// if counter.stack.Len() == 0 {
	// 	top = 0
	// } else {
	// 	top = counter.stack.Top()
	// }
	top := counter.stack.Top()

	topBlock := counter.blocks[top]
	topBlock.cost = topBlock.cost + value
	topBlock.index = append(topBlock.index, index)
}

func InjectCounter(disassembly []disasm.Instr, module *wasm.Module, rule Gas) []disasm.Instr {
	_, _, gasIndex := utils.GetIndex(module)
	if gasIndex == -1 {
		return disassembly
	}
	counter := NewCounter()
	counter.Begin(0)
	for i, instr := range disassembly {
		switch instr.Op.Code {
		case ops.Block, ops.Loop, ops.If:
			//instruction_cost = rules.process(instruction)?;
			// instrCost := 1 //instr的gas消耗规则
			// counter.Increment(uint32(instrCost), i)

			// Begin new block. The cost of the following opcodes until `End` or `Else` will
			// be included into this block.
			counter.Begin(i + 1)
		case ops.Br, ops.BrIf, ops.BrTable:
			counter.Finalize()
			// instrCost := 1 //instr的gas消耗规则
			// counter.Increment(uint32(instrCost), i)
			counter.Begin(i + 1)
		case ops.End:
			counter.Finalize()
		case ops.Else:
			counter.Finalize()
			counter.Begin(i + 1)
		default:
			instrCost := rule.GasCost(instr.Op.Code)
			counter.Increment(instrCost, i)
		}
	}

	for _, v := range counter.blocks {
		if len(v.index) > 0 {
			v.starPos = v.index[0]
		}
	}
	sort.Sort(counter.blocks)
	offset := 0
	for _, v := range counter.blocks {
		pos := v.starPos + offset
		constOp, _ := ops.New(ops.I64Const)
		constInstr := disasm.Instr{Op: constOp, Immediates: []interface{}{int64(v.cost)}}
		callOp, _ := ops.New(ops.Call)
		callInstr := disasm.Instr{Op: callOp, Immediates: []interface{}{uint32(gasIndex)}}
		//disassembly=append(disassembly[0:pos],)
		res := common.Insert(disassembly, pos, []disasm.Instr{callInstr})
		disassembly = res.([]disasm.Instr)
		res = common.Insert(disassembly, pos, []disasm.Instr{constInstr})
		disassembly = res.([]disasm.Instr)
		offset += 2
	}
	return disassembly
}
