package main

import (
    "testing"
    "strings"
)

func TestClearNotExistDirs(t *testing.T) {
    orig := isValidPath
    defer func() { isValidPath = orig }()
    isValidPath = func(p string) bool {
        return !strings.HasSuffix(p, "not-exist")
    }
    entries := []*Entry{
        &Entry{"/foo/bar", 10},
        &Entry{"/foo/not-exist", 10},
        &Entry{"/tmp", 10},
        &Entry{"/not-exist", 10},
    }
    result := clearNotExistDirs(entries)
    var output []string
    for _, r := range result {
        output = append(output, r.Path)
    }
    expected := []string{"/foo/bar", "/tmp"}
    assertItemsEqual(t, output, expected)
}
