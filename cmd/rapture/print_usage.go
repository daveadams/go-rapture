package main

import (
	"fmt"
	"os"
)

func PrintUsageAndExit(xv int) {
	fmt.Fprintf(os.Stderr, `Usage: rapture <command> [<args> ...]

Commands:

  whoami
    prints the IAM ARN of the currently active identity

  assume <role>
    attempts to assume the role given (either an ARN or an alias)

  resume
    reverts to the prior credentials

  role
    manages role ARN aliases

  account
    manages account aliases

  version
    prints the current version
`)

	os.Exit(xv)
}
