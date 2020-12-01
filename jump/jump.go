package jump

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	VERSION            = "0.7.16"
	separator          = "__"
	maxCompleteOptions = 9
	defaultWeight      = 20.0
)

func GetNCandidate(args []string, index int, defaultPath string) string {
	entries := LoadEntries(EnsureDataPath())
	candidates := getCandidates(entries, args, index)
	if len(candidates) == index {
		return candidates[index-1]
	}
	return defaultPath
}

func ParseCompleteOption(s string) (needle string, index int, path string) {
	parts := strings.SplitN(s, separator, 3)
	n := len(parts)
	if n == 1 {
		needle = s
	} else {
		needle = parts[0]
		_index, err := strconv.Atoi(parts[1])
		if err != nil {
			index = 0
		} else {
			index = _index
			if n == 3 {
				path = parts[2]
			}
		}
	}
	return
}

func ClearNotExistDirs(entries entryList) (entryList, bool) {
	var changed bool
	var result entryList
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

func preprocessPath(path string) (string, error) {
	// normalize the input
	path = strings.TrimSuffix(path, string(os.PathSeparator))
	return filepath.Abs(path)
}

func AddPath(pathToAdd string) error {
	path, err := preprocessPath(pathToAdd)
	if err != nil {
		return err
	}
	if !isValidPath(path) {
		return nil
	}
	oldEntries := LoadEntries(EnsureDataPath())
	oldEntries.Age()
	newEntries := oldEntries.Update(path, defaultWeight)
	newEntries.Save(EnsureDataPath())
	return nil
}

func ShowAutoCompleteOptions(arg string) {
	needle, index, path := ParseCompleteOption(arg)
	if path != "" {
		fmt.Println(path)
	} else if index != 0 {
		path = GetNCandidate([]string{needle}, index, "")
		if path != "" {
			fmt.Println(path)
		}
	} else {
		entries := LoadEntries(EnsureDataPath())
		candidates := getCandidates(entries, []string{needle}, maxCompleteOptions)
		for i, path := range candidates {
			parts := []string{needle, strconv.Itoa(i + 1), path}
			fmt.Println(strings.Join(parts, separator))
		}
	}
}

func EnsureDataPath() string {
	usr, _ := user.Current()
	dir := filepath.Join(usr.HomeDir, ".local/share/shonenjump")
	if err := os.MkdirAll(dir, 0740); err != nil {
		panic(err)
	}
	return filepath.Join(dir, "shonenjump.txt")
}
