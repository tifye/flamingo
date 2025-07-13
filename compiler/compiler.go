package compiler

import (
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/tifye/flamingo/ast"
	"github.com/tifye/flamingo/lexer"
	"github.com/tifye/flamingo/parser"
)

func Compile(input string) string {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	root := p.Parse()

	w := &walker{
		outb:    &strings.Builder{},
		renders: make([]string, 0),
	}
	ast.Walk(w, root)
	for _, r := range w.renders {
		fmt.Fprintln(w.outb, r)
	}
	return w.outb.String()
}

type walker struct {
	idCounter atomic.Int32
	parComp   *ast.Component
	parCompId string
	curComp   *ast.Component
	curCompId string
	outb      *strings.Builder
	renders   []string
}

func (w *walker) Visit(n ast.Node) ast.Visitor {
	switch nt := n.(type) {
	case *ast.Component:
		if w.curComp != nil {
			w.parComp = w.curComp
			w.parCompId = w.curCompId
		}

		w.curComp = nt
		w.curCompId = fmt.Sprintf("%s%d", w.curComp.Name.Name, w.idCounter.Add(1))
		w.write("%s := renderer.NewComponent(\"%s\")\n", w.curCompId, nt.Name.Name)

		if w.parComp != nil {
			w.renders = append(w.renders, fmt.Sprintf("renderer.Append(%s, %s)", w.parCompId, w.curCompId))
		} else {
			w.renders = append(w.renders, fmt.Sprintf("renderer.Render(%s)", w.curCompId))
		}

		return w
	case *ast.Attr:
		w.write("%s.SetAttribute(\"%s\", \"%s\")\n", w.curCompId, nt.Name.Name, nt.ValueLit)
		return w
	case *ast.Text:
		w.write("%s.SetAttribute(\"innerText\", \"%s\")\n", w.curCompId, nt.Lit)
		return w
	case *ast.Fragment, *ast.Ident, *ast.Root:
		return w
	}

	return w
}

func (w *walker) write(format string, a ...any) {
	fmt.Fprintf(w.outb, format, a...)
}
