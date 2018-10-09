package main

import (
	"fmt"

	"github.com/daveadams/go-rapture/config"
	"github.com/daveadams/go-rapture/log"
	"github.com/daveadams/go-rapture/session"
)

func CommandAssume(cmd string, args []string) int {
	log.Tracef("[main] CommandAssume(cmd='%s', args=%s)", cmd, args)

	if !shgen.Wrapped() {
		shgen.ErrEcho("ERROR: You must run this command using the shell wrapper")
		return 1
	}

	if len(args) != 1 {
		shgen.ErrEcho("Usage: rapture assume <role>")
		return 1
	}

	roles, err := config.LoadRoles()
	if err != nil {
		shgen.ErrEchof("WARNING: could not load role alias config: %s", err)
	}

	roleName := args[0]
	arn := roleName
	if val, ok := roles[arn]; ok {
		arn = val
	}

	sess, _, err := session.CurrentSession()
	if err != nil {
		shgen.ErrEchof("ERROR: Could not load current Rapture session: %s", err)
		return 1
	}

	cc, err := sess.CredentialsForRole(arn)
	if err != nil {
		if err == session.ErrBaseCredsExpired {
			shgen.ErrEcho("INFO: Base credentials have expired, renewing...")
			// TODO: refactor this madness and eliminate the bashism
			shgen.Pass(fmt.Sprintf("rapture init && rapture assume %s\n", roleName))
			return 0
		}
		shgen.ErrEchof("ERROR: Could not assume role '%s': %s", args[0], err)
		return 1
	}

	cc.Creds.ExportToEnvironment(shgen)
	shgen.Echof("Assumed role '%s'", cc.RoleArn)

	shgen.Export("RAPTURE_ROLE", roleName)
	shgen.Export("RAPTURE_ASSUMED_ROLE_ARN", cc.RoleArn)

	return 0
}
