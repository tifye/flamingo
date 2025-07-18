package compiler

import (
	"fmt"
	gtoken "go/token"
	"io"
	"sync/atomic"

	"github.com/tifye/flamingo/ast"
	"github.com/tifye/flamingo/lexer"
	"github.com/tifye/flamingo/parser"
)

func CompileFile(pkg string, file *gtoken.File, input string, output io.Writer) error {
	l := lexer.NewLexer(file, input)
	p := parser.NewParser(l)
	root := p.Parse()

	if len(p.Errors()) > 0 {
		return fmt.Errorf("one or more parser errors: %s", p.Errors())
	}

	imports := [...]string{
		"github.com/tifye/flamingo/render",
		// "github.com/tifye/flamingo/web",
	}

	w := &walker{
		output:  output,
		renders: make([]string, 0),
	}

	fmt.Fprintf(output, "package %s\n\n", pkg)
	fmt.Fprint(output, "import (\n")
	for _, imp := range imports {
		fmt.Fprintf(output, "\t\"%s\"\n", imp)
	}
	fmt.Fprint(output, ")\n\n")
	fmt.Fprintf(output, "func %s(renderer render.Renderer) {\n", file.Name())

	ast.Walk(w, root)
	fmt.Fprint(output, "\n")
	for _, r := range w.renders {
		fmt.Fprintln(w.output, r)
	}
	fmt.Fprint(output, "}")

	return nil
}

type walker struct {
	idCounter atomic.Int32
	parComp   *ast.Component
	parCompId string
	curComp   *ast.Component
	curCompId string
	output    io.Writer
	renders   []string
}

func (w *walker) Visit(n ast.Node) ast.Visitor {
	switch nt := n.(type) {
	case *ast.Component:
		if w.curComp != nil {
			w.parComp = w.curComp
			w.parCompId = w.curCompId
			w.write("\n")
		}

		w.curComp = nt
		w.curCompId = fmt.Sprintf("%s%d", w.curComp.Name.Name, w.idCounter.Add(1))
		w.write("\t%s := renderer.NewComponent(\"%s\")\n", w.curCompId, nt.Name.Name)

		if w.parComp != nil {
			w.renders = append(w.renders, fmt.Sprintf("\trenderer.Append(%s, %s)", w.parCompId, w.curCompId))
		} else {
			w.renders = append(w.renders, fmt.Sprintf("\trenderer.Render(%s)", w.curCompId))
		}

		return w
	case *ast.Attr:
		w.write("\t%s.SetAttribute(\"%s\", \"%s\")\n", w.curCompId, nt.Name.Name, nt.ValueLit)
		return w
	case *ast.Text:
		w.write("\t%s.SetAttribute(\"innerText\", `%s`)\n", w.curCompId, nt.Lit)
		return w
	case *ast.Fragment, *ast.Ident, *ast.Root:
		return w
	}

	return w
}

func (w *walker) write(format string, a ...any) {
	fmt.Fprintf(w.output, format, a...)
}
