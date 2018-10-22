package main

import (
	"github.com/daveadams/go-rapture/log"
	"github.com/daveadams/go-rapture/session"
)

func CommandResume(cmd string, args []string) int {
	log.Tracef("[main] CommandResume(cmd='%s', args=%s)", cmd, args)

	RequireWrap()

	sess, _, err := session.CurrentSession()
	if err != nil {
		shgen.ErrEchof("ERROR: Could not load current Rapture session: %s", err)
		return 1
	}

	sess.BaseCreds.ExportToEnvironment(shgen)

	shgen.Echof("Resumed base credentials")

	shgen.Unset(session.AssumedRoleAliasEnvVar)
	shgen.Unset(session.AssumedRoleArnEnvVar)

	return 0
}
