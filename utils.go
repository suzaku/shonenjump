package main

import (
	"os"
)

var isValidPath = func(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	return true
}
