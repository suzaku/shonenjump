package main

import "testing"

func TestAnywhere(t *testing.T) {
    entries := []*Entry{
        &Entry{"/foo/bar/baz", 10},
        &Entry{"/foo/bazar", 10},
        &Entry{"/tmp", 10},
        &Entry{"/foo/gxxbazabc", 10},
    }
    result := matchAnywhere(entries, []string{"foo", "baz"})
    expected := []string{
        "/foo/bar/baz",
        "/foo/bazar",
        "/foo/gxxbazabc",
    }
    assertItemsEqual(t, result, expected)
}

func TestFuzzy(t *testing.T) {
    entries := []*Entry{
        &Entry{"/foo/bar/baz", 10},
        &Entry{"/foo/bazar", 10},
        &Entry{"/tmp", 10},
        &Entry{"/foo/gxxbazabc", 10},
    }
    result := matchFuzzy(entries, []string{"baz"})
    expected := []string{
        "/foo/bar/baz",
        "/foo/bazar",
    }
    assertItemsEqual(t, result, expected)
}

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
    assertItemsEqual(t, result, expected)
}

func assertItemsEqual(t *testing.T, result []string, expected []string) {
    if len(result) != len(expected) {
        t.Errorf("Expected %v, got %v", expected, result)
    }
    for i, r := range result {
        if expected[i] != r {
            t.Errorf("Got unexpected element in index %d: %v", i, r)
        }
    }
}
