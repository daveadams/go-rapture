package main

import (
	"os"
)

type CommandMap map[string]func(string, []string) int

var cmdMap CommandMap = CommandMap{
	"account": CommandAccount,
	"acct":    CommandAccount,
	"alias":   CommandRole,
	"assume":  CommandAssume,
	"check":   CommandCheck,
	"init":    CommandInit,
	"resume":  CommandResume,
	"role":    CommandRole,
	"setup":   CommandSetup,
	"version": CommandVersion,
	"whoami":  CommandWhoami,
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
