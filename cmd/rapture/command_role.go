package main

import (
	"sort"

	"github.com/daveadams/go-rapture/config"
	"github.com/daveadams/go-rapture/validation"
)

func CommandRole(cmd string, args []string) int {
	if len(args) == 0 {
		PrintRoleUsage()
		return 1
	}

	subcmd := args[0]
	subargs := args[1:]

	switch subcmd {
	case "ls", "list":
		return CommandRoleList()

	case "set", "add", "new", "create":
		if len(subargs) != 2 {
			shgen.ErrEchof("Usage: rapture %s %s <alias> <arn>", cmd, subcmd)
			return 1
		}
		return CommandRoleAdd(subargs[0], subargs[1])

	case "rm", "remove", "del", "delete":
		if len(subargs) != 1 {
			shgen.ErrEchof("Usage: rapture %s %s <alias>", cmd, subcmd)
			return 1
		}
		return CommandRoleRemove(subargs[0])

	default:
		if len(subargs) != 0 {
			shgen.ErrEchof("Usage: rapture %s <alias>", cmd)
			return 1
		}
		// treat 'subcmd' as an role alias
		return CommandRoleShow(subcmd)
	}
}

func PrintRoleUsage() {
	shgen.ErrEcho(`Usage: rapture role <command> [<args> ...]

Commands:

  ls
    lists all currently defined roles

  set <alias> <arn>
    creates or updates a role named <alias> for the value <arn>

  rm <alias>
    removes the role <alias>

  <alias>
    prints the value of <arn> for the role with alias <alias>
`)
}

func CommandRoleList() int {
	roleMap, err := config.LoadRoles()
	if err != nil {
		shgen.ErrEchof("ERROR: Could not load roles: %s", err)
		return 1
	}

	// sort role aliases alphabetically
	var aliases []string
	for alias, _ := range roleMap {
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)

	for _, alias := range aliases {
		shgen.Echof("%s %s", roleMap[alias], alias)
	}
	return 0
}

func CommandRoleAdd(a, arn string) int {
	roleMap, err := config.LoadRoles()
	if err != nil {
		shgen.ErrEchof("ERROR: Could not load roles: %s", err)
		return 1
	}

	if !validation.IsValidIamRoleArn(arn) {
		shgen.ErrEchof("ERROR: '%s' is not a valid IAM role ARN", arn)
		return 1
	}

	roleMap[a] = arn

	err = config.WriteRoles(roleMap)
	if err != nil {
		shgen.ErrEchof("ERROR: Could not write roles to disk: %s", err)
		return 1
	}

	shgen.Echof("OK: Added alias '%s' for role %s", a, arn)

	return 0
}

func CommandRoleRemove(a string) int {
	roleMap, err := config.LoadRoles()
	if err != nil {
		shgen.ErrEchof("ERROR: Could not load roles: %s", err)
		return 1
	}

	if _, ok := roleMap[a]; !ok {
		shgen.ErrEchof("ERROR: Role alias '%s' was not found.", a)
		return 1
	}

	delete(roleMap, a)

	err = config.WriteRoles(roleMap)
	if err != nil {
		shgen.ErrEchof("ERROR: Could not write roles to disk: %s", err)
		return 1
	}

	shgen.Echof("OK: Removed role alias '%s'", a)

	return 0
}

func CommandRoleShow(a string) int {
	roleMap, err := config.LoadRoles()
	if err != nil {
		shgen.ErrEchof("ERROR: Could not load roles: %s", err)
		return 1
	}

	if arn, ok := roleMap[a]; ok {
		shgen.Echo(arn)
		return 0
	} else {
		shgen.ErrEchof("ERROR: Role alias '%s' was not found.", a)
		return 1
	}
}
