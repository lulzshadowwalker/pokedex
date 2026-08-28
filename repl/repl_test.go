package repl_test

import (
	"testing"

	"github.com/lulzshadowwalker/pokedex/repl"
)

func TestSplit(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
	}

	for _, c := range cases {
		actual := repl.Split(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("got %d items want %d items", len(actual), len(c.expected))
			continue
		}

		for i, e := range c.expected {
			if e != actual[i] {
				t.Errorf("got %q want %q", actual[i], e)
				t.Fail()
			}
		}
	}
}
