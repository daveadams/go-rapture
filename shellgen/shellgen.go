package shellgen

import (
	"os"
	"path/filepath"

	shellquote "github.com/kballard/go-shellquote"
)

type Generator interface {
	Wrapped() bool
	Set(string, string)
	Export(string, string)
	Unset(string)
	Echo(string)
	Echof(string, ...interface{})
	ErrEcho(string)
	ErrEchof(string, ...interface{})
	Pass(string)
	Passf(string, ...interface{})
	Run([]string)
	Print()
}

func shellQuote(s string) string {
	return shellquote.Join(s)
}

func NewGenerator() Generator {
	if os.Getenv("_rapture_wrap") == "true" {
		return NewGeneratorForShell(os.Getenv("SHELL"))
	} else {
		// don't return executable shell code, just print
		return &TerminalGenerator{}
	}
}

func NewGeneratorForShell(shell string) Generator {
	if shell == "" {
		shell = "/bin/bash"
	}

	switch filepath.Base(shell) {
	case "bash", "zsh":
		return &BashGenerator{}
	case "fish":
		return &FishGenerator{}
	default:
		return nil
	}
}
