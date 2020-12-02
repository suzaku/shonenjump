package main

import (
	"bufio"
	"io/ioutil"
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

	entries.Save(fileName)

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

	entries.Save(fileName)

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
