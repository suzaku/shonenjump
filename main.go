package main

import (
	"github.com/suzaku/shonenjump/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
