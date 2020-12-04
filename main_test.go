package main

import (
	"testing"
)

func TestParseCompleteOption(t *testing.T) {
	tests := []struct {
		input  string
		needle string
		index  int
		path   string
	}{
		{"abc__7__/home/tester/abc", "abc", 7, "/home/tester/abc"},
		{"abc__7__", "abc", 7, ""},
		{"abc__7", "abc", 7, ""},
		{"abc__", "abc", 0, ""},
	}
	for _, test := range tests {
		needle, index, path := parseCompleteOption(test.input)
		if !(test.needle == needle && test.index == index && test.path == path) {
			t.Errorf("Unexpected parse result for %s: (%s, %d, %s)", test.input, needle, index, path)
		}
	}
}
