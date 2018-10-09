package main

import (
	"sort"

	"github.com/daveadams/go-rapture/config"
	"github.com/daveadams/go-rapture/validation"
)

func CommandAccount(cmd string, args []string) int {
	if len(args) == 0 {
		PrintAccountUsage()
		return 1
	}

	subcmd := args[0]
	subargs := args[1:]

	switch subcmd {
	case "ls", "list":
		return CommandAccountList()

	case "set", "add", "new", "create":
		if len(subargs) != 2 {
			shgen.ErrEchof("Usage: rapture %s %s <alias> <id>", cmd, subcmd)
			return 1
		}
		return CommandAccountAdd(subargs[0], subargs[1])

	case "rm", "remove", "del", "delete":
		if len(subargs) != 1 {
			shgen.ErrEchof("Usage: rapture %s %s <alias>", cmd, subcmd)
			return 1
		}
		return CommandAccountRemove(subargs[0])

	default:
		if len(subargs) != 0 {
			shgen.ErrEcho("Usage: rapture account <alias>")
			return 1
		}
		// treat 'subcmd' as an account alias
		return CommandAccountShow(subcmd)
	}
}

func PrintAccountUsage() {
	shgen.ErrEcho(`Usage: rapture account <command> [<args> ...]

Commands:

  ls
    lists all currently defined accounts

  set <account> <id>
    creates or updates an account named <account> for the value <id>

  rm <account>
    removes the account <account>

  <account>
    prints the value of <id> for alias <account>
`)
}

func CommandAccountList() int {
	acctMap, err := config.LoadAccounts()
	if err != nil {
		shgen.ErrEchof("ERROR: Could not load accounts: %s", err)
		return 1
	}

	// sort account aliases alphabetically
	var aliases []string
	for alias, _ := range acctMap {
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)

	for _, alias := range aliases {
		shgen.Echof("%s %s", acctMap[alias], alias)
	}
	return 0
}

func CommandAccountAdd(a, id string) int {
	acctMap, err := config.LoadAccounts()
	if err != nil {
		shgen.ErrEchof("ERROR: Could not load accounts: %s", err)
		return 1
	}

	if !validation.IsValidAwsAccountId(id) {
		shgen.ErrEchof("ERROR: '%s' is not a valid AWS account ID", id)
		return 1
	}

	acctMap[a] = id

	err = config.WriteAccounts(acctMap)
	if err != nil {
		shgen.ErrEchof("ERROR: Could not write accounts to disk: %s", err)
		return 1
	}

	shgen.Echof("OK: Added alias '%s' for account %s", a, id)

	return 0
}

func CommandAccountRemove(a string) int {
	acctMap, err := config.LoadAccounts()
	if err != nil {
		shgen.ErrEchof("ERROR: Could not load accounts: %s", err)
		return 1
	}

	if _, ok := acctMap[a]; !ok {
		shgen.ErrEchof("ERROR: Account alias '%s' was not found.", a)
		return 1
	}

	delete(acctMap, a)

	err = config.WriteAccounts(acctMap)
	if err != nil {
		shgen.ErrEchof("ERROR: Could not write accounts to disk: %s", err)
		return 1
	}

	shgen.Echof("OK: Removed account alias '%s'", a)

	return 0
}

func CommandAccountShow(a string) int {
	acctMap, err := config.LoadAccounts()
	if err != nil {
		shgen.ErrEchof("ERROR: Could not load accounts: %s", err)
		return 1
	}

	if id, ok := acctMap[a]; ok {
		shgen.Echo(id)
		return 0
	} else {
		shgen.ErrEchof("ERROR: Account alias '%s' was not found.", a)
		return 1
	}
}
