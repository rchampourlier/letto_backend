package util

import (
	"fmt"
	"regexp"
	"unicode"
)

// Parameterize converts the given string to snake case (acronyms are converted to lower-case
// and preceded by the specified separator rune). Spaces are converted to the separator too.
func Parameterize(s string, sep rune) string {
	in := []rune(s)
	isLower := func(idx int) bool {
		return idx >= 0 && idx < len(in) && unicode.IsLower(in[idx])
	}

	out := make([]rune, 0, len(in)+len(in)/2)
	for i, r := range in {
		if unicode.IsUpper(r) {
			r = unicode.ToLower(r)
			if i > 0 && in[i-1] != sep && (isLower(i-1) || isLower(i+1)) {
				out = append(out, sep)
			}
		}
		out = append(out, r)
	}

	re := regexp.MustCompile(`[\s_-]+`)
	str := fmt.Sprintf("%c", sep)
	return re.ReplaceAllString(string(out), str)
}
