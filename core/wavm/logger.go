package wavm

import (
	"fmt"
	"math/big"
	"time"

	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core/vm"
	"github.com/vntchain/go-vnt/core/vm/interface"
)

type WasmLogger struct {
}

func NewWasmLogger() *WasmLogger {
	return &WasmLogger{}
}

func (l *WasmLogger) CaptureStart(from common.Address, to common.Address, call bool, input []byte, gas uint64, value *big.Int) error {
	return nil
}
func (l *WasmLogger) CaptureState(env vm.VM, pc uint64, op vm.OPCode, gas, cost uint64, memory *vm.Memory, stack *vm.Stack, contract inter.Contract, depth int, err error) error {
	fmt.Printf("CaptureState \n")
	return nil
}
func (l *WasmLogger) CaptureLog(env vm.VM, msg string) error {
	return nil
}
func (l *WasmLogger) CaptureFault(env vm.VM, pc uint64, op vm.OPCode, gas, cost uint64, memory *vm.Memory, stack *vm.Stack, contract inter.Contract, depth int, err error) error {
	return nil
}
func (l *WasmLogger) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) error {
	return nil
}
