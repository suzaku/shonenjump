package main

import (
	"github.com/renstrom/fuzzysearch/fuzzy"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var isValidPath = func(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	return true
}

type matcher func([]*Entry, []string) []string

func bestGuess(entries []*Entry, args []string) string {
	candidates := getCandidates(entries, args, 1)
	if len(candidates) > 0 {
		return candidates[0]
	} else {
		return "."
	}
}

func matchConsecutive(entries []*Entry, args []string) []string {
	nArgs := len(args)
	var matches []string

loop_entries:
	for _, e := range entries {
		parts := strings.Split(e.Path, string(os.PathSeparator))
		parts = parts[1:]
		for i, j := len(parts)-1, nArgs-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
			if !strings.Contains(
				strings.ToLower(parts[i]),
				strings.ToLower(args[j]),
			) {
				continue loop_entries
			}
		}
		matches = append(matches, e.Path)
	}
	return matches
}

func matchFuzzy(entries []*Entry, args []string) []string {
	var matches []string
	// Only match the last part
	arg := args[len(args)-1]
	distanceThreshold := len(arg) * 2
	for _, e := range entries {
		_, lastPart := filepath.Split(e.Path)
		rank := fuzzy.RankMatch(arg, lastPart)
		if rank == -1 {
			continue
		}
		if rank < distanceThreshold {
			matches = append(matches, e.Path)
		}
	}
	return matches
}

func matchAnywhere(entries []*Entry, args []string) []string {
	var matches []string
	any := ".*"
	regexParts := []string{"(?i)", any, strings.Join(args, any), any}
	regex := strings.Join(regexParts, "")
	pattern, err := regexp.Compile(regex)

	if err != nil {
		return matches
	}

	for _, e := range entries {
		if pattern.Match([]byte(e.Path)) {
			matches = append(matches, e.Path)
		}
	}

	return matches
}

func getCandidates(entries []*Entry, args []string, limit int) []string {
	var candidates []string
	matchers := []matcher{matchConsecutive, matchFuzzy, matchAnywhere}
	for _, m := range matchers {
		paths := m(entries, args)
		if len(paths) > 0 {
			for _, p := range paths {
				if !isValidPath(p) {
					continue
				}
				candidates = append(candidates, p)
				if len(candidates) >= limit {
					return candidates
				}
			}
		}
	}
	return candidates
}
