package gen

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/go-clang/bootstrap/clang"
)

// Defines all available Go types
const (
	GoByte      = "byte"
	GoInt8      = "int8"
	GoUInt8     = "uint8"
	GoInt16     = "int16"
	GoUInt16    = "uint16"
	GoInt32     = "int32"
	GoUInt32    = "uint32"
	GoInt64     = "int64"
	GoUInt64    = "uint64"
	GoFloat32   = "float32"
	GoFloat64   = "float64"
	GoBool      = "bool"
	GoInterface = "interface"
	GoPointer   = "unsafe.Pointer"
)

// Defines all available C types
const (
	CChar      = "char"
	CSChar     = "schar"
	CUChar     = "uchar"
	CShort     = "short"
	CUShort    = "ushort"
	CInt       = "int"
	CUInt      = "uint"
	CLongInt   = "long"
	CULongInt  = "ulong"
	CLongLong  = "longlong"
	CULongLong = "ulonglong"
	CFloat     = "float"
	CDouble    = "double"
)

// Type represents a generation type
type Type struct {
	CName   string
	CGoName string
	GoName  string

	PointerLevel      int
	IsPrimitive       bool
	IsArray           bool
	ArraySize         int64
	IsEnumLiteral     bool
	IsFunctionPointer bool
	IsReturnArgument  bool
	IsSlice           bool
	LengthOfSlice     string

	IsPointerComposition bool
}

func typeFromClangType(cType clang.Type) (Type, error) {
	typ := Type{
		CName: cType.Spelling(),

		PointerLevel:      0,
		IsPrimitive:       true,
		IsArray:           false,
		ArraySize:         -1,
		IsEnumLiteral:     false,
		IsFunctionPointer: false,
	}

	switch cType.Kind() {
	case clang.Type_Char_S:
		typ.CGoName = CSChar
		typ.GoName = GoInt8
	case clang.Type_Char_U:
		typ.CGoName = CUChar
		typ.GoName = GoUInt8
	case clang.Type_Int:
		typ.CGoName = CInt
		typ.GoName = GoInt32
	case clang.Type_Short:
		typ.CGoName = CShort
		typ.GoName = GoInt16
	case clang.Type_UShort:
		typ.CGoName = CUShort
		typ.GoName = GoUInt16
	case clang.Type_UInt:
		typ.CGoName = CUInt
		typ.GoName = GoUInt32
	case clang.Type_Long:
		typ.CGoName = CLongInt
		typ.GoName = GoInt64
	case clang.Type_ULong:
		typ.CGoName = CULongInt
		typ.GoName = GoUInt64
	case clang.Type_LongLong:
		typ.CGoName = CLongLong
		typ.GoName = GoInt64
	case clang.Type_ULongLong:
		typ.CGoName = CULongLong
		typ.GoName = GoUInt64
	case clang.Type_Float:
		typ.CGoName = CFloat
		typ.GoName = GoFloat32
	case clang.Type_Double:
		typ.CGoName = CDouble
		typ.GoName = GoFloat64
	case clang.Type_Bool:
		typ.GoName = GoBool
	case clang.Type_Void:
		// TODO Does not exist in Go, what should we do with it? https://github.com/go-clang/gen/issues/50
		typ.CGoName = "void"
		typ.GoName = "void"
	case clang.Type_ConstantArray:
		subTyp, err := typeFromClangType(cType.ArrayElementType())
		if err != nil {
			return Type{}, err
		}

		typ.CGoName = subTyp.CGoName
		typ.GoName = subTyp.GoName
		typ.PointerLevel += subTyp.PointerLevel
		typ.IsArray = true
		typ.ArraySize = cType.ArraySize()
	case clang.Type_Typedef:
		typ.IsPrimitive = false

		typeStr := cType.Spelling()
		if typeStr == "CXString" { // TODO eliminate CXString from the generic code https://github.com/go-clang/gen/issues/25
			typeStr = "cxstring"
		} else if typeStr == "time_t" {
			typ.CGoName = typeStr
			typeStr = "time.Time"

			typ.IsPrimitive = true
		} else {
			typeStr = TrimLanguagePrefix(cType.Declaration().Type().Spelling())
		}

		typ.CGoName = cType.Declaration().Type().Spelling()
		typ.GoName = typeStr

		if cType.CanonicalType().Kind() == clang.Type_Enum {
			typ.IsEnumLiteral = true
			typ.IsPrimitive = true
		}
	case clang.Type_Pointer:
		typ.PointerLevel++

		if cType.PointeeType().CanonicalType().Kind() == clang.Type_FunctionProto {
			typ.IsFunctionPointer = true
		}

		subTyp, err := typeFromClangType(cType.PointeeType())
		if err != nil {
			return Type{}, err
		}

		typ.CGoName = subTyp.CGoName
		typ.GoName = subTyp.GoName
		typ.PointerLevel += subTyp.PointerLevel
		typ.IsPrimitive = subTyp.IsPrimitive
	case clang.Type_Record:
		typ.CGoName = cType.Declaration().Type().Spelling()
		typ.GoName = TrimLanguagePrefix(typ.CGoName)
		typ.IsPrimitive = false
	case clang.Type_FunctionProto:
		typ.IsFunctionPointer = true
		typ.CGoName = cType.Declaration().Type().Spelling()
		typ.GoName = TrimLanguagePrefix(typ.CGoName)
	case clang.Type_Enum:
		typ.GoName = TrimLanguagePrefix(cType.Declaration().DisplayName())
		typ.IsEnumLiteral = true
		typ.IsPrimitive = true
	case clang.Type_Unexposed: // There is a bug in clang for enums the kind is set to unexposed dunno why, bug persists since 2013 https://llvm.org/bugs/show_bug.cgi?id=15089
		subTyp, err := typeFromClangType(cType.CanonicalType())
		if err != nil {
			return Type{}, err
		}

		typ.CGoName = subTyp.CGoName
		typ.GoName = subTyp.GoName
		typ.PointerLevel += subTyp.PointerLevel
		typ.IsPrimitive = subTyp.IsPrimitive
	default:
		return Type{}, fmt.Errorf("unhandled type %q of kind %q", cType.Spelling(), cType.Kind().Spelling())
	}

	return typ, nil
}

func ArrayNameFromLength(lengthCName string) string {
	if pan := strings.TrimPrefix(lengthCName, "num_"); len(pan) != len(lengthCName) {
		return pan
	} else if pan := strings.TrimPrefix(lengthCName, "num"); len(pan) != len(lengthCName) {
		return pan
	} else if pan := strings.TrimPrefix(lengthCName, "Num"); len(pan) != len(lengthCName) && unicode.IsUpper(rune(pan[0])) {
		return pan
	} else if pan := strings.TrimSuffix(lengthCName, "_size"); len(pan) != len(lengthCName) {
		return pan
	}

	return ""
}

func IsInteger(typ *Type) bool {
	switch typ.GoName {
	case GoInt8, GoUInt8, GoInt16, GoUInt16, GoInt32, GoUInt32, GoInt64, GoUInt64:
		return true
	}

	return false
}
