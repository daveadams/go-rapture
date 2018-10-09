package shellgen

import (
	"fmt"
)

type BashGenerator struct {
	basegen
}

func (g *BashGenerator) Set(name, value string) {
	g.sb.WriteString(fmt.Sprintf("%s=%s\n", name, shellQuote(value)))
}

func (g *BashGenerator) Export(name, value string) {
	g.sb.WriteString(fmt.Sprintf("export %s=%s\n", name, shellQuote(value)))
}

func (g *BashGenerator) Unset(name string) {
	g.sb.WriteString(fmt.Sprintf("unset %s\n", name))
}

func (g *BashGenerator) Echo(content string) {
	g.sb.WriteString(fmt.Sprintf("echo %s\n", shellQuote(content)))
}

func (g *BashGenerator) Echof(format string, args ...interface{}) {
	g.sb.WriteString(fmt.Sprintf("echo %s\n", shellQuote(fmt.Sprintf(format, args...))))
}

func (g *BashGenerator) ErrEcho(content string) {
	g.sb.WriteString(fmt.Sprintf("echo %s >&2\n", shellQuote(content)))
}

func (g *BashGenerator) ErrEchof(format string, args ...interface{}) {
	g.sb.WriteString(fmt.Sprintf("echo %s >&2\n", shellQuote(fmt.Sprintf(format, args...))))
}
