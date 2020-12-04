package jump

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestEntryListSave(t *testing.T) {
	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	rawEntries := []*entry{
		{filepath.Join(dir, "b"), 10},
		{filepath.Join(dir, "a"), 20},
		{filepath.Join(dir, "c"), 15},
	}
	for _, e := range rawEntries {
		if err := os.MkdirAll(e.val, 0664); err != nil {
			t.Fatal(err)
		}
	}
	// Append a non-exist dir that should be ignored
	rawEntries = append(rawEntries, &entry{val: "non-exist", score: 15})
	entries := entryList(rawEntries)

	fileName := filepath.Join(dir, "testEntries")

	if err := entries.Save(fileName); err != nil {
		log.Fatal(err)
	}

	entriesFile, err := os.Open(fileName)
	if err != nil {
		t.Fatal(err)
	}
	scanner := bufio.NewScanner(entriesFile)
	var results []string
	for scanner.Scan() {
		line := scanner.Text()
		results = append(results, line)
	}
	if len(results) != len(entries)-1 {
		t.Errorf("Incorrect number of entries saved: %q", results)
	}
	for i, r := range results {
		if r != entries[i].String() {
			t.Errorf("Entry %d saved incorrectly: %q", i, results[i])
		}
		if err := os.Remove(entries[i].val); err != nil {
			t.Fatal(err)
		}
	}

	if err := entries.Save(fileName); err != nil {
		log.Fatal(err)
	}

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	if len(content) != 0 {
		t.Errorf("Expected empty content, got: %v", string(content))
	}
}

func TestEntryListSort(t *testing.T) {
	rawEntries := []*entry{
		{"b", 10},
		{"a", 20},
		{"c", 15},
	}
	entries := entryList(rawEntries)
	entries.Sort()
	expected := []string{"a", "c", "b"}
	for i, e := range entries {
		if expected[i] != e.val {
			t.Errorf("Item %d not in place, expected %s, got %s", i, expected[i], e.val)
		}
	}
}

func TestEntryListUpdate(t *testing.T) {
	entries := entryList{
		&entry{"/path_b", 10},
		&entry{"/path_a", 0},
	}
	entries = entries.Update("/path_a", 1)
	if entries[0].score != 10 || entries[1].score != 1 {
		t.Errorf("Invalid update: %v", entries)
	}
	entries = entries.Update("/path_c", 1)
	if len(entries) != 3 {
		t.Errorf("New entry not created: %d", len(entries))
	}
}

func TestEntryListAge(t *testing.T) {
	entries := entryList{
		&entry{"a", 20},
		&entry{"b", 10},
		&entry{"c", 0},
	}
	entries.Age()
	expected := []float64{18.0, 9.0, 0}
	for i, e := range entries {
		if e.score != expected[i] {
			t.Errorf("Score not updated correctly, expect %f, get %f", expected[i], e.score)
		}
	}
}

func TestString(t *testing.T) {
	e := &entry{"/etc/init", 10.1234}
	if e.String() != "10.12\t/etc/init" {
		t.Errorf("Wrong string representation: %s", e.String())
	}
}

func TestUpdateEntryScore(t *testing.T) {
	e := &entry{"/etc/init", 0}
	e.updateScore(10)
	if e.score != 10 {
		t.Errorf("Entity score is wrong: %f", e.score)
	}
	e.updateScore(10)
	if e.score-14.14 < 0.001 {
		t.Errorf("Entity score is wrong: %f", e.score)
	}
}

func TestLoadEntries(t *testing.T) {
	t.Run("Should make sure the entries are sorted by score", func(t *testing.T) {
		f, err := ioutil.TempFile("", "entries")
		if err != nil {
			t.Fatal(err)
		}
		content := []byte("20\t/a\n22\t/a/b\n12.5\t/c\n")
		err = ioutil.WriteFile(f.Name(), content, 0666)
		if err != nil {
			t.Fatal(err)
		}
		store := NewStore(f.Name())
		entries, err := store.ReadEntries()
		if err != nil {
			t.Fatal(err)
		}
		if len(entries) != 3 {
			t.Errorf("Expect 3 entries, got: %v ", entries)
		}
		paths := make([]string, len(entries))
		for i, e := range entries {
			paths[i] = e.val
		}
		assertItemsEqual(t, paths, []string{"/a/b", "/a", "/c"})
	})
}

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

func TestClearNotExistDirs(t *testing.T) {
	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		basename string
		create   bool
	}{
		{basename: "abc", create: true},
		{basename: "bat", create: false},
		{basename: "world", create: false},
		{basename: "super", create: true},
		{basename: "wrong", create: false},
	}
	var expected []string
	var entries entryList
	for _, c := range cases {
		e := &entry{val: filepath.Join(dir, c.basename)}
		if c.create {
			if err := os.MkdirAll(e.val, 0644); err != nil {
				log.Fatal(err)
			}
			expected = append(expected, e.val)
		}
		entries = append(entries, e)
	}
	result, changed := clearNotExistDirs(entries)
	var output []string
	for _, r := range result {
		output = append(output, r.val)
	}
	assertItemsEqual(t, output, expected)
	if !changed {
		t.Error("Empty dirs get deleted, but changed is false.")
	}
}
