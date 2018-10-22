package shellgen

import (
	"fmt"
)

type FishGenerator struct {
	basegen
}

func (g *FishGenerator) Set(name, value string) {
	g.sb.WriteString(fmt.Sprintf("set -g %s %s;\n", name, shellQuote(value)))
}

func (g *FishGenerator) Export(name, value string) {
	g.sb.WriteString(fmt.Sprintf("set -gx %s %s;\n", name, shellQuote(value)))
}

func (g *FishGenerator) Unset(name string) {
	g.sb.WriteString(fmt.Sprintf("set -e %s;\n", name))
}

func (g *FishGenerator) Echo(content string) {
	g.sb.WriteString(fmt.Sprintf("echo %s;\n", shellQuote(content)))
}

func (g *FishGenerator) Echof(format string, args ...interface{}) {
	g.sb.WriteString(fmt.Sprintf("echo %s;\n", shellQuote(fmt.Sprintf(format, args...))))
}

func (g *FishGenerator) ErrEcho(content string) {
	g.sb.WriteString(fmt.Sprintf("echo %s >&2;\n", shellQuote(content)))
}

func (g *FishGenerator) ErrEchof(format string, args ...interface{}) {
	g.sb.WriteString(fmt.Sprintf("echo %s >&2;\n", shellQuote(fmt.Sprintf(format, args...))))
}
