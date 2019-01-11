package common

type VmType byte

var (
	VmTypeNone = VmType(0x00)
	VmTypeEvm  = VmType(0x01)
	VmTypeWavm = VmType(0x02)
)
