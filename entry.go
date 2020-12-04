package main

import (
	"bufio"
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

type Store struct {
	path string
}

func (s Store) AddPath(pathToAdd string) error {
	path, err := preprocessPath(pathToAdd)
	if err != nil {
		return err
	}
	if !isValidPath(path) {
		return fmt.Errorf("invalid path: %v", path)
	}
	oldEntries, err := s.ReadEntries()
	if err != nil {
		return err
	}
	oldEntries.Age()
	newEntries := oldEntries.Update(path, defaultWeight)
	return s.saveEntries(newEntries)
}

func (s Store) ReadEntries() (entryList, error) {
	var entries entryList
	file, err := os.Open(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return entries, nil
		}
		return entries, err
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
	return entries, nil
}

func (s Store) Cleanup() error {
	entries, err := s.ReadEntries()
	if err != nil {
		return err
	}
	entries, changed := clearNotExistDirs(entries)
	if changed {
		return entries.Save(s.path)
	}
	return nil
}

func (s Store) GetNthCandidate(args []string, index int, defaultPath string) (string, error) {
	entries, err := s.ReadEntries()
	if err != nil {
		return "", err
	}
	candidates := getCandidates(entries, args, index)
	if len(candidates) == index {
		return candidates[index-1], nil
	}
	return defaultPath, nil
}

func (s Store) saveEntries(entries entryList) error {
	return entries.Save(s.path)
}

func clearNotExistDirs(entries entryList) (result entryList, changed bool) {
	for _, e := range entries {
		if isValidPath(e.val) {
			result = append(result, e)
		} else {
			log.Printf("Directory %s no longer exists", e.val)
			changed = true
		}
	}
	return result, changed
}

func NewStore(dataPath string) Store {
	return Store{
		path: dataPath,
	}
}

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

func (entries entryList) Sort() {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].score > entries[j].score
	})
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

func (entries entryList) Save(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0740); err != nil {
		return err
	}

	tempfile, err := ioutil.TempFile(filepath.Dir(path), "shonenjump")
	if err != nil {
		return err
	}
	defer os.Remove(tempfile.Name())

	writer := bufio.NewWriter(tempfile)
	for _, e := range entries {
		if !isValidPath(e.val) {
			continue
		}
		if _, err := writer.WriteString(e.String() + "\n"); err != nil {
			return err
		}
	}
	writer.Flush()

	if err := tempfile.Close(); err != nil {
		return err
	}

	if err = os.Rename(tempfile.Name(), path); err != nil {
		return err
	}

	return nil
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
