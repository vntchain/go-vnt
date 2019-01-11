package gen

import (
	"go/ast"
	"strings"

	"github.com/go-clang/bootstrap/clang"
)

// Enum represents a generation enum
type Enum struct {
	IncludeFiles includeFiles

	Name           string
	CName          string
	CNameIsTypeDef bool
	Receiver       Receiver
	Comment        string
	UnderlyingType string

	Items []EnumItem

	Methods []interface{}
}

// EnumItem represents a generation enum item
type EnumItem struct {
	Name    string
	CName   string
	Comment string
	Value   uint64
}

func handleEnumCursor(cursor clang.Cursor, cname string, cnameIsTypeDef bool) *Enum {
	e := Enum{
		CName:          cname,
		CNameIsTypeDef: cnameIsTypeDef,
		Comment:        CleanDoxygenComment(cursor.RawCommentText()),

		IncludeFiles: newIncludeFiles(),

		Items: []EnumItem{},
	}

	e.Name = TrimLanguagePrefix(e.CName)

	e.Receiver.Name = commonReceiverName(e.Name)
	e.Receiver.Type.GoName = e.Name
	e.Receiver.Type.CGoName = e.CName
	if cnameIsTypeDef {
		e.Receiver.Type.CGoName = e.CName
	} else {
		e.Receiver.Type.CGoName = "enum_" + e.CName
	}

	enumNamePrefix := e.Name
	enumNamePrefix = strings.TrimSuffix(enumNamePrefix, "Kind")
	enumNamePrefix = strings.SplitN(enumNamePrefix, "_", 2)[0]

	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		switch cursor.Kind() {
		case clang.Cursor_EnumConstantDecl:
			ei := EnumItem{
				CName:   cursor.Spelling(),
				Comment: CleanDoxygenComment(cursor.RawCommentText()), // TODO We are always using the same comment if there is none, see "TypeKind" https://github.com/go-clang/gen/issues/58
				Value:   cursor.EnumConstantDeclUnsignedValue(),
			}
			ei.Name = TrimLanguagePrefix(ei.CName)

			// Check if the first item has an enum prefix
			if len(e.Items) == 0 {
				eis := strings.SplitN(ei.Name, "_", 2)
				if len(eis) == 2 {
					enumNamePrefix = ""
				}
			}

			// Add the enum prefix to the item
			if enumNamePrefix != "" {
				ei.Name = strings.TrimSuffix(ei.Name, enumNamePrefix)

				if !strings.HasPrefix(ei.Name, enumNamePrefix) {
					ei.Name = enumNamePrefix + "_" + ei.Name
				}
			}

			e.Items = append(e.Items, ei)
		default:
			panic(cursor.Kind())
		}

		return clang.ChildVisit_Continue
	})

	if strings.HasSuffix(e.Name, "Error") {
		e.UnderlyingType = "int32"
	} else {
		e.UnderlyingType = "uint32"
	}

	return &e
}

func (e *Enum) ContainsMethod(name string) bool {
	for _, m := range e.Methods {
		switch m := m.(type) {
		case *Function:
			if m.Name == name {
				return true
			}
		case string:
			if strings.Contains(m, ") "+name+"()") {
				return true
			}
		}
	}

	return false
}

func (e *Enum) generate() error {
	f := newFile(strings.ToLower(e.Name))
	f.Enums = append(f.Enums, e)

	return f.generate()
}

func (e *Enum) addEnumStringMethods() error {
	if !e.ContainsMethod("Spelling") {
		if err := e.addEnumSpellingMethod(); err != nil {
			return err
		}
	}
	if !e.ContainsMethod("String") {
		if err := e.addSpellingMethodAlias("String"); err != nil {
			return err
		}
	}
	if strings.HasSuffix(e.Name, "Error") && !e.ContainsMethod("Error") {
		if err := e.addSpellingMethodAlias("Error"); err != nil {
			return err
		}
	}

	return nil
}

func (e *Enum) addEnumSpellingMethod() error {
	f := newFunction("Spelling", e.Name, "", "", Type{GoName: "string"})
	fa := newASTFunc(f)

	fa.generateReceiver()

	fa.addReturnType("", Type{
		GoName: "string",
	})

	switchStmt := doSwitchStmt(&ast.Ident{Name: f.Receiver.Name})
	fa.Body.List = append(fa.Body.List, switchStmt)

	m := make(map[uint64]*ast.CaseClause)

	for _, enumerator := range e.Items {
		/* 	EnumItems might have the same value:
		 	e.g.,
		 	enum Example {
				aValue = 1
				bValue = aValue
			}

			Translating each EnumItem in its own case would result in a compliation error
			(https://code.google.com/p/go/issues/detail?id=4524).
			Thus, all EnumItems with the same value need to be pooled.
		*/
		if caseClause, ok := m[enumerator.Value]; !ok {
			c := []ast.Expr{&ast.Ident{Name: enumerator.Name}}

			ret := &ast.ReturnStmt{
				Results: []ast.Expr{
					doStringLit(strings.Replace(enumerator.Name, "_", "=", 1)),
				},
			}

			b := []ast.Stmt{ret}

			caseClause = doCaseClause(c, b)
			switchStmt.Body.List = append(switchStmt.Body.List, caseClause)
			m[enumerator.Value] = caseClause
		} else {
			retStr := caseClause.Body[0].(*ast.ReturnStmt).Results[0].(*ast.BasicLit).Value
			retStr = retStr[0:len(retStr)-1] + ", " + enumerator.Name[strings.Index(enumerator.Name, "_")+1:] + "\""
			caseClause.Body[0].(*ast.ReturnStmt).Results[0].(*ast.BasicLit).Value = retStr
		}
	}

	fa.addReturnItem(doCall(
		"fmt",
		"Sprintf",
		doStringLit(f.Receiver.Type.GoName+" unkown %d"),
		doCast("int", &ast.Ident{Name: f.Receiver.Name}),
	))

	fa.addEmptyLine()
	fa.addStatement(fa.ret)

	e.Methods = append(e.Methods, generateFunctionString(fa))

	return nil
}

func (e *Enum) addSpellingMethodAlias(name string) error {
	returnType := Type{
		GoName: "string",
	}

	f := newFunction(name, e.Name, "", "", returnType)
	fa := newASTFunc(f)

	fa.generateReceiver()

	fa.addReturnType("", Type{
		GoName: "string",
	})

	fa.addReturnItem(doCall(
		e.Receiver.Name,
		"Spelling",
	))

	fa.addStatement(fa.ret)

	e.Methods = append(e.Methods, generateFunctionString(fa))

	return nil
}
