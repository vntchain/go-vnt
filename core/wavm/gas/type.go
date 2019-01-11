package gas

type Metering int

const (
	MeteringRegular Metering = iota
	MeteringForbidden
	MeteringFixed
)

type InstructionType int

const (
	InstructionTypeBit InstructionType = iota
	InstructionTypeAdd
	InstructionTypeMul
	InstructionTypeDiv
	InstructionTypeLoad
	InstructionTypeStore
	InstructionTypeConst
	InstructionTypeFloatConst
	InstructionTypeLocal
	InstructionTypeGlobal
	InstructionTypeControlFlow
	InstructionTypeIntegerComparsion
	InstructionTypeFloatComparsion
	InstructionTypeFloat
	InstructionTypeConversion
	InstructionTypeFloatConversion
	InstructionTypeReinterpretation
	InstructionTypeUnreachable
	InstructionTypeNop
	InstructionTypeCurrentMemory
	InstructionTypeGrowMemory
)
