package main

import "testing"

func BenchmarkGetCandidates(b *testing.B) {
	orig := isValidPath
	defer func() { isValidPath = orig }()
	isValidPath = func(p string) bool {
		return true
	}
	entries := generateEntries()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		getCandidates(entries, []string{"foo", "bar"}, maxCompleteOptions)
	}
}

func generateEntries() []*entry {
	paths := []string{
		"/home/tester", "/home/tester/projects",
		"/foo/bar/baz", "/foo/bazar",
		"/tmp", "/foo/gxxbazabc",
		"/tmp/abc", "/tmp/def",
	}
	var entries []*entry
	for _, p := range paths {
		entries = append(entries, &entry{p, 1.0})
	}
	return entries
}

func TestGetCandidatesShouldRemoveDuplication(t *testing.T) {
	orig := isValidPath
	defer func() { isValidPath = orig }()
	isValidPath = func(p string) bool {
		return true
	}

	orig1, orig2, orig3 := matchConsecutive, matchFuzzy, matchAnywhere
	var dummyMatcher = func(entries []*entry, args []string) []string {
		return []string{"path1", "path2"}
	}
	matchConsecutive = dummyMatcher
	matchFuzzy = dummyMatcher
	matchAnywhere = dummyMatcher
	defer func() {
		matchConsecutive, matchFuzzy, matchAnywhere = orig1, orig2, orig3
	}()

	entries := []*entry{{"path1", 10}}
	result := getCandidates(entries, []string{"foo"}, 4)
	expected := []string{"path1", "path2"}
	assertItemsEqual(t, result, expected)
}

func TestGetCandidates(t *testing.T) {
	orig := isValidPath
	defer func() { isValidPath = orig }()
	isValidPath = func(p string) bool {
		return true
	}

	paths := []string{
		"/home/tester", "/home/tester/projects",
		"/foo/bar/baz", "/foo/bazar",
		"/tmp", "/foo/gxxbazabc",
		"/tmp/abc", "/tmp/def",
	}
	var entries []*entry
	for _, p := range paths {
		entries = append(entries, &entry{p, 1.0})
	}

	result := getCandidates(entries, []string{"foo", "bar"}, 2)
	expected := []string{
		"/foo/bazar",
		"/foo/bar/baz",
	}
	assertItemsEqual(t, result, expected)
}

func TestAnywhere(t *testing.T) {
	entries := []*entry{
		{"/foo/bar/baz", 10},
		{"/foo/bazar", 10},
		{"/tmp", 10},
		{"/foo/gxxbazabc", 10},
	}
	result := matchAnywhere(entries, []string{"foo", "baz"})
	expected := []string{
		"/foo/bar/baz",
		"/foo/bazar",
		"/foo/gxxbazabc",
	}
	assertItemsEqual(t, result, expected)
}

func TestExactName(t *testing.T) {
	entries := []*entry{
		{"/app/open/tidb", 10},
		{"/app/open/redis", 10},
		{"/foo/redis-sdk/bazar", 10},
		{"/tmp", 10},
		{"/foo/tidb/gxxbazabc", 10},
	}
	t.Run("Should returns empty result if the number of args is not exactly one", func(t *testing.T) {
		result := matchExactName(entries, []string{"tidb", "baz"})
		assertItemsEqual(t, result, []string{})
	})
	t.Run("Should only match last part of name", func(t *testing.T) {
		result := matchExactName(entries, []string{"tidb"})
		assertItemsEqual(t, result, []string{"/app/open/tidb"})
	})
}

func TestFuzzy(t *testing.T) {
	entries := []*entry{
		{"/foo/bar/baz", 10},
		{"/foo/bazar", 10},
		{"/tmp", 10},
		{"/foo/gxxbazabc", 10},
	}
	result := matchFuzzy(entries, []string{"baz"})
	expected := []string{
		"/foo/bar/baz",
		"/foo/bazar",
	}
	assertItemsEqual(t, result, expected)
}

func TestConsecutive(t *testing.T) {
	entries := []*entry{
		{"/foo/bar/baz", 10},
		{"/foo/baz/moo", 10},
		{"/moo/foo/Baz", 10},
		{"/foo/bazar", 10},
		{"/foo/xxbaz", 10},
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
