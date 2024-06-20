//go:build !solution

package spacecollapse

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func CollapseSpaces(input string) string {
	var builder strings.Builder
	builder.Grow(len(input))

	lastWasSpace := false

	for i := 0; i < len(input); {
		r, size := utf8.DecodeRuneInString(input[i:])
		if r == utf8.RuneError && size == 1 {
			builder.WriteRune(utf8.RuneError)
			i++
		} else {
			if unicode.IsSpace(r) {
				if !lastWasSpace {
					builder.WriteRune(' ')
					lastWasSpace = true
				}
			} else {
				builder.WriteRune(r)
				lastWasSpace = false
			}
			i += size
		}
	}

	return builder.String()
}
