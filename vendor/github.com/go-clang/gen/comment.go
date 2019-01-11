package gen

import (
	"regexp"
	"strings"
)

var (
	reReplaceCComments  = regexp.MustCompile(`[ \t]*\n[ \t]*\*[ \t]*`)
	reReplaceMultilines = regexp.MustCompile(`[ \t]*\n[ \t]*`)
)

func CleanDoxygenComment(comment string) string {
	// Remove C style comment
	comment = strings.TrimPrefix(comment, "/**")
	comment = strings.TrimSuffix(comment, "*/")
	comment = strings.TrimPrefix(comment, "//")
	comment = reReplaceCComments.ReplaceAllString(comment, "\n")

	// Replace some tags
	comment = strings.Replace(comment, "\\brief ", "", -1)
	comment = strings.Replace(comment, "\\c ", "", -1)
	comment = strings.Replace(comment, "\\param ", "Parameter ", -1)
	comment = strings.Replace(comment, "\\returns ", "Returns ", -1)

	// If there is no empty line we definitely do the comment in one line
	if !strings.Contains(comment, "\n\n") {
		comment = strings.Replace(comment, "\n", " ", -1)
	}

	// Replace spaces
	comment = strings.Replace(comment, "  ", " ", -1)
	comment = strings.TrimSpace(comment)

	if comment == "" {
		return comment
	}

	// This might be a bug in the parsing by clang. The nearest comment is added to an item as its comment even tough it is not directly the item's comment
	if strings.HasPrefix(comment, "\\defgroup") {
		return ""
	}

	// Indent multiline comments
	if strings.ContainsRune(comment, '\n') {
		comment = reReplaceMultilines.ReplaceAllString(comment, "\n\t")

		return "/*\n\t" + comment + "\n*/"
	} else {
		return "// " + comment
	}
}
