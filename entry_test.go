package main

import (
	"testing"
)

func TestEntryListSort(t *testing.T) {
	rawEntries := []*entry{
		&entry{"b", 10},
		&entry{"a", 20},
	}
	entries := entryList(rawEntries)
	entries.Sort()
	expected := []string{"a", "b"}
	for i, e := range rawEntries {
		if expected[i] != e.Path {
			t.Errorf("Item %d not in place, expected %s, got %s", expected[i], e.Path)
		}
	}
}

func TestEntryListFilter(t *testing.T) {
	entries := entryList{
		&entry{"/path_b", 10},
		&entry{"/path_a", 0},
	}
	nonZeroScore := func(e *entry) bool { return e.Score > 0 }
	nonZeroEntries := entries.Filter(nonZeroScore)
	if len(nonZeroEntries) != 1 {
		t.Errorf("Entries not filtered correctly: %v", nonZeroEntries)
	}
	if nonZeroEntries[0].Path != "/path_b" {
		t.Errorf("Incorrect entry left after filtering: %v", nonZeroEntries)
	}
}

func TestEntryListUpdate(t *testing.T) {
	entries := entryList{
		&entry{"/path_b", 10},
		&entry{"/path_a", 0},
	}
	entries = entries.Update("/path_a", 1)
	if entries[0].Score != 10 || entries[1].Score != 1 {
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
		if e.Score != expected[i] {
			t.Errorf("Score not updated correctly, expect %f, get %f", expected[i], e.Score)
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
	if e.Score != 10 {
		t.Errorf("Entity score is wrong: %f", e.Score)
	}
	e.updateScore(10)
	if e.Score-14.14 < 0.001 {
		t.Errorf("Entity score is wrong: %f", e.Score)
	}
}
