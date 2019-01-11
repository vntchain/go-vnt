package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-clang/gen"
	genclang "github.com/go-clang/gen/clang"
)

func main() {
	api := &gen.API{
		PrepareFunctionName:     prepareFunctionName,
		PrepareFunction:         prepareFunction,
		FilterFunction:          filterFunction,
		FilterFunctionParameter: filterFunctionParameter,
		FixedFunctionName:       fixedFunctionName,

		PrepareStructMembers:     prepareStructMembers,
		FilterStructMemberGetter: filterStructMemberGetter,
	}

	err := genclang.Cmd(os.Args[1:], api)
	if err != nil {
		fmt.Println(err)

		os.Exit(1)
	}
}

func prepareFunctionName(g *gen.Generation, f *gen.Function) string {
	fname := f.Name

	fname = strings.TrimPrefix(fname, "clang_")

	// Trim some whitelisted prefixes by their function name
	if fn := strings.TrimPrefix(fname, "indexLoc_"); len(fn) != len(fname) {
		fname = fn
	} else if fn := strings.TrimPrefix(fname, "index_"); len(fn) != len(fname) {
		fname = fn
	} else if fn := strings.TrimPrefix(fname, "Location_"); len(fn) != len(fname) {
		fname = fn
	} else if fn := strings.TrimPrefix(fname, "Range_"); len(fn) != len(fname) {
		fname = fn
	} else if fn := strings.TrimPrefix(fname, "remap_"); len(fn) != len(fname) {
		fname = fn
	}

	// Trim some whitelisted prefixes by their types
	if len(f.Parameters) > 0 && g.IsEnumOrStruct(f.Parameters[0].Type.GoName) {
		switch f.Parameters[0].Type.GoName {
		case "CodeCompleteResults":
			fname = strings.TrimPrefix(fname, "codeComplete")
		case "CompletionString":
			if f.CName == "clang_getNumCompletionChunks" {
				fname = "NumChunks"
			} else {
				fname = strings.TrimPrefix(fname, "getCompletion")
			}
		case "SourceRange":
			fname = strings.TrimPrefix(fname, "getRange")
		}
	}

	return fname
}

func fixedFunctionName(f *gen.Function) string {
	// Needs to be renamed manually since clang_getTranslationUnitCursor will conflict with clang_getCursor
	if f.CName == "clang_getTranslationUnitCursor" {
		return "TranslationUnitCursor"
	}

	return ""
}

func prepareFunction(f *gen.Function) {
	for i := range f.Parameters {
		p := &f.Parameters[i]

		if f.CName == "clang_getRemappingsFromFileList" {
			switch p.CName {
			case "filePaths":
				p.Type.IsSlice = true
			case "numFiles":
				p.Type.LengthOfSlice = "filePaths"
			}

			continue
		}

		// Whiteflag types that are return arguments
		if p.Type.PointerLevel == 1 && (p.Type.GoName == "File" || p.Type.GoName == "FileUniqueID" || p.Type.GoName == "IdxClientFile" || p.Type.GoName == "cxstring" || p.Type.GoName == gen.GoInt32 || p.Type.GoName == gen.GoUInt32 || p.Type.GoName == "CompilationDatabase_Error" || p.Type.GoName == "PlatformAvailability" || p.Type.GoName == "SourceRange" || p.Type.GoName == "LoadDiag_Error") {
			p.Type.IsReturnArgument = true
		}
		if p.Type.PointerLevel == 2 && (p.Type.GoName == "Token" || p.Type.GoName == "Cursor") {
			p.Type.IsReturnArgument = true
		}

		if f.CName == "clang_disposeOverriddenCursors" && p.CName == "overridden" {
			p.Type.IsSlice = true
		}

		// If this is an array length parameter we need to find its partner
		paCName := gen.ArrayNameFromLength(p.CName)

		if paCName != "" {
			for j := range f.Parameters {
				pa := &f.Parameters[j]

				if strings.ToLower(pa.CName) == strings.ToLower(paCName) {
					if pa.Type.GoName == "struct CXUnsavedFile" || pa.Type.GoName == "UnsavedFile" {
						pa.Type.GoName = "UnsavedFile"
						pa.Type.CGoName = "struct_CXUnsavedFile"
					} else if pa.Type.CGoName == gen.CSChar && pa.Type.PointerLevel == 2 {
					} else if pa.Type.GoName == "CompletionResult" {
					} else if pa.Type.GoName == "Token" {
					} else if pa.Type.GoName == "Cursor" {
					} else {
						break
					}

					p.Type.LengthOfSlice = pa.Name
					pa.Type.IsSlice = true

					if pa.Type.IsReturnArgument && p.Type.PointerLevel > 0 {
						p.Type.IsReturnArgument = true
					}

					break
				}
			}
		}
	}

	for i := range f.Parameters {
		p := &f.Parameters[i]

		if p.Type.CGoName == gen.CSChar && p.Type.PointerLevel == 2 && !p.Type.IsSlice {
			p.Type.IsReturnArgument = true
		}
	}
}

