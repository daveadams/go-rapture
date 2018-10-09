package log

import (
	"fmt"
	"os"
)

var debug bool
var trace bool

func init() {
	if os.Getenv("RAPTURE_DEBUG") != "" {
		debug = true
	}

	if os.Getenv("RAPTURE_TRACE") != "" {
		debug = true
		trace = true
	}
}

func Debug(s string) {
	if debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] %s\n", s)
	}
}

func Debugf(s string, args ...interface{}) {
	if debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] %s\n", fmt.Sprintf(s, args...))
	}
}

func Trace(s string) {
	if trace {
		fmt.Fprintf(os.Stderr, "[TRACE] %s\n", s)
	}
}

func Tracef(s string, args ...interface{}) {
	if debug {
		fmt.Fprintf(os.Stderr, "[TRACE] %s\n", fmt.Sprintf(s, args...))
	}
}
