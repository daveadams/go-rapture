package main

import (
	"github.com/daveadams/go-rapture/log"
)

func CommandCheck(cmd string, args []string) int {
	log.Tracef("[main] CommandCheck(cmd='%s', args=%s)", cmd, args)

	if !shgen.Wrapped() {
		shgen.ErrEcho("ERROR: Rapture is not correctly wrapped by your shell")
		return 1
	}

	shgen.Echo("OK: Rapture is set up correctly")
	return 0
}
