package main

import (
	"os"

	"github.com/daveadams/go-rapture/config"
	"github.com/daveadams/go-rapture/log"
)

func CommandConfig(cmd string, args []string) int {
	log.Tracef("[main] CommandConfig(cmd='%s', args=%s)", cmd, args)

	defaults := config.DefaultConfig()
	conf := config.GetConfig()
	raw, exists, err := config.RawConfig()

	fn := DisplayFilename(config.ConfigFilename())
	if exists {
		shgen.Echof("Config file: %s", fn)
	} else {
		shgen.Echof("INFO: Config file not found at %s, using defaults", fn)
	}
	shgen.Echo("")

	if err != nil {
		shgen.Echof("WARN: An error occurred reading the config file:\n  %s\n", err)
		shgen.Echo("INFO: Defaults will be used instead")
	}

	awsDefaultRegionEnvValue := "<unset>"
	if v := os.Getenv("AWS_DEFAULT_REGION"); v != "" {
		awsDefaultRegionEnvValue = v
	}

	awsRegionEnvValue := "<unset>"
	if v := os.Getenv("AWS_REGION"); v != "" {
		awsRegionEnvValue = v
	}

	shgen.Echo("AWS Region to use for API calls:")
	shgen.Echof("  AWS_DEFAULT_REGION: %s", awsDefaultRegionEnvValue)
	shgen.Echof("          AWS_REGION: %s", awsRegionEnvValue)
	shgen.Echof("             default: %s", config.DefaultAwsRegion)
	shgen.Echof("               value: %s", conf.Region())
	shgen.Echo("")

	cid := raw.Identifier
	if cid == "" {
		cid = "<unset>"
	}

	shgen.Echo("STS Assumed Role Session Identifier:")
	shgen.Echo("      key: identifier")
	shgen.Echof("  default: %s (from USER env var)", defaults.Identifier)
	shgen.Echof("     file: %s", cid)
	shgen.Echof("    value: %s", conf.Identifier)
	shgen.Echo("")

	shgen.Echo("STS Session Duration:")
	shgen.Echo("      key: session_duration")
	shgen.Echof("  default: %d", defaults.SessionDuration)
	if raw.SessionDuration == 0 {
		shgen.Echo("     file: <unset>")
	} else {
		shgen.Echof("     file: %d", raw.SessionDuration)
	}
	shgen.Echof("    value: %d", conf.SessionDuration)

	if conf.SessionDuration < 900 {
		shgen.Echo("WARNING: session_duration must be at least 900")
	} else if raw.SessionDuration > 3600 {
		shgen.Echo("WARNING: assuming roles from temporary credentials must have a session_duration of 3600 or lower")
	}
	shgen.Echo("")

	cvault := raw.DefaultVault
	if cvault == "" {
		cvault = "<unset>"
	}

	shgen.Echo("Default Vaulted Vault Name:")
	shgen.Echo("      key: default_vault")
	shgen.Echof("  default: %s", defaults.DefaultVault)
	shgen.Echof("     file: %s", cvault)
	shgen.Echof("    value: %s", conf.DefaultVault)

	shgen.Echof("\nTo change any of these settings, edit %s", fn)

	return 0
}
