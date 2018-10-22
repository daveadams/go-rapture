package main

import (
	"github.com/daveadams/go-rapture/log"
)

// Assumes the specified role and runs the specified command
func CommandWrap(cmd string, args []string) int {
	log.Tracef("[main] CommandWrap(cmd='%s', args=%s)", cmd, args)

	RequireWrap()

	if len(args) < 2 {
		shgen.ErrEcho("Usage: rapture wrap <role> <command>...")
		return 1
	}

	// save current env vars
	env := StashEnvironment()

	if _, err := LoadCredentials(args[0]); err != nil {
		shgen.ErrEchof("ERROR: %s", err)
		return 1
	}

	shgen.Run(args[1:])

	// restore previous values
	RestoreEnvironment(env)

	return 0
}
