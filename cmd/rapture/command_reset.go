package main

import (
	"github.com/daveadams/go-rapture/log"
)

// fully resets the rapture environment
func CommandReset(cmd string, args []string) int {
	log.Tracef("[main] CommandReset(cmd='%s', args=%s)", cmd, args)

	RequireWrap()

	shgen.Echo("Rapture Reset COMING SOON")

	return 0
}
