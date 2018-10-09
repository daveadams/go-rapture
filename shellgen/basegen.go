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

func (g *basegen) Print() {
	fmt.Print(g.sb.String())
}
