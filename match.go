package main

import (
	"github.com/lithammer/fuzzysearch/fuzzy"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type matcher func([]*entry, []string) []string

func bestGuess(entries []*entry, args []string) string {
	candidates := getCandidates(entries, args, 1)
	if len(candidates) > 0 {
		return candidates[0]
	}
	return "."
}

var matchConsecutive = func(entries []*entry, args []string) []string {
	nArgs := len(args)
	var matches []string

loop_entries:
	for _, e := range entries {
		parts := strings.Split(e.val, string(os.PathSeparator))
		parts = parts[1:]
		for i, j := len(parts)-1, nArgs-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
			if !strings.Contains(
				strings.ToLower(parts[i]),
				strings.ToLower(args[j]),
			) {
				continue loop_entries
			}
		}
		matches = append(matches, e.val)
	}
	return matches
}

var matchFuzzy = func(entries []*entry, args []string) []string {
	var matches []string
	// Only match the last part
	arg := args[len(args)-1]
	distanceThreshold := len(arg) * 2
	for _, e := range entries {
		_, lastPart := filepath.Split(e.val)
		rank := fuzzy.RankMatch(arg, lastPart)
		if rank == -1 {
			continue
		}
		if rank < distanceThreshold {
			matches = append(matches, e.val)
		}
	}
	return matches
}

var matchAnywhere = func(entries []*entry, args []string) []string {
	var matches []string
	any := ".*"
	regexParts := []string{"(?i)", any, strings.Join(args, any), any}
	regex := strings.Join(regexParts, "")
	pattern, err := regexp.Compile(regex)

	if err != nil {
		return matches
	}

	for _, e := range entries {
		if pattern.Match([]byte(e.val)) {
			matches = append(matches, e.val)
		}
	}

	return matches
}

func getCandidates(entries []*entry, args []string, limit int) []string {
	var candidates []string
	seen := make(map[string]bool)
	matchers := []matcher{matchConsecutive, matchFuzzy, matchAnywhere}
	for _, m := range matchers {
		paths := m(entries, args)
		if len(paths) > 0 {
			for _, p := range paths {
				if seen[p] || !isValidPath(p) {
					continue
				}
				candidates = append(candidates, p)
				seen[p] = true
				if len(candidates) >= limit {
					return candidates
				}
			}
		}
	}
	return candidates
}
