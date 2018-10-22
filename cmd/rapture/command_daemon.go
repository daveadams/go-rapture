package main

import (
	"os"

	"github.com/daveadams/go-rapture/log"
)

// Starts the rapture daemon to renew creds in the background
func CommandDaemon(cmd string, args []string) int {
	log.Tracef("[main] CommandDaemon(cmd='%s', args=%s)", cmd, args)

	if !shgen.Wrapped() {
		shgen.ErrEcho("ERROR: You must run this command using the shell wrapper")
		return 1
	}

	ppid := os.Getppid()

	shgen.Echof("Parent PID is %d", ppid)
	shgen.Echo("Rapture Daemon COMING SOON")

	return 0
}
