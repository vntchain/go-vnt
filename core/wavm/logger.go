package wavm

import (
	"io"
	"math/big"
	"time"

	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/common/math"
	"github.com/vntchain/go-vnt/core/vm"
	"github.com/vntchain/go-vnt/core/vm/interface"
)

type WasmLogger struct {
	cfg       vm.LogConfig
	logs      []StructLog
	debugLogs []DebugLog
}

type StructLog struct {
	Pc         uint64                      `json:"pc"`
	Op         vm.OPCode                   `json:"op"`
	Gas        uint64                      `json:"gas"`
	GasCost    uint64                      `json:"gasCost"`
	Memory     []byte                      `json:"-"`
	MemorySize int                         `json:"-"`
	Stack      []*big.Int                  `json:"-"`
	Storage    map[common.Hash]common.Hash `json:"-"`
	Depth      int                         `json:"depth"`
	Err        error                       `json:"error"`
}

type DebugLog struct {
	PrintMsg string `json:"printMsg"`
}

// overrides for gencodec
type structLogMarshaling struct {
	Gas         math.HexOrDecimal64
	GasCost     math.HexOrDecimal64
	OpName      string `json:"opName"` // adds call to OpName() in MarshalJSON
	ErrorString string `json:"error"`  // adds call to ErrorString() in MarshalJSON
}

func (s *StructLog) OpName() string {
	return s.Op.String()
}

func (s *StructLog) ErrorString() string {
	if s.Err != nil {
		return s.Err.Error()
	}
	return ""
}

func NewWasmLogger(cfg *vm.LogConfig) *WasmLogger {
	logger := &WasmLogger{}
	if cfg != nil {
		logger.cfg = *cfg
	}
	return logger
}

func (l *WasmLogger) CaptureStart(from common.Address, to common.Address, call bool, input []byte, gas uint64, value *big.Int) error {
	return nil
}
func (l *WasmLogger) CaptureState(env vm.VM, pc uint64, op vm.OPCode, gas, cost uint64, memory *vm.Memory, stack *vm.Stack, contract inter.Contract, depth int, err error) error {
	// check if already accumulated the specified number of logs
	if l.cfg.Limit != 0 && l.cfg.Limit <= len(l.logs) {
		return vm.ErrTraceLimitReached
	}

	// create a new snaptshot of the VM.
	log := StructLog{pc, op, gas, cost, nil, 0, nil, nil, depth, err}

	l.logs = append(l.logs, log)
	return nil
}
func (l *WasmLogger) CaptureLog(env vm.VM, msg string) error {
	log := DebugLog{msg}
	l.debugLogs = append(l.debugLogs, log)
	return nil
}
func (l *WasmLogger) CaptureFault(env vm.VM, pc uint64, op vm.OPCode, gas, cost uint64, memory *vm.Memory, stack *vm.Stack, contract inter.Contract, depth int, err error) error {
	return nil
}
func (l *WasmLogger) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) error {
	return nil
}

// WasmLogger returns the captured log entries.
func (l *WasmLogger) StructLogs() []StructLog { return l.logs }

// DebugLogs returns the captured debug log entries.
func (l *WasmLogger) DebugLogs() []DebugLog { return l.debugLogs }

// WriteTrace writes a formatted trace to the given writer
func WriteTrace(writer io.Writer, logs []StructLog) {

}
