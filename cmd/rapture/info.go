package main

import (
	"fmt"
	"os"

	"github.com/daveadams/go-rapture/config"
)

var quiet bool

func init() {
	if os.Getenv("RAPTURE_QUIET") != "" || config.GetConfig().Quiet {
		quiet = true
	}
}

// info for unimportant but immediate informational messages
// print these messages immediately to stderr unless quiet is true
func info(s string) {
	if !quiet {
		fmt.Fprint(os.Stderr, s)
	}
}

func infof(s string, args ...interface{}) {
	if !quiet {
		fmt.Fprintf(os.Stderr, s, args...)
	}
}
