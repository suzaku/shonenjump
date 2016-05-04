package main

import (
    "strings"
    "os"
)

func bestGuess(entries []*Entry, args []string) string {
    paths := matchConsecutive(entries, args)
    if len(paths) > 0 {
        return paths[0]
    } else {
        return ""
    }
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
