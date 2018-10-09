package main

import (
	"os"
	"path/filepath"

	"github.com/daveadams/go-rapture/log"
	"github.com/daveadams/go-rapture/shellgen"
)

func CommandSetup(cmd string, args []string) int {
	log.Tracef("[main] CommandEnv(cmd='%s', args=%s)", cmd, args)

	shbin := os.Getenv("SHELL")

	// this command won't be run while wrapped, but we still want to generate safe shell code
	shgen = shellgen.NewGeneratorForShell(shbin)

	binpath, _ := os.Executable()
	shgen.Set("_rapture_bin", binpath)

	switch filepath.Base(shbin) {
	case "bash":
		setupBash()
	case "zsh":
		setupZsh()
	case "fish":
		setupFish()
	default:
		shgen.ErrEchof("ERROR: Rapture does not support inline setup for your shell ('%s')", shbin)
		return 1
	}

	return 0
}

func setupBash() {
	log.Trace("[main] setupBash()")

	shgen.Pass(`rapture() {
    eval "$(
        export _rapture_session_id _rapture_session_key _rapture_session_salt _rapture_wrap=true
        "${_rapture_bin}" "$@"
    )"
}
`)
}

func setupZsh() {
	log.Trace("[main] setupZsh()")
	shgen.Echo("TODO")
}

func setupFish() {
	log.Trace("[main] setupFish()")
	shgen.Echo("TODO")
}
