package main

import (
	"os"
	"path/filepath"

	"github.com/daveadams/go-rapture/log"
	"github.com/daveadams/go-rapture/shellgen"
)

func CommandShellInit(cmd string, args []string) int {
	log.Tracef("[main] CommandShellInit(cmd='%s', args=%s)", cmd, args)

	shbin := os.Getenv("SHELL")

	shgen = shellgen.NewGeneratorForShell(shbin)

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

	binpath, _ := os.Executable()

	shgen.Passf(`rapture() {
    eval "$(
        export _rapture_session_id _rapture_session_key _rapture_session_salt _rapture_wrap=true
        "%s" "$@"
    )"
}
`, binpath)
}

func setupZsh() {
	log.Trace("[main] setupZsh()")
	setupBash()
}

func setupFish() {
	log.Trace("[main] setupFish()")

	binpath, _ := os.Executable()
	shgen.Passf(`function rapture;
    set -l IFS;
    set -lx _rapture_session_id $_rapture_session_id;
    set -lx _rapture_session_key $_rapture_session_key;
    set -lx _rapture_session_salt $_rapture_session_salt;
    set -lx _rapture_wrap true;
    %s $argv |source;
end
`, binpath)
}
