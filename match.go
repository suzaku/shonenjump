package main

import (
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

type matcher func([]*entry, []string) []string

func bestGuess(entries []*entry, args []string) string {
	candidates := getCandidates(entries, args, 1)
	if len(candidates) > 0 {
		return candidates[0]
	}
	return "."
}

var matchExactName = func(entries []*entry, args []string) (matches []string) {
	if len(args) != 1 {
		return
	}
	q := args[0]
	for _, e := range entries {
		if _, name := path.Split(e.val); name == q {
			matches = append(matches, e.val)
		}
	}
	return
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
		diff := calculateDiff(arg, lastPart)
		if diff == -1 {
			continue
		}
		if diff < distanceThreshold {
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
		if pattern.MatchString(e.val) {
			matches = append(matches, e.val)
		}
	}

	return matches
}

func getCandidates(entries []*entry, args []string, limit int) []string {
	candidates := make([]string, 0, limit)
	seen := make(map[string]bool)
	matchers := []matcher{matchExactName, matchConsecutive, matchFuzzy, matchAnywhere}
	for _, m := range matchers {
		paths := m(entries, args)
		if len(paths) == 0 {
			continue
		}
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
	return candidates
}

func calculateDiff(source, target string) int {
	lenDiff := len(target) - len(source)
	if lenDiff < 0 {
		return -1
	}
	if lenDiff == 0 && source == target {
		return 0
	}

	var runeDiff int

	for _, r := range source {
		i := strings.IndexRune(target, r)
		if i == -1 {
			return -1
		}
		runeDiff += utf8.RuneCountInString(target[:i])
		target = target[i+utf8.RuneLen(r):]
	}

	// Count up remaining char
	runeDiff += utf8.RuneCountInString(target)
	return runeDiff
}
