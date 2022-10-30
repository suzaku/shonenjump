package jump

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntryListSort(t *testing.T) {
	rawEntries := []*entry{
		{"b", 10},
		{"a", 20},
		{"c", 15},
	}
	entries := EntryList(rawEntries)
	entries.Sort()
	expected := []string{"a", "c", "b"}
	for i, e := range entries {
		assert.Equal(t, expected[i], e.val)
	}
}

func TestEntryListUpdate(t *testing.T) {
	entries := EntryList{
		&entry{"/path_b", 10},
		&entry{"/path_a", 0},
	}
	entries = entries.Update("/path_a", 1)
	assert.Equal(t, float64(10), entries[0].score)
	assert.Equal(t, float64(1), entries[1].score)

	entries = entries.Update("/path_c", 1)
	assert.Len(t, entries, 3)
}

func TestEntryListAge(t *testing.T) {
	entries := EntryList{
		&entry{"a", 20},
		&entry{"b", 10},
		&entry{"c", 0},
	}
	entries.Age()
	expected := []float64{18.0, 9.0, 0}
	for i, e := range entries {
		assert.Equal(t, expected[i], e.score)
	}
}

func TestString(t *testing.T) {
	e := &entry{"/etc/init", 10.1234}
	assert.Equal(t, "10.12\t/etc/init", e.String())
}

func TestUpdateEntryScore(t *testing.T) {
	e := &entry{"/etc/init", 0}
	e.updateScore(10)
	assert.Equal(t, float64(10), e.score)

	e.updateScore(10)
	assert.InDelta(t, 14.14, e.score, 0.01)
}

func TestLoadEntries(t *testing.T) {
	t.Run("Should make sure the entries are sorted by score", func(t *testing.T) {
		f, err := os.CreateTemp("", "entries")
		assert.Nil(t, err)

		content := []byte("20\t/a\n22\t/a/b\n12.5\t/c\n")
		err = os.WriteFile(f.Name(), content, 0666)
		assert.Nil(t, err)

		store := NewStore(f.Name())
		entries, err := store.ReadEntries()
		assert.Nil(t, err)

		assert.Len(t, entries, 3)
		paths := make([]string, len(entries))
		for i, e := range entries {
			paths[i] = e.val
		}
		assert.Equal(t, []string{"/a/b", "/a", "/c"}, paths)
	})
}

func TestPreprocessPath(t *testing.T) {
	path, err := preprocessPath("/abc/")
	assert.Nil(t, err)
	assert.Equal(t, "/abc", path, "Trailing slash is not removed in path")
	path, err = preprocessPath("abc")
	assert.Nil(t, err)
	pwd, err := os.Getwd()
	assert.Nil(t, err)
	assert.Equal(t, filepath.Join(pwd, "abc"), path)
}

func TestClearNotExistDirs(t *testing.T) {
	dir, err := os.MkdirTemp("", "test")
	assert.Nil(t, err)
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
	var entries EntryList
	for _, c := range cases {
		e := &entry{val: filepath.Join(dir, c.basename)}
		if c.create {
			err := os.MkdirAll(e.val, 0644)
			assert.Nil(t, err)
			expected = append(expected, e.val)
		}
		entries = append(entries, e)
	}
	result, changed := clearNotExistDirs(entries)
	var output []string
	for _, r := range result {
		output = append(output, r.val)
	}
	assert.Equal(t, expected, output)
	assert.True(t, changed, "Empty dirs get deleted, but changed is false.")
}
