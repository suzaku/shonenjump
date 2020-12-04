package jump

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

const (
	defaultWeight = 20.0
)

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

func preprocessPath(path string) (string, error) {
	// normalize the input
	path = strings.TrimSuffix(path, string(os.PathSeparator))
	return filepath.Abs(path)
}
