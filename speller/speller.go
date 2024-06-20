//go:build !solution

package speller

import (
	"strings"
)

var (
	ones      = []string{"", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
	teens     = []string{"ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen"}
	tens      = []string{"", "ten", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety"}
	thousands = []string{"", "thousand", "million", "billion", "trillion"}
)

func Spell(n int64) string {
	if n == 0 {
		return "zero"
	}
	if n < 0 {
		return "minus " + Spell(-n)
	}

	var words []string
	for unitIndex := 0; n > 0; unitIndex++ {
		if numPart := n % 1000; numPart > 0 {
			words = append([]string{convertThreeDigits(numPart) + " " + thousands[unitIndex]}, words...)
		}
		n /= 1000
	}

	return strings.TrimSpace(strings.Join(words, " "))
}

func convertThreeDigits(num int64) string {
	var result strings.Builder

	if num >= 100 {
		result.WriteString(ones[num/100] + " hundred ")
		num %= 100
	}

	if num >= 20 {
		if num%10 > 0 {
			result.WriteString(tens[num/10] + "-" + ones[num%10])
		} else {
			result.WriteString(tens[num/10])
		}
	} else if num >= 10 {
		result.WriteString(teens[num-10])
	} else if num > 0 {
		result.WriteString(ones[num])
	}

	return strings.TrimSpace(result.String())
}
