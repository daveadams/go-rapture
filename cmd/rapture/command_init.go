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

	var vars map[string]string
	var err error
	switch config.GetConfig().InitMethod {
	case "vaulted":
		vars, err = initWithVaulted(args)
	default:
		shgen.ErrEchof("ERROR: init_method '%s' is unknown. Cannot continue.")
	}

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
	shgen.Unset(session.AssumedRoleExpirationEnvVar)

	for varname, value := range vars {
		// set the variable in the shell
		shgen.Export(varname, value)
		// also set the variable in this process so we can cache the base creds
		os.Setenv(varname, value)
	}

	sess, _, err := session.CurrentSession()
	if err != nil {
		shgen.ErrEchof("ERROR: %s", err)
		return 1
	}

	err = sess.Save(shgen)
	if err != nil {
		shgen.ErrEchof("ERROR: Failed to save session: %s", err)
	}

	return PrintWhoami()
}

func initWithVaulted(args []string) (map[string]string, error) {
	if !vaulted.Installed() {
		return nil, fmt.Errorf("can't find 'vaulted' in your path")
	}

	vault := config.GetConfig().DefaultVault

	// if VAULTED_ENV is set, use that as the default
	if val, ok := os.LookupEnv("VAULTED_ENV"); ok {
		vault = val
	}

	names, err := vaulted.New().ListVaults()
	if err != nil {
		return nil, fmt.Errorf("Could not load list of vaults: %s", err)
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
		if len(args) == 0 {
			if len(names) > 0 {
				fn := DisplayFilename(config.ConfigFilename())
				return nil, fmt.Errorf("No default vault is defined. Consider setting 'default_vault' in %s.", fn)
			} else {
				return nil, fmt.Errorf("No Vaulted vaults are available. Please use Vaulted to store your base credentials.")
			}
		} else {
			return nil, fmt.Errorf("No such vault '%s'", vault)
		}
	}

	infof("Initializing vaulted env '%s':\n", vault)
	return vaulted.LoadVault(vault)
}
