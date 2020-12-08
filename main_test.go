package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, test.needle, needle)
		assert.Equal(t, test.index, index)
		assert.Equal(t, test.path, path)
	}
}
