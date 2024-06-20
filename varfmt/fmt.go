//go:build !solution

package varfmt

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func Sprintf(format string, args ...interface{}) string {
	var result strings.Builder
	var cnt, num int
	numSize := -1

	argsList := make([]string, 0, len(args))
	argsSize := 0
	for _, arg := range args {
		s := fmt.Sprint(arg)
		argsList = append(argsList, s)
		argsSize += len(s)
	}
	result.Grow(len(format) + argsSize)

	for len(format) > 0 {
		r, size := utf8.DecodeRuneInString(format)
		format = format[size:]
		switch true {
		case r == '{':
			numSize = 0
		case r == '}':
			index := num
			if numSize <= 0 {
				index = cnt
			}
			result.WriteString(argsList[index])
			num = 0
			numSize = -1
			cnt++
		case numSize >= 0:
			num = num*10 + (int(r) - '0')
			numSize++
		default:
			result.WriteRune(r)
		}
	}

	return result.String()
}
