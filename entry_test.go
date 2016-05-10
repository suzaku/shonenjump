package main

import (
	"strings"
	"testing"
)

func TestDecrementScoreOfEntries(t *testing.T) {
	entries := []*entry{
		&entry{"a", 20},
		&entry{"b", 10},
		&entry{"c", 0},
	}
	decrementScoreOfEntries(entries)
	expected := []float64{19.0, 9.0, 0}
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
		output = append(output, r.Path)
	}
	expected := []string{"/foo/bar", "/tmp"}
	assertItemsEqual(t, output, expected)
}
