package main

import (
    "strings"
    "os"
    "path/filepath"
    "github.com/renstrom/fuzzysearch/fuzzy"
    "regexp"
)

func bestGuess(entries []*Entry, args []string) string {
    paths := matchConsecutive(entries, args)
    if len(paths) > 0 {
        return paths[0]
    }

    paths = matchFuzzy(entries, args)
    if len(paths) > 0 {
        return paths[0]
    }

    paths = matchAnywhere(entries, args)
    if len(paths) > 0 {
        return paths[0]
    }

    return "."
}

func matchConsecutive(entries []*Entry, args []string) []string {
    nArgs := len(args)
    var matches []string

    loop_entries:
    for _, e := range entries {
        parts := strings.Split(e.Path, string(os.PathSeparator))
        parts = parts[1:]
        for i, j := len(parts)-1, nArgs-1;
            i >= 0 && j >= 0;
            i, j = i-1, j-1 {
            if !strings.Contains(parts[i], args[j]) {
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
    arg := args[len(args) - 1]
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
