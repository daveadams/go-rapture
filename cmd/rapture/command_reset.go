package main

import (
	"github.com/daveadams/go-rapture/log"
	"github.com/daveadams/go-rapture/session"
)

// fully resets the rapture environment
func CommandReset(cmd string, args []string) int {
	log.Tracef("[main] CommandReset(cmd='%s', args=%s)", cmd, args)

	RequireWrap()

	// clear rapture session env vars
	shgen.Unset(session.IDEnvVar)
	shgen.Unset(session.KeyEnvVar)
	shgen.Unset(session.SaltEnvVar)

	// clear exported env vars related to the assumed role
	shgen.Unset(session.AssumedRoleArnEnvVar)
	shgen.Unset(session.AssumedRoleAliasEnvVar)
	shgen.Unset(session.AssumedRoleExpirationEnvVar)

	// clear Vaulted env vars
	shgen.Unset("VAULTED_ENV")
	shgen.Unset("VAULTED_ENV_EXPIRATION")

	// clear AWS credential env vars (controversial?)
	shgen.Unset("AWS_ACCESS_KEY_ID")
	shgen.Unset("AWS_SECRET_ACCESS_KEY")
	shgen.Unset("AWS_SESSION_TOKEN")
	shgen.Unset("AWS_SECURITY_TOKEN")

	shgen.Echo("All Rapture-related environment variables have been unset.")

	return 0
}
