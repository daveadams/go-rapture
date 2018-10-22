package main

import (
	"fmt"
	"os"

	"github.com/daveadams/go-rapture/config"
	"github.com/daveadams/go-rapture/log"
	"github.com/daveadams/go-rapture/session"
	"github.com/daveadams/go-rapture/vaulted"
)

func CommandInit(cmd string, args []string) int {
	log.Tracef("main: CommandInit(cmd='%s', args=%s)", cmd, args)

	RequireWrap()

	if !vaulted.Installed() {
		shgen.ErrEcho("ERROR: can't find 'vaulted' in your path")
		return 1
	}

	vault := config.GetConfig().DefaultVault

	// if VAULTED_ENV is set, use that as the default
	if val, ok := os.LookupEnv("VAULTED_ENV"); ok {
		vault = val
	}

	names, err := vaulted.New().ListVaults()
	if err != nil {
		shgen.ErrEchof("ERROR: Could not load list of vaults: %s", err)
		return 1
	}

	// if a vault name is specified, override the default
	if len(args) > 0 {
		vault = args[0]
	}

	exists := false
	for _, vn := range names {
		if vn == vault {
			exists = true
			break
		}
	}

	if !exists {
		shgen.ErrEchof("ERROR: No such vault '%s'", vault)
		return 1
	}

	// use fmt to print this immediately to stderr
	fmt.Fprintf(os.Stderr, "Initializing vaulted env '%s':\n", vault)
	vars, err := vaulted.LoadVault(vault)
	if err != nil {
		shgen.ErrEchof("ERROR: %s", err)
		return 1
	}

	// clear out existing session vars, we're starting fresh
	os.Unsetenv(session.IDEnvVar)
	os.Unsetenv(session.KeyEnvVar)
	os.Unsetenv(session.SaltEnvVar)

	// remove any rapture assumed-role env vars
	shgen.Unset(session.AssumedRoleArnEnvVar)
	shgen.Unset(session.AssumedRoleAliasEnvVar)

	for varname, value := range vars {
		// set the variable in the shell
		shgen.Export(varname, value)
		// also set the variable in this process so we can cache the base creds
		os.Setenv(varname, value)
	}

	sess, isnew, err := session.CurrentSession()
	if err != nil {
		shgen.ErrEchof("ERROR: %s", err)
		return 1
	}
	if isnew {
		log.Debug("Started a new Rapture session")
	} else {
		log.Debug("Using an existing session, that seems wrong...")
	}

	err = sess.Save(shgen)
	if err != nil {
		shgen.ErrEchof("ERROR: Failed to cache base credentials: %s", err)
	}

	return PrintWhoami()
}
