package main

import (
    "fmt"
    "os"
    "bufio"
    "io/ioutil"
    "strings"
    "log"
    "strconv"
    "sort"
    "math"
    "bytes"
    "path/filepath"
)

func saveEntries(entries []*Entry, path string) {
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

func updateEntriesWithPath(entries []*Entry, path string, weight float64) []*Entry {
    path = strings.TrimSuffix(path, string(os.PathSeparator))
    path, err := filepath.Abs(path)
    if err != nil {
        log.Fatal(err)
    }
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

func (e Entry) String() string {
    return fmt.Sprintf("%f\t%s", e.Score, e.Path)
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
    var entries []*Entry
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
            fmt.Errorf("Failed to parse score from line: %v", line)
        }
        entries = append(entries, &entry)
    }
    sortEntriesByScore(entries)
    return entries
}
