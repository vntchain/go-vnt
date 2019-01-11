package gen

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/go-clang/bootstrap/clang"
)

type HeaderFile struct {
	Lookup

	api *API

	Filename string
	Path     string

	Enums     []*Enum
	Functions []*Function
	Structs   []*Struct
}

func newHeaderFile(a *API, name string, dir string) *HeaderFile {
	return &HeaderFile{
		api: a,

		Path:     dir,
		Filename: name,

		Lookup: NewLookup(),
	}
}

func (h *HeaderFile) prepareFile() error {
	/*
		Hide all "void *" fields of structs by replacing the type with "uintptr_t".

		To paraphrase the original go-clang source code:
			Not hiding these fields confuses the Go GC during garbage collection and
			pointer scanning, making it think the heap/stack has been somehow corrupted.

		I do not know how the original author debugged this, but one thing: Thank you!
	*/
	findStructsRe := regexp.MustCompile(`(?s)struct[\s\w]+{.+?}`)
	f, err := ioutil.ReadFile(h.FullPath())
	if err != nil {
		return fmt.Errorf("Cannot read Index.h: %v", err)
	}
	voidPointerReplacements := map[string]string{}
	findVoidPointerRe := regexp.MustCompile(`(?:const\s+)?void\s*\*\s*(\w+(\[\d+\])?;)`)
	for _, s := range findStructsRe.FindAll(f, -1) {
		s2 := findVoidPointerRe.ReplaceAll(s, []byte("uintptr_t $1"))
		if len(s) != len(s2) {
			voidPointerReplacements[string(s)] = string(s2)
		}
	}
	fs := string(f)
	for s, r := range voidPointerReplacements {
		fs = strings.Replace(fs, s, r, -1)
	}
	if incl := "#include <stdint.h>"; !strings.HasPrefix(fs, incl) { // Include for uintptr_t
		fs = "#include <stdint.h>\n\n" + fs
	}
	err = ioutil.WriteFile(h.FullPath(), []byte(fs), 0600)
	if err != nil {
		return fmt.Errorf("Cannot write %s: %v", h.FullPath(), err)
	}

	return nil
}

func (h *HeaderFile) handleFile(cursor clang.Cursor) {
	/*
		TODO mark the enum https://github.com/go-clang/gen/issues/40
			typedef enum CXChildVisitResult (*CXCursorVisitor)(CXCursor cursor, CXCursor parent, CXClientData client_data);
		as manually implemented
	*/
	// TODO report other enums like callbacks that they are not implemented https://github.com/go-clang/gen/issues/51

	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		// Only handle code of the current file
		sourceFile, _, _, _ := cursor.Location().FileLocation()
		isCurrentFile := sourceFile.Name() == h.FullPath()

		if !strings.HasPrefix(sourceFile.Name(), h.Path) {
			return clang.ChildVisit_Continue
		}

		cname := cursor.Spelling()
		cnameIsTypeDef := false

		if parentCName := parent.Spelling(); parent.Kind() == clang.Cursor_TypedefDecl && parentCName != "" {
			cname = parentCName
			cnameIsTypeDef = true
		}

		switch cursor.Kind() {
		case clang.Cursor_EnumDecl:
			if cname == "" {
				break
			}

			e := handleEnumCursor(cursor, cname, cnameIsTypeDef)
			e.IncludeFiles.addIncludeFile(sourceFile.Name())

			if _, ok := h.HasEnum(e.Name); !ok {
				h.Enums = append(h.Enums, e)
			}

		case clang.Cursor_FunctionDecl:
			if !isCurrentFile {
				return clang.ChildVisit_Continue
			}

			f := handleFunctionCursor(cursor)
			if f != nil {
				f.IncludeFiles.addIncludeFile(sourceFile.Name())
				h.Functions = append(h.Functions, f)
			}
		case clang.Cursor_StructDecl:
			if cname == "" {
				break
			}

			s := handleStructCursor(cursor, cname, cnameIsTypeDef)
			s.api = h.api
			s.IncludeFiles.addIncludeFile(sourceFile.Name())

			if _, ok := h.HasStruct(s.Name); !ok {
				h.RegisterStruct(s)
				h.Structs = append(h.Structs, s)
			}
		case clang.Cursor_TypedefDecl:
			underlyingType := cursor.TypedefDeclUnderlyingType().Spelling()
			underlyingStructType := strings.TrimSuffix(strings.TrimPrefix(underlyingType, "struct "), " *")

			if s, ok := h.HasStruct(underlyingStructType); ok && !s.CNameIsTypeDef && strings.HasPrefix(underlyingType, "struct "+s.CName) {
				// Sometimes the typedef is not a parent of the struct but a sibling
				sn := handleStructCursor(cursor, cname, true)
				sn.api = h.api
				sn.IncludeFiles.addIncludeFile(sourceFile.Name())

				if sn.Comment == "" {
					sn.Comment = s.Comment
				}
				sn.Members = s.Members
				sn.Methods = s.Methods

				h.RemoveStruct(s)
				h.RegisterStruct(sn)

				for i, si := range h.Structs {
					if si == s {
						h.Structs[i] = sn

						break
					}
				}
			} else if underlyingType == "void *" {
				s := handleStructCursor(cursor, cname, true)
				s.api = h.api
				s.IncludeFiles.addIncludeFile(sourceFile.Name())

				if _, ok := h.HasStruct(s.Name); !ok {
					h.RegisterStruct(s)
					h.Structs = append(h.Structs, s)
				}
			}
		}

		return clang.ChildVisit_Recurse
	})
}

func (h *HeaderFile) parse(clangArguments []string) error {
	if err := h.prepareFile(); err != nil {
		return err
	}

	// Parse the header file to analyse everything we need to know
	idx := clang.NewIndex(0, 1)
	defer idx.Dispose()

	tu := idx.ParseTranslationUnit(h.FullPath(), clangArguments, nil, 0)
	defer tu.Dispose()

	if !tu.IsValid() {
		return fmt.Errorf("Cannot parse Index.h")
	}

	for _, diag := range tu.Diagnostics() {
		switch diag.Severity() {
		case clang.Diagnostic_Error:
			return fmt.Errorf("Diagnostic error in Index.h: %s", diag.Spelling())
		case clang.Diagnostic_Fatal:
			return fmt.Errorf("Diagnostic fatal in Index.h: %s", diag.Spelling())
		}
	}

	h.handleFile(tu.TranslationUnitCursor())

	return nil
}

func (h *HeaderFile) FullPath() string {
	path := h.Path
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	return path + h.Filename
}
