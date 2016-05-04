package main

import (
    "fmt"
    "flag"
)

func main() {
    config := getConfig()
    dataPath := config.getDataPath()
    fmt.Println(dataPath)
    pathToAdd := flag.String("add", "", "Add this path")
    flag.Parse()
    if *pathToAdd != "" {
        entries := loadEntries(dataPath)
        weight := 10.0

        entries = updateEntriesWithPath(entries, *pathToAdd, weight)

        saveEntries(entries, dataPath)
    } else if flag.NArg() > 0 {
        args := flag.Args()
        entries := loadEntries(dataPath)
        fmt.Println(bestGuess(entries, args))
    } else {
        flag.Usage()
    }
}
