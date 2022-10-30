package jump

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Store struct {
	path string
}

func NewStore(dataPath string) Store {
	return Store{
		path: dataPath,
	}
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

func (s Store) ReadEntries() (EntryList, error) {
	var entries EntryList
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
	if entries != nil {
		entries.Sort()
	}
	return entries, nil
}

func (s Store) topEntry() (entry, error) {
	var ent entry

	file, err := os.Open(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return ent, nil
		}
		return ent, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		ent, err := parseEntry(line)
		if err != nil {
			log.Printf("Failed to parse score from line: %v", line)
			continue
		}
		return ent, nil
	}
	return ent, nil
}

func (s Store) Cleanup() error {
	entries, err := s.ReadEntries()
	if err != nil {
		return err
	}
	entries, changed := clearNotExistDirs(entries)
	if changed {
		return s.saveEntries(entries)
	}
	return nil
}

func (s Store) GetNthCandidate(args []string, index int, defaultPath string) (string, error) {
	entries, err := s.ReadEntries()
	if err != nil {
		return "", err
	}
	candidates := GetCandidates(entries, args, index)
	if len(candidates) == index {
		return candidates[index-1], nil
	}
	return defaultPath, nil
}

func (s Store) GetTopPath(defaultPath string) (string, error) {
	ent, err := s.topEntry()
	if err != nil {
		return "", err
	}
	return ent.val, nil
}

func (s Store) saveEntries(entries EntryList) error {
	folder := filepath.Dir(s.path)
	if err := os.MkdirAll(folder, 0740); err != nil {
		return err
	}

	tempfile, err := os.CreateTemp(folder, "shonenjump")
	if err != nil {
		return err
	}
	defer os.Remove(tempfile.Name())

	writer := bufio.NewWriter(tempfile)
	for _, e := range entries {
		if !isValidPath(e.val) {
			continue
		}
		if _, err := fmt.Fprintln(writer, e); err != nil {
			return err
		}
	}
	if err := writer.Flush(); err != nil {
		return err
	}

	if err := tempfile.Close(); err != nil {
		return err
	}

	if err = os.Rename(tempfile.Name(), s.path); err != nil {
		return err
	}

	return nil
}
