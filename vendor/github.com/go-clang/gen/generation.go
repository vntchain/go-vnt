package gen

import (
	"fmt"
	"strings"
	"unicode"
)

type Generation struct {
	api *API

	enums     []*Enum
	functions []*Function
	structs   []*Struct

	Lookup
}

func NewGeneration(a *API) *Generation {
	gen := &Generation{
		api: a,

		Lookup: NewLookup(),
	}

	return gen
}

func (g *Generation) AddHeaderFiles(headerFiles []*HeaderFile) {
	for _, h := range headerFiles {
		for _, e := range h.Enums {
			g.enums = append(g.enums, e)
			g.RegisterEnum(e)
		}

		for _, s := range h.Structs {
			g.structs = append(g.structs, s)
			g.RegisterStruct(s)
		}

		for _, f := range h.Functions {
			g.functions = append(g.functions, f)
		}
	}
}

func (g *Generation) Generate() error {
	// Prepare all functions
	clangFile := newFile("clang")

	for _, f := range g.functions {
		fname := f.Name
		if g.api.PrepareFunctionName != nil {
			fname = g.api.PrepareFunctionName(g, f)
			f.Name = fname
		}

		if g.api.FilterFunction != nil && !g.api.FilterFunction(f) {
			continue
		}

		// Prepare the parameters
		for i := range f.Parameters {
			p := &f.Parameters[i]

			if n, ok := g.LookupNonTypedef(p.Type.CGoName); ok {
				p.Type.GoName = n
			}
			if e, ok := g.HasEnum(p.Type.GoName); ok {
				p.CName = e.Receiver.CName
				// TODO remove the receiver... and copy only names here to preserve the original pointers and so https://github.com/go-clang/gen/issues/52
				p.Type.GoName = e.Receiver.Type.GoName
				p.Type.CGoName = e.Receiver.Type.CGoName
				p.Type.CGoName = e.Receiver.Type.CGoName
			}
		}

		if g.api.PrepareFunction != nil {
			g.api.PrepareFunction(f)
		}

		// Prepare the return argument
		if n, ok := g.LookupNonTypedef(f.ReturnType.CGoName); ok {
			f.ReturnType.GoName = n
		}
		if e, ok := g.HasEnum(f.ReturnType.GoName); ok {
			f.ReturnType.CGoName = e.Receiver.Type.CGoName
		}

		// Prepare the receiver
		var rt Receiver
		if len(f.Parameters) > 0 {
			rt.Name = commonReceiverName(f.Parameters[0].Type.GoName)
			rt.CName = f.Parameters[0].CName
			rt.Type = f.Parameters[0].Type
		} else {
			if e, ok := g.HasEnum(f.ReturnType.GoName); ok {
				rt.Type = e.Receiver.Type
			} else if s, ok := g.HasStruct(f.ReturnType.GoName); ok {
				rt.Type.GoName = s.Name
			}
		}

		// Check upfront if we can handle a function
		found := false

		for _, p := range f.Parameters {
			if g.api.FilterFunctionParameter != nil && !g.api.FilterFunctionParameter(p) {
				continue
			}

			// Return arguments and slices are always ok since we mark them earlier
			if p.Type.IsReturnArgument || p.Type.IsSlice {
				continue
			}

			if (!g.IsEnumOrStruct(p.Type.GoName) && !p.Type.IsPrimitive) || p.Type.PointerLevel != 0 {
				found = true

				fmt.Printf("Cannot handle parameter %s -> %#v\n", f.CName, p)

				break
			}
		}

		// If we find a heuristic to add the function, add it!
		added := false
		if !found {
			added = g.addBasicMethods(f, fname, "", rt)

			if !added {
				if s := strings.Split(f.Name, "_"); len(s) == 2 {
					if s[0] == rt.Type.GoName {
						rtc := rt
						rtc.Name = s[0]

						added = g.addBasicMethods(f, s[1], "", rtc)
					} else {
						if s[0] != "" {
							s[0] += "_"
						}

						added = g.addBasicMethods(f, strings.Join(s[1:], ""), s[0], rt)
					}
				}
			}

			if !added {
				fname = TrimCommonFunctionName(fname, rt.Type)
				if fn := strings.TrimPrefix(fname, f.ReturnType.GoName+"_"); len(fn) != len(fname) {
					fname = fn
				}

				added = g.addMethod(f, fname, "", rt)

				if !added && g.IsEnumOrStruct(f.ReturnType.GoName) {
					rtc := rt
					rtc.Type = f.ReturnType

					fname = TrimCommonFunctionNamePrefix(fname)
					if fname == "" {
						fname = f.ReturnType.GoName
					}

					if strings.HasPrefix(f.CName, "clang_create") || strings.HasPrefix(f.CName, "clang_get") || strings.HasSuffix(f.CName, "_create") {
						fname = "New" + fname
					}

					added = g.addMethod(f, fname, "", rtc)
				}
				if !added {
					f.Name = UpperFirstCharacter(f.Name)

					clangFile.Functions = append(clangFile.Functions, f)

					added = true
				}
			}
		}

		if !added {
			fmt.Println("Unused function:", f.Name)
		}
	}

	for _, e := range g.enums {
		if err := e.addEnumStringMethods(); err != nil {
			return fmt.Errorf("Cannot generate enum string methods: %v", err)
		}

		for i, m := range e.Methods {
			e.Methods[i] = g.generateMethod(e.Name, m)

			switch m := m.(type) {
			case *Function:
				if len(m.Parameters) > 0 && !m.Parameters[0].Type.IsSlice && m.Parameters[0].Type.GoName == e.Name {
					m.Receiver = Receiver{
						Name: commonReceiverName(e.Name),
						Type: m.Parameters[0].Type,
					}

					g.setIsPointerComposition(&m.Receiver.Type)
				}
				for i := range m.Parameters {
					g.setIsPointerComposition(&m.Parameters[i].Type)
				}

				e.IncludeFiles.unifyIncludeFiles(m.IncludeFiles)

				e.Methods[i] = m.generate()
			}
		}

		if err := e.generate(); err != nil {
			return fmt.Errorf("Cannot generate enum: %v", err)
		}
	}

	for _, s := range g.structs {
		if err := s.addMemberGetters(); err != nil {
			return fmt.Errorf("Cannot generate struct member getters: %v", err)
		}

		for i, m := range s.Methods {
			s.Methods[i] = g.generateMethod(s.Name, m)

			switch m := m.(type) {
			case *Function:
				s.IncludeFiles.unifyIncludeFiles(m.IncludeFiles)
			}
		}

		if err := s.generate(); err != nil {
			return fmt.Errorf("Cannot generate struct: %v", err)
		}
	}

	if len(clangFile.Functions) > 0 {
		for _, m := range clangFile.Functions {
			switch m := m.(type) {
			case *Function:
				clangFile.IncludeFiles.unifyIncludeFiles(m.IncludeFiles)
			}
		}

		for i, m := range clangFile.Functions {
			clangFile.Functions[i] = g.generateMethod("", m)
		}

		if err := clangFile.generate(); err != nil {
			return fmt.Errorf("Cannot generate clang file: %v", err)
		}
	}

	return nil
}