func filterFunction(f *gen.Function) bool {
	// Some functions are not compiled in the library see https://lists.launchpad.net/desktop-packages/msg75835.html for a never resolved bug report
	if f.CName == "clang_CompileCommand_getMappedSourceContent" || f.CName == "clang_CompileCommand_getMappedSourcePath" || f.CName == "clang_CompileCommand_getNumMappedSources" {
		fmt.Printf("Ignore function %q because it is not compiled within libClang\n", f.CName)

		return false
	}
	// Some functions can not be handled automatically by us
	if f.CName == "clang_executeOnThread" || f.CName == "clang_getInclusions" {
		fmt.Printf("Ignore function %q because it cannot be handled automatically\n", f.CName)

		return false
	}
	// Some functions are simply manually implemented
	if f.CName == "clang_annotateTokens" || f.CName == "clang_getCursorPlatformAvailability" || f.CName == "clang_visitChildren" {
		fmt.Printf("Ignore function %q because it is manually implemented\n", f.CName)

		return false
	}

	// TODO if this function is from CXString.h we ignore it https://github.com/go-clang/gen/issues/25
	for i := range f.IncludeFiles {
		if strings.HasSuffix(i, "CXString.h") {
			return false
		}
	}

	return true
}

func filterFunctionParameter(p gen.FunctionParameter) bool {
	// These pointers are ok
	if p.Type.PointerLevel == 1 && (p.Type.CGoName == gen.CSChar || p.Type.GoName == "UnsavedFile" || p.Type.GoName == "CodeCompleteResults" || p.Type.GoName == "CursorKind" || p.Type.GoName == "IdxContainerInfo" || p.Type.GoName == "IdxDeclInfo" || p.Type.GoName == "IndexerCallbacks" || p.Type.GoName == "TranslationUnit" || p.Type.GoName == "IdxEntityInfo" || p.Type.GoName == "IdxAttrInfo") {
		return false
	}

	return true
}

func prepareStructMembers(s *gen.Struct) {
	for _, m := range s.Members {
		if (strings.HasPrefix(m.CName, "has") || strings.HasPrefix(m.CName, "is")) && m.Type.GoName == gen.GoInt32 {
			m.Type.GoName = gen.GoBool
		}

		// If this is an array length parameter we need to find its partner
		maCName := gen.ArrayNameFromLength(m.CName)

		if maCName != "" {
			for _, ma := range s.Members {
				if strings.ToLower(ma.CName) == strings.ToLower(maCName) {
					m.Type.LengthOfSlice = ma.CName
					ma.Type.IsSlice = true
					ma.Type.LengthOfSlice = m.CName // TODO wrong usage but needed for the getter generation... maybe refactor this LengthOfSlice alltogether? https://github.com/go-clang/gen/issues/49

					break
				}
			}
		}
	}

	prepareStructMembersArrayStruct(s)
}

// prepareStructMembersArrayStruct checks if the struct has two member variables, one is an array and the other a plain int/uint with size/length/count/len as its name because then this should be an array struct, and we connect them to handle a slice.
func prepareStructMembersArrayStruct(s *gen.Struct) {
	if len(s.Members) != 2 {
		return
	}

	if !arrayLengthCombination(&s.Members[0].Type, &s.Members[1].Type) && !arrayLengthCombination(&s.Members[1].Type, &s.Members[0].Type) {
		return
	}

	// If one of the members is already marked as array/slice another heuristic has already covered both members.
	if s.Members[0].Type.IsArray || s.Members[1].Type.IsArray ||
		s.Members[0].Type.IsSlice || s.Members[1].Type.IsSlice {
		return
	}

	var a *gen.StructMember
	var c *gen.StructMember

	if s.Members[0].Type.PointerLevel == 1 {
		a = s.Members[0]
		c = s.Members[1]
	} else {
		c = s.Members[0]
		a = s.Members[1]
	}

	lengthName := strings.ToLower(c.CName)
	if lengthName != "count" && lengthName != "len" && lengthName != "length" && lengthName != "size" {
		return
	}

	c.Type.LengthOfSlice = a.CName
	a.Type.IsSlice = true
	a.Type.LengthOfSlice = c.CName // TODO wrong usage but needed for the getter generation... maybe refactor this LengthOfSlice alltogether? https://github.com/go-clang/gen/issues/49
}

func arrayLengthCombination(a *gen.Type, b *gen.Type) bool {
	return a.PointerLevel == 1 && b.PointerLevel == 0 &&
		!gen.IsInteger(a) && gen.IsInteger(b)
}

func filterStructMemberGetter(m *gen.StructMember) bool {
	// We do not want getters to *int_data members
	if strings.HasSuffix(m.CName, "int_data") {
		return false
	}

	return true
}
