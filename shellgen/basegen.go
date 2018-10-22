package shellgen

import (
	"fmt"
	"strings"
)

type basegen struct {
	sb strings.Builder
}

func (g *basegen) Wrapped() bool {
	return true
}

// simply pass-through the given string to the shell
func (g *basegen) Pass(s string) {
	g.sb.WriteString(s)
}

func (g *basegen) Passf(format string, args ...interface{}) {
	g.sb.WriteString(fmt.Sprintf(format, args...))
}

func (g *basegen) Print() {
	fmt.Print(g.sb.String())
}

func (g *basegen) Run(argv []string) {
	actual := make([]string, len(argv))
	for i, arg := range argv {
		actual[i] = shellQuote(arg)
	}

	g.Pass(strings.Join(actual, " "))
	g.Pass("\n")
}