func (g *Generation) generateMethod(receiverName string, m interface{}) string {
	switch m := m.(type) {
	case *Function:
		if g.api.FixedFunctionName != nil {
			if fname := g.api.FixedFunctionName(m); fname != "" {
				m.Name = fname
			}
		}

		if len(m.Parameters) > 0 && !m.Parameters[0].Type.IsSlice && m.Parameters[0].Type.GoName == receiverName {
			m.Receiver = Receiver{
				Name: commonReceiverName(receiverName),
				Type: m.Parameters[0].Type,
			}

			g.setIsPointerComposition(&m.Receiver.Type)
		}
		for i := range m.Parameters {
			g.setIsPointerComposition(&m.Parameters[i].Type)
		}
		g.setIsPointerComposition(&m.ReturnType)

		return m.generate()
	case string:
		return m
	default:
		panic(fmt.Sprintf("Cannot handle %v in handleMethod", m))
	}
}

func (g *Generation) setIsPointerComposition(typ *Type) {
	if s, ok := g.HasStruct(typ.GoName); ok && s.IsPointerComposition {
		typ.IsPointerComposition = true
	}
}

func (g *Generation) addMethod(f *Function, fname string, fnamePrefix string, rt Receiver) bool {
	fname = UpperFirstCharacter(fnamePrefix + UpperFirstCharacter(fname))

	if e, ok := g.HasEnum(rt.Type.GoName); ok {
		f.Name = fname

		e.Methods = append(e.Methods, f)

		return true
	} else if s, ok := g.HasStruct(rt.Type.GoName); ok && s.CName != "CXString" {
		f.Name = fname

		if !rt.Type.IsSlice && rt.Type.PointerLevel > 0 {
			s.IsPointerComposition = true
		}

		s.Methods = append(s.Methods, f)

		return true
	}

	return false
}

func (g *Generation) addBasicMethods(f *Function, fname string, fnamePrefix string, rt Receiver) bool {
	if len(f.Parameters) == 0 && g.IsEnumOrStruct(f.ReturnType.GoName) {
		rtc := rt
		rtc.Type = f.ReturnType

		fname = TrimCommonFunctionNamePrefix(fname)
		if fname == "" {
			fname = f.ReturnType.GoName
		}

		if strings.HasPrefix(f.CName, "clang_create") || strings.HasPrefix(f.CName, "clang_get") || strings.HasSuffix(f.CName, "_create") {
			fname = "New" + fname
		}

		return g.addMethod(f, fname, fnamePrefix, rtc)
	} else if (fname[0] == 'i' && fname[1] == 's' && unicode.IsUpper(rune(fname[2]))) || (fname[0] == 'h' && fname[1] == 'a' && fname[2] == 's' && unicode.IsUpper(rune(fname[3]))) {
		f.ReturnType.GoName = "bool"

		return g.addMethod(f, fname, fnamePrefix, rt)
	} else if len(f.Parameters) == 1 && g.IsEnumOrStruct(f.Parameters[0].Type.GoName) && strings.HasPrefix(fname, "dispose") && f.ReturnType.GoName == "void" {
		fname = "Dispose"

		return g.addMethod(f, fname, fnamePrefix, rt)
	} else if len(f.Parameters) == 2 && strings.HasPrefix(fname, "equal") && g.IsEnumOrStruct(f.Parameters[0].Type.GoName) && f.Parameters[0].Type == f.Parameters[1].Type {
		fname = "Equal"
		f.Parameters[0].Name = commonReceiverName(f.Parameters[0].Type.GoName)
		f.Parameters[1].Name = f.Parameters[0].Name + "2"

		f.ReturnType.GoName = "bool"

		return g.addMethod(f, fname, fnamePrefix, rt)
	}

	return false
}
