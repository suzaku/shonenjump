package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Entry correspond to a line in the data file
type entry struct {
	Path  string
	Score float64
}

func (e *entry) updateScore(weight float64) float64 {
	e.Score = math.Sqrt(math.Pow(e.Score, 2) + math.Pow(weight, 2))
	return e.Score
}

func (e entry) String() string {
	return fmt.Sprintf("%.2f\t%s", e.Score, e.Path)
}

type entryList []*entry

func (a entryList) Len() int {
	return len(a)
}

func (a entryList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a entryList) Less(i, j int) bool {
	return a[i].Score < a[j].Score
}

func (a entryList) Sort() {
	sort.Sort(sort.Reverse(a))
}

func (entries entryList) Update(path string, weight float64) entryList {
	// normalize the input
	path = strings.TrimSuffix(path, string(os.PathSeparator))
	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	var ent *entry
	for _, e := range entries {
		if e.Path == path {
			ent = e
			break
		}
	}
	if ent == nil {
		ent = &entry{path, 0}
		entries = append(entries, ent)
	}
	ent.updateScore(weight)

	entries.Sort()

	return entries
}

func (entries entryList) Filter(f func(*entry) bool) entryList {
	var filtered entryList
	for _, e := range entries {
		if f(e) {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

func (entries entryList) Save(path string) {
	tempfile, err := ioutil.TempFile("", "shonenjump")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tempfile.Name())

	var buffer bytes.Buffer
	for _, e := range entries {
		buffer.WriteString(e.String() + "\n")
	}
	if _, err := tempfile.Write(buffer.Bytes()); err != nil {
		log.Fatal(err)
	}
	if err := tempfile.Close(); err != nil {
		log.Fatal(err)
	}

	if err = os.MkdirAll(filepath.Dir(path), 0740); err != nil {
		log.Fatal(err)
	}

	if err = os.Rename(tempfile.Name(), path); err != nil {
		log.Fatal(err)
	}
}

// As entries get older, their scores become lower.
func (entries entryList) Age() {
	for _, e := range entries {
		e.Score = math.Max(e.Score-1, 0)
	}
}

func parseEntry(s string) (ent entry, err error) {
	parts := strings.Split(s, "\t")
	score, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return
	}
	ent = entry{parts[1], score}
	return ent, nil
}

func loadEntries(path string) []*entry {
	var entries []*entry
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return entries
		}
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		entry, err := parseEntry(line)
		if err != nil {
			log.Printf("Failed to parse score from line: %v", line)
			continue
		}
		entries = append(entries, &entry)
	}
	entryList(entries).Sort()
	return entries
}
