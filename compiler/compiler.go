package compiler

import (
	"fmt"
	gtoken "go/token"
	"io"
	"io/fs"
	"os"
	"strings"
	"sync/atomic"

	"github.com/tifye/flamingo/assert"
	"github.com/tifye/flamingo/ast"
	"github.com/tifye/flamingo/lexer"
	"github.com/tifye/flamingo/parser"
)

func CompileDir(pkg string, fset *gtoken.FileSet, path string, output func(fs.FileInfo) (io.WriteCloser, error)) error {
	assert.AssertNotNil(output)

	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("read dir: %s", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".flamingo") {
			continue
		}

		finfo, err := entry.Info()
		if err != nil {
			return fmt.Errorf("read file info for %s: %s", entry.Name(), err)
		}

		w, err := output(finfo)
		if err != nil {
			return err
		}

		inputb, err := os.ReadFile(entry.Name())
		if err != nil {
			return err
		}

		name := strings.TrimSuffix(entry.Name(), ".flamingo")
		file := fset.AddFile(name, fset.Base(), len(inputb))
		l := lexer.NewLexer(file, string(inputb))
		p := parser.NewParser(l)
		root := p.Parse()
		if len(p.Errors()) > 0 {
			return fmt.Errorf("one or more parser errors: %s", p.Errors())
		}

		if err := CompileFile(pkg, name, root, string(inputb), w); err != nil {
			_ = w.Close()
			return err
		}
		_ = w.Close()
	}

	return nil
}

func CompileFile(pkg string, file string, root *ast.Root, input string, output io.Writer) error {
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
	fmt.Fprintf(output, "func %s(renderer render.Renderer) {\n", file)

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
