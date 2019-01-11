package gen

import (
	"strings"
	"unicode"
)

func LowerFirstCharacter(s string) string {
	r := []rune(s)

	r[0] = unicode.ToLower(r[0])

	return string(r)
}

func UpperFirstCharacter(s string) string {
	r := []rune(s)

	r[0] = unicode.ToUpper(r[0])

	return string(r)
}

var goKeywordReplacements = map[string]string{
	"range": "r",
}

func ReplaceGoKeywords(s string) string {
	if r, ok := goKeywordReplacements[s]; ok {
		return r
	}

	return ""
}

func TrimCommonFunctionName(name string, typ Type) string {
	name = TrimCommonFunctionNamePrefix(name)

	if fn := strings.TrimPrefix(name, typ.GoName+"_"); len(fn) != len(name) {
		name = fn
	} else if fn := strings.TrimPrefix(name, typ.GoName); len(fn) != len(name) {
		name = fn
	}
	if tkn := strings.TrimSuffix(typ.GoName, "Kind"); len(tkn) != len(typ.GoName) {
		if fn := strings.TrimPrefix(name, tkn+"_"); len(fn) != len(name) {
			name = fn
		} else if fn := strings.TrimPrefix(name, tkn); len(fn) != len(name) {
			name = fn
		}
	}

	name = TrimCommonFunctionNamePrefix(name)

	// If the function name is empty at this point, it is a constructor
	if name == "" {
		name = typ.GoName
	}

	return name
}

func TrimCommonFunctionNamePrefix(name string) string {
	name = strings.TrimPrefix(name, "create")
	name = strings.TrimPrefix(name, "get")
	if len(name) > 4 && unicode.IsUpper(rune(name[3])) {
		name = strings.TrimPrefix(name, "Get")
	}

	name = TrimLanguagePrefix(name)

	return name
}

func TrimLanguagePrefix(name string) string {
	name = strings.TrimPrefix(name, "CX_CXX")
	name = strings.TrimPrefix(name, "CXX")
	name = strings.TrimPrefix(name, "CX")
	name = strings.TrimPrefix(name, "ObjC")
	name = strings.TrimPrefix(name, "_")

	return name
}

func commonReceiverName(s string) string {
	var n []rune

	for _, c := range s {
		if unicode.IsUpper(c) {
			n = append(n, unicode.ToLower(c))
		}
	}

	return string(n)
}
