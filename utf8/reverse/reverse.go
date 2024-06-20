//go:build !solution

package reverse

import (
	"strings"
	"unicode/utf8"
)

func Reverse(input string) string {
	var builder strings.Builder
	builder.Grow(len(input))

	for i := len(input); i > 0; {
		r, size := utf8.DecodeLastRuneInString(input[:i])
		builder.WriteRune(r)
		i -= size
	}

	return builder.String()
}
