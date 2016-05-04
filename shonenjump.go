package main

import (
    "fmt"
    "flag"
    "strings"
    "strconv"
)

const separator = "__"

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
        needle := arg
        entries := loadEntries(dataPath)
        candidates := getCandidates(entries, []string{arg}, 9)
        for i, path := range candidates {
            parts := []string{needle, strconv.Itoa(i + 1), path}
            fmt.Println(strings.Join(parts, separator))
        }
    } else if flag.NArg() > 0 {
        args := flag.Args()
        entries := loadEntries(dataPath)
        fmt.Println(bestGuess(entries, args))
    } else {
        flag.Usage()
    }
}
