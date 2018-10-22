package main

import (
	"github.com/daveadams/go-rapture/log"
)

func CommandCheck(cmd string, args []string) int {
	log.Tracef("main: CommandCheck(cmd='%s', args=%s)", cmd, args)

	RequireWrap()

	shgen.Echo("OK: Rapture is set up correctly")
	return 0
}
