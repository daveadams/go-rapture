package main

import (
	"fmt"
	"os"
)

func PrintUsageAndExit(xv int) {
	fmt.Fprintf(os.Stderr, `Usage: rapture <command> [<args> ...]

Commands:

  init <vault-name>
    Initializes a rapture session in this shell by loading credentials from
    the Vaulted vault <vault-name>.

  whoami
    Prints the IAM ARN of the currently active identity.

  assume <role>
    Attempts to assume the role given (either an ARN or an alias) from the
    current session's base credentials.

  wrap <role> <command ...>
    Assumes the role <role> and executes <command>. Then restores the existing
    credentials.

  resume
    Reverts from an assumed role to the base credentials.

  check
    Checks the status of Rapture's shell integration.

  config
    Prints the current configuration.

  version
    Prints the current version.
`)

	os.Exit(xv)
}
