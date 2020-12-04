package jump

import (
	"bufio"
	"fmt"
	"log"
	"os"
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
	candidates := GetCandidates(entries, args, index)
	if len(candidates) == index {
		return candidates[index-1], nil
	}
	return defaultPath, nil
}

func (s Store) saveEntries(entries entryList) error {
	return entries.Save(s.path)
}
