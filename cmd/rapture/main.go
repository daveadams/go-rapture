package main

import (
	"os"
)

type CommandMap map[string]func(string, []string) int

var cmdMap CommandMap = CommandMap{
	"account":    CommandAccount,
	"acct":       CommandAccount,
	"alias":      CommandRole,
	"assume":     CommandAssume,
	"check":      CommandCheck,
	"clean":      CommandClean,
	"config":     CommandConfig,
	"daemon":     CommandDaemon,
	"init":       CommandInit,
	"refresh":    CommandRefresh,
	"reset":      CommandReset,
	"resume":     CommandResume,
	"role":       CommandRole,
	"shell-init": CommandShellInit,
	"version":    CommandVersion,
	"whoami":     CommandWhoami,
	"wrap":       CommandWrap,
}

func main() {
	if len(os.Args) == 1 {
		PrintUsageAndExit(1)
	}

	if cmdFunc, ok := cmdMap[os.Args[1]]; ok {
		rv := cmdFunc(os.Args[1], os.Args[2:])
		shgen.Print()
		os.Exit(rv)
	} else {
		PrintUsageAndExit(1)
	}
}
