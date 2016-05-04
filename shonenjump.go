package main

import (
    "fmt"
    "os"
    "bufio"
    "strings"
    "log"
    "strconv"
    "sort"
    "math"
)

func updateEntriesWithPath(entries []*Entry, path string, weight float64) []*Entry {
    var entry *Entry
    for _, e := range entries {
        if e.Path == path {
            entry = e
            break
        }
    }
    if entry == nil {
        entry = &Entry{path, 0}
        entries = append(entries, entry)
    }

    entry.updateScore(weight)

    sortEntriesByScore(entries)

    return entries
}

type Entry struct {
    Path string
    Score float64
}

func (e *Entry) updateScore(weight float64) float64 {
    e.Score = math.Sqrt(math.Pow(e.Score, 2) + math.Pow(weight, 2))
    return e.Score
}

type ByScore []*Entry

func (a ByScore) Len() int {
    return len(a)
}

func (a ByScore) Swap(i, j int) {
    a[i], a[j] = a[j], a[i]
}

func (a ByScore) Less(i, j int) bool {
    return a[i].Score < a[j].Score
}

func parseEntry(s string) (entry Entry, err error) {
    parts := strings.Split(s, "\t")
    score, err := strconv.ParseFloat(parts[0], 64)
    if err != nil {
        return
    }
    entry = Entry{parts[1], score}
    return entry, nil
}

func sortEntriesByScore(entries []*Entry) {
    sort.Sort(sort.Reverse(ByScore(entries)))
}

func loadEntries(path string) []*Entry {
    file, err := os.Open(path)
    if err != nil {
        log.Fatal("Failed to open data file")
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    var entries []*Entry
    for scanner.Scan() {
        line := scanner.Text()
        entry, err := parseEntry(line)
        if err != nil {
            fmt.Errorf("Failed to parse score from line: %v", line)
        }
        entries = append(entries, &entry)
    }
    sortEntriesByScore(entries)
    return entries
}

func main() {
    dataPath := "/Users/satoru/Library/autojump/autojump.txt"
    entries := loadEntries(dataPath)
    path := "/tmp"
    weight := 10.0

    entries = updateEntriesWithPath(entries, path, weight)
    for _, e := range entries {
        fmt.Println(e.Score, e.Path)
    }
}
