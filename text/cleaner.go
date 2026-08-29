package text

import (
	"strings"
	"unicode"
)

func Clean(s string) string { // @note to be deprecated. Taken into Options
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if unicode.IsLetter(r) {
			b.WriteRune(unicode.ToLower(r))
		}
	}
	return b.String()
}
