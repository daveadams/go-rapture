package shellgen

import (
	"fmt"
	"os"
)

// if rapture is run without the shell wrapper, we want to print something sensible
type TerminalGenerator struct{}

func (g *TerminalGenerator) Wrapped() bool {
	return false
}

// setting vars doesn't work for this generator
func (g *TerminalGenerator) Set(name, value string)    {}
func (g *TerminalGenerator) Export(name, value string) {}
func (g *TerminalGenerator) Unset(name string)         {}

func (g *TerminalGenerator) Echo(content string) {
	fmt.Printf("%s\n", content)
}

func (g *TerminalGenerator) Echof(format string, args ...interface{}) {
	fmt.Printf("%s\n", fmt.Sprintf(format, args...))
}

func (g *TerminalGenerator) ErrEcho(content string) {
	fmt.Fprintf(os.Stderr, "%s\n", content)
}

func (g *TerminalGenerator) ErrEchof(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s\n", fmt.Sprintf(format, args...))
}

func (g *TerminalGenerator) Pass(s string) {
	fmt.Print(s)
}

func (g *TerminalGenerator) Passf(format string, args ...interface{}) {
	fmt.Printf("%s\n", fmt.Sprintf(format, args...))
}

// nothing to do here
func (g *TerminalGenerator) Print()            {}
func (g *TerminalGenerator) Run(argv []string) {}
