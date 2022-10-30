package jump

import (
	"bufio"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveEntries(t *testing.T) {
	dir, err := os.MkdirTemp("", "test")
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
		err := os.MkdirAll(e.val, 0664)
		assert.Nil(t, err)
	}
	// Append a non-exist dir that should be ignored
	rawEntries = append(rawEntries, &entry{val: "non-exist", score: 15})
	entries := EntryList(rawEntries)

	fileName := filepath.Join(dir, "testEntries")
	store := NewStore(fileName)

	err = store.saveEntries(entries)
	assert.Nil(t, err)

	entriesFile, err := os.Open(fileName)
	assert.Nil(t, err)

	scanner := bufio.NewScanner(entriesFile)
	var results []string
	for scanner.Scan() {
		line := scanner.Text()
		results = append(results, line)
	}
	assert.Equal(t, len(entries)-1, len(results), "Incorrect number of entries saved")

	for i, r := range results {
		assert.Equal(t, entries[i].String(), r)
		err := os.Remove(entries[i].val)
		assert.Nil(t, err)
	}

	err = store.saveEntries(entries)
	assert.Nil(t, err)

	content, err := os.ReadFile(fileName)
	assert.Nil(t, err)

	assert.Empty(t, content)
}
