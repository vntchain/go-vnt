package wavm

import "github.com/vntchain/go-vnt/core/vm"

// Config are the configuration options for the Interpreter
type Config struct {
	// Debug enabled debugging Interpreter options
	Debug bool
	// Tracer is the op code logger
	Tracer                   vm.Tracer
	NoRecursion              bool
	MaxMemoryPages           int
	MaxTableSize             int
	MaxValueSlots            int
	MaxCallStackDepth        int
	DefaultMemoryPages       int
	DefaultTableSize         int
	GasLimit                 uint64
	DisableFloatingPoint     bool
	ReturnOnGasLimitExceeded bool
}
