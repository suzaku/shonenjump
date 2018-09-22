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
	val   string
	score float64
}

func (e *entry) updateScore(weight float64) float64 {
	e.score = math.Sqrt(math.Pow(e.score, 2) + math.Pow(weight, 2))
	return e.score
}

func (e entry) String() string {
	return fmt.Sprintf("%.2f\t%s", e.score, e.val)
}

type entryList []*entry

func (entries entryList) Len() int {
	return len(entries)
}

func (entries entryList) Swap(i, j int) {
	entries[i], entries[j] = entries[j], entries[i]
}

func (entries entryList) Less(i, j int) bool {
	return entries[i].score < entries[j].score
}

func (entries entryList) Sort() {
	sort.Sort(sort.Reverse(entries))
}

func (entries entryList) Update(val string, weight float64) entryList {
	var ent *entry
	for _, e := range entries {
		if e.val == val {
			ent = e
			break
		}
	}
	if ent == nil {
		ent = &entry{val, 0}
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
	tempfile, err := ioutil.TempFile(filepath.Dir(path), "shonenjump")
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
		delta := math.Ceil(e.score / 10)
		e.score = math.Max(e.score-delta, 0)
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

func loadEntries(path string) entryList {
	var entries entryList
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
	entries.Sort()
	return entries
}
