package main

import (
	"github.com/daveadams/go-rapture/log"
)

// CommandRefresh implements the refreshing of credentials.
func CommandRefresh(cmd string, args []string) int {
	log.Tracef("[main] CommandRefresh(cmd='%s', args=%s)", cmd, args)

	RequireWrap()

	sess, err := CurrentSession()
	if err != nil {
		shgen.ErrEchof("ERROR: Could not load current Rapture session: %s", err)
		return 1
	}
	sess.BaseCreds.ExportToEnvironment(shgen)

	if cc, err := LoadCredentialsWithForce(sess.AssumedRoleAlias, true); err != nil {
		shgen.ErrEchof("ERROR: %s", err)
		return 1
	} else {
		shgen.Echof("Refreshed role '%s'", cc.RoleArn)
		return 0
	}
}
