package main

import (
	"github.com/daveadams/go-rapture/log"
)

func CommandAssume(cmd string, args []string) int {
	log.Tracef("[main] CommandAssume(cmd='%s', args=%s)", cmd, args)

	RequireWrap()

	if len(args) != 1 {
		shgen.ErrEcho("Usage: rapture assume <role>")
		return 1
	}

	if cc, err := LoadCredentials(args[0]); err != nil {
		shgen.ErrEchof("ERROR: %s", err)
		return 1
	} else {
		shgen.Echof("Assumed role '%s'", cc.RoleArn)
		return 0
	}
}
