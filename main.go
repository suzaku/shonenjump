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

	"github.com/suzaku/shonenjump/jump"
)

const (
	version   = "0.7.19"
	separator = "__"
)

func ensureDataPath() string {
	usr, _ := user.Current()
	dir := filepath.Join(usr.HomeDir, ".local/share/shonenjump")
	if err := os.MkdirAll(dir, 0740); err != nil {
		panic(err)
	}
	return filepath.Join(dir, "shonenjump.txt")
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

func main() {
	pathToAdd := flag.String("add", "", "Add this path")
	complete := flag.Bool("complete", false, "Used for tab completion")
	purge := flag.Bool("purge", false, "Remove non-existent paths from database")
	stat := flag.Bool("stat", false, "Show information about recorded paths")
	ver := flag.Bool("version", false, "Show version of shonenjump")
	flag.Parse()
	dataPath := ensureDataPath()
	store := jump.NewStore(dataPath)
	if *pathToAdd != "" {
		if err := store.AddPath(*pathToAdd); err != nil {
			log.Fatal(err)
		}
	} else if *complete {
		args := flag.Args()
		var arg string
		if len(args) > 0 {
			arg = args[0]
		} else {
			arg = ""
		}
		showAutoCompleteOptions(store, arg)
	} else if *purge {
		if err := store.Cleanup(); err != nil {
			log.Fatal(err)
		}
	} else if *stat {
		entries, err := store.ReadEntries()
		if err != nil {
			log.Fatal(err)
		}
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
				path, err := store.GetNthCandidate([]string{needle}, index, ".")
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(path)
				return
			}
			args = []string{needle}
		}
		entries, err := store.ReadEntries()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(jump.BestGuess(entries, args))
	} else {
		flag.Usage()
	}
}

func showAutoCompleteOptions(store jump.Store, arg string) {
	needle, index, path := parseCompleteOption(arg)
	if path != "" {
		fmt.Println(path)
	} else if index != 0 {
		path, err := store.GetNthCandidate([]string{needle}, index, "")
		if err != nil {
			log.Fatal(err)
		}
		if path != "" {
			fmt.Println(path)
		}
	} else {
		entries, err := store.ReadEntries()
		if err != nil {
			log.Fatal(err)
		}
		candidates := jump.GetCandidates(entries, []string{needle}, jump.MaxCompleteOptions)
		for i, path := range candidates {
			parts := []string{needle, strconv.Itoa(i + 1), path}
			fmt.Println(strings.Join(parts, separator))
		}
	}
}
