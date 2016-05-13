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

const separator = "__"
const maxCompleteOptions = 9
const defaultWeight = 20.0

var dataPath string

func init() {
	usr, _ := user.Current()
	dir := filepath.Join(usr.HomeDir, ".local/share/shonenjump")

	dataPath = filepath.Join(dir, "shonenjump.txt")
}

func getNCandidate(args []string, index int, defaultPath string) string {
	entries := loadEntries(dataPath)
	candidates := getCandidates(entries, args, index)
	if len(candidates) == index {
		return candidates[index-1]
	}
	return defaultPath
}

func parseCompleteOption(s string) (string, int, string) {
	needle := ""
	index := 0
	path := ""

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

	return needle, index, path
}

func clearNotExistDirs(entries entryList) entryList {
	isValid := func(e *entry) bool {
		if isValidPath(e.Val) {
			return true
		}
		log.Printf("Directory %s no longer exists", e.Val)
		return false
	}
	return entries.Filter(isValid)
}

func main() {
	pathToAdd := flag.String("add", "", "Add this path")
	complete := flag.Bool("complete", false, "Used for tab completion")
	purge := flag.Bool("purge", false, "Remove non-existent paths from database")
	stat := flag.Bool("stat", false, "Show information about recorded paths")
	flag.Parse()
	if *pathToAdd != "" {
		oldEntries := loadEntries(dataPath)

		path, err := preprocessPath(*pathToAdd)
		if err != nil {
			log.Fatal(err)
		}

		oldEntries.Age()
		newEntries := oldEntries.Update(path, defaultWeight)
		newEntries.Save(dataPath)
	} else if *complete {
		args := flag.Args()
		var arg string
		if len(args) > 0 {
			arg = args[0]
		} else {
			arg = ""
		}
		needle, index, path := parseCompleteOption(arg)
		if path != "" {
			fmt.Println(path)
		} else if index != 0 {
			path = getNCandidate([]string{needle}, index, "")
			if path != "" {
				fmt.Println(path)
			}
		} else {
			entries := loadEntries(dataPath)
			candidates := getCandidates(entries, []string{needle}, maxCompleteOptions)
			for i, path := range candidates {
				parts := []string{needle, strconv.Itoa(i + 1), path}
				fmt.Println(strings.Join(parts, separator))
			}
		}
	} else if *purge {
		entries := loadEntries(dataPath)
		entries = clearNotExistDirs(entries)
		entries.Save(dataPath)
	} else if *stat {
		entries := loadEntries(dataPath)
		for _, e := range entries {
			fmt.Println(e)
		}
	} else if flag.NArg() > 0 {
		args := flag.Args()
		if len(args) == 1 {
			needle, index, path := parseCompleteOption(args[0])
			if path != "" {
				fmt.Println(path)
				return
			} else if index != 0 {
				path = getNCandidate([]string{needle}, index, ".")
				fmt.Println(path)
				return
			} else {
				args = []string{needle}
			}
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
