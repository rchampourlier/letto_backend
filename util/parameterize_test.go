package util_test

import (
	"strings"
	"testing"

	"gitlab.com/letto/letto_backend/util"
)

type ParameterizeTest struct {
	input  string
	output string
	sep    rune
}

var parameterizeTests = []ParameterizeTest{
	{"a small label", "a-small-label", '-'},
	{"a label with UpperCase", "a-label-with-upper-case", '-'},
	{"downCaseUp_underscore", "down-case-up-underscore", '-'},
	{"multiple  Spaces", "multiple-spaces", '-'},
}

func TestParameterize(t *testing.T) {
	for _, test := range parameterizeTests {
		if util.Parameterize(test.input, test.sep) != test.output {
			t.Errorf("Parameterize(%q) -> %q, want %q", test.input, util.Parameterize(test.input, '-'), test.output)
		}
	}
}

func BenchmarkParameterize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range parameterizeTests {
			util.Parameterize(test.input, '-')
		}
	}
}

func BenchmarkToLower(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range parameterizeTests {
			strings.ToLower(test.input)
		}
	}
}
