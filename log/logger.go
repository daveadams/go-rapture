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

func DebugEnabled() bool {
	return debug
}

func TraceEnabled() bool {
	return trace
}

func TraceEnvironment() {
	if trace {
		for _, v := range os.Environ() {
			Trace(v)
		}
	}
}

func DebugEnvironment() {
	if debug {
		for _, v := range os.Environ() {
			Debug(v)
		}
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
	if trace {
		fmt.Fprintf(os.Stderr, "[TRACE] %s\n", fmt.Sprintf(s, args...))
	}
}
