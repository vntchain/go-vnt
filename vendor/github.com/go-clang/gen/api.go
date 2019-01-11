package gen

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type API struct {
	// ClangArguments holds the command line arguments for Clang
	ClangArguments []string

	// PrepareFunctionName returns a prepared function name for further processing
	PrepareFunctionName func(g *Generation, f *Function) string
	// PrepareFunction prepares a function for further processing
	PrepareFunction func(f *Function)
	// FilterFunction determines if a function is generateable
	FilterFunction func(f *Function) bool
	// FilterFunctionParameter determines if a function parameter is generateable
	FilterFunctionParameter func(p FunctionParameter) bool
	// FixedFunctionName returns an unempty string if a function needs to receive a specific name
	FixedFunctionName func(f *Function) string

	// PrepareStructMembers is called before adding struct member getters
	PrepareStructMembers func(s *Struct)
	// FilterStructMemberGetter determines if a getter should be generated for a member
	FilterStructMemberGetter func(m *StructMember) bool
}

func (a *API) HandleDirectory(dir string) ([]*HeaderFile, error) {
	headers, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("Cannot list clang-c directory: %v", err)
	}

	var headerFiles []*HeaderFile

	for _, hf := range headers {
		if hf.IsDir() || !strings.HasSuffix(hf.Name(), ".h") {
			continue
		}

		h := newHeaderFile(a, hf.Name(), dir)

		if err := h.parse(a.ClangArguments); err != nil {
			return nil, fmt.Errorf("Cannot handle header file %q: %v", h.FullPath(), err)
		}

		headerFiles = append(headerFiles, h)
	}

	return headerFiles, nil
}
