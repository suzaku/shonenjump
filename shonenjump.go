package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

const separator = "__"

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

func main() {
	config := getConfig()
	dataPath := config.getDataPath()
	pathToAdd := flag.String("add", "", "Add this path")
	complete := flag.Bool("complete", false, "Used for tab completion")
	flag.Parse()
	if *pathToAdd != "" {
		entries := loadEntries(dataPath)
		weight := 10.0

		entries = updateEntriesWithPath(entries, *pathToAdd, weight)

		saveEntries(entries, dataPath)
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
			entries := loadEntries(dataPath)
			candidates := getCandidates(entries, []string{needle}, index)
			if len(candidates) == index {
				fmt.Println(candidates[index-1])
			}
		} else {
			entries := loadEntries(dataPath)
			candidates := getCandidates(entries, []string{needle}, 9)
			for i, path := range candidates {
				parts := []string{needle, strconv.Itoa(i + 1), path}
				fmt.Println(strings.Join(parts, separator))
			}
		}
	} else if flag.NArg() > 0 {
		args := flag.Args()
		entries := loadEntries(dataPath)
		fmt.Println(bestGuess(entries, args))
	} else {
		flag.Usage()
	}
}
