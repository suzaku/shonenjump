package main

import "testing"

func TestConsecutive(t *testing.T) {
    entries := []*Entry{
        &Entry{"/foo/bar/baz", 10},
        &Entry{"/foo/baz/moo", 10},
        &Entry{"/moo/foo/Baz", 10},
        &Entry{"/foo/bazar", 10},
        &Entry{"/foo/xxbaz", 10},
    }
    result := matchConsecutive(entries, []string{"foo", "baz"})
    expected := []string{
        "/moo/foo/Baz",
        "/foo/bazar",
        "/foo/xxbaz",
    }
    if len(result) != len(expected) {
        t.Errorf("Expected %v, got %v", expected, result)
    }
    for i, r := range result {
        if expected[i] != r {
            t.Errorf("Got unexpected element in index %d: %v", i, r)
        }
    }
}
