package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPreprocessPath(t *testing.T) {
	path, err := preprocessPath("/abc/")
	if err != nil {
		t.Error(err)
	}
	if path != "/abc" {
		t.Errorf("Trailing slash is not removed in path: %s", path)
	}
	path, err = preprocessPath("abc")
	if err != nil {
		t.Error(err)
	}
	pwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	if path != filepath.Join(pwd, "abc") {
		t.Errorf("Relative path not converted to absolute path: %s", path)
	}
}

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

func TestClearNotExistDirs(t *testing.T) {
	orig := isValidPath
	defer func() { isValidPath = orig }()
	isValidPath = func(p string) bool {
		return !strings.HasSuffix(p, "not-exist")
	}
	entries := []*entry{
		{"/foo/bar", 10},
		{"/foo/not-exist", 10},
		{"/tmp", 10},
		{"/not-exist", 10},
	}
	result := clearNotExistDirs(entries)
	var output []string
	for _, r := range result {
		output = append(output, r.Val)
	}
	expected := []string{"/foo/bar", "/tmp"}
	assertItemsEqual(t, output, expected)
}
