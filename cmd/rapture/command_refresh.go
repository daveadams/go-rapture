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

	if sess.AssumedRoleArn == "" {
		shgen.ErrEcho("WARN: No currently assumed role. Nothing to refresh.")
		return 1
	}

	if cc, err := LoadCredentialsWithForce(sess.AssumedRoleAlias, true); err != nil {
		shgen.ErrEchof("ERROR: %s", err)
		return 1
	} else {
		shgen.Echof("Refreshed role '%s'", cc.RoleArn)
		return 0
	}
}
