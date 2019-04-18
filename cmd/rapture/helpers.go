package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/daveadams/go-rapture/config"
	"github.com/daveadams/go-rapture/log"
	"github.com/daveadams/go-rapture/session"
)

func CurrentSession() (*session.RaptureSession, error) {
	log.Trace("main: CurrentSession()")

	if !session.CurrentSessionExists() {
		if ec := session.ReadCredentialsFromEnvironment(); ec.Valid() {
			// if AWS credentials exist, use those as the base creds
			info("No current rapture session exists. Using existing credentials from environment.\n")
		} else {
			// otherwise, attempt to init using defaults
			info("No current rapture session exists! Initializing:\n")
			initResult := CommandInit("init", []string{})
			if initResult != 0 {
				return nil, fmt.Errorf("Failed to initialize base credentials")
			}
		}
	}

	sess, _, err := session.CurrentSession()
	if err != nil {
		if err == session.ErrBaseCredsExpired {
			info("Base credentials have expired! Re-initializing:\n")

			// hacky?
			initResult := CommandInit("init", []string{})
			if initResult != 0 {
				return nil, fmt.Errorf("Failed to re-initialize base credentials")
			}

			// retry loading session
			sess, _, err = session.CurrentSession()
		}

		if err != nil {
			shgen.ErrEchof("ERROR: Could not load current Rapture session: %s", err)
			return nil, err
		}
	}

	if err := sess.Save(shgen); err != nil {
		return nil, err
	}

	return sess, nil
}

func LoadCredentials(id string) (*session.CachedCredentials, error) {
	log.Tracef("main: LoadCredentials(id='%s')", id)
	return LoadCredentialsWithForce(id, false)
}

func LoadCredentialsWithForce(id string, forceRefresh bool) (*session.CachedCredentials, error) {
	log.Tracef("main: LoadCredentialsWithForce(id='%s', force_refresh='%t')", id, forceRefresh)
	roles, err := config.LoadRoles()
	if err != nil {
		shgen.ErrEchof("WARNING: could not load role alias config: %s", err)
	}

	roleName := id
	arn := roleName
	if val, ok := roles[arn]; ok {
		arn = val
	}

	sess, err := CurrentSession()
	if err != nil {
		return nil, err
	}

	cc, err := sess.GetCredentialsForRole(arn, forceRefresh)
	if err != nil {
		return nil, fmt.Errorf("Could not assume role '%s': %s", id, err)
	}

	cc.Creds.ExportToEnvironment(shgen)

	shgen.Export(session.AssumedRoleAliasEnvVar, roleName)
	shgen.Export(session.AssumedRoleArnEnvVar, cc.RoleArn)

	return cc, nil
}

// returns a map of current relevant environment vars for later restoration by RestoreEnvironment
func StashEnvironment() map[string]string {
	log.Trace("main: StashEnvironment()")

	vars := []string{
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY",
		"AWS_SECURITY_TOKEN",
		"AWS_SESSION_TOKEN",
		"RAPTURE_ROLE",
		"RAPTURE_ASSUMED_ROLE_ARN",
		"RAPTURE_ROLE_EXPIRATION",
	}

	rv := make(map[string]string, len(vars))

	for _, v := range vars {
		rv[v] = os.Getenv(v)
	}

	return rv
}

// restores a map of environment vars both in the current process, and to shgen
func RestoreEnvironment(env map[string]string) {
	log.Tracef("main: RestoreEnvironment(env=%s)", env)

	for name, value := range env {
		os.Setenv(name, value)
		shgen.Export(name, value)
	}
}

// Shortens a full filepath if possible
func DisplayFilename(fp string) string {
	home := os.Getenv("HOME")
	if home == "" {
		return fp
	}

	rel, err := filepath.Rel(home, fp)
	if err != nil {
		return fp
	}
	if strings.Contains(rel, "../") {
		return fp
	}

	return filepath.Join("~", rel)
}

// Check for wrap status, and fail with an explanatory message if we aren't
func RequireWrap() {
	log.Trace("main: RequireWrap()")

	if shgen.Wrapped() {
		return
	}

	shgen.ErrEcho("ERROR: Rapture is not correctly wrapped by your shell.")

	shbin, ok := os.LookupEnv("SHELL")
	if !ok {
		shgen.ErrEcho("\nNOTICE: Your SHELL environment variable is not set corrrectly.")
		os.Exit(1)
	}

	var f, cmd string
	shname := filepath.Base(shbin)

	switch shname {
	case "bash":
		f = "~/.bash_profile or ~/.bashrc"
		cmd = "eval \"$( command rapture shell-init )\""

	case "zsh":
		f = "~/.zshrc"
		cmd = "eval \"$( command rapture shell-init )\""

	case "fish":
		f = "~/.config/fish/fish.config"
		cmd = "eval ( command rapture shell-init )"

	default:
		shgen.ErrEchof("\nNOTICE: Rapture does not support your shell ('%s').", shname)
		os.Exit(1)
	}

	shgen.ErrEchof(`
To set up Rapture for the '%s' shell, add the following line to your shell
startup configuration file (probably %s):

    %s

Then start a new shell and run 'rapture check' to verify that it worked.`, shname, f, cmd)

	os.Exit(1)
}
