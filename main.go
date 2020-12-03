package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	version            = "0.7.18"
	separator          = "__"
	maxCompleteOptions = 9
	defaultWeight      = 20.0
)

func ensureDataPath() string {
	usr, _ := user.Current()
	dir := filepath.Join(usr.HomeDir, ".local/share/shonenjump")
	if err := os.MkdirAll(dir, 0740); err != nil {
		panic(err)
	}
	return filepath.Join(dir, "shonenjump.txt")
}

func getNCandidate(args []string, index int, defaultPath string) string {
	entries := loadEntries(ensureDataPath())
	candidates := getCandidates(entries, args, index)
	if len(candidates) == index {
		return candidates[index-1]
	}
	return defaultPath
}

func parseCompleteOption(s string) (needle string, index int, path string) {
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

func clearNotExistDirs(entries entryList) (entryList, bool) {
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

func main() {
	pathToAdd := flag.String("add", "", "Add this path")
	complete := flag.Bool("complete", false, "Used for tab completion")
	purge := flag.Bool("purge", false, "Remove non-existent paths from database")
	stat := flag.Bool("stat", false, "Show information about recorded paths")
	ver := flag.Bool("version", false, "Show version of shonenjump")
	flag.Parse()
	dataPath := ensureDataPath()
	if *pathToAdd != "" {
		addPath(*pathToAdd)
	} else if *complete {
		args := flag.Args()
		var arg string
		if len(args) > 0 {
			arg = args[0]
		} else {
			arg = ""
		}
		showAutoCompleteOptions(arg)
	} else if *purge {
		entries := loadEntries(dataPath)
		entries, changed := clearNotExistDirs(entries)
		if changed {
			entries.Save(dataPath)
		}
	} else if *stat {
		entries := loadEntries(dataPath)
		for _, e := range entries {
			fmt.Println(e)
		}
	} else if *ver {
		fmt.Println(version)
	} else if flag.NArg() > 0 {
		args := flag.Args()
		if len(args) == 1 {
			needle, index, path := parseCompleteOption(args[0])
			if path != "" {
				fmt.Println(path)
				return
			}
			if index != 0 {
				path = getNCandidate([]string{needle}, index, ".")
				fmt.Println(path)
				return
			}
			args = []string{needle}
		}
		entries := loadEntries(dataPath)
		fmt.Println(bestGuess(entries, args))
	} else {
		flag.Usage()
	}
}

func preprocessPath(path string) (string, error) {
	// normalize the input
	path = strings.TrimSuffix(path, string(os.PathSeparator))
	return filepath.Abs(path)
}

func addPath(pathToAdd string) {
	path, err := preprocessPath(pathToAdd)
	if err != nil {
		log.Fatal(err)
	}
	if !isValidPath(path) {
		return
	}
	oldEntries := loadEntries(ensureDataPath())
	oldEntries.Age()
	newEntries := oldEntries.Update(path, defaultWeight)
	newEntries.Save(ensureDataPath())
}

func showAutoCompleteOptions(arg string) {
	needle, index, path := parseCompleteOption(arg)
	if path != "" {
		fmt.Println(path)
	} else if index != 0 {
		path = getNCandidate([]string{needle}, index, "")
		if path != "" {
			fmt.Println(path)
		}
	} else {
		entries := loadEntries(ensureDataPath())
		candidates := getCandidates(entries, []string{needle}, maxCompleteOptions)
		for i, path := range candidates {
			parts := []string{needle, strconv.Itoa(i + 1), path}
			fmt.Println(strings.Join(parts, separator))
		}
	}
}
